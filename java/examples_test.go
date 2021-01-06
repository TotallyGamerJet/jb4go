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

func Test_Longs(t *testing.T) {
	// Longs.class tests to make sure that all the operators on longs function properly
	src := transpile("../examples/Longs.class")
	v := eval(src, "main.longs_Longs_Call_J_J")
	f, ok := v.(func(int64) int64)
	if !ok {
		t.Failed()
	}
	g := func(l int64) int64 {
		l = -l
		l = l + 5
		l = l - 8
		l = l / 12
		l = l % 2
		l = l << 8
		l = l >> 2
		l = int64(uint64(l) >> 3)
		l = l | 0b100101101
		l = l ^ 21
		return l
	}
	const num = -503
	if f(num) != g(num) {
		t.Failed()
	}
}

func Test_SimpleArray(t *testing.T) {
	// SimpleArray.class tests creating, loading and storing in an int array
	src := transpile("../examples/SimpleArray.class")
	v := eval(src, "main.simple_SimpleArray_array_I_I")
	f, ok := v.(func(int32) int32)
	if !ok {
		t.Failed()
	}
	g := func(offset int32) int32 {
		var x [5]int32
		for i := 1; i < len(x); i++ {
			x[i] = offset - x[i-1]
		}
		return x[4]
	}
	const num = 27
	if f(num) != g(num) {
		t.Failed()
	}
}

func Test_BoolFunc(t *testing.T) {
	// BoolFunc.class test to make sure booleans function properly
	src := transpile("../examples/BoolFunc.class")
	v := eval(src, "main.bool_BoolFunc_bool_ZZII_I")
	f, ok := v.(func(int32, int32, int32, int32) int32)
	if !ok {
		t.Failed()
	}
	g := func(b, b2, x, y int32) int32 {
		if b != 0 {
			return x
		}
		if b2 == 0 {
			return y
		}
		if b != 0 && b2 != 0 {
			return x
		}
		if b != 0 || b2 != 0 {
			return y
		}
		return x + y
	}
	const (
		b  = 0 // 0 & 1 only its a bool!
		b2 = 1 // 0 & 1 only its a bool!
		x  = 5
		y  = 7
	)
	if f(b, b2, x, y) != g(b, b2, x, y) {
		t.Failed()
	}
}

func Test_ShortLoop(t *testing.T) {
	// ShortLoop.class test to make sure shorts function properly
	src := transpile("../examples/ShortLoop.class")
	v := eval(src, "main.shorts_ShortLoop_loop_S_S")
	f, ok := v.(func(int32) int32)
	if !ok {
		t.Failed()
	}
	g := func(x int32) int32 {
		var y = x
		for ; y < 100; y++ {
			x--
		}
		return x
	}
	const num = 7
	if f(num) != g(num) {
		t.Failed()
	}
}

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
		x = -x
		return int32(uint32(((((((x-15)+5)<<3)>>2)&1)|62680)^4104) >> 2)
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
