package transformer

import (
	"fmt"
	"github.com/totallygamerjet/jb4go/parser"
	"strconv"
	"strings"
)

const (
	localName = "@local"
	prelude   = `var currentLabel = 0
	controlFlowLoop:
	for {
		switch currentLabel {`
	envoi = `default:
			panic("unexpected control flow")
		}
	}`
)

type stack []struct {
	v string
	t string
}

func (s *stack) push(v, t string) {
	*s = append(*s, struct {
		v string
		t string
	}{v: v, t: t})
}

func (s *stack) pop() (string, string) {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v.v, v.t
}

type ssaInstruction struct {
	Op          opcode
	Type        string
	Dest        string
	Value       string
	Args        []string // the first arg is the receive if is a method
	Func        string
	FDesc       string // describes the function params and return
	HasReceiver bool   // does this method have a receiver
	//Funcs []string
	//Labels []string
}

func (i ssaInstruction) String() (s string) {
	if len(i.Dest) > 0 {
		if i.Dest != "_" && !strings.HasPrefix(i.Dest, "@") {
			s += "var "
		}
		s += i.Dest + " " + i.Type + " = " + i.Value
	}
	if i.Func != "" {
		if i.HasReceiver {
			s += fmt.Sprintf("%s.%s(%s)", i.Args[0], i.Func, i.Args[1:])
		} else {
			s += fmt.Sprintf("%s(%s)", i.Func, i.Args)
		}
		s += "//" + i.Op.String()
	} else if len(i.Args) > 0 {
		s += fmt.Sprintf("%s %s", i.Op.String(), i.Args)
	} else {
		s += " " + i.Op.String()
	}
	return s
}

type ssaL []ssaInstruction

func (s *ssaL) add(i ssaInstruction) {
	*s = append(*s, i)
}

const (
	intJ    = "int"
	objRefJ = "ObjRef"
)

// creates an intermediate form of the code
func createIntermediate(blocks []basicBlock, class parser.RawClass) []ssaL {
	stack := make(stack, 0)
	nextVar := getUniqueCounter("v")
	//var inter = []string{prelude}
	var cfg []ssaL
	var ssa = make(ssaL, 0)
	for _, block := range blocks {
		ssa.add(ssaInstruction{Func: "case", Args: []string{strconv.Itoa(block[0].loc)}})
		for _, inst := range block {
			ssaI := ssaInstruction{Op: inst.opcode}
			switch inst.opcode {
			case nop: //ignore
				continue
			case aload_0, aload_1, aload_2, aload_3:
				v := nextVar()
				ssaI.Type = objRefJ
				ssaI.Dest = v
				ssaI.Value = localName + strconv.Itoa(int(inst.opcode-aload_0))
				stack.push(v, objRefJ)
			case iload:
				v := nextVar()
				ssaI.Type = intJ
				ssaI.Dest = v
				ssaI.Value = localName + strconv.Itoa(int(inst.operands[0]))
				stack.push(v, intJ)
			case iload_0, iload_1, iload_2, iload_3:
				v := nextVar()
				ssaI.Type = intJ
				ssaI.Dest = v
				ssaI.Value = localName + strconv.Itoa(int(inst.opcode-iload_0))
				stack.push(v, intJ)
			case astore_0, astore_1, astore_2, astore_3:
				ssaI.Dest = localName + strconv.Itoa(int(inst.opcode-astore_0))
				ssaI.Value, ssaI.Type = stack.pop()
			case istore:
				ssaI.Dest = localName + strconv.Itoa(int(inst.operands[0]))
				ssaI.Value, ssaI.Type = stack.pop()
			case istore_0, istore_1, istore_2, istore_3:
				ssaI.Dest = localName + strconv.Itoa(int(inst.opcode-istore_0))
				ssaI.Value, ssaI.Type = stack.pop()
			case invokespecial, invokevirtual, invokestatic:
				c, n, t := class.GetMethodRef(inst.index())
				p, _ := translateParams(t)
				var params []string
				for j := 0; j < len(p); j++ {
					v, _ := stack.pop()
					params = append([]string{v}, params...)
				}
				if inst.opcode != invokestatic {
					r, _ := stack.pop() // the receiver
					ssaI.HasReceiver = true
					ssaI.Args = append([]string{r}, params...) // start with receiver
				} else {
					ssaI.Args = params
				}
				m := c + "." + n + ":" + t
				ssaI.Func = n
				ssaI.FDesc = m
				if strings.HasSuffix(t, ")V") { //  ends in void
					ssaI.Dest = "_"
				} else {
					v := nextVar()
					ssaI.Dest = v
					ssaI.Type = getJavaType(t[strings.LastIndex(t, ")")+1:])
					stack.push(v, ssaI.Type)
				}
			case new_:
				v := nextVar()
				name := class.GetClass(inst.index())
				ssaI.Type = name
				ssaI.Dest = v
				ssaI.Args = []string{name}
				stack.push(v, ssaI.Type)
			case dup:
				v := nextVar()
				s, t := stack.pop()
				stack.push(s, t)
				stack.push(v, t)
				ssaI.Type = t
				ssaI.Dest = v
				ssaI.Value = s
			case irem:
				v := nextVar()
				i2, _ := stack.pop()
				i1, _ := stack.pop()
				ssaI.Type = intJ
				ssaI.Dest = v
				ssaI.Args = []string{i1, i2}
				stack.push(v, ssaI.Type)
			case return_:
			case ireturn:
				var p string
				p, ssaI.Type = stack.pop()
				ssaI.Args = []string{p}
			case areturn:
				var p string
				p, ssaI.Type = stack.pop()
				ssaI.Args = []string{p}
			case getstatic:
				v := nextVar()
				var n, f string
				n, f, ssaI.Type = class.GetFieldRef(inst.index())
				ssaI.Type = getJavaType(ssaI.Type) // cleans up the type if has L...;
				ssaI.Dest = v
				ssaI.Args = []string{n, f}
				stack.push(v, ssaI.Type)
			case ldc:
				v := nextVar()
				ssaI.Value, ssaI.Type = class.GetConstant(int(inst.operands[0]))
				ssaI.Dest = v
				stack.push(v, ssaI.Type)
			case putfield:
				ssaI.Value, _ = stack.pop()
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				ssaI.Type = t
				ssaI.Dest = ref + "." + f
			case getfield:
				v := nextVar()
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				ssaI.Type = t
				ssaI.Dest = v
				ssaI.Value = ref + "." + f
				stack.push(v, ssaI.Type)
			case ifge:
				v, _ := stack.pop()
				ssaI.Args = []string{v, strconv.Itoa(inst.index() + inst.loc)}
			case ifne:
				v, _ := stack.pop()
				ssaI.Args = []string{v, strconv.Itoa(inst.index() + inst.loc)}
			case goto_:
				ssaI.Args = []string{strconv.Itoa(inst.index() + inst.loc)}
			default:
				panic("unknown opcode: " + inst.opcode.String())
			}
			ssa.add(ssaI)
			fmt.Println(ssaI)
		}
		cfg = append(cfg, ssa)
		ssa = make(ssaL, 0)
	}
	if len(stack) != 0 { // sanity check to make sure nothing went wrong
		panic(fmt.Sprintf("stack isn't empty: %s", stack))
	}
	return cfg
}
