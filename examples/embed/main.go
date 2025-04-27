package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/messiashenrique/skingo"
)

//go:embed templates/* components/*
var templateFS embed.FS

func main() {
	// Makes a new TemplateSet with layout template
	ts := skingo.NewTemplateSet("layout")

	// Analyze the templates in the "templates" and "components" directories
	// Use ParseFS instead of ParseDir
	// if err := ts.ParseDirs("templates", "components"); err != nil {
	if err := ts.ParseFS(templateFS, "templates", "components"); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Handler for Home Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := ts.Execute(w, "home", map[string]interface{}{
			"Title":   "Embedded Templates",
			"Content": "Welcome to Skingo!",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
