package dss

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PackageTarget represents a framework-specific output target.
type PackageTarget string

const (
	TargetCSS            PackageTarget = "css"
	TargetTailwind       PackageTarget = "tailwind"
	TargetShadCN         PackageTarget = "shadcn"
	TargetMkDocsMaterial PackageTarget = "mkdocs-material"
	TargetSCSS           PackageTarget = "scss"
	TargetJSON           PackageTarget = "json"
	TargetW3C            PackageTarget = "w3c"
)

// PackageGeneratorOptions configures NPM package generation.
type PackageGeneratorOptions struct {
	// OutputDir is the directory to write the package to.
	OutputDir string

	// Targets specifies which framework outputs to generate.
	Targets []PackageTarget

	// Scope is the NPM scope (e.g., "@plexusone").
	// If empty, derived from meta.
	Scope string

	// PackageName is the package name (default: "design-tokens").
	PackageName string

	// DryRun previews output without writing files.
	DryRun bool

	// IncludeReadme generates a README.md file.
	IncludeReadme bool
}

// DefaultPackageOptions returns sensible defaults for package generation.
func DefaultPackageOptions() PackageGeneratorOptions {
	return PackageGeneratorOptions{
		Targets:       []PackageTarget{TargetCSS, TargetTailwind},
		PackageName:   "design-tokens",
		IncludeReadme: true,
	}
}

// GeneratePackage generates a complete NPM package from the design system.
func (ds *DesignSystem) GeneratePackage(opts PackageGeneratorOptions) error {
	if opts.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	// Create output directory
	if !opts.DryRun {
		if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
	}

	// Determine package scope and name
	scope := opts.Scope
	if scope == "" {
		scope = ds.deriveScope()
	}

	pkgName := opts.PackageName
	if pkgName == "" {
		pkgName = "design-tokens"
	}

	fullName := pkgName
	if scope != "" {
		fullName = scope + "/" + pkgName
	}

	// Track generated files for exports
	var exports []string

	// Generate each target
	for _, target := range opts.Targets {
		switch target {
		case TargetCSS:
			if err := ds.generateCSSPackage(opts); err != nil {
				return fmt.Errorf("generating CSS: %w", err)
			}
			exports = append(exports, "./css")

		case TargetTailwind:
			if err := ds.generateTailwindPackage(opts); err != nil {
				return fmt.Errorf("generating Tailwind: %w", err)
			}
			exports = append(exports, "./tailwind")

		case TargetShadCN:
			if err := ds.generateShadCNPackage(opts); err != nil {
				return fmt.Errorf("generating ShadCN: %w", err)
			}
			exports = append(exports, "./shadcn")

		case TargetMkDocsMaterial:
			if err := ds.generateMkDocsPackage(opts); err != nil {
				return fmt.Errorf("generating MkDocs: %w", err)
			}
			exports = append(exports, "./mkdocs")

		case TargetSCSS:
			if err := ds.generateSCSSPackage(opts); err != nil {
				return fmt.Errorf("generating SCSS: %w", err)
			}
			exports = append(exports, "./scss")

		case TargetJSON:
			if err := ds.generateJSONPackage(opts); err != nil {
				return fmt.Errorf("generating JSON: %w", err)
			}

		case TargetW3C:
			if err := ds.generateW3CPackage(opts); err != nil {
				return fmt.Errorf("generating W3C: %w", err)
			}
		}
	}

	// Generate package.json
	if err := ds.generatePackageJSON(opts, fullName, exports); err != nil {
		return fmt.Errorf("generating package.json: %w", err)
	}

	// Generate index files
	if err := ds.generateIndexFiles(opts); err != nil {
		return fmt.Errorf("generating index files: %w", err)
	}

	// Generate README
	if opts.IncludeReadme {
		if err := ds.generateReadme(opts, fullName); err != nil {
			return fmt.Errorf("generating README: %w", err)
		}
	}

	return nil
}

// deriveScope extracts NPM scope from meta.
func (ds *DesignSystem) deriveScope() string {
	// Try to extract from repository URL
	if ds.Meta.Repository != "" {
		parts := strings.Split(ds.Meta.Repository, "/")
		if len(parts) >= 2 {
			org := parts[len(parts)-2]
			if strings.HasPrefix(org, "github.com") {
				org = parts[len(parts)-1]
			}
			return "@" + strings.ToLower(org)
		}
	}

	// Fall back to lowercase name
	name := strings.ToLower(ds.Meta.Name)
	name = strings.ReplaceAll(name, " ", "")
	return "@" + name
}

// generateCSSPackage generates CSS custom properties.
func (ds *DesignSystem) generateCSSPackage(opts PackageGeneratorOptions) error {
	cssDir := filepath.Join(opts.OutputDir, "css")
	if !opts.DryRun {
		if err := os.MkdirAll(cssDir, 0755); err != nil {
			return err
		}
	}

	cssOpts := DefaultCSSOptions()
	cssOpts.Format = "css-vars"
	cssOpts.IncludeComments = true

	css, err := ds.GenerateCSS(cssOpts)
	if err != nil {
		return err
	}

	if !opts.DryRun {
		return os.WriteFile(filepath.Join(cssDir, "tokens.css"), []byte(css), 0600)
	}
	return nil
}

// generateTailwindPackage generates Tailwind v4 preset.
func (ds *DesignSystem) generateTailwindPackage(opts PackageGeneratorOptions) error {
	twDir := filepath.Join(opts.OutputDir, "tailwind")
	if !opts.DryRun {
		if err := os.MkdirAll(twDir, 0755); err != nil {
			return err
		}
	}

	// Generate theme.css (Tailwind v4 @theme block)
	cssOpts := DefaultCSSOptions()
	cssOpts.Format = "tailwind4"
	cssOpts.IncludeComments = true

	themeCSS, err := ds.GenerateCSS(cssOpts)
	if err != nil {
		return err
	}

	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(twDir, "theme.css"), []byte(themeCSS), 0600); err != nil {
			return err
		}
	}

	// Generate preset.js
	preset := ds.generateTailwindPreset()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(twDir, "preset.js"), []byte(preset), 0600); err != nil {
			return err
		}
	}

	// Generate preset.d.ts
	presetTypes := ds.generateTailwindPresetTypes()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(twDir, "preset.d.ts"), []byte(presetTypes), 0600); err != nil {
			return err
		}
	}

	return nil
}

// generateTailwindPreset generates the Tailwind preset JavaScript.
func (ds *DesignSystem) generateTailwindPreset() string {
	var b strings.Builder

	b.WriteString("/** @type {import('tailwindcss').Config} */\n")
	b.WriteString("export default {\n")
	b.WriteString("  theme: {\n")
	b.WriteString("    extend: {\n")

	// Colors
	if len(ds.Foundations.Colors) > 0 {
		b.WriteString("      colors: {\n")
		for _, c := range ds.Foundations.Colors {
			varName := fmt.Sprintf("var(--color-%s)", normalizeID(c.ID))
			b.WriteString(fmt.Sprintf("        '%s': '%s',\n", normalizeID(c.ID), varName))
		}
		b.WriteString("      },\n")
	}

	// Font families
	if ds.Foundations.Typography != nil && len(ds.Foundations.Typography.FontFamilies) > 0 {
		b.WriteString("      fontFamily: {\n")
		for _, ff := range ds.Foundations.Typography.FontFamilies {
			stack := ff.Stack
			if stack == "" {
				stack = ff.Value
			}
			// Parse font stack into array
			fonts := strings.Split(stack, ",")
			var fontArr []string
			for _, f := range fonts {
				fontArr = append(fontArr, "'"+strings.TrimSpace(f)+"'")
			}
			b.WriteString(fmt.Sprintf("        '%s': [%s],\n", normalizeID(ff.ID), strings.Join(fontArr, ", ")))
		}
		b.WriteString("      },\n")
	}

	// Border radius
	if len(ds.Foundations.BorderRadius) > 0 {
		b.WriteString("      borderRadius: {\n")
		for _, br := range ds.Foundations.BorderRadius {
			varName := fmt.Sprintf("var(--radius-%s)", normalizeID(br.ID))
			b.WriteString(fmt.Sprintf("        '%s': '%s',\n", normalizeID(br.ID), varName))
		}
		b.WriteString("      },\n")
	}

	b.WriteString("    },\n")
	b.WriteString("  },\n")
	b.WriteString("}\n")

	return b.String()
}

// generateTailwindPresetTypes generates TypeScript types for the preset.
func (ds *DesignSystem) generateTailwindPresetTypes() string {
	var b strings.Builder

	b.WriteString("import type { Config } from 'tailwindcss'\n\n")
	b.WriteString("declare const config: Partial<Config>\n")
	b.WriteString("export default config\n")

	return b.String()
}

// generateShadCNPackage generates ShadCN/UI theme files.
func (ds *DesignSystem) generateShadCNPackage(opts PackageGeneratorOptions) error {
	shadcnDir := filepath.Join(opts.OutputDir, "shadcn")
	if !opts.DryRun {
		if err := os.MkdirAll(shadcnDir, 0755); err != nil {
			return err
		}
	}

	// Generate theme.css
	themeCSS := ds.generateShadCNThemeCSS()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(shadcnDir, "theme.css"), []byte(themeCSS), 0600); err != nil {
			return err
		}
	}

	// Generate colors.json
	colorsJSON := ds.generateShadCNColors()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(shadcnDir, "colors.json"), []byte(colorsJSON), 0600); err != nil {
			return err
		}
	}

	return nil
}

// generateShadCNThemeCSS generates ShadCN-compatible CSS variables.
func (ds *DesignSystem) generateShadCNThemeCSS() string {
	var b strings.Builder

	b.WriteString("/* ")
	b.WriteString(ds.Meta.Name)
	b.WriteString(" - ShadCN Theme */\n\n")

	b.WriteString("@layer base {\n")
	b.WriteString("  :root {\n")

	// Map colors to ShadCN semantic names
	colorMap := ds.mapColorsToShadCN()
	for name, value := range colorMap {
		b.WriteString(fmt.Sprintf("    --%s: %s;\n", name, value))
	}

	b.WriteString("  }\n\n")

	// Dark mode
	b.WriteString("  .dark {\n")
	darkColorMap := ds.mapColorsToShadCNDark()
	for name, value := range darkColorMap {
		b.WriteString(fmt.Sprintf("    --%s: %s;\n", name, value))
	}
	b.WriteString("  }\n")
	b.WriteString("}\n")

	return b.String()
}

// mapColorsToShadCN maps design system colors to ShadCN semantic names.
func (ds *DesignSystem) mapColorsToShadCN() map[string]string {
	result := make(map[string]string)

	// Default ShadCN structure
	result["background"] = "0 0% 100%"
	result["foreground"] = "222.2 84% 4.9%"
	result["card"] = "0 0% 100%"
	result["card-foreground"] = "222.2 84% 4.9%"
	result["popover"] = "0 0% 100%"
	result["popover-foreground"] = "222.2 84% 4.9%"
	result["primary"] = "221.2 83.2% 53.3%"
	result["primary-foreground"] = "210 40% 98%"
	result["secondary"] = "210 40% 96%"
	result["secondary-foreground"] = "222.2 47.4% 11.2%"
	result["muted"] = "210 40% 96%"
	result["muted-foreground"] = "215.4 16.3% 46.9%"
	result["accent"] = "210 40% 96%"
	result["accent-foreground"] = "222.2 47.4% 11.2%"
	result["destructive"] = "0 84.2% 60.2%"
	result["destructive-foreground"] = "210 40% 98%"
	result["border"] = "214.3 31.8% 91.4%"
	result["input"] = "214.3 31.8% 91.4%"
	result["ring"] = "221.2 83.2% 53.3%"
	result["radius"] = "0.5rem"

	// Override with actual colors from the design system
	for _, c := range ds.Foundations.Colors {
		hsl := hexToHSL(c.Value)
		switch c.ID {
		case "primary", "cyan":
			result["primary"] = hsl
		case "background", "bg":
			result["background"] = hsl
		case "foreground", "text":
			result["foreground"] = hsl
		case "secondary", "purple", "accent":
			result["secondary"] = hsl
			result["accent"] = hsl
		case "border":
			result["border"] = hsl
			result["input"] = hsl
		case "error", "danger", "destructive":
			result["destructive"] = hsl
		case "muted", "text-muted":
			result["muted-foreground"] = hsl
		}
	}

	// Border radius
	for _, br := range ds.Foundations.BorderRadius {
		if br.ID == "default" || br.ID == "md" {
			result["radius"] = br.Value
			break
		}
	}

	return result
}

// mapColorsToShadCNDark maps colors for ShadCN dark mode.
func (ds *DesignSystem) mapColorsToShadCNDark() map[string]string {
	result := make(map[string]string)

	// Default dark mode values
	result["background"] = "222.2 84% 4.9%"
	result["foreground"] = "210 40% 98%"
	result["card"] = "222.2 84% 4.9%"
	result["card-foreground"] = "210 40% 98%"
	result["popover"] = "222.2 84% 4.9%"
	result["popover-foreground"] = "210 40% 98%"
	result["primary"] = "217.2 91.2% 59.8%"
	result["primary-foreground"] = "222.2 47.4% 11.2%"
	result["secondary"] = "217.2 32.6% 17.5%"
	result["secondary-foreground"] = "210 40% 98%"
	result["muted"] = "217.2 32.6% 17.5%"
	result["muted-foreground"] = "215 20.2% 65.1%"
	result["accent"] = "217.2 32.6% 17.5%"
	result["accent-foreground"] = "210 40% 98%"
	result["destructive"] = "0 62.8% 30.6%"
	result["destructive-foreground"] = "210 40% 98%"
	result["border"] = "217.2 32.6% 17.5%"
	result["input"] = "217.2 32.6% 17.5%"
	result["ring"] = "224.3 76.3% 48%"

	// Override with actual dark colors if the design system is dark-first
	for _, c := range ds.Foundations.Colors {
		hsl := hexToHSL(c.Value)
		switch c.ID {
		case "dark", "bg", "background":
			result["background"] = hsl
			result["card"] = hsl
			result["popover"] = hsl
		case "text", "foreground":
			result["foreground"] = hsl
			result["card-foreground"] = hsl
			result["popover-foreground"] = hsl
		case "primary", "cyan":
			result["primary"] = hsl
		case "secondary", "purple":
			result["secondary"] = hsl
		case "border":
			result["border"] = hsl
			result["input"] = hsl
		}
	}

	return result
}

// generateShadCNColors generates colors.json for ShadCN.
func (ds *DesignSystem) generateShadCNColors() string {
	colors := make(map[string]interface{})

	for _, c := range ds.Foundations.Colors {
		colors[c.ID] = map[string]string{
			"DEFAULT": c.Value,
		}
	}

	data, _ := json.MarshalIndent(colors, "", "  ")
	return string(data)
}

// generateMkDocsPackage generates MkDocs Material theme files.
func (ds *DesignSystem) generateMkDocsPackage(opts PackageGeneratorOptions) error {
	mkdocsDir := filepath.Join(opts.OutputDir, "mkdocs")
	if !opts.DryRun {
		if err := os.MkdirAll(mkdocsDir, 0755); err != nil {
			return err
		}
	}

	// Generate extra.css
	cssOpts := DefaultCSSOptions()
	cssOpts.Format = "mkdocs-material"
	cssOpts.IncludeComments = true

	extraCSS, err := ds.GenerateCSS(cssOpts)
	if err != nil {
		return err
	}

	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(mkdocsDir, "extra.css"), []byte(extraCSS), 0600); err != nil {
			return err
		}
	}

	// Generate palette.yml
	paletteYML := ds.generateMkDocsPalette()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(mkdocsDir, "palette.yml"), []byte(paletteYML), 0600); err != nil {
			return err
		}
	}

	return nil
}

// generateMkDocsPalette generates MkDocs palette configuration.
func (ds *DesignSystem) generateMkDocsPalette() string {
	return `# MkDocs Material Palette Configuration
# Add this to your mkdocs.yml under theme.palette

palette:
  - scheme: slate
    primary: custom
    accent: custom
    toggle:
      icon: material/brightness-4
      name: Switch to light mode
  - scheme: default
    primary: custom
    accent: custom
    toggle:
      icon: material/brightness-7
      name: Switch to dark mode
`
}

// generateSCSSPackage generates SCSS variables.
func (ds *DesignSystem) generateSCSSPackage(opts PackageGeneratorOptions) error {
	scssDir := filepath.Join(opts.OutputDir, "scss")
	if !opts.DryRun {
		if err := os.MkdirAll(scssDir, 0755); err != nil {
			return err
		}
	}

	cssOpts := DefaultCSSOptions()
	cssOpts.Format = "scss"
	cssOpts.IncludeComments = true

	scss, err := ds.GenerateCSS(cssOpts)
	if err != nil {
		return err
	}

	if !opts.DryRun {
		return os.WriteFile(filepath.Join(scssDir, "_variables.scss"), []byte(scss), 0600)
	}
	return nil
}

// generateJSONPackage generates raw JSON tokens.
func (ds *DesignSystem) generateJSONPackage(opts PackageGeneratorOptions) error {
	tokens := map[string]interface{}{
		"name":        ds.Meta.Name,
		"version":     ds.Meta.Version,
		"foundations": ds.Foundations,
	}

	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}

	if !opts.DryRun {
		return os.WriteFile(filepath.Join(opts.OutputDir, "tokens.json"), data, 0600)
	}
	return nil
}

// generateW3CPackage generates W3C Design Tokens format.
func (ds *DesignSystem) generateW3CPackage(opts PackageGeneratorOptions) error {
	w3cOpts := W3CTokensOptions{
		IncludeDescriptions: true,
		IncludeExtensions:   true,
	}

	w3cTokens, err := ds.GenerateW3CTokens(w3cOpts)
	if err != nil {
		return err
	}

	if !opts.DryRun {
		return os.WriteFile(filepath.Join(opts.OutputDir, "tokens.w3c.json"), []byte(w3cTokens), 0600)
	}
	return nil
}

// generatePackageJSON generates the package.json file.
func (ds *DesignSystem) generatePackageJSON(opts PackageGeneratorOptions, fullName string, exports []string) error {
	pkg := map[string]interface{}{
		"name":        fullName,
		"version":     ds.Meta.Version,
		"description": fmt.Sprintf("Design tokens for %s", ds.Meta.Name),
		"main":        "index.js",
		"module":      "index.mjs",
		"types":       "index.d.ts",
		"exports": map[string]interface{}{
			".": map[string]string{
				"import":  "./index.mjs",
				"require": "./index.js",
				"types":   "./index.d.ts",
			},
		},
		"files": []string{
			"css",
			"tailwind",
			"shadcn",
			"mkdocs",
			"scss",
			"*.js",
			"*.mjs",
			"*.d.ts",
			"*.json",
		},
		"keywords": []string{
			"design-tokens",
			"design-system",
			"css",
			"tailwind",
			ds.Meta.Name,
		},
		"license":     "MIT",
		"generatedBy": "systemspec-designsystem",
	}

	// Add exports for each target
	exportsMap := pkg["exports"].(map[string]interface{})
	for _, exp := range exports {
		switch exp {
		case "./css":
			exportsMap["./css"] = "./css/tokens.css"
		case "./tailwind":
			exportsMap["./tailwind"] = map[string]string{
				"import": "./tailwind/preset.js",
				"types":  "./tailwind/preset.d.ts",
			}
		case "./shadcn":
			exportsMap["./shadcn"] = "./shadcn/theme.css"
		case "./mkdocs":
			exportsMap["./mkdocs"] = "./mkdocs/extra.css"
		case "./scss":
			exportsMap["./scss"] = "./scss/_variables.scss"
		}
	}

	// Add repository if available
	if ds.Meta.Repository != "" {
		pkg["repository"] = map[string]string{
			"type": "git",
			"url":  ds.Meta.Repository,
		}
	}

	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	if !opts.DryRun {
		return os.WriteFile(filepath.Join(opts.OutputDir, "package.json"), data, 0600)
	}
	return nil
}

// generateIndexFiles generates the main entry point files.
func (ds *DesignSystem) generateIndexFiles(opts PackageGeneratorOptions) error {
	// Generate index.js (CommonJS)
	indexJS := ds.generateIndexJS()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(opts.OutputDir, "index.js"), []byte(indexJS), 0600); err != nil {
			return err
		}
	}

	// Generate index.mjs (ESM)
	indexMJS := ds.generateIndexMJS()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(opts.OutputDir, "index.mjs"), []byte(indexMJS), 0600); err != nil {
			return err
		}
	}

	// Generate index.d.ts (TypeScript)
	indexDTS := ds.generateIndexDTS()
	if !opts.DryRun {
		if err := os.WriteFile(filepath.Join(opts.OutputDir, "index.d.ts"), []byte(indexDTS), 0600); err != nil {
			return err
		}
	}

	return nil
}

// generateIndexJS generates CommonJS entry point.
func (ds *DesignSystem) generateIndexJS() string {
	var b strings.Builder

	b.WriteString("'use strict';\n\n")
	b.WriteString(fmt.Sprintf("// %s Design Tokens\n\n", ds.Meta.Name))

	// Export colors
	b.WriteString("const colors = {\n")
	for _, c := range ds.Foundations.Colors {
		b.WriteString(fmt.Sprintf("  '%s': '%s',\n", normalizeID(c.ID), c.Value))
	}
	b.WriteString("};\n\n")

	// Export spacing
	if ds.Foundations.Spacing != nil {
		b.WriteString("const spacing = {\n")
		for _, s := range ds.Foundations.Spacing.Scale {
			b.WriteString(fmt.Sprintf("  '%s': '%s',\n", normalizeID(s.ID), s.Value))
		}
		b.WriteString("};\n\n")
	}

	// Export border radius
	if len(ds.Foundations.BorderRadius) > 0 {
		b.WriteString("const borderRadius = {\n")
		for _, br := range ds.Foundations.BorderRadius {
			b.WriteString(fmt.Sprintf("  '%s': '%s',\n", normalizeID(br.ID), br.Value))
		}
		b.WriteString("};\n\n")
	}

	b.WriteString("module.exports = {\n")
	b.WriteString("  colors,\n")
	if ds.Foundations.Spacing != nil {
		b.WriteString("  spacing,\n")
	}
	if len(ds.Foundations.BorderRadius) > 0 {
		b.WriteString("  borderRadius,\n")
	}
	b.WriteString("};\n")

	return b.String()
}

// generateIndexMJS generates ESM entry point.
func (ds *DesignSystem) generateIndexMJS() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("// %s Design Tokens\n\n", ds.Meta.Name))

	// Export colors
	b.WriteString("export const colors = {\n")
	for _, c := range ds.Foundations.Colors {
		b.WriteString(fmt.Sprintf("  '%s': '%s',\n", normalizeID(c.ID), c.Value))
	}
	b.WriteString("};\n\n")

	// Export spacing
	if ds.Foundations.Spacing != nil {
		b.WriteString("export const spacing = {\n")
		for _, s := range ds.Foundations.Spacing.Scale {
			b.WriteString(fmt.Sprintf("  '%s': '%s',\n", normalizeID(s.ID), s.Value))
		}
		b.WriteString("};\n\n")
	}

	// Export border radius
	if len(ds.Foundations.BorderRadius) > 0 {
		b.WriteString("export const borderRadius = {\n")
		for _, br := range ds.Foundations.BorderRadius {
			b.WriteString(fmt.Sprintf("  '%s': '%s',\n", normalizeID(br.ID), br.Value))
		}
		b.WriteString("};\n\n")
	}

	return b.String()
}

// generateIndexDTS generates TypeScript declarations.
func (ds *DesignSystem) generateIndexDTS() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("// %s Design Tokens\n\n", ds.Meta.Name))

	// Colors type
	b.WriteString("export declare const colors: {\n")
	for _, c := range ds.Foundations.Colors {
		b.WriteString(fmt.Sprintf("  readonly '%s': '%s';\n", normalizeID(c.ID), c.Value))
	}
	b.WriteString("};\n\n")

	// Spacing type
	if ds.Foundations.Spacing != nil {
		b.WriteString("export declare const spacing: {\n")
		for _, s := range ds.Foundations.Spacing.Scale {
			b.WriteString(fmt.Sprintf("  readonly '%s': '%s';\n", normalizeID(s.ID), s.Value))
		}
		b.WriteString("};\n\n")
	}

	// Border radius type
	if len(ds.Foundations.BorderRadius) > 0 {
		b.WriteString("export declare const borderRadius: {\n")
		for _, br := range ds.Foundations.BorderRadius {
			b.WriteString(fmt.Sprintf("  readonly '%s': '%s';\n", normalizeID(br.ID), br.Value))
		}
		b.WriteString("};\n\n")
	}

	// Export color names as union type
	b.WriteString("export type ColorName = ")
	var colorNames []string
	for _, c := range ds.Foundations.Colors {
		colorNames = append(colorNames, fmt.Sprintf("'%s'", normalizeID(c.ID)))
	}
	b.WriteString(strings.Join(colorNames, " | "))
	b.WriteString(";\n")

	return b.String()
}

// generateReadme generates README.md for the package.
func (ds *DesignSystem) generateReadme(opts PackageGeneratorOptions, fullName string) error {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s\n\n", fullName))
	b.WriteString(fmt.Sprintf("Design tokens for %s.\n\n", ds.Meta.Name))

	if ds.Meta.Description != "" {
		b.WriteString(fmt.Sprintf("> %s\n\n", ds.Meta.Description))
	}

	b.WriteString("## Installation\n\n")
	b.WriteString("```bash\n")
	b.WriteString(fmt.Sprintf("npm install %s\n", fullName))
	b.WriteString("```\n\n")

	b.WriteString("## Usage\n\n")

	// CSS
	b.WriteString("### CSS Variables\n\n")
	b.WriteString("```css\n")
	b.WriteString(fmt.Sprintf("@import '%s/css/tokens.css';\n\n", fullName))
	b.WriteString(".my-element {\n")
	b.WriteString("  color: var(--color-primary);\n")
	b.WriteString("}\n")
	b.WriteString("```\n\n")

	// Tailwind
	b.WriteString("### Tailwind CSS\n\n")
	b.WriteString("```javascript\n")
	b.WriteString("// tailwind.config.js\n")
	b.WriteString(fmt.Sprintf("import preset from '%s/tailwind'\n\n", fullName))
	b.WriteString("export default {\n")
	b.WriteString("  presets: [preset],\n")
	b.WriteString("  content: ['./src/**/*.{js,ts,jsx,tsx}'],\n")
	b.WriteString("}\n")
	b.WriteString("```\n\n")

	// JavaScript
	b.WriteString("### JavaScript/TypeScript\n\n")
	b.WriteString("```javascript\n")
	b.WriteString(fmt.Sprintf("import { colors, spacing } from '%s'\n\n", fullName))
	b.WriteString("console.log(colors.primary) // #...\n")
	b.WriteString("```\n\n")

	// Tokens
	b.WriteString("## Available Tokens\n\n")

	b.WriteString("### Colors\n\n")
	b.WriteString("| Token | Value |\n")
	b.WriteString("|-------|-------|\n")
	for _, c := range ds.Foundations.Colors {
		b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", c.ID, c.Value))
	}
	b.WriteString("\n")

	b.WriteString("---\n\n")
	b.WriteString("Generated by [systemspec-designsystem](https://github.com/plexusone/systemspec-designsystem)\n")

	if !opts.DryRun {
		return os.WriteFile(filepath.Join(opts.OutputDir, "README.md"), []byte(b.String()), 0600)
	}
	return nil
}

// hexToHSL converts a hex color to HSL format for ShadCN.
func hexToHSL(hex string) string {
	// Remove # prefix
	hex = strings.TrimPrefix(hex, "#")

	// Parse RGB values
	if len(hex) != 6 {
		return "0 0% 0%"
	}

	r := float64(hexToDec(hex[0:2])) / 255.0
	g := float64(hexToDec(hex[2:4])) / 255.0
	b := float64(hexToDec(hex[4:6])) / 255.0

	max := maxFloat(r, maxFloat(g, b))
	min := minFloat(r, minFloat(g, b))

	l := (max + min) / 2.0

	var h, s float64

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min

		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}

		h /= 6
	}

	return fmt.Sprintf("%.1f %.1f%% %.1f%%", h*360, s*100, l*100)
}

func hexToDec(hex string) int {
	var result int
	_, _ = fmt.Sscanf(hex, "%x", &result)
	return result
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// ParseTargets parses a comma-separated list of target names.
func ParseTargets(targets string) []PackageTarget {
	if targets == "" {
		return []PackageTarget{TargetCSS, TargetTailwind}
	}

	parts := strings.Split(targets, ",")
	var result []PackageTarget

	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		switch p {
		case "css":
			result = append(result, TargetCSS)
		case "tailwind":
			result = append(result, TargetTailwind)
		case "shadcn":
			result = append(result, TargetShadCN)
		case "mkdocs-material", "mkdocs":
			result = append(result, TargetMkDocsMaterial)
		case "scss":
			result = append(result, TargetSCSS)
		case "json":
			result = append(result, TargetJSON)
		case "w3c":
			result = append(result, TargetW3C)
		case "all":
			return []PackageTarget{
				TargetCSS,
				TargetTailwind,
				TargetShadCN,
				TargetMkDocsMaterial,
				TargetSCSS,
				TargetJSON,
				TargetW3C,
			}
		}
	}

	return result
}
