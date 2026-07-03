package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/messiashenrique/skingo"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	ts := skingo.NewTemplateSet("layout")
	if err := ts.ParseFS(templateFS, "templates"); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := ts.Execute(w, "home", map[string]interface{}{
			"Title":   "Public Home",
			"Message": "Rendered with the default layout.",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if err := ts.ExecuteWithLayout(w, "admin", "home", map[string]interface{}{
			"Title":   "Admin Home",
			"Message": "The same template rendered with the admin layout.",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
