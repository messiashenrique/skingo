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
Component with named parameters
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
      <!-- Using component with positional parameters -->
      {{ comp "button.html" .buttonText }}
    </div>
  </div>
</template>

<style>
  .card {
    border: 0.0625rem solid #e2e8f0;
    border-radius: 0.5rem;
    overflow: hidden;
    margin-bottom: 1rem;
    box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.1);
  }

  .card-header {
    background-color: #f7fafc;
    padding: 0.5rem;
    border-bottom: 0.0625rem solid #e2e8f0;
  }

  .card-header h3 {
    margin: 0;
    font-size: 1.25rem;
  }

  .card-body {
    padding: 0.5rem 1rem;
  }

  .card-footer {
    padding: 0.25rem;
    background-color: #f7fafc;
    border-top: 0.0625rem solid #e2e8f0;
  }
</style>
```

### Using a component
```html
Using the components on the Home Page and also nested components.
<!-- templates/home.html -->
<template>
  <div class="container">
    <h1>{{.Title}}</h1>
    <p>{{.Content}}</p>

    <!-- Using components with named parameters -->
    {{ comp "card.html" (dict 
      "title" "Card Example" 
      "content" "This is an example of a card component with a button." 
      "buttonText" "Read more"
    ) }}
    
    {{ comp "card.html" (dict 
      "title" "Other Card" 
      "content" "Components can be easily reused with different content." 
      "buttonText" "Find out more"
    ) }}
    
    <!-- Using component with positional parameters and optional 2nd parameter -->
    {{ comp "button.html" "Click me!" "green" }}
    
  </div>
</template>
```

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

## Roadmap for Development

| Stage | Description | Priority | Status |
|-------|-----------|------------|--------|
| **Tests** | Implementation of comprehensive unit tests | High | 🔄 In progress |
| **Performance Optimization** | Refactoring to improve rendering efficiency | High | 📅 Planned |
| **Full Documentation** | Detailed documentation with examples for each feature | High | 🔄 In progress |
| **HTMX Integration** | Improved support for HTMX with dedicated helpers | High | 📅 Planned |
| **Advanced Examples** | Repository with more complex examples and real use cases | Medium | 📅 Planned |
| **Hot Reload** | Support for hot reload during development | Medium | 🔮 Considering |
| **Parameter Validation** | Parameter validation system for components | Medium | 📅 Planned |
| **Benchmarks** | Performance comparison with other solutions | Medium | 📅 Planned |
| **CSS/JS Minification** | Automatic minification of CSS and JS in production | Medium | 📅 Planned |
| **Extensions for Tools** | Plugins for IDEs and integrations with development tools | Low | 🔮 Considering |
| **Server Side Rendering** | Implementation of SSR optimized for SPAs | Low | 🔮 Considering |
| **Integrated Design System** | Base components to facilitate the creation of consistent interfaces | Low | 🔮 Considering |

### Caption
- 🔄 **In progress**: Development has started
- 📅 **Planned**: Planned for implementation soon
- 🔮 **Considering**: Being considered for the future

## License
MIT







