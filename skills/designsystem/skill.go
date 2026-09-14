// Package designsystem provides an omniskill-based skill for design system operations.
// It exposes tools for reading spec, generating guidance, and validating implementations.
package designsystem

import (
	"context"

	"github.com/plexusone/systemspec-designsystem/internal/omniskill/skill"
	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

// Skill provides design system operations as MCP tools.
type Skill struct {
	service *dss.Service
}

// New creates a new design system skill from a service.
func New(service *dss.Service) *Skill {
	return &Skill{
		service: service,
	}
}

// Name returns the skill identifier.
func (s *Skill) Name() string {
	return "designsystem"
}

// Description returns a human-readable description.
func (s *Skill) Description() string {
	return "Design system specification tools for reading specs, generating guidance, and validating implementations"
}

// Init performs any initialization needed by the skill.
func (s *Skill) Init(_ context.Context) error {
	return nil
}

// Close cleans up any resources.
func (s *Skill) Close() error {
	return nil
}

// Tools returns all tools provided by this skill.
func (s *Skill) Tools() []skill.Tool {
	return []skill.Tool{
		// Spec reading tools
		s.getComponentTool(),
		s.listComponentsTool(),
		s.getTokenTool(),
		s.listTokensTool(),
		s.getPatternTool(),
		s.listPatternsTool(),
		s.getMetaTool(),

		// Implementation guidance tools
		s.generatePromptTool(),
		s.getVariantsTool(),
		s.getPropsTool(),
		s.getAntiPatternsTool(),

		// Validation tools
		s.validateFileTool(),
		s.validateDirectoryTool(),
		s.checkColorsTool(),
		s.checkSpacingTool(),

		// Fix tools
		s.fixFileTool(),
		s.suggestFixesTool(),
		s.fixColorsTool(),
		s.fixSpacingTool(),
		s.fixAccessibilityTool(),
		s.fixDirectoryTool(),

		// Lint tools
		s.lintSpecTool(),
		s.listLintRulesTool(),
		s.checkAgentReadinessTool(),

		// Validator delegation tools
		s.listValidatorsTool(),
		s.getValidatorTool(),
		s.getValidationRequirementsTool(),
		s.getValidatorInvocationTool(),

		// Compliance and release tools
		s.generateComplianceReportTool(),
		s.checkReleaseGateTool(),
		s.runFixLoopTool(),
		s.fixAndVerifyTool(),
		s.getComplianceCertificateTool(),

		// Accessibility integration tools (Phase 1-4)
		s.getAccessibilityRequirementsTool(),
		s.getA11yAntiPatternsTool(),
		s.suggestContrastTokenTool(),
		s.getComponentFixContextTool(),
	}
}

// Service returns the underlying service for direct SDK access.
func (s *Skill) Service() *dss.Service {
	return s.service
}
