package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
	"github.com/spf13/cobra"
)

var (
	lintSpecRules    []string
	lintSpecMinScore int
	lintSpecJSON     bool
	lintSpecVerbose  bool
)

var lintSpecCmd = &cobra.Command{
	Use:   "lint-spec",
	Short: "Lint design system spec for completeness",
	Long: `Lint the design system specification for completeness and best practices.

Checks include:
  - Meta fields (name, version, description)
  - Component definitions (variants, props, LLM context)
  - LLM context completeness (intent, anti-patterns, allowed contexts)
  - Token references and descriptions
  - Cross-reference validation
  - Accessibility requirements
  - Theming contracts

Returns a completeness score (0-100) and coverage metrics.`,
	Example: `  # Lint current directory
  dss lint-spec

  # Lint with minimum score requirement
  dss lint-spec --min-score 80

  # Lint specific rules only
  dss lint-spec --rules component-has-llm-context,llm-has-intent

  # JSON output for CI
  dss lint-spec --json

  # Show all issues including info level
  dss lint-spec --verbose

  # Show available rules
  dss lint-spec rules`,
	RunE: runLintSpec,
}

var listLintSpecRulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List available spec lint rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		rules := dss.AvailableLintRules()

		// Sort rules for consistent output
		ids := make([]string, 0, len(rules))
		for id := range rules {
			ids = append(ids, id)
		}
		sort.Strings(ids)

		fmt.Println("Available Spec Lint Rules:")
		fmt.Println()
		for _, id := range ids {
			fmt.Printf("  %-30s %s\n", id, rules[id])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintSpecCmd)
	lintSpecCmd.AddCommand(listLintSpecRulesCmd)

	lintSpecCmd.Flags().StringSliceVar(&lintSpecRules, "rules", nil, "Specific rules to check (comma-separated)")
	lintSpecCmd.Flags().IntVar(&lintSpecMinScore, "min-score", 0, "Minimum acceptable score (0-100)")
	lintSpecCmd.Flags().BoolVar(&lintSpecJSON, "json", false, "Output as JSON")
	lintSpecCmd.Flags().BoolVarP(&lintSpecVerbose, "verbose", "v", false, "Show all issues including info level")
}

func runLintSpec(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	service := dss.NewService(ds)
	opts := &dss.LintOptions{
		Rules:              lintSpecRules,
		MinScore:           lintSpecMinScore,
		IncludeSuggestions: true,
	}

	result := service.LintSpec(cmd.Context(), opts)

	if lintSpecJSON {
		return outputLintSpecJSON(result)
	}

	return outputLintSpecHuman(result, lintSpecVerbose, lintSpecMinScore)
}

func outputLintSpecJSON(result *dss.LintResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputLintSpecHuman(result *dss.LintResult, verbose bool, minScore int) error {
	// Score header
	scoreColor := "\033[32m" // green
	if result.Score < 50 {
		scoreColor = "\033[31m" // red
	} else if result.Score < 80 {
		scoreColor = "\033[33m" // yellow
	}
	fmt.Printf("\nSpec Completeness Score: %s%d/100\033[0m\n\n", scoreColor, result.Score)

	// Coverage metrics
	fmt.Println("Coverage:")
	fmt.Printf("  Components with LLM context: %.0f%%\n", result.Coverage.ComponentsWithLLMContext)
	fmt.Printf("  Components with variants:    %.0f%%\n", result.Coverage.ComponentsWithVariants)
	fmt.Printf("  Components with props:       %.0f%%\n", result.Coverage.ComponentsWithProps)
	fmt.Printf("  Tokens with descriptions:    %.0f%%\n", result.Coverage.TokensWithDescriptions)
	fmt.Printf("  Tokens referenced:           %.0f%%\n", result.Coverage.TokensReferenced)
	fmt.Println()

	// Issues summary
	fmt.Printf("Issues: %d errors, %d warnings, %d info\n\n",
		result.Summary.Errors, result.Summary.Warnings, result.Summary.Infos)

	// Issue details
	if len(result.Issues) > 0 {
		// Group by severity
		errors := filterSpecBySeverity(result.Issues, "error")
		warnings := filterSpecBySeverity(result.Issues, "warning")
		infos := filterSpecBySeverity(result.Issues, "info")

		if len(errors) > 0 {
			fmt.Println("\033[31mErrors:\033[0m")
			for _, issue := range errors {
				printSpecIssue(issue)
			}
			fmt.Println()
		}

		if len(warnings) > 0 {
			fmt.Println("\033[33mWarnings:\033[0m")
			for _, issue := range warnings {
				printSpecIssue(issue)
			}
			fmt.Println()
		}

		if verbose && len(infos) > 0 {
			fmt.Println("\033[36mInfo:\033[0m")
			for _, issue := range infos {
				printSpecIssue(issue)
			}
			fmt.Println()
		}
	}

	// Check minimum score
	if minScore > 0 && result.Score < minScore {
		fmt.Printf("\033[31m✗ Score %d is below minimum required score of %d\033[0m\n", result.Score, minScore)
		os.Exit(1)
	}

	// Final status
	if result.Summary.Errors > 0 {
		fmt.Println("\033[31m✗ Spec lint failed with errors\033[0m")
		os.Exit(1)
	}

	if result.Score >= 80 {
		fmt.Println("\033[32m✓ Spec is agent-ready\033[0m")
	} else {
		fmt.Println("\033[33m⚠ Spec could be improved for better agent support\033[0m")
	}

	return nil
}

func filterSpecBySeverity(issues []dss.LintIssue, severity string) []dss.LintIssue {
	var filtered []dss.LintIssue
	for _, issue := range issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func printSpecIssue(issue dss.LintIssue) {
	component := ""
	if issue.Component != "" {
		component = fmt.Sprintf(" [%s]", issue.Component)
	}
	fmt.Printf("  %s%s: %s\n", issue.Path, component, issue.Message)
	if issue.Suggestion != "" {
		fmt.Printf("    → %s\n", issue.Suggestion)
	}
}
