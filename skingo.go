package skingo

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// Template represents a template with separate HTML, CSS and JS.
// Each template can contain HTML, CSS and JavaScript code that is
// automatically unified and scoped.
type Template struct {
	Name       string
	HTML       string
	CSS        string
	JS         string
	tmpl       *template.Template
	scopeClass string
}

// Layout represents a template for a layout
type Layout struct {
	HTML string
	tmpl *template.Template
}

// TemplateSet é um conjunto de templates
type TemplateSet struct {
	templates     map[string]*Template
	layout        *Layout
	layoutName    string
	masterTmpl    *template.Template
	templateHTML  map[string]string
	mu            sync.Mutex
	usedTemplates map[string]bool  // Track which templates have been used
	customFuncs   template.FuncMap // Stores custom functions
}

// defaultFuncs contém as funções padrão disponíveis em todos os templates
var defaultFuncs = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"mod": func(a, b int) int { return a % b },
	"mul": func(a, b int) int { return a * b },
	"sub": func(a, b int) int { return a - b },
	"toJson": func(v interface{}) string {
		b, err := json.Marshal(v)
		if err != nil {
			return "{}"
		}
		return string(b)
	},
	"hasField": func(v interface{}, field string) bool {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return false
		}
		f := val.FieldByName(field)
		return f.IsValid() && !f.IsZero()
	},
}

// NewTemplateSet creates a new template set using the specified template
// as the layout. The layout must contain <head> and <body> tags
// where the CSS and JS will be automatically injected.
func NewTemplateSet(layoutName string) *TemplateSet {
	ts := &TemplateSet{
		templates:     make(map[string]*Template),
		layout:        nil,
		layoutName:    layoutName,
		masterTmpl:    template.New("master"),
		templateHTML:  make(map[string]string),
		usedTemplates: make(map[string]bool),
		customFuncs:   make(template.FuncMap),
	}

	// Apply default functions immediately
	ts.masterTmpl.Funcs(defaultFuncs)

	return ts
}

// AddFuncs adds custom functions to the template set.
// These functions will be available in all templates.
// Note: This method should be called before ParseDirs.
func (ts *TemplateSet) AddFuncs(funcMap template.FuncMap) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Save the custom functions for later use
	for name, fn := range funcMap {
		ts.customFuncs[name] = fn
	}

	// Apply them to the master template
	ts.masterTmpl.Funcs(funcMap)
}

// generateScopeClass build a scope class based on the template name and returns
func generateScopeClass(name string) string {
	// build a hash basead in template name
	hash := md5.Sum([]byte(name))
	// Return the first six characters of the hash
	return fmt.Sprintf("s-%x", hash)[:8]
}

// ParseFile analyze a file and extract HTML, CSS and JS
func (ts *TemplateSet) ParseFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	name := filepath.Base(filename)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	// Processes layout in a special way
	if name == ts.layoutName {
		return ts.parseLayoutFile(string(content))
	}

	t := &Template{
		Name:       name,
		scopeClass: generateScopeClass(name),
	}

	// Extract the HTML, CSS and JS from template tags
	htmlRegex := regexp.MustCompile(`(?s)<template>(.*?)</template>`)
	if matches := htmlRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		t.HTML = matches[1]
	}

	// Extract the CSS from tags style
	cssRegex := regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	if matches := cssRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		css := matches[1]
		t.HTML = fmt.Sprintf(`<div class="%s">%s</div>`, t.scopeClass, t.HTML)
		t.CSS = simpleScopedCSS(css, t.scopeClass)
	} else {
		t.HTML = fmt.Sprintf(`<div class="%s">%s</div>`, t.scopeClass, t.HTML)
	}

	// Extract the JS from tags script
	jsRegex := regexp.MustCompile(`(?s)<script>(.*?)</script>`)
	if matches := jsRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		t.JS = matches[1]
	}

	// Stores the template for later processing
	ts.templates[t.Name] = t
	ts.templateHTML[t.Name] = t.HTML

	return nil
}

// simpleScopedCSS prefix the scope class to all CSS selectors
func simpleScopedCSS(css string, scopeClass string) string {
	// We divide CSS into blocks of rules
	cssBlocks := strings.Split(css, "}")
	result := make([]string, 0, len(cssBlocks))

	for _, block := range cssBlocks {
		if strings.TrimSpace(block) == "" {
			continue
		}

		parts := strings.SplitN(block, "{", 2)
		if len(parts) < 2 {
			// Not a valid CSS block
			result = append(result, block)
			continue
		}

		selectors := strings.Split(parts[0], ",")
		transformedSelectors := make([]string, 0, len(selectors))

		for _, sel := range selectors {
			sel = strings.TrimSpace(sel)
			transformedSelectors = append(transformedSelectors, fmt.Sprintf(".%s %s", scopeClass, sel))
		}

		result = append(result, fmt.Sprintf("%s {%s", strings.Join(transformedSelectors, ", "), parts[1]))
	}

	return strings.Join(result, "}\n") + "}"
}

// parseLayoutFile processes a layout template file
func (ts *TemplateSet) parseLayoutFile(content string) error {
	layout := &Layout{
		HTML: content,
	}

	// Insert the style tag for the template before the </head>
	headCloseIndex := strings.Index(layout.HTML, "</head>")
	if headCloseIndex == -1 {
		return fmt.Errorf("layout template must contain </head> tag")
	}

	layout.HTML = layout.HTML[:headCloseIndex] +
		"\n\t<style>{{ .CSS }}</style>\n" +
		layout.HTML[headCloseIndex:]

	// Insert the script tag for the template before the </body>
	bodyCloseIndex := strings.Index(layout.HTML, "</body>")
	if bodyCloseIndex == -1 {
		return fmt.Errorf("layout template must contain </body> tag")
	}

	layout.HTML = layout.HTML[:bodyCloseIndex] +
		"\n\t<script>{{ .JS }}</script>\n" +
		layout.HTML[bodyCloseIndex:]

	ts.layout = layout

	return nil
}

// ParseDirs parses all HTML/template files in the given directories.
// The method processes files with the .html or .tmpl extension, extracting components
// that contain <template>, <style>, and <script> tags.
//
// For each file, the content inside the <template> tag is extracted as HTML.
// The content inside the <style> tag is extracted as CSS and automatically scoped using unique classes to avoid conflicts.
// The content inside the <script> tag is extracted as JavaScript.
//
// The method requires that a layout template (defined when creating the TemplateSet)
// be found in at least one of the given directories.
//
// After processing, the templates are available for rendering via
// the Execute method, with their CSS styles and JS scripts automatically included
// in the appropriate places in the layout.
//
// Returns an error if any directory cannot be read, if any template
// cannot be parsed, or if the layout template is not found.
func (ts *TemplateSet) ParseDirs(dirs ...string) error {
	layoutFound := false

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", dir, err)
		}

		// First pass: read all files and extract HTML, CSS and JS
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) == ".html" || filepath.Ext(file.Name()) == ".tmpl" {
				name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				if name == ts.layoutName {
					layoutFound = true
				}
				if err := ts.ParseFile(filepath.Join(dir, file.Name())); err != nil {
					return fmt.Errorf("error parsing file %s: %w", file.Name(), err)
				}
			}
		}
	}

	if !layoutFound {
		return fmt.Errorf("layout template '%s' not found in any of the provided directories", ts.layoutName)
	}

	type compCall struct {
		Args []interface{}
		Name string
	}

	// Component call stack for handling nested components
	var compStack []compCall

	var compMu sync.Mutex

	// Globals functions for all templates
	internalFuncs := template.FuncMap{
		"_register_template": func(name string) string {
			ts.mu.Lock()
			defer ts.mu.Unlock()
			ts.usedTemplates[name] = true
			return ""
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict needs key and value pairs as arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"param": func(index int) interface{} {
			compMu.Lock()
			defer compMu.Unlock()

			if len(compStack) == 0 {
				return nil
			}

			current := compStack[len(compStack)-1]
			if index < 0 || index >= len(current.Args) {
				return nil
			}
			return current.Args[index]
		},
		"paramOr": func(index int, defaultValue interface{}) interface{} {
			compMu.Lock()
			defer compMu.Unlock()

			if len(compStack) == 0 {
				return defaultValue
			}

			current := compStack[len(compStack)-1]
			if index < 0 || index >= len(current.Args) {
				return defaultValue
			}
			if current.Args[index] == nil {
				return defaultValue
			}
			return current.Args[index]
		},
		"comp": func(templateName string, args ...interface{}) (template.HTML, error) {
			name := strings.TrimSuffix(templateName, ".html")

			ts.mu.Lock()
			ts.usedTemplates[name] = true
			ts.mu.Unlock()

			compMu.Lock()
			compStack = append(compStack, compCall{
				Args: args,
				Name: name,
			})
			compMu.Unlock()

			// Ensures stack removal when finished
			defer func() {
				compMu.Lock()
				if len(compStack) > 0 {
					compStack = compStack[:len(compStack)-1]
				}
				compMu.Unlock()
			}()

			var buf strings.Builder

			var data interface{}

			if len(args) == 1 {
				if mapData, ok := args[0].(map[string]interface{}); ok {
					data = mapData
				} else {
					// Caso simples: um único valor que não é um mapa
					// Mapeamos para índice numérico
					data = map[string]interface{}{
						"0": args[0],
					}
				}
			} else {
				dataMap := make(map[string]interface{})
				for i, arg := range args {
					dataMap[fmt.Sprintf("%d", i)] = arg
				}
				data = dataMap
			}

			tmplName := name
			if !strings.HasSuffix(tmplName, ".html") {
				tmplName = tmplName + ".html"
			}

			if err := ts.masterTmpl.ExecuteTemplate(&buf, tmplName, data); err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
	}

	// Add internal functions
	ts.masterTmpl.Funcs(internalFuncs)

	// Second pass: create the templates and allow references between them
	for name, html := range ts.templateHTML {
		templateName := name
		if !strings.HasSuffix(templateName, ".html") {
			templateName = name + ".html"
		}

		// We modified the HTML to register the template when it is executed
		registeredHTML := "{{_register_template \"" + name + "\"}}" + html

		_, err := ts.masterTmpl.New(templateName).Parse(registeredHTML)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", name, err)
		}

		ts.templates[name].tmpl = ts.masterTmpl.Lookup(templateName)
	}

	// Prepare the layout template with all functions
	layoutFuncs := template.FuncMap{}

	// Combine default functions
	for name, fn := range defaultFuncs {
		layoutFuncs[name] = fn
	}

	// Add custom functions
	for name, fn := range ts.customFuncs {
		layoutFuncs[name] = fn
	}

	layoutTmpl := template.New(ts.layoutName)
	layoutTmpl.Funcs(layoutFuncs)

	layoutTmpl, err := layoutTmpl.Parse(ts.layout.HTML)
	if err != nil {
		return err
	}
	ts.layout.tmpl = layoutTmpl

	return nil
}

// ParseDir invokes ParseDirs, but with a unique directory. See more in ParseDirs
func (ts *TemplateSet) ParseDir(dir string) error {
	return ts.ParseDirs(dir)
}

// Execute renders a specific template using the configured layout.
// The method combines the HTML content of the requested template with the layout,
// automatically injecting all CSS and JavaScript associated with the templates
// used (including referenced components) into the appropriate places in the layout.
//
// The 'name' parameter must match the name of a previously parsed template
// (without extension).
//
// The 'data' parameter contains the data that will be passed to the template.
// It can be any type supported by the html/template package, such as map, struct,
// or nil if no data is required.
//
// During rendering, the system automatically tracks which templates
// are used (including components referenced via the comp function) and includes
// only the CSS and JavaScript of the templates actually used.
//
// The resulting HTML is written in Writer 'w'.
//
// Returns an error if the requested template does not exist, if the layout is
// not defined, or if an error occurs during template execution.
func (ts *TemplateSet) Execute(w io.Writer, name string, data interface{}) error {
	_, ok := ts.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	if ts.layout == nil {
		return fmt.Errorf("layout template not defined")
	}

	// Clean the usedTemplates list.
	ts.mu.Lock()
	ts.usedTemplates = make(map[string]bool)
	ts.mu.Unlock()

	// Creates a buffer to capture the template output
	var contentBuf strings.Builder

	// Use masterTmpl to execute the template
	err := ts.masterTmpl.ExecuteTemplate(&contentBuf, name+".html", data)
	if err != nil {
		return err
	}

	var allCSS strings.Builder
	var allJS strings.Builder

	ts.mu.Lock()
	for templateName := range ts.usedTemplates {
		if template, ok := ts.templates[templateName]; ok {
			if template.CSS != "" {
				allCSS.WriteString(template.CSS)
				allCSS.WriteString("\n")
			}
			if template.JS != "" {
				allJS.WriteString(template.JS)
				allJS.WriteString("\n")
			}
		}
	}
	ts.mu.Unlock()

	// Prepare the data for layout
	layoutData := map[string]interface{}{
		"Content": template.HTML(contentBuf.String()),
		"CSS":     template.CSS(allCSS.String()),
		"JS":      template.JS(allJS.String()),
		"Data":    data,
	}

	// Execute the layout template with the prepared data
	return ts.layout.tmpl.Execute(w, layoutData)
}

// ExecuteIsolated renders a template directly, without using the configured layout.
// This method is ideal for use with 'HTMX', Ajax requests, or any scenario
// where only an HTML fragment is needed, without the full page structure.
//
// The 'filename' parameter must be the full path to a template file.
// Unlike the Execute method, the file is read and parsed on demand, and does not
// need to have been previously processed by ParseDirs.
//
// If the file contains <template> tags, only the content inside those
// tags will be processed. Otherwise, the entire file is treated as a template.
//
// The 'data' parameter contains the data to be passed to the template.
//
// Note that, unlike the Execute method, CSS and JavaScript are not included
// in the result, only the raw HTML is rendered. There is also no tracking
// of components used.
//
// Returns an error if the file cannot be read or if an error occurs during
// template execution.
func (ts *TemplateSet) ExecuteIsolated(w io.Writer, filename string, data interface{}) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	name := filepath.Base(filename)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	var htmlContent string
	htmlRegex := regexp.MustCompile(`(?s)<template>(.*?)</template>`)
	if matches := htmlRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		htmlContent = matches[1]
	} else {
		htmlContent = string(content)
	}

	isolatedTmpl := template.New(name + "_isolated")

	// Adicionar funções padrão
	isolatedTmpl.Funcs(defaultFuncs)

	// Adicionar funções customizadas
	isolatedTmpl.Funcs(ts.customFuncs)

	parsedTmpl, err := isolatedTmpl.Parse(htmlContent)
	if err != nil {
		return fmt.Errorf("error parsing isolated template: %w", err)
	}

	// Execute the isolated template with data
	return parsedTmpl.Execute(w, data)
}
