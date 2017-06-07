package tloader

import (
	"fmt"
	"path/filepath"
	"text/template"
)

// TextLoader knows how to load a set of page templates, along with an
// optional base template.
type TextLoader struct {
	// Provides context for pages to be loaded in.  Optional.
	BaseName string
	// Value that can be passed to filePath.Glob to find all the page templates.
	PagesGlob string
	funcMap   map[string]interface{}
}

// FuncMap applies a function map to the loader and returns the loader
// (to enable function chaining).
func (t TextLoader) FuncMap(funcMap map[string]interface{}) TextLoader {
	t.funcMap = funcMap
	return t
}

// Load loads all pages found in the pages glob into the fields of v.
// We expect v to be a struct and all public fields to be of type
// *text/template.Template and will assign pages to each field by
// matching the lowercase version of the field name with the lowercase
// basename of the file.
func (t TextLoader) Load(v interface{}) error {
	pages, err := filepath.Glob(t.PagesGlob)
	if err != nil {
		return fmt.Errorf("globbing %q, %v", t.PagesGlob, err)
	}

	var root *template.Template
	if t.BaseName != "" {
		root = template.New(filepath.Base(t.BaseName))
		if t.funcMap != nil {
			root.Funcs(t.funcMap)
		}
		if _, err = root.ParseFiles(t.BaseName); err != nil {
			return fmt.Errorf("parsing %q: %v", t.BaseName, err)
		}
	}
	ldr := textPageLoader{root, t.funcMap}

	tByName := map[string]interface{}{}
	for _, p := range pages {
		t, err := ldr.loadPage(p)
		if err != nil {
			return fmt.Errorf("parsing %q: %v", p, err)
		}
		tByName[fixName(filepath.Base(p))] = t
	}

	return fillStruct(tByName, v)
}

type textPageLoader struct {
	root    *template.Template // can be nil
	funcMap map[string]interface{}
}

func (l textPageLoader) loadPage(page string) (*template.Template, error) {
	if l.root == nil {
		t := template.New("")
		if l.funcMap != nil {
			t.Funcs(l.funcMap)
		}
		if _, err := t.ParseFiles(page); err != nil {
			return nil, err
		}
		return t.Lookup(filepath.Base(page)), nil
	}
	t, err := l.root.Clone()
	if err != nil {
		return nil, err
	}
	return t.ParseFiles(page)
}
