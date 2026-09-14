package dss

// Foundations contains all design tokens for the system.
type Foundations struct {
	// Colors defines the color palette and semantic color tokens.
	Colors []ColorToken `json:"colors,omitempty"`

	// Typography defines font families, sizes, weights, and type scale.
	Typography *Typography `json:"typography,omitempty"`

	// Spacing defines the spacing scale for margins, padding, and gaps.
	Spacing *SpacingScale `json:"spacing,omitempty"`

	// Elevation defines shadow/depth tokens for layering.
	Elevation []ElevationToken `json:"elevation,omitempty"`

	// Motion defines animation and transition tokens.
	Motion *MotionSystem `json:"motion,omitempty"`

	// Grid defines the layout grid system.
	Grid *GridSystem `json:"grid,omitempty"`

	// Breakpoints defines responsive breakpoints.
	Breakpoints []Breakpoint `json:"breakpoints,omitempty"`

	// BorderRadius defines corner radius tokens.
	BorderRadius []BorderRadiusToken `json:"borderRadius,omitempty"`

	// BorderWidth defines border width tokens.
	BorderWidth []BorderWidthToken `json:"borderWidth,omitempty"`

	// Opacity defines opacity scale tokens.
	Opacity []OpacityToken `json:"opacity,omitempty"`

	// ZIndex defines z-index layer tokens.
	ZIndex []ZIndexToken `json:"zIndex,omitempty"`

	// Densities defines spacing density variants (e.g. comfortable, compact).
	Densities []DensityToken `json:"densities,omitempty"`
}

// DensityToken defines one spacing density variant.
type DensityToken struct {
	// ID is a unique identifier (e.g., "comfortable", "compact").
	ID string `json:"id"`

	// Scale is the spacing multiplier for this density (e.g., 1.0, 0.75).
	Scale float64 `json:"scale"`

	// Description explains when to use this density.
	Description string `json:"description,omitempty"`

	// SpacingOverrides replaces specific spacing-token values (keyed by
	// spacing token ID) for designs where a linear multiplier is not enough.
	SpacingOverrides map[string]string `json:"spacingOverrides,omitempty"`
}

// ColorToken represents a single color in the design system.
type ColorToken struct {
	// ID is a unique identifier (e.g., "primary-500", "danger-base").
	ID string `json:"id"`

	// Value is the color value in hex, rgb, rgba, or hsl format.
	Value string `json:"value"`

	// Semantic indicates the color's purpose (primary, secondary, danger, etc.).
	Semantic string `json:"semantic,omitempty"`

	// Usage describes when and how to use this color (LLM optimization).
	Usage string `json:"usage,omitempty"`

	// Contrast provides accessibility contrast information.
	Contrast *Contrast `json:"contrast,omitempty"`

	// LightModeValue is an alternate value for light mode (if different).
	LightModeValue string `json:"lightModeValue,omitempty"`

	// DarkModeValue is an alternate value for dark mode.
	DarkModeValue string `json:"darkModeValue,omitempty"`

	// Modes generalizes per-mode values beyond light/dark: a map of mode ID
	// (see DesignSystem.Modes) to color value. Explicit entries here win
	// over the LightModeValue/DarkModeValue sugar fields.
	Modes map[string]string `json:"modes,omitempty"`
}

// EffectiveModes folds the LightModeValue/DarkModeValue sugar fields into
// the generalized Modes map. Explicit Modes entries take precedence. The
// returned map is a copy; an empty result means the token has no per-mode
// values.
func (c ColorToken) EffectiveModes() map[string]string {
	modes := map[string]string{}
	if c.LightModeValue != "" {
		modes["light"] = c.LightModeValue
	}
	if c.DarkModeValue != "" {
		modes["dark"] = c.DarkModeValue
	}
	for mode, value := range c.Modes {
		modes[mode] = value
	}
	return modes
}

// Contrast provides WCAG contrast ratio information.
type Contrast struct {
	// OnLight is the recommended text color when this color is on a light background.
	OnLight string `json:"onLight,omitempty"`

	// OnDark is the recommended text color when this color is on a dark background.
	OnDark string `json:"onDark,omitempty"`

	// Ratio is the contrast ratio against the primary background.
	Ratio float64 `json:"ratio,omitempty"`

	// WCAGLevel indicates the WCAG compliance level (AA, AAA).
	WCAGLevel string `json:"wcagLevel,omitempty" jsonschema:"enum=AA,enum=AAA"`
}

// Typography defines the type system.
type Typography struct {
	// FontFamilies lists available font stacks.
	FontFamilies []FontFamily `json:"fontFamilies"`

	// FontSizes defines the font size scale.
	FontSizes []FontSize `json:"fontSizes"`

	// FontWeights defines available font weights.
	FontWeights []FontWeight `json:"fontWeights"`

	// LineHeights defines line height scale.
	LineHeights []LineHeight `json:"lineHeights"`

	// LetterSpacings defines letter spacing scale.
	LetterSpacings []LetterSpacing `json:"letterSpacings,omitempty"`

	// TypeScale defines complete type styles (heading, body, etc.).
	TypeScale []TypeStyle `json:"typeScale,omitempty"`
}

// FontFamily represents a font stack.
type FontFamily struct {
	// ID is a unique identifier (e.g., "sans", "serif", "mono").
	ID string `json:"id"`

	// Name is a display name for the font family.
	Name string `json:"name"`

	// Stack is the CSS font-family fallback stack (e.g., "'Inter', system-ui, sans-serif").
	// This is the preferred field for font family values.
	Stack string `json:"stack,omitempty"`

	// Value is an alias for Stack (deprecated, use Stack instead).
	Value string `json:"value,omitempty"`

	// Weights lists available font weights (e.g., [400, 500, 600, 700]).
	Weights []int `json:"weights,omitempty"`

	// Source is the URL to load the font (e.g., Google Fonts URL).
	Source string `json:"source,omitempty"`

	// Usage describes when to use this font family.
	Usage string `json:"usage,omitempty"`
}

// FontSize represents a font size token.
type FontSize struct {
	// ID is a unique identifier (e.g., "xs", "sm", "base", "lg").
	ID string `json:"id"`

	// Value is the size value with unit (e.g., "1rem", "16px").
	Value string `json:"value"`

	// PixelValue is the equivalent pixel value for reference.
	PixelValue int `json:"pixelValue,omitempty"`
}

// FontWeight represents a font weight token.
type FontWeight struct {
	// ID is a unique identifier (e.g., "normal", "medium", "bold").
	ID string `json:"id"`

	// Value is the numeric weight (100-900).
	Value int `json:"value" jsonschema:"minimum=100,maximum=900"`
}

// LineHeight represents a line height token.
type LineHeight struct {
	// ID is a unique identifier (e.g., "tight", "normal", "relaxed").
	ID string `json:"id"`

	// Value is the line height value (unitless or with unit).
	Value string `json:"value"`
}

// LetterSpacing represents a letter spacing token.
type LetterSpacing struct {
	// ID is a unique identifier.
	ID string `json:"id"`

	// Value is the letter spacing value (e.g., "-0.02em", "0.05em").
	Value string `json:"value"`
}

// TypeStyle defines a complete typographic style combining multiple tokens.
type TypeStyle struct {
	// ID is a unique identifier (e.g., "h1", "body-lg", "caption").
	ID string `json:"id"`

	// Name is a display name for the style.
	Name string `json:"name"`

	// FontFamily references a FontFamily ID.
	FontFamily string `json:"fontFamily"`

	// FontSize references a FontSize ID.
	FontSize string `json:"fontSize"`

	// FontWeight references a FontWeight ID.
	FontWeight string `json:"fontWeight"`

	// LineHeight references a LineHeight ID.
	LineHeight string `json:"lineHeight"`

	// LetterSpacing references a LetterSpacing ID.
	LetterSpacing string `json:"letterSpacing,omitempty"`

	// Usage describes when to use this type style (LLM optimization).
	Usage string `json:"usage,omitempty"`
}

// SpacingScale defines the spacing system.
type SpacingScale struct {
	// BaseUnit is the fundamental spacing unit (e.g., "4px", "0.25rem").
	BaseUnit string `json:"baseUnit"`

	// Scale defines the spacing tokens.
	Scale []SpacingToken `json:"scale"`
}

// SpacingToken represents a single spacing value.
type SpacingToken struct {
	// ID is a unique identifier (e.g., "0", "1", "2", "4", "8").
	ID string `json:"id"`

	// Value is the spacing value with unit.
	Value string `json:"value"`

	// PixelValue is the equivalent pixel value for reference.
	PixelValue int `json:"pixelValue,omitempty"`
}

// ElevationToken represents a shadow/depth level.
type ElevationToken struct {
	// ID is a unique identifier (e.g., "sm", "md", "lg", "xl").
	ID string `json:"id"`

	// Value is the CSS box-shadow value.
	Value string `json:"value"`

	// Usage describes when to use this elevation level.
	Usage string `json:"usage,omitempty"`
}

// MotionSystem defines animation and transition tokens.
type MotionSystem struct {
	// Durations defines timing duration tokens.
	Durations []DurationToken `json:"durations"`

	// Easings defines easing function tokens.
	Easings []EasingToken `json:"easings"`
}

// DurationToken represents an animation duration.
type DurationToken struct {
	// ID is a unique identifier (e.g., "fast", "normal", "slow").
	ID string `json:"id"`

	// Value is the duration in milliseconds or CSS format.
	Value string `json:"value"`

	// Usage describes when to use this duration.
	Usage string `json:"usage,omitempty"`
}

// EasingToken represents an easing function.
type EasingToken struct {
	// ID is a unique identifier (e.g., "easeIn", "easeOut", "spring").
	ID string `json:"id"`

	// Value is the CSS timing function or cubic-bezier value.
	Value string `json:"value"`

	// Usage describes when to use this easing.
	Usage string `json:"usage,omitempty"`
}

// GridSystem defines the layout grid.
type GridSystem struct {
	// Columns is the number of columns (e.g., 12).
	Columns int `json:"columns"`

	// GutterWidth is the space between columns.
	GutterWidth string `json:"gutterWidth"`

	// MarginWidth is the outer margin.
	MarginWidth string `json:"marginWidth"`

	// MaxWidth is the maximum container width.
	MaxWidth string `json:"maxWidth,omitempty"`
}

// Breakpoint defines a responsive breakpoint.
type Breakpoint struct {
	// ID is a unique identifier (e.g., "sm", "md", "lg", "xl").
	ID string `json:"id"`

	// MinWidth is the minimum viewport width for this breakpoint.
	MinWidth string `json:"minWidth"`

	// MaxWidth is the maximum viewport width (optional).
	MaxWidth string `json:"maxWidth,omitempty"`

	// Columns overrides the grid columns at this breakpoint.
	Columns int `json:"columns,omitempty"`
}

// BorderRadiusToken represents a border radius value.
type BorderRadiusToken struct {
	// ID is a unique identifier (e.g., "none", "sm", "md", "lg", "full").
	ID string `json:"id"`

	// Value is the border radius with unit.
	Value string `json:"value"`
}

// BorderWidthToken represents a border width value.
type BorderWidthToken struct {
	// ID is a unique identifier (e.g., "none", "thin", "thick").
	ID string `json:"id"`

	// Value is the border width with unit.
	Value string `json:"value"`
}

// OpacityToken represents an opacity level.
type OpacityToken struct {
	// ID is a unique identifier (e.g., "0", "50", "100").
	ID string `json:"id"`

	// Value is the opacity value (0-1 or percentage).
	Value string `json:"value"`
}

// ZIndexToken represents a z-index layer.
type ZIndexToken struct {
	// ID is a unique identifier (e.g., "base", "dropdown", "modal", "toast").
	ID string `json:"id"`

	// Value is the z-index number.
	Value int `json:"value"`

	// Usage describes when to use this layer.
	Usage string `json:"usage,omitempty"`
}
