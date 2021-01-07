package transformer

import (
	"encoding/binary"
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
			case aload:
				inst.Type = params[inst.operands[0]]
				inst.Args = []string{"a" + localName + strconv.Itoa(int(inst.operands[0]))}
				stack.push(inst.Args[0], inst.Type)
			case iload:
				inst.Type = intJ
				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.operands[0]))}
				stack.push(inst.Args[0], inst.Type)
			case fload:
				inst.Type = floatJ
				inst.Args = []string{"f" + localName + strconv.Itoa(int(inst.operands[0]))}
				stack.push(inst.Args[0], inst.Type)
			case dload:
				inst.Type = doubleJ
				inst.Args = []string{"d" + localName + strconv.Itoa(int(inst.operands[0]))}
				stack.push(inst.Args[0], inst.Type)
			case aload_0, aload_1, aload_2, aload_3:
				inst.Type = params[int(inst.Op-aload_0)]
				inst.Args = []string{"a" + localName + strconv.Itoa(int(inst.Op-aload_0))}
				stack.push(inst.Args[0], inst.Type)
			case dload_0, dload_1, dload_2, dload_3:
				inst.Type = doubleJ
				inst.Args = []string{"d" + localName + strconv.Itoa(int(inst.Op-dload_0))}
				stack.push(inst.Args[0], inst.Type)
			case iload_0, iload_1, iload_2, iload_3:
				inst.Type = intJ
				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.Op-iload_0))}
				stack.push(inst.Args[0], inst.Type)
			case lload_0, lload_1, lload_2, lload_3:
				inst.Type = longJ
				inst.Args = []string{"l" + localName + strconv.Itoa(int(inst.Op-lload_0))}
				stack.push(inst.Args[0], inst.Type)
			case aaload:
				v := nextVar()
				idx, _ := stack.pop()
				ref, _ := stack.pop()
				inst.Dest = v
				inst.Args = []string{ref, idx}
				inst.Type = "*java_lang_Object"
				stack.push(v, inst.Type)
			case iaload:
				idx, _ := stack.pop()
				ref, _ := stack.pop()
				stack.push(ref+"["+idx+"]", intJ)
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
			case istore, dstore, astore, fstore:
				var prefix string
				switch inst.Op {
				case istore:
					prefix = "i"
				case dstore:
					prefix = "d"
				case astore:
					prefix = "a"
				case fstore:
					prefix = "f"
				}
				inst.Dest = prefix + localName + strconv.Itoa(int(inst.operands[0]))
				inst.Value, inst.Type = stack.pop()
				if int(inst.operands[0]) >= len(params) {
					var p = make([]string, int(inst.operands[0])+1)
					copy(p, params)
					params = p
				}
				params[int(inst.operands[0])] = inst.Type
			case istore_0, istore_1, istore_2, istore_3:
				inst.Dest = "i" + localName + strconv.Itoa(int(inst.Op-istore_0))
				inst.Value, inst.Type = stack.pop()
			case lstore_0, lstore_1, lstore_2, lstore_3:
				inst.Dest = "l" + localName + strconv.Itoa(int(inst.Op-lstore_0))
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
				s, t := stack.pop()
				stack.push(s, t)
				stack.push(s, t)
			case i2d, l2d:
				s, _ := stack.pop()
				stack.push("float64("+s+")", doubleJ)
			case i2s:
				s, _ := stack.pop()
				stack.push("int32(int16("+s+"))", intJ)
			case i2l:
				s, _ := stack.pop()
				stack.push("int64("+s+")", longJ)
			case d2i, f2i:
				s, _ := stack.pop()
				stack.push("int32("+s+")", intJ)
			case i2c:
				s, _ := stack.pop()
				stack.push("int32(uint16("+s+"))", charJ)
			case d2f:
				s, _ := stack.pop()
				stack.push("float32("+s+")", floatJ)
			case i2b:
				s, _ := stack.pop()
				stack.push("int32(int8("+s+"))", intJ)
			case iadd, dadd, ladd:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"+"+i2, t)
			case isub, lsub:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"-"+i2, t)
			case imul, dmul, lmul, fmul:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"*"+i2, t)
			case idiv, ddiv, ldiv:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"/"+i2, t)
			case irem, lrem:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"%"+i2, t)
			case ineg, lneg:
				i1, t := stack.pop()
				stack.push("-"+i1, t)
			case ishl, lshl:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"<<"+i2, t)
			case ishr, lshr:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+">>"+i2, t)
			case iushr:
				i2, _ := stack.pop()
				i1, t := stack.pop() // go uses logical shift when types are unsigned
				stack.push(fmt.Sprintf("int32(uint32(%s)>>uint32(%s))", i1, i2), t)
			case lushr:
				i2, _ := stack.pop()
				i1, t := stack.pop() // go uses logical shift when types are unsigned
				stack.push(fmt.Sprintf("int64(uint64(%s)>>uint64(%s))", i1, i2), t)
			case iand, land:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"&"+i2, t)
			case ior, lor:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"|"+i2, t)
			case ixor, lxor:
				i2, _ := stack.pop()
				i1, t := stack.pop()
				stack.push(i1+"^"+i2, t)
			case return_:
			case ireturn, dreturn, areturn, lreturn:
				inst.Args = make([]string, 1)
				inst.Args[0], inst.Type = stack.pop()
			case getstatic:
				n, f, t := class.GetFieldRef(inst.index())
				inst.Type = getJavaType(t) // cleans up the type if has L...;
				inst.Args = []string{ValidateName(n) + "_" + f}
				stack.push(inst.Args[0], inst.Type)
			case ldc2_w:
				var a = [1]string{}
				inst.Args = (a)[:]
				inst.Args[0], inst.Type = class.GetConstant(inst.index())
				stack.push(inst.Args[0], inst.Type)
			case ldc:
				var v string
				v, inst.Type = class.GetConstant(int(inst.operands[0]))
				switch inst.Type { // do any necessary conversions for this constant
				case "java/lang/String":
					inst.Args = []string{fmt.Sprintf("newString(%s)", v)}
				case intJ, floatJ:
					inst.Args = []string{v}
				default:
					panic("unknown type: " + inst.Type)
				}
				stack.push(inst.Args[0], inst.Type)
			case iconst_m1, iconst_0, iconst_1, iconst_2, iconst_3, iconst_4, iconst_5:
				stack.push(strconv.Itoa(int(inst.Op)-int(iconst_0)), intJ)
			case lconst_0, lconst_1:
				stack.push(strconv.Itoa(int(inst.Op-lconst_0)), longJ)
			case dconst_0, dconst_1:
				stack.push(strconv.Itoa(int(inst.Op-dconst_0)), doubleJ)
			case aconst_null:
				stack.push("nil", objRefJ)
			case bipush:
				stack.push(strconv.Itoa(int(inst.operands[0])), intJ)
			case sipush:
				stack.push(strconv.Itoa(inst.index()), intJ)
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
			case ifge, ifne, ifgt, ifle, iflt, ifeq:
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
			case if_icmpge, if_icmplt, if_icmpgt, if_icmple, if_icmpne:
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
				ref, _ := stack.pop()
				stack.push("int32(len("+ref+"))", inst.Type)
			case goto_:
				inst.Args = []string{strconv.Itoa(inst.index() + inst.Loc)}
			case iinc:
				//TODO: make sure this is always an int and not some other type
				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.operands[0])), strconv.Itoa(int(inst.operands[1]))}
			case lookupswitch:
				key, _ := stack.pop()
				defaultOffset := int(binary.BigEndian.Uint32(inst.operands))
				npairs := int(binary.BigEndian.Uint32(inst.operands[4:]))
				inst.Args = []string{key, strconv.Itoa(defaultOffset + inst.Loc)}
				var i = 8 // starts at 8 bc defaultOffset and npairs
				for npairs > 0 {
					match := int(binary.BigEndian.Uint32(inst.operands[i:]))
					i += 4
					offset := int(binary.BigEndian.Uint32(inst.operands[i:]))
					i += 4
					inst.Args = append(inst.Args, strconv.Itoa(match), strconv.Itoa(offset+inst.Loc))
					npairs--
				}
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
