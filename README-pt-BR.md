
<h1 align="center">
  <picture>
    <img height="72" alt="Skingo" src="docs/static/img/skingo-logo.svg">
  </picture>
</h1>

🌏 **[English](README.md)** | 🇧🇷 Português

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

## Layout 

O Skingo permite a utilização flexível de layout. Assim, o único ponto obrigatório é definir a variável `{{ .Yield }}` como ponto de entrada para a renderização dos templates que utilizarem esse layout.

Os código CSS e JavaScript declarados no layout terão escopo global.

Um exemplo de layout pode ser visto a seguir:

### Definindo um Layout 
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
Para definir o arquivo acima como layout, basta inserir o nome do arquivo na chamada de criação do conjunto de templates, fazendo `ts := skingo.NewTemplateSet("layout")`. 

Não de esqueça de incluir na função `ParseDirs` o diretório onde está localizado o arquivo do layout.


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

O Skingo vai, de forma inteligente, determinar os escopos de CSS e criar automaticamente classes que auxiliam na estilização de cada componente, respeitando os estilos específicos em primeiro lugar.

Se mais de um elemento sem pai (sem um contêiner) forem declarados entre as tags `<template><template>`, o Skingo criará de forma automática um cointêiner (`<div>`) para envolvê-los e assim separar inteligentemente os estilos entre os diversos componentes, respeitando cada escopo. 

Para evitar esse comportamento acima, basta adicionar o atributo `unwrap` na tag "template", dessa forma: `<template unwrap>`.

### Exemplo com Filesystem Embutido
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
    // Cria um novo conjunto de templates com "layout" como o template de layout
    ts := skingo.NewTemplateSet("layout")
    
    // Analisa os templates no filesystem embutido
    if err := ts.ParseFS(templateFS, "templates/pages", "templates/components"); err != nil {
        log.Fatalf("Error parsing templates: %v", err)
    }
    
    // Handler para a página inicial
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := ts.Execute(w, "home", map[string]interface{}{
            "Title": "Home Page",
            "Content": "Welcome to Skingo with embedded templates!",
        }); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    
    // Handler para requisições HTMX que precisam apenas de fragmentos
    http.HandleFunc("/fragment", func(w http.ResponseWriter, r *http.Request) {
        if err := ts.ExecuteIsolatedFS(w, templateFS, "templates/fragments/partial.html", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    
    log.Println("Servidor rodando em http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
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

### ParseFS
```go
func (ts *TemplateSet) ParseFS(filesystem fs.FS, roots ...string) error
```
Analisa todos os arquivos HTML/template em um sistema de arquivos embutido (embedded filesystem).

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
* **Nota:** `ExecuteIsolated` não faz separação de escopo CSS. Portanto, o recomendado é que os estilos sejam declarados globalmente.

Embora o `ExecuteIsolated` carregue o template sob demanda, ele usa o armazenamento em cache para, caso precise executar novamente o template, ele já esteja em memória, otimizando assim a performance.

### ExecuteIsolatedFS
```go
func (ts *TemplateSet) ExecuteIsolatedFS(w io.Writer, filesystem fs.FS, fsPath string, data interface{}) error
```
Renderiza um template diretamente de um sistema de arquivos embutido, sem usar o layout configurado.

Este método é semelhante ao ExecuteIsolated, mas funciona com sistemas de arquivos embutidos.
É ideal para uso com 'HTMX', requisições Ajax, ou qualquer cenário onde apenas um fragmento HTML
é necessário.

O parâmetro 'fsPath' deve ser o caminho dentro do sistema de arquivos.

## Funções de Template

O Skingo oferece diversas funções auxiliares para uso nos templates.

### Funções Padrão

O Skingo inclui as seguintes funções padrão disponíveis em todos os templates:

| Função | Descrição | Exemplo |
|--------|-----------|---------|
| `add` | Soma dois números | `{{add 3 5}}` → `8` |
| `sub` | Subtrai dois números | `{{sub 10 4}}` → `6` |
| `mul` | Multiplica dois números | `{{mul 3 5}}` → `15` |
| `mod` | Retorna o resto da divisão | `{{mod 10 3}}` → `1` |
| `comp` | Invoca um componente passando parâmetros | `{{comp "card" "Black Card"}}` |
| `dict` | Cria um mapa de chave/valor | `{{comp "button" (dict "text" "Clique")}}` |
| `param` | Acessa um parâmetro posicional | `{{param 0}}` |
| `paramOr` | Acessa um parâmetro posicional com valor padrão | `{{paramOr 1 "Padrão"}}` |
| `toJson` | Converte um valor para JSON | `{{toJson .user}}` → `{"name":"João"}` |

### Adicionando Funções Customizadas

Você pode adicionar suas próprias funções para uso nos templates:

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
* **Nota**: Este método deve ser chamado antes de `ParseDirs`.

## Roteiro de Desenvolvimento

| Etapa | Descrição | Prioridade | Status |
|-------|-----------|------------|--------|
| **Testes** | Implementação de testes unitários abrangentes | Alta | 🔄 Em progresso |
| **Otimização de Performance** | Refatoração para melhorar a eficiência de renderização | Alta | 📅 Planejado |
| **Documentação Completa** | Documentação detalhada com exemplos para cada funcionalidade | Alta | 🔄 Em progresso |
| **Integração HTMX** | Suporte aprimorado para HTMX com helpers dedicados | Alta | 📅 Planejado |
| **Exemplos Avançados** | Repositório com exemplos mais complexos e casos de uso reais | Média | 📅 Planejado |
| **Hot Reload** | Suporte para hot reload durante o desenvolvimento | Média | 🔮 Considerando |
| **Validação de Parâmetros** | Sistema de validação de parâmetros para componentes | Média | 📅 Planejado |
| **Benchmarks** | Comparativo de performance com outras soluções | Média | 📅 Planejado |
| **Minificação CSS/JS** | Minificação automática de CSS e JS em produção | Média | 📅 Planejado |
| **Extensões para Ferramentas** | Plugins para IDEs e integrações com ferramentas de desenvolvimento | Baixa | 🔮 Considerando |
| **Server Side Rendering** | Implementação de SSR otimizado para SPAs | Baixa | 🔮 Considerando |
| **Design System Integrado** | Componentes base para facilitar a criação de interfaces consistentes | Baixa | 🔮 Considerando |

### Legenda
- 🔄 Em progresso: Desenvolvimento iniciado
- 📅 Planejado: Planejado para implementação em breve
- 🔮 Considerando: Sendo considerado para o futuro

## Licença
MIT







