package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about the design system",
	Long: `Display information about the design system specification.

Shows metadata, foundation counts, component counts, and validation status.

Examples:
  dss info
  dss info -d ./design-system`,
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(_ *cobra.Command, _ []string) error {
	dir := getSpecDir()
	ctx := context.Background()

	// Load design system and create service
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}
	service := dss.NewService(ds)

	// Get metadata via service
	meta := service.GetMeta(ctx)
	fmt.Printf("Design System: %s\n", meta.Name)
	fmt.Printf("Version: %s\n", meta.Version)
	if meta.Description != "" {
		fmt.Printf("Description: %s\n", meta.Description)
	}
	fmt.Println()

	// Principles (access underlying DS for now, could add to service later)
	underlyingDS := service.DesignSystem()
	if len(underlyingDS.Principles) > 0 {
		fmt.Printf("Principles: %d\n", len(underlyingDS.Principles))
		for _, p := range underlyingDS.Principles {
			fmt.Printf("  - %s\n", p.Name)
		}
		fmt.Println()
	}

	// Foundations - show token counts
	fmt.Println("Foundations:")
	printTokenCounts(ctx, service)
	fmt.Println()

	// Components via service
	components := service.ListComponents(ctx)
	if len(components) > 0 {
		fmt.Printf("Components: %d\n", len(components))
		for _, c := range components {
			// Get full component for variant count
			comp, err := service.GetComponent(ctx, c.ID)
			if err == nil && len(comp.Variants) > 0 {
				fmt.Printf("  - %s (%d variants)\n", c.Name, len(comp.Variants))
			} else {
				fmt.Printf("  - %s\n", c.Name)
			}
		}
		fmt.Println()
	}

	// Patterns via service
	patterns := service.ListPatterns(ctx)
	if len(patterns) > 0 {
		fmt.Printf("Patterns: %d\n", len(patterns))
		for _, p := range patterns {
			fmt.Printf("  - %s\n", p.Name)
		}
		fmt.Println()
	}

	// Accessibility
	a11y := service.GetAccessibility(ctx)
	if a11y != nil && a11y.WCAGLevel != "" {
		fmt.Printf("Accessibility Target: WCAG %s\n", a11y.WCAGLevel)
	}

	// Validation
	if err := underlyingDS.Validate(); err != nil {
		fmt.Printf("\n⚠ Validation Issues: %v\n", err)
	} else {
		fmt.Println("✓ Spec validation passed")
	}

	return nil
}

// printTokenCounts displays counts for each token type.
func printTokenCounts(ctx context.Context, service *dss.Service) {
	tokenTypes := []struct {
		name     string
		typeName string
	}{
		{"Colors", "color"},
		{"Spacing", "spacing"},
		{"Typography", "typography"},
		{"Elevation", "elevation"},
		{"Border Radius", "borderRadius"},
		{"Breakpoints", "breakpoint"},
	}

	for _, tt := range tokenTypes {
		tokens, err := service.ListTokens(ctx, tt.typeName)
		if err == nil && len(tokens) > 0 {
			fmt.Printf("  %s: %d tokens\n", tt.name, len(tokens))
		}
	}
}
