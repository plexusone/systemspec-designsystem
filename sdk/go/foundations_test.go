package dss

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEffectiveModes(t *testing.T) {
	tests := []struct {
		name  string
		token ColorToken
		want  map[string]string
	}{
		{
			name:  "no mode values",
			token: ColorToken{ID: "plain", Value: "#000"},
			want:  map[string]string{},
		},
		{
			name:  "sugar fields fold in",
			token: ColorToken{ID: "primary", Value: "#0af", LightModeValue: "#0df", DarkModeValue: "#068"},
			want:  map[string]string{"light": "#0df", "dark": "#068"},
		},
		{
			name: "explicit modes win over sugar and extend beyond light/dark",
			token: ColorToken{
				ID: "primary", Value: "#0af",
				LightModeValue: "#0df",
				Modes:          map[string]string{"light": "#override", "high-contrast": "#fff"},
			},
			want: map[string]string{"light": "#override", "high-contrast": "#fff"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.EffectiveModes()
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("mode %q = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestEffectiveDefaults(t *testing.T) {
	token := ThemeToken{
		ID: "background", CSSProperty: "--x-bg",
		DefaultLight: "#fff", DefaultDark: "#111",
		Defaults: map[string]string{"dark": "#000", "high-contrast": "#0f0"},
	}
	got := token.EffectiveDefaults()
	want := map[string]string{"light": "#fff", "dark": "#000", "high-contrast": "#0f0"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("default %q = %q, want %q", k, got[k], v)
		}
	}
}

func TestLoadModesAndDensitiesFromDirectory(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, content string) {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("meta.json", `{"name": "Test", "version": "1.0.0"}`)
	write("modes.json", `["light", "dark", "high-contrast"]`)
	write("foundations/colors.json", `[
		{"id": "primary", "value": "#0af", "modes": {"light": "#0df", "dark": "#068"}}
	]`)
	write("foundations/densities.json", `[
		{"id": "comfortable", "scale": 1.0},
		{"id": "compact", "scale": 0.75, "spacingOverrides": {"4": "0.65rem"}}
	]`)

	ds, err := LoadDesignSystem(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds.Modes) != 3 || ds.Modes[2] != "high-contrast" {
		t.Errorf("Modes = %v", ds.Modes)
	}
	if got := ds.Foundations.Colors[0].EffectiveModes()["dark"]; got != "#068" {
		t.Errorf("dark mode value = %q", got)
	}
	if len(ds.Foundations.Densities) != 2 {
		t.Fatalf("Densities = %v", ds.Foundations.Densities)
	}
	compact := ds.Foundations.Densities[1]
	if compact.Scale != 0.75 || compact.SpacingOverrides["4"] != "0.65rem" {
		t.Errorf("compact density = %+v", compact)
	}
}

func TestLoadModesAndDensitiesFromSingleFile(t *testing.T) {
	dir := t.TempDir()
	doc := `{
		"meta": {"name": "Test", "version": "1.0.0"},
		"modes": ["light", "dark"],
		"foundations": {
			"colors": [{"id": "primary", "value": "#0af", "darkModeValue": "#068"}],
			"densities": [{"id": "compact", "scale": 0.8}]
		}
	}`
	path := filepath.Join(dir, "design-system.json")
	if err := os.WriteFile(path, []byte(doc), 0o600); err != nil {
		t.Fatal(err)
	}
	ds, err := LoadDesignSystem(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds.Modes) != 2 || ds.Foundations.Densities[0].Scale != 0.8 {
		t.Errorf("single-file load: modes=%v densities=%+v", ds.Modes, ds.Foundations.Densities)
	}
}
