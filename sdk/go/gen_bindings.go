package dss

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BindingFormat specifies the output format for generated bindings.
type BindingFormat string

const (
	FormatCSS        BindingFormat = "css"
	FormatTypeScript BindingFormat = "typescript"
	FormatSCSS       BindingFormat = "scss"
)

// BindingOptions configures binding generation.
type BindingOptions struct {
	// Format specifies the output format.
	Format BindingFormat

	// SpecDir is the directory containing the design system spec.
	// Used to resolve local component references.
	SpecDir string

	// DefaultStrategy is the fallback strategy when not specified in bindings.
	DefaultStrategy string
}

// GeneratedBinding represents the output for a single component's theme bindings.
type GeneratedBinding struct {
	// Component is the component ID.
	Component string

	// CSS is the generated CSS/SCSS content.
	CSS string

	// Warnings are non-fatal issues encountered during generation.
	Warnings []string
}

// GenerateBindings generates theme bindings for all configured components.
func GenerateBindings(ds *DesignSystem, opts BindingOptions) ([]GeneratedBinding, error) {
	if opts.Format == "" {
		opts.Format = FormatCSS
	}
	if opts.DefaultStrategy == "" {
		opts.DefaultStrategy = "explicit"
	}

	var results []GeneratedBinding

	for _, binding := range ds.ThemeBindings {
		result, err := generateBinding(ds, &binding, opts)
		if err != nil {
			return nil, fmt.Errorf("binding for %s: %w", binding.Component, err)
		}
		results = append(results, *result)
	}

	return results, nil
}

// generateBinding generates bindings for a single component.
func generateBinding(ds *DesignSystem, binding *ThemeBindings, opts BindingOptions) (*GeneratedBinding, error) {
	result := &GeneratedBinding{
		Component: binding.Component,
	}

	// Find the component's theming contract
	contract, err := findContract(ds, binding, opts.SpecDir)
	if err != nil {
		return nil, err
	}

	if contract == nil {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("component '%s' has no theming contract", binding.Component))
		return result, nil
	}

	strategy := binding.Strategy
	if strategy == "" {
		strategy = opts.DefaultStrategy
	}

	// Build token mappings
	mappings := buildMappings(ds, contract, binding, strategy)

	// Generate output
	switch opts.Format {
	case FormatCSS:
		result.CSS = generateCSS(contract.Prefix, mappings, binding.ThemeMode)
	case FormatSCSS:
		result.CSS = generateSCSS(contract.Prefix, mappings, binding.ThemeMode)
	case FormatTypeScript:
		result.CSS = generateTypeScript(binding.Component, contract.Prefix, mappings)
	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return result, nil
}

// TokenBindingResult represents a resolved token mapping.
type TokenBindingResult struct {
	CSSProperty string
	Value       string
	Transform   string
	Source      string // "explicit", "semantic", "inherit"
}

// findContract finds the theming contract for a component.
func findContract(ds *DesignSystem, binding *ThemeBindings, specDir string) (*ThemingContract, error) {
	// Check local components first
	for i := range ds.Components {
		if ds.Components[i].ID == binding.Component {
			return ds.Components[i].ThemingContract, nil
		}
	}

	// Fetch from specUrl if provided
	if binding.SpecURL != "" {
		return fetchContract(binding.SpecURL, specDir)
	}

	return nil, nil
}

// fetchContract fetches a theming contract from a URL or local file.
func fetchContract(specURL, baseDir string) (*ThemingContract, error) {
	var data []byte
	var err error

	if strings.HasPrefix(specURL, "http://") || strings.HasPrefix(specURL, "https://") {
		// HTTP fetch
		// #nosec G107 -- URL from design system spec, user-controlled configuration
		resp, err := http.Get(specURL)
		if err != nil {
			return nil, fmt.Errorf("fetch %s: %w", specURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch %s: status %d", specURL, resp.StatusCode)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
	} else {
		// Local file
		path := specURL
		if !filepath.IsAbs(path) && baseDir != "" {
			path = filepath.Join(baseDir, path)
		}

		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
	}

	// Parse component JSON
	var component Component
	if err := json.Unmarshal(data, &component); err != nil {
		return nil, fmt.Errorf("parse component: %w", err)
	}

	return component.ThemingContract, nil
}

// buildMappings builds the final token-to-value mappings.
func buildMappings(ds *DesignSystem, contract *ThemingContract, binding *ThemeBindings, strategy string) []TokenBindingResult {
	// Build lookup from explicit mappings
	explicitMappings := make(map[string]TokenMapping)
	for _, m := range binding.Mappings {
		explicitMappings[m.To] = m
	}

	// Build semantic token lookup from design system
	semanticTokens := buildSemanticLookup(ds, binding.ThemeMode)

	var results []TokenBindingResult

	for _, token := range contract.Tokens {
		var result TokenBindingResult
		result.CSSProperty = token.CSSProperty

		// Try explicit mapping first
		if m, ok := explicitMappings[token.ID]; ok {
			result.Value = resolveTokenValue(ds, m.From, binding.ThemeMode)
			result.Transform = m.Transform
			result.Source = "explicit"
			results = append(results, result)
			continue
		}

		// Handle based on strategy
		switch strategy {
		case "explicit":
			// Skip unmapped tokens
			continue

		case "semantic":
			// Try semantic auto-mapping
			if token.Semantic != "" {
				if value, ok := semanticTokens[token.Semantic]; ok {
					result.Value = value
					result.Source = "semantic"
					results = append(results, result)
					continue
				}
			}
			// Fall back to component defaults
			result.Value = getDefault(token, binding.ThemeMode)
			result.Source = "inherit"
			if result.Value != "" {
				results = append(results, result)
			}

		case "inherit":
			// Use component defaults
			result.Value = getDefault(token, binding.ThemeMode)
			result.Source = "inherit"
			if result.Value != "" {
				results = append(results, result)
			}
		}
	}

	return results
}

// buildSemanticLookup builds a map of semantic names to token values.
func buildSemanticLookup(ds *DesignSystem, themeMode string) map[string]string {
	lookup := make(map[string]string)

	// Map colors by common semantic names
	for _, color := range ds.Foundations.Colors {
		name := strings.ToLower(color.ID)
		// Map common patterns
		switch {
		case strings.Contains(name, "primary"):
			lookup["primary"] = modeValue(color, themeMode)
		case strings.Contains(name, "secondary"):
			lookup["secondary"] = modeValue(color, themeMode)
		case strings.Contains(name, "accent"):
			lookup["accent"] = modeValue(color, themeMode)
		case strings.Contains(name, "danger") || strings.Contains(name, "error"):
			lookup["danger"] = modeValue(color, themeMode)
		case strings.Contains(name, "warning"):
			lookup["warning"] = modeValue(color, themeMode)
		case strings.Contains(name, "success"):
			lookup["success"] = modeValue(color, themeMode)
		case strings.Contains(name, "info"):
			lookup["info"] = modeValue(color, themeMode)
		case strings.Contains(name, "neutral") || strings.Contains(name, "gray"):
			if _, exists := lookup["neutral"]; !exists {
				lookup["neutral"] = modeValue(color, themeMode)
			}
		case strings.Contains(name, "surface") || strings.Contains(name, "background"):
			if _, exists := lookup["surface"]; !exists {
				lookup["surface"] = modeValue(color, themeMode)
			}
		case strings.Contains(name, "text"):
			if strings.Contains(name, "muted") {
				lookup["text-muted"] = modeValue(color, themeMode)
			} else if strings.Contains(name, "inverse") {
				lookup["text-inverse"] = modeValue(color, themeMode)
			} else if _, exists := lookup["text"]; !exists {
				lookup["text"] = modeValue(color, themeMode)
			}
		case strings.Contains(name, "border"):
			lookup["border"] = modeValue(color, themeMode)
		case strings.Contains(name, "focus"):
			lookup["focus"] = modeValue(color, themeMode)
		case strings.Contains(name, "disabled"):
			lookup["disabled"] = modeValue(color, themeMode)
		}
	}

	return lookup
}

// resolveTokenValue resolves a token reference to its value.
func resolveTokenValue(ds *DesignSystem, tokenRef string, themeMode string) string {
	// Support dot notation (colors.primary-500) or direct ID (primary-500)
	parts := strings.SplitN(tokenRef, ".", 2)
	var tokenID string
	if len(parts) == 2 {
		tokenID = parts[1]
	} else {
		tokenID = tokenRef
	}

	// Search colors (mode-aware: prefer the token's value for the active mode)
	for _, color := range ds.Foundations.Colors {
		if color.ID == tokenID || color.ID == tokenRef {
			return modeValue(color, themeMode)
		}
	}

	// If not found, return as CSS variable reference
	return fmt.Sprintf("var(--%s)", strings.ReplaceAll(tokenRef, ".", "-"))
}

// getDefault returns the appropriate default value for a token.
func getDefault(token ThemeToken, themeMode string) string {
	defaults := token.EffectiveDefaults()
	if themeMode != "" {
		return defaults[themeMode]
	}
	// Prefer dark for PlexusOne consistency
	if v := defaults["dark"]; v != "" {
		return v
	}
	return defaults["light"]
}

// generateCSS generates CSS output.
func generateCSS(prefix string, mappings []TokenBindingResult, themeMode string) string {
	if len(mappings) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("/* Theme bindings for ")
	sb.WriteString(prefix)
	sb.WriteString(" */\n")

	selector := ":root"
	if themeMode == "dark" {
		selector = ":root[data-theme='dark'], .dark"
	} else if themeMode == "light" {
		selector = ":root[data-theme='light'], .light"
	}

	sb.WriteString(selector)
	sb.WriteString(" {\n")

	// Sort for consistent output
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].CSSProperty < mappings[j].CSSProperty
	})

	for _, m := range mappings {
		value := m.Value
		if m.Transform != "" {
			value = applyTransform(value, m.Transform)
		}
		sb.WriteString("  ")
		sb.WriteString(m.CSSProperty)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString(";")
		if m.Source != "explicit" {
			sb.WriteString(" /* ")
			sb.WriteString(m.Source)
			sb.WriteString(" */")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateSCSS generates SCSS output.
func generateSCSS(prefix string, mappings []TokenBindingResult, themeMode string) string {
	if len(mappings) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("// Theme bindings for ")
	sb.WriteString(prefix)
	sb.WriteString("\n\n")

	// Generate SCSS variables
	for _, m := range mappings {
		varName := strings.TrimPrefix(m.CSSProperty, "--")
		value := m.Value
		if m.Transform != "" {
			value = applyTransform(value, m.Transform)
		}
		sb.WriteString("$")
		sb.WriteString(varName)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString(";\n")
	}

	sb.WriteString("\n// CSS custom properties\n")
	sb.WriteString(generateCSS(prefix, mappings, themeMode))

	return sb.String()
}

// generateTypeScript generates TypeScript output.
func generateTypeScript(component, prefix string, mappings []TokenBindingResult) string {
	if len(mappings) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("// Theme bindings for ")
	sb.WriteString(component)
	sb.WriteString("\n\n")

	sb.WriteString("export const ")
	sb.WriteString(toCamelCase(component))
	sb.WriteString("Theme = {\n")

	for _, m := range mappings {
		// Convert CSS property to JS property name
		propName := strings.TrimPrefix(m.CSSProperty, prefix+"-")
		propName = toCamelCase(propName)

		value := m.Value
		if m.Transform != "" {
			value = applyTransform(value, m.Transform)
		}

		sb.WriteString("  ")
		sb.WriteString(propName)
		sb.WriteString(": '")
		sb.WriteString(value)
		sb.WriteString("',\n")
	}

	sb.WriteString("} as const;\n\n")

	// Generate CSS variable mapping
	sb.WriteString("export const ")
	sb.WriteString(toCamelCase(component))
	sb.WriteString("CSSVars = {\n")

	for _, m := range mappings {
		propName := strings.TrimPrefix(m.CSSProperty, prefix+"-")
		propName = toCamelCase(propName)

		sb.WriteString("  ")
		sb.WriteString(propName)
		sb.WriteString(": '")
		sb.WriteString(m.CSSProperty)
		sb.WriteString("',\n")
	}

	sb.WriteString("} as const;\n")

	return sb.String()
}

// applyTransform applies a CSS transform to a value.
func applyTransform(value, transform string) string {
	if transform == "" {
		return value
	}
	// Simple transform application
	return fmt.Sprintf("%s(%s)", transform, value)
}

// toCamelCase converts a kebab-case string to camelCase.
func toCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// WriteBindings writes all generated bindings to a writer.
func WriteBindings(w io.Writer, bindings []GeneratedBinding) error {
	for i, b := range bindings {
		if b.CSS == "" {
			continue
		}
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, b.CSS); err != nil {
			return err
		}
	}
	return nil
}
