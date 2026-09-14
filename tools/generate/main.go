// Package main generates JSON Schema files from Go types.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Default output directory
	outputDir := "schema"
	if len(os.Args) > 1 {
		outputDir = filepath.Clean(os.Args[1])
	}

	// Ensure output directories exist
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "foundations"),
		filepath.Join(outputDir, "components"),
	}
	for _, dir := range dirs {
		// #nosec G703 -- CLI tool with user-provided output path; path is cleaned above
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	// Generate schemas
	schemas := []struct {
		name   string
		path   string
		schema any
	}{
		{
			name:   "DesignSystem",
			path:   filepath.Join(outputDir, "design-system.schema.json"),
			schema: dss.GenerateSchema(),
		},
		{
			name:   "Meta",
			path:   filepath.Join(outputDir, "meta.schema.json"),
			schema: dss.GenerateMetaSchema(),
		},
		{
			name:   "Foundations",
			path:   filepath.Join(outputDir, "foundations", "foundations.schema.json"),
			schema: dss.GenerateFoundationsSchema(),
		},
		{
			name:   "Component",
			path:   filepath.Join(outputDir, "components", "component.schema.json"),
			schema: dss.GenerateComponentSchema(),
		},
		{
			name:   "Pattern",
			path:   filepath.Join(outputDir, "patterns.schema.json"),
			schema: dss.GeneratePatternSchema(),
		},
		{
			name:   "Template",
			path:   filepath.Join(outputDir, "templates.schema.json"),
			schema: dss.GenerateTemplateSchema(),
		},
		{
			name:   "Accessibility",
			path:   filepath.Join(outputDir, "accessibility.schema.json"),
			schema: dss.GenerateAccessibilitySchema(),
		},
		{
			name:   "Governance",
			path:   filepath.Join(outputDir, "governance.schema.json"),
			schema: dss.GenerateGovernanceSchema(),
		},
	}

	for _, s := range schemas {
		data, err := json.MarshalIndent(s.schema, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal %s schema: %w", s.name, err)
		}

		if err := os.WriteFile(s.path, data, 0600); err != nil {
			return fmt.Errorf("write %s schema: %w", s.name, err)
		}

		fmt.Printf("Generated %s\n", s.path)
	}

	fmt.Println("\nSchema generation complete!")
	return nil
}
