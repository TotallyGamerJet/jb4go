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
				inst.Type = params[int(inst.Op-aload_0)]
				inst.Args = []string{"a" + localName + strconv.Itoa(int(inst.Op-aload_0))}
				stack.push(inst.Args[0], inst.Type)
			case aload:
				v := nextVar()
				inst.Type = params[inst.operands[0]]
				inst.Dest = v
				inst.Value = "a" + localName + strconv.Itoa(int(inst.operands[0]))
				stack.push(v, inst.Type)
			case iload:
				v := nextVar()
				inst.Type = intJ
				inst.Dest = v
				inst.Value = "i" + localName + strconv.Itoa(int(inst.operands[0]))
				stack.push(v, inst.Type)
			case dload:
				v := nextVar()
				inst.Type = doubleJ
				inst.Dest = v
				inst.Value = "d" + localName + strconv.Itoa(int(inst.operands[0]))
				stack.push(v, inst.Type)
			case dload_0, dload_1, dload_2, dload_3:
				v := nextVar()
				inst.Type = doubleJ
				inst.Dest = v
				inst.Value = "d" + localName + strconv.Itoa(int(inst.Op-dload_0))
				stack.push(v, inst.Type)
			case iload_0, iload_1, iload_2, iload_3:
				v := nextVar()
				inst.Type = intJ
				inst.Dest = v
				inst.Value = "i" + localName + strconv.Itoa(int(inst.Op-iload_0))
				stack.push(v, inst.Type)
			case aaload:
				v := nextVar()
				idx, _ := stack.pop()
				ref, _ := stack.pop()
				inst.Dest = v
				inst.Args = []string{ref, idx}
				inst.Type = "*java_lang_Object"
				stack.push(v, inst.Type)
			case iaload:
				v := nextVar()
				idx, _ := stack.pop()
				ref, _ := stack.pop()
				inst.Dest = v
				inst.Args = []string{ref, idx}
				inst.Type = intJ
				stack.push(v, inst.Type)
			case astore_0, astore_1, astore_2, astore_3:
				inst.Dest = "a" + localName + strconv.Itoa(int(inst.Op-astore_0))
				inst.Value, inst.Type = stack.pop()
				if int(inst.Op-astore_0) >= len(params) {
					var p [4]string
					copy(p[:], params)
					params = p[:]
				}
				params[int(inst.Op-astore_0)] = inst.Type
			case aastore, iastore:
				val, _ := stack.pop()
				index, _ := stack.pop()
				ref, _ := stack.pop()
				inst.Args = []string{ref, index, val}
			case istore, dstore, astore:
				var prefix string
				switch inst.Op {
				case istore:
					prefix = "i"
				case dstore:
					prefix = "d"
				case astore:
					prefix = "a"
				}
				inst.Dest = prefix + localName + strconv.Itoa(int(inst.operands[0]))
				inst.Value, inst.Type = stack.pop()
				if int(inst.operands[0]) >= len(params) {
					var p = make([]string, int(inst.operands[0])+1)
					copy(p[:], params)
					params = p[:]
				}
				params[int(inst.operands[0])] = inst.Type
			case istore_0, istore_1, istore_2, istore_3:
				inst.Dest = "i" + localName + strconv.Itoa(int(inst.Op-istore_0))
				inst.Value, inst.Type = stack.pop()
			case dstore_0, dstore_1, dstore_2, dstore_3:
				inst.Dest = "d" + localName + strconv.Itoa(int(inst.Op-dstore_0))
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
						r += ".super" // + ValidateName(c)
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
			case irem, iadd, isub, imul, idiv, ishl, ishr, iand, ior:
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
					inst.Value = fmt.Sprintf("newString(%s)", inst.Value)
				}
				inst.Dest = v
				stack.push(v, inst.Type)
			case iconst_m1, iconst_0, iconst_1, iconst_2, iconst_3, iconst_4, iconst_5:
				v := nextVar()
				inst.Dest = v
				inst.Value = strconv.Itoa(int(inst.Op - iconst_0))
				inst.Type = intJ
				stack.push(v, inst.Type)
			case dconst_0, dconst_1:
				v := nextVar()
				inst.Dest = v
				inst.Value = strconv.Itoa(int(inst.Op - dconst_0))
				inst.Type = doubleJ
				stack.push(v, inst.Type)
			case putfield:
				v, _ := stack.pop()
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				inst.Type = t
				inst.Args = []string{ref, f, v}
			case getfield:
				v := nextVar()
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				inst.Type = getJavaType(t) // confirm that is a java type
				inst.Dest = v
				inst.Args = []string{ref, f}
				stack.push(v, inst.Type)
			case pop:
				ref, _ := stack.pop()
				inst.Dest = "_"
				inst.Args = []string{ref}
			case bipush:
				v := nextVar()
				inst.Dest = v
				inst.Value = strconv.Itoa(int(inst.operands[0]))
				inst.Type = intJ
				stack.push(v, inst.Type)
			case ifge, ifne, ifgt, ifle, iflt:
				v, _ := stack.pop()
				inst.Args = []string{v, strconv.Itoa(inst.index() + inst.Loc)}
			case dcmpg:
				v := nextVar()
				a2, _ := stack.pop()
				a1, _ := stack.pop()
				inst.Dest = v
				inst.Type = intJ
				inst.Args = []string{a1, a2}
				stack.push(v, inst.Type)
			case if_icmpge, if_icmplt, if_icmpgt:
				v2, _ := stack.pop()
				v1, _ := stack.pop()
				inst.Args = []string{v1, v2, strconv.Itoa(inst.index() + inst.Loc)}
			case anewarray:
				v := nextVar()
				size, _ := stack.pop()
				c := class.GetClass(inst.index())
				inst.Dest = v
				inst.Type = "[]" + getGoType(c)
				inst.Args = []string{size, c}
				stack.push(v, inst.Type)
			case newarray:
				v := nextVar()
				size, _ := stack.pop()
				t := arrayTypeCodes(int(inst.operands[0]))
				inst.Dest = v
				inst.Type = "[]" + getGoType(t)
				inst.Args = []string{size, t}
				stack.push(v, inst.Type)
			case arraylength:
				v := nextVar()
				ref, _ := stack.pop()
				inst.Dest = v
				inst.Type = intJ
				inst.Args = []string{ref}
				stack.push(v, inst.Type)
			case goto_:
				inst.Args = []string{strconv.Itoa(inst.index() + inst.Loc)}
			case iinc:
				//TODO: make sure this is always an int and not some other type
				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.operands[0])), strconv.Itoa(int(inst.operands[1]))}
			default:
				panic("unknown Op: " + inst.Op.String())
			}
			block[i] = inst // update the instruction
		}
	}
	if len(stack) != 0 { // sanity check to make sure nothing went wrong
		panic(fmt.Sprintf("stack isn't empty: %s", stack))
	}
}
