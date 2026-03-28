package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/messiashenrique/skingo"
	"github.com/messiashenrique/skingo/uikit"
)

//go:embed templates/*.html
var appFS embed.FS

func main() {
	// Create a new template set with "layout" as the layout template
	ts := skingo.NewTemplateSet("layout")

	// Register the uikit catalog and theme system
	if err := uikit.RegisterCatalog(ts); err != nil {
		log.Fatalf("Error registering catalog: %v", err)
	}

	// Register the theme system (light theme by default)
	if err := uikit.RegisterTheme(ts, "light"); err != nil {
		log.Fatalf("Error registering theme: %v", err)
	}

	// Enable component validation
	ts.SetComponentValidation(skingo.ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: true,
	})

	// Parse local templates and uikit components using ParseManyFS
	err := ts.ParseManyFS(
		skingo.ParseFSSource{Filesystem: appFS, Roots: []string{"templates"}},
		uikit.Source(),
	)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Single route with query parameter for theme selection
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read theme from query parameter, default to "light"
		theme := r.URL.Query().Get("theme")
		if theme != "dark" && theme != "light" {
			theme = "light"
		}

		// Set the global theme
		uikit.SetGlobalTheme(theme)

		// Determine titles based on theme
		title := "Themed Components - Light Mode"
		description := "This page demonstrates the light theme with component variants"
		if theme == "dark" {
			title = "Themed Components - Dark Mode"
			description = "This page demonstrates the dark theme with component variants"
		}

		if err := ts.Execute(w, "home", map[string]interface{}{
			"Title":          title,
			"Description":    description,
			"Theme":          theme,
			"ThemeVarsInput": theme, // Pass theme to template for client JS
		}); err != nil {
			log.Printf("Error rendering page: %v", err)
		}
	})

	log.Println("Theme-aware server running on http://localhost:8080")
	log.Println("  Light theme: http://localhost:8080?theme=light")
	log.Println("  Dark theme:  http://localhost:8080?theme=dark")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
