package dss

import (
	"context"
	"strings"
	"testing"
)

func TestLintSpec_EmptyDesignSystem(t *testing.T) {
	ds := &DesignSystem{}
	result := ds.Lint()

	// Should have errors for missing meta
	if result.Summary.Errors < 2 {
		t.Errorf("Expected at least 2 errors for missing meta, got %d", result.Summary.Errors)
	}

	// Score should be low
	if result.Score > 80 {
		t.Errorf("Expected low score for empty design system, got %d", result.Score)
	}
}

func TestLintSpec_MinimalValidDesignSystem(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
	}
	result := ds.Lint()

	// Should have no errors for meta
	hasMetaError := false
	for _, issue := range result.Issues {
		if issue.Rule == "meta-required" && issue.Severity == "error" {
			hasMetaError = true
		}
	}
	if hasMetaError {
		t.Error("Should not have meta-required errors when meta is valid")
	}
}

func TestLintSpec_ComponentWithoutLLMContext(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
			},
		},
	}

	result := ds.Lint()

	// Should have warning for missing LLM context
	hasLLMWarning := false
	for _, issue := range result.Issues {
		if issue.Rule == "component-has-llm-context" && issue.Component == "button" {
			hasLLMWarning = true
		}
	}
	if !hasLLMWarning {
		t.Error("Expected warning for component without LLM context")
	}
}

func TestLintSpec_ComponentWithLLMContextMissingIntent(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				LLM:  &LLMContext{}, // Empty LLM context
			},
		},
	}

	result := ds.Lint()

	// Should have error for missing intent
	hasIntentError := false
	for _, issue := range result.Issues {
		if issue.Rule == "llm-has-intent" && issue.Severity == "error" {
			hasIntentError = true
		}
	}
	if !hasIntentError {
		t.Error("Expected error for LLM context without intent")
	}
}

func TestLintSpec_ComponentWithFullLLMContext(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				Variants: []Variant{
					{ID: "default"},
					{ID: "secondary"},
				},
				Props: []Prop{
					{Name: "disabled", Type: "boolean"},
				},
				LLM: &LLMContext{
					Intent:          "Trigger user actions",
					AllowedContexts: []string{"forms", "dialogs"},
					AntiPatterns:    []string{"Multiple primary buttons"},
				},
			},
		},
	}

	result := ds.Lint()

	// Should have no component-related warnings
	componentWarnings := 0
	for _, issue := range result.Issues {
		if issue.Component == "button" && issue.Severity != "info" {
			componentWarnings++
		}
	}
	if componentWarnings > 0 {
		t.Errorf("Expected no warnings for well-defined component, got %d", componentWarnings)
	}

	// Coverage should be 100% for component metrics
	if result.Coverage.ComponentsWithLLMContext != 100 {
		t.Errorf("Expected 100%% LLM context coverage, got %.1f%%", result.Coverage.ComponentsWithLLMContext)
	}
	if result.Coverage.ComponentsWithVariants != 100 {
		t.Errorf("Expected 100%% variants coverage, got %.1f%%", result.Coverage.ComponentsWithVariants)
	}
}

func TestLintSpec_InvalidTokenReference(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:         "button",
				Name:       "Button",
				TokensUsed: []string{"nonexistent-token"},
			},
		},
	}

	result := ds.Lint()

	// Should have error for invalid token reference
	hasTokenError := false
	for _, issue := range result.Issues {
		if issue.Rule == "token-references-valid" {
			hasTokenError = true
		}
	}
	if !hasTokenError {
		t.Error("Expected error for invalid token reference")
	}
}

func TestLintSpec_ValidTokenReference(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#3B82F6"},
			},
		},
		Components: []Component{
			{
				ID:         "button",
				Name:       "Button",
				TokensUsed: []string{"primary"},
			},
		},
	}

	result := ds.Lint()

	// Should have no error for valid token reference
	hasTokenError := false
	for _, issue := range result.Issues {
		if issue.Rule == "token-references-valid" {
			hasTokenError = true
		}
	}
	if hasTokenError {
		t.Error("Should not have error for valid token reference")
	}
}

func TestLintSpec_OrphanToken(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#3B82F6"},
				{ID: "unused", Value: "#FF0000"},
			},
		},
		Components: []Component{
			{
				ID:         "button",
				Name:       "Button",
				TokensUsed: []string{"primary"},
			},
		},
	}

	result := ds.Lint()

	// Should have info for orphan token
	hasOrphanInfo := false
	for _, issue := range result.Issues {
		if issue.Rule == "no-orphan-tokens" && issue.Path == "foundations.colors[1]" {
			hasOrphanInfo = true
		}
	}
	if !hasOrphanInfo {
		t.Error("Expected info for orphan token")
	}
}

func TestLintSpec_InvalidComponentUses(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "card",
				Name: "Card",
				Uses: []string{"nonexistent-component"},
			},
		},
	}

	result := ds.Lint()

	// Should have error for invalid uses reference
	hasUsesError := false
	for _, issue := range result.Issues {
		if issue.Rule == "component-uses-valid" {
			hasUsesError = true
		}
	}
	if !hasUsesError {
		t.Error("Expected error for invalid component uses reference")
	}
}

func TestLintSpec_ValidComponentUses(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
			},
			{
				ID:   "card",
				Name: "Card",
				Uses: []string{"button"},
			},
		},
	}

	result := ds.Lint()

	// Should have no error for valid uses reference
	hasUsesError := false
	for _, issue := range result.Issues {
		if issue.Rule == "component-uses-valid" {
			hasUsesError = true
		}
	}
	if hasUsesError {
		t.Error("Should not have error for valid component uses reference")
	}
}

func TestLintSpec_MissingAccessibility(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
	}

	result := ds.Lint()

	// Should have warning for missing accessibility
	hasA11yWarning := false
	for _, issue := range result.Issues {
		if issue.Rule == "accessibility-defined" {
			hasA11yWarning = true
		}
	}
	if !hasA11yWarning {
		t.Error("Expected warning for missing accessibility section")
	}
}

func TestLintSpec_ThemingContractValidation(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "", // Missing prefix
				},
			},
		},
	}

	result := ds.Lint()

	// Should have error for missing theming contract prefix
	hasPrefixError := false
	for _, issue := range result.Issues {
		if issue.Rule == "theming-contract-valid" && issue.Severity == "error" {
			hasPrefixError = true
		}
	}
	if !hasPrefixError {
		t.Error("Expected error for missing theming contract prefix")
	}
}

func TestLintSpec_ThemingContractInvalidPrefix(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				ThemingContract: &ThemingContract{
					Prefix: "btn", // Should start with --
				},
			},
		},
	}

	result := ds.Lint()

	// Should have warning for invalid prefix format
	hasPrefixWarning := false
	for _, issue := range result.Issues {
		if issue.Rule == "theming-contract-valid" && issue.Severity == "warning" {
			hasPrefixWarning = true
		}
	}
	if !hasPrefixWarning {
		t.Error("Expected warning for prefix not starting with '--'")
	}
}

func TestLintSpec_SpecificRulesOnly(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				// Missing LLM context, variants, props
			},
		},
	}

	service := NewService(ds)
	result := service.LintSpec(context.Background(), &LintOptions{
		Rules: []string{"component-has-llm-context"}, // Only this rule
	})

	// Should only have LLM context issues
	for _, issue := range result.Issues {
		if issue.Rule != "component-has-llm-context" {
			t.Errorf("Expected only component-has-llm-context issues, got %s", issue.Rule)
		}
	}
}

func TestLintSpec_ScoreCalculation(t *testing.T) {
	// Well-defined system should have high score
	goodDS := &DesignSystem{
		Meta: Meta{
			Name:        "Test System",
			Version:     "1.0.0",
			Description: "A test design system",
		},
		Components: []Component{
			{
				ID:   "button",
				Name: "Button",
				Variants: []Variant{
					{ID: "default"},
				},
				Props: []Prop{
					{Name: "disabled", Type: "boolean"},
				},
				LLM: &LLMContext{
					Intent:          "Trigger actions",
					AllowedContexts: []string{"forms"},
					AntiPatterns:    []string{"Multiple primary buttons"},
				},
			},
		},
		Accessibility: &Accessibility{
			WCAGLevel: "AA",
		},
	}

	result := goodDS.Lint()
	if result.Score < 80 {
		t.Errorf("Expected score >= 80 for well-defined system, got %d", result.Score)
	}

	// Poorly defined system should have lower score
	badDS := &DesignSystem{}
	badResult := badDS.Lint()
	// Empty system has: missing name (-10), missing version (-10), missing description (-1), missing accessibility (-3) = 76
	if badResult.Score > 80 {
		t.Errorf("Expected score <= 80 for empty system, got %d", badResult.Score)
	}
	// Should have at least 2 errors
	if badResult.Summary.Errors < 2 {
		t.Errorf("Expected at least 2 errors for empty system, got %d", badResult.Summary.Errors)
	}
}

func TestLintSpec_CoverageCalculation(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Components: []Component{
			{
				ID:       "button",
				Name:     "Button",
				Variants: []Variant{{ID: "default"}},
				Props:    []Prop{{Name: "disabled", Type: "boolean"}},
				LLM:      &LLMContext{Intent: "Actions"},
			},
			{
				ID:   "card",
				Name: "Card",
				// No variants, props, or LLM
			},
		},
	}

	result := ds.Lint()

	// 50% of components have LLM context
	if result.Coverage.ComponentsWithLLMContext != 50 {
		t.Errorf("Expected 50%% LLM coverage, got %.1f%%", result.Coverage.ComponentsWithLLMContext)
	}

	// 50% have variants
	if result.Coverage.ComponentsWithVariants != 50 {
		t.Errorf("Expected 50%% variants coverage, got %.1f%%", result.Coverage.ComponentsWithVariants)
	}

	// 50% have props
	if result.Coverage.ComponentsWithProps != 50 {
		t.Errorf("Expected 50%% props coverage, got %.1f%%", result.Coverage.ComponentsWithProps)
	}
}

func TestAvailableLintRules(t *testing.T) {
	rules := AvailableLintRules()

	expectedRules := []string{
		"meta-required",
		"component-has-variants",
		"component-has-props",
		"component-has-llm-context",
		"llm-has-intent",
		"llm-has-anti-patterns",
		"llm-has-allowed-contexts",
		"tokens-have-descriptions",
		"token-references-valid",
		"no-orphan-tokens",
		"component-uses-valid",
		"accessibility-defined",
		"theming-contract-valid",
		"validators-configured",
		"validator-tool-required",
		"validator-type-valid",
	}

	for _, rule := range expectedRules {
		if _, ok := rules[rule]; !ok {
			t.Errorf("Expected rule %s in available rules", rule)
		}
	}
}

func TestLintSpec_ValidatorsConfigured_MissingValidator(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Accessibility: &Accessibility{
			WCAGLevel: "AA",
		},
		// No Validators configured
	}

	result := ds.Lint()

	// Should have warning for missing accessibility validator
	hasValidatorWarning := false
	for _, issue := range result.Issues {
		if issue.Rule == "validators-configured" {
			hasValidatorWarning = true
		}
	}
	if !hasValidatorWarning {
		t.Error("Expected warning for accessibility requirements without validator")
	}
}

func TestLintSpec_ValidatorsConfigured_WithValidator(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Accessibility: &Accessibility{
			WCAGLevel: "AA",
		},
		Validators: &Validators{
			Accessibility: &AccessibilityValidator{
				Tool: "agent-a11y",
				Type: ValidatorTypeMCP,
			},
		},
	}

	result := ds.Lint()

	// Should NOT have warning for missing validator
	hasValidatorWarning := false
	for _, issue := range result.Issues {
		if issue.Rule == "validators-configured" {
			hasValidatorWarning = true
		}
	}
	if hasValidatorWarning {
		t.Error("Should not have validator warning when validator is configured")
	}
}

func TestLintSpec_ValidatorToolRequired(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Validators: &Validators{
			Accessibility: &AccessibilityValidator{
				// Missing Tool
				Type: ValidatorTypeMCP,
			},
		},
	}

	result := ds.Lint()

	// Should have error for missing tool
	hasToolError := false
	for _, issue := range result.Issues {
		if issue.Rule == "validator-tool-required" && issue.Severity == "error" {
			hasToolError = true
		}
	}
	if !hasToolError {
		t.Error("Expected error for validator missing tool name")
	}
}

func TestLintSpec_ValidatorTypeValid(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Validators: &Validators{
			Accessibility: &AccessibilityValidator{
				Tool: "agent-a11y",
				Type: ValidatorType("invalid"), // Invalid type
			},
		},
	}

	result := ds.Lint()

	// Should have error for invalid type
	hasTypeError := false
	for _, issue := range result.Issues {
		if issue.Rule == "validator-type-valid" && issue.Severity == "error" {
			hasTypeError = true
		}
	}
	if !hasTypeError {
		t.Error("Expected error for invalid validator type")
	}
}

func TestLintSpec_ValidatorAllValid(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{
			Name:    "Test System",
			Version: "1.0.0",
		},
		Accessibility: &Accessibility{
			WCAGLevel: "AA",
		},
		Validators: &Validators{
			Accessibility: &AccessibilityValidator{
				Tool:      "agent-a11y",
				Type:      ValidatorTypeMCP,
				Command:   "agent-a11y mcp serve",
				Standards: []string{"WCAG2.1-AA"},
				Required:  true,
			},
			APIStyle: &APIStyleValidator{
				Tool:    "spectral",
				Type:    ValidatorTypeCLI,
				Command: "npx spectral lint",
				Ruleset: ".spectral.yaml",
			},
			Custom: []CustomValidator{
				{
					ID:     "security",
					Name:   "Security Scanner",
					Domain: "security",
					Tool:   "snyk",
					Type:   ValidatorTypeCLI,
				},
			},
		},
	}

	result := ds.Lint()

	// Should have no validator-related errors
	for _, issue := range result.Issues {
		if issue.Rule == "validators-configured" ||
			issue.Rule == "validator-tool-required" ||
			issue.Rule == "validator-type-valid" {
			t.Errorf("Unexpected validator issue: %s - %s", issue.Rule, issue.Message)
		}
	}
}

func TestModeCompleteness(t *testing.T) {
	ds := &DesignSystem{
		Meta:  Meta{Name: "Test", Version: "1.0.0"},
		Modes: []string{"light", "dark"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "complete", Value: "#0af", Modes: map[string]string{"light": "#0df", "dark": "#068"}},
				{ID: "partial", Value: "#f00", DarkModeValue: "#800"},
				{ID: "undeclared", Value: "#0f0", Modes: map[string]string{"light": "#afa", "dark": "#050", "sepia": "#dd0"}},
				{ID: "modeless", Value: "#123"},
			},
		},
	}
	result := ds.Lint()

	var partialMissing, undeclaredMode, wrongFlags int
	for _, issue := range result.Issues {
		if issue.Rule != "mode-completeness" {
			continue
		}
		switch {
		case containsAll(issue.Message, "'partial'", "missing declared mode 'light'"):
			partialMissing++
		case containsAll(issue.Message, "'undeclared'", "undeclared mode 'sepia'"):
			undeclaredMode++
		case containsAll(issue.Message, "'complete'") || containsAll(issue.Message, "'modeless'"):
			wrongFlags++
		}
	}
	if partialMissing != 1 || undeclaredMode != 1 || wrongFlags != 0 {
		t.Errorf("mode-completeness issues wrong: partial=%d undeclared=%d wrongFlags=%d (%+v)",
			partialMissing, undeclaredMode, wrongFlags, result.Issues)
	}
}

func TestModeCompletenessUndeclaredDocument(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Test", Version: "1.0.0"},
		Foundations: Foundations{
			Colors: []ColorToken{{ID: "c", Value: "#000", DarkModeValue: "#111"}},
		},
	}
	result := ds.Lint()
	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "mode-completeness" &&
			containsAll(issue.Message, "declares no modes list") && issue.Severity == "info" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected info about missing modes declaration, got %+v", result.Issues)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
