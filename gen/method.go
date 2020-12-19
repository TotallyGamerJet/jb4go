package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/totallygamerjet/jb4go/transformer"
	"strconv"
	"strings"
)

func transpileMethod(met transformer.Method, className string, file *jen.File) {
	fName := translateMethodName(met.Name, met.Params, met.Return, met.IsPublic)
	f := file.Func()
	if !met.IsStatic {
		f.Params(jen.Id("this").Id("*" + className))
	}
	renderParams(f.Id(fName), met)
	if met.Return != "void" {
		toGoType(met.Return, f)
	}
	f.Block(
	// TODO:
	)
}

// translateMethodName takes in the Java version name, list of params, return type and if its public
// and converts it to a Go compatible string that adheres to the Java Bytecode Transpiler
// specification. Parameters and return type must be valid java builtin types or Objects
func translateMethodName(name string, params []string, ret string, isPublic bool) string {
	if name == "<init>" { // a constructor
		name = "new"
	}
	name += "_" // start of params
	for _, v := range params {
		name += TranslateIdent(v)
	}
	name += "_" + TranslateIdent(ret) // for return
	name = validateName(name, isPublic)
	return name
}

func renderParams(file *jen.Statement, met transformer.Method) {
	var params []jen.Code
	var n int // keep track of arg number
	var arrayCount int
	for _, v := range met.Params {
		if v == "[" {
			arrayCount++
			continue
		}
		s := jen.Id("arg" + strconv.Itoa(n))
		if arrayCount > 0 {
			// TODO:
			strings.Repeat("[]", arrayCount)
		} else {
			toGoType(v, s)
		}
		params = append(params, s)
		n++
	}
	file.Params(params...)
}
