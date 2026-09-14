package dss

import (
	"context"
	"fmt"
	"strings"
)

// LintResult contains all lint findings for a design system spec.
type LintResult struct {
	// Score is the completeness score (0-100)
	Score int `json:"score"`

	// Issues contains all lint findings
	Issues []LintIssue `json:"issues"`

	// Summary provides aggregate counts
	Summary LintSummary `json:"summary"`

	// Coverage provides completeness metrics
	Coverage LintCoverage `json:"coverage"`
}

// LintIssue represents a single lint finding.
type LintIssue struct {
	// Path is the JSON path to the issue (e.g., "components[0].llm.intent")
	Path string `json:"path"`

	// Severity is "error", "warning", or "info"
	Severity string `json:"severity"`

	// Rule is the rule ID (e.g., "llm-context-required")
	Rule string `json:"rule"`

	// Message describes the issue
	Message string `json:"message"`

	// Component is the component ID if applicable
	Component string `json:"component,omitempty"`

	// Suggestion provides guidance on how to fix
	Suggestion string `json:"suggestion,omitempty"`
}

// LintSummary provides aggregate counts.
type LintSummary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Infos    int `json:"infos"`
}

// LintCoverage provides completeness metrics.
type LintCoverage struct {
	// ComponentsWithLLMContext is the percentage of components with LLM context
	ComponentsWithLLMContext float64 `json:"componentsWithLlmContext"`

	// ComponentsWithVariants is the percentage of components with variants
	ComponentsWithVariants float64 `json:"componentsWithVariants"`

	// ComponentsWithProps is the percentage of components with props
	ComponentsWithProps float64 `json:"componentsWithProps"`

	// TokensWithDescriptions is the percentage of tokens with descriptions
	TokensWithDescriptions float64 `json:"tokensWithDescriptions"`

	// TokensReferenced is the percentage of tokens referenced by components
	TokensReferenced float64 `json:"tokensReferenced"`
}

// LintOptions configures lint behavior.
type LintOptions struct {
	// Rules limits linting to specific rules (empty = all)
	Rules []string `json:"rules,omitempty"`

	// MinScore is the minimum acceptable score (default: 0)
	MinScore int `json:"minScore,omitempty"`

	// IncludeSuggestions includes fix suggestions in output
	IncludeSuggestions bool `json:"includeSuggestions,omitempty"`
}

// DefaultLintOptions returns sensible defaults.
func DefaultLintOptions() *LintOptions {
	return &LintOptions{
		MinScore:           0,
		IncludeSuggestions: true,
	}
}

// LintSpec checks the design system spec for completeness and best practices.
func (s *Service) LintSpec(ctx context.Context, opts *LintOptions) *LintResult {
	if opts == nil {
		opts = DefaultLintOptions()
	}

	result := &LintResult{
		Issues: []LintIssue{},
	}

	linter := &specLinter{
		ds:     s.ds,
		opts:   opts,
		result: result,
	}

	linter.lint(ctx)
	linter.calculateCoverage()
	linter.calculateScore()
	linter.updateSummary()

	return result
}

// Lint is a convenience method on DesignSystem.
func (ds *DesignSystem) Lint() *LintResult {
	service := NewService(ds)
	return service.LintSpec(context.Background(), nil)
}

// specLinter handles linting of a design system spec.
type specLinter struct {
	ds     *DesignSystem
	opts   *LintOptions
	result *LintResult
}

func (l *specLinter) lint(_ context.Context) {
	rules := l.opts.Rules
	allRules := len(rules) == 0

	// Meta validation
	if allRules || l.hasRule(rules, "meta-required") {
		l.checkMetaRequired()
	}

	// Mode validation
	if allRules || l.hasRule(rules, "mode-completeness") {
		l.checkModeCompleteness()
	}

	// Component validation
	if allRules || l.hasRule(rules, "component-has-variants") {
		l.checkComponentVariants()
	}
	if allRules || l.hasRule(rules, "component-has-props") {
		l.checkComponentProps()
	}
	if allRules || l.hasRule(rules, "component-has-llm-context") {
		l.checkComponentLLMContext()
	}
	if allRules || l.hasRule(rules, "llm-has-intent") {
		l.checkLLMIntent()
	}
	if allRules || l.hasRule(rules, "llm-has-anti-patterns") {
		l.checkLLMAntiPatterns()
	}
	if allRules || l.hasRule(rules, "llm-has-allowed-contexts") {
		l.checkLLMAllowedContexts()
	}

	// Token validation
	if allRules || l.hasRule(rules, "tokens-have-descriptions") {
		l.checkTokenDescriptions()
	}
	if allRules || l.hasRule(rules, "token-references-valid") {
		l.checkTokenReferences()
	}
	if allRules || l.hasRule(rules, "no-orphan-tokens") {
		l.checkOrphanTokens()
	}

	// Cross-reference validation
	if allRules || l.hasRule(rules, "component-uses-valid") {
		l.checkComponentUsesReferences()
	}

	// Accessibility validation
	if allRules || l.hasRule(rules, "accessibility-defined") {
		l.checkAccessibilityDefined()
	}

	// Theming validation
	if allRules || l.hasRule(rules, "theming-contract-valid") {
		l.checkThemingContracts()
	}

	// Validator configuration
	if allRules || l.hasRule(rules, "validators-configured") {
		l.checkValidatorsConfigured()
	}
	if allRules || l.hasRule(rules, "validator-tool-required") {
		l.checkValidatorToolRequired()
	}
	if allRules || l.hasRule(rules, "validator-type-valid") {
		l.checkValidatorTypeValid()
	}
}

func (l *specLinter) hasRule(rules []string, wanted ...string) bool {
	for _, r := range rules {
		for _, w := range wanted {
			if r == w {
				return true
			}
		}
	}
	return false
}

func (l *specLinter) addIssue(path, rule, message, severity string, component string) {
	issue := LintIssue{
		Path:      path,
		Rule:      rule,
		Message:   message,
		Severity:  severity,
		Component: component,
	}

	if l.opts.IncludeSuggestions {
		issue.Suggestion = l.getSuggestion(rule)
	}

	l.result.Issues = append(l.result.Issues, issue)
}

func (l *specLinter) getSuggestion(rule string) string {
	suggestions := map[string]string{
		"mode-completeness":         "Declare the document's modes list and give every mode-aware token a value for each declared mode",
		"meta-required":             "Add name and version to meta.json",
		"component-has-variants":    "Define variants array with at least one variant",
		"component-has-props":       "Define props array with component properties",
		"component-has-llm-context": "Add llm field with intent, allowedContexts, and antiPatterns",
		"llm-has-intent":            "Add intent field describing when to use this component",
		"llm-has-anti-patterns":     "Add antiPatterns array with common misuse patterns",
		"llm-has-allowed-contexts":  "Add allowedContexts array with valid usage contexts",
		"tokens-have-descriptions":  "Add description field to token definitions",
		"token-references-valid":    "Ensure tokensUsed references valid token IDs",
		"no-orphan-tokens":          "Reference token in component tokensUsed or remove if unused",
		"component-uses-valid":      "Ensure uses array references valid component IDs",
		"accessibility-defined":     "Add accessibility section with WCAG level and requirements",
		"theming-contract-valid":    "Ensure themingContract has prefix and tokens defined",
		"validators-configured":     "Add validators section with external validation tools (e.g., agent-a11y for accessibility)",
		"validator-tool-required":   "Specify tool name for validator (e.g., 'agent-a11y', 'spectral')",
		"validator-type-valid":      "Use valid validator type: mcp, cli, npm, api, or library",
	}
	return suggestions[rule]
}

// Meta checks
func (l *specLinter) checkMetaRequired() {
	if l.ds.Meta.Name == "" {
		l.addIssue("meta.name", "meta-required", "Design system name is required", "error", "")
	}
	if l.ds.Meta.Version == "" {
		l.addIssue("meta.version", "meta-required", "Design system version is required", "error", "")
	}
	if l.ds.Meta.Description == "" {
		l.addIssue("meta.description", "meta-required", "Design system description is recommended", "info", "")
	}
}

// Component checks
func (l *specLinter) checkComponentVariants() {
	for i, c := range l.ds.Components {
		if len(c.Variants) == 0 {
			l.addIssue(
				fmt.Sprintf("components[%d].variants", i),
				"component-has-variants",
				fmt.Sprintf("Component '%s' has no variants defined", c.ID),
				"warning",
				c.ID,
			)
		}
	}
}

func (l *specLinter) checkComponentProps() {
	for i, c := range l.ds.Components {
		if len(c.Props) == 0 {
			l.addIssue(
				fmt.Sprintf("components[%d].props", i),
				"component-has-props",
				fmt.Sprintf("Component '%s' has no props defined", c.ID),
				"info",
				c.ID,
			)
		}
	}
}

func (l *specLinter) checkComponentLLMContext() {
	for i, c := range l.ds.Components {
		if c.LLM == nil {
			l.addIssue(
				fmt.Sprintf("components[%d].llm", i),
				"component-has-llm-context",
				fmt.Sprintf("Component '%s' missing LLM context for AI code generation", c.ID),
				"warning",
				c.ID,
			)
		}
	}
}

func (l *specLinter) checkLLMIntent() {
	for i, c := range l.ds.Components {
		if c.LLM != nil && c.LLM.Intent == "" {
			l.addIssue(
				fmt.Sprintf("components[%d].llm.intent", i),
				"llm-has-intent",
				fmt.Sprintf("Component '%s' LLM context missing intent", c.ID),
				"error",
				c.ID,
			)
		}
	}
}

func (l *specLinter) checkLLMAntiPatterns() {
	for i, c := range l.ds.Components {
		if c.LLM != nil && len(c.LLM.AntiPatterns) == 0 {
			l.addIssue(
				fmt.Sprintf("components[%d].llm.antiPatterns", i),
				"llm-has-anti-patterns",
				fmt.Sprintf("Component '%s' LLM context should document anti-patterns", c.ID),
				"warning",
				c.ID,
			)
		}
	}
}

func (l *specLinter) checkLLMAllowedContexts() {
	for i, c := range l.ds.Components {
		if c.LLM != nil && len(c.LLM.AllowedContexts) == 0 {
			l.addIssue(
				fmt.Sprintf("components[%d].llm.allowedContexts", i),
				"llm-has-allowed-contexts",
				fmt.Sprintf("Component '%s' LLM context should specify allowed contexts", c.ID),
				"info",
				c.ID,
			)
		}
	}
}

// Token checks
func (l *specLinter) checkTokenDescriptions() {
	// Check colors - ColorToken uses "Usage" field for description
	for i, c := range l.ds.Foundations.Colors {
		if c.Usage == "" {
			l.addIssue(
				fmt.Sprintf("foundations.colors[%d].usage", i),
				"tokens-have-descriptions",
				fmt.Sprintf("Color token '%s' missing usage description", c.ID),
				"info",
				"",
			)
		}
	}

	// Check elevation - ElevationToken has Usage field
	for i, e := range l.ds.Foundations.Elevation {
		if e.Usage == "" {
			l.addIssue(
				fmt.Sprintf("foundations.elevation[%d].usage", i),
				"tokens-have-descriptions",
				fmt.Sprintf("Elevation token '%s' missing usage description", e.ID),
				"info",
				"",
			)
		}
	}
}

// checkModeCompleteness verifies that per-mode token values line up with the
// document's declared modes (DesignSystem.Modes): tokens that define some
// mode values should cover every declared mode, mode keys should be
// declared, and documents using modes should declare them.
func (l *specLinter) checkModeCompleteness() {
	declared := map[string]bool{}
	for _, m := range l.ds.Modes {
		declared[m] = true
	}

	anyModeUse := false
	for i, c := range l.ds.Foundations.Colors {
		modes := c.EffectiveModes()
		if len(modes) == 0 {
			continue
		}
		anyModeUse = true

		if len(declared) > 0 {
			for _, m := range l.ds.Modes {
				if _, ok := modes[m]; !ok {
					l.addIssue(
						fmt.Sprintf("foundations.colors[%d].modes", i),
						"mode-completeness",
						fmt.Sprintf("Color token '%s' defines mode values but is missing declared mode '%s'", c.ID, m),
						"warning",
						"",
					)
				}
			}
		}
		for m := range modes {
			if len(declared) > 0 && !declared[m] {
				l.addIssue(
					fmt.Sprintf("foundations.colors[%d].modes", i),
					"mode-completeness",
					fmt.Sprintf("Color token '%s' defines value for undeclared mode '%s' (declare it in the document's modes list)", c.ID, m),
					"warning",
					"",
				)
			}
		}
	}

	if anyModeUse && len(declared) == 0 {
		l.addIssue(
			"modes",
			"mode-completeness",
			"Tokens define per-mode values but the document declares no modes list",
			"info",
			"",
		)
	}
}

func (l *specLinter) checkTokenReferences() {
	// Build set of valid token IDs
	validTokens := l.buildTokenSet()

	// Check component tokensUsed references
	for i, c := range l.ds.Components {
		for j, tokenRef := range c.TokensUsed {
			if !validTokens[tokenRef] {
				l.addIssue(
					fmt.Sprintf("components[%d].tokensUsed[%d]", i, j),
					"token-references-valid",
					fmt.Sprintf("Component '%s' references unknown token '%s'", c.ID, tokenRef),
					"error",
					c.ID,
				)
			}
		}
	}
}

func (l *specLinter) checkOrphanTokens() {
	// Build set of referenced tokens
	referencedTokens := make(map[string]bool)
	for _, c := range l.ds.Components {
		for _, tokenRef := range c.TokensUsed {
			referencedTokens[tokenRef] = true
		}
	}

	// Check for orphan color tokens
	for i, c := range l.ds.Foundations.Colors {
		if !referencedTokens[c.ID] {
			l.addIssue(
				fmt.Sprintf("foundations.colors[%d]", i),
				"no-orphan-tokens",
				fmt.Sprintf("Color token '%s' is not referenced by any component", c.ID),
				"info",
				"",
			)
		}
	}
}

func (l *specLinter) buildTokenSet() map[string]bool {
	tokens := make(map[string]bool)

	// Colors
	for _, c := range l.ds.Foundations.Colors {
		tokens[c.ID] = true
	}

	// Spacing
	if l.ds.Foundations.Spacing != nil {
		for _, s := range l.ds.Foundations.Spacing.Scale {
			tokens[s.ID] = true
		}
	}

	// Typography
	if l.ds.Foundations.Typography != nil {
		for _, f := range l.ds.Foundations.Typography.FontFamilies {
			tokens[f.ID] = true
		}
		for _, s := range l.ds.Foundations.Typography.FontSizes {
			tokens[s.ID] = true
		}
	}

	// Elevation
	for _, e := range l.ds.Foundations.Elevation {
		tokens[e.ID] = true
	}

	// Border radius
	for _, b := range l.ds.Foundations.BorderRadius {
		tokens[b.ID] = true
	}

	return tokens
}

// Cross-reference checks
func (l *specLinter) checkComponentUsesReferences() {
	// Build set of valid component IDs
	validComponents := make(map[string]bool)
	for _, c := range l.ds.Components {
		validComponents[c.ID] = true
	}

	// Check uses references
	for i, c := range l.ds.Components {
		for j, usedID := range c.Uses {
			if !validComponents[usedID] {
				l.addIssue(
					fmt.Sprintf("components[%d].uses[%d]", i, j),
					"component-uses-valid",
					fmt.Sprintf("Component '%s' references unknown component '%s' in uses", c.ID, usedID),
					"error",
					c.ID,
				)
			}
		}
	}
}

// Accessibility checks
func (l *specLinter) checkAccessibilityDefined() {
	if l.ds.Accessibility == nil {
		l.addIssue(
			"accessibility",
			"accessibility-defined",
			"Design system should define accessibility requirements",
			"warning",
			"",
		)
		return
	}

	if l.ds.Accessibility.WCAGLevel == "" {
		l.addIssue(
			"accessibility.wcagLevel",
			"accessibility-defined",
			"Accessibility should specify WCAG target level (A, AA, AAA)",
			"warning",
			"",
		)
	}
}

// Theming checks
func (l *specLinter) checkThemingContracts() {
	for i, c := range l.ds.Components {
		if c.ThemingContract != nil {
			if c.ThemingContract.Prefix == "" {
				l.addIssue(
					fmt.Sprintf("components[%d].themingContract.prefix", i),
					"theming-contract-valid",
					fmt.Sprintf("Component '%s' theming contract missing prefix", c.ID),
					"error",
					c.ID,
				)
			} else if !strings.HasPrefix(c.ThemingContract.Prefix, "--") {
				l.addIssue(
					fmt.Sprintf("components[%d].themingContract.prefix", i),
					"theming-contract-valid",
					fmt.Sprintf("Component '%s' theming contract prefix should start with '--'", c.ID),
					"warning",
					c.ID,
				)
			}

			if len(c.ThemingContract.Tokens) == 0 {
				l.addIssue(
					fmt.Sprintf("components[%d].themingContract.tokens", i),
					"theming-contract-valid",
					fmt.Sprintf("Component '%s' theming contract has no tokens", c.ID),
					"warning",
					c.ID,
				)
			}
		}
	}
}

// Validator checks
func (l *specLinter) checkValidatorsConfigured() {
	// If accessibility requirements are defined, check for accessibility validator
	if l.ds.Accessibility != nil && l.ds.Accessibility.WCAGLevel != "" {
		if l.ds.Validators == nil || l.ds.Validators.Accessibility == nil {
			l.addIssue(
				"validators.accessibility",
				"validators-configured",
				fmt.Sprintf("Accessibility requirements defined (WCAG %s) but no validator configured", l.ds.Accessibility.WCAGLevel),
				"warning",
				"",
			)
		}
	}
}

func (l *specLinter) checkValidatorToolRequired() {
	if l.ds.Validators == nil {
		return
	}

	// Check accessibility validator
	if v := l.ds.Validators.Accessibility; v != nil {
		if v.Tool == "" {
			l.addIssue(
				"validators.accessibility.tool",
				"validator-tool-required",
				"Accessibility validator missing tool name",
				"error",
				"",
			)
		}
	}

	// Check API style validator
	if v := l.ds.Validators.APIStyle; v != nil {
		if v.Tool == "" {
			l.addIssue(
				"validators.apiStyle.tool",
				"validator-tool-required",
				"API style validator missing tool name",
				"error",
				"",
			)
		}
	}

	// Check custom validators
	for i, v := range l.ds.Validators.Custom {
		if v.Tool == "" {
			l.addIssue(
				fmt.Sprintf("validators.custom[%d].tool", i),
				"validator-tool-required",
				fmt.Sprintf("Custom validator '%s' missing tool name", v.ID),
				"error",
				"",
			)
		}
	}
}

func (l *specLinter) checkValidatorTypeValid() {
	if l.ds.Validators == nil {
		return
	}

	validTypes := map[ValidatorType]bool{
		ValidatorTypeMCP:     true,
		ValidatorTypeCLI:     true,
		ValidatorTypeNPM:     true,
		ValidatorTypeAPI:     true,
		ValidatorTypeLibrary: true,
	}

	// Check accessibility validator
	if v := l.ds.Validators.Accessibility; v != nil {
		if v.Type != "" && !validTypes[v.Type] {
			l.addIssue(
				"validators.accessibility.type",
				"validator-type-valid",
				fmt.Sprintf("Accessibility validator has invalid type '%s'", v.Type),
				"error",
				"",
			)
		}
	}

	// Check API style validator
	if v := l.ds.Validators.APIStyle; v != nil {
		if v.Type != "" && !validTypes[v.Type] {
			l.addIssue(
				"validators.apiStyle.type",
				"validator-type-valid",
				fmt.Sprintf("API style validator has invalid type '%s'", v.Type),
				"error",
				"",
			)
		}
	}

	// Check custom validators
	for i, v := range l.ds.Validators.Custom {
		if v.Type != "" && !validTypes[v.Type] {
			l.addIssue(
				fmt.Sprintf("validators.custom[%d].type", i),
				"validator-type-valid",
				fmt.Sprintf("Custom validator '%s' has invalid type '%s'", v.ID, v.Type),
				"error",
				"",
			)
		}
	}
}

// Coverage calculation
func (l *specLinter) calculateCoverage() {
	totalComponents := len(l.ds.Components)
	if totalComponents == 0 {
		return
	}

	// Components with LLM context
	withLLM := 0
	for _, c := range l.ds.Components {
		if c.LLM != nil {
			withLLM++
		}
	}
	l.result.Coverage.ComponentsWithLLMContext = float64(withLLM) / float64(totalComponents) * 100

	// Components with variants
	withVariants := 0
	for _, c := range l.ds.Components {
		if len(c.Variants) > 0 {
			withVariants++
		}
	}
	l.result.Coverage.ComponentsWithVariants = float64(withVariants) / float64(totalComponents) * 100

	// Components with props
	withProps := 0
	for _, c := range l.ds.Components {
		if len(c.Props) > 0 {
			withProps++
		}
	}
	l.result.Coverage.ComponentsWithProps = float64(withProps) / float64(totalComponents) * 100

	// Tokens with descriptions (using Usage field)
	totalTokens := len(l.ds.Foundations.Colors)
	tokensWithDesc := 0
	for _, c := range l.ds.Foundations.Colors {
		if c.Usage != "" {
			tokensWithDesc++
		}
	}
	if totalTokens > 0 {
		l.result.Coverage.TokensWithDescriptions = float64(tokensWithDesc) / float64(totalTokens) * 100
	}

	// Tokens referenced
	referencedTokens := make(map[string]bool)
	for _, c := range l.ds.Components {
		for _, t := range c.TokensUsed {
			referencedTokens[t] = true
		}
	}
	allTokens := l.buildTokenSet()
	if len(allTokens) > 0 {
		l.result.Coverage.TokensReferenced = float64(len(referencedTokens)) / float64(len(allTokens)) * 100
	}
}

// Score calculation
func (l *specLinter) calculateScore() {
	// Start with 100
	score := 100

	// Deduct points for issues
	for _, issue := range l.result.Issues {
		switch issue.Severity {
		case "error":
			score -= 10
		case "warning":
			score -= 3
		case "info":
			score -= 1
		}
	}

	// Add bonus for good coverage
	coverage := l.result.Coverage
	avgCoverage := (coverage.ComponentsWithLLMContext +
		coverage.ComponentsWithVariants +
		coverage.ComponentsWithProps) / 3

	if avgCoverage >= 80 {
		score += 5
	}

	// Clamp to 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	l.result.Score = score
}

func (l *specLinter) updateSummary() {
	for _, issue := range l.result.Issues {
		switch issue.Severity {
		case "error":
			l.result.Summary.Errors++
		case "warning":
			l.result.Summary.Warnings++
		case "info":
			l.result.Summary.Infos++
		}
	}
}

// AvailableLintRules returns all available lint rule IDs with descriptions.
func AvailableLintRules() map[string]string {
	return map[string]string{
		"meta-required":             "Design system must have name and version",
		"mode-completeness":         "Mode-aware tokens must cover every declared mode",
		"component-has-variants":    "Components should define variants",
		"component-has-props":       "Components should define props",
		"component-has-llm-context": "Components should have LLM context for AI generation",
		"llm-has-intent":            "LLM context must have intent field",
		"llm-has-anti-patterns":     "LLM context should document anti-patterns",
		"llm-has-allowed-contexts":  "LLM context should specify allowed contexts",
		"tokens-have-descriptions":  "Tokens should have descriptions",
		"token-references-valid":    "Token references must resolve to valid tokens",
		"no-orphan-tokens":          "Tokens should be referenced by components",
		"component-uses-valid":      "Component uses references must be valid",
		"accessibility-defined":     "Design system should define accessibility requirements",
		"theming-contract-valid":    "Theming contracts must be properly configured",
		"validators-configured":     "External validators should be configured for requirements",
		"validator-tool-required":   "Validators must specify a tool name",
		"validator-type-valid":      "Validator type must be mcp, cli, npm, api, or library",
	}
}
