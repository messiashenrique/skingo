<h1 align="center">
  <picture>
    <img height="72" alt="Skingo" src="docs/static/img/skingo-logo.svg">
  </picture>
</h1>

🌏 English | 🇧🇷 **[Português](README-pt-BR.md)**

# skingo
Simple Proposal for Using HTML Templates in Go

Skingo is a Go package that extends the standard `html/template` libray with component functionality, CSS scoping, JS auto-inclusion, and more.

Skingo was inspired by the simple and clean way of interfacing HTML, CSS, and JS that Vue.js pages and components use.

## Features

- 🧩 Reusable component system
- 🎨 Automatic CSS scoping to avoid conflicts
- 📦 Automatic CSS and JS grouping
- 🔍 Smart dependency tracking
- 🚀 Template layouts

## Installation

```bash
go get github.com/messiashenrique/skingo
```

## How to use

### Basic example
```go
//main.go
package main

import (
    "log"
    "net/http"
    "github.com/messiashenrique/skingo"
)

func main() {
    // Makes a new template set with "layout" as the layout template
    ts := skingo.NewTemplateSet("layout")
    
    // Analyze the templates in the "templates" directory
    if err := ts.ParseDirs("templates"); err != nil {
        log.Fatalf("Error parsing templates: %v", err)
    }
    
    // Handler to home page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := ts.Execute(w, "home", map[string]interface{}{
            "Title": "Home Page",
            "Content": "Welcome to Skingo!",
        }); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Layout

Skingo allows flexible layout usage. Therefore, the only mandatory item is to define the `{{ .Yield }}` variable as the entry point for rendering templates that use this layout.

The CSS and JavaScript codes declared in the layout will have global scope.

An example of a layout can be seen below:

### Defining a Layout
```html
<!-- templates/layout.html -->
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Skingo</title>
</head>
<body>
	{{ .Yield }}
</body>
</html>
```
To define the above file as a layout, simply insert the file name into the template set creation call by doing `ts := skingo.NewTemplateSet("layout")`.

Don't forget to include the directory where the layout file is located in the `ParseDirs` function.

## Components

Skingo lets you create reusable components that encapsulate HTML, CSS, and JavaScript.

### Defining a component
Component with positional parameters and optional 2nd parameter
```html
<!-- templates/button.html -->
<template>
  <button class="btn {{ paramOr 1 "blue"}}">{{ param 0 }}</button>
</template>

<style>
  .btn {
    margin: 0.5rem 0;
    padding: 0.5rem 1rem;
    color: white;
    border-radius: 0.25rem;
    border: none;
    cursor: pointer;
  }

  .blue {
    background-color: #3490dc;
  }
  
  .green {
    background-color: #019001;
  }
</style>

<script>
  console.log("Loaded button!");
</script>
```
### Writing Components

Components integrate seamlessly with templates. Here's a button component that leverages helper syntax:

```html
<!-- templates/button.html -->
<template>
  <button class="btn {{ paramOr 1 "blue"}}">{{ param 0 }}</button>
</template>

<style>
  .btn {
    padding: 0.5rem 1rem;
    color: white;
    border-radius: 0.25rem;
    border: none;
    cursor: pointer;
  }
  .blue {
    background-color: #3490dc;
  }
  .green {
    background-color: #019001;
  }
</style>
```

And a card component that nests other components:

```html
<!-- templates/card.html -->
<template>
  <div class="card">
    <div class="card-header">
      <h3>{{.title}}</h3>
    </div>
    <div class="card-body">
      <p>{{.content}}</p>
    </div>
    <div class="card-footer">
      <!-- Using helper to nest button component -->
      {{ button .buttonText }}
    </div>
  </div>
</template>

<style>
  .card { border: 1px solid #e2e8f0; border-radius: 0.5rem; }
  .card-header { background-color: #f7fafc; padding: 0.5rem; }
  .card-body { padding: 0.5rem 1rem; }
  .card-footer { background-color: #f7fafc; }
</style>
```

### Using a component (Helper Syntax)

Skingo automatically generates helper functions for all registered components, providing a clean and intuitive syntax:

```html
<!-- templates/home.html -->
<template>
  <div class="container">
    <h1>{{.Title}}</h1>
    <p>{{.Content}}</p>

    <!-- Using component helpers with named parameters -->
    {{ card (dict 
      "title" "Card Example" 
      "content" "This is an example of a card component with a button." 
      "buttonText" "Read more"
    ) }}
    
    <!-- Using positional parameters -->
    {{ button "Click me!" "green" }}
  </div>
</template>
```

**Instead of:** `{{ comp "card.html" (dict "title" "...") }}`  
**You can now write:** `{{ card (dict "title" "...") }}`

Helper functions are automatically generated based on:
- Component name (derived from template filename or registered component name)
- Registered component metadata
- The `comp` function is still available as a fallback

Skingo will intelligently determine the CSS scopes and automatically create classes that help in styling each component, respecting first the specific styles.

If more than one element without a parent (without a container) are declared between the `<template><template>` tags, Skingo will automatically create a container (`<div>`) to wrap them and thus intelligently separate the styles between the different components, respecting each scope.

To avoid this behavior above, simply add the `unwrap` attribute to the "template" tag, like this: `<template unwrap>`.

### Example with Embedded Filesystem
```go
//main.go
package main

import (
    "embed"
    "log"
    "net/http"
    "github.com/messiashenrique/skingo"
)

//go:embed templates/**/*.html
var templateFS embed.FS

func main() {
    // Create a new template set with "layout" as the layout template
    ts := skingo.NewTemplateSet("layout")
    
    // Parse templates in the embedded filesystem
    if err := ts.ParseFS(templateFS, "templates/pages", "templates/components"); err != nil {
        log.Fatalf("Error parsing templates: %v", err)
    }
    
    // Handler for the home page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := ts.Execute(w, "home", map[string]interface{}{
            "Title": "Home Page",
            "Content": "Welcome to Skingo with embedded templates!",
        }); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    
    // Handler for HTMX requests that only need fragments
    http.HandleFunc("/fragment", func(w http.ResponseWriter, r *http.Request) {
        if err := ts.ExecuteIsolatedFS(w, templateFS, "templates/fragments/partial.html", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## API

### NewTemplateSet
```go
func NewTemplateSet(layoutName string) *TemplateSet
```
Makes a new template set using the specified template as the layout.

### ParseDirs
```go
func (ts *TemplateSet) ParseDirs(dirs ...string) error
```
Parses all HTML/templates files in the specified directories.

### ParseFS

```go
func (ts *TemplateSet) ParseFS(filesystem fs.FS, roots ...string) error
```
Parses all HTML/template files in an embedded filesystem.

### ParseManyFS

```go
type ParseFSSource struct {
  Filesystem fs.FS
  Roots      []string
}

func (ts *TemplateSet) ParseManyFS(sources ...ParseFSSource) error
```

Parses templates from multiple filesystems in one pass. This is useful for hybrid setups where your app templates and a component catalog are embedded in different packages.

At least one source must provide the configured layout template.

### Execute
```go
func (ts *TemplateSet) Execute(w io.Writer, name string, data interface{}) error
```
Renders the specified template using the configured layout.

### ExecuteIsolated
```go
func (ts *TemplateSet) ExecuteIsolated(w io.Writer, filename string, data interface{}) error
```
Renders a template in isolation, without using the layout. Useful for HTMX and Ajax requests.
* **Note:** `ExecuteIsolated` does not separate CSS scope. Therefore, it is recommended that styles be declared globally.

Although `ExecuteIsolated` load the template on demand, it uses caching so that if it needs to execute the template again, it is already in memory, thus optimizing performance.

### ExecuteIsolatedFS
```go
func (ts *TemplateSet) ExecuteIsolatedFS(w io.Writer, filesystem fs.FS, fsPath string, data interface{}) error
```
Renders a template directly from an embedded filesystem, without using the configured layout.

This method is similar to ExecuteIsolated, but works with embedded filesystems.

It is ideal for use with 'HTMX', Ajax requests, or any scenario where only a single HTML fragment
is required.

The 'fsPath' parameter should be the path within the filesystem.

### Component Catalog Metadata
Skingo now provides optional catalog metadata registration APIs as a foundation for reusable UI packs.
These APIs do not change rendering behavior and are useful for docs, tooling and future validation.

```go
func (ts *TemplateSet) RegisterComponentMeta(name string, meta ComponentMeta) error
func (ts *TemplateSet) RegisterComponentCatalog(catalogName string, components map[string]ComponentMeta) error
func (ts *TemplateSet) RegisterComponentCatalogJSON(catalogName string, manifest []byte) error
func (ts *TemplateSet) RegisterComponentCatalogFile(catalogName string, filename string) error
func (ts *TemplateSet) RegisterComponentCatalogFS(catalogName string, filesystem fs.FS, manifestPath string) error
func (ts *TemplateSet) ListComponents() []ComponentInfo
func (ts *TemplateSet) GetComponentMeta(name string) (ComponentMeta, bool)
```

JSON manifest format example:

```json
{
  "components": {
    "button": {
      "description": "Clickable action trigger",
      "version": "1.0.0",
      "variants": ["solid", "outline", "ghost"],
      "dependencies": ["icon"],
      "params": [
        {
          "name": "label",
          "type": "string",
          "required": true,
          "description": "Button label"
        }
      ]
    }
  }
}
```

### Component Validation (Optional)
You can enable runtime validation for component calls based on registered metadata.

```go
type ComponentValidationOptions struct {
  Enabled     bool
  StrictTypes bool
}

func (ts *TemplateSet) SetComponentValidation(options ComponentValidationOptions)
func (ts *TemplateSet) EnableComponentValidation(enabled bool)
func (ts *TemplateSet) GetComponentValidation() ComponentValidationOptions
```

Validation behavior:
- `Enabled=false` (default): no validation.
- `Enabled=true`: validates required params.
- `StrictTypes=true` (default): validates basic declared types (`string`, `bool`, `int`, `float`, `number`, `[]string`, `[]map[string]string`, `map[string]interface{}`).
- If the component has a `variant` param and metadata includes `variants`, the value is validated against allowed variants.

Example:

```go
ts.SetComponentValidation(skingo.ComponentValidationOptions{
  Enabled:     true,
  StrictTypes: true,
})
```

### Hybrid Catalog Example

Skingo includes an initial reusable UI catalog package in `uikit` with pre-built components:
- `SkButton` - Styled button with variants (primary, outline, ghost)
- `SkInput` - Form input with label support
- `SkBadge` - Status badge with semantic variants (success, warning, danger)
- `SkInfo` - Alert/info box with variants (info, success, error)
- `SkCard` - Container card with header, content, and optional footer action

**Example usage:**
```go
import (
  "github.com/messiashenrique/skingo"
  "github.com/messiashenrique/skingo/uikit"
)

func main() {
  ts := skingo.NewTemplateSet("layout")
  
  // Register the uikit catalog
  if err := uikit.RegisterCatalog(ts); err != nil {
    log.Fatal(err)
  }
  
  // Enable optional validation
  ts.SetComponentValidation(skingo.ComponentValidationOptions{
    Enabled:     true,
    StrictTypes: true,
  })
  
  // Parse and use
  if err := ts.ParseDirs("templates"); err != nil {
    log.Fatal(err)
  }
  
  // SkButton, SkInput, SkCard helpers are now available in templates
}
```

**In templates:**
```html
{{ SkButton "Click me" "primary" }}
{{ SkInput (dict "name" "email" "label" "Email") }}
{{ SkCard (dict "title" "My Card" "content" "Content here") }}
{{ SkInfo (dict "title" "Info" "message" "Hello!" "variant" "success") }}
{{ SkBadge "Active" "success" }}
```

See the hybrid integration example in `examples/hybrid`.
  "github.com/messiashenrique/skingo"
  "github.com/messiashenrique/skingo/uikit"
)

ts := skingo.NewTemplateSet("layout")

_ = uikit.RegisterCatalog(ts)

err := ts.ParseManyFS(
  skingo.ParseFSSource{Filesystem: appFS, Roots: []string{"templates"}},
  uikit.Source(),
)
```

## Template Functions

Skingo offers several auxiliary functions for use in templates.

### Default Functions

Skingo includes the following standard functions available in all templates:

| Function | Description | Example |
|--------|-----------|---------|
| `add` | Adds two numbers | `{{add 3 5}}` → `8` |
| `sub` | Subtracts two numbers | `{{sub 10 4}}` → `6` |
| `mul` | Multiplies two numbers | `{{mul 3 5}}` → `15` |
| `mod` | Returns the remainder of the division | `{{mod 10 3}}` → `1` |
| `addFloat` | Adds two floating point numbers | `{{addFloat 3.0 3.1}}` → `6.1` |
| `subFloat` | Subtract two floating point numbers | `{{subFloat 7.3 3.1}}` → `4.2` |
| `mulFloat` | Multiplies two floating point numbers | `{{mulFloat 3.0 7.1}}` → `21.3` |
| `divFloat` | Divides two floating point numbers | `{{divFloat 24.6 3.0}}` → `8.2` |
| `comp` | Invokes a component passing parameters | `{{comp "card" "Black Card"}}` |
| `dict` | Creates a key/value map | `{{comp "button" (dict "text" "Click")}}` |
| `param` | Accesses a positional parameter | `{{param 0}}` |
| `paramOr` | Accesses a positional parameter with default value | `{{paramOr 1 "Default"}}` |
| `toJson` | Converts a value to JSON | `{{toJson .user}}` → `{"name":"John"}` |

### Adding Custom Functions

You can add your own functions for use in templates:

```go
ts := skingo.NewTemplateSet("layout")

ts.AddFuncs(template.FuncMap{
    "uppercase": strings.ToUpper,
    "lowercase": strings.ToLower,
    "formatDate": func(date time.Time) string {
        return date.Format("02/01/2006")
    },
})
```
* **Note**: This method should be called before `ParseDirs`.

## Component Testing

Skingo includes built-in APIs for testing component metadata and rendering:

```go
func TestComponentMetadata(t *testing.T) {
    ts := skingo.NewTemplateSet("layout")
    
    // Register a component catalog
    if err := skingo.RegisterComponentCatalogJSON(ts, "mycomponents", []byte(`{
        "components": {
            "button": {
                "description": "Clickable button",
                "variables": [
                    {"name": "label", "type": "string", "required": true}
                ]
            }
        }
    }`)); err != nil {
        t.Fatal(err)
    }
    
    // Test metadata retrieval
    meta, ok := ts.GetComponentMeta("button")
    if !ok {
        t.Fatal("Component metadata not found")
    }
    
    if meta.Description != "Clickable button" {
        t.Errorf("Expected 'Clickable button', got %s", meta.Description)
    }
}

func TestComponentValidation(t *testing.T) {
    ts := skingo.NewTemplateSet("layout")
    
    // Enable validation
    ts.SetComponentValidation(skingo.ComponentValidationOptions{
        Enabled:     true,
        StrictTypes: true,
    })
    
    // Validation will now check required params and types during component execution
}
```

Run tests with standard Go tooling:
```bash
go test ./...          # Run all tests
go test -v .           # Run tests in current package with verbose output
go test -cover ./...   # Run tests with coverage report
```

See `skingo_test.go` for comprehensive examples of testing:
- Component metadata registration and retrieval
- Multi-filesystem catalog parsing with `ParseManyFS`
- Validation options and type checking
- Component helper function generation

## Roadmap for Development

| Stage | Description | Priority | Status |
|-------|-----------|------------|--------|
| **Tests** | Implementation of comprehensive unit tests | High | ✅ Complete |
| **Performance Optimization** | Refactoring to improve rendering efficiency | High | 📅 Planned |
| **Full Documentation** | Detailed documentation with examples for each feature | High | 🔄 In progress |
| **HTMX Integration** | Improved support for HTMX with dedicated helpers | High | 📅 Planned |
| **Themed Variants** | Component variants with light/dark/custom theme support | High | 📅 Planned |
| **Design Tokens** | Centralized design token system for uikit components | High | 📅 Planned |
| **Advanced Examples** | Repository with more complex examples and real use cases | Medium | 📅 Planned |
| **Hot Reload** | Support for hot reload during development | Medium | 🔮 Considering |
| **Benchmarks** | Performance comparison with other solutions | Medium | 📅 Planned |
| **CSS/JS Minification** | Automatic minification of CSS and JS in production | Medium | 📅 Planned |
| **Extensions for Tools** | Plugins for IDEs and integrations with development tools | Low | 🔮 Considering |
| **Server Side Rendering** | Implementation of SSR optimized for SPAs | Low | 🔮 Considering |


### Caption
- 🔄 **In progress**: Development has started
- 📅 **Planned**: Planned for implementation soon
- 🔮 **Considering**: Being considered for the future

## License
MIT







