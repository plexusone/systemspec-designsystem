package dss

import (
	"encoding/json"
	"fmt"
)

// W3CTokensOptions configures W3C Design Tokens output.
type W3CTokensOptions struct {
	// Prefix for token names (e.g., "mde" → "mde.color.primary")
	Prefix string

	// IncludeDescriptions adds $description to tokens
	IncludeDescriptions bool

	// IncludeExtensions adds $extensions with DSS metadata
	IncludeExtensions bool
}

// DefaultW3CTokensOptions returns sensible defaults.
func DefaultW3CTokensOptions() W3CTokensOptions {
	return W3CTokensOptions{
		Prefix:              "",
		IncludeDescriptions: true,
		IncludeExtensions:   true,
	}
}

// W3CTokenFile represents the root of a W3C Design Tokens file.
type W3CTokenFile struct {
	Schema  string                 `json:"$schema,omitempty"`
	Color   map[string]*W3CToken   `json:"color,omitempty"`
	Space   map[string]*W3CToken   `json:"space,omitempty"`
	Size    map[string]*W3CToken   `json:"size,omitempty"`
	Radius  map[string]*W3CToken   `json:"radius,omitempty"`
	Shadow  map[string]*W3CToken   `json:"shadow,omitempty"`
	Font    map[string]interface{} `json:"font,omitempty"`
	Motion  map[string]interface{} `json:"motion,omitempty"`
	ZIndex  map[string]*W3CToken   `json:"zIndex,omitempty"`
	Density map[string]*W3CToken   `json:"density,omitempty"`
}

// W3CToken represents a single W3C Design Token.
type W3CToken struct {
	Value       interface{}            `json:"$value"`
	Type        string                 `json:"$type,omitempty"`
	Description string                 `json:"$description,omitempty"`
	Extensions  map[string]interface{} `json:"$extensions,omitempty"`
}

// W3CFontFamily represents a W3C font family token group.
type W3CFontFamily struct {
	Value       interface{}            `json:"$value"`
	Type        string                 `json:"$type,omitempty"`
	Description string                 `json:"$description,omitempty"`
	Extensions  map[string]interface{} `json:"$extensions,omitempty"`
}

// GenerateW3CTokens generates W3C Design Tokens Community Group format.
// See: https://design-tokens.github.io/community-group/format/
func (ds *DesignSystem) GenerateW3CTokens(opts W3CTokensOptions) (string, error) {
	if ds.Foundations.isEmpty() {
		return "", fmt.Errorf("design system has no foundations defined")
	}

	tokens := &W3CTokenFile{
		Schema: "https://design-tokens.github.io/community-group/format/",
	}

	f := ds.Foundations

	// Colors
	if len(f.Colors) > 0 {
		tokens.Color = make(map[string]*W3CToken)
		for _, c := range f.Colors {
			token := &W3CToken{
				Value: c.Value,
				Type:  "color",
			}
			if opts.IncludeDescriptions && c.Usage != "" {
				token.Description = c.Usage
			}
			if opts.IncludeExtensions {
				token.Extensions = map[string]interface{}{
					"com.plexusone.dss": map[string]interface{}{
						"id":       c.ID,
						"category": "color",
					},
				}
				if modes := c.EffectiveModes(); len(modes) > 0 {
					token.Extensions["com.plexusone.dss"].(map[string]interface{})["modes"] = modes
				}
			}
			tokens.Color[normalizeTokenName(c.ID)] = token
		}
	}

	if len(f.Densities) > 0 {
		tokens.Density = make(map[string]*W3CToken)
		for _, d := range f.Densities {
			token := &W3CToken{Value: d.Scale, Type: "number"}
			if opts.IncludeDescriptions && d.Description != "" {
				token.Description = d.Description
			}
			if opts.IncludeExtensions {
				ext := map[string]interface{}{"id": d.ID, "category": "density"}
				if len(d.SpacingOverrides) > 0 {
					ext["spacingOverrides"] = d.SpacingOverrides
				}
				token.Extensions = map[string]interface{}{"com.plexusone.dss": ext}
			}
			tokens.Density[normalizeTokenName(d.ID)] = token
		}
	}

	// Spacing
	if f.Spacing != nil && len(f.Spacing.Scale) > 0 {
		tokens.Space = make(map[string]*W3CToken)
		for _, s := range f.Spacing.Scale {
			token := &W3CToken{
				Value: s.Value,
				Type:  "dimension",
			}
			if opts.IncludeExtensions {
				token.Extensions = map[string]interface{}{
					"com.plexusone.dss": map[string]interface{}{
						"id":       s.ID,
						"category": "spacing",
					},
				}
			}
			tokens.Space[normalizeTokenName(s.ID)] = token
		}
	}

	// Border Radius
	if len(f.BorderRadius) > 0 {
		tokens.Radius = make(map[string]*W3CToken)
		for _, br := range f.BorderRadius {
			token := &W3CToken{
				Value: br.Value,
				Type:  "dimension",
			}
			if opts.IncludeExtensions {
				token.Extensions = map[string]interface{}{
					"com.plexusone.dss": map[string]interface{}{
						"id":       br.ID,
						"category": "radius",
					},
				}
			}
			tokens.Radius[normalizeTokenName(br.ID)] = token
		}
	}

	// Elevation/Shadows
	if len(f.Elevation) > 0 {
		tokens.Shadow = make(map[string]*W3CToken)
		for _, e := range f.Elevation {
			token := &W3CToken{
				Value: e.Value,
				Type:  "shadow",
			}
			if opts.IncludeDescriptions && e.Usage != "" {
				token.Description = e.Usage
			}
			if opts.IncludeExtensions {
				token.Extensions = map[string]interface{}{
					"com.plexusone.dss": map[string]interface{}{
						"id":       e.ID,
						"category": "elevation",
					},
				}
			}
			tokens.Shadow[normalizeTokenName(e.ID)] = token
		}
	}

	// Typography
	if f.Typography != nil {
		tokens.Font = make(map[string]interface{})

		// Font Families
		if len(f.Typography.FontFamilies) > 0 {
			families := make(map[string]*W3CToken)
			for _, ff := range f.Typography.FontFamilies {
				value := ff.Stack
				if value == "" {
					value = ff.Value
				}
				token := &W3CToken{
					Value: value,
					Type:  "fontFamily",
				}
				if opts.IncludeDescriptions && ff.Usage != "" {
					token.Description = ff.Usage
				}
				families[normalizeTokenName(ff.ID)] = token
			}
			tokens.Font["family"] = families
		}

		// Font Sizes
		if len(f.Typography.FontSizes) > 0 {
			sizes := make(map[string]*W3CToken)
			for _, fs := range f.Typography.FontSizes {
				token := &W3CToken{
					Value: fs.Value,
					Type:  "dimension",
				}
				sizes[normalizeTokenName(fs.ID)] = token
			}
			tokens.Font["size"] = sizes
		}

		// Font Weights
		if len(f.Typography.FontWeights) > 0 {
			weights := make(map[string]*W3CToken)
			for _, fw := range f.Typography.FontWeights {
				token := &W3CToken{
					Value: fw.Value,
					Type:  "fontWeight",
				}
				weights[normalizeTokenName(fw.ID)] = token
			}
			tokens.Font["weight"] = weights
		}

		// Line Heights
		if len(f.Typography.LineHeights) > 0 {
			lineHeights := make(map[string]*W3CToken)
			for _, lh := range f.Typography.LineHeights {
				token := &W3CToken{
					Value: lh.Value,
					Type:  "number",
				}
				lineHeights[normalizeTokenName(lh.ID)] = token
			}
			tokens.Font["lineHeight"] = lineHeights
		}
	}

	// Motion
	if f.Motion != nil {
		tokens.Motion = make(map[string]interface{})

		// Durations
		if len(f.Motion.Durations) > 0 {
			durations := make(map[string]*W3CToken)
			for _, d := range f.Motion.Durations {
				token := &W3CToken{
					Value: d.Value,
					Type:  "duration",
				}
				if opts.IncludeDescriptions && d.Usage != "" {
					token.Description = d.Usage
				}
				durations[normalizeTokenName(d.ID)] = token
			}
			tokens.Motion["duration"] = durations
		}

		// Easings
		if len(f.Motion.Easings) > 0 {
			easings := make(map[string]*W3CToken)
			for _, e := range f.Motion.Easings {
				token := &W3CToken{
					Value: e.Value,
					Type:  "cubicBezier",
				}
				if opts.IncludeDescriptions && e.Usage != "" {
					token.Description = e.Usage
				}
				easings[normalizeTokenName(e.ID)] = token
			}
			tokens.Motion["easing"] = easings
		}
	}

	// Z-Index
	if len(f.ZIndex) > 0 {
		tokens.ZIndex = make(map[string]*W3CToken)
		for _, z := range f.ZIndex {
			token := &W3CToken{
				Value: z.Value,
				Type:  "number",
			}
			if opts.IncludeDescriptions && z.Usage != "" {
				token.Description = z.Usage
			}
			tokens.ZIndex[normalizeTokenName(z.ID)] = token
		}
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling W3C tokens: %w", err)
	}

	return string(data), nil
}

// normalizeTokenName converts IDs to W3C-friendly camelCase names.
func normalizeTokenName(id string) string {
	// For now, just return as-is (IDs should already be valid)
	// Could add camelCase conversion if needed
	return id
}
