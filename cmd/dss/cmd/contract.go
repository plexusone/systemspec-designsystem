package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var (
	contractJSON bool
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Display and validate theming contracts",
	Long: `Display and validate component theming contracts.

Theming contracts define the customizable CSS tokens for a component,
enabling consistent theming across different applications.

Examples:
  dss contract show button
  dss contract validate
  dss contract show button --json`,
}

var contractShowCmd = &cobra.Command{
	Use:   "show <component>",
	Short: "Display a component's theming contract",
	Long: `Display the theming contract for a specific component.

Shows all customizable tokens, their CSS properties, semantic meanings,
and default values for light and dark modes.

Examples:
  dss contract show button
  dss contract show button --json
  dss contract show button -d ./design-system`,
	Args: cobra.ExactArgs(1),
	RunE: runContractShow,
}

var contractValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all theming contracts",
	Long: `Validate all component theming contracts for completeness.

Checks:
- All tokens have defaultLight and defaultDark values
- CSS properties follow the prefix convention
- Semantic values are from the allowed set
- No duplicate token IDs

Examples:
  dss contract validate
  dss contract validate --json
  dss contract validate -d ./design-system`,
	RunE: runContractValidate,
}

func init() {
	rootCmd.AddCommand(contractCmd)
	contractCmd.AddCommand(contractShowCmd)
	contractCmd.AddCommand(contractValidateCmd)

	contractCmd.PersistentFlags().BoolVar(&contractJSON, "json", false, "output as JSON instead of table")
}

func runContractShow(cmd *cobra.Command, args []string) error {
	componentID := args[0]
	dir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	// Find component
	var component *dss.Component
	for i := range ds.Components {
		if ds.Components[i].ID == componentID {
			component = &ds.Components[i]
			break
		}
	}

	if component == nil {
		return fmt.Errorf("component '%s' not found", componentID)
	}

	if component.ThemingContract == nil {
		fmt.Printf("Component '%s' has no theming contract defined.\n", componentID)
		return nil
	}

	contract := component.ThemingContract

	if contractJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(contract)
	}

	// Table output
	fmt.Printf("Theming Contract: %s\n", component.Name)
	fmt.Printf("Prefix: %s\n", contract.Prefix)
	if contract.Description != "" {
		fmt.Printf("Description: %s\n", contract.Description)
	}
	fmt.Println()

	if len(contract.Tokens) == 0 {
		fmt.Println("No tokens defined.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCSS PROPERTY\tSEMANTIC\tLIGHT DEFAULT\tDARK DEFAULT")
	fmt.Fprintln(w, "--\t------------\t--------\t-------------\t------------")

	for _, token := range contract.Tokens {
		semantic := token.Semantic
		if semantic == "" {
			semantic = "-"
		}
		defaultLight := token.DefaultLight
		if defaultLight == "" {
			defaultLight = "-"
		}
		defaultDark := token.DefaultDark
		if defaultDark == "" {
			defaultDark = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			token.ID, token.CSSProperty, semantic, defaultLight, defaultDark)
	}
	w.Flush()

	return nil
}

func runContractValidate(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	// Validate all contracts
	reports := dss.ValidateAllContracts(ds)

	if len(reports) == 0 {
		fmt.Println("No components with theming contracts found.")
		return nil
	}

	if contractJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(reports)
	}

	// Table output
	totalErrors := dss.TotalErrors(reports)
	totalWarnings := dss.TotalWarnings(reports)
	allPassed := totalErrors == 0

	for _, report := range reports {
		if report.Passed && len(report.Warnings) == 0 {
			fmt.Printf("✓ %s: OK\n", report.ComponentID)
		} else if report.Passed {
			fmt.Printf("⚠ %s: %d warnings\n", report.ComponentID, len(report.Warnings))
		} else {
			fmt.Printf("✗ %s: %d errors, %d warnings\n",
				report.ComponentID, len(report.Errors), len(report.Warnings))
		}

		// Print errors
		for _, issue := range report.Errors {
			if issue.TokenID != "" {
				fmt.Printf("    ERROR [%s] %s: %s\n", issue.TokenID, issue.Field, issue.Message)
			} else {
				fmt.Printf("    ERROR %s: %s\n", issue.Field, issue.Message)
			}
		}

		// Print warnings
		for _, issue := range report.Warnings {
			if issue.TokenID != "" {
				fmt.Printf("    WARN  [%s] %s: %s\n", issue.TokenID, issue.Field, issue.Message)
			} else {
				fmt.Printf("    WARN  %s: %s\n", issue.Field, issue.Message)
			}
		}
	}

	fmt.Println()
	if allPassed {
		fmt.Printf("Validation passed: %d contracts validated, %d warnings\n",
			len(reports), totalWarnings)
		return nil
	}

	return fmt.Errorf("validation failed: %d errors, %d warnings", totalErrors, totalWarnings)
}
