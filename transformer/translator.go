package transformer

import (
	"fmt"
	"strings"
)

type GoFile struct {
	FileName string
	Package  string
	Imports  []string
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
	g.Package = "main" // the generated code should be runnable
	sT := ValidateName(class.SuperName)
	g.Struct = Struct{
		Name:  ValidateName(class.Name),
		Embed: sT,
	}
	g.FileName = g.Struct.Name + ".go" // I expect the filename to be the same as the class name
	for _, v := range class.Fields {
		var f [2]string
		f[0] = "E_" + ValidateName(v.Name) // fields are prefixed with an export tag to make accessible
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
			m.Receiver = "a" + nextArg() + " *java_lang_Object"
		}
		for _, v2 := range v.Params {
			t := getGoType(v2.type_)
			if v2.isArray {
				t = "[]" + t
			}
			var prefix string
			switch t { //TODO: add more prefixes
			case "int32":
				prefix = "i"
			case "uint16":
				prefix = "c"
			case "float64":
				prefix = "d"
			default:
				prefix = "a"
			}
			m.Params = append(m.Params, [3]string{prefix + nextArg(), t})
			if t == "float64" { // doubles and longs take up two argument slots
				_ = nextArg()
			}
		}
		m.Code = translateCode(v.Code)
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
	} else {
		name = ValidateName(name)
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

// translateCode takes instructions in basic blocks and converts them to valid Go and returns it as a string
func translateCode(blocks []basicBlock) string {
	exists := make(map[string]bool)
	vars := strings.Builder{} // stores all the variables at the top of the function
	b := strings.Builder{}    // code goes in here
	for _, block := range blocks {
		b.WriteString(fmt.Sprintf("label%d:\n", block[0].Loc))
		for _, inst := range block {
			if inst.Op == nop { // ignore
				continue
			}
			if inst.Dest != "" && inst.Dest != "_" {
				if strings.Contains(inst.Dest, localName) {
					inst.Dest = strings.ReplaceAll(inst.Dest, localName, "arg")
				}
				var t string
				if inst.Type != "" {
					if inst.Op == anewarray {
						t = inst.Type
					} else {
						t = getGoType(inst.Type)
					}
				}
				if _, ok := exists[inst.Dest]; !strings.Contains(inst.Dest, ".") && !ok { // ignore fields in the var list
					vars.WriteString(fmt.Sprintf("var %s %s\n", inst.Dest, t))
					exists[inst.Dest] = true
				}
				b.WriteString(fmt.Sprintf("%s = ", inst.Dest))
			}
			switch {
			case inst.Value != "":
				if strings.Contains(inst.Value, localName) {
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
				sName := ValidateName(inst.FDesc[:strings.Index(inst.FDesc, ".")])
				if inst.HasReceiver {
					var methodCall string
					switch ret {
					case charJ, intJ:
						methodCall = "callMethodInt"
					case doubleJ:
						methodCall = "callMethodDouble"
					case voidJ: // TODO: more return types?
						methodCall = "callMethod"
					default:
						methodCall = "callMethodObject"
					}
					b.WriteString(fmt.Sprintf("%s(\"%s\", %s)", methodCall, translateMethodName(inst.Func, sName, ret, !inst.HasReceiver, true, params), p.String()))
				} else {
					b.WriteString(fmt.Sprintf("%s(%s)", translateMethodName(inst.Func, sName, ret, !inst.HasReceiver, true, params), p.String()))
				}
			default: // any complicated instructions go here
				switch inst.Op {
				case getstatic:
					b.WriteString(fmt.Sprintf("%s_%s", ValidateName(inst.Args[0]), inst.Args[1]))
				case getfield:
					var getF string
					switch inst.Type { //TODO: check the type and call the right method?
					case charJ, intJ:
						getF = "getFieldInt"
					case doubleJ:
						getF = "getFieldDouble"
					default:
						getF = "getFieldObject"
					}
					b.WriteString(fmt.Sprintf("%s.%s(\"E_%s\")", inst.Args[0], getF, inst.Args[1]))
				case putfield:
					b.WriteString(fmt.Sprintf("%s.setField(\"E_%s\", %s)", inst.Args[0], inst.Args[1], inst.Args[2]))
				case return_, areturn, ireturn, dreturn:
					b.WriteString("return ")
					if len(inst.Args) > 0 {
						b.WriteString(inst.Args[0])
					}
				case ifge:
					b.WriteString(fmt.Sprintf("if %s >= 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case ifgt:
					b.WriteString(fmt.Sprintf("if %s > 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case ifle:
					b.WriteString(fmt.Sprintf("if %s <= 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case iflt:
					b.WriteString(fmt.Sprintf("if %s < 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case ifne:
					b.WriteString(fmt.Sprintf("if %s != 0 { goto label%s }", inst.Args[0], inst.Args[1]))
				case if_icmpge:
					b.WriteString(fmt.Sprintf("if %s >= %s { goto label%s }", inst.Args[0], inst.Args[1], inst.Args[2]))
				case if_icmplt:
					b.WriteString(fmt.Sprintf("if %s < %s { goto label%s }", inst.Args[0], inst.Args[1], inst.Args[2]))
				case if_icmpgt:
					b.WriteString(fmt.Sprintf("if %s > %s { goto label%s }", inst.Args[0], inst.Args[1], inst.Args[2]))
				case dcmpg:
					b.WriteString(fmt.Sprintf("func(x, y float64) int32 {if x > y {return 1;} else if x == y {return 0;} else if x < y {return -1;}; return 1;}(%s, %s)", inst.Args[0], inst.Args[1]))
				case goto_:
					b.WriteString(fmt.Sprintf("goto label%s", inst.Args[0]))
				case new_:
					b.WriteString(fmt.Sprintf("new_%s()", ValidateName(inst.Args[0])))
				case irem:
					b.WriteString(fmt.Sprintf("%s %% %s", inst.Args[0], inst.Args[1]))
				case iadd, dadd:
					b.WriteString(fmt.Sprintf("%s + %s", inst.Args[0], inst.Args[1]))
				case isub:
					b.WriteString(fmt.Sprintf("%s - %s", inst.Args[0], inst.Args[1]))
				case imul, dmul:
					b.WriteString(fmt.Sprintf("%s * %s", inst.Args[0], inst.Args[1]))
				case idiv, ddiv:
					b.WriteString(fmt.Sprintf("%s / %s", inst.Args[0], inst.Args[1]))
				case i2d:
					b.WriteString(fmt.Sprintf("float64(%s)", inst.Args[0]))
				case iinc:
					b.WriteString(fmt.Sprintf("%s += %s", strings.ReplaceAll(inst.Args[0], localName, "arg"), inst.Args[1]))
				case pop:
					b.WriteString(fmt.Sprintf("_ = %s", inst.Args[0]))
				case newarray:
					b.WriteString(fmt.Sprintf("make([]%s, %s)", getGoType(inst.Args[1]), inst.Args[0]))
				case anewarray:
					b.WriteString(fmt.Sprintf("make([]%s, %s)", getGoType(inst.Args[1]), inst.Args[0]))
				case arraylength:
					b.WriteString(fmt.Sprintf("int32(len(%s))", inst.Args[0]))
				case aaload, iaload:
					b.WriteString(fmt.Sprintf("%s[%s]", inst.Args[0], inst.Args[1]))
				case aastore, iastore:
					b.WriteString(fmt.Sprintf("%s[%s] = %s", inst.Args[0], inst.Args[1], inst.Args[2]))
				default:
					panic(inst.String())
				}
			}
			b.WriteRune('\n')
		}
	}
	return vars.String() + b.String()
}
