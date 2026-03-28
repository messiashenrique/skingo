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
	ts := skingo.NewTemplateSet("layout")

	if err := uikit.RegisterCatalog(ts); err != nil {
		log.Fatalf("Error loading UI catalog metadata: %v", err)
	}

	ts.SetComponentValidation(skingo.ComponentValidationOptions{
		Enabled:     true,
		StrictTypes: true,
	})

	if err := ts.ParseManyFS(
		skingo.ParseFSSource{Filesystem: appFS, Roots: []string{"templates"}},
		uikit.Source(),
	); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	for _, component := range ts.ListComponents() {
		log.Printf("Registered: %s [%s]", component.Name, component.Catalog)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Skingo Hybrid Catalog",
		}
		if err := ts.Execute(w, "home", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
