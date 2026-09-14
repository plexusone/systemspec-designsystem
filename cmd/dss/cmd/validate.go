package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var (
	jsonOutput bool
)

// ComplianceReport holds validation results (backward compatible format).
type ComplianceReport struct {
	Passed   []string          `json:"passed"`
	Warnings []ComplianceIssue `json:"warnings"`
	Errors   []ComplianceIssue `json:"errors"`
}

// ComplianceIssue represents a single compliance violation.
type ComplianceIssue struct {
	Component string `json:"component,omitempty"`
	File      string `json:"file"`
	Line      int    `json:"line,omitempty"`
	Rule      string `json:"rule"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
}

var validateCmd = &cobra.Command{
	Use:   "validate [components-dir]",
	Short: "Validate component implementations against spec",
	Long: `Validate that React/TypeScript component implementations comply with the design system spec.

Checks for:
  - Hardcoded colors (should use CSS variables)
  - Non-standard spacing values
  - Accessibility issues (missing alt, aria-label)
  - Anti-patterns defined in the spec
  - Variant values matching component specs

Examples:
  dss validate ./src/components
  dss validate -d ./design-system ./web/src/components
  dss validate --json ./src/components`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&jsonOutput, "json", false, "output results as JSON")

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	componentsDir := args[0]
	specDir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(specDir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	if !jsonOutput {
		fmt.Printf("Validating against: %s v%s\n", ds.Meta.Name, ds.Meta.Version)
		fmt.Printf("Components directory: %s\n\n", componentsDir)
	}

	// Create service and validate
	service := dss.NewService(ds)
	result, err := service.ValidateDirectory(context.Background(), componentsDir, nil)
	if err != nil {
		return fmt.Errorf("validating directory: %w", err)
	}

	// Convert to ComplianceReport for backward compatibility
	report := convertToComplianceReport(result)

	// Output report
	if jsonOutput {
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		printReport(report)
	}

	// Exit with error if there are errors
	if len(report.Errors) > 0 {
		return fmt.Errorf("validation found %d errors", len(report.Errors))
	}

	return nil
}

// convertToComplianceReport converts ValidationResult to ComplianceReport.
func convertToComplianceReport(result *dss.ValidationResult) *ComplianceReport {
	report := &ComplianceReport{
		Passed:   []string{},
		Warnings: []ComplianceIssue{},
		Errors:   []ComplianceIssue{},
	}

	for _, v := range result.Violations {
		issue := ComplianceIssue{
			Component: v.Component,
			File:      v.File,
			Line:      v.Line,
			Rule:      v.Rule,
			Message:   v.Message,
			Severity:  v.Severity,
		}

		switch v.Severity {
		case "error":
			report.Errors = append(report.Errors, issue)
		case "warning":
			report.Warnings = append(report.Warnings, issue)
		default:
			report.Warnings = append(report.Warnings, issue)
		}
	}

	// Add passed count based on files checked
	if result.Files > 0 && len(report.Errors) == 0 {
		report.Passed = append(report.Passed, fmt.Sprintf("%d files validated", result.Files))
	}

	return report
}

func printReport(report *ComplianceReport) {
	fmt.Print("=== Compliance Report ===\n\n")

	if len(report.Passed) > 0 {
		fmt.Println("✓ Passed:")
		for _, p := range report.Passed {
			fmt.Printf("  %s\n", p)
		}
		fmt.Println()
	}

	if len(report.Warnings) > 0 {
		fmt.Printf("⚠ Warnings (%d):\n", len(report.Warnings))
		for _, w := range report.Warnings {
			loc := w.File
			if w.Line > 0 {
				loc = fmt.Sprintf("%s:%d", w.File, w.Line)
			}
			fmt.Printf("  [%s] %s\n    %s\n", w.Rule, loc, w.Message)
		}
		fmt.Println()
	}

	if len(report.Errors) > 0 {
		fmt.Printf("✗ Errors (%d):\n", len(report.Errors))
		for _, e := range report.Errors {
			loc := e.File
			if e.Line > 0 {
				loc = fmt.Sprintf("%s:%d", e.File, e.Line)
			}
			fmt.Printf("  [%s] %s\n    %s\n", e.Rule, loc, e.Message)
		}
		fmt.Println()
	}

	// Summary
	total := len(report.Passed) + len(report.Warnings) + len(report.Errors)
	fmt.Printf("Summary: %d checks (%d passed, %d warnings, %d errors)\n",
		total, len(report.Passed), len(report.Warnings), len(report.Errors))

	if len(report.Errors) == 0 && len(report.Warnings) == 0 {
		fmt.Println("\n✓ All components comply with design system!")
	} else if len(report.Errors) == 0 {
		fmt.Println("\n⚠ Components have warnings but no blocking errors")
	} else {
		fmt.Println("\n✗ Components have compliance errors that should be fixed")
	}
}
