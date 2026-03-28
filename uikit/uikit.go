package uikit

import (
	"embed"
	"io/fs"

	"github.com/messiashenrique/skingo"
)

const CatalogName = "skingo-ui-core"

//go:embed templates/components/*.html catalog.json
var embeddedAssets embed.FS

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
