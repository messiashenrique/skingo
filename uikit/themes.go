package uikit

import (
	"fmt"
	"strings"
	"text/template"
)

// ThemeManager manages theme switching and CSS variable injection
type ThemeManager struct {
	currentTheme string
	tokens       *DesignTokens
	cssVariables string
}

// NewThemeManager creates a new theme manager with light theme by default
func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{
		currentTheme: "light",
		tokens:       LightTheme(),
	}
	tm.generateCSSVariables()
	return tm
}

// SetTheme changes the current theme and regenerates CSS variables
func (tm *ThemeManager) SetTheme(name string) error {
	tokens := GetTheme(name)
	if tokens == nil {
		return fmt.Errorf("unknown theme: %s", name)
	}
	tm.currentTheme = name
	tm.tokens = tokens
	tm.generateCSSVariables()
	return nil
}

// GetCurrentTheme returns the name of the current theme
func (tm *ThemeManager) GetCurrentTheme() string {
	return tm.currentTheme
}

// GetTokens returns the current design tokens
func (tm *ThemeManager) GetTokens() *DesignTokens {
	return tm.tokens
}

// generateCSSVariables creates CSS custom properties from design tokens
func (tm *ThemeManager) generateCSSVariables() {
	var sb strings.Builder
	sb.WriteString(":root {\n")

	// Color variables
	sb.WriteString(fmt.Sprintf("  --color-primary: %s;\n", tm.tokens.Colors.Primary))
	sb.WriteString(fmt.Sprintf("  --color-secondary: %s;\n", tm.tokens.Colors.Secondary))
	sb.WriteString(fmt.Sprintf("  --color-success: %s;\n", tm.tokens.Colors.Success))
	sb.WriteString(fmt.Sprintf("  --color-warning: %s;\n", tm.tokens.Colors.Warning))
	sb.WriteString(fmt.Sprintf("  --color-error: %s;\n", tm.tokens.Colors.Error))
	sb.WriteString(fmt.Sprintf("  --color-info: %s;\n", tm.tokens.Colors.Info))
	sb.WriteString(fmt.Sprintf("  --color-background: %s;\n", tm.tokens.Colors.Background))
	sb.WriteString(fmt.Sprintf("  --color-surface: %s;\n", tm.tokens.Colors.Surface))
	sb.WriteString(fmt.Sprintf("  --color-border: %s;\n", tm.tokens.Colors.Border))
	sb.WriteString(fmt.Sprintf("  --color-text: %s;\n", tm.tokens.Colors.Text))
	sb.WriteString(fmt.Sprintf("  --color-text-muted: %s;\n", tm.tokens.Colors.TextMuted))
	sb.WriteString(fmt.Sprintf("  --color-primary-light: %s;\n", tm.tokens.Colors.PrimaryLight))
	sb.WriteString(fmt.Sprintf("  --color-secondary-light: %s;\n", tm.tokens.Colors.SecondaryLight))
	sb.WriteString(fmt.Sprintf("  --color-success-light: %s;\n", tm.tokens.Colors.SuccessLight))
	sb.WriteString(fmt.Sprintf("  --color-warning-light: %s;\n", tm.tokens.Colors.WarningLight))
	sb.WriteString(fmt.Sprintf("  --color-error-light: %s;\n", tm.tokens.Colors.ErrorLight))
	sb.WriteString(fmt.Sprintf("  --color-info-light: %s;\n", tm.tokens.Colors.InfoLight))
	sb.WriteString(fmt.Sprintf("  --color-primary-outline: %s;\n", tm.tokens.Colors.PrimaryOutline))
	sb.WriteString(fmt.Sprintf("  --color-secondary-outline: %s;\n", tm.tokens.Colors.SecondaryOutline))
	sb.WriteString(fmt.Sprintf("  --color-success-outline: %s;\n", tm.tokens.Colors.SuccessOutline))
	sb.WriteString(fmt.Sprintf("  --color-warning-outline: %s;\n", tm.tokens.Colors.WarningOutline))
	sb.WriteString(fmt.Sprintf("  --color-error-outline: %s;\n", tm.tokens.Colors.ErrorOutline))
	sb.WriteString(fmt.Sprintf("  --color-info-outline: %s;\n", tm.tokens.Colors.InfoOutline))

	// Spacing variables
	sb.WriteString(fmt.Sprintf("  --spacing-xs: %s;\n", tm.tokens.Spacing.XS))
	sb.WriteString(fmt.Sprintf("  --spacing-sm: %s;\n", tm.tokens.Spacing.SM))
	sb.WriteString(fmt.Sprintf("  --spacing-md: %s;\n", tm.tokens.Spacing.MD))
	sb.WriteString(fmt.Sprintf("  --spacing-lg: %s;\n", tm.tokens.Spacing.LG))
	sb.WriteString(fmt.Sprintf("  --spacing-xl: %s;\n", tm.tokens.Spacing.XL))
	sb.WriteString(fmt.Sprintf("  --spacing-xxl: %s;\n", tm.tokens.Spacing.XXL))

	// Typography variables
	sb.WriteString(fmt.Sprintf("  --font-family: %s;\n", tm.tokens.Typography.FontFamily))
	sb.WriteString(fmt.Sprintf("  --font-size-xs: %s;\n", tm.tokens.Typography.FontSize.XS))
	sb.WriteString(fmt.Sprintf("  --font-size-sm: %s;\n", tm.tokens.Typography.FontSize.SM))
	sb.WriteString(fmt.Sprintf("  --font-size-md: %s;\n", tm.tokens.Typography.FontSize.MD))
	sb.WriteString(fmt.Sprintf("  --font-size-lg: %s;\n", tm.tokens.Typography.FontSize.LG))
	sb.WriteString(fmt.Sprintf("  --font-size-xl: %s;\n", tm.tokens.Typography.FontSize.XL))
	sb.WriteString(fmt.Sprintf("  --font-size-xxl: %s;\n", tm.tokens.Typography.FontSize.XXL))
	sb.WriteString(fmt.Sprintf("  --line-height-tight: %s;\n", tm.tokens.Typography.LineHeight.Tight))
	sb.WriteString(fmt.Sprintf("  --line-height-normal: %s;\n", tm.tokens.Typography.LineHeight.Normal))
	sb.WriteString(fmt.Sprintf("  --line-height-relaxed: %s;\n", tm.tokens.Typography.LineHeight.Relaxed))
	sb.WriteString(fmt.Sprintf("  --font-weight-regular: %s;\n", tm.tokens.Typography.FontWeight.Regular))
	sb.WriteString(fmt.Sprintf("  --font-weight-medium: %s;\n", tm.tokens.Typography.FontWeight.Medium))
	sb.WriteString(fmt.Sprintf("  --font-weight-semibold: %s;\n", tm.tokens.Typography.FontWeight.Semibold))
	sb.WriteString(fmt.Sprintf("  --font-weight-bold: %s;\n", tm.tokens.Typography.FontWeight.Bold))

	// Border variables
	sb.WriteString(fmt.Sprintf("  --border-radius-none: %s;\n", tm.tokens.Border.Radius.None))
	sb.WriteString(fmt.Sprintf("  --border-radius-sm: %s;\n", tm.tokens.Border.Radius.SM))
	sb.WriteString(fmt.Sprintf("  --border-radius-md: %s;\n", tm.tokens.Border.Radius.MD))
	sb.WriteString(fmt.Sprintf("  --border-radius-lg: %s;\n", tm.tokens.Border.Radius.LG))
	sb.WriteString(fmt.Sprintf("  --border-radius-xl: %s;\n", tm.tokens.Border.Radius.XL))
	sb.WriteString(fmt.Sprintf("  --border-width-thin: %s;\n", tm.tokens.Border.Width.Thin))
	sb.WriteString(fmt.Sprintf("  --border-width-base: %s;\n", tm.tokens.Border.Width.Base))
	sb.WriteString(fmt.Sprintf("  --border-width-thick: %s;\n", tm.tokens.Border.Width.Thick))

	// Shadow variables
	sb.WriteString(fmt.Sprintf("  --shadow-sm: %s;\n", tm.tokens.Shadows.SM))
	sb.WriteString(fmt.Sprintf("  --shadow-md: %s;\n", tm.tokens.Shadows.MD))
	sb.WriteString(fmt.Sprintf("  --shadow-lg: %s;\n", tm.tokens.Shadows.LG))
	sb.WriteString(fmt.Sprintf("  --shadow-xl: %s;\n", tm.tokens.Shadows.XL))

	// Component-specific variables
	sb.WriteString(fmt.Sprintf("  --button-padding-v: %s;\n", tm.tokens.Components.Button.PaddingVertical))
	sb.WriteString(fmt.Sprintf("  --button-padding-h: %s;\n", tm.tokens.Components.Button.PaddingHorizontal))
	sb.WriteString(fmt.Sprintf("  --button-height: %s;\n", tm.tokens.Components.Button.MinHeight))
	sb.WriteString(fmt.Sprintf("  --button-font-size: %s;\n", tm.tokens.Components.Button.FontSize))
	sb.WriteString(fmt.Sprintf("  --input-padding-v: %s;\n", tm.tokens.Components.Input.PaddingVertical))
	sb.WriteString(fmt.Sprintf("  --input-padding-h: %s;\n", tm.tokens.Components.Input.PaddingHorizontal))
	sb.WriteString(fmt.Sprintf("  --input-height: %s;\n", tm.tokens.Components.Input.MinHeight))
	sb.WriteString(fmt.Sprintf("  --card-padding: %s;\n", tm.tokens.Components.Card.Padding))
	sb.WriteString(fmt.Sprintf("  --card-border-radius: %s;\n", tm.tokens.Components.Card.BorderRadius))
	sb.WriteString(fmt.Sprintf("  --badge-padding-v: %s;\n", tm.tokens.Components.Badge.PaddingVertical))
	sb.WriteString(fmt.Sprintf("  --badge-padding-h: %s;\n", tm.tokens.Components.Badge.PaddingHorizontal))

	sb.WriteString("}\n")
	tm.cssVariables = sb.String()
}

// GetCSSVariablesStyle returns the CSS variables as an HTML style element
func (tm *ThemeManager) GetCSSVariablesStyle() string {
	return "<style>\n" + tm.cssVariables + "</style>"
}

// GetCSSVariablesString returns the raw CSS variables string
func (tm *ThemeManager) GetCSSVariablesString() string {
	return tm.cssVariables
}

// TokensToFuncMap converts theme tokens to template function map for easy access in templates
func TokensToFuncMap(tm *ThemeManager) template.FuncMap {
	tokens := tm.GetTokens()
	return template.FuncMap{
		// Access color tokens
		"colorPrimary":    func() string { return tokens.Colors.Primary },
		"colorSecondary":  func() string { return tokens.Colors.Secondary },
		"colorSuccess":    func() string { return tokens.Colors.Success },
		"colorWarning":    func() string { return tokens.Colors.Warning },
		"colorError":      func() string { return tokens.Colors.Error },
		"colorInfo":       func() string { return tokens.Colors.Info },
		"colorBackground": func() string { return tokens.Colors.Background },
		"colorSurface":    func() string { return tokens.Colors.Surface },
		"colorBorder":     func() string { return tokens.Colors.Border },
		"colorText":       func() string { return tokens.Colors.Text },
		"colorTextMuted":  func() string { return tokens.Colors.TextMuted },

		// Access spacing
		"spacingXS": func() string { return tokens.Spacing.XS },
		"spacingSM": func() string { return tokens.Spacing.SM },
		"spacingMD": func() string { return tokens.Spacing.MD },
		"spacingLG": func() string { return tokens.Spacing.LG },
		"spacingXL": func() string { return tokens.Spacing.XL },

		// Access typography
		"fontSize": func(key string) string {
			switch key {
			case "sm":
				return tokens.Typography.FontSize.SM
			case "md":
				return tokens.Typography.FontSize.MD
			case "lg":
				return tokens.Typography.FontSize.LG
			case "xl":
				return tokens.Typography.FontSize.XL
			default:
				return tokens.Typography.FontSize.MD
			}
		},
		"fontWeight": func(key string) string {
			switch key {
			case "medium":
				return tokens.Typography.FontWeight.Medium
			case "semibold":
				return tokens.Typography.FontWeight.Semibold
			case "bold":
				return tokens.Typography.FontWeight.Bold
			default:
				return tokens.Typography.FontWeight.Regular
			}
		},
		"borderRadius": func(key string) string {
			switch key {
			case "sm":
				return tokens.Border.Radius.SM
			case "md":
				return tokens.Border.Radius.MD
			case "lg":
				return tokens.Border.Radius.LG
			case "xl":
				return tokens.Border.Radius.XL
			default:
				return tokens.Border.Radius.MD
			}
		},
	}
}
