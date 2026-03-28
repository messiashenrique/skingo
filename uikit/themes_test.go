package uikit

import (
	"strings"
	"testing"
)

// TestThemeManagerLightTheme tests light theme creation and tokens
func TestThemeManagerLightTheme(t *testing.T) {
	tm := NewThemeManager()

	if tm.GetCurrentTheme() != "light" {
		t.Errorf("Expected default theme to be 'light', got '%s'", tm.GetCurrentTheme())
	}

	tokens := tm.GetTokens()
	if tokens == nil {
		t.Fatal("Expected tokens to not be nil")
	}

	if tokens.Colors.Primary == "" {
		t.Error("Expected primary color to be set")
	}

	if tokens.Colors.Background != "#ffffff" {
		t.Errorf("Expected background to be white, got %s", tokens.Colors.Background)
	}
}

// TestThemeManagerDarkTheme tests dark theme switching and tokens
func TestThemeManagerDarkTheme(t *testing.T) {
	tm := NewThemeManager()

	err := tm.SetTheme("dark")
	if err != nil {
		t.Fatalf("SetTheme failed: %v", err)
	}

	if tm.GetCurrentTheme() != "dark" {
		t.Errorf("Expected theme to be 'dark', got '%s'", tm.GetCurrentTheme())
	}

	tokens := tm.GetTokens()
	if tokens.Colors.Background != "#0f172a" {
		t.Errorf("Expected dark background, got %s", tokens.Colors.Background)
	}
}

// TestCSSVariablesGeneration tests CSS variables are properly generated
func TestCSSVariablesGeneration(t *testing.T) {
	tm := NewThemeManager()

	cssVars := tm.GetCSSVariablesString()
	if cssVars == "" {
		t.Fatal("Expected CSS variables to not be empty")
	}

	// Check for key CSS variables
	expectedVars := []string{
		"--color-primary",
		"--color-success",
		"--spacing-md",
		"--font-size-lg",
		"--border-radius-md",
		"--shadow-md",
	}

	for _, varName := range expectedVars {
		if !strings.Contains(cssVars, varName) {
			t.Errorf("Expected CSS variable '%s' not found in output", varName)
		}
	}
}

// TestDesignTokensStructure tests that design tokens are properly structured
func TestDesignTokensStructure(t *testing.T) {
	tokens := LightTheme()

	if len(tokens.Spacing.XS) == 0 {
		t.Error("Expected spacing tokens to be set")
	}

	if len(tokens.Border.Radius.MD) == 0 {
		t.Error("Expected border radius tokens to be set")
	}

	if len(tokens.Typography.FontSize.LG) == 0 {
		t.Error("Expected font size tokens to be set")
	}

	if len(tokens.Shadows.MD) == 0 {
		t.Error("Expected shadow tokens to be set")
	}

	if len(tokens.Components.Button.MinHeight) == 0 {
		t.Error("Expected button component tokens to be set")
	}
}

// TestGetTheme tests the GetTheme helper function
func TestGetTheme(t *testing.T) {
	lightTokens := GetTheme("light")
	if lightTokens == nil {
		t.Fatal("Expected light theme tokens to not be nil")
	}

	darkTokens := GetTheme("dark")
	if darkTokens == nil {
		t.Fatal("Expected dark theme tokens to not be nil")
	}

	// Verify they're different
	if lightTokens.Colors.Background == darkTokens.Colors.Background {
		t.Error("Expected light and dark themes to have different backgrounds")
	}

	// Default should be light
	defaultTokens := GetTheme("unknown")
	if defaultTokens.Colors.Background != lightTokens.Colors.Background {
		t.Error("Expected unknown theme to default to light theme")
	}
}
