package designsystem

import (
	"context"
	"testing"

	dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func TestSkillMetadata(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
	}
	service := dss.NewService(ds)
	skill := New(service)

	if skill.Name() != "designsystem" {
		t.Errorf("expected name 'designsystem', got '%s'", skill.Name())
	}

	if skill.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestSkillInit(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
	}
	service := dss.NewService(ds)
	skill := New(service)

	if err := skill.Init(context.Background()); err != nil {
		t.Errorf("Init failed: %v", err)
	}

	if err := skill.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestSkillTools(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Components: []dss.Component{
			{ID: "button", Name: "Button"},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)

	tools := skill.Tools()

	if len(tools) == 0 {
		t.Fatal("expected tools to be registered")
	}

	// Check for expected tools
	expectedTools := []string{
		"get_component",
		"list_components",
		"get_token",
		"list_tokens",
		"get_pattern",
		"list_patterns",
		"get_meta",
		"generate_prompt",
		"get_variants",
		"get_props",
		"get_anti_patterns",
		"validate_file",
		"validate_directory",
		"check_colors",
		"check_spacing",
		"fix_file",
		"suggest_fixes",
		"fix_colors",
		"fix_spacing",
		"fix_accessibility",
		"fix_directory",
	}

	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name()] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool '%s' not found", expected)
		}
	}
}

func TestGetComponentTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Components: []dss.Component{
			{ID: "button", Name: "Button", Description: "A clickable button"},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the get_component tool
	var getTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "get_component" {
			getTool = tool.Execute
			break
		}
	}

	if getTool == nil {
		t.Fatal("get_component tool not found")
	}

	// Test with valid component
	result, err := getTool(ctx, map[string]any{"id": "button"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	comp, ok := result.(*dss.Component)
	if !ok {
		t.Fatalf("expected *dss.Component, got %T", result)
	}

	if comp.ID != "button" {
		t.Errorf("expected ID 'button', got '%s'", comp.ID)
	}

	// Test with invalid component
	_, err = getTool(ctx, map[string]any{"id": "nonexistent"})
	if err == nil {
		t.Error("expected error for nonexistent component")
	}
}

func TestListComponentsTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Components: []dss.Component{
			{ID: "button", Name: "Button"},
			{ID: "input", Name: "Input"},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the list_components tool
	var listTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "list_components" {
			listTool = tool.Execute
			break
		}
	}

	if listTool == nil {
		t.Fatal("list_components tool not found")
	}

	result, err := listTool(ctx, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	components, ok := result.([]dss.ComponentSummary)
	if !ok {
		t.Fatalf("expected []dss.ComponentSummary, got %T", result)
	}

	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}
}

func TestGetTokenTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Foundations: dss.Foundations{
			Colors: []dss.ColorToken{
				{ID: "primary-500", Value: "#3B82F6"},
			},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the get_token tool
	var getTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "get_token" {
			getTool = tool.Execute
			break
		}
	}

	if getTool == nil {
		t.Fatal("get_token tool not found")
	}

	result, err := getTool(ctx, map[string]any{
		"type": "color",
		"name": "primary-500",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	token, ok := result.(dss.ColorToken)
	if !ok {
		t.Fatalf("expected dss.ColorToken, got %T", result)
	}

	if token.Value != "#3B82F6" {
		t.Errorf("expected value '#3B82F6', got '%s'", token.Value)
	}
}

func TestGetMetaTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{
			Name:        "Test System",
			Version:     "2.0.0",
			Description: "A test design system",
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the get_meta tool
	var getTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "get_meta" {
			getTool = tool.Execute
			break
		}
	}

	if getTool == nil {
		t.Fatal("get_meta tool not found")
	}

	result, err := getTool(ctx, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	meta, ok := result.(dss.Meta)
	if !ok {
		t.Fatalf("expected dss.Meta, got %T", result)
	}

	if meta.Name != "Test System" {
		t.Errorf("expected name 'Test System', got '%s'", meta.Name)
	}
}

func TestGeneratePromptTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Components: []dss.Component{
			{ID: "button", Name: "Button"},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the generate_prompt tool
	var genTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "generate_prompt" {
			genTool = tool.Execute
			break
		}
	}

	if genTool == nil {
		t.Fatal("generate_prompt tool not found")
	}

	result, err := genTool(ctx, map[string]any{
		"format":             "markdown",
		"include_components": true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	prompt, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}

	if len(prompt) == 0 {
		t.Error("expected non-empty prompt")
	}
}

func TestGetVariantsTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Components: []dss.Component{
			{
				ID:   "button",
				Name: "Button",
				Variants: []dss.Variant{
					{ID: "primary", Name: "Primary"},
					{ID: "secondary", Name: "Secondary"},
				},
			},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the get_variants tool
	var getTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "get_variants" {
			getTool = tool.Execute
			break
		}
	}

	if getTool == nil {
		t.Fatal("get_variants tool not found")
	}

	result, err := getTool(ctx, map[string]any{"component_id": "button"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variants, ok := result.([]dss.Variant)
	if !ok {
		t.Fatalf("expected []dss.Variant, got %T", result)
	}

	if len(variants) != 2 {
		t.Errorf("expected 2 variants, got %d", len(variants))
	}
}

func TestGetAntiPatternsTool(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
		Components: []dss.Component{
			{
				ID:   "button",
				Name: "Button",
				LLM: &dss.LLMContext{
					AntiPatterns: []string{
						"Don't use for destructive actions",
						"Avoid multiple primary buttons",
					},
				},
			},
		},
	}
	service := dss.NewService(ds)
	skill := New(service)
	ctx := context.Background()

	// Find the get_anti_patterns tool
	var getTool func(context.Context, map[string]any) (any, error)
	for _, tool := range skill.Tools() {
		if tool.Name() == "get_anti_patterns" {
			getTool = tool.Execute
			break
		}
	}

	if getTool == nil {
		t.Fatal("get_anti_patterns tool not found")
	}

	result, err := getTool(ctx, map[string]any{"component_id": "button"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}

	antiPatterns, ok := resultMap["anti_patterns"].([]string)
	if !ok {
		t.Fatalf("expected []string for anti_patterns, got %T", resultMap["anti_patterns"])
	}

	if len(antiPatterns) != 2 {
		t.Errorf("expected 2 anti-patterns, got %d", len(antiPatterns))
	}
}

func TestServiceAccessor(t *testing.T) {
	ds := &dss.DesignSystem{
		Meta: dss.Meta{Name: "Test", Version: "1.0.0"},
	}
	service := dss.NewService(ds)
	skill := New(service)

	if skill.Service() != service {
		t.Error("Service() should return the underlying service")
	}
}
