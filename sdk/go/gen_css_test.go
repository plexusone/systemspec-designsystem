package dss

import (
	"strings"
	"testing"
)

func cssFixture() *DesignSystem {
	return &DesignSystem{
		Meta:  Meta{Name: "Test", Version: "1.0.0"},
		Modes: []string{"light", "dark"},
		Foundations: Foundations{
			Colors: []ColorToken{
				{ID: "primary", Value: "#0af", Modes: map[string]string{"light": "#0df", "dark": "#068"}},
				{ID: "surface", Value: "#111", LightModeValue: "#fff"},
				{ID: "static", Value: "#123"},
			},
			Spacing: &SpacingScale{
				BaseUnit: "0.25rem",
				Scale:    []SpacingToken{{ID: "4", Value: "1rem"}},
			},
			Densities: []DensityToken{
				{ID: "comfortable", Scale: 1.0},
				{ID: "compact", Scale: 0.75, SpacingOverrides: map[string]string{"4": "0.65rem"}},
			},
		},
	}
}

func TestGenerateCSSVarsModesAndDensity(t *testing.T) {
	ds := cssFixture()
	out, err := ds.GenerateCSS(CSSGeneratorOptions{Format: "css-vars"})
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"--color-primary: #0af;",
		"[data-mode=\"dark\"] {\n  --color-primary: #068;\n}",
		"[data-mode=\"light\"] {\n  --color-primary: #0df;\n  --color-surface: #fff;\n}",
		":root {\n  --density: 1;\n}",
		"[data-density=\"compact\"] {\n  --density: 0.75;\n  --spacing-4: 0.65rem;\n}",
		"[data-density=\"comfortable\"] {\n  --density: 1;\n}",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("css-vars output missing %q\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, "--color-static: #123;\n}") && strings.Count(out, "--color-static") != 1 {
		t.Errorf("mode-less token should appear only in :root")
	}
}

func TestGenerateCSSVarsPrefix(t *testing.T) {
	ds := cssFixture()
	out, err := ds.GenerateCSS(CSSGeneratorOptions{Format: "css-vars", Prefix: "pm"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"--pm-color-primary: #0af;",
		"--pm-density: 1;",
		"--pm-spacing-4: 0.65rem;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("prefixed output missing %q", want)
		}
	}
}

func TestGenerateTailwindModesAndDensity(t *testing.T) {
	ds := cssFixture()
	out, err := ds.GenerateCSS(CSSGeneratorOptions{Format: "tailwind4"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"[data-mode=\"dark\"]",
		"--color-primary: #068;",
		"[data-density=\"compact\"]",
		"--spacing-4: 0.65rem;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("tailwind output missing %q", want)
		}
	}
}

func TestGenerateSCSSModesAndDensity(t *testing.T) {
	ds := cssFixture()
	out, err := ds.GenerateCSS(CSSGeneratorOptions{Format: "scss"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"$color-primary--dark: #068;",
		"$color-surface--light: #fff;",
		"$density-compact: 0.75;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("scss output missing %q\n---\n%s", want, out)
		}
	}
}

// TestGenerateCSSBackwardCompatible guards that documents without modes or
// densities emit no mode/density constructs at all.
func TestGenerateCSSBackwardCompatible(t *testing.T) {
	ds := &DesignSystem{
		Meta: Meta{Name: "Plain", Version: "1.0.0"},
		Foundations: Foundations{
			Colors:  []ColorToken{{ID: "primary", Value: "#0af"}},
			Spacing: &SpacingScale{BaseUnit: "0.25rem", Scale: []SpacingToken{{ID: "4", Value: "1rem"}}},
		},
	}
	for _, format := range []string{"tailwind4", "css-vars", "scss"} {
		out, err := ds.GenerateCSS(CSSGeneratorOptions{Format: format})
		if err != nil {
			t.Fatal(err)
		}
		for _, banned := range []string{"data-mode", "data-density", "--density", "$density", "// Modes"} {
			if strings.Contains(out, banned) {
				t.Errorf("%s output for mode-free document contains %q", format, banned)
			}
		}
	}
}
