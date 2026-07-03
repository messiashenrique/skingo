package skingo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
)

const testLayout = `<!DOCTYPE html>
<html>
<head><title>test</title></head>
<body>{{ .Yield }}</body>
</html>`

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}
	return path
}

func newTestFS(files map[string]string) fs.FS {
	testFS := fstest.MapFS{}
	for name, content := range files {
		testFS[name] = &fstest.MapFile{Data: []byte(content)}
	}
	return testFS
}

func TestExecuteIsolatedExtractsTemplateContent(t *testing.T) {
	dir := t.TempDir()
	fragment := writeTestFile(t, dir, "fragment.html", `<template unwrap><p>Hello {{ .Name }}</p></template>`)

	ts := NewTemplateSet("layout")
	var out strings.Builder
	if err := ts.ExecuteIsolated(&out, fragment, map[string]string{"Name": "Skingo"}); err != nil {
		t.Fatalf("ExecuteIsolated returned error: %v", err)
	}

	if got, want := out.String(), "<p>Hello Skingo</p>"; got != want {
		t.Fatalf("unexpected isolated output: got %q want %q", got, want)
	}
}

func TestExecuteIsolatedFSExtractsTemplateContent(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"templates/fragment.html": `<template unwrap><p>Hello {{ .Name }}</p></template>`,
	})

	ts := NewTemplateSet("layout")
	var out strings.Builder
	if err := ts.ExecuteIsolatedFS(&out, testFS, "templates/fragment.html", map[string]string{"Name": "Skingo"}); err != nil {
		t.Fatalf("ExecuteIsolatedFS returned error: %v", err)
	}

	if got, want := out.String(), "<p>Hello Skingo</p>"; got != want {
		t.Fatalf("unexpected isolated fs output: got %q want %q", got, want)
	}
}

func TestParseDirsRejectsDuplicateTemplateNames(t *testing.T) {
	dirA := t.TempDir()
	dirB := t.TempDir()
	writeTestFile(t, dirA, "layout.html", testLayout)
	writeTestFile(t, dirA, "card.html", `<template><div>A</div></template>`)
	writeTestFile(t, dirB, "card.html", `<template><div>B</div></template>`)

	ts := NewTemplateSet("layout")
	err := ts.ParseDirs(dirA, dirB)
	if err == nil {
		t.Fatal("expected duplicate template name error")
	}
	if !strings.Contains(err.Error(), `duplicate template name "card"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFSRejectsDuplicateTemplateNames(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"pages/layout.html":     testLayout,
		"pages/card.html":       `<template><div>A</div></template>`,
		"components/card.html":  `<template><div>B</div></template>`,
		"components/other.html": `<template><div>C</div></template>`,
	})

	ts := NewTemplateSet("layout")
	err := ts.ParseFS(testFS, "pages", "components")
	if err == nil {
		t.Fatal("expected duplicate template name error")
	}
	if !strings.Contains(err.Error(), `duplicate template name "card"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNestedComponentsAndParams(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": testLayout,
		"templates/page.html": `<template>
{{ comp "card" (dict "label" "Save") }}
</template>`,
		"templates/card.html": `<template>
<section>{{ comp "button" .label "green" }}</section>
</template>`,
		"templates/button.html": `<template>
<button class="{{ paramOr 1 "blue" }}">{{ param 0 }}</button>
</template>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		t.Fatalf("ParseFS returned error: %v", err)
	}

	html, err := ts.ExecuteString("page", nil)
	if err != nil {
		t.Fatalf("ExecuteString returned error: %v", err)
	}

	if !strings.Contains(html, `<button class="green">Save</button>`) {
		t.Fatalf("expected nested component output, got:\n%s", html)
	}
}

func TestComponentCSSAndJSIncludedOnce(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": testLayout,
		"templates/page.html": `<template>
{{ comp "button" "A" }}
{{ comp "button" "B" }}
</template>`,
		"templates/button.html": `<template>
<button class="btn">{{ param 0 }}</button>
</template>
<style>
.btn { color: red; }
</style>
<script>
console.log("button");
</script>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		t.Fatalf("ParseFS returned error: %v", err)
	}

	html, err := ts.ExecuteString("page", nil)
	if err != nil {
		t.Fatalf("ExecuteString returned error: %v", err)
	}

	if got := strings.Count(html, "color: red"); got != 1 {
		t.Fatalf("expected button CSS once, got %d occurrences in:\n%s", got, html)
	}
	if got := strings.Count(html, `console.log("button");`); got != 1 {
		t.Fatalf("expected button JS once, got %d occurrences in:\n%s", got, html)
	}
}

func TestExecuteWithLayoutUsesRequestedLayout(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": `<!DOCTYPE html>
<html>
<head><title>default</title></head>
<body><main class="default">{{ .Yield }}</main></body>
</html>`,
		"templates/admin.html": `<!DOCTYPE html>
<html>
<head><title>admin</title></head>
<body><aside>Admin</aside><main class="admin">{{ .Yield }}</main></body>
</html>`,
		"templates/page.html": `<template><h1>{{ .Title }}</h1></template>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		t.Fatalf("ParseFS returned error: %v", err)
	}

	defaultHTML, err := ts.ExecuteString("page", map[string]string{"Title": "Dashboard"})
	if err != nil {
		t.Fatalf("ExecuteString returned error: %v", err)
	}
	if !strings.Contains(defaultHTML, `class="default"`) || strings.Contains(defaultHTML, "<aside>Admin</aside>") {
		t.Fatalf("expected default layout output, got:\n%s", defaultHTML)
	}

	adminHTML, err := ts.ExecuteStringWithLayout("admin", "page", map[string]string{"Title": "Dashboard"})
	if err != nil {
		t.Fatalf("ExecuteStringWithLayout returned error: %v", err)
	}
	if !strings.Contains(adminHTML, `class="admin"`) || !strings.Contains(adminHTML, "<aside>Admin</aside>") {
		t.Fatalf("expected admin layout output, got:\n%s", adminHTML)
	}
}

func TestExecuteWithLayoutReturnsErrorForUnknownLayout(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": testLayout,
		"templates/page.html":   `<template><h1>{{ .Title }}</h1></template>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		t.Fatalf("ParseFS returned error: %v", err)
	}

	err := ts.ExecuteWithLayout(&strings.Builder{}, "missing", "page", nil)
	if err == nil {
		t.Fatal("expected unknown layout error")
	}
	if !strings.Contains(err.Error(), "layout template missing not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteConcurrent(t *testing.T) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": testLayout,
		"templates/page.html": `<template>
{{ comp "button" .Label .Color }}
</template>`,
		"templates/button.html": `<template>
<button class="{{ paramOr 1 "blue" }}">{{ param 0 }}</button>
</template>
<style>
.blue { color: blue; }
.green { color: green; }
</style>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		t.Fatalf("ParseFS returned error: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 250; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			label := fmt.Sprintf("Button %d", i)
			color := "blue"
			if i%2 == 0 {
				color = "green"
			}
			html, err := ts.ExecuteString("page", map[string]string{
				"Label": label,
				"Color": color,
			})
			if err != nil {
				t.Errorf("ExecuteString returned error: %v", err)
				return
			}
			if !strings.Contains(html, label) || !strings.Contains(html, " "+color+`">`) {
				t.Errorf("unexpected concurrent output for %q/%q:\n%s", label, color, html)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkExecuteSimple(b *testing.B) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": testLayout,
		"templates/page.html":   `<template><main><h1>{{ .Title }}</h1></main></template>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		b.Fatalf("ParseFS returned error: %v", err)
	}

	for i := 0; i < b.N; i++ {
		var out strings.Builder
		if err := ts.Execute(&out, "page", map[string]string{"Title": "Hello"}); err != nil {
			b.Fatalf("Execute returned error: %v", err)
		}
	}
}

func BenchmarkExecuteNestedComponents(b *testing.B) {
	testFS := newTestFS(map[string]string{
		"templates/layout.html": testLayout,
		"templates/page.html":   `<template>{{ comp "card" (dict "label" "Save") }}</template>`,
		"templates/card.html":   `<template><section>{{ comp "button" .label "green" }}</section></template>`,
		"templates/button.html": `<template><button class="{{ paramOr 1 "blue" }}">{{ param 0 }}</button></template>`,
	})

	ts := NewTemplateSet("layout")
	if err := ts.ParseFS(testFS, "templates"); err != nil {
		b.Fatalf("ParseFS returned error: %v", err)
	}

	for i := 0; i < b.N; i++ {
		var out strings.Builder
		if err := ts.Execute(&out, "page", nil); err != nil {
			b.Fatalf("Execute returned error: %v", err)
		}
	}
}

func BenchmarkExecuteIsolatedCacheHit(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "fragment.html")
	if err := os.WriteFile(path, []byte(`<template><p>{{ .Text }}</p></template>`), 0o644); err != nil {
		b.Fatalf("writing test file: %v", err)
	}

	ts := NewTemplateSet("layout")
	var warmup strings.Builder
	if err := ts.ExecuteIsolated(&warmup, path, map[string]string{"Text": "Hello"}); err != nil {
		b.Fatalf("warmup ExecuteIsolated returned error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out strings.Builder
		if err := ts.ExecuteIsolated(&out, path, map[string]string{"Text": "Hello"}); err != nil {
			b.Fatalf("ExecuteIsolated returned error: %v", err)
		}
	}
}

func BenchmarkExecuteIsolatedCacheMiss(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "fragment.html")
	if err := os.WriteFile(path, []byte(`<template><p>{{ .Text }}</p></template>`), 0o644); err != nil {
		b.Fatalf("writing test file: %v", err)
	}

	ts := NewTemplateSet("layout")
	for i := 0; i < b.N; i++ {
		ts.ClearIsolatedCache()
		var out strings.Builder
		if err := ts.ExecuteIsolated(&out, path, map[string]string{"Text": "Hello"}); err != nil {
			b.Fatalf("ExecuteIsolated returned error: %v", err)
		}
	}
}

func BenchmarkParseFSManyTemplates(b *testing.B) {
	files := map[string]string{
		"templates/layout.html": testLayout,
	}
	for i := 0; i < 100; i++ {
		files[fmt.Sprintf("templates/page-%03d.html", i)] = fmt.Sprintf(`<template><p>Page %d</p></template>`, i)
	}
	testFS := newTestFS(files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts := NewTemplateSet("layout")
		if err := ts.ParseFS(testFS, "templates"); err != nil {
			b.Fatalf("ParseFS returned error: %v", err)
		}
	}
}
