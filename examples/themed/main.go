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

	// Route for light theme
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uikit.SetGlobalTheme("light")
		if err := ts.Execute(w, "home", map[string]interface{}{
			"Title":       "Themed Components - Light Mode",
			"Description": "This page demonstrates the light theme with component variants",
			"Theme":       "light",
		}); err != nil {
			log.Printf("Error rendering light theme: %v", err)
		}
	})

	// Route for dark theme
	http.HandleFunc("/dark", func(w http.ResponseWriter, r *http.Request) {
		uikit.SetGlobalTheme("dark")
		if err := ts.Execute(w, "home", map[string]interface{}{
			"Title":       "Themed Components - Dark Mode",
			"Description": "This page demonstrates the dark theme with component variants",
			"Theme":       "dark",
		}); err != nil {
			log.Printf("Error rendering dark theme: %v", err)
		}
	})

	log.Println("Theme-aware server running on http://localhost:8080")
	log.Println("  Light theme: http://localhost:8080")
	log.Println("  Dark theme:  http://localhost:8080/dark")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
