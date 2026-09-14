package designsystem

import (
	"context"
	"fmt"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func (s *Skill) validateFileTool() skill.Tool {
	return skill.NewTool(
		"validate_file",
		"Validate a file against the design system spec. Checks for hardcoded colors, spacing, accessibility issues, and anti-patterns.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to validate",
				Required:    true,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to check (default: all). Options: no-hardcoded-colors, use-spacing-scale, img-alt-required, button-accessible-name, valid-variant, single-primary-button, no-nested-cards",
				Required:    false,
			},
			"include_context": {
				Type:        "boolean",
				Description: "Include code snippets in violation reports",
				Required:    false,
				Default:     false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, ok := params["path"].(string)
			if !ok {
				return nil, fmt.Errorf("path parameter is required")
			}

			opts := dss.DefaultValidateOptions()

			if rules, ok := params["rules"].([]any); ok {
				opts.Rules = toStringSlice(rules)
			}
			if v, ok := params["include_context"].(bool); ok {
				opts.IncludeContext = v
			}

			return s.service.ValidateFile(ctx, path, opts)
		},
	)
}

func (s *Skill) validateDirectoryTool() skill.Tool {
	return skill.NewTool(
		"validate_directory",
		"Validate all files in a directory against the design system spec",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the directory to validate",
				Required:    true,
			},
			"extensions": {
				Type:        "array",
				Description: "File extensions to check (default: .tsx, .jsx, .ts, .js, .css)",
				Required:    false,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to check (default: all)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, ok := params["path"].(string)
			if !ok {
				return nil, fmt.Errorf("path parameter is required")
			}

			opts := dss.DefaultValidateOptions()

			if extensions, ok := params["extensions"].([]any); ok {
				opts.Extensions = toStringSlice(extensions)
			}
			if rules, ok := params["rules"].([]any); ok {
				opts.Rules = toStringSlice(rules)
			}

			return s.service.ValidateDirectory(ctx, path, opts)
		},
	)
}

func (s *Skill) checkColorsTool() skill.Tool {
	return skill.NewTool(
		"check_colors",
		"Check a file for hardcoded colors that should use design tokens. Returns violations for hex colors, rgb(), hsl() values.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to check",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, ok := params["path"].(string)
			if !ok {
				return nil, fmt.Errorf("path parameter is required")
			}

			opts := &dss.ValidateOptions{
				Rules:          []string{"no-hardcoded-colors"},
				IncludeContext: true,
			}

			return s.service.ValidateFile(ctx, path, opts)
		},
	)
}

func (s *Skill) checkSpacingTool() skill.Tool {
	return skill.NewTool(
		"check_spacing",
		"Check a file for hardcoded spacing values that should use the spacing scale",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to check",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, ok := params["path"].(string)
			if !ok {
				return nil, fmt.Errorf("path parameter is required")
			}

			opts := &dss.ValidateOptions{
				Rules:          []string{"use-spacing-scale"},
				IncludeContext: true,
			}

			return s.service.ValidateFile(ctx, path, opts)
		},
	)
}

// toStringSlice converts []any to []string.
func toStringSlice(arr []any) []string {
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
