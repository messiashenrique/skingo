package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/messiashenrique/skingo"
)

func main() {
	ts := skingo.NewTemplateSet("layout")

	// Add custom	functions to the template set, available in all templates
	ts.AddFuncs(template.FuncMap{
		"uppercase": strings.ToUpper,
		"lowercase": strings.ToLower,
		"formatPrice": func(price float64) string {
			return fmt.Sprintf("US$ %.2f", price)
		},
	})

	// Analyze the templates in the "templates" directory
	if err := ts.ParseDirs("templates"); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Handler for Home Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title":   "Home Page - Skingo",
			"Content": "Welcome to Skingo!",
			"Price":   159.90,
		}

		if err := ts.Execute(w, "home", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
