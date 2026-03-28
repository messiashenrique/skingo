package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/messiashenrique/skingo"
	"github.com/messiashenrique/skingo/uikit"
)

func main() {
	// Create a new template set
	ts := skingo.NewTemplateSet("layout")

	// Register the uikit catalog
	if err := uikit.RegisterCatalog(ts); err != nil {
		log.Fatal(err)
	}

	// List registered components BEFORE parsing to see what's registered
	components := ts.ListComponents()
	fmt.Printf("Components registered BEFORE parsing: %d\n", len(components))
	for _, comp := range components {
		fmt.Printf("  - %s [%s]\n", comp.Name, comp.Catalog)
	}

	// Enable component validation to help catch errors
	ts.EnableComponentValidation(true)

	// Parse templates from the current directory
	if err := ts.ParseDirs("./templates"); err != nil {
		log.Fatal(err)
	}

	// List registered components AFTER parsing
	components = ts.ListComponents()
	fmt.Printf("Components registered AFTER parsing: %d\n", len(components))
	for _, comp := range components {
		fmt.Printf("  - %s [%s]\n", comp.Name, comp.Catalog)
	}

	// Set up HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Execute the home template
		err := ts.Execute(w, "home", map[string]interface{}{
			"title": "Component Helpers Demo",
		})
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v\n", err)
			return
		}
	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
