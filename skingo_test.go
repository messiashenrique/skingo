package skingo

import (
	"testing"
	"testing/fstest"
)

// TestComponentMetaRegistration tests that component metadata can be registered and retrieved.
func TestComponentMetaRegistration(t *testing.T) {
	ts := NewTemplateSet("layout")

	meta := ComponentMeta{
		Name:        "button",
		Description: "Test button component",
		Version:     "1.0.0",
		Variants:    []string{"solid", "outline"},
		Params: []ComponentParam{
			{
				Name:        "label",
				Type:        "string",
				Required:    true,
				Description: "Button label",
			},
		},
	}

	err := ts.RegisterComponentMeta("button", meta)
	if err != nil {
		t.Fatalf("RegisterComponentMeta failed: %v", err)
	}

	retrieved, exists := ts.GetComponentMeta("button")
	if !exists {
		t.Fatal("Component metadata not found after registration")
	}

	if retrieved.Name != "button" {
		t.Errorf("Expected name 'button', got '%s'", retrieved.Name)
	}

	if len(retrieved.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(retrieved.Variants))
	}
}

// TestComponentCatalogRegistration tests JSON catalog registration.
func TestComponentCatalogRegistration(t *testing.T) {
	ts := NewTemplateSet("layout")

	catalogJSON := `{
  "components": {
    "button": {
      "name": "button",
      "description": "Test button",
      "version": "1.0.0",
      "variants": ["solid", "outline"],
      "params": [
        {
          "name": "label",
          "type": "string",
          "required": true,
          "description": "Button label"
        }
      ]
    },
    "input": {
      "name": "input",
      "description": "Test input",
      "version": "1.0.0",
      "params": [
        {
          "name": "placeholder",
          "type": "string",
          "required": false,
          "description": "Placeholder text"
        }
      ]
    }
  }
}`

	err := ts.RegisterComponentCatalogJSON("test-catalog", []byte(catalogJSON))
	if err != nil {
		t.Fatalf("RegisterComponentCatalogJSON failed: %v", err)
	}

	components := ts.ListComponents()
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	names := make(map[string]bool)
	for _, c := range components {
		names[c.Name] = true
	}

	if !names["button"] {
		t.Error("Expected 'button' component not found")
	}
	if !names["input"] {
		t.Error("Expected 'input' component not found")
	}
}

// TestComponentValidation tests parameter validation for components.
func TestComponentValidation(t *testing.T) {
	ts := NewTemplateSet("layout")

	meta := ComponentMeta{
		Name:        "button",
		Description: "Test button",
		Params: []ComponentParam{
			{
				Name:        "label",
				Type:        "string",
				Required:    true,
				Description: "Button label",
			},
		},
		Variants: []string{"solid", "outline", "ghost"},
	}

	ts.RegisterComponentMeta("button", meta)
	ts.SetComponentValidation(ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: true,
	})

	// Test 1: Valid call with required parameter
	validArgs := []interface{}{"Click me"}
	err := ts.validateComponentCall("button", validArgs)
	if err != nil {
		t.Errorf("Expected valid call to pass, got error: %v", err)
	}

	// Test 2: Missing required parameter
	ts.SetComponentValidation(ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: false,
	})
	invalidArgs := []interface{}{}
	err = ts.validateComponentCall("button", invalidArgs)
	if err == nil {
		t.Error("Expected validation error for missing required param, got nil")
	}

	// Test 3: Valid variant with named parameters (using dict)
	validDictArgs := []interface{}{map[string]interface{}{
		"label":   "Click",
		"variant": "outline",
	}}
	err = ts.validateComponentCall("button", validDictArgs)
	if err != nil {
		t.Errorf("Expected valid variant, got error: %v", err)
	}
}

// TestParseManyFS tests parsing from multiple filesystem sources.
func TestParseManyFS(t *testing.T) {
	ts := NewTemplateSet("layout")

	fsMap := fstest.MapFS{
		"templates/layout.html": {
			Data: []byte(`<!DOCTYPE html><html><head><title>Test</title></head><body>{{ .Yield }}</body></html>`),
		},
		"templates/home.html": {
			Data: []byte(`<template><h1>Home</h1></template>`),
		},
		"components/button.html": {
			Data: []byte(`<template><button>{{ param 0 }}</button></template>`),
		},
		"catalog.json": {
			Data: []byte(`{"components": {"button": {"name": "button", "description": "Test button", "version": "1.0.0", "params": [{"name": "label", "type": "string", "required": true}]}}}`),
		},
	}

	catalogData := fsMap["catalog.json"].Data
	err := ts.RegisterComponentCatalogJSON("test-catalog", catalogData)
	if err != nil {
		t.Fatalf("RegisterComponentCatalogJSON failed: %v", err)
	}

	sources := []ParseFSSource{
		{
			Filesystem: fsMap,
			Roots:      []string{"templates"},
		},
		{
			Filesystem: fsMap,
			Roots:      []string{"components"},
		},
	}

	err = ts.ParseManyFS(sources...)
	if err != nil {
		t.Fatalf("ParseManyFS failed: %v", err)
	}

	templates := ts.templates
	if _, hasHome := templates["home"]; !hasHome {
		t.Error("Expected 'home' template not found")
	}

	if _, hasButton := templates["button"]; !hasButton {
		t.Error("Expected 'button' template not found")
	}

	if meta, exists := ts.GetComponentMeta("button"); !exists {
		t.Error("Expected button metadata not found after ParseManyFS")
	} else if meta.Name != "button" {
		t.Errorf("Expected metadata name 'button', got '%s'", meta.Name)
	}
}

// TestComponentValidationOptions tests configuration of validation settings.
func TestComponentValidationOptions(t *testing.T) {
	ts := NewTemplateSet("layout")

	options := ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: true,
	}

	ts.SetComponentValidation(options)
	retrieved := ts.GetComponentValidation()

	if !retrieved.Enabled {
		t.Error("Expected validation to be enabled")
	}

	if !retrieved.StrictTypes {
		t.Error("Expected strict types to be enabled")
	}

	ts.EnableComponentValidation(false)
	retrieved = ts.GetComponentValidation()
	if retrieved.Enabled {
		t.Error("Expected validation to be disabled")
	}
}

// TestListComponents tests component listing functionality.
func TestListComponents(t *testing.T) {
	ts := NewTemplateSet("layout")

	components := ts.ListComponents()
	if len(components) != 0 {
		t.Errorf("Expected 0 components initially, got %d", len(components))
	}

	for i := 1; i <= 3; i++ {
		name := "comp" + string(rune(48+i))
		ts.RegisterComponentMeta(name, ComponentMeta{
			Name:        name,
			Description: "Test component " + string(rune(48+i)),
		})
	}

	components = ts.ListComponents()
	if len(components) != 3 {
		t.Errorf("Expected 3 components, got %d", len(components))
	}
}

// TestCatalogNameNormalization tests catalog name validation.
func TestCatalogNameNormalization(t *testing.T) {
	ts := NewTemplateSet("layout")

	components := make(map[string]ComponentMeta)
	components["button"] = ComponentMeta{Name: "button"}

	err := ts.RegisterComponentCatalog("", components)
	if err == nil {
		t.Error("Expected error for empty catalog name")
	}

	err = ts.RegisterComponentCatalog("test", make(map[string]ComponentMeta))
	if err == nil {
		t.Error("Expected error for empty components")
	}

	err = ts.RegisterComponentCatalog("test-catalog", components)
	if err != nil {
		t.Errorf("Error registering valid catalog: %v", err)
	}
}

// BenchmarkValidation benchmarks component validation performance.
func BenchmarkValidation(b *testing.B) {
	ts := NewTemplateSet("layout")

	meta := ComponentMeta{
		Name: "button",
		Params: []ComponentParam{
			{Name: "label", Type: "string", Required: true},
			{Name: "variant", Type: "string", Required: false},
		},
		Variants: []string{"solid", "outline", "ghost"},
	}

	ts.RegisterComponentMeta("button", meta)
	ts.SetComponentValidation(ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: true,
	})

	args := []interface{}{"Click", "outline"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts.validateComponentCall("button", args)
	}
}
