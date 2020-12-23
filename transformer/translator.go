package transformer

import (
	"strings"
)

type GoFile struct {
	FileName string
	Package  string
	Imports  [][2]string // import_path, alias
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
	Params   [][3]string //name, type, (import if available)
	//TODO: code
	Return string
}

func Translate(class JClass) (g GoFile, err error) {
	g.Package = class.Name[:strings.LastIndex(class.Name, "/")]
	g.FileName = g.Package + "_" + strings.ReplaceAll(class.SrcFileName, ".java", ".go")
	g.Imports = [][2]string{{"github.com/totallygamerjet/jb4go/java", "."}} //TODO: get other imports
	sT, _ := getGoType(class.SuperName)
	g.Struct = Struct{
		Name:  ValidateName(class.Name, class.IsPublic),
		Embed: sT,
	}
	for _, v := range class.Fields {
		var f [2]string
		f[0] = ValidateName(v.Name, v.IsPublic)
		f[1], _ = getGoType(v.Type)
		if v.IsStatic {
			g.Globals = append(g.Globals, f)
		} else {
			g.Struct.Fields = append(g.Struct.Fields, f)
		}
	}
	for _, v := range class.Methods {
		m := Method{
			Name: translateMethodName(v, g.Struct),
		}
		if v.Return != "void" {
			m.Return, _ = getGoType(v.Return)
		}
		nextArg := getUniqueCounter("arg")
		if !v.IsStatic {
			m.Receiver = nextArg() + " *" + g.Struct.Name
		}
		var isArray = false
		for _, v2 := range v.Params {
			if isArray {
				t, i := getGoType(v2) //type and import
				m.Params = append(m.Params, [3]string{nextArg(), "[]" + t, i})
				isArray = false
				continue
			}
			if v2 == "[" {
				isArray = true
				continue
			}
			t, i := getGoType(v2)
			m.Params = append(m.Params, [3]string{nextArg(), t, i})
		}
		//TODO: code
		g.Methods = append(g.Methods, m)
	}
	return g, nil
}

func translateMethodName(v JMethod, s Struct) (name string) {
	name = v.Name
	if name == "<init>" {
		name = "init"
	}
	if v.IsStatic {
		name = s.Name + "_" + name
	} else {
		name = ValidateName(name, v.IsPublic)
	}

	name += "_"
	for _, p := range v.Params {
		name += TranslateIdent(p)
	}
	name += "_" + TranslateIdent(v.Return)
	return name
}
