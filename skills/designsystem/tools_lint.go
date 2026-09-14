package designsystem

import (
	"context"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func (s *Skill) lintSpecTool() skill.Tool {
	return skill.NewTool(
		"lint_spec",
		"Lint the design system spec for completeness and best practices. Returns a score (0-100), issues, and coverage metrics.",
		skill.Parameters{
			"rules": {
				Type:        "array",
				Description: "Specific rules to check (default: all). Available: meta-required, component-has-variants, component-has-props, component-has-llm-context, llm-has-intent, llm-has-anti-patterns, llm-has-allowed-contexts, tokens-have-descriptions, token-references-valid, no-orphan-tokens, component-uses-valid, accessibility-defined, theming-contract-valid, validators-configured, validator-tool-required, validator-type-valid",
				Required:    false,
			},
			"min_score": {
				Type:        "number",
				Description: "Minimum acceptable score (0-100). Returns error if score is below this threshold.",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			opts := &dss.LintOptions{
				IncludeSuggestions: true,
			}

			if rules, ok := params["rules"].([]any); ok {
				for _, r := range rules {
					if rs, ok := r.(string); ok {
						opts.Rules = append(opts.Rules, rs)
					}
				}
			}

			if minScore, ok := params["min_score"].(float64); ok {
				opts.MinScore = int(minScore)
			}

			result := s.service.LintSpec(ctx, opts)

			return map[string]any{
				"score":    result.Score,
				"issues":   result.Issues,
				"summary":  result.Summary,
				"coverage": result.Coverage,
				"passed":   result.Score >= opts.MinScore,
			}, nil
		},
	)
}

func (s *Skill) listLintRulesTool() skill.Tool {
	return skill.NewTool(
		"list_lint_rules",
		"List all available lint rules with descriptions.",
		skill.Parameters{},
		func(ctx context.Context, params map[string]any) (any, error) {
			rules := dss.AvailableLintRules()

			// Convert to array format for better readability
			ruleList := make([]map[string]string, 0, len(rules))
			for id, desc := range rules {
				ruleList = append(ruleList, map[string]string{
					"id":          id,
					"description": desc,
				})
			}

			return map[string]any{
				"rules": ruleList,
				"count": len(rules),
			}, nil
		},
	)
}

func (s *Skill) checkAgentReadinessTool() skill.Tool {
	return skill.NewTool(
		"check_agent_readiness",
		"Check if the design system spec is ready for AI agent code generation. Verifies that components have LLM context with intent, anti-patterns, and allowed contexts. Also checks that external validators are configured for requirements delegation.",
		skill.Parameters{
			"component_id": {
				Type:        "string",
				Description: "Optional: Check readiness for a specific component. If not provided, checks all components.",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			// Run lint with agent-readiness rules (LLM context + validators)
			opts := &dss.LintOptions{
				Rules: []string{
					"component-has-llm-context",
					"llm-has-intent",
					"llm-has-anti-patterns",
					"llm-has-allowed-contexts",
					"validators-configured",
				},
				IncludeSuggestions: true,
			}

			result := s.service.LintSpec(ctx, opts)

			// Filter by component if specified
			componentID := ""
			if id, ok := params["component_id"].(string); ok {
				componentID = id
			}

			var filteredIssues []dss.LintIssue
			if componentID != "" {
				for _, issue := range result.Issues {
					if issue.Component == componentID {
						filteredIssues = append(filteredIssues, issue)
					}
				}
			} else {
				filteredIssues = result.Issues
			}

			// Calculate agent readiness score based on LLM context coverage
			readiness := result.Coverage.ComponentsWithLLMContext
			ready := readiness >= 80 && result.Summary.Errors == 0

			return map[string]any{
				"ready":                ready,
				"readiness_score":      readiness,
				"llm_context_coverage": result.Coverage.ComponentsWithLLMContext,
				"issues":               filteredIssues,
				"errors":               result.Summary.Errors,
				"warnings":             result.Summary.Warnings,
				"recommendation":       getReadinessRecommendation(readiness, result.Summary.Errors),
			}, nil
		},
	)
}

func getReadinessRecommendation(readiness float64, errors int) string {
	if errors > 0 {
		return "Fix LLM context errors before using with AI agents. Components with LLM context must have 'intent' defined."
	}
	if readiness < 50 {
		return "Low agent readiness. Add LLM context (intent, allowedContexts, antiPatterns) to more components."
	}
	if readiness < 80 {
		return "Moderate agent readiness. Consider adding LLM context to remaining components for best results."
	}
	return "Good agent readiness. Components have LLM context for AI code generation."
}
