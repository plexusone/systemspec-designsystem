package dss

import (
	"fmt"
	"sort"
	"strings"
)

// LLMPromptOptions configures the LLM context prompt generation.
type LLMPromptOptions struct {
	// Format: "markdown" (default), "xml", "json"
	Format string

	// IncludeFoundations includes color/spacing/typography guidance
	IncludeFoundations bool

	// IncludeComponents includes component usage guidelines
	IncludeComponents bool

	// IncludePatterns includes pattern recommendations
	IncludePatterns bool

	// IncludeAccessibility includes a11y requirements
	IncludeAccessibility bool

	// IncludeAntiPatterns emphasizes what NOT to do
	IncludeAntiPatterns bool

	// MaxExamples limits code examples per component (0 = all)
	MaxExamples int
}

// DefaultLLMPromptOptions returns comprehensive defaults.
func DefaultLLMPromptOptions() LLMPromptOptions {
	return LLMPromptOptions{
		Format:               "markdown",
		IncludeFoundations:   true,
		IncludeComponents:    true,
		IncludePatterns:      true,
		IncludeAccessibility: true,
		IncludeAntiPatterns:  true,
		MaxExamples:          3,
	}
}

// GenerateLLMPrompt generates a Claude-friendly context prompt from the design system.
func (ds *DesignSystem) GenerateLLMPrompt(opts LLMPromptOptions) (string, error) {
	switch opts.Format {
	case "markdown":
		return ds.generateMarkdownPrompt(opts)
	case "xml":
		return ds.generateXMLPrompt(opts)
	default:
		return ds.generateMarkdownPrompt(opts)
	}
}

func (ds *DesignSystem) generateMarkdownPrompt(opts LLMPromptOptions) (string, error) {
	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("# %s Design System\n\n", ds.Meta.Name))
	if ds.Meta.Description != "" {
		b.WriteString(ds.Meta.Description)
		b.WriteString("\n\n")
	}

	// Design Principles
	if len(ds.Principles) > 0 {
		b.WriteString("## Design Principles\n\n")
		for _, p := range ds.Principles {
			b.WriteString(fmt.Sprintf("### %s\n", p.Name))
			b.WriteString(p.Description)
			b.WriteString("\n\n")
		}
	}

	// Foundations
	if opts.IncludeFoundations {
		writeFoundationsSection(&b, ds.Foundations)
	}

	// Components
	if opts.IncludeComponents && len(ds.Components) > 0 {
		b.WriteString("## Components\n\n")
		for _, c := range ds.Components {
			writeComponentSection(&b, c, opts)
		}
	}

	// Patterns
	if opts.IncludePatterns && len(ds.Patterns) > 0 {
		b.WriteString("## Patterns\n\n")
		for _, p := range ds.Patterns {
			writePatternSection(&b, p, opts)
		}
	}

	// Accessibility
	if opts.IncludeAccessibility && ds.Accessibility != nil {
		writeAccessibilitySection(&b, *ds.Accessibility)
	}

	// Global Anti-Patterns
	if opts.IncludeAntiPatterns {
		writeAntiPatternsSection(&b, ds)
	}

	return b.String(), nil
}

func writeFoundationsSection(b *strings.Builder, f Foundations) {
	b.WriteString("## Design Tokens\n\n")

	// Colors with semantic meaning
	if len(f.Colors) > 0 {
		b.WriteString("### Colors\n\n")
		b.WriteString("| Token | Value | Modes | Usage |\n")
		b.WriteString("|-------|-------|-------|-------|\n")
		for _, c := range f.Colors {
			usage := c.Usage
			if usage == "" {
				usage = c.Semantic
			}
			modes := c.EffectiveModes()
			modeCol := "—"
			if len(modes) > 0 {
				keys := make([]string, 0, len(modes))
				for m := range modes {
					keys = append(keys, m)
				}
				sort.Strings(keys)
				parts := make([]string, 0, len(keys))
				for _, m := range keys {
					parts = append(parts, fmt.Sprintf("%s: `%s`", m, modes[m]))
				}
				modeCol = strings.Join(parts, ", ")
			}
			b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s |\n", c.ID, c.Value, modeCol, usage))
		}
		b.WriteString("\n")
	}

	// Densities
	if len(f.Densities) > 0 {
		b.WriteString("### Densities\n\n")
		b.WriteString("| Density | Scale | Description |\n")
		b.WriteString("|---------|-------|-------------|\n")
		for _, d := range f.Densities {
			b.WriteString(fmt.Sprintf("| `%s` | %v | %s |\n", d.ID, d.Scale, d.Description))
		}
		b.WriteString("\n")
	}

	// Spacing
	if f.Spacing != nil && len(f.Spacing.Scale) > 0 {
		b.WriteString("### Spacing Scale\n\n")
		b.WriteString(fmt.Sprintf("Base unit: `%s`\n\n", f.Spacing.BaseUnit))
		b.WriteString("| Token | Value |\n")
		b.WriteString("|-------|-------|\n")
		for _, s := range f.Spacing.Scale {
			b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", s.ID, s.Value))
		}
		b.WriteString("\n")
	}

	// Border Radius
	if len(f.BorderRadius) > 0 {
		b.WriteString("### Border Radius\n\n")
		for _, br := range f.BorderRadius {
			b.WriteString(fmt.Sprintf("- `%s`: %s\n", br.ID, br.Value))
		}
		b.WriteString("\n")
	}
}

func writeComponentSection(b *strings.Builder, c Component, opts LLMPromptOptions) {
	b.WriteString(fmt.Sprintf("### %s\n\n", c.Name))

	if c.Description != "" {
		b.WriteString(c.Description)
		b.WriteString("\n\n")
	}

	// LLM Context
	if c.LLM != nil {
		if c.LLM.Intent != "" {
			b.WriteString(fmt.Sprintf("**Intent:** %s\n\n", c.LLM.Intent))
		}

		if len(c.LLM.AllowedContexts) > 0 {
			b.WriteString("**Use in:** ")
			b.WriteString(strings.Join(c.LLM.AllowedContexts, ", "))
			b.WriteString("\n\n")
		}

		if len(c.LLM.ForbiddenContexts) > 0 {
			b.WriteString("**Avoid in:** ")
			b.WriteString(strings.Join(c.LLM.ForbiddenContexts, ", "))
			b.WriteString("\n\n")
		}
	}

	// Variants
	if len(c.Variants) > 0 {
		b.WriteString("**Variants:**\n")
		for _, v := range c.Variants {
			defaultMarker := ""
			if v.IsDefault {
				defaultMarker = " (default)"
			}
			desc := v.Description
			if desc == "" {
				desc = v.Name
			}
			b.WriteString(fmt.Sprintf("- `%s`%s: %s\n", v.ID, defaultMarker, desc))
		}
		b.WriteString("\n")
	}

	// Props
	if len(c.Props) > 0 {
		b.WriteString("**Props:**\n")
		for _, p := range c.Props {
			required := ""
			if p.Required {
				required = " (required)"
			}
			b.WriteString(fmt.Sprintf("- `%s`: %s%s", p.Name, p.Type, required))
			if p.Description != "" {
				b.WriteString(fmt.Sprintf(" - %s", p.Description))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Example Usage
	if c.LLM != nil && len(c.LLM.ExampleUsage) > 0 {
		b.WriteString("**Examples:**\n```tsx\n")
		maxExamples := len(c.LLM.ExampleUsage)
		if opts.MaxExamples > 0 && opts.MaxExamples < maxExamples {
			maxExamples = opts.MaxExamples
		}
		for i := 0; i < maxExamples; i++ {
			b.WriteString(c.LLM.ExampleUsage[i])
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
	}

	// Anti-patterns
	if opts.IncludeAntiPatterns && c.LLM != nil && len(c.LLM.AntiPatterns) > 0 {
		b.WriteString("**Don't:**\n")
		for _, ap := range c.LLM.AntiPatterns {
			b.WriteString(fmt.Sprintf("- %s\n", ap))
		}
		b.WriteString("\n")
	}

	// Accessibility
	if c.Accessibility != nil {
		b.WriteString("**Accessibility:**\n")
		if c.Accessibility.Role != "" {
			b.WriteString(fmt.Sprintf("- Role: `%s`\n", c.Accessibility.Role))
		}
		if len(c.Accessibility.KeyboardSupport) > 0 {
			b.WriteString("- Keyboard: ")
			keys := make([]string, len(c.Accessibility.KeyboardSupport))
			for i, k := range c.Accessibility.KeyboardSupport {
				keys[i] = k.Key
			}
			b.WriteString(strings.Join(keys, ", "))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
}

func writePatternSection(b *strings.Builder, p Pattern, _ LLMPromptOptions) {
	b.WriteString(fmt.Sprintf("### %s\n\n", p.Name))

	if p.Description != "" {
		b.WriteString(p.Description)
		b.WriteString("\n\n")
	}

	if p.LLM != nil {
		if p.LLM.Intent != "" {
			b.WriteString(fmt.Sprintf("**Intent:** %s\n\n", p.LLM.Intent))
		}

		if len(p.LLM.ExampleUsage) > 0 {
			b.WriteString("**Example:**\n```tsx\n")
			b.WriteString(p.LLM.ExampleUsage[0])
			b.WriteString("\n```\n\n")
		}
	}

	if len(p.Components) > 0 {
		b.WriteString("**Components used:** ")
		compIDs := make([]string, len(p.Components))
		for i, c := range p.Components {
			compIDs[i] = c.ComponentID
		}
		b.WriteString(strings.Join(compIDs, ", "))
		b.WriteString("\n\n")
	}
}

func writeAccessibilitySection(b *strings.Builder, a Accessibility) {
	b.WriteString("## Accessibility Requirements\n\n")

	b.WriteString(fmt.Sprintf("**Target Level:** WCAG %s %s\n\n", a.WCAGVersion, a.WCAGLevel))

	if a.ColorContrast != nil {
		b.WriteString("### Color Contrast\n")
		b.WriteString(fmt.Sprintf("- Normal text ratio: %.1f:1\n", a.ColorContrast.NormalTextRatio))
		b.WriteString(fmt.Sprintf("- Large text ratio: %.1f:1\n", a.ColorContrast.LargeTextRatio))
		b.WriteString("\n")
	}

	if a.Keyboard != nil {
		b.WriteString("### Keyboard Navigation\n")
		if a.Keyboard.FocusVisible {
			b.WriteString("- All interactive elements must have visible focus indicators\n")
		}
		if a.Keyboard.SkipLinks {
			b.WriteString("- Skip links required for main content\n")
		}
		if a.Keyboard.NoKeyboardTrap {
			b.WriteString("- No keyboard traps allowed\n")
		}
		b.WriteString("\n")
	}

	if len(a.Guidelines) > 0 {
		b.WriteString("### Guidelines\n")
		for _, g := range a.Guidelines {
			b.WriteString(fmt.Sprintf("- **%s**: %s\n", g.Title, g.Description))
		}
		b.WriteString("\n")
	}
}

func writeAntiPatternsSection(b *strings.Builder, ds *DesignSystem) {
	var antiPatterns []string

	// Collect anti-patterns from all components
	for _, c := range ds.Components {
		if c.LLM != nil {
			for _, ap := range c.LLM.AntiPatterns {
				antiPatterns = append(antiPatterns, fmt.Sprintf("[%s] %s", c.Name, ap))
			}
		}
	}

	// Collect from patterns
	for _, p := range ds.Patterns {
		if p.LLM != nil {
			for _, ap := range p.LLM.AntiPatterns {
				antiPatterns = append(antiPatterns, fmt.Sprintf("[%s] %s", p.Name, ap))
			}
		}
	}

	if len(antiPatterns) > 0 {
		b.WriteString("## Anti-Patterns to Avoid\n\n")
		for _, ap := range antiPatterns {
			b.WriteString(fmt.Sprintf("- %s\n", ap))
		}
		b.WriteString("\n")
	}
}

func (ds *DesignSystem) generateXMLPrompt(opts LLMPromptOptions) (string, error) {
	var b strings.Builder

	b.WriteString("<design-system>\n")
	b.WriteString(fmt.Sprintf("  <name>%s</name>\n", ds.Meta.Name))
	if ds.Meta.Description != "" {
		b.WriteString(fmt.Sprintf("  <description>%s</description>\n", ds.Meta.Description))
	}

	// Components
	if opts.IncludeComponents && len(ds.Components) > 0 {
		b.WriteString("  <components>\n")
		for _, c := range ds.Components {
			b.WriteString(fmt.Sprintf("    <component id=\"%s\">\n", c.ID))
			b.WriteString(fmt.Sprintf("      <name>%s</name>\n", c.Name))
			if c.LLM != nil {
				if c.LLM.Intent != "" {
					b.WriteString(fmt.Sprintf("      <intent>%s</intent>\n", c.LLM.Intent))
				}
				if len(c.LLM.AllowedContexts) > 0 {
					b.WriteString(fmt.Sprintf("      <allowed-contexts>%s</allowed-contexts>\n",
						strings.Join(c.LLM.AllowedContexts, ", ")))
				}
				if len(c.LLM.ForbiddenContexts) > 0 {
					b.WriteString(fmt.Sprintf("      <forbidden-contexts>%s</forbidden-contexts>\n",
						strings.Join(c.LLM.ForbiddenContexts, ", ")))
				}
				if len(c.LLM.AntiPatterns) > 0 {
					b.WriteString("      <anti-patterns>\n")
					for _, ap := range c.LLM.AntiPatterns {
						b.WriteString(fmt.Sprintf("        <pattern>%s</pattern>\n", ap))
					}
					b.WriteString("      </anti-patterns>\n")
				}
			}
			b.WriteString("    </component>\n")
		}
		b.WriteString("  </components>\n")
	}

	b.WriteString("</design-system>\n")

	return b.String(), nil
}
