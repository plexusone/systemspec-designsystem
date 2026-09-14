package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/plexusone/systemspec-designsystem/sdk/go/visual"
	"github.com/spf13/cobra"
)

var (
	// visual-test flags
	visualTestsDir     string
	visualBaselinesDir string
	visualOutputDir    string
	visualParallel     int
	visualThreshold    float64
	visualViewports    []string
	visualTestIDs      []string
	visualBaseline     string
	visualJSONOutput   bool
)

var visualCmd = &cobra.Command{
	Use:   "visual",
	Short: "Visual regression testing commands",
	Long:  "Commands for visual regression testing of design system components.",
}

var visualTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run visual regression tests",
	Long: `Run visual regression tests comparing current screenshots against baselines.

Examples:
  dss visual test                                # Run all tests
  dss visual test --tests button,card            # Run specific tests
  dss visual test --viewports desktop,mobile     # Run specific viewports
  dss visual test --threshold 0.01               # Set diff threshold (1%)
  dss visual test --json                         # Output results as JSON`,
	RunE: runVisualTest,
}

var visualBaselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Manage visual test baselines",
	Long:  "Commands for generating, updating, and managing baseline images.",
}

var visualBaselineGenerateCmd = &cobra.Command{
	Use:   "generate <version>",
	Short: "Generate new baselines",
	Long: `Generate baseline images for a new version.

Examples:
  dss visual baseline generate v1.0.0
  dss visual baseline generate v1.1.0 --tests button,card`,
	Args: cobra.ExactArgs(1),
	RunE: runVisualBaselineGenerate,
}

var visualBaselineUpdateCmd = &cobra.Command{
	Use:   "update <version>",
	Short: "Update specific test baselines",
	Long: `Update baseline images for specific tests.

Examples:
  dss visual baseline update v1.0.0 --tests button
  dss visual baseline update v1.0.0 --tests button,card`,
	Args: cobra.ExactArgs(1),
	RunE: runVisualBaselineUpdate,
}

var visualBaselineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available baseline versions",
	Long:  "List all available baseline versions.",
	RunE:  runVisualBaselineList,
}

var visualBaselinePruneCmd = &cobra.Command{
	Use:   "prune <version>",
	Short: "Remove a baseline version",
	Long: `Remove a baseline version and all its images.

Examples:
  dss visual baseline prune v0.9.0`,
	Args: cobra.ExactArgs(1),
	RunE: runVisualBaselinePrune,
}

func init() {
	// Add visual command to root
	rootCmd.AddCommand(visualCmd)

	// Add subcommands to visual
	visualCmd.AddCommand(visualTestCmd)
	visualCmd.AddCommand(visualBaselineCmd)

	// Add subcommands to visual baseline
	visualBaselineCmd.AddCommand(visualBaselineGenerateCmd)
	visualBaselineCmd.AddCommand(visualBaselineUpdateCmd)
	visualBaselineCmd.AddCommand(visualBaselineListCmd)
	visualBaselineCmd.AddCommand(visualBaselinePruneCmd)

	// visual test flags
	visualTestCmd.Flags().StringVar(&visualTestsDir, "tests-dir", "", "directory containing visual-tests.yaml (default: <spec-dir>/visual)")
	visualTestCmd.Flags().StringVar(&visualBaselinesDir, "baselines-dir", "", "directory containing baselines (default: <spec-dir>/visual/baselines)")
	visualTestCmd.Flags().StringVar(&visualOutputDir, "output-dir", "", "directory for test results (default: <spec-dir>/visual/results)")
	visualTestCmd.Flags().IntVar(&visualParallel, "parallel", 4, "number of parallel workers")
	visualTestCmd.Flags().Float64Var(&visualThreshold, "threshold", 0, "override diff threshold (0 = use test defaults)")
	visualTestCmd.Flags().StringSliceVar(&visualViewports, "viewports", nil, "specific viewports to test")
	visualTestCmd.Flags().StringSliceVar(&visualTestIDs, "tests", nil, "specific tests to run")
	visualTestCmd.Flags().StringVar(&visualBaseline, "baseline", "latest", "baseline version to compare against")
	visualTestCmd.Flags().BoolVar(&visualJSONOutput, "json", false, "output results as JSON")

	// visual baseline generate flags
	visualBaselineGenerateCmd.Flags().StringVar(&visualTestsDir, "tests-dir", "", "directory containing visual-tests.yaml")
	visualBaselineGenerateCmd.Flags().StringVar(&visualBaselinesDir, "baselines-dir", "", "directory for baselines")
	visualBaselineGenerateCmd.Flags().StringSliceVar(&visualTestIDs, "tests", nil, "specific tests to generate (empty = all)")

	// visual baseline update flags
	visualBaselineUpdateCmd.Flags().StringVar(&visualTestsDir, "tests-dir", "", "directory containing visual-tests.yaml")
	visualBaselineUpdateCmd.Flags().StringVar(&visualBaselinesDir, "baselines-dir", "", "directory for baselines")
	visualBaselineUpdateCmd.Flags().StringSliceVar(&visualTestIDs, "tests", nil, "specific tests to update (required)")
	if err := visualBaselineUpdateCmd.MarkFlagRequired("tests"); err != nil {
		panic(err)
	}

	// visual baseline list flags
	visualBaselineListCmd.Flags().StringVar(&visualBaselinesDir, "baselines-dir", "", "directory containing baselines")
	visualBaselineListCmd.Flags().BoolVar(&visualJSONOutput, "json", false, "output as JSON")

	// visual baseline prune flags
	visualBaselinePruneCmd.Flags().StringVar(&visualBaselinesDir, "baselines-dir", "", "directory containing baselines")
}

func runVisualTest(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Resolve directories
	specDir := getSpecDir()
	testsDir := visualTestsDir
	if testsDir == "" {
		testsDir = filepath.Join(specDir, "visual")
	}
	baselinesDir := visualBaselinesDir
	if baselinesDir == "" {
		baselinesDir = filepath.Join(testsDir, "baselines")
	}
	outputDir := visualOutputDir
	if outputDir == "" {
		outputDir = filepath.Join(testsDir, "results")
	}

	// Create service
	service := visual.NewService(visual.ServiceOptions{
		TestsDir:     testsDir,
		BaselinesDir: baselinesDir,
		OutputDir:    outputDir,
		Parallel:     visualParallel,
	})

	// Run tests
	report, err := service.RunTests(ctx, &visual.TestOptions{
		BaselineVersion: visualBaseline,
		TestIDs:         visualTestIDs,
		Viewports:       visualViewports,
		Threshold:       visualThreshold,
	})
	if err != nil {
		return fmt.Errorf("failed to run tests: %w", err)
	}

	// Output results
	if visualJSONOutput {
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		printVisualTestReport(report)
	}

	// Exit with error if tests failed
	if report.Summary.Failed > 0 || report.Summary.Errors > 0 {
		os.Exit(1)
	}

	return nil
}

func printVisualTestReport(report *visual.VisualTestReport) {
	fmt.Printf("\nVisual Regression Test Results\n")
	fmt.Printf("==============================\n")
	fmt.Printf("Baseline: %s\n", report.BaselineVersion)
	fmt.Printf("Duration: %s\n", report.Duration)
	fmt.Printf("\n")

	// Summary
	fmt.Printf("Summary:\n")
	fmt.Printf("  Total:   %d\n", report.Summary.Total)
	fmt.Printf("  Passed:  %d\n", report.Summary.Passed)
	fmt.Printf("  Failed:  %d\n", report.Summary.Failed)
	fmt.Printf("  Skipped: %d\n", report.Summary.Skipped)
	fmt.Printf("  Errors:  %d\n", report.Summary.Errors)
	fmt.Printf("\n")

	// Failed tests
	var failed []visual.VisualTestResult
	for _, r := range report.Results {
		if r.Status == visual.TestStatusFailed {
			failed = append(failed, r)
		}
	}

	if len(failed) > 0 {
		fmt.Printf("Failed Tests:\n")
		for _, r := range failed {
			fmt.Printf("  - %s @ %s (diff: %.2f%%)\n", r.TestID, r.Viewport, r.DiffPercent*100)
			if r.DiffPath != "" {
				fmt.Printf("    Diff: %s\n", r.DiffPath)
			}
		}
		fmt.Printf("\n")
	}

	// Errors
	var errors []visual.VisualTestResult
	for _, r := range report.Results {
		if r.Status == visual.TestStatusError {
			errors = append(errors, r)
		}
	}

	if len(errors) > 0 {
		fmt.Printf("Errors:\n")
		for _, r := range errors {
			fmt.Printf("  - %s @ %s: %s\n", r.TestID, r.Viewport, r.Error)
		}
		fmt.Printf("\n")
	}
}

func runVisualBaselineGenerate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	version := args[0]

	// Resolve directories
	specDir := getSpecDir()
	testsDir := visualTestsDir
	if testsDir == "" {
		testsDir = filepath.Join(specDir, "visual")
	}
	baselinesDir := visualBaselinesDir
	if baselinesDir == "" {
		baselinesDir = filepath.Join(testsDir, "baselines")
	}

	// Create service
	service := visual.NewService(visual.ServiceOptions{
		TestsDir:     testsDir,
		BaselinesDir: baselinesDir,
	})

	// Generate baseline
	result, err := service.GenerateBaseline(ctx, version)
	if err != nil {
		return fmt.Errorf("failed to generate baseline: %w", err)
	}

	fmt.Printf("Generated baseline %s\n", result.Version)
	fmt.Printf("  Tests: %d\n", result.TestCount)
	fmt.Printf("  Path: %s\n", result.Path)
	if result.Errors > 0 {
		fmt.Printf("  Errors: %d\n", result.Errors)
	}

	return nil
}

func runVisualBaselineUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	version := args[0]

	// Resolve directories
	specDir := getSpecDir()
	testsDir := visualTestsDir
	if testsDir == "" {
		testsDir = filepath.Join(specDir, "visual")
	}
	baselinesDir := visualBaselinesDir
	if baselinesDir == "" {
		baselinesDir = filepath.Join(testsDir, "baselines")
	}

	// Create service
	service := visual.NewService(visual.ServiceOptions{
		TestsDir:     testsDir,
		BaselinesDir: baselinesDir,
	})

	// Update baseline
	result, err := service.UpdateBaseline(ctx, version, visualTestIDs)
	if err != nil {
		return fmt.Errorf("failed to update baseline: %w", err)
	}

	fmt.Printf("Updated baseline %s\n", result.Version)
	fmt.Printf("  Updated tests: %d\n", len(visualTestIDs))
	fmt.Printf("  Total tests: %d\n", result.TestCount)
	if result.Errors > 0 {
		fmt.Printf("  Errors: %d\n", result.Errors)
	}

	return nil
}

func runVisualBaselineList(cmd *cobra.Command, args []string) error {
	// Resolve directories
	specDir := getSpecDir()
	baselinesDir := visualBaselinesDir
	if baselinesDir == "" {
		baselinesDir = filepath.Join(specDir, "visual", "baselines")
	}

	// Create service
	service := visual.NewService(visual.ServiceOptions{
		BaselinesDir: baselinesDir,
	})

	// List versions
	versions, err := service.ListBaselineVersions()
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	if visualJSONOutput {
		data, err := json.MarshalIndent(versions, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		if len(versions) == 0 {
			fmt.Println("No baseline versions found.")
		} else {
			fmt.Println("Available baseline versions:")
			for _, v := range versions {
				manifest, _ := service.GetBaselineManifest(v)
				if manifest != nil {
					fmt.Printf("  %s (%d tests, %s)\n", v, manifest.TestCount, manifest.CreatedAt.Format("2006-01-02"))
				} else {
					fmt.Printf("  %s\n", v)
				}
			}
		}
	}

	return nil
}

func runVisualBaselinePrune(cmd *cobra.Command, args []string) error {
	version := args[0]

	// Resolve directories
	specDir := getSpecDir()
	baselinesDir := visualBaselinesDir
	if baselinesDir == "" {
		baselinesDir = filepath.Join(specDir, "visual", "baselines")
	}

	// Create manager directly for prune
	manager := visual.NewBaselineManager(baselinesDir)

	// Check if version exists
	if !manager.VersionExists(version) {
		return fmt.Errorf("baseline version %s not found", version)
	}

	// Prune
	if err := manager.PruneVersion(version); err != nil {
		return fmt.Errorf("failed to prune version: %w", err)
	}

	fmt.Printf("Removed baseline version %s\n", version)
	return nil
}
