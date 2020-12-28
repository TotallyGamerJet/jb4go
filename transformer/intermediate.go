package transformer

import (
	"fmt"
	"github.com/totallygamerjet/jb4go/parser"
	"strconv"
	"strings"
)

const (
	localName = "@local"
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

const (
	doubleJ = "double"
	intJ    = "int"
	objRefJ = "ObjRef"
)

// creates an intermediate form of the code
func createIntermediate(blocks []basicBlock, class parser.RawClass, params []string) {
	stack := make(stack, 0)
	nextVar := getUniqueCounter("v")
	for _, block := range blocks {
		for i, inst := range block {
			switch inst.Op {
			case nop: //ignore
				continue
			case aload_0, aload_1, aload_2, aload_3:
				v := nextVar()
				if len(params) > 0 {
					inst.Type = params[int(inst.Op-aload_0)]
				} else {
					inst.Type = objRefJ //TODO: remove this once all stores will add to params
				}
				inst.Dest = v
				inst.Value = localName + strconv.Itoa(int(inst.Op-aload_0))
				stack.push(v, inst.Type)
			case iload:
				v := nextVar()
				inst.Type = intJ
				inst.Dest = v
				inst.Value = localName + strconv.Itoa(int(inst.operands[0]))
				stack.push(v, inst.Type)
			case dload:
				v := nextVar()
				inst.Type = doubleJ
				inst.Dest = v
				inst.Value = localName + strconv.Itoa(int(inst.operands[0]))
				stack.push(v, inst.Type)
			case dload_0, dload_1, dload_2, dload_3:
				v := nextVar()
				inst.Type = doubleJ
				inst.Dest = v
				inst.Value = localName + strconv.Itoa(int(inst.Op-dload_0))
				stack.push(v, inst.Type)
			case iload_0, iload_1, iload_2, iload_3:
				v := nextVar()
				inst.Type = intJ
				inst.Dest = v
				inst.Value = localName + strconv.Itoa(int(inst.Op-iload_0))
				stack.push(v, inst.Type)
			case astore_0, astore_1, astore_2, astore_3:
				inst.Dest = localName + strconv.Itoa(int(inst.Op-astore_0))
				inst.Value, inst.Type = stack.pop()
				if int(inst.Op-astore_0) >= len(params) {
					var p [4]string
					copy(p[:], params)
					params = p[:]
				}
				params[int(inst.Op-astore_0)] = inst.Type
			case istore, dstore:
				inst.Dest = localName + strconv.Itoa(int(inst.operands[0]))
				inst.Value, inst.Type = stack.pop()
			case istore_0, istore_1, istore_2, istore_3:
				inst.Dest = localName + strconv.Itoa(int(inst.Op-istore_0))
				inst.Value, inst.Type = stack.pop()
			case invokespecial, invokevirtual, invokestatic:
				c, n, t := class.GetMethodRef(inst.index())
				p, _ := translateParams(t)
				var params []string
				for j := 0; j < len(p); j++ {
					v, _ := stack.pop()
					params = append([]string{v}, params...)
				}
				if inst.Op != invokestatic {
					r, t := stack.pop()                     // the receiver
					if inst.Op == invokespecial && c != t { // invokespecial calls inits so add the super class
						r += "." + ValidateName(c)
					}
					inst.HasReceiver = true
					inst.Args = append([]string{r}, params...) // start with receiver
				} else {
					inst.Args = params
				}
				m := c + "." + n + ":" + t
				inst.Func = n
				inst.FDesc = m
				if strings.HasSuffix(t, ")V") { //  ends in void
					inst.Dest = "_"
				} else {
					v := nextVar()
					inst.Dest = v
					inst.Type = getJavaType(t[strings.LastIndex(t, ")")+1:])
					stack.push(v, inst.Type)
				}
			case new_:
				v := nextVar()
				name := class.GetClass(inst.index())
				inst.Type = name
				inst.Dest = v
				inst.Args = []string{name}
				stack.push(v, inst.Type)
			case dup:
				v := nextVar()
				s, t := stack.pop()
				stack.push(s, t)
				stack.push(v, t)
				inst.Type = t
				inst.Dest = v
				inst.Value = s
			case i2d:
				v := nextVar()
				s, _ := stack.pop()
				inst.Type = doubleJ
				inst.Dest = v
				inst.Args = []string{s}
				stack.push(v, inst.Type)
			case irem, iadd, isub:
				v := nextVar()
				i2, _ := stack.pop()
				i1, _ := stack.pop()
				inst.Type = intJ
				inst.Dest = v
				inst.Args = []string{i1, i2}
				stack.push(v, inst.Type)
			case dmul, ddiv, dadd:
				v := nextVar()
				i2, _ := stack.pop()
				i1, _ := stack.pop()
				inst.Type = doubleJ
				inst.Dest = v
				inst.Args = []string{i1, i2}
				stack.push(v, inst.Type)
			case return_:
			case ireturn, dreturn, areturn:
				var p string
				p, inst.Type = stack.pop()
				inst.Args = []string{p}
			case getstatic:
				v := nextVar()
				var n, f string
				n, f, inst.Type = class.GetFieldRef(inst.index())
				inst.Type = getJavaType(inst.Type) // cleans up the type if has L...;
				inst.Dest = v
				inst.Args = []string{n, f}
				stack.push(v, inst.Type)
			case ldc2_w:
				v := nextVar()
				inst.Value, inst.Type = class.GetConstant(inst.index())
				inst.Dest = v
				stack.push(v, inst.Type)
			case ldc:
				v := nextVar()
				inst.Value, inst.Type = class.GetConstant(int(inst.operands[0]))
				switch inst.Type { // do any necessary conversions for this constant
				case "java/lang/String":
					inst.Value = fmt.Sprintf("New_string_G(%s)", inst.Value)
				}
				inst.Dest = v
				stack.push(v, inst.Type)
			case iconst_m1, iconst_0, iconst_1, iconst_2, iconst_3, iconst_4, iconst_5:
				v := nextVar()
				inst.Dest = v
				inst.Value = strconv.Itoa(int(inst.Op - iconst_0))
				inst.Type = intJ
				stack.push(v, inst.Type)
			case putfield:
				inst.Value, _ = stack.pop()
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				inst.Type = t
				inst.Dest = ref + "." + f
			case getfield:
				v := nextVar()
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				inst.Type = getJavaType(t) // confirm that is a java type
				inst.Dest = v
				inst.Value = ref + "." + f
				stack.push(v, inst.Type)
			case ifge:
				v, _ := stack.pop()
				inst.Args = []string{v, strconv.Itoa(inst.index() + inst.Loc)}
			case ifne:
				v, _ := stack.pop()
				inst.Args = []string{v, strconv.Itoa(inst.index() + inst.Loc)}
			case goto_:
				inst.Args = []string{strconv.Itoa(inst.index() + inst.Loc)}
			default:
				panic("unknown Op: " + inst.Op.String())
			}
			block[i] = inst // update the instruction
			//fmt.Println(inst)
		}
	}
	if len(stack) != 0 { // sanity check to make sure nothing went wrong
		panic(fmt.Sprintf("stack isn't empty: %s", stack))
	}
}
