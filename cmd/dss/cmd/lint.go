package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var (
	lintFormat     string
	lintJsonOutput bool
	lintFix        bool
)

// LintReport holds linting results
type LintReport struct {
	File       string      `json:"file"`
	Format     string      `json:"format"`
	Violations []Violation `json:"violations"`
	Summary    LintSummary `json:"summary"`
}

// Violation represents a single lint violation
type Violation struct {
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"` // error, warning, info
	Message  string `json:"message"`
	Expected string `json:"expected,omitempty"`
	Found    string `json:"found,omitempty"`
	Fixable  bool   `json:"fixable,omitempty"`
}

// LintSummary summarizes the lint results
type LintSummary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
	Passed   int `json:"passed"`
}

var lintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Lint CSS/code files against design system spec",
	Long: `Lint files against the design system specification.

Performs deterministic checks to verify CSS and code files comply with
the design system tokens and rules.

Supported formats:
  - tailwind4:        Tailwind v4 CSS with @theme block
  - mkdocs-material:  MkDocs Material theme CSS
  - css-vars:         Plain CSS with custom properties
  - tsx:              React/TypeScript components

Examples:
  # Lint a Tailwind CSS file
  dss lint --format tailwind4 ./src/index.css

  # Lint MkDocs theme CSS
  dss lint --format mkdocs-material ./docs/stylesheets/extra.css

  # Lint React components
  dss lint --format tsx ./src/components/Button.tsx

  # Output as JSON for CI integration
  dss lint --json --format tailwind4 ./src/index.css`,
	Args: cobra.ExactArgs(1),
	RunE: runLint,
}

func init() {
	lintCmd.Flags().StringVar(&lintFormat, "format", "", "file format: tailwind4, mkdocs-material, css-vars, tsx (auto-detected if not specified)")
	lintCmd.Flags().BoolVar(&lintJsonOutput, "json", false, "output results as JSON")
	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "automatically fix violations where possible")

	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	dir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Auto-detect format if not specified
	format := lintFormat
	if format == "" {
		format = detectFormat(filePath, string(content))
	}

	report := &LintReport{
		File:       filePath,
		Format:     format,
		Violations: []Violation{},
	}

	// Run format-specific linting
	switch format {
	case "tailwind4":
		lintTailwind4(string(content), ds, report)
	case "mkdocs-material", "mkdocs":
		lintMkDocsMaterial(string(content), ds, report)
	case "css-vars", "css":
		lintCSSVars(string(content), ds, report)
	case "tsx", "jsx":
		lintTSX(string(content), ds, report)
	default:
		return fmt.Errorf("unknown format: %s (supported: tailwind4, mkdocs-material, css-vars, tsx)", format)
	}

	// Calculate summary
	for _, v := range report.Violations {
		switch v.Severity {
		case "error":
			report.Summary.Errors++
		case "warning":
			report.Summary.Warnings++
		case "info":
			report.Summary.Info++
		}
	}

	// Output results
	if lintJsonOutput {
		jsonData, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		printLintReport(report)
	}

	// Exit with error if there are errors
	if report.Summary.Errors > 0 {
		return fmt.Errorf("lint found %d errors", report.Summary.Errors)
	}

	return nil
}

func detectFormat(filePath, content string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".tsx", ".jsx":
		return "tsx"
	case ".css":
		if strings.Contains(content, "@theme") {
			return "tailwind4"
		}
		if strings.Contains(content, "--md-primary") || strings.Contains(content, "data-md-color-scheme") {
			return "mkdocs-material"
		}
		return "css-vars"
	}

	return "css-vars"
}

// lintTailwind4 checks Tailwind v4 CSS against design system
func lintTailwind4(content string, ds *dss.DesignSystem, report *LintReport) {
	lines := strings.Split(content, "\n")

	// Build expected tokens map
	expectedColors := make(map[string]string)
	for _, c := range ds.Foundations.Colors {
		varName := fmt.Sprintf("--color-%s", normalizeTokenID(c.ID))
		expectedColors[varName] = c.Value
	}

	// Check @theme block exists
	if !strings.Contains(content, "@theme") {
		report.Violations = append(report.Violations, Violation{
			Rule:     "tailwind4-theme-block",
			Severity: "error",
			Message:  "Missing @theme block - Tailwind v4 requires @theme for custom properties",
		})
	}

	// Check each color token
	for varName, expectedValue := range expectedColors {
		pattern := regexp.MustCompile(fmt.Sprintf(`%s:\s*([^;]+);`, regexp.QuoteMeta(varName)))
		match := pattern.FindStringSubmatch(content)

		if match == nil {
			report.Violations = append(report.Violations, Violation{
				Rule:     "missing-color-token",
				Severity: "warning",
				Message:  fmt.Sprintf("Missing color token: %s", varName),
				Expected: expectedValue,
			})
		} else {
			foundValue := strings.TrimSpace(match[1])
			if !colorsEqual(foundValue, expectedValue) {
				// Find line number
				line := findLineNumber(lines, varName)
				report.Violations = append(report.Violations, Violation{
					Line:     line,
					Rule:     "color-mismatch",
					Severity: "error",
					Message:  fmt.Sprintf("Color token %s has incorrect value", varName),
					Expected: expectedValue,
					Found:    foundValue,
					Fixable:  true,
				})
			} else {
				report.Summary.Passed++
			}
		}
	}

	// Check for hardcoded colors outside @theme
	checkHardcodedColorsInCSS(content, lines, ds, report)

	// Check font families
	if ds.Foundations.Typography != nil {
		for _, ff := range ds.Foundations.Typography.FontFamilies {
			varName := fmt.Sprintf("--font-%s", normalizeTokenID(ff.ID))
			expectedValue := ff.Stack
			if expectedValue == "" {
				expectedValue = ff.Value
			}

			pattern := regexp.MustCompile(fmt.Sprintf(`%s:\s*([^;]+);`, regexp.QuoteMeta(varName)))
			match := pattern.FindStringSubmatch(content)

			if match == nil && expectedValue != "" {
				report.Violations = append(report.Violations, Violation{
					Rule:     "missing-font-token",
					Severity: "warning",
					Message:  fmt.Sprintf("Missing font token: %s", varName),
					Expected: expectedValue,
				})
			}
		}
	}
}

// lintMkDocsMaterial checks MkDocs Material CSS
func lintMkDocsMaterial(content string, ds *dss.DesignSystem, report *LintReport) {
	lines := strings.Split(content, "\n")

	// Map design system colors to MkDocs Material variables
	colorMapping := map[string]string{
		"--md-primary-fg-color": "cyan",
		"--md-accent-fg-color":  "purple",
		"--md-default-bg-color": "background",
		"--md-default-fg-color": "foreground",
		"--md-code-bg-color":    "background-subtle",
	}

	dsColors := make(map[string]string)
	for _, c := range ds.Foundations.Colors {
		dsColors[c.ID] = c.Value
	}

	// Check color mappings
	for mdVar, dsColorID := range colorMapping {
		expectedValue, ok := dsColors[dsColorID]
		if !ok {
			continue
		}

		pattern := regexp.MustCompile(fmt.Sprintf(`%s:\s*([^;]+);`, regexp.QuoteMeta(mdVar)))
		match := pattern.FindStringSubmatch(content)

		if match != nil {
			foundValue := strings.TrimSpace(match[1])
			if !colorsEqual(foundValue, expectedValue) {
				line := findLineNumber(lines, mdVar)
				report.Violations = append(report.Violations, Violation{
					Line:     line,
					Rule:     "mkdocs-color-mismatch",
					Severity: "error",
					Message:  fmt.Sprintf("MkDocs variable %s should match design system '%s'", mdVar, dsColorID),
					Expected: expectedValue,
					Found:    foundValue,
					Fixable:  true,
				})
			} else {
				report.Summary.Passed++
			}
		}
	}

	// Check for slate scheme (dark mode default)
	if !strings.Contains(content, `[data-md-color-scheme="slate"]`) {
		report.Violations = append(report.Violations, Violation{
			Rule:     "mkdocs-dark-mode",
			Severity: "warning",
			Message:  "Missing dark mode (slate) styles - design system is dark-first",
		})
	}
}

// lintCSSVars checks plain CSS custom properties
func lintCSSVars(content string, ds *dss.DesignSystem, report *LintReport) {
	lines := strings.Split(content, "\n")

	// Check for hardcoded colors
	checkHardcodedColorsInCSS(content, lines, ds, report)
}

// lintTSX checks React/TypeScript components
func lintTSX(content string, ds *dss.DesignSystem, report *LintReport) {
	lines := strings.Split(content, "\n")

	// Check for hardcoded colors
	hexPattern := regexp.MustCompile(`(?:color|background|border).*?["']?(#[0-9a-fA-F]{3,8})["']?`)
	for i, line := range lines {
		matches := hexPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// Skip if it's in a comment
				if strings.Contains(line[:strings.Index(line, match[0])], "//") {
					continue
				}
				report.Violations = append(report.Violations, Violation{
					Line:     i + 1,
					Rule:     "no-hardcoded-colors",
					Severity: "warning",
					Message:  fmt.Sprintf("Hardcoded color '%s' - use CSS variable or design token", match[1]),
					Found:    match[1],
					Fixable:  false,
				})
			}
		}
	}

	// Check for hardcoded pixel spacing
	pxPattern := regexp.MustCompile(`(?:margin|padding|gap|top|left|right|bottom).*?:\s*["']?(\d+)px["']?`)
	validSpacing := map[string]bool{
		"0": true, "1": true, "2": true, "4": true, "8": true,
		"12": true, "16": true, "20": true, "24": true, "32": true,
		"40": true, "48": true, "64": true, "80": true, "96": true,
	}

	for i, line := range lines {
		matches := pxPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 && !validSpacing[match[1]] {
				report.Violations = append(report.Violations, Violation{
					Line:     i + 1,
					Rule:     "use-spacing-scale",
					Severity: "warning",
					Message:  fmt.Sprintf("Non-standard spacing '%spx' - use design system spacing scale", match[1]),
					Found:    match[1] + "px",
				})
			}
		}
	}

	// Check component variants against spec
	for _, comp := range ds.Components {
		variantPattern := regexp.MustCompile(fmt.Sprintf(`<%s[^>]*variant=["']([^"']+)["']`, comp.Name))
		matches := variantPattern.FindAllStringSubmatch(content, -1)

		validVariants := make(map[string]bool)
		for _, v := range comp.Variants {
			validVariants[v.ID] = true
		}

		for _, match := range matches {
			if len(match) > 1 && !validVariants[match[1]] {
				line := findLineNumber(lines, match[0])
				report.Violations = append(report.Violations, Violation{
					Line:     line,
					Rule:     "invalid-variant",
					Severity: "error",
					Message:  fmt.Sprintf("Invalid variant '%s' for %s component", match[1], comp.Name),
					Found:    match[1],
				})
			}
		}
	}

	// Check anti-patterns
	checkAntiPatternsInTSX(content, lines, ds, report)
}

func checkHardcodedColorsInCSS(_ string, lines []string, ds *dss.DesignSystem, report *LintReport) {
	// Build set of allowed colors from design system
	allowedColors := make(map[string]bool)
	for _, c := range ds.Foundations.Colors {
		allowedColors[strings.ToLower(c.Value)] = true
	}

	// Find hardcoded hex colors outside of variable definitions
	hexPattern := regexp.MustCompile(`:\s*(#[0-9a-fA-F]{3,8})`)

	for i, line := range lines {
		// Skip if it's a variable definition line
		if strings.Contains(line, "--color-") || strings.Contains(line, "--plexus-") || strings.Contains(line, "--md-") {
			continue
		}

		matches := hexPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				color := strings.ToLower(match[1])
				if !allowedColors[color] {
					report.Violations = append(report.Violations, Violation{
						Line:     i + 1,
						Rule:     "no-hardcoded-colors",
						Severity: "warning",
						Message:  fmt.Sprintf("Hardcoded color '%s' not in design system", match[1]),
						Found:    match[1],
					})
				}
			}
		}
	}
}

func checkAntiPatternsInTSX(content string, _ []string, _ *dss.DesignSystem, report *LintReport) {
	// Multiple primary buttons
	primaryBtnPattern := regexp.MustCompile(`<Button[^>]*variant=["'](?:default|primary)["']`)
	matches := primaryBtnPattern.FindAllString(content, -1)
	if len(matches) > 1 {
		report.Violations = append(report.Violations, Violation{
			Rule:     "single-primary-button",
			Severity: "warning",
			Message:  fmt.Sprintf("Found %d primary buttons - design system recommends one per view", len(matches)),
		})
	}

	// Nested cards
	nestedCardPattern := regexp.MustCompile(`<Card[^>]*>[\s\S]*?<Card`)
	if nestedCardPattern.MatchString(content) {
		report.Violations = append(report.Violations, Violation{
			Rule:     "no-nested-cards",
			Severity: "warning",
			Message:  "Nested cards detected - design system recommends avoiding nesting",
		})
	}
}

func colorsEqual(a, b string) bool {
	// Normalize colors for comparison
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	// Handle shorthand hex (#fff vs #ffffff)
	if len(a) == 4 && len(b) == 7 {
		a = expandHex(a)
	} else if len(b) == 4 && len(a) == 7 {
		b = expandHex(b)
	}

	return a == b
}

func expandHex(short string) string {
	if len(short) != 4 || short[0] != '#' {
		return short
	}
	return fmt.Sprintf("#%c%c%c%c%c%c", short[1], short[1], short[2], short[2], short[3], short[3])
}

func findLineNumber(lines []string, needle string) int {
	for i, line := range lines {
		if strings.Contains(line, needle) {
			return i + 1
		}
	}
	return 0
}

func printLintReport(report *LintReport) {
	fmt.Printf("Linting: %s (format: %s)\n\n", report.File, report.Format)

	if len(report.Violations) == 0 {
		fmt.Println("✓ No violations found")
		return
	}

	// Group by severity
	var errors, warnings, infos []Violation
	for _, v := range report.Violations {
		switch v.Severity {
		case "error":
			errors = append(errors, v)
		case "warning":
			warnings = append(warnings, v)
		case "info":
			infos = append(infos, v)
		}
	}

	if len(errors) > 0 {
		fmt.Printf("✗ Errors (%d):\n", len(errors))
		for _, v := range errors {
			printViolation(v)
		}
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Printf("⚠ Warnings (%d):\n", len(warnings))
		for _, v := range warnings {
			printViolation(v)
		}
		fmt.Println()
	}

	if len(infos) > 0 {
		fmt.Printf("ℹ Info (%d):\n", len(infos))
		for _, v := range infos {
			printViolation(v)
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("Summary: %d errors, %d warnings, %d passed\n",
		report.Summary.Errors, report.Summary.Warnings, report.Summary.Passed)
}

// normalizeTokenID converts token IDs to CSS-friendly format.
func normalizeTokenID(id string) string {
	return strings.ToLower(strings.ReplaceAll(id, "_", "-"))
}

func printViolation(v Violation) {
	loc := ""
	if v.Line > 0 {
		loc = fmt.Sprintf(":%d", v.Line)
		if v.Column > 0 {
			loc = fmt.Sprintf(":%d:%d", v.Line, v.Column)
		}
	}

	fmt.Printf("  [%s]%s %s\n", v.Rule, loc, v.Message)
	if v.Expected != "" && v.Found != "" {
		fmt.Printf("    expected: %s\n", v.Expected)
		fmt.Printf("    found:    %s\n", v.Found)
	} else if v.Found != "" {
		fmt.Printf("    found: %s\n", v.Found)
	}
}
