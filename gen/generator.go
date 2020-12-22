package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/totallygamerjet/jb4go/transformer"
)

func Generate(g transformer.GoFile) error {
	f := jen.NewFile(g.Package)
	for _, v := range g.Imports {
		f.ImportAlias(v[0], v[1])
	}
	for _, v := range g.Globals {
		f.Var().Id(v[0]).Id(v[1])
	}
	var code = []jen.Code{jen.Id(g.Struct.Embed)}
	for _, v := range g.Struct.Fields {
		code = append(code, jen.Id(v[0]).Id(v[1]))
	}
	f.Type().Id(g.Struct.Name).Struct(code...)
	for _, v := range g.Methods {
		fn := f.Func()
		if v.Receiver != "" {
			fn.Parens(jen.Id(v.Receiver))
		}
		fn.Id(v.Name)
		var params []jen.Code
		for _, p := range v.Params {
			params = append(params, genType(p))
		}
		fn.Params(params...)
		if v.Return != "" {
			fn.Id(v.Return)
		}
		fn.Block(
		//TODO: write code
		)
		f.Line()
	}
	return f.Save(g.FileName)
}

func genType(p [3]string) *jen.Statement {
	s := jen.Id(p[0])
	if p[2] != "" { // check if this type has an import
		s.Qual(p[2], p[1])
		return s
	}
	switch p[1] {
	case "int32":
		s.Int32()
	default:
		panic("not implemented")
	}
	return s
}
