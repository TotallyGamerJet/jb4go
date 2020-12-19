package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/totallygamerjet/jb4go/transformer"
)

func renderFields(name string, class transformer.Class, file *jen.File) {
	//TODO: get super
	codes := []jen.Code{jen.Qual("github.com/totallygamerjet/jb4go/java", "Java_lang_Object")}
	for _, f := range class.Fields {
		name := validateName(f.Name, f.IsPublic)
		if f.IsStatic {
			panic("not implemented")
		} else {
			s := jen.Id(name)
			toGoType(f.Type, s)
			codes = append(codes, s)
		}
	}
	file.Type().Id(name).Struct(
		codes...,
	)
}
