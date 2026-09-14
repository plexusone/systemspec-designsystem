package designsystem

import (
	"context"
	"fmt"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func (s *Skill) generatePromptTool() skill.Tool {
	return skill.NewTool(
		"generate_prompt",
		"Generate an LLM context prompt for the design system. Use this to get comprehensive guidance for implementing UI components.",
		skill.Parameters{
			"format": {
				Type:        "string",
				Description: "Output format: markdown or xml",
				Required:    false,
				Default:     "markdown",
				Enum:        []any{"markdown", "xml"},
			},
			"include_foundations": {
				Type:        "boolean",
				Description: "Include design tokens (colors, spacing, typography)",
				Required:    false,
				Default:     true,
			},
			"include_components": {
				Type:        "boolean",
				Description: "Include component definitions and usage",
				Required:    false,
				Default:     true,
			},
			"include_patterns": {
				Type:        "boolean",
				Description: "Include pattern recommendations",
				Required:    false,
				Default:     true,
			},
			"include_accessibility": {
				Type:        "boolean",
				Description: "Include accessibility requirements",
				Required:    false,
				Default:     true,
			},
			"include_anti_patterns": {
				Type:        "boolean",
				Description: "Include anti-patterns to avoid",
				Required:    false,
				Default:     true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			opts := dss.DefaultPromptOptions()

			if format, ok := params["format"].(string); ok {
				opts.Format = format
			}
			if v, ok := params["include_foundations"].(bool); ok {
				opts.IncludeFoundations = v
			}
			if v, ok := params["include_components"].(bool); ok {
				opts.IncludeComponents = v
			}
			if v, ok := params["include_patterns"].(bool); ok {
				opts.IncludePatterns = v
			}
			if v, ok := params["include_accessibility"].(bool); ok {
				opts.IncludeAccessibility = v
			}
			if v, ok := params["include_anti_patterns"].(bool); ok {
				opts.IncludeAntiPatterns = v
			}

			return s.service.GenerateLLMPrompt(ctx, opts)
		},
	)
}

func (s *Skill) getVariantsTool() skill.Tool {
	return skill.NewTool(
		"get_variants",
		"Get all valid variants for a component with descriptions and default information",
		skill.Parameters{
			"component_id": {
				Type:        "string",
				Description: "Component ID (e.g., 'button', 'input')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			id, ok := params["component_id"].(string)
			if !ok {
				return nil, fmt.Errorf("component_id parameter is required")
			}
			return s.service.GetComponentVariants(ctx, id)
		},
	)
}

func (s *Skill) getPropsTool() skill.Tool {
	return skill.NewTool(
		"get_props",
		"Get prop definitions for a component including types, defaults, constraints, and required flags",
		skill.Parameters{
			"component_id": {
				Type:        "string",
				Description: "Component ID (e.g., 'button', 'input')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			id, ok := params["component_id"].(string)
			if !ok {
				return nil, fmt.Errorf("component_id parameter is required")
			}
			return s.service.GetComponentProps(ctx, id)
		},
	)
}

func (s *Skill) getAntiPatternsTool() skill.Tool {
	return skill.NewTool(
		"get_anti_patterns",
		"Get anti-patterns to avoid for a component. These are common mistakes that should not be made.",
		skill.Parameters{
			"component_id": {
				Type:        "string",
				Description: "Component ID (e.g., 'button', 'input')",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			id, ok := params["component_id"].(string)
			if !ok {
				return nil, fmt.Errorf("component_id parameter is required")
			}
			antiPatterns, err := s.service.GetComponentAntiPatterns(ctx, id)
			if err != nil {
				return nil, err
			}

			// Return structured response
			return map[string]any{
				"component_id":  id,
				"anti_patterns": antiPatterns,
			}, nil
		},
	)
}
