package tloader_test

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ginabythebay/gkit/tloader"
)

func TestHtmlBase(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "htmltest")
	ok(t, err)
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"content.basetmpl",
			`<html><head><title>Title</title></head><body>{{template "content" .}}</body></html>`,
		},
		file{
			"one.tmpl",
			`{{define "content"}}Hi {{.Name}}, this is page one.{{end}}`,
		},
		file{
			"two.tmpl",
			`{{define "content"}}Hi {{ toLower .Name}}, this is page two.{{end}}`,
		},
	}

	ok(t, writeAll(tmpDir, files))

	ourTemplates := struct {
		One *template.Template
		Two *template.Template
	}{}

	loader := tloader.HTMLLoader{
		BaseName:  filepath.Join(tmpDir, "content.basetmpl"),
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.FuncMap(funcMap)
	ok(t, loader.Load(&ourTemplates))

	data := struct {
		Name string
	}{
		"Bob",
	}

	equals(t, "<html><head><title>Title</title></head><body>Hi Bob, this is page one.</body></html>", execute(t, ourTemplates.One, data))
	equals(t, "<html><head><title>Title</title></head><body>Hi bob, this is page two.</body></html>", execute(t, ourTemplates.Two, data))
}

func TestHtmlNoBase(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "htmltest")
	ok(t, err)
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"one.tmpl",
			`<html>Hi {{.Name}}, this is page one.</html>`,
		},
		file{
			"two.tmpl",
			`<html>Hi {{ toLower .Name}}, this is page two.</html>`,
		},
	}

	ok(t, writeAll(tmpDir, files))

	ourTemplates := struct {
		One *template.Template
		Two *template.Template
	}{}

	loader := tloader.HTMLLoader{
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.FuncMap(funcMap)
	ok(t, loader.Load(&ourTemplates))

	data := struct {
		Name string
	}{
		"Bob",
	}

	equals(t, "<html>Hi Bob, this is page one.</html>", execute(t, ourTemplates.One, data))
	equals(t, "<html>Hi bob, this is page two.</html>", execute(t, ourTemplates.Two, data))
}

// Similar to TestNoBase, but we have an additional field in
// ourTemplates, below that doesn't have a corresponding file, which
// causes an error during loading.
func TestHtmlMissingFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "htmltest")
	ok(t, err)
	defer os.RemoveAll(tmpDir)

	files := []file{
		file{
			"one.tmpl",
			`<html>Hi {{.Name}}, this is page one.</html>`,
		},
		file{
			"two.tmpl",
			`<html>Hi {{ toLower .Name}}, this is page two.</html>`,
		},
	}

	ok(t, writeAll(tmpDir, files))

	ourTemplates := struct {
		One   *template.Template
		Two   *template.Template
		Three *template.Template
	}{}

	loader := tloader.HTMLLoader{
		PagesGlob: filepath.Join(tmpDir, "*.tmpl"),
	}.FuncMap(funcMap)
	// This should fail because there is no file corresponding to field Three
	err = loader.Load(&ourTemplates)
	assert(t, err != nil, "Failed to detect error, error should not have been nil")
	assert(t, strings.Contains(err.Error(), "three"), "Unexpected error %q.  Should have contained mention of %q", err.Error(), "three")
}
