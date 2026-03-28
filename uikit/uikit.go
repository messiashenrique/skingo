package uikit

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/messiashenrique/skingo"
)

const CatalogName = "skingo-ui-core"

//go:embed templates/components/*.html catalog.json
var embeddedAssets embed.FS

// ThemeManager instance for managing themes globally
var defaultThemeManager = NewThemeManager()

// FS returns the embedded filesystem containing UI templates and metadata.
func FS() fs.FS {
	return embeddedAssets
}

// Roots returns template roots that should be parsed for component templates.
func Roots() []string {
	return []string{"templates/components"}
}

// Source returns a ParseFSSource ready to be used with skingo.ParseManyFS.
func Source() skingo.ParseFSSource {
	return skingo.ParseFSSource{
		Filesystem: embeddedAssets,
		Roots:      Roots(),
	}
}

// RegisterCatalog registers the embedded component manifest in a TemplateSet.
func RegisterCatalog(ts *skingo.TemplateSet) error {
	return ts.RegisterComponentCatalogFS(CatalogName, embeddedAssets, "catalog.json")
}

// RegisterTheme registers the theme manager and injects CSS variables into the TemplateSet.
// This adds theme-related functions to templates and injects CSS variables.
func RegisterTheme(ts *skingo.TemplateSet, themeName string) error {
	// Set the theme
	if err := defaultThemeManager.SetTheme(themeName); err != nil {
		return err
	}

	// Add theme-related functions
	ts.AddFuncs(template.FuncMap{
		"themeVars": func() template.HTML {
			return template.HTML(defaultThemeManager.GetCSSVariablesString())
		},
		"currentTheme": func() string {
			return defaultThemeManager.GetCurrentTheme()
		},
		"setTheme": func(name string) error {
			return defaultThemeManager.SetTheme(name)
		},
		"getThemeTokens": func() *DesignTokens {
			return defaultThemeManager.GetTokens()
		},
	})

	// Add token accessor functions
	ts.AddFuncs(TokensToFuncMap(defaultThemeManager))

	return nil
}

// SetGlobalTheme changes the global theme
func SetGlobalTheme(name string) error {
	return defaultThemeManager.SetTheme(name)
}

// GetGlobalTheme returns the current global theme
func GetGlobalTheme() string {
	return defaultThemeManager.GetCurrentTheme()
}

// GetThemeManager returns the default theme manager instance
func GetThemeManager() *ThemeManager {
	return defaultThemeManager
}
