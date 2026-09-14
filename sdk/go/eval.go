package dss

import (
	"context"
	"fmt"

	"github.com/plexusone/structured-evaluation/rubric"
)

// ReviewTypeDSS is the review type identifier for DSS evaluations.
const ReviewTypeDSS = "dss-spec"

// EvalOptions configures evaluation behavior.
type EvalOptions struct {
	// IncludeIssues includes detailed issue list (default: true)
	IncludeIssues bool `json:"includeIssues,omitempty"`

	// MinScore is the minimum acceptable score (1-5, default: 4)
	MinScore int `json:"minScore,omitempty"`

	// Categories to evaluate (default: all)
	Categories []string `json:"categories,omitempty"`

	// Verbose enables verbose output
	Verbose bool `json:"verbose,omitempty"`
}

// DefaultEvalOptions returns sensible defaults.
func DefaultEvalOptions() *EvalOptions {
	return &EvalOptions{
		IncludeIssues: true,
		MinScore:      4, // maps to "Good" in 1-5 scale
		Categories:    []string{"completeness", "agent-readiness", "accessibility", "documentation"},
		Verbose:       false,
	}
}

// DSS coverage section names for CoverageReport.
const (
	CoverageSectionComponents    = "components"
	CoverageSectionFoundations   = "foundations"
	CoverageSectionPatterns      = "patterns"
	CoverageSectionAccessibility = "accessibility"
)

// Category weights for evaluation scoring (0.0-1.0)
var evalCategoryWeights = map[string]float64{
	"completeness":    0.25,
	"agent-readiness": 0.30,
	"accessibility":   0.25,
	"documentation":   0.20,
}

// Evaluate runs a complete evaluation of the design system spec.
// Returns a structured-evaluation rubric.Rubric for consistency with other evaluations.
func (s *Service) Evaluate(ctx context.Context, opts *EvalOptions) (*rubric.Rubric, error) {
	if opts == nil {
		opts = DefaultEvalOptions()
	}

	// Create rubric with DSS review type
	docName := s.ds.Meta.Name
	if docName == "" {
		docName = "systemspec-designsystem"
	}
	r := rubric.NewRubric(ReviewTypeDSS, docName)
	r.Metadata.DocumentVersion = s.ds.Meta.Version
	r.Metadata.GeneratedBy = "systemspec-designsystem/eval"

	// Evaluate each category
	for _, catName := range opts.Categories {
		var catResult *rubric.CategoryResult

		switch catName {
		case "completeness":
			catResult = s.evaluateCompleteness(ctx)
		case "agent-readiness":
			catResult = s.evaluateAgentReadiness(ctx)
		case "accessibility":
			catResult = s.evaluateAccessibility(ctx)
		case "documentation":
			catResult = s.evaluateDocumentation(ctx)
		default:
			continue
		}

		if catResult != nil {
			r.AddCategoryResult(*catResult)
		}
	}

	// Calculate coverage and store using structured-evaluation's CoverageReport
	coverage := s.calculateCoverage(ctx)
	r.SetCoverage(coverage)
	r.SetRubricInfo("dss-eval", "1.0")

	// Set pass criteria with minimum score threshold
	r.PassCriteria = rubric.DefaultPassCriteria()
	if opts.MinScore > 0 {
		r.PassCriteria.MinIntScore = rubric.ParseIntegerScore(opts.MinScore)
	}

	// Finalize the rubric (computes decision, next steps, summary)
	r.Finalize(nil, "dss eval --spec ./path/to/spec")

	// Append coverage summary
	compSection := coverage.GetSection(CoverageSectionComponents)
	foundSection := coverage.GetSection(CoverageSectionFoundations)
	pattSection := coverage.GetSection(CoverageSectionPatterns)
	a11ySection := coverage.GetSection(CoverageSectionAccessibility)
	r.Summary = fmt.Sprintf("%s Coverage: components %d%%, foundations %d%%, patterns %d%%, accessibility %d%%, overall %d%%.",
		r.Summary,
		compSection.Percentage,
		foundSection.Percentage,
		pattSection.Percentage,
		a11ySection.Percentage,
		coverage.Overall,
	)

	return r, nil
}

// evaluateCompleteness checks for required fields and completeness.
//
//nolint:unparam // ctx reserved for future logging/cancellation
func (s *Service) evaluateCompleteness(ctx context.Context) *rubric.CategoryResult {
	catResult := rubric.NewCategoryResult("completeness", rubric.ScorePass, "")

	var checks, passed int

	// Check meta completeness
	if s.ds.Meta.Name == "" {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-meta-name",
			Category:       "completeness",
			Severity:       rubric.SeverityHigh,
			Title:          "Design system name is required",
			Description:    "The meta.name field is missing",
			Location:       "meta.name",
			Recommendation: "Add a name field to meta.json",
		})
	} else {
		passed++
	}
	checks++

	if s.ds.Meta.Version == "" {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-meta-version",
			Category:       "completeness",
			Severity:       rubric.SeverityHigh,
			Title:          "Design system version is required",
			Description:    "The meta.version field is missing",
			Location:       "meta.version",
			Recommendation: "Add a version field to meta.json (e.g., '1.0.0')",
		})
	} else {
		passed++
	}
	checks++

	if s.ds.Meta.Description == "" {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-meta-description",
			Category:       "completeness",
			Severity:       rubric.SeverityMedium,
			Title:          "Design system description is recommended",
			Description:    "The meta.description field is missing",
			Location:       "meta.description",
			Recommendation: "Add a description field to meta.json",
		})
	} else {
		passed++
	}
	checks++

	// Check foundations
	if len(s.ds.Foundations.Colors) == 0 {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-colors",
			Category:       "completeness",
			Severity:       rubric.SeverityHigh,
			Title:          "At least one color token is required",
			Description:    "No color tokens defined in foundations",
			Location:       "foundations.colors",
			Recommendation: "Add color tokens to foundations/colors.json",
		})
	} else {
		passed++
	}
	checks++

	if s.ds.Foundations.Spacing == nil || len(s.ds.Foundations.Spacing.Scale) == 0 {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-spacing",
			Category:       "completeness",
			Severity:       rubric.SeverityHigh,
			Title:          "At least one spacing token is required",
			Description:    "No spacing tokens defined in foundations",
			Location:       "foundations.spacing",
			Recommendation: "Add spacing tokens to foundations/spacing.json",
		})
	} else {
		passed++
	}
	checks++

	if s.ds.Foundations.Typography == nil || s.ds.Foundations.Typography.FontFamilies == nil || len(s.ds.Foundations.Typography.FontFamilies) == 0 {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-typography",
			Category:       "completeness",
			Severity:       rubric.SeverityMedium,
			Title:          "Typography font families are recommended",
			Description:    "No font families defined in typography",
			Location:       "foundations.typography.fontFamilies",
			Recommendation: "Add typography definitions to foundations/typography.json",
		})
	} else {
		passed++
	}
	checks++

	// Check components
	if len(s.ds.Components) == 0 {
		catResult.AddFinding(rubric.Finding{
			ID:             "completeness-components",
			Category:       "completeness",
			Severity:       rubric.SeverityHigh,
			Title:          "At least one component is required",
			Description:    "No components defined",
			Location:       "components",
			Recommendation: "Add component definitions to components/",
		})
	} else {
		passed++

		// Check each component for required fields
		for _, comp := range s.ds.Components {
			if comp.Name == "" {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("completeness-component-%s-name", comp.ID),
					Category:       "completeness",
					Severity:       rubric.SeverityHigh,
					Title:          "Component name is required",
					Description:    fmt.Sprintf("Component '%s' is missing a name", comp.ID),
					Location:       fmt.Sprintf("components.%s.name", comp.ID),
					Recommendation: "Add a name field to the component",
				})
			} else {
				passed++
			}
			checks++

			if len(comp.Props) == 0 {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("completeness-component-%s-props", comp.ID),
					Category:       "completeness",
					Severity:       rubric.SeverityMedium,
					Title:          fmt.Sprintf("Component '%s' has no props defined", comp.Name),
					Description:    "Components should have props to be useful",
					Location:       fmt.Sprintf("components.%s.props", comp.ID),
					Recommendation: "Add props array with component properties",
				})
			} else {
				passed++
			}
			checks++
		}
	}
	checks++

	// Calculate score (1-5 scale)
	intScore := percentageToIntScore(passed, checks)
	catResult.SetIntScore(intScore)
	catResult.SetConfidence(0.9)
	catResult.Reasoning = fmt.Sprintf("Checked %d completeness requirements, %d passed", checks, passed)

	return catResult
}

// evaluateAgentReadiness checks for LLM-friendly specifications.
//
//nolint:unparam // ctx reserved for future logging/cancellation
func (s *Service) evaluateAgentReadiness(ctx context.Context) *rubric.CategoryResult {
	catResult := rubric.NewCategoryResult("agent-readiness", rubric.ScorePass, "")

	var checks, passed int

	// Check for LLMContext on components
	for _, comp := range s.ds.Components {
		if comp.LLM == nil {
			catResult.AddFinding(rubric.Finding{
				ID:             fmt.Sprintf("agent-readiness-%s-llm-context", comp.ID),
				Category:       "agent-readiness",
				Severity:       rubric.SeverityMedium,
				Title:          fmt.Sprintf("Component '%s' is missing LLM context", comp.Name),
				Description:    "LLM context helps AI agents understand how to use the component",
				Location:       fmt.Sprintf("components.%s.llm", comp.ID),
				Recommendation: "Add llm field with intent, allowedContexts, and antiPatterns",
			})
		} else {
			// Check quality of LLM context
			if comp.LLM.Intent == "" {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("agent-readiness-%s-llm-intent", comp.ID),
					Category:       "agent-readiness",
					Severity:       rubric.SeverityMedium,
					Title:          fmt.Sprintf("Component '%s' LLM context missing intent", comp.Name),
					Description:    "The intent field describes the component's primary purpose",
					Location:       fmt.Sprintf("components.%s.llm.intent", comp.ID),
					Recommendation: "Add intent describing the component's primary purpose",
				})
			} else {
				passed++
			}
			checks++

			if len(comp.LLM.AntiPatterns) == 0 {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("agent-readiness-%s-llm-antipatterns", comp.ID),
					Category:       "agent-readiness",
					Severity:       rubric.SeverityLow,
					Title:          fmt.Sprintf("Component '%s' has no anti-patterns defined", comp.Name),
					Description:    "Anti-patterns help agents avoid common mistakes",
					Location:       fmt.Sprintf("components.%s.llm.antiPatterns", comp.ID),
					Recommendation: "Add antiPatterns array with common mistakes to avoid",
				})
			} else {
				passed++
			}
			checks++

			if len(comp.LLM.ExampleUsage) == 0 {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("agent-readiness-%s-llm-examples", comp.ID),
					Category:       "agent-readiness",
					Severity:       rubric.SeverityLow,
					Title:          fmt.Sprintf("Component '%s' has no example usage", comp.Name),
					Description:    "Examples help agents understand correct usage patterns",
					Location:       fmt.Sprintf("components.%s.llm.exampleUsage", comp.ID),
					Recommendation: "Add exampleUsage array with code snippets",
				})
			} else {
				passed++
			}
			checks++
		}
	}

	// Check for anti-patterns at spec level
	if len(s.ds.Principles) == 0 {
		catResult.AddFinding(rubric.Finding{
			ID:             "agent-readiness-principles",
			Category:       "agent-readiness",
			Severity:       rubric.SeverityMedium,
			Title:          "Design principles help agents understand the design system philosophy",
			Description:    "Principles provide high-level guidance for consistent decisions",
			Location:       "principles",
			Recommendation: "Add principles array with design guidelines",
		})
	} else {
		passed++
	}
	checks++

	// Check for patterns
	if len(s.ds.Patterns) == 0 {
		catResult.AddFinding(rubric.Finding{
			ID:             "agent-readiness-patterns",
			Category:       "agent-readiness",
			Severity:       rubric.SeverityLow,
			Title:          "Patterns help agents understand component combinations",
			Description:    "Patterns document common usage scenarios",
			Location:       "patterns",
			Recommendation: "Add patterns array with usage patterns",
		})
	} else {
		passed++
	}
	checks++

	// Calculate score (1-5 scale)
	intScore := percentageToIntScore(passed, checks)
	catResult.SetIntScore(intScore)
	catResult.SetConfidence(0.85)
	catResult.Reasoning = fmt.Sprintf("Checked %d agent-readiness requirements, %d passed", checks, passed)

	return catResult
}

// evaluateAccessibility checks for accessibility requirements.
//
//nolint:unparam // ctx reserved for future logging/cancellation
func (s *Service) evaluateAccessibility(ctx context.Context) *rubric.CategoryResult {
	catResult := rubric.NewCategoryResult("accessibility", rubric.ScorePass, "")

	var checks, passed int

	// Check global accessibility settings
	if s.ds.Accessibility == nil {
		catResult.AddFinding(rubric.Finding{
			ID:             "accessibility-section",
			Category:       "accessibility",
			Severity:       rubric.SeverityHigh,
			Title:          "Accessibility section is required",
			Description:    "Global accessibility requirements should be defined",
			Location:       "accessibility",
			Recommendation: "Add accessibility.json with WCAG guidelines and requirements",
		})
	} else {
		passed++
		checks++

		// Check WCAG level
		if s.ds.Accessibility.WCAGLevel == "" {
			catResult.AddFinding(rubric.Finding{
				ID:             "accessibility-wcag-level",
				Category:       "accessibility",
				Severity:       rubric.SeverityMedium,
				Title:          "WCAG compliance level should be specified",
				Description:    "Specify the target WCAG compliance level",
				Location:       "accessibility.wcagLevel",
				Recommendation: "Set wcagLevel to 'AA' (recommended) or 'AAA'",
			})
		} else {
			passed++
		}
		checks++

		// Check keyboard requirements
		if s.ds.Accessibility.Keyboard == nil {
			catResult.AddFinding(rubric.Finding{
				ID:             "accessibility-keyboard",
				Category:       "accessibility",
				Severity:       rubric.SeverityMedium,
				Title:          "Keyboard navigation requirements should be specified",
				Description:    "Keyboard accessibility is essential for WCAG compliance",
				Location:       "accessibility.keyboard",
				Recommendation: "Add keyboard section with focus management guidelines",
			})
		} else {
			passed++
		}
		checks++

		// Check screen reader requirements
		if s.ds.Accessibility.ScreenReader == nil {
			catResult.AddFinding(rubric.Finding{
				ID:             "accessibility-screenreader",
				Category:       "accessibility",
				Severity:       rubric.SeverityMedium,
				Title:          "Screen reader requirements should be specified",
				Description:    "Screen reader support is essential for WCAG compliance",
				Location:       "accessibility.screenReader",
				Recommendation: "Add screenReader section with ARIA and landmark guidelines",
			})
		} else {
			passed++
		}
		checks++
	}
	checks++

	// Check component-level accessibility
	for _, comp := range s.ds.Components {
		if comp.Accessibility == nil {
			catResult.AddFinding(rubric.Finding{
				ID:             fmt.Sprintf("accessibility-component-%s", comp.ID),
				Category:       "accessibility",
				Severity:       rubric.SeverityMedium,
				Title:          fmt.Sprintf("Component '%s' should have accessibility requirements", comp.Name),
				Description:    "Each component should define its accessibility requirements",
				Location:       fmt.Sprintf("components.%s.accessibility", comp.ID),
				Recommendation: "Add accessibility section with role, aria attributes, keyboard interactions",
			})
		} else {
			passed++

			// Check for role
			if comp.Accessibility.Role == "" {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("accessibility-component-%s-role", comp.ID),
					Category:       "accessibility",
					Severity:       rubric.SeverityLow,
					Title:          fmt.Sprintf("Component '%s' should specify an ARIA role", comp.Name),
					Description:    "ARIA roles help assistive technologies understand the component",
					Location:       fmt.Sprintf("components.%s.accessibility.role", comp.ID),
					Recommendation: "Add role field with appropriate ARIA role",
				})
			} else {
				passed++
			}
			checks++

			// Check for keyboard interactions
			if len(comp.Accessibility.KeyboardSupport) == 0 {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("accessibility-component-%s-keyboard", comp.ID),
					Category:       "accessibility",
					Severity:       rubric.SeverityLow,
					Title:          fmt.Sprintf("Component '%s' should define keyboard interactions", comp.Name),
					Description:    "Keyboard support ensures the component is accessible",
					Location:       fmt.Sprintf("components.%s.accessibility.keyboardSupport", comp.ID),
					Recommendation: "Add keyboardSupport array with key bindings",
				})
			} else {
				passed++
			}
			checks++
		}
		checks++
	}

	// Calculate score (1-5 scale)
	intScore := percentageToIntScore(passed, checks)
	catResult.SetIntScore(intScore)
	catResult.SetConfidence(0.9)
	catResult.Reasoning = fmt.Sprintf("Checked %d accessibility requirements, %d passed", checks, passed)

	return catResult
}

// evaluateDocumentation checks for documentation quality.
//
//nolint:unparam // ctx reserved for future logging/cancellation
func (s *Service) evaluateDocumentation(ctx context.Context) *rubric.CategoryResult {
	catResult := rubric.NewCategoryResult("documentation", rubric.ScorePass, "")

	var checks, passed int

	// Check meta documentation
	if s.ds.Meta.Description == "" {
		catResult.AddFinding(rubric.Finding{
			ID:             "documentation-meta-description",
			Category:       "documentation",
			Severity:       rubric.SeverityMedium,
			Title:          "Design system should have a description",
			Description:    "A description helps users understand the design system's purpose",
			Location:       "meta.description",
			Recommendation: "Add a description explaining the design system's purpose",
		})
	} else {
		passed++
	}
	checks++

	// Check component documentation
	for _, comp := range s.ds.Components {
		if comp.Description == "" {
			catResult.AddFinding(rubric.Finding{
				ID:             fmt.Sprintf("documentation-component-%s-description", comp.ID),
				Category:       "documentation",
				Severity:       rubric.SeverityMedium,
				Title:          fmt.Sprintf("Component '%s' should have a description", comp.Name),
				Description:    "Component descriptions help users understand when to use them",
				Location:       fmt.Sprintf("components.%s.description", comp.ID),
				Recommendation: "Add a description explaining the component's purpose",
			})
		} else {
			passed++
		}
		checks++

		// Check prop documentation
		for _, prop := range comp.Props {
			if prop.Description == "" {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("documentation-component-%s-prop-%s", comp.ID, prop.Name),
					Category:       "documentation",
					Severity:       rubric.SeverityLow,
					Title:          fmt.Sprintf("Prop '%s' in '%s' should have a description", prop.Name, comp.Name),
					Description:    "Prop descriptions help users understand how to configure components",
					Location:       fmt.Sprintf("components.%s.props.%s.description", comp.ID, prop.Name),
					Recommendation: "Add a description explaining the prop's purpose",
				})
			} else {
				passed++
			}
			checks++
		}

		// Check variant documentation
		for _, variant := range comp.Variants {
			if variant.Description == "" {
				catResult.AddFinding(rubric.Finding{
					ID:             fmt.Sprintf("documentation-component-%s-variant-%s", comp.ID, variant.ID),
					Category:       "documentation",
					Severity:       rubric.SeverityLow,
					Title:          fmt.Sprintf("Variant '%s' in '%s' should have a description", variant.Name, comp.Name),
					Description:    "Variant descriptions help users choose the right variant",
					Location:       fmt.Sprintf("components.%s.variants.%s.description", comp.ID, variant.ID),
					Recommendation: "Add a description explaining when to use this variant",
				})
			} else {
				passed++
			}
			checks++
		}
	}

	// Check pattern documentation
	for _, pattern := range s.ds.Patterns {
		if pattern.Description == "" {
			catResult.AddFinding(rubric.Finding{
				ID:             fmt.Sprintf("documentation-pattern-%s", pattern.ID),
				Category:       "documentation",
				Severity:       rubric.SeverityMedium,
				Title:          fmt.Sprintf("Pattern '%s' should have a description", pattern.Name),
				Description:    "Pattern descriptions help users understand use cases",
				Location:       fmt.Sprintf("patterns.%s.description", pattern.ID),
				Recommendation: "Add a description explaining the pattern's use case",
			})
		} else {
			passed++
		}
		checks++
	}

	// Calculate score (1-5 scale)
	intScore := percentageToIntScore(passed, checks)
	catResult.SetIntScore(intScore)
	catResult.SetConfidence(0.9)
	catResult.Reasoning = fmt.Sprintf("Checked %d documentation requirements, %d passed", checks, passed)

	return catResult
}

// calculateCoverage computes spec coverage metrics using structured-evaluation's CoverageReport.
//
//nolint:unparam // ctx reserved for future logging/cancellation
func (s *Service) calculateCoverage(ctx context.Context) *rubric.CoverageReport {
	cr := rubric.NewCoverageReport()

	// Components coverage
	var compMissing []string
	compComplete := 0
	for _, comp := range s.ds.Components {
		if isComponentComplete(comp) {
			compComplete++
		} else {
			compMissing = append(compMissing, comp.ID)
		}
	}
	cr.SetSection(CoverageSectionComponents, len(s.ds.Components), compComplete, compMissing)

	// Foundations coverage
	foundationItems := 0
	foundationComplete := 0
	var foundMissing []string

	if len(s.ds.Foundations.Colors) > 0 {
		foundationItems++
		foundationComplete++
	} else {
		foundationItems++
		foundMissing = append(foundMissing, "colors")
	}

	if s.ds.Foundations.Spacing != nil && len(s.ds.Foundations.Spacing.Scale) > 0 {
		foundationItems++
		foundationComplete++
	} else {
		foundationItems++
		foundMissing = append(foundMissing, "spacing")
	}

	if s.ds.Foundations.Typography != nil && s.ds.Foundations.Typography.FontFamilies != nil && len(s.ds.Foundations.Typography.FontFamilies) > 0 {
		foundationItems++
		foundationComplete++
	} else {
		foundationItems++
		foundMissing = append(foundMissing, "typography")
	}

	if len(s.ds.Foundations.Elevation) > 0 {
		foundationItems++
		foundationComplete++
	} else {
		foundationItems++
		foundMissing = append(foundMissing, "elevation")
	}

	cr.SetSection(CoverageSectionFoundations, foundationItems, foundationComplete, foundMissing)

	// Patterns coverage
	var pattMissing []string
	pattComplete := 0
	for _, pattern := range s.ds.Patterns {
		if isPatternComplete(pattern) {
			pattComplete++
		} else {
			pattMissing = append(pattMissing, pattern.ID)
		}
	}
	cr.SetSection(CoverageSectionPatterns, len(s.ds.Patterns), pattComplete, pattMissing)

	// Accessibility coverage
	a11yTotal := 0
	a11yComplete := 0
	var a11yMissing []string
	if s.ds.Accessibility != nil {
		a11yTotal = 1
		if isAccessibilityComplete(s.ds.Accessibility) {
			a11yComplete = 1
		} else {
			a11yMissing = append(a11yMissing, "accessibility")
		}
	}
	cr.SetSection(CoverageSectionAccessibility, a11yTotal, a11yComplete, a11yMissing)

	// Compute overall using weighted average
	// Components and Accessibility are more important
	weights := map[string]float64{
		CoverageSectionComponents:    2.0,
		CoverageSectionFoundations:   1.0,
		CoverageSectionPatterns:      1.0,
		CoverageSectionAccessibility: 2.0,
	}
	cr.ComputeOverallWeighted(weights)

	return cr
}

// Helper functions

// percentageToIntScore converts a percentage (0-100) to IntegerScore (1-5).
func percentageToIntScore(passed, checks int) rubric.IntegerScore {
	if checks == 0 {
		return rubric.ScoreExcellent
	}
	percentage := (passed * 100) / checks

	switch {
	case percentage >= 90:
		return rubric.ScoreExcellent
	case percentage >= 75:
		return rubric.ScoreGood
	case percentage >= 50:
		return rubric.ScoreAcceptable
	case percentage >= 25:
		return rubric.ScoreMajorRevisions
	default:
		return rubric.ScoreUnacceptable
	}
}

func isComponentComplete(comp Component) bool {
	if comp.Name == "" || comp.Description == "" {
		return false
	}
	if len(comp.Props) == 0 {
		return false
	}
	if comp.Accessibility == nil {
		return false
	}
	return true
}

func isPatternComplete(pattern Pattern) bool {
	if pattern.Name == "" || pattern.Description == "" {
		return false
	}
	if len(pattern.Components) == 0 {
		return false
	}
	return true
}

func isAccessibilityComplete(a11y *Accessibility) bool {
	if a11y.WCAGLevel == "" {
		return false
	}
	return true
}
