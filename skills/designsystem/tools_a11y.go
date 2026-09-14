package designsystem

import (
	"context"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
)

// getAccessibilityRequirementsTool returns the tool for getting accessibility requirements.
func (s *Skill) getAccessibilityRequirementsTool() skill.Tool {
	return skill.NewTool(
		"get_accessibility_requirements",
		"Get accessibility requirements for a component including required props, keyboard interactions, focus management, and WCAG criteria",
		skill.Parameters{
			"component": {
				Type:        "string",
				Description: "Component ID (e.g., 'button', 'input', 'modal')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			componentID, _ := params["component"].(string)
			return s.service.GetAccessibilityRequirements(ctx, componentID)
		},
	)
}

// getA11yAntiPatternsTool returns the tool for getting accessibility anti-patterns.
func (s *Skill) getA11yAntiPatternsTool() skill.Tool {
	return skill.NewTool(
		"get_a11y_anti_patterns",
		"Get accessibility anti-patterns to avoid for a component or rule. Returns bad examples, good examples, and WCAG criteria.",
		skill.Parameters{
			"component": {
				Type:        "string",
				Description: "Component ID to get anti-patterns for (optional)",
				Required:    false,
			},
			"rule_id": {
				Type:        "string",
				Description: "Rule ID like 'color-contrast', 'missing-label', 'keyboard', 'focus' (optional)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			componentID, _ := params["component"].(string)
			ruleID, _ := params["rule_id"].(string)
			return s.service.GetAntiPatterns(ctx, componentID, ruleID)
		},
	)
}

// suggestContrastTokenTool returns the tool for suggesting contrast-compliant tokens.
func (s *Skill) suggestContrastTokenTool() skill.Tool {
	return skill.NewTool(
		"suggest_contrast_token",
		"Suggest color tokens that meet contrast requirements against a background color",
		skill.Parameters{
			"background": {
				Type:        "string",
				Description: "Background color as hex (e.g., '#ffffff', '#1a1a2e')",
				Required:    true,
			},
			"min_contrast": {
				Type:        "number",
				Description: "Minimum contrast ratio (default: 4.5 for AA normal text)",
				Required:    false,
				Default:     4.5,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			background, _ := params["background"].(string)
			minContrast := 4.5
			if mc, ok := params["min_contrast"].(float64); ok && mc > 0 {
				minContrast = mc
			}
			return s.service.SuggestContrastToken(ctx, background, minContrast)
		},
	)
}

// getComponentFixContextTool returns the tool for getting fix context.
func (s *Skill) getComponentFixContextTool() skill.Tool {
	return skill.NewTool(
		"get_component_fix_context",
		"Get full context needed to fix accessibility issues in a component including file patterns, props to add, styles to check, and available tokens",
		skill.Parameters{
			"component": {
				Type:        "string",
				Description: "Component ID (e.g., 'button', 'input')",
				Required:    true,
			},
			"issue_type": {
				Type:        "string",
				Description: "Issue type: 'color-contrast', 'missing-label', 'keyboard', 'focus'",
				Required:    true,
				Enum:        []any{"color-contrast", "missing-label", "keyboard", "focus"},
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			componentID, _ := params["component"].(string)
			issueType, _ := params["issue_type"].(string)
			return s.service.GetComponentFixContext(ctx, componentID, issueType)
		},
	)
}
