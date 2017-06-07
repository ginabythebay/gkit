package tloader_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var funcMap = map[string]interface{}{
	"toLower": strings.ToLower,
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type file struct {
	name    string
	content string
}

func (f file) write(dir string) error {
	fn := filepath.Join(dir, f.name)
	return ioutil.WriteFile(fn, []byte(f.content), 0666)
}

func writeAll(dir string, files []file) error {
	for _, f := range files {
		if err := f.write(dir); err != nil {
			return err
		}
	}
	return nil
}

type executor interface {
	Execute(io.Writer, interface{}) error
}

func execute(t *testing.T, tpl executor, data interface{}) string {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}
	return buf.String()
}
