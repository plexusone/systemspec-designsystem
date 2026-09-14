package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

var coverageJsonOutput bool

// CoverageReport holds spec completeness results
type CoverageReport struct {
	DesignSystem string             `json:"designSystem"`
	Version      string             `json:"version"`
	Categories   []CategoryCoverage `json:"categories"`
	Summary      CoverageSummary    `json:"summary"`
}

// CategoryCoverage represents coverage for a DSS category
type CategoryCoverage struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Required    bool           `json:"required"`
	Items       []ItemCoverage `json:"items"`
	Score       float64        `json:"score"`  // 0-100
	Status      string         `json:"status"` // complete, partial, missing
}

// ItemCoverage represents coverage for a specific item
type ItemCoverage struct {
	Name     string `json:"name"`
	Present  bool   `json:"present"`
	Count    int    `json:"count,omitempty"`
	Required bool   `json:"required"`
	Notes    string `json:"notes,omitempty"`
}

// CoverageSummary provides overall stats
type CoverageSummary struct {
	TotalCategories    int     `json:"totalCategories"`
	CompleteCategories int     `json:"completeCategories"`
	PartialCategories  int     `json:"partialCategories"`
	MissingCategories  int     `json:"missingCategories"`
	OverallScore       float64 `json:"overallScore"`
	Grade              string  `json:"grade"` // A, B, C, D, F
}

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Check design system spec completeness",
	Long: `Analyze a design system specification and report on its completeness
against the full DSS schema.

Reports coverage for:
  - Meta (metadata, maintainers)
  - Principles (design philosophy)
  - Foundations (colors, typography, spacing, elevation, motion, grid, breakpoints)
  - Components (variants, states, props, slots, constraints)
  - Patterns (multi-component compositions)
  - Templates (page layouts)
  - Content (voice, tone, terminology)
  - Accessibility (WCAG, contrast, keyboard)
  - Governance (versioning, deprecation)

Examples:
  dss coverage
  dss coverage -d ./design-system
  dss coverage --json`,
	RunE: runCoverage,
}

func init() {
	coverageCmd.Flags().BoolVar(&coverageJsonOutput, "json", false, "output as JSON")
	rootCmd.AddCommand(coverageCmd)
}

func runCoverage(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()

	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	report := analyzeCoverage(ds)

	if coverageJsonOutput {
		jsonData, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		printCoverageReport(report)
	}

	return nil
}

func analyzeCoverage(ds *dss.DesignSystem) *CoverageReport {
	report := &CoverageReport{
		DesignSystem: ds.Meta.Name,
		Version:      ds.Meta.Version,
		Categories:   []CategoryCoverage{},
	}

	// Meta
	meta := CategoryCoverage{
		Name:        "Meta",
		Description: "System metadata and maintainer info",
		Required:    true,
		Items: []ItemCoverage{
			{Name: "name", Present: ds.Meta.Name != "", Required: true},
			{Name: "version", Present: ds.Meta.Version != "", Required: true},
			{Name: "description", Present: ds.Meta.Description != "", Required: false},
			{Name: "maintainers", Present: len(ds.Meta.Maintainers) > 0, Required: false, Count: len(ds.Meta.Maintainers)},
		},
	}
	calculateCategoryScore(&meta)
	report.Categories = append(report.Categories, meta)

	// Principles
	principles := CategoryCoverage{
		Name:        "Principles",
		Description: "Design philosophy and guiding principles",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "principles", Present: len(ds.Principles) > 0, Required: true, Count: len(ds.Principles)},
		},
	}
	// Check if principles have examples
	hasExamples := false
	for _, p := range ds.Principles {
		if len(p.Examples) > 0 {
			hasExamples = true
			break
		}
	}
	principles.Items = append(principles.Items, ItemCoverage{
		Name: "examples", Present: hasExamples, Required: false,
	})
	calculateCategoryScore(&principles)
	report.Categories = append(report.Categories, principles)

	// Foundations
	f := ds.Foundations
	foundations := CategoryCoverage{
		Name:        "Foundations",
		Description: "Design tokens (colors, typography, spacing, etc.)",
		Required:    true,
		Items: []ItemCoverage{
			{Name: "colors", Present: len(f.Colors) > 0, Required: true, Count: len(f.Colors)},
			{Name: "typography", Present: f.Typography != nil && len(f.Typography.FontFamilies) > 0, Required: true},
			{Name: "spacing", Present: f.Spacing != nil && len(f.Spacing.Scale) > 0, Required: true},
			{Name: "elevation", Present: len(f.Elevation) > 0, Required: false, Count: len(f.Elevation)},
			{Name: "motion", Present: f.Motion != nil, Required: false},
			{Name: "grid", Present: f.Grid != nil, Required: false},
			{Name: "breakpoints", Present: len(f.Breakpoints) > 0, Required: false, Count: len(f.Breakpoints)},
			{Name: "borderRadius", Present: len(f.BorderRadius) > 0, Required: false, Count: len(f.BorderRadius)},
			{Name: "borderWidth", Present: len(f.BorderWidth) > 0, Required: false, Count: len(f.BorderWidth)},
			{Name: "opacity", Present: len(f.Opacity) > 0, Required: false, Count: len(f.Opacity)},
			{Name: "zIndex", Present: len(f.ZIndex) > 0, Required: false, Count: len(f.ZIndex)},
		},
	}
	calculateCategoryScore(&foundations)
	report.Categories = append(report.Categories, foundations)

	// Components
	components := CategoryCoverage{
		Name:        "Components",
		Description: "UI component specifications",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "components", Present: len(ds.Components) > 0, Required: true, Count: len(ds.Components)},
		},
	}
	// Check component completeness
	hasVariants, hasStates, hasProps, hasSlots, hasLLM := false, false, false, false, false
	for _, c := range ds.Components {
		if len(c.Variants) > 0 {
			hasVariants = true
		}
		if len(c.States) > 0 {
			hasStates = true
		}
		if len(c.Props) > 0 {
			hasProps = true
		}
		if len(c.Slots) > 0 {
			hasSlots = true
		}
		if c.LLM != nil {
			hasLLM = true
		}
	}
	components.Items = append(components.Items,
		ItemCoverage{Name: "variants", Present: hasVariants, Required: false},
		ItemCoverage{Name: "states", Present: hasStates, Required: false},
		ItemCoverage{Name: "props", Present: hasProps, Required: false},
		ItemCoverage{Name: "slots", Present: hasSlots, Required: false},
		ItemCoverage{Name: "llmContext", Present: hasLLM, Required: false},
	)
	calculateCategoryScore(&components)
	report.Categories = append(report.Categories, components)

	// Patterns
	patterns := CategoryCoverage{
		Name:        "Patterns",
		Description: "Multi-component compositions",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "patterns", Present: len(ds.Patterns) > 0, Required: true, Count: len(ds.Patterns)},
		},
	}
	calculateCategoryScore(&patterns)
	report.Categories = append(report.Categories, patterns)

	// Templates
	templates := CategoryCoverage{
		Name:        "Templates",
		Description: "Page layout templates",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "templates", Present: len(ds.Templates) > 0, Required: true, Count: len(ds.Templates)},
		},
	}
	calculateCategoryScore(&templates)
	report.Categories = append(report.Categories, templates)

	// Content
	content := CategoryCoverage{
		Name:        "Content",
		Description: "Voice, tone, and content guidelines",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "content", Present: ds.Content != nil, Required: true},
		},
	}
	if ds.Content != nil {
		content.Items = append(content.Items,
			ItemCoverage{Name: "voiceGuidelines", Present: ds.Content.Voice != nil, Required: false},
			ItemCoverage{Name: "toneGuidelines", Present: len(ds.Content.Tone) > 0, Required: false},
			ItemCoverage{Name: "terminology", Present: ds.Content.Terminology != nil, Required: false},
		)
	}
	calculateCategoryScore(&content)
	report.Categories = append(report.Categories, content)

	// Accessibility
	accessibility := CategoryCoverage{
		Name:        "Accessibility",
		Description: "WCAG compliance and requirements",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "accessibility", Present: ds.Accessibility != nil, Required: true},
		},
	}
	if ds.Accessibility != nil {
		accessibility.Items = append(accessibility.Items,
			ItemCoverage{Name: "wcagLevel", Present: ds.Accessibility.WCAGLevel != "", Required: false},
			ItemCoverage{Name: "colorContrast", Present: ds.Accessibility.ColorContrast != nil, Required: false},
			ItemCoverage{Name: "keyboard", Present: ds.Accessibility.Keyboard != nil, Required: false},
			ItemCoverage{Name: "screenReader", Present: ds.Accessibility.ScreenReader != nil, Required: false},
		)
	}
	calculateCategoryScore(&accessibility)
	report.Categories = append(report.Categories, accessibility)

	// Governance
	governance := CategoryCoverage{
		Name:        "Governance",
		Description: "Versioning and deprecation policies",
		Required:    false,
		Items: []ItemCoverage{
			{Name: "governance", Present: ds.Governance != nil, Required: true},
		},
	}
	if ds.Governance != nil {
		governance.Items = append(governance.Items,
			ItemCoverage{Name: "versioning", Present: ds.Governance.Versioning != nil, Required: false},
			ItemCoverage{Name: "deprecation", Present: ds.Governance.Deprecation != nil, Required: false},
			ItemCoverage{Name: "contribution", Present: ds.Governance.Contribution != nil, Required: false},
		)
	}
	calculateCategoryScore(&governance)
	report.Categories = append(report.Categories, governance)

	// Calculate summary
	calculateSummary(report)

	return report
}

func calculateCategoryScore(cat *CategoryCoverage) {
	if len(cat.Items) == 0 {
		cat.Score = 0
		cat.Status = "missing"
		return
	}

	total := 0.0
	present := 0.0

	for _, item := range cat.Items {
		weight := 1.0
		if item.Required {
			weight = 2.0 // Required items count double
		}
		total += weight
		if item.Present {
			present += weight
		}
	}

	if total == 0 {
		cat.Score = 0
	} else {
		cat.Score = (present / total) * 100
	}

	if cat.Score >= 80 {
		cat.Status = "complete"
	} else if cat.Score > 0 {
		cat.Status = "partial"
	} else {
		cat.Status = "missing"
	}
}

func calculateSummary(report *CoverageReport) {
	totalScore := 0.0
	for _, cat := range report.Categories {
		report.Summary.TotalCategories++
		switch cat.Status {
		case "complete":
			report.Summary.CompleteCategories++
		case "partial":
			report.Summary.PartialCategories++
		case "missing":
			report.Summary.MissingCategories++
		}
		totalScore += cat.Score
	}

	if report.Summary.TotalCategories > 0 {
		report.Summary.OverallScore = totalScore / float64(report.Summary.TotalCategories)
	}

	// Assign grade
	switch {
	case report.Summary.OverallScore >= 90:
		report.Summary.Grade = "A"
	case report.Summary.OverallScore >= 80:
		report.Summary.Grade = "B"
	case report.Summary.OverallScore >= 70:
		report.Summary.Grade = "C"
	case report.Summary.OverallScore >= 60:
		report.Summary.Grade = "D"
	default:
		report.Summary.Grade = "F"
	}
}

func printCoverageReport(report *CoverageReport) {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         DSS COVERAGE REPORT                                ║")
	fmt.Println("╠════════════════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Design System: %-60s ║\n", report.DesignSystem)
	fmt.Printf("║ Version: %-66s ║\n", report.Version)
	fmt.Println("╠════════════════════════════════════════════════════════════════════════════╣")

	for _, cat := range report.Categories {
		icon := "🔴"
		if cat.Status == "complete" {
			icon = "🟢"
		} else if cat.Status == "partial" {
			icon = "🟡"
		}

		name := fmt.Sprintf("%-15s", cat.Name)
		score := fmt.Sprintf("%5.1f%%", cat.Score)
		fmt.Printf("║ %s %s %s ", icon, name, score)

		// Show item details
		details := []string{}
		for _, item := range cat.Items {
			if item.Present {
				if item.Count > 0 {
					details = append(details, fmt.Sprintf("%s:%d", item.Name, item.Count))
				} else {
					details = append(details, "✓"+item.Name)
				}
			}
		}
		detailStr := ""
		if len(details) > 0 {
			detailStr = fmt.Sprintf("(%s)", joinMax(details, ", ", 40))
		}
		fmt.Printf("%-42s║\n", detailStr)
	}

	fmt.Println("╠════════════════════════════════════════════════════════════════════════════╣")

	gradeIcon := "🔴"
	if report.Summary.Grade == "A" || report.Summary.Grade == "B" {
		gradeIcon = "🟢"
	} else if report.Summary.Grade == "C" {
		gradeIcon = "🟡"
	}

	summary := fmt.Sprintf("Grade: %s  Score: %.1f%%  (%d complete, %d partial, %d missing)",
		report.Summary.Grade, report.Summary.OverallScore,
		report.Summary.CompleteCategories, report.Summary.PartialCategories, report.Summary.MissingCategories)
	fmt.Printf("║ %s %-73s║\n", gradeIcon, summary)
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")

	// Print recommendations
	fmt.Println("\nRecommendations:")
	for _, cat := range report.Categories {
		if cat.Status == "missing" {
			fmt.Printf("  • Add %s: %s\n", cat.Name, cat.Description)
		} else if cat.Status == "partial" {
			missing := []string{}
			for _, item := range cat.Items {
				if !item.Present {
					missing = append(missing, item.Name)
				}
			}
			if len(missing) > 0 {
				fmt.Printf("  • %s: add %s\n", cat.Name, joinMax(missing, ", ", 50))
			}
		}
	}
}

func joinMax(items []string, sep string, maxLen int) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		if len(result)+len(item) > maxLen {
			result += "..."
			break
		}
		result += item
	}
	return result
}
