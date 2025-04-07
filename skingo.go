package skingo

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
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

// TemplateSet represents a set of templates
type TemplateSet struct {
	templates     map[string]*Template
	layout        *Layout
	layoutName    string
	masterTmpl    *template.Template
	templateHTML  map[string]string
	mu            sync.Mutex
	usedTemplates map[string]bool               // Track which templates have been used
	customFuncs   template.FuncMap              // Stores custom functions
	isolatedCache map[string]*template.Template // Cache of isolated templates
	cacheMu       sync.RWMutex                  // Specific mutex for cache
}

const (
	uniqueOpenToken      = "___GO_TEMPLATE_OPEN___"
	uniqueCloseToken     = "___GO_TEMPLATE_CLOSE___"
	ElementTypeNormal    = 0 // Normal element
	ElementTypeSingle    = 1 // Single element
	ElementTypeContainer = 2 // Root Container
)

var (
	htmlRegex     = regexp.MustCompile(`(?s)<template([^>]*)>(.*?)</template>`)
	cssRegex      = regexp.MustCompile(`(?s)<style([^>]*)>(.*?)</style>`)
	jsRegex       = regexp.MustCompile(`(?s)<script>(.*?)</script>`)
	classRegex    = regexp.MustCompile(`class\s*=\s*["']([^"']*)["']`)
	openTagRegex  = regexp.MustCompile(`^\s*<[^>]+>`)
	unwrapRegex   = regexp.MustCompile(`unwrap`)
	firstTagRegex = regexp.MustCompile(`^\s*<([a-zA-Z][a-zA-Z0-9]*)([^>]*)>`)
)

// defaultFuncs contains the default functions available in all templates
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
}

// NewTemplateSet creates a new template set using the specified template
// as the layout. The layout must contain <head> and <body> tags
// where the CSS and JS will be automatically injected.
// In addition, it is mandatory to inform the entry point of the templates that will be
// rendered in the layout, defining the '{{ .Yield }}' variable.
func NewTemplateSet(layoutName string) *TemplateSet {
	ts := &TemplateSet{
		templates:     make(map[string]*Template),
		layout:        nil,
		layoutName:    layoutName,
		masterTmpl:    template.New("master"),
		templateHTML:  make(map[string]string),
		usedTemplates: make(map[string]bool),
		customFuncs:   make(template.FuncMap),
		isolatedCache: make(map[string]*template.Template),
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

// parseFile analyze a file and extract HTML, CSS and JS
func (ts *TemplateSet) parseFile(filename string) error {

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
	if matches := htmlRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		templateAttrs := matches[1]
		templateContent := matches[2]
		trimmedContent := strings.TrimSpace(templateContent)

		// Verify if has unwrap attribute
		unwrap := unwrapRegex.MatchString(templateAttrs)

		t.HTML = trimmedContent

		// First, temporarily replace the {{ }} delimiters so as not to interfere with parsing
		safeContent := strings.ReplaceAll(trimmedContent, "{{", uniqueOpenToken)
		safeContent = strings.ReplaceAll(safeContent, "}}", uniqueCloseToken)

		// Verify if it starts with an opening tag and find which it is
		hasRootElement := false
		isSingleElement := false
		isRootContainer := false // Flag for identifying a root container
		rootTagName := ""
		rootClasses := []string{} // Store the classes of the root element

		// Regex for finding the first opening tag
		if firstTagMatch := firstTagRegex.FindStringSubmatch(safeContent); len(firstTagMatch) > 2 {
			tagName := firstTagMatch[1]
			rootTagName = tagName

			// Extract classes from root element
			rootAttributes := firstTagMatch[2]
			if classMatches := classRegex.FindStringSubmatch(rootAttributes); len(classMatches) > 1 {
				classStr := classMatches[1]
				// Split classes by space and append
				rootClasses = append(rootClasses, strings.Fields(classStr)...)
			}

			// Verify if it ends with the corresponding closing tag
			closeTagPattern := fmt.Sprintf(`</\s*%s\s*>\s*$`, regexp.QuoteMeta(tagName))
			closeTagRegex := regexp.MustCompile(closeTagPattern)

			if closeTagRegex.MatchString(safeContent) {
				hasRootElement = true

				// Verify if it's a single element (without other elements between the tags)
				innerContent := safeContent
				innerContent = openTagRegex.ReplaceAllString(innerContent, "")
				closeTagRegex := regexp.MustCompile(`</\s*[^>]+>\s*$`)
				innerContent = closeTagRegex.ReplaceAllString(innerContent, "")

				if !strings.Contains(innerContent, "<") {
					isSingleElement = true
				} else {
					isRootContainer = true
				}
			}
		}

		// Extract the CSS
		var css string
		if cssMatches := cssRegex.FindStringSubmatch(string(content)); len(cssMatches) > 2 {
			css = cssMatches[2]
		}

		// If there is no CSS, we don't need to do anything with the scope
		if css == "" {
		} else if unwrap || hasRootElement {
			if hasRootElement {
				// Verify if there is a class attribute, adding our class in various possible situations
				if strings.Contains(t.HTML, "class=\"") {
					t.HTML = strings.Replace(t.HTML, "class=\"", fmt.Sprintf("class=\"%s ", t.scopeClass), 1)
				} else if strings.Contains(t.HTML, "class='") {
					t.HTML = strings.Replace(t.HTML, "class='", fmt.Sprintf("class='%s ", t.scopeClass), 1)
				} else if strings.Contains(t.HTML, "class={{") {
					t.HTML = strings.Replace(t.HTML, "class={{", fmt.Sprintf("class=\"%s {{", t.scopeClass), 1)
				} else {
					// Without class attribute, we need to add before the >
					lastPos := -1
					depth := 0
					for i, char := range t.HTML {
						if char == '{' {
							depth++
						} else if char == '}' {
							depth--
						} else if char == '>' && depth == 0 {
							lastPos = i
							break
						}
					}

					if lastPos != -1 {
						t.HTML = t.HTML[:lastPos] + fmt.Sprintf(" class=\"%s\"", t.scopeClass) + t.HTML[lastPos:]
					}
				}

				// Process CSS according to element type
				var elementType int
				if isSingleElement || unwrap {
					elementType = ElementTypeSingle
				} else if isRootContainer {
					elementType = ElementTypeContainer
				} else {
					elementType = ElementTypeNormal
				}

				t.CSS = scopedCSS(css, t.scopeClass, rootTagName, rootClasses, elementType)
			} else {
				// Whithout root element, but with unwrap, we use a custom selector instead of class
				t.HTML = fmt.Sprintf(`<div class="%s" style="display:contents">%s</div>`, t.scopeClass, t.HTML)
				t.CSS = containedScopedCSS(css, t.scopeClass)
			}
		} else {
			// Default case: wrap with div
			t.HTML = fmt.Sprintf(`<div class="%s">%s</div>`, t.scopeClass, t.HTML)
			t.CSS = containedScopedCSS(css, t.scopeClass)
		}
	}

	// Extract the JS from tags script
	if matches := jsRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		t.JS = matches[1]
	}

	// Stores the template for later processing
	ts.templates[t.Name] = t
	ts.templateHTML[t.Name] = t.HTML

	return nil
}

// scopedCSS creates CSS scope for elements inside a container
// (for example, when elements are inside a div with the scope class)
func scopedCSS(css string, scopeClass string, rootElementTag string, rootClasses []string, elementType int) string {
	cssBlocks := strings.Split(css, "}")
	var scopedCSS strings.Builder

	for _, block := range cssBlocks {
		if strings.TrimSpace(block) == "" {
			continue
		}

		parts := strings.SplitN(block, "{", 2)
		if len(parts) != 2 {
			continue
		}

		selectors := parts[0]
		declarations := parts[1]

		// Split multiple selectors (separated by commas)
		selectorList := strings.Split(selectors, ",")
		var scopedSelectors []string

		for _, selector := range selectorList {
			selector = strings.TrimSpace(selector)
			if selector == "" {
				continue
			}

			if selector == rootElementTag {
				// Is it the root element, add the class directly
				scopedSelectors = append(scopedSelectors, fmt.Sprintf("%s.%s", selector, scopeClass))
			} else if strings.HasPrefix(selector, ".") {
				// Extract the class name without the dot
				className := selector[1:]

				// Verify if it's a single element or the class is in the root element
				useDirectScope := false

				if elementType == ElementTypeSingle {
					// For single elements, all classes are treated without space
					useDirectScope = true
				} else {
					// For other types, check if the class is in the root element
					for _, rootClass := range rootClasses {
						if rootClass == className {
							useDirectScope = true
							break
						}
					}
				}

				if useDirectScope {
					// Without espace: ".class" -> ".s-xxxxx.class"
					scopedSelectors = append(scopedSelectors, fmt.Sprintf(".%s%s", scopeClass, selector))
				} else {
					// With espace: ".class" -> ".s-xxxxx .class"
					scopedSelectors = append(scopedSelectors, fmt.Sprintf(".%s %s", scopeClass, selector))
				}
			} else if strings.HasPrefix(selector, ":") {
				// Is a pseudo-class
				if rootElementTag != "" {
					scopedSelectors = append(scopedSelectors, fmt.Sprintf("%s.%s%s", rootElementTag, scopeClass, selector))
				} else {
					scopedSelectors = append(scopedSelectors, fmt.Sprintf(".%s%s", scopeClass, selector))
				}
			} else if strings.Contains(selector, " ") || strings.Contains(selector, ">") ||
				strings.Contains(selector, "+") || strings.Contains(selector, "~") {
				// Is a selector with children or siblings
				scopedSelectors = append(scopedSelectors, fmt.Sprintf(".%s %s", scopeClass, selector))
			} else {
				// Is other element
				scopedSelectors = append(scopedSelectors, fmt.Sprintf(".%s %s", scopeClass, selector))
			}
		}

		scopedSelector := strings.Join(scopedSelectors, ", ")
		scopedCSS.WriteString(scopedSelector)
		scopedCSS.WriteString(" {")
		scopedCSS.WriteString(declarations)
		scopedCSS.WriteString("}\n")
	}

	return scopedCSS.String()
}

// containedScopedCSS creates CSS scope for elements inside a container
// (for example, when elements are inside a div with the scope class)
func containedScopedCSS(css string, scopeClass string) string {
	// Split CSS into blocks of rules
	cssBlocks := strings.Split(css, "}")
	var scopedCSS strings.Builder

	for _, block := range cssBlocks {
		if strings.TrimSpace(block) == "" {
			continue
		}

		// Split into selector and declarations
		parts := strings.SplitN(block, "{", 2)
		if len(parts) != 2 {
			continue
		}

		selectors := parts[0]
		declarations := parts[1]

		// Split multiple selectors (separated by commas)
		selectorList := strings.Split(selectors, ",")
		var scopedSelectors []string

		for _, selector := range selectorList {
			selector = strings.TrimSpace(selector)
			if selector == "" {
				continue
			}

			// For any type of selector, we use the scope class as the ancestor
			// This works for elements (h1, p, a) and for classes (.btn, .blue)
			scopedSelectors = append(scopedSelectors, fmt.Sprintf(".%s %s", scopeClass, selector))
		}

		// Merge the transformed selectors
		scopedSelector := strings.Join(scopedSelectors, ", ")
		scopedCSS.WriteString(scopedSelector)
		scopedCSS.WriteString(" {")
		scopedCSS.WriteString(declarations)
		scopedCSS.WriteString("}\n")
	}

	return scopedCSS.String()
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
// The method processes files with the .html or .tmpl extension, extracting
// components that contain <template>, <style>, and <script> tags.
//
// For each file, the content inside the <template> tag is extracted as HTML.
// The content inside the <style> tag is extracted as CSS and automatically
// scoped using unique classes to avoid conflicts.
// The content inside the <script> tag is extracted as JavaScript.
//
// The method requires that a layout template (defined when creating the
// TemplateSet) be found in at least one of the given directories.
//
// After processing, the templates are available for rendering via
// the Execute method, with their CSS styles and JS scripts automatically
// included in the appropriate places in the layout.
//
// Returns an error if any directory cannot be read, if any template
// cannot be parsed, or if the layout template is not found.
// ParseDirs parses all HTML/template files in the given directories.
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
				if err := ts.parseFile(filepath.Join(dir, file.Name())); err != nil {
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

	// Add internal functions to layout - especialmente 'comp'
	for name, fn := range internalFuncs {
		// Adicionar apenas funções úteis para o layout
		if name == "comp" || name == "dict" || name == "param" || name == "paramOr" {
			layoutFuncs[name] = fn
		}
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

// ParseDir invokes ParseDirs, but with a unique directory.
func (ts *TemplateSet) ParseDir(dir string) error {
	return ts.ParseDirs(dir)
}

// Execute renders a specific template using the configured layout.
// The method combines the HTML content of the requested template with the
// layout, automatically injecting all CSS and JavaScript associated with the
// templates used (including referenced components) into the appropriate places
// in the layout.
//
// The 'name' parameter must match the name of a previously parsed template
// (without extension).
//
// The 'data' parameter contains the data that will be passed to the template.
// It can be any type supported by the html/template package, such as map,
// struct, or nil if no data is required.
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

	// Pre-parse the layout to find all component calls
	layoutContent := ts.layout.HTML
	compRegex := regexp.MustCompile(`{{[^}]*comp\s+"?([^"\s}]+)"?`)
	matches := compRegex.FindAllStringSubmatch(layoutContent, -1)

	ts.mu.Lock()
	for _, match := range matches {
		if len(match) > 1 {
			compName := match[1]
			compName = strings.TrimSuffix(compName, ".html")
			ts.usedTemplates[compName] = true
		}
	}
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
		"Yield": template.HTML(contentBuf.String()),
		"CSS":   template.CSS(allCSS.String()),
		"JS":    template.JS(allJS.String()),
		"Data":  data,
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

	ts.cacheMu.RLock()
	cachedTmpl, exists := ts.isolatedCache[filename]
	ts.cacheMu.RUnlock()

	if exists {
		return cachedTmpl.Execute(w, data) // Use the cached template
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	name := filepath.Base(filename)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	var htmlContent string
	if matches := htmlRegex.FindStringSubmatch(string(content)); len(matches) > 1 {
		htmlContent = matches[1]
	} else {
		htmlContent = string(content)
	}

	isolatedTmpl := template.New(name + "_isolated")
	isolatedTmpl.Funcs(defaultFuncs)   // Add default functions
	isolatedTmpl.Funcs(ts.customFuncs) // Add custom functions

	parsedTmpl, err := isolatedTmpl.Parse(htmlContent)
	if err != nil {
		return fmt.Errorf("error parsing isolated template: %w", err)
	}

	// Add to cache
	ts.cacheMu.Lock()
	ts.isolatedCache[filename] = parsedTmpl
	ts.cacheMu.Unlock()

	// Execute the isolated template with data
	return parsedTmpl.Execute(w, data)
}
