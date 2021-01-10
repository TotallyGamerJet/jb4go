package transformer

type GoFile struct {
	FileName string
	Package  string
	Imports  []string
	Globals  []Field //TODO: handle const vs var
	Struct   Struct
	Methods  []Method
}

type Struct struct {
	Name   string
	Embed  string // the super object
	Fields []Field
}

type Field struct {
	Name, Type, Value string
}

type Method struct {
	Name     string
	Receiver string      // empty if static
	Params   [][3]string //Name, type, (import if available)
	Code     string
	Return   string
}

func Translate(class JClass) (g GoFile, err error) {
	g.Package = "main" // the generated code should be runnable
	sT := ValidateName(class.SuperName)
	g.Struct = Struct{
		Name:  ValidateName(class.Name),
		Embed: sT,
	}
	g.FileName = g.Struct.Name + ".go" // I expect the filename to be the same as the class Name
	for _, v := range class.Fields {
		var f = Field{
			Name: ValidateName(v.Name), // fields are prefixed with an export tag to make accessible
			Type: getGoType(v.Type).type_,
		}
		if v.IsStatic {
			f.Name = g.Struct.Name + "_" + f.Name
			f.Value = v.Value
			g.Globals = append(g.Globals, f)
		} else {
			f.Name = "E_" + f.Name
			g.Struct.Fields = append(g.Struct.Fields, f)
		}
	}
	for _, v := range class.Methods {
		m := Method{
			Name: translateMethodName(v.Name, g.Struct.Name, v.Return, v.IsStatic, v.IsPublic, v.Params),
		}
		if v.IsAbstract { // ignore abstract methods
			m.Receiver = "nil"
			g.Methods = append(g.Methods, m)
			continue
		}
		if v.Return != "void" {
			m.Return = getGoType(v.Return).type_
		}
		nextArg := getUniqueCounter("arg")
		if !v.IsStatic {
			m.Receiver = "a" + nextArg() + " *java_lang_Object"
		}
		for _, v2 := range v.Params {
			t := getGoType(v2.String())
			var prefix = getPrefix(t)
			m.Params = append(m.Params, [3]string{prefix + nextArg(), t.String()})
			if t.type_ == "float64" { // doubles and longs take up two argument slots
				_ = nextArg()
			}
		}
		m.Code = v.Code
		g.Methods = append(g.Methods, m)
		if v.Name == "main" && v.IsStatic && v.Return == "void" { // add a real main method to call the java generated one
			g.Imports = append(g.Imports, "os")
			g.Methods = append(g.Methods, Method{
				Name: "main",
				Code: getMain(m.Name, len(m.Params)),
			})
		}
	}
	return g, nil
}

func getMain(name string, paramsN int) string {
	if paramsN > 0 {
		return `args := make([]*java_lang_Object, len(os.Args))
	for i, v := range os.Args {
		args[i] = newString(v)
	}
	` + name + `(args)
`
	}
	return name + "()\n"
}

func translateMethodName(mName, sName, Return string, isStatic, _ bool, params []nameAndType) (name string) {
	name = mName
	if name == "<init>" {
		name = "init"
	}
	if isStatic {
		name = sName + "_" + name
	}
	name += "_"
	for _, p := range params {
		if p.isArray {
			name += TranslateIdent("[")
		}
		name += TranslateIdent(p.type_)
	}
	name += "_" + TranslateIdent(Return)
	return name
}
