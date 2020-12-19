package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/totallygamerjet/jb4go/transformer"
	"strings"
)

func Transpile(class transformer.Class) error {
	className := class.Name[strings.LastIndex(class.Name, "/")+1:] //class.Name
	className = validateName(className, class.IsPublic)
	packge := strings.ReplaceAll(class.Name[:strings.LastIndex(class.Name, "/")], "/", "_")
	f := jen.NewFile(packge)
	renderFields(className, class, f)

	for _, v := range class.Methods {
		transpileMethod(v, className, f)
		f.Line()
	}
	return f.Save(strings.ReplaceAll(class.SrcFileName, ".java", ".go"))
}
