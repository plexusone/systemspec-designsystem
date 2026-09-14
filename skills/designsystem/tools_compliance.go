package designsystem

import (
	"context"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func (s *Skill) generateComplianceReportTool() skill.Tool {
	return skill.NewTool(
		"generate_compliance_report",
		"Generate a comprehensive compliance report for release. Returns score, status, category breakdown, and issues list. Use this before release to verify design system compliance.",
		skill.Parameters{
			"directory": {
				Type:        "string",
				Description: "Directory to validate",
				Required:    true,
			},
			"min_score": {
				Type:        "number",
				Description: "Minimum acceptable score (default: 80)",
				Required:    false,
			},
			"blocking_categories": {
				Type:        "array",
				Description: "Categories that block release on failure (default: ['colors', 'accessibility'])",
				Required:    false,
			},
			"include_issues": {
				Type:        "boolean",
				Description: "Include detailed issue list (default: true)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			dir := params["directory"].(string)

			opts := dss.DefaultComplianceOptions()

			if minScore, ok := params["min_score"].(float64); ok {
				opts.MinScore = int(minScore)
			}

			if blocking, ok := params["blocking_categories"].([]any); ok {
				opts.BlockingCategories = make([]string, 0, len(blocking))
				for _, b := range blocking {
					if bs, ok := b.(string); ok {
						opts.BlockingCategories = append(opts.BlockingCategories, bs)
					}
				}
			}

			if include, ok := params["include_issues"].(bool); ok {
				opts.IncludeIssues = include
			}

			report, err := s.service.GenerateComplianceReport(ctx, dir, opts)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"status":       report.Status,
				"score":        report.Score,
				"timestamp":    report.Timestamp,
				"designSystem": report.DesignSystem,
				"target":       report.Target,
				"categories":   report.Categories,
				"summary":      report.Summary,
				"issues":       report.Issues,
			}, nil
		},
	)
}

func (s *Skill) checkReleaseGateTool() skill.Tool {
	return skill.NewTool(
		"check_release_gate",
		"Check if the codebase passes the release gate. Returns approved/rejected with reasons. Use this as a final check before creating a release.",
		skill.Parameters{
			"directory": {
				Type:        "string",
				Description: "Directory to validate",
				Required:    true,
			},
			"min_score": {
				Type:        "number",
				Description: "Minimum acceptable score (default: 80)",
				Required:    false,
			},
			"require_zero_errors": {
				Type:        "boolean",
				Description: "Require zero error-level issues (default: true)",
				Required:    false,
			},
			"allow_warnings": {
				Type:        "boolean",
				Description: "Allow release with warnings (default: true)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			dir := params["directory"].(string)

			opts := dss.DefaultReleaseGateOptions()

			if minScore, ok := params["min_score"].(float64); ok {
				opts.MinScore = int(minScore)
			}

			if requireZero, ok := params["require_zero_errors"].(bool); ok {
				opts.RequireZeroErrors = requireZero
			}

			if allowWarn, ok := params["allow_warnings"].(bool); ok {
				opts.AllowWarnings = allowWarn
			}

			gate, err := s.service.CheckReleaseGate(ctx, dir, opts)
			if err != nil {
				return nil, err
			}

			result := map[string]any{
				"approved":       gate.Approved,
				"reason":         gate.Reason,
				"score":          gate.Score,
				"blockingIssues": gate.BlockingIssues,
				"warnings":       gate.Warnings,
			}

			if gate.Certificate != nil {
				result["certificate"] = map[string]any{
					"id":        gate.Certificate.ID,
					"timestamp": gate.Certificate.Timestamp,
					"hash":      gate.Certificate.Hash,
					"status":    gate.Certificate.Status,
				}
			}

			return result, nil
		},
	)
}

func (s *Skill) runFixLoopTool() skill.Tool {
	return skill.NewTool(
		"run_fix_loop",
		"Run an automated fix-validate loop until convergence or max iterations. Fixes violations, re-validates, and repeats until no more fixes can be applied. Use this to automatically fix all auto-fixable issues.",
		skill.Parameters{
			"directory": {
				Type:        "string",
				Description: "Directory to fix",
				Required:    true,
			},
			"max_iterations": {
				Type:        "number",
				Description: "Maximum fix-validate cycles (default: 3)",
				Required:    false,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "Preview fixes without applying (default: false)",
				Required:    false,
			},
			"stop_on_regression": {
				Type:        "boolean",
				Description: "Stop if fixes introduce new violations (default: true)",
				Required:    false,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to fix (default: all fixable rules)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			dir := params["directory"].(string)

			opts := dss.DefaultFixLoopOptions()

			if maxIter, ok := params["max_iterations"].(float64); ok {
				opts.MaxIterations = int(maxIter)
			}

			if dryRun, ok := params["dry_run"].(bool); ok {
				opts.DryRun = dryRun
			}

			if stopReg, ok := params["stop_on_regression"].(bool); ok {
				opts.StopOnRegression = stopReg
			}

			if rules, ok := params["rules"].([]any); ok {
				for _, r := range rules {
					if rs, ok := r.(string); ok {
						opts.Rules = append(opts.Rules, rs)
					}
				}
			}

			result, err := s.service.RunFixLoop(ctx, dir, opts)
			if err != nil {
				// Return partial result with error info
				return map[string]any{
					"error":             err.Error(),
					"converged":         false,
					"iterations":        result.Iterations,
					"initialViolations": result.InitialViolations,
					"finalViolations":   result.FinalViolations,
					"fixedCount":        result.FixedCount,
					"iterationDetails":  result.IterationDetails,
					"remainingIssues":   result.RemainingIssues,
				}, nil
			}

			return map[string]any{
				"converged":         result.Converged,
				"iterations":        result.Iterations,
				"initialViolations": result.InitialViolations,
				"finalViolations":   result.FinalViolations,
				"fixedCount":        result.FixedCount,
				"unfixableCount":    result.UnfixableCount,
				"iterationDetails":  result.IterationDetails,
				"remainingIssues":   result.RemainingIssues,
				"status":            result.Status,
			}, nil
		},
	)
}

func (s *Skill) fixAndVerifyTool() skill.Tool {
	return skill.NewTool(
		"fix_and_verify",
		"Fix a single file and verify the fix resolved the violations. Returns before/after comparison and verification status. Use this for targeted fixes with confirmation.",
		skill.Parameters{
			"file": {
				Type:        "string",
				Description: "File to fix",
				Required:    true,
			},
			"rules": {
				Type:        "array",
				Description: "Specific rules to fix (default: all fixable rules)",
				Required:    false,
			},
			"dry_run": {
				Type:        "boolean",
				Description: "Preview fixes without applying (default: false)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			file := params["file"].(string)

			opts := dss.DefaultFixOptions()

			if dryRun, ok := params["dry_run"].(bool); ok {
				opts.DryRun = dryRun
			}

			if rules, ok := params["rules"].([]any); ok {
				for _, r := range rules {
					if rs, ok := r.(string); ok {
						opts.Rules = append(opts.Rules, rs)
					}
				}
			}

			result, err := s.service.FixAndVerify(ctx, file, opts)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"file":             result.File,
				"violationsBefore": result.ViolationsBefore,
				"fixesApplied":     result.FixesApplied,
				"violationsAfter":  result.ViolationsAfter,
				"improvement":      result.Improvement,
				"regressions":      result.Regressions,
				"verified":         result.Verified,
				"reason":           result.Reason,
			}, nil
		},
	)
}

func (s *Skill) getComplianceCertificateTool() skill.Tool {
	return skill.NewTool(
		"get_compliance_certificate",
		"Generate a compliance certificate after passing the release gate. The certificate provides proof of compliance that can be included in release artifacts.",
		skill.Parameters{
			"directory": {
				Type:        "string",
				Description: "Directory that was validated",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			dir := params["directory"].(string)

			// Generate compliance report
			report, err := s.service.GenerateComplianceReport(ctx, dir, nil)
			if err != nil {
				return nil, err
			}

			// Generate certificate
			cert := s.service.GenerateCertificate(report)

			return map[string]any{
				"id":           cert.ID,
				"timestamp":    cert.Timestamp,
				"designSystem": cert.DesignSystem,
				"target":       cert.Target,
				"status":       cert.Status,
				"score":        cert.Score,
				"hash":         cert.Hash,
				"categories":   cert.Categories,
			}, nil
		},
	)
}
