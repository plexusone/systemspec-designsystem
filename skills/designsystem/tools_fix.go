package designsystem

import (
	"context"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

// parseFixOptions extracts FixOptions from tool parameters.
func parseFixOptions(params map[string]any) *dss.FixOptions {
	opts := &dss.FixOptions{}

	if rules, ok := params["rules"].([]any); ok {
		for _, r := range rules {
			if rs, ok := r.(string); ok {
				opts.Rules = append(opts.Rules, rs)
			}
		}
	}

	if dryRun, ok := params["dry_run"].(bool); ok {
		opts.DryRun = dryRun
	}

	return opts
}

func (s *Skill) fixFileTool() skill.Tool {
	return skill.NewTool(
		"fix_file",
		"Fix design system violations in a file. Replaces hardcoded colors with CSS variables, fixes spacing values, and adds missing accessibility attributes.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to fix",
				Required:    true,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to fix (default: all). Options: no-hardcoded-colors, use-spacing-scale, img-alt-required, button-accessible-name",
				Required:    false,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "If true, return fixes without applying them (default: false)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, _ := params["path"].(string)
			return s.service.FixFile(ctx, path, parseFixOptions(params))
		},
	)
}

func (s *Skill) suggestFixesTool() skill.Tool {
	return skill.NewTool(
		"suggest_fixes",
		"Suggest fixes for design system violations without applying them. Returns the original content, suggested fixes, and what the fixed content would look like.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to analyze",
				Required:    true,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to check (default: all). Options: no-hardcoded-colors, use-spacing-scale, img-alt-required, button-accessible-name",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, _ := params["path"].(string)

			opts := &dss.FixOptions{}

			if rules, ok := params["rules"].([]any); ok {
				for _, r := range rules {
					if rs, ok := r.(string); ok {
						opts.Rules = append(opts.Rules, rs)
					}
				}
			}

			return s.service.SuggestFixes(ctx, path, opts)
		},
	)
}

func (s *Skill) fixColorsTool() skill.Tool {
	return skill.NewTool(
		"fix_colors",
		"Fix hardcoded colors in a file by replacing them with CSS variables from the design system. Matches exact colors and suggests closest matches for similar colors.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to fix",
				Required:    true,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "If true, return fixes without applying them (default: false)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, _ := params["path"].(string)

			opts := &dss.FixOptions{
				Rules: []string{"no-hardcoded-colors"},
			}

			if dryRun, ok := params["dry_run"].(bool); ok {
				opts.DryRun = dryRun
			}

			return s.service.FixFile(ctx, path, opts)
		},
	)
}

func (s *Skill) fixSpacingTool() skill.Tool {
	return skill.NewTool(
		"fix_spacing",
		"Fix hardcoded spacing values in a file by replacing them with spacing tokens from the design system. Converts pixel values to CSS variables.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to fix",
				Required:    true,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "If true, return fixes without applying them (default: false)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, _ := params["path"].(string)

			opts := &dss.FixOptions{
				Rules: []string{"use-spacing-scale"},
			}

			if dryRun, ok := params["dry_run"].(bool); ok {
				opts.DryRun = dryRun
			}

			return s.service.FixFile(ctx, path, opts)
		},
	)
}

func (s *Skill) fixAccessibilityTool() skill.Tool {
	return skill.NewTool(
		"fix_accessibility",
		"Fix accessibility issues in a file. Adds missing alt attributes to images and aria-label attributes to icon-only buttons.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the file to fix",
				Required:    true,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "If true, return fixes without applying them (default: false)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, _ := params["path"].(string)

			opts := &dss.FixOptions{
				Rules: []string{"img-alt-required", "button-accessible-name"},
			}

			if dryRun, ok := params["dry_run"].(bool); ok {
				opts.DryRun = dryRun
			}

			return s.service.FixFile(ctx, path, opts)
		},
	)
}

func (s *Skill) fixDirectoryTool() skill.Tool {
	return skill.NewTool(
		"fix_directory",
		"Fix design system violations in all files in a directory. Returns a list of all files that were fixed with details of each fix applied.",
		skill.Parameters{
			"path": {
				Type:        "string",
				Description: "Path to the directory to fix",
				Required:    true,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to fix (default: all)",
				Required:    false,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "If true, return fixes without applying them (default: false)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			path, _ := params["path"].(string)
			return s.service.FixDirectory(ctx, path, parseFixOptions(params))
		},
	)
}
