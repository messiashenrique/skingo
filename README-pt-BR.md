---
<h1 align="center">
  <picture>
    <img height="125" alt="Skingo" src="https://raw.githubusercontent.com/messiashenrique/skingo/refs/heads/main/docs/static/img/skingo-logo.svg">
  </picture>
</h1>
---

# skingo
Simples proposta para usar modelos HTML em Go

Skingo é um pacote Go que estende o pacote `html/template` padrão com funcionalidades de componentes, escopo de CSS, inclusão automática de JS e muito mais.

Skingo foi inspirado na forma simples e clara de serapação entre HTML, CSS e JS adotada nas páginas e componentes do Vue.js.

## Características

- 🧩 Sistema de componentes reutilizáveis
- 🎨 Escopo automático de CSS para evitar conflitos
- 📦 Agrupamento automático de CSS e JS
- 🔍 Rastreamento inteligente de dependências
- 🚀 Template layouts

## Instalação

```bash
go get github.com/messiashenrique/skingo
```

## Como usar

### Exemplo básico
```go
//main.go
package main

import (
    "log"
    "net/http"
    "github.com/messiashenrique/skingo"
)

func main() {
    // Cria um novo conjunto de templates com "layout" como template de layout
    ts := skingo.NewTemplateSet("layout")
    
    // Analisa os templates no diretório "templates"
    if err := ts.ParseDirs("templates"); err != nil {
        log.Fatalf("Erro ao analisar templates: %v", err)
    }
    
    // Handler para a página inicial
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := ts.Execute(w, "home", map[string]interface{}{
            "Title": "Página Inicial",
            "Content": "Bem-vindo ao Skingo!",
        }); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    
    log.Println("Servidor rodando em http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Componentes

Skingo permite criar componentes reutilizáveis que encapsulam HTML, CSS e JavaScript.

### Definindo um componente
Componente com parâmtros posicionais e 2º parâmetro opcional
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
  console.log("Botão carregado!");
</script>
```
Componente com parâmetros nomeados
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

### Usando um componente
Usando os componetes na Página principal e também componentes aninhados.
```html
<!-- templates/home.html -->
<template>
  <div class="container">
    <h1>{{.Title}}</h1>
    <p>{{.Content}}</p>
    
    <!-- Usando componentes com parâmetros nomeados -->
    {{ comp "card.html" (dict 
      "title" "Exemplo de Card" 
      "content" "Este é um exemplo de um componente de card com um botão." 
      "buttonText" "Ler mais"
    ) }}
    
    {{ comp "card.html" (dict 
      "title" "Outro Card" 
      "content" "Os componentes podem ser facilmente reutilizados com diferentes conteúdos." 
      "buttonText" "Saiba mais"
    ) }}

    <!-- Usando componente com parâmetros posicionais e 2º parâmtro opicional -->
    {{ comp "button.html" "Clique-me!" "green" }}
  </div>
</template>
```

## API

### NewTemplateSet
```go
func NewTemplateSet(layoutName string) *TemplateSet
```
Cria um novo conjunto de templates usando o template especificado como layout.

### ParseDirs
```go
func (ts *TemplateSet) ParseDirs(dirs ...string) error
```
Analisa todos os arquivos HTML/templates nos diretórios especificados.

### Execute
```go
func (ts *TemplateSet) Execute(w io.Writer, name string, data interface{}) error
```
Renderiza o template especificado usando o layout configurado.

### ExecuteIsolated
```go
func (ts *TemplateSet) ExecuteIsolated(w io.Writer, filename string, data interface{}) error
```
Renderiza um template de forma isolada, sem usar o layout. Útil para HTMX e requisições Ajax.
**Nota:** `ExecuteIsolated` não faz separação de escopos JS e CSS. Portanto, o recomendado é que os estilos sejam declarados globalmente.

## Licença
MIT







