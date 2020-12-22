package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/totallygamerjet/jb4go/transformer"
)

func Generate(g transformer.GoFile) error {
	f := jen.NewFile(g.Package)
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
			params = append(params, jen.Id(p[0]).Id(p[1]))
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
