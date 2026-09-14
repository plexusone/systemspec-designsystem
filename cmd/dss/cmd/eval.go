package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/plexusone/structured-evaluation/rubric"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var (
	evalJSONOutput bool
	evalMinScore   int
	evalCategories []string
	evalVerbose    bool
)

var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate design system spec completeness and quality",
	Long: `Evaluate a design system specification for completeness, agent-readiness,
accessibility, and documentation quality.

Categories evaluated (with weights):
  - completeness (25%):    Required fields, foundations, components
  - agent-readiness (30%): LLM context, anti-patterns, examples
  - accessibility (25%):   WCAG compliance, keyboard, screen reader
  - documentation (20%):   Descriptions, usage guidance

Score scale (1-5):
  5 - Excellent: Exceeds expectations
  4 - Good: Meets expectations well
  3 - Acceptable: Minimum requirements met
  2 - Major Revisions: Significant work needed
  1 - Unacceptable: Does not meet requirements

Examples:
  dss eval
  dss eval --spec ./specs/v3 --json
  dss eval --min-score 4
  dss eval --categories completeness,accessibility`,
	RunE: runEval,
}

func init() {
	evalCmd.Flags().BoolVar(&evalJSONOutput, "json", false, "output results as JSON")
	evalCmd.Flags().IntVar(&evalMinScore, "min-score", 0, "minimum acceptable score 1-5 (exit 1 if below)")
	evalCmd.Flags().StringSliceVar(&evalCategories, "categories", nil, "categories to evaluate (default: all)")
	evalCmd.Flags().BoolVarP(&evalVerbose, "verbose", "v", false, "show detailed issue information")

	rootCmd.AddCommand(evalCmd)
}

func runEval(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()
	ctx := context.Background()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	// Create service
	service := dss.NewService(ds)

	// Build options
	opts := dss.DefaultEvalOptions()
	opts.Verbose = evalVerbose

	if len(evalCategories) > 0 {
		opts.Categories = evalCategories
	}

	// Run evaluation
	result, err := service.Evaluate(ctx, opts)
	if err != nil {
		return fmt.Errorf("evaluating spec: %w", err)
	}

	// Output results
	if evalJSONOutput {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		printEvalResult(result, evalVerbose)
	}

	// Check minimum score
	if evalMinScore > 0 && int(result.IntScore) < evalMinScore {
		return fmt.Errorf("score %d is below minimum %d", result.IntScore, evalMinScore)
	}

	return nil
}

func printEvalResult(result *rubric.Rubric, verbose bool) {
	fmt.Printf("=== Design System Evaluation ===\n\n")
	fmt.Printf("System: %s v%s\n", result.Metadata.Document, result.Metadata.DocumentVersion)
	fmt.Printf("Evaluated: %s\n\n", result.Metadata.GeneratedAt.Format("2006-01-02 15:04:05"))

	// Overall score
	fmt.Printf("Overall Score: %d/5 (%s)\n", result.IntScore, result.IntScore.String())
	fmt.Printf("Decision: %s\n\n", result.OverallDecision)

	// Category breakdown
	fmt.Println("Category Breakdown:")
	fmt.Println("-------------------")
	for _, cat := range result.Categories {
		bar := intScoreBar(int(cat.IntScore))
		fmt.Printf("  %-18s %s %d/5 (%s)\n", cat.Category, bar, cat.IntScore, cat.Score)
	}
	fmt.Println()

	// Findings summary
	if len(result.Findings) > 0 {
		counts := rubric.CountFindings(result.Findings)

		fmt.Printf("Findings: %d total", counts.Total)
		if counts.Critical > 0 || counts.High > 0 {
			fmt.Printf(" (%d critical, %d high, %d medium, %d low, %d info)",
				counts.Critical, counts.High, counts.Medium, counts.Low, counts.Info)
		}
		fmt.Println()

		if verbose {
			fmt.Println()
			printEvalFindings(result.Findings)
		} else if counts.Critical > 0 || counts.High > 0 || counts.Medium > 0 {
			fmt.Println("(use --verbose to see details)")
		}
	} else {
		fmt.Println("No findings!")
	}

	// Summary
	if result.Summary != "" {
		fmt.Printf("\n%s\n", result.Summary)
	}

	// Next steps
	if len(result.NextSteps.Immediate) > 0 {
		fmt.Println("\nImmediate Actions:")
		for _, action := range result.NextSteps.Immediate {
			fmt.Printf("  - %s\n", action.Action)
		}
	}

	if verbose && len(result.NextSteps.Recommended) > 0 {
		fmt.Println("\nRecommended Actions:")
		for _, action := range result.NextSteps.Recommended {
			fmt.Printf("  - %s\n", action.Action)
		}
	}

	// Score interpretation
	fmt.Println()
	switch result.IntScore {
	case rubric.ScoreExcellent:
		fmt.Println("Excellent! This spec is well-prepared for agent-driven development.")
	case rubric.ScoreGood:
		fmt.Println("Good spec quality. A few improvements would enhance agent readiness.")
	case rubric.ScoreAcceptable:
		fmt.Println("Adequate spec. Consider addressing findings for better agent support.")
	case rubric.ScoreMajorRevisions:
		fmt.Println("Needs improvement. Multiple areas require attention.")
	case rubric.ScoreUnacceptable:
		fmt.Println("Incomplete spec. Significant work needed before agent use.")
	}
}

func printEvalFindings(findings []rubric.Finding) {
	fmt.Println("Findings:")
	fmt.Println("---------")

	// Group by severity
	var critical, high, medium, low, info []rubric.Finding
	for _, finding := range findings {
		switch finding.Severity {
		case rubric.SeverityCritical:
			critical = append(critical, finding)
		case rubric.SeverityHigh:
			high = append(high, finding)
		case rubric.SeverityMedium:
			medium = append(medium, finding)
		case rubric.SeverityLow:
			low = append(low, finding)
		case rubric.SeverityInfo:
			info = append(info, finding)
		}
	}

	// Print by severity (critical/high first)
	if len(critical) > 0 {
		fmt.Println("\nCritical:")
		for _, f := range critical {
			printFinding(f)
		}
	}

	if len(high) > 0 {
		fmt.Println("\nHigh:")
		for _, f := range high {
			printFinding(f)
		}
	}

	if len(medium) > 0 {
		fmt.Println("\nMedium:")
		for _, f := range medium {
			printFinding(f)
		}
	}

	if len(low) > 0 {
		fmt.Println("\nLow:")
		for _, f := range low {
			printFinding(f)
		}
	}

	if len(info) > 0 {
		fmt.Println("\nInfo:")
		for _, f := range info {
			printFinding(f)
		}
	}
}

func printFinding(f rubric.Finding) {
	fmt.Printf("  [%s] %s\n", f.ID, f.Location)
	fmt.Printf("    %s\n", f.Title)
	if f.Description != "" && f.Description != f.Title {
		fmt.Printf("    %s\n", f.Description)
	}
	if f.Recommendation != "" {
		fmt.Printf("    -> %s\n", f.Recommendation)
	}
}

func intScoreBar(score int) string {
	filled := score
	empty := 5 - score

	bar := "["
	for i := 0; i < filled; i++ {
		bar += "#"
	}
	for i := 0; i < empty; i++ {
		bar += " "
	}
	bar += "]"

	return bar
}
