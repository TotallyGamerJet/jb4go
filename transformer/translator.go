package transformer

import (
	"strconv"
	"strings"
)

type GoFile struct {
	FileName string
	Package  string
	Globals  [][2]string //name type //TODO: handle const vs var
	Struct   Struct
	Methods  []Method
}

type Struct struct {
	Name   string
	Embed  string      // the super object
	Fields [][2]string //name, type
}

type Method struct {
	Name     string
	Receiver string      // empty if static
	Params   [][2]string //name type
	//TODO: code
	Return string
}

func Translate(class JClass) (g GoFile, err error) {
	g.FileName = strings.ReplaceAll(class.SrcFileName, ".java", ".go")
	g.Package = class.Name[:strings.LastIndex(class.Name, "/")]
	g.Struct = Struct{
		Name:  ValidateName(class.Name, class.IsPublic),
		Embed: ValidateName(class.SuperName, true), // is this always true?
	}
	for _, v := range class.Fields {
		var f [2]string
		f[0] = ValidateName(v.Name, v.IsPublic)
		f[1] = getGoType(v.Type)
		if v.IsStatic {
			g.Globals = append(g.Globals, f)
		} else {
			g.Struct.Fields = append(g.Struct.Fields, f)
		}
	}
	for _, v := range class.Methods {
		m := Method{
			Name: translateMethodName(v),
		}
		if v.Return != "void" {
			m.Return = getGoType(v.Return)
		}
		var argN = 0
		if !v.IsStatic {
			m.Receiver = "arg" + strconv.Itoa(argN) + " *" + g.Struct.Name
			argN++
		}
		var isArray = false
		for _, v2 := range v.Params {
			if isArray {
				m.Params = append(m.Params, [2]string{"arg" + strconv.Itoa(argN), "[]" + getGoType(v2)})
				isArray = false
				argN++
				continue
			}
			if v2 == "[" {
				isArray = true
				continue
			}
			m.Params = append(m.Params, [2]string{"arg" + strconv.Itoa(argN), getGoType(v2)})
			argN++
		}
		//TODO: code
		g.Methods = append(g.Methods, m)
	}
	return g, nil
}

func translateMethodName(v JMethod) (name string) {
	name = v.Name
	if name == "<init>" {
		name = "init"
	}
	name = ValidateName(name, v.IsPublic)
	name += "_"
	for _, p := range v.Params {
		name += TranslateIdent(p)
	}
	name += "_" + TranslateIdent(v.Return)
	return name
}
