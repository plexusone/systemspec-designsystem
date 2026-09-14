package dss

// ThemingContract defines the theming interface for a component.
// It specifies which CSS custom properties can be customized and their semantic meaning.
type ThemingContract struct {
	// Prefix is the CSS custom property prefix for this component (e.g., "--btn", "--card").
	Prefix string `json:"prefix"`

	// Description provides context about the component's theming approach.
	Description string `json:"description,omitempty"`

	// Tokens lists all customizable theme tokens for the component.
	Tokens []ThemeToken `json:"tokens"`
}

// ThemeToken represents a single customizable token in a theming contract.
type ThemeToken struct {
	// ID is a unique identifier for the token within the component (e.g., "background", "text-color").
	ID string `json:"id"`

	// CSSProperty is the full CSS custom property name (e.g., "--btn-background").
	CSSProperty string `json:"cssProperty"`

	// Semantic indicates the token's semantic purpose for auto-mapping.
	// Common values: "primary", "secondary", "accent", "danger", "warning", "success",
	// "neutral", "surface", "text", "text-muted", "border".
	Semantic string `json:"semantic,omitempty"`

	// Description explains what this token controls.
	Description string `json:"description,omitempty"`

	// DefaultLight is the default value for light mode.
	DefaultLight string `json:"defaultLight,omitempty"`

	// DefaultDark is the default value for dark mode.
	DefaultDark string `json:"defaultDark,omitempty"`

	// Defaults generalizes per-mode defaults beyond light/dark: a map of
	// mode ID to default value. Explicit entries win over the
	// DefaultLight/DefaultDark sugar fields.
	Defaults map[string]string `json:"defaults,omitempty"`

	// DensitySensitive marks tokens whose value should scale with the
	// active density (see Foundations.Densities).
	DensitySensitive bool `json:"densitySensitive,omitempty"`
}

// EffectiveDefaults folds the DefaultLight/DefaultDark sugar fields into
// the generalized Defaults map. Explicit Defaults entries take precedence.
func (t ThemeToken) EffectiveDefaults() map[string]string {
	defaults := map[string]string{}
	if t.DefaultLight != "" {
		defaults["light"] = t.DefaultLight
	}
	if t.DefaultDark != "" {
		defaults["dark"] = t.DefaultDark
	}
	for mode, value := range t.Defaults {
		defaults[mode] = value
	}
	return defaults
}

// ThemeBindings maps application design tokens to a component's theming contract.
type ThemeBindings struct {
	// Component is the ID of the component being themed (e.g., "button", "card").
	Component string `json:"component"`

	// SpecURL is an optional URL to fetch the component's spec (for external components).
	// Can be a local file path or HTTP URL.
	SpecURL string `json:"specUrl,omitempty" jsonschema:"format=uri"`

	// ThemeMode specifies which mode these bindings apply to — any mode the
	// design system declares (see DesignSystem.Modes), or empty for all.
	ThemeMode string `json:"themeMode,omitempty"`

	// Strategy defines how to handle unmapped tokens.
	// - "explicit": Only use defined mappings, skip unmapped tokens
	// - "semantic": Auto-map by semantic field, fall back to defaults
	// - "inherit": Use component defaults for all unmapped tokens
	Strategy string `json:"strategy,omitempty" jsonschema:"enum=explicit,enum=semantic,enum=inherit"`

	// Mappings defines explicit token-to-token mappings.
	Mappings []TokenMapping `json:"mappings"`
}

// TokenMapping maps an application design token to a component token.
type TokenMapping struct {
	// From is the application's token (e.g., "brand-primary-500", "colors.primary").
	From string `json:"from"`

	// To is the component's token ID (e.g., "background", "text-color").
	To string `json:"to"`

	// Transform is an optional CSS function to apply (e.g., "rgb", "hsl", "lighten(10%)").
	Transform string `json:"transform,omitempty"`
}

// ValidSemantics lists the allowed semantic values for ThemeToken.Semantic.
var ValidSemantics = []string{
	"primary",
	"secondary",
	"accent",
	"danger",
	"warning",
	"success",
	"info",
	"neutral",
	"surface",
	"background",
	"text",
	"text-muted",
	"text-inverse",
	"border",
	"focus",
	"disabled",
	"shadow",
}

// IsValidSemantic checks if a semantic value is in the allowed set.
func IsValidSemantic(semantic string) bool {
	if semantic == "" {
		return true // Empty is allowed
	}
	for _, s := range ValidSemantics {
		if s == semantic {
			return true
		}
	}
	return false
}
