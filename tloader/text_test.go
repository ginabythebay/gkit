package tloader_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/ginabythebay/gkit/tloader"
)

func ExampleTextLoader() {
	tmpDir, err := ioutil.TempDir("", "textexample")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"content.basetmpl",
			`Subject: testing\n{{template "content" .}}`,
		},
		file{
			"one.tmpl",
			`{{define "content"}}Hi {{.Name}}, this is email one.{{end}}`,
		},
		file{
			"two.tmpl",
			`{{define "content"}}Hi {{ toLower .Name}}, this is email two.{{end}}`,
		},
	}

	if err = writeAll(tmpDir, files); err != nil {
		panic(err)
	}

	ourTemplates := struct {
		One *template.Template
		Two *template.Template
	}{}

	funcMap := map[string]interface{}{
		"toLower": strings.ToLower,
	}

	err = tloader.TextLoader{
		BaseName:  filepath.Join(tmpDir, "content.basetmpl"),
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.
		FuncMap(funcMap).
		Load(&ourTemplates)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{
		"Bob",
	}

	if err = ourTemplates.One.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
	fmt.Println()
	if err = ourTemplates.Two.Execute(os.Stdout, data); err != nil {
		panic(err)
	}

	// Output:
	// Subject: testing\nHi Bob, this is email one.
	// Subject: testing\nHi bob, this is email two.
}

func TestTextBase(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "texttest")
	ok(t, err)
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"content.basetmpl",
			`Subject: testing\n{{template "content" .}}`,
		},
		file{
			"one.tmpl",
			`{{define "content"}}Hi {{.Name}}, this is email one.{{end}}`,
		},
		file{
			"two.tmpl",
			`{{define "content"}}Hi {{ toLower .Name}}, this is email two.{{end}}`,
		},
	}

	ok(t, writeAll(tmpDir, files))

	ourTemplates := struct {
		One *template.Template
		Two *template.Template
	}{}

	loader := tloader.TextLoader{
		BaseName:  filepath.Join(tmpDir, "content.basetmpl"),
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.FuncMap(funcMap)
	ok(t, loader.Load(&ourTemplates))

	data := struct {
		Name string
	}{
		"Bob",
	}

	equals(t, `Subject: testing\nHi Bob, this is email one.`, execute(t, ourTemplates.One, data))
	equals(t, `Subject: testing\nHi bob, this is email two.`, execute(t, ourTemplates.Two, data))
}

func TestTextNoBase(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "texttest")
	ok(t, err)
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"one.tmpl",
			`Hi {{.Name}}, this is email one.`,
		},
		file{
			"two.tmpl",
			`Hi {{ toLower .Name}}, this is email two.`,
		},
	}

	ok(t, writeAll(tmpDir, files))

	ourTemplates := struct {
		One *template.Template
		Two *template.Template
	}{}

	loader := tloader.TextLoader{
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.FuncMap(funcMap)
	ok(t, loader.Load(&ourTemplates))

	data := struct {
		Name string
	}{
		"Bob",
	}

	equals(t, "Hi Bob, this is email one.", execute(t, ourTemplates.One, data))
	equals(t, "Hi bob, this is email two.", execute(t, ourTemplates.Two, data))
}

// Similar to TestNoBase, but we have an additional field in
// ourTemplates, below that doesn't have a corresponding file, which
// causes an error during loading.
func TestTextMissingFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "texttest")
	ok(t, err)
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"one.tmpl",
			`Hi {{.Name}}, this is email one.`,
		},
		file{
			"two.tmpl",
			`Hi {{ toLower .Name}}, this is email two.`,
		},
	}

	ok(t, writeAll(tmpDir, files))

	ourTemplates := struct {
		One   *template.Template
		Two   *template.Template
		Three *template.Template
	}{}

	loader := tloader.TextLoader{
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.FuncMap(funcMap)
	// This should fail because there is no file corresponding to field Three
	err = loader.Load(&ourTemplates)
	assert(t, err != nil, "Failed to detect error, error should not have been nil")
	assert(t, strings.Contains(err.Error(), "three"), "Unexpected error %q.  Should have contained mention of %q", err.Error(), "three")
}
