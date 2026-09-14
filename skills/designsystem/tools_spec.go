package designsystem

import (
	"context"
	"fmt"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
)

func (s *Skill) getComponentTool() skill.Tool {
	return skill.NewTool(
		"get_component",
		"Get a component definition including variants, props, states, events, slots, accessibility requirements, and LLM context",
		skill.Parameters{
			"id": {
				Type:        "string",
				Description: "Component ID (e.g., 'button', 'input', 'modal')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			id, ok := params["id"].(string)
			if !ok {
				return nil, fmt.Errorf("id parameter is required and must be a string")
			}
			return s.service.GetComponent(ctx, id)
		},
	)
}

func (s *Skill) listComponentsTool() skill.Tool {
	return skill.NewTool(
		"list_components",
		"List all components in the design system with ID, name, category, and description",
		skill.Parameters{},
		func(ctx context.Context, _ map[string]any) (any, error) {
			return s.service.ListComponents(ctx), nil
		},
	)
}

func (s *Skill) getTokenTool() skill.Tool {
	return skill.NewTool(
		"get_token",
		"Get a design token value (color, spacing, typography, elevation, etc.)",
		skill.Parameters{
			"type": {
				Type:        "string",
				Description: "Token type: color, spacing, typography, elevation, borderRadius, breakpoint",
				Required:    true,
				Enum:        []any{"color", "spacing", "typography", "elevation", "borderRadius", "breakpoint"},
			},
			"name": {
				Type:        "string",
				Description: "Token name/ID (e.g., 'primary-500', 'spacing-4', 'h1')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			tokenType, ok := params["type"].(string)
			if !ok {
				return nil, fmt.Errorf("type parameter is required")
			}
			name, ok := params["name"].(string)
			if !ok {
				return nil, fmt.Errorf("name parameter is required")
			}
			return s.service.GetToken(ctx, tokenType, name)
		},
	)
}

func (s *Skill) listTokensTool() skill.Tool {
	return skill.NewTool(
		"list_tokens",
		"List all tokens of a given type with ID and value",
		skill.Parameters{
			"type": {
				Type:        "string",
				Description: "Token type: color, spacing, typography, elevation, borderRadius, breakpoint",
				Required:    true,
				Enum:        []any{"color", "spacing", "typography", "elevation", "borderRadius", "breakpoint"},
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			tokenType, ok := params["type"].(string)
			if !ok {
				return nil, fmt.Errorf("type parameter is required")
			}
			return s.service.ListTokens(ctx, tokenType)
		},
	)
}

func (s *Skill) getPatternTool() skill.Tool {
	return skill.NewTool(
		"get_pattern",
		"Get a pattern definition including components, layout, variations, and accessibility requirements",
		skill.Parameters{
			"id": {
				Type:        "string",
				Description: "Pattern ID (e.g., 'form-validation', 'data-table')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			id, ok := params["id"].(string)
			if !ok {
				return nil, fmt.Errorf("id parameter is required")
			}
			return s.service.GetPattern(ctx, id)
		},
	)
}

func (s *Skill) listPatternsTool() skill.Tool {
	return skill.NewTool(
		"list_patterns",
		"List all patterns in the design system with ID, name, category, and description",
		skill.Parameters{},
		func(ctx context.Context, _ map[string]any) (any, error) {
			return s.service.ListPatterns(ctx), nil
		},
	)
}

func (s *Skill) getMetaTool() skill.Tool {
	return skill.NewTool(
		"get_meta",
		"Get design system metadata including name, version, description, maturity level, and maintainers",
		skill.Parameters{},
		func(ctx context.Context, _ map[string]any) (any, error) {
			return s.service.GetMeta(ctx), nil
		},
	)
}
