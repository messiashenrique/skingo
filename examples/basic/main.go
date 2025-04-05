package main

import (
	"log"
	"net/http"

	"github.com/messiashenrique/skingo"
)

func main() {
	// Makes a new TemplateSet with layout template
	ts := skingo.NewTemplateSet("layout")

	// Analyze the templates in the "templates" directory
	if err := ts.ParseDirs("templates"); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Handler for Home Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := ts.Execute(w, "home", map[string]interface{}{"Title": "Home Page"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handler for About Page
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		if err := ts.Execute(w, "about", map[string]interface{}{"Title": "About Page"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
