package designsystem

import (
	"context"
	"fmt"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func (s *Skill) listValidatorsTool() skill.Tool {
	return skill.NewTool(
		"list_validators",
		"List all configured external validators for this design system. Returns accessibility, API style, and custom validators with their configuration.",
		skill.Parameters{},
		func(ctx context.Context, params map[string]any) (any, error) {
			ds := s.service.DesignSystem()
			if ds.Validators == nil {
				return map[string]any{
					"configured": false,
					"message":    "No external validators configured",
					"validators": []any{},
				}, nil
			}

			validators := []map[string]any{}

			// Accessibility validator
			if v := ds.Validators.Accessibility; v != nil {
				validators = append(validators, map[string]any{
					"domain":    "accessibility",
					"tool":      v.Tool,
					"type":      v.Type,
					"command":   v.Command,
					"standards": v.Standards,
					"required":  v.Required,
				})
			}

			// API style validator
			if v := ds.Validators.APIStyle; v != nil {
				validators = append(validators, map[string]any{
					"domain":   "api",
					"tool":     v.Tool,
					"type":     v.Type,
					"command":  v.Command,
					"ruleset":  v.Ruleset,
					"required": v.Required,
				})
			}

			// Custom validators
			for _, v := range ds.Validators.Custom {
				validators = append(validators, map[string]any{
					"domain":   v.Domain,
					"id":       v.ID,
					"name":     v.Name,
					"tool":     v.Tool,
					"type":     v.Type,
					"command":  v.Command,
					"required": v.Required,
				})
			}

			return map[string]any{
				"configured": len(validators) > 0,
				"validators": validators,
				"count":      len(validators),
			}, nil
		},
	)
}

func (s *Skill) getValidatorTool() skill.Tool {
	return skill.NewTool(
		"get_validator",
		"Get configuration for a specific validator by domain (accessibility, api, or custom domain).",
		skill.Parameters{
			"domain": {
				Type:        "string",
				Description: "Validator domain: 'accessibility', 'api', or custom domain name",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			domain := params["domain"].(string)
			ds := s.service.DesignSystem()

			if ds.Validators == nil {
				return nil, fmt.Errorf("no validators configured")
			}

			switch domain {
			case "accessibility":
				if v := ds.Validators.Accessibility; v != nil {
					return map[string]any{
						"domain":    "accessibility",
						"tool":      v.Tool,
						"type":      v.Type,
						"command":   v.Command,
						"standards": v.Standards,
						"checks":    v.Checks,
						"required":  v.Required,
					}, nil
				}
				return nil, fmt.Errorf("accessibility validator not configured")

			case "api":
				if v := ds.Validators.APIStyle; v != nil {
					return map[string]any{
						"domain":   "api",
						"tool":     v.Tool,
						"type":     v.Type,
						"command":  v.Command,
						"ruleset":  v.Ruleset,
						"rules":    v.Rules,
						"required": v.Required,
					}, nil
				}
				return nil, fmt.Errorf("API style validator not configured")

			default:
				// Look in custom validators
				for _, v := range ds.Validators.Custom {
					if v.Domain == domain || v.ID == domain {
						return map[string]any{
							"domain":   v.Domain,
							"id":       v.ID,
							"name":     v.Name,
							"tool":     v.Tool,
							"type":     v.Type,
							"command":  v.Command,
							"config":   v.Config,
							"required": v.Required,
						}, nil
					}
				}
				return nil, fmt.Errorf("validator not found for domain: %s", domain)
			}
		},
	)
}

func (s *Skill) getValidationRequirementsTool() skill.Tool {
	return skill.NewTool(
		"get_validation_requirements",
		"Get all validation requirements defined in the design system spec. Returns accessibility requirements, API style requirements, and any custom validation domains with their corresponding validators.",
		skill.Parameters{},
		func(ctx context.Context, params map[string]any) (any, error) {
			ds := s.service.DesignSystem()

			requirements := map[string]any{}

			// Accessibility requirements
			if a := ds.Accessibility; a != nil {
				req := map[string]any{
					"wcagLevel":   a.WCAGLevel,
					"wcagVersion": a.WCAGVersion,
				}

				if a.ColorContrast != nil {
					req["colorContrast"] = map[string]any{
						"normalTextRatio": a.ColorContrast.NormalTextRatio,
						"largeTextRatio":  a.ColorContrast.LargeTextRatio,
						"nonTextRatio":    a.ColorContrast.NonTextRatio,
					}
				}

				if a.Keyboard != nil {
					req["keyboard"] = map[string]any{
						"focusVisible":   a.Keyboard.FocusVisible,
						"focusOrder":     a.Keyboard.FocusOrder,
						"noKeyboardTrap": a.Keyboard.NoKeyboardTrap,
					}
				}

				// Add validator info if configured
				if ds.Validators != nil && ds.Validators.Accessibility != nil {
					v := ds.Validators.Accessibility
					req["validator"] = map[string]any{
						"tool":    v.Tool,
						"type":    v.Type,
						"command": v.Command,
					}
				}

				requirements["accessibility"] = req
			}

			// API style requirements (from Governance or dedicated field)
			if ds.Validators != nil && ds.Validators.APIStyle != nil {
				v := ds.Validators.APIStyle
				requirements["api"] = map[string]any{
					"validator": map[string]any{
						"tool":    v.Tool,
						"type":    v.Type,
						"command": v.Command,
						"ruleset": v.Ruleset,
					},
				}
			}

			// Custom validation domains
			if ds.Validators != nil {
				for _, v := range ds.Validators.Custom {
					requirements[v.Domain] = map[string]any{
						"id":   v.ID,
						"name": v.Name,
						"validator": map[string]any{
							"tool":    v.Tool,
							"type":    v.Type,
							"command": v.Command,
						},
					}
				}
			}

			return map[string]any{
				"requirements": requirements,
				"domains":      getKeys(requirements),
			}, nil
		},
	)
}

func (s *Skill) getValidatorInvocationTool() skill.Tool {
	return skill.NewTool(
		"get_validator_invocation",
		"Get the invocation details for a validator. Returns how to invoke the validator (MCP server command, CLI command, etc.) based on its type.",
		skill.Parameters{
			"domain": {
				Type:        "string",
				Description: "Validator domain: 'accessibility', 'api', or custom domain name",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			domain := params["domain"].(string)
			ds := s.service.DesignSystem()

			if ds.Validators == nil {
				return nil, fmt.Errorf("no validators configured")
			}

			var tool string
			var vType dss.ValidatorType
			var command string
			var extraConfig map[string]any

			switch domain {
			case "accessibility":
				if v := ds.Validators.Accessibility; v != nil {
					tool = v.Tool
					vType = v.Type
					command = v.Command
					extraConfig = map[string]any{
						"standards": v.Standards,
						"checks":    v.Checks,
					}
				} else {
					return nil, fmt.Errorf("accessibility validator not configured")
				}

			case "api":
				if v := ds.Validators.APIStyle; v != nil {
					tool = v.Tool
					vType = v.Type
					command = v.Command
					extraConfig = map[string]any{
						"ruleset": v.Ruleset,
						"rules":   v.Rules,
					}
				} else {
					return nil, fmt.Errorf("API style validator not configured")
				}

			default:
				found := false
				for _, v := range ds.Validators.Custom {
					if v.Domain == domain || v.ID == domain {
						tool = v.Tool
						vType = v.Type
						command = v.Command
						extraConfig = v.Config
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("validator not found for domain: %s", domain)
				}
			}

			invocation := map[string]any{
				"tool":   tool,
				"type":   vType,
				"config": extraConfig,
			}

			switch vType {
			case dss.ValidatorTypeMCP:
				invocation["method"] = "mcp"
				invocation["instruction"] = fmt.Sprintf("Start MCP server with: %s", command)
				invocation["usage"] = "Connect to the MCP server and use its tools for validation"

			case dss.ValidatorTypeCLI:
				invocation["method"] = "cli"
				invocation["instruction"] = fmt.Sprintf("Run CLI command: %s", command)
				invocation["usage"] = "Execute the command and parse output for validation results"

			case dss.ValidatorTypeNPM:
				invocation["method"] = "npm"
				invocation["instruction"] = fmt.Sprintf("Run via npm/npx: npx %s", tool)
				invocation["usage"] = "Use npx to run the validator package"

			case dss.ValidatorTypeAPI:
				invocation["method"] = "api"
				invocation["instruction"] = fmt.Sprintf("Call API endpoint: %s", command)
				invocation["usage"] = "Make HTTP requests to the validation API"

			case dss.ValidatorTypeLibrary:
				invocation["method"] = "library"
				invocation["instruction"] = fmt.Sprintf("Import Go library: %s", command)
				invocation["usage"] = "Use the library's API for validation"
			}

			return invocation, nil
		},
	)
}

func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
