package main

import (
	"github.com/totallygamerjet/jb4go/parser"
	"github.com/totallygamerjet/jb4go/transformer"
	"github.com/traefik/yaegi/interp"
	"io"
	"os"
	"strings"
	"testing"
	"text/template"
)

func Test_IntsConst(t *testing.T) {
	// IntConst.class tests to make sure all the constant opcodes work properly
	src := transpile("../examples/IntConst.class")
	v := eval(src, "main.ints_IntConst_Call__I")
	f, ok := v.(func() int32)
	if !ok {
		t.Failed()
	}
	g := func() int32 {
		return -1 + 0 + 1 + 2 + 3 + 4 + 5
	}
	if f() != g() {
		t.Failed()
	}
}

func Test_Ints(t *testing.T) {
	// Ints.class tests to make sure that all the operators on ints function properly
	src := transpile("../examples/Ints.class")
	v := eval(src, "main.ints_Ints_Call_I_I")
	f, ok := v.(func(int32) int32)
	if !ok {
		t.Failed()
	}
	g := func(x int32) int32 {
		return x - 10*3/2 + 5%6<<3>>2&1 | 21 ^ 8
	}
	const num = 12
	if f(num) != g(num) {
		t.Failed()
	}
}

func transpile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	raw, err := parser.Parse(file)
	if err != nil {
		panic(err)
	}
	class, err := transformer.Simplify(raw)
	if err != nil {
		panic(err)
	}
	class.Methods = class.Methods[1:]
	gFile, err := transformer.Translate(class)
	if err != nil {
		panic(err)
	}
	var b = strings.Builder{}
	err = generate(gFile, &b)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func generate(g transformer.GoFile, w io.Writer) error {
	const temp = `package {{ .Package }}

{{if gt (len .Imports) 0 -}}
import (
	{{- range .Imports}}
	"{{ . }}"
	{{- end}}
)
{{- end}}
{{ range .Methods -}}
{{- if not .Receiver}}
func {{.Name}}({{ range .Params }}{{index . 0}} {{index . 1}}, {{end}}) {{.Return}} {
	{{.Code -}}
}
{{- end}}
{{end}}`
	t, err := template.New("gofile").Parse(temp)
	if err != nil {
		return err
	}
	return t.Execute(w, g)
}

func eval(src, method string) interface{} {
	i := interp.New(interp.Options{GoPath: "./"})

	_, err := i.Eval(src)
	if err != nil {
		panic(err)
	}

	v, err := i.Eval(method)
	if err != nil {
		panic(err)
	}
	return v.Interface()
}
