package transformer

import (
	"fmt"
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
	Code     string
	Return   string
}

func Translate(class JClass) (g GoFile, err error) {
	g.Package = class.Name[:strings.LastIndex(class.Name, "/")]
	g.FileName = g.Package + "_" + strings.ReplaceAll(class.SrcFileName, ".java", ".go")
	sT := getGoType(class.SuperName)
	g.Struct = Struct{
		Name:  ValidateName(class.Name),
		Embed: sT,
	}
	for _, v := range class.Fields {
		var f [2]string
		f[0] = ValidateName(v.Name)
		f[1] = getGoType(v.Type)
		if v.IsStatic {
			g.Globals = append(g.Globals, f)
		} else {
			g.Struct.Fields = append(g.Struct.Fields, f)
		}
	}
	for _, v := range class.Methods {
		m := Method{
			Name: translateMethodName(v.Name, g.Struct.Name, v.Return, v.IsStatic, v.IsPublic, v.Params),
		}
		if v.Return != "void" {
			m.Return = getGoType(v.Return)
		}
		nextArg := getUniqueCounter("arg")
		if !v.IsStatic {
			m.Receiver = nextArg() + " *" + g.Struct.Name
		}
		var isArray = false
		for _, v2 := range v.Params {
			if isArray {
				t := getGoType(v2) //type and import
				m.Params = append(m.Params, [3]string{nextArg(), "[]" + t})
				isArray = false
				continue
			}
			if v2 == "[" {
				isArray = true
				continue
			}
			t := getGoType(v2)
			m.Params = append(m.Params, [3]string{nextArg(), t})
		}
		m.Code = translateCode(v.Code)
		g.Methods = append(g.Methods, m)
		fmt.Println(m.Code)

	}
	return g, nil
}

func translateMethodName(mName, sName, Return string, isStatic, _ bool, params []string) (name string) {
	name = mName
	if name == "<init>" {
		name = "init"
	}
	if isStatic {
		name = sName + "_" + name
	} else {
		name = ValidateName(name)
	}
	name += "_"
	for _, p := range params {
		name += TranslateIdent(p)
	}
	name += "_" + TranslateIdent(Return)
	return name
}

// translateCode takes instructions in basic blocks and converts them to valid Go and returns it as a string
func translateCode(blocks []basicBlock) string {
	vars := strings.Builder{}
	b := strings.Builder{}
	for _, block := range blocks {
		b.WriteString(fmt.Sprintf("label%d:\n", block[0].Loc))
		for _, inst := range block {
			if inst.Op == nop { // ignore
				continue
			}
			if inst.Dest != "" && inst.Dest != "_" {
				if strings.HasPrefix(inst.Dest, "@") {
					inst.Dest = strings.ReplaceAll(inst.Dest, localName, "arg")
				}
				var t string
				if inst.Type != "" {
					t = getGoType(inst.Type)
				}
				vars.WriteString(fmt.Sprintf("var %s %s\n", inst.Dest, t))
				b.WriteString(fmt.Sprintf("%s = ", inst.Dest))
			}
			switch {
			case inst.Value != "":
				if strings.HasPrefix(inst.Value, "@") {
					inst.Value = strings.ReplaceAll(inst.Value, localName, "arg")
				}
				b.WriteString(inst.Value)
			case inst.Func != "":
				if inst.HasReceiver {
					b.WriteString(inst.Args[0] + ".")
					inst.Args = inst.Args[1:]
				}
				var p strings.Builder
				for _, v := range inst.Args {
					p.WriteString(v + ",")
				}
				params, ret := translateParams(inst.FDesc[strings.Index(inst.FDesc, ":")+1:])
				sName := ValidateName(inst.FDesc[:strings.Index(inst.FDesc, ".")]) //TODO: not always true
				b.WriteString(fmt.Sprintf("%s(%s)", translateMethodName(inst.Func, sName, ret, !inst.HasReceiver, true, params), p.String()))
			default:
				switch inst.Op {
				case getstatic:
					b.WriteString(fmt.Sprintf("%s_%s", ValidateName(inst.Args[0]), inst.Args[1]))
				case return_, areturn, ireturn:
					b.WriteString("return ")
					if len(inst.Args) > 0 {
						b.WriteString(inst.Args[0])
					}
				case ifge:
					b.WriteString(fmt.Sprintf("if %s >= 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case ifne:
					b.WriteString(fmt.Sprintf("if %s != 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case goto_:
					b.WriteString(fmt.Sprintf("goto label%s", inst.Args[0]))
				case new_:
					b.WriteString(fmt.Sprintf("new(%s)", ValidateName(inst.Args[0])))
				case irem:
					b.WriteString(fmt.Sprintf("%s %% %s", inst.Args[0], inst.Args[1]))
				default:
					b.WriteString("#####" + inst.String())
				}
			}
			b.WriteRune('\n')
		}
	}
	return vars.String() + b.String()
}
