---
<h1 align="center">
  <picture>
    <img height="125" alt="Skingo" src="https://raw.githubusercontent.com/messiashenrique/skingo/refs/heads/main/docs/static/img/skingo-logo.svg">
  </picture>
</h1>
---

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
**Note:** `ExecuteIsolated` does not separate JS and CSS scopes. Therefore, it is recommended that styles be declared globally.

## License
MIT







