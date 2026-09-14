package dss

// DesignSystem is the root type representing a complete design system specification.
type DesignSystem struct {
	// Meta contains system-level metadata.
	Meta Meta `json:"meta"`

	// Principles defines design philosophy and brand values.
	Principles []Principle `json:"principles,omitempty"`

	// Foundations contains design tokens (colors, typography, spacing, etc.).
	Foundations Foundations `json:"foundations"`

	// Modes declares the discrete presentation modes this design system
	// defines (e.g. ["light", "dark", "high-contrast"]). Tokens provide
	// per-mode values via their modes maps; the mode-completeness lint rule
	// checks coverage against this list.
	Modes []string `json:"modes,omitempty"`

	// Components defines reusable UI elements.
	Components []Component `json:"components,omitempty"`

	// Patterns defines multi-component UX solutions.
	Patterns []Pattern `json:"patterns,omitempty"`

	// Templates defines page-level layouts.
	Templates []Template `json:"templates,omitempty"`

	// Content defines voice, tone, and content guidelines.
	Content *Content `json:"content,omitempty"`

	// Accessibility defines system-wide accessibility requirements.
	Accessibility *Accessibility `json:"accessibility,omitempty"`

	// Governance defines versioning, contribution, and deprecation policies.
	Governance *Governance `json:"governance,omitempty"`

	// ThemeBindings maps application tokens to component theming contracts.
	ThemeBindings []ThemeBindings `json:"themeBindings,omitempty"`

	// Validators references external validation tools for requirements delegation.
	// DSS defines requirements (e.g., WCAG AA in Accessibility); Validators specifies
	// which external tools perform the actual validation (e.g., agent-a11y, spectral).
	Validators *Validators `json:"validators,omitempty"`
}

// Validate performs basic validation on the design system.
func (ds *DesignSystem) Validate() error {
	if ds.Meta.Name == "" {
		return &ValidationError{Field: "meta.name", Message: "name is required"}
	}
	if ds.Meta.Version == "" {
		return &ValidationError{Field: "meta.version", Message: "version is required"}
	}
	return nil
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
