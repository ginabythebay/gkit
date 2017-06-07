package tloader

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

func fixName(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name)) // trim extension
	name = strings.ToLower(name)
	return name
}

func fillStruct(namedValues map[string]interface{},
	container interface{}) error {
	s := reflect.ValueOf(container).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		name := strings.ToLower(typeOfT.Field(i).Name)
		t, ok := namedValues[name]
		if !ok {
			return fmt.Errorf("unable to find expected template %q", name)
		}
		f.Set(reflect.ValueOf(t))
	}
	return nil
}
