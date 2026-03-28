package main

import (
	"log"
	"net/http"

	"github.com/messiashenrique/skingo"
)

func main() {
	// Makes a new TemplateSet with layout template
	ts := skingo.NewTemplateSet("layout")

	// Optional: register component metadata catalog used by docs/tooling.
	if err := ts.RegisterComponentCatalogFile("skingo-ui", "components/catalog.json"); err != nil {
		log.Fatalf("Error loading component catalog: %v", err)
	}

	// Optional: validate required params and basic types for component calls.
	ts.SetComponentValidation(skingo.ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: true,
	})

	// Analyze the templates in the "templates" and "components" directories
	if err := ts.ParseDirs("templates", "components"); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	for _, component := range ts.ListComponents() {
		log.Printf("Component: %s (catalog=%s, version=%s)", component.Name, component.Catalog, component.Version)
	}

	// Handler for Home Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := ts.Execute(w, "home", map[string]interface{}{
			"Title":   "Home Page",
			"Content": "Welcome to Skingo!",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
