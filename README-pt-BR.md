
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

Componentes integram-se perfeitamente com templates. Aqui está um componente de botão que aproveita a sintaxe de helpers:

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

E um componente de card que aninha outros componentes:

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
      <!-- Usando helper para aninhar componente de botão -->
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

### Usando um componente (Sintaxe com Helpers)

Skingo gera automaticamente funções helpers para todos os componentes registrados, oferecendo uma sintaxe limpa e intuitiva:

```html
<!-- templates/home.html -->
<template>
  <div class="container">
    <h1>{{.Title}}</h1>
    <p>{{.Content}}</p>

    <!-- Usando helpers de componentes com parâmetros nomeados -->
    {{ card (dict 
      "title" "Exemplo de Card" 
      "content" "Este é um exemplo de componente card com um botão." 
      "buttonText" "Ler mais"
    ) }}
    
    <!-- Usando parâmetros posicionais -->
    {{ button "Clique aqui!" "green" }}
  </div>
</template>
```

**Em vez de:** `{{ comp "card.html" (dict "title" "...") }}`  
**Você pode escrever:** `{{ card (dict "title" "...") }}`

Os helpers são gerados automaticamente com base em:
- Nome do componente (derivado do nome do arquivo ou nome registrado)
- Metadados do componente registrado
- A função `comp` continua disponível como alternativa

O Skingo vai, de forma inteligente, determinar os escopos de CSS e criar automaticamente classes que auxiliam na estilização de cada componente, respeitando os estilos específicos em primeiro lugar.

Se mais de um elemento sem pai (sem um contêiner) forem declarados entre as tags `<template><template>`, o Skingo criará de forma automática um cointêiner (`<div>`) para envolvê-los e assim separar inteligentemente os estilos entre os diversos componentes, respeitando cada escopo. 

Para evitar esse comportamento acima, basta adicionar o atributo `unwrap` na tag "template", dessa forma: `<template unwrap>`.

### Exemplo de Catálogo Híbrido

Skingo inclui um catálogo de interface de usuário inicial reutilizável no pacote `uikit` com componentes pré-construídos:
- `SkButton` - Botão estilizado com variantes (primary, outline, ghost)
- `SkInput` - Input de formulário com suporte a label
- `SkBadge` - Badge de status com variantes semânticas (success, warning, danger)
- `SkInfo` - Caixa de alerta/informação com variantes (info, success, error)
- `SkCard` - Card contêiner com header, conteúdo e ação opcional de rodapé

**Exemplo de uso:**
```go
import (
  "github.com/messiashenrique/skingo"
  "github.com/messiashenrique/skingo/uikit"
)

func main() {
  ts := skingo.NewTemplateSet("layout")
  
  // Registra o catálogo uikit
  if err := uikit.RegisterCatalog(ts); err != nil {
    log.Fatal(err)
  }
  
  // Habilita validação opcional
  ts.SetComponentValidation(skingo.ComponentValidationOptions{
    Enabled:     true,
    StrictTypes: true,
  })
  
  // Analisa e usa
  if err := ts.ParseDirs("templates"); err != nil {
    log.Fatal(err)
  }
  
  // Os helpers SkButton, SkInput, SkCard agora estão disponíveis em templates
}
```

**Em templates:**
```html
{{ SkButton "Clique aqui" "primary" }}
{{ SkInput (dict "name" "email" "label" "Email") }}
{{ SkCard (dict "title" "Meu Card" "content" "Conteúdo aqui") }}
{{ SkInfo (dict "title" "Informação" "message" "Olá!" "variant" "success") }}
{{ SkBadge "Ativo" "success" }}
```

Ver exemplo de integração híbrida em `examples/hybrid`.
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

## Testes de Componentes

Skingo inclui APIs integradas para testar metadados de componentes e renderização:

```go
func TestMetadadosComponente(t *testing.T) {
    ts := skingo.NewTemplateSet("layout")
    
    // Registra um catálogo de componentes
    if err := skingo.RegisterComponentCatalogJSON(ts, "meuscomponentes", []byte(`{
        "components": {
            "button": {
                "description": "Botão clicável",
                "variables": [
                    {"name": "label", "type": "string", "required": true}
                ]
            }
        }
    }`)); err != nil {
        t.Fatal(err)
    }
    
    // Testa recuperação de metadados
    meta, ok := ts.GetComponentMeta("button")
    if !ok {
        t.Fatal("Metadados do componente não encontrados")
    }
    
    if meta.Description != "Botão clicável" {
        t.Errorf("Esperado 'Botão clicável', obtive %s", meta.Description)
    }
}

func TestValidacaoComponente(t *testing.T) {
    ts := skingo.NewTemplateSet("layout")
    
    // Habilita validação
    ts.SetComponentValidation(skingo.ComponentValidationOptions{
        Enabled:     true,
        StrictTypes: true,
    })
    
    // A validação agora verifica parâmetros obrigatórios e tipos durante execução
}
```

Executar testes com ferramentas Go padrão:
```bash
go test ./...          # Executa todos os testes
go test -v .           # Executa testes com saída verbosa
go test -cover ./...   # Executa testes com relatório de cobertura
```

Ver `skingo_test.go` para exemplos abrangentes de testes:
- Registro e recuperação de metadados de componentes
- Análise de múltiplos filesystems com `ParseManyFS`
- Opções de validação e verificação de tipos
- Geração automática de funções helpers de componentes

## Roadmap para Desenvolvimento

| Etapa | Descrição | Prioridade | Status |
|-------|-----------|------------|--------|
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

### Metadados de Catálogo de Componentes
O Skingo agora possui APIs opcionais para registro de metadados de catálogos de componentes.
Essas APIs não alteram o comportamento de renderização e servem de base para documentação, tooling e validação futura.

```go
func (ts *TemplateSet) RegisterComponentMeta(name string, meta ComponentMeta) error
func (ts *TemplateSet) RegisterComponentCatalog(catalogName string, components map[string]ComponentMeta) error
func (ts *TemplateSet) RegisterComponentCatalogJSON(catalogName string, manifest []byte) error
func (ts *TemplateSet) RegisterComponentCatalogFile(catalogName string, filename string) error
func (ts *TemplateSet) RegisterComponentCatalogFS(catalogName string, filesystem fs.FS, manifestPath string) error
func (ts *TemplateSet) ListComponents() []ComponentInfo
func (ts *TemplateSet) GetComponentMeta(name string) (ComponentMeta, bool)
```

Exemplo de manifesto JSON:

```json
{
  "components": {
    "button": {
      "description": "Disparador de ação clicável",
      "version": "1.0.0",
      "variants": ["solid", "outline", "ghost"],
      "dependencies": ["icon"],
      "params": [
        {
          "name": "label",
          "type": "string",
          "required": true,
          "description": "Rótulo do botão"
        }
      ]
    }
  }
}
```

### Validação de Componentes (Opcional)
Você pode habilitar validação em runtime das chamadas de componente com base nos metadados registrados.

```go
type ComponentValidationOptions struct {
  Enabled     bool
  StrictTypes bool
}

func (ts *TemplateSet) SetComponentValidation(options ComponentValidationOptions)
func (ts *TemplateSet) EnableComponentValidation(enabled bool)
func (ts *TemplateSet) GetComponentValidation() ComponentValidationOptions
```

Comportamento da validação:
- `Enabled=false` (padrão): sem validação.
- `Enabled=true`: valida parâmetros obrigatórios.
- `StrictTypes=true` (padrão): valida tipos básicos declarados (`string`, `bool`, `int`, `float`, `number`, `[]string`, `[]map[string]string`, `map[string]interface{}`).
- Se o componente tiver parâmetro `variant` e os metadados tiverem `variants`, o valor é validado contra as variantes permitidas.

Exemplo:

```go
ts.SetComponentValidation(skingo.ComponentValidationOptions{
  Enabled:     true,
  StrictTypes: true,
})
```

### Exemplo de Catálogo Híbrido
O Skingo agora inclui um pacote inicial de catálogo UI reutilizável em `uikit` com componentes:
- `SkButton`
- `SkInput`
- `SkBadge`
- `SkInfo`
- `SkCard`

Veja o exemplo de integração híbrida em `examples/hybrid`.

```go
import (
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
| `addFloat` | Soma dois número do tipo Float | `{{addFloat 3.0 3.1}}` → `6.1` |
| `subFloat` | Subtrai dois número do tipo Float | `{{subFloat 7.3 3.1}}` → `4.2` |
| `mulFloat` | Multiplica dois número do tipo Float | `{{mulFloat 3.0 7.1}}` → `21.3` |
| `divFloat` | Divide dois número do tipo Float | `{{divFloat 24.6 3.0}}` → `8.2` |
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
| **Testes** | Implementação de testes unitários abrangentes | Alta | ✅ Completo |
| **Otimização de Performance** | Refatoração para melhorar a eficiência de renderização | Alta | 📅 Planejado |
| **Documentação Completa** | Documentação detalhada com exemplos para cada funcionalidade | Alta | 🔄 Em progresso |
| **Integração HTMX** | Suporte aprimorado para HTMX com helpers dedicados | Alta | 📅 Planejado |
| **Variantes Temáticas** | Suporte a variantes de componentes com light/dark/custom themes | Alta | 📅 Planejado |
| **Tokens de Design** | Sistema centralizado de tokens de design para componentes uikit | Alta | 📅 Planejado |
| **Exemplos Avançados** | Repositório com exemplos mais complexos e casos de uso reais | Média | 📅 Planejado |
| **Hot Reload** | Suporte para hot reload durante o desenvolvimento | Média | 🔮 Considerando |
| **Benchmarks** | Comparativo de performance com outras soluções | Média | 📅 Planejado |
| **Minificação CSS/JS** | Minificação automática de CSS e JS em produção | Média | 📅 Planejado |
| **Extensões para Ferramentas** | Plugins para IDEs e integrações com ferramentas de desenvolvimento | Baixa | 🔮 Considerando |
| **Server Side Rendering** | Implementação de SSR otimizado para SPAs | Baixa | 🔮 Considerando |

### Legenda
- 🔄 Em progresso: Desenvolvimento iniciado
- 📅 Planejado: Planejado para implementação em breve
- 🔮 Considerando: Sendo considerado para o futuro

## Licença
MIT







