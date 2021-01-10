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

func translate(blocks []basicBlock, cfg map[int][]int, stackSize, localSize int, params []nameAndType, raw parser.RawClass) string {
	vB := newVarBuilder(params, localSize)
	stack := newStack(stackSize)
	var code strings.Builder
	translateBlock(&stack, &vB, blocks, cfg, raw, &code)
	return vB.String() + "\n" + code.String()
}

type varBuilder struct {
	marker     int
	marked     bool
	localTable []nameAndType
	vars       map[string]struct{}
	counter    int
	b          strings.Builder
}

func newVarBuilder(params []nameAndType, localSize int) varBuilder {
	table := make([]nameAndType, localSize)
	var idx = 0
	for _, v := range params {
		t := getGoType(v.String())
		table[idx] = t
		if t.type_ == "float64" || t.type_ == "int64" {
			idx += 2
		} else {
			idx++
		}
	}
	return varBuilder{vars: make(map[string]struct{}), localTable: table}
}

func (v *varBuilder) String() string {
	return v.b.String()
}

func (v *varBuilder) mark() {
	v.marker = v.counter
	v.marked = true
}

func (v *varBuilder) reset() {
	if !v.marked {
		return
	}
	v.counter = v.marker
	v.marked = false
}

func (v *varBuilder) newVar(t nameAndType) string {
	vName := "v" + strconv.Itoa(v.counter)
	vName = getPrefix(t) + vName
	if _, ok := v.vars[vName]; !ok {
		v.b.WriteString(fmt.Sprintf("var %s %s\n", vName, t))
		v.vars[vName] = struct{}{}
	}
	v.counter++
	return vName
}

func (v *varBuilder) addLocal(loc int, n string, t nameAndType) {
	if v.localTable[loc] == t {
		return
	}
	v.localTable[loc] = t
	v.b.WriteString(fmt.Sprintf("var %s %s\n", n, t))
}

func (v *varBuilder) typeOf(loc int) nameAndType {
	return v.localTable[loc]
}

type stack struct {
	i      int
	marked bool
	marker int
	data   []struct {
		v string
		t nameAndType
	}
}

func (s *stack) mark() {
	s.marked = true
	s.marker = s.i
}

func (s *stack) reset() {
	if !s.marked {
		return
	}
	s.i = s.marker
	s.marked = false
}

func (s *stack) push(v string, t nameAndType) {
	s.data[s.i] = struct {
		v string
		t nameAndType
	}{v: v, t: t}
	s.i++
}

func (s *stack) pop() (string, nameAndType) {
	s.i--
	val := s.data[s.i]
	return val.v, val.t
}

func (s *stack) length() int {
	return s.i
}

func newStack(size int) stack {
	return stack{data: make([]struct {
		v string
		t nameAndType
	}, size)}
}

// map compare opcodes to their symbol
var compare = map[opcode]string{
	ifeq:      "==",
	ifne:      "!=",
	iflt:      "<",
	ifge:      ">=",
	if_icmplt: "<",
	if_icmpge: ">=",
	if_icmpgt: ">",
}

func translateBlock(stack *stack, vB *varBuilder, blocks []basicBlock, cfg map[int][]int, class parser.RawClass, code *strings.Builder) {
	nextVar := func(t string, isArray bool) string {
		v := vB.newVar(nameAndType{type_: t, isArray: isArray})
		stack.push(v, nameAndType{type_: t, isArray: isArray})
		return v
	}
	sym := func(symbol string) {
		v2, _ := stack.pop()
		v1, t := stack.pop()
		v := nextVar(t.type_, t.isArray)
		code.WriteString(fmt.Sprintf("%s = %s %s %s\n", v, v1, symbol, v2))
	}
	set := func(dest, val string) {
		code.WriteString(fmt.Sprintf("%s = %s\n", dest, val))
	}
	load := func(prefix string, loc int) {
		t := vB.typeOf(loc)
		v := nextVar(t.type_, t.isArray)
		n := prefix + "arg" + strconv.Itoa(loc)
		set(v, n)
	}
	store := func(prefix string, loc int) {
		v, t := stack.pop()
		n := prefix + "arg" + strconv.Itoa(loc)
		vB.addLocal(loc, n, t)
		set(n, v)
	}
	cast := func(newType, format string) {
		v1, _ := stack.pop()
		v := nextVar(newType, false)
		code.WriteString(fmt.Sprintf("%s = "+format+"\n", v, v1))
	}
nextBlock:
	for id, block := range blocks {
		if id == 0 {
			// TODO: figure out if we need it
		} else {
			code.WriteString(fmt.Sprintf("label%d:\n", block.getLineStart()))
		}
		var last instruction
		for _, inst := range block {
			if inst.Op == nop {
				continue
			}
			switch inst.Op {
			case aload:
				load("a", int(inst.operands[0]))
			case astore:
				store("a", int(inst.operands[0]))
			case aload_0, aload_1, aload_2, aload_3:
				load("a", int(inst.Op-aload_0))
			case astore_0, astore_1, astore_2, astore_3:
				store("a", int(inst.Op-astore_0))
			case iload:
				load("i", int(inst.operands[0]))
			case istore:
				store("i", int(inst.operands[0]))
			case iload_0, iload_1, iload_2, iload_3:
				load("i", int(inst.Op-iload_0))
			case istore_0, istore_1, istore_2, istore_3:
				store("i", int(inst.Op-istore_0))
			case lload_0, lload_1, lload_2, lload_3:
				load("l", int(inst.Op-lload_0))
			case lstore_0, lstore_1, lstore_2, lstore_3:
				store("l", int(inst.Op-lstore_0))
			case dload:
				load("d", int(inst.operands[0]))
			case dstore:
				store("d", int(inst.operands[0]))
			case dload_0, dload_1, dload_2, dload_3:
				load("d", int(inst.Op-dload_0))
			case dstore_0, dstore_1, dstore_2, dstore_3:
				store("d", int(inst.Op-dstore_0))
			case fload:
				load("f", int(inst.operands[0]))
			case fstore:
				store("f", int(inst.operands[0]))
			case iaload, aaload:
				idx, _ := stack.pop()
				ref, t := stack.pop()
				v := nextVar(t.type_, false)
				set(v, ref+"["+idx+"]")
			case iastore, aastore:
				val, _ := stack.pop()
				idx, _ := stack.pop()
				ref, _ := stack.pop()
				code.WriteString(fmt.Sprintf("%s[%s] = %s\n", ref, idx, val))
			case anewarray:
				c := class.GetClass(inst.index())
				t := getGoType(c)
				size, _ := stack.pop()
				v := nextVar(t.type_, true)
				code.WriteString(fmt.Sprintf("%s = make([]%s, %s)\n", v, t, size))
			case newarray:
				t := getGoType(arrayTypeCodes(int(inst.operands[0])))
				size, _ := stack.pop()
				v := nextVar(t.type_, true)
				code.WriteString(fmt.Sprintf("%s = make([]%s, %s)\n", v, t, size))
			case arraylength:
				ref, _ := stack.pop()
				v := nextVar("int32", false)
				set(v, "int32(len("+ref+"))")
			case iconst_m1, iconst_0, iconst_1, iconst_2, iconst_3, iconst_4, iconst_5:
				v := nextVar("int32", false)
				set(v, strconv.Itoa(int(inst.Op)-int(iconst_0)))
			case dconst_0, dconst_1:
				set(nextVar("float64", false), strconv.Itoa(int(inst.Op-dconst_0)))
			case bipush:
				v := nextVar("int32", false)
				set(v, strconv.Itoa(int(inst.operands[0])))
			case sipush:
				v := nextVar("int32", false)
				set(v, strconv.Itoa(inst.index()))
			case new_:
				set(nextVar("*java_lang_Object", false), "new_"+ValidateName(class.GetClass(inst.index()))+"()")
			case ldc:
				c, t := class.GetConstant(int(inst.operands[0]))
				switch t {
				case intJ:
					set(nextVar("int32", false), c)
				case floatJ:
					set(nextVar("float32", false), c)
				case "java/lang/String":
					set(nextVar("*java_lang_Object", false), "newString("+c+")")
				default:
					panic("unknown type: " + t)
				}
			case ldc2_w:
				c, t := class.GetConstant(inst.index())
				switch t {
				case longJ:
					set(nextVar("int64", false), c)
				case doubleJ:
					set(nextVar("float64", false), c)
				default:
					panic("unknown type: " + t)
				}
			case dup:
				v, t := stack.pop()
				stack.push(v, t)
				set(nextVar(t.type_, t.isArray), v)
			case pop:
				v, _ := stack.pop()
				set("_", v)
			case iinc:
				code.WriteString(fmt.Sprintf("iarg%d += %d\n", inst.operands[0], inst.operands[1]))
			case iadd, ladd, dadd:
				sym("+")
			case isub, lsub:
				sym("-")
			case imul, lmul, fmul, dmul:
				sym("*")
			case idiv, ldiv, ddiv:
				sym("/")
			case irem, lrem:
				sym("%")
			case ineg, lneg:
				v1, t := stack.pop()
				v := nextVar(t.type_, t.isArray)
				code.WriteString(fmt.Sprintf("%s = -%s\n", v, v1))
			case ishl, lshl:
				sym("<<")
			case ishr, lshr:
				sym(">>")
			case iushr:
				v2, _ := stack.pop()
				v1, t := stack.pop()
				v := nextVar(t.type_, t.isArray)
				code.WriteString(fmt.Sprintf("%s = int32(uint32(%s)>>uint32(%s))\n", v, v1, v2))
			case lushr:
				v2, _ := stack.pop()
				v1, t := stack.pop()
				v := nextVar(t.type_, t.isArray)
				code.WriteString(fmt.Sprintf("%s = int64(uint64(%s)>>uint64(%s))\n", v, v1, v2))
			case iand, land:
				sym("&")
			case ior, lor:
				sym("|")
			case ixor, lxor:
				sym("^")
			case i2l:
				cast("int64", "int64(%s)")
			case i2d, l2d:
				cast("float64", "float64(%s)")
			case f2i:
				cast("int32", "int32(%s)")
			case d2f:
				cast("float32", "float32(%s)")
			case i2b:
				cast("int32", "int32(int8(%s))")
			case i2c:
				cast("int32", "int32(uint16(%s))")
			case i2s:
				cast("int32", "int32(int16(%s))")
			case putfield:
				_, f, _ := class.GetFieldRef(inst.index())
				val, _ := stack.pop()
				ref, _ := stack.pop()
				code.WriteString(fmt.Sprintf("%s.setField(\"%s\", %s)\n", ref, "E_"+f, val))
			case getfield:
				ref, _ := stack.pop()
				_, f, t := class.GetFieldRef(inst.index())
				g := getGoType(getJavaType(t))
				var callMethod = "getField"
				switch g.type_ {
				case "float64":
					callMethod = "getFieldDouble"
				case "int32":
					callMethod = "getFieldInt"
				default:
					callMethod = "getFieldObject"
				}
				v := nextVar(g.type_, g.isArray)
				code.WriteString(fmt.Sprintf("%s = %s.%s(\"%s\")\n", v, ref, callMethod, "E_"+f))
			case putstatic:
				panic("not implemented")
			case getstatic:
				n, f, t := class.GetFieldRef(inst.index())
				g := getGoType(t)
				set(nextVar(g.type_, g.isArray), ValidateName(n+"_"+f))
			case invokevirtual:
				c, n, t := class.GetMethodRef(inst.index())
				p, ret := translateParams(t)
				n = translateMethodName(n, c, ret, false, true, p)
				var params string
				for i := 0; i < len(p); i++ {
					v, _ := stack.pop()
					params = v + "," + params
				}
				rec, _ := stack.pop()
				var callMethod string
				var v = "_"
				switch ret {
				case intJ:
					callMethod = "callMethodInt"
					v = nextVar("int32", false)
				case doubleJ:
					callMethod = "callMethodDouble"
					v = nextVar("float64", false)
				case voidJ:
					callMethod = "callMethod"
				default:
					callMethod = "callMethodObject"
					v = nextVar("*java_lang_Object", false)
				}
				code.WriteString(fmt.Sprintf("%s = %s.%s(\"%s\", %s)\n", v, rec, callMethod, ValidateName(n), params))
			case invokespecial:
				c, n, t := class.GetMethodRef(inst.index())
				if n == "<init>" {
					n = "init"
				}
				p, ret := translateParams(t)
				n = translateMethodName(n, c, ret, false, true, p)
				var params string
				for i := 0; i < len(p); i++ {
					v, _ := stack.pop()
					params = v + "," + params
				}
				rec, _ := stack.pop()
				if n == "init" { // TODO: determine if should use super
					rec += ".super"
				}
				var callMethod string
				var v = "_"
				switch ret {
				case intJ:
					callMethod = "callMethodInt"
					v = nextVar("int32", false)
				case doubleJ:
					callMethod = "callMethodDouble"
					v = nextVar("float64", false)
				case voidJ:
					callMethod = "callMethod"
				default:
					callMethod = "callMethodObject"
					v = nextVar("*java_lang_Object", false)
				}
				code.WriteString(fmt.Sprintf("%s = %s.%s(\"%s\", %s)\n", v, rec, callMethod, ValidateName(n), params))
			case invokestatic:
				c, n, t := class.GetMethodRef(inst.index())
				p, ret := translateParams(t)
				n = translateMethodName(n, c, ret, true, true, p)
				var params string
				for i := 0; i < len(p); i++ {
					v, _ := stack.pop()
					params = v + "," + params
				}
				var v = "_"
				switch ret {
				case intJ:
					v = nextVar("int32", false)
				case doubleJ:
					v = nextVar("float64", false)
				case voidJ:
				default:
					v = nextVar("*java_lang_Object", false)
				}
				code.WriteString(fmt.Sprintf("%s = %s(%s)\n", v, ValidateName(n), params))
			case dcmpg:
				v2, _ := stack.pop()
				v1, _ := stack.pop()
				v := nextVar("int32", false)
				code.WriteString(fmt.Sprintf("%s = func(x, y float64) int32 {if x > y {return 1;} else if x == y {return 0;} else if x < y {return -1;}; return 1;}(%s, %s)\n", v, v1, v2))
			case goto_:
				code.WriteString(fmt.Sprintf("goto label%d\n", inst.index()+inst.Loc))
			case ifge, ifeq, ifne, iflt:
				v, _ := stack.pop()
				code.WriteString(fmt.Sprintf("if %s %s 0 { goto label%d } else { goto label%d }\n", v, compare[inst.Op], inst.index()+inst.Loc, blocks[cfg[id][1]].getLineStart()))
				vB.mark()
				stack.mark()
				continue nextBlock
			case if_icmpge, if_icmplt, if_icmpgt:
				v2, _ := stack.pop()
				v1, _ := stack.pop()
				code.WriteString(fmt.Sprintf("if %s %s %s { goto label%d } else { goto label%d }\n", v1, compare[inst.Op], v2, inst.index()+inst.Loc, blocks[cfg[id][1]].getLineStart()))
				vB.mark()
				stack.mark()
				continue nextBlock
			case lookupswitch:
				key, _ := stack.pop()
				code.WriteString("switch " + key + " {\n")
				defaultOffset := int(binary.BigEndian.Uint32(inst.operands))
				npairs := int(binary.BigEndian.Uint32(inst.operands[4:]))
				var i = 8
				for npairs > 0 {
					match := int(binary.BigEndian.Uint32(inst.operands[i:]))
					i += 4
					offset := int(binary.BigEndian.Uint32(inst.operands[i:]))
					i += 4
					npairs--
					code.WriteString(fmt.Sprintf("case %d: goto label%d\n", match, offset+inst.Loc))
				}
				code.WriteString(fmt.Sprintf("default: goto label%d\n}\n", defaultOffset+inst.Loc))
			case return_:
				code.WriteString("return\n")
			case ireturn, lreturn, areturn, dreturn:
				v, _ := stack.pop()
				code.WriteString("return " + v + "\n")
			default:
				panic(fmt.Sprintf("unknown opcode: %s", inst.Op))
			}
			last = inst
		}
		stack.reset()
		vB.reset()
		if last.Op == goto_ || last.Op == goto_w {
			continue
		}
		succs := cfg[id]
		if len(succs) <= 0 {
			continue
		}
		code.WriteString(fmt.Sprintf("goto label%d\n", blocks[succs[0]].getLineStart()))
	}
	if stack.length() != 0 {
		panic("stack isn't empty")
	}
}

// creates an intermediate form of the code
//func createIntermediate(blocks []basicBlock, class parser.RawClass, params []string) error {
//	stack := make(stack, 0)
//	nextVar := getUniqueCounter("v")
//	for _, block := range blocks {
//		for i, inst := range block {
//			switch inst.Op {
//			case nop: //ignore
//				continue
//			case aload:
//				inst.Type = params[inst.operands[0]]
//				inst.Args = []string{"a" + localName + strconv.Itoa(int(inst.operands[0]))}
//				stack.push(inst.Args[0], inst.Type)
//			case iload:
//				inst.Type = intJ
//				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.operands[0]))}
//				stack.push(inst.Args[0], inst.Type)
//			case fload:
//				inst.Type = floatJ
//				inst.Args = []string{"f" + localName + strconv.Itoa(int(inst.operands[0]))}
//				stack.push(inst.Args[0], inst.Type)
//			case dload:
//				inst.Type = doubleJ
//				inst.Args = []string{"d" + localName + strconv.Itoa(int(inst.operands[0]))}
//				stack.push(inst.Args[0], inst.Type)
//			case aload_0, aload_1, aload_2, aload_3:
//				inst.Type = params[int(inst.Op-aload_0)]
//				inst.Args = []string{"a" + localName + strconv.Itoa(int(inst.Op-aload_0))}
//				stack.push(inst.Args[0], inst.Type)
//			case dload_0, dload_1, dload_2, dload_3:
//				inst.Type = doubleJ
//				inst.Args = []string{"d" + localName + strconv.Itoa(int(inst.Op-dload_0))}
//				stack.push(inst.Args[0], inst.Type)
//			case iload_0, iload_1, iload_2, iload_3:
//				inst.Type = intJ
//				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.Op-iload_0))}
//				stack.push(inst.Args[0], inst.Type)
//			case lload_0, lload_1, lload_2, lload_3:
//				inst.Type = longJ
//				inst.Args = []string{"l" + localName + strconv.Itoa(int(inst.Op-lload_0))}
//				stack.push(inst.Args[0], inst.Type)
//			case aaload:
//				v := nextVar()
//				idx, _ := stack.pop()
//				ref, _ := stack.pop()
//				inst.Dest = v
//				inst.Args = []string{ref, idx}
//				inst.Type = "*java_lang_Object"
//				stack.push(v, inst.Type)
//			case iaload:
//				idx, _ := stack.pop()
//				ref, _ := stack.pop()
//				stack.push(ref+"["+idx+"]", intJ)
//			case caload:
//				idx, _ := stack.pop()
//				ref, _ := stack.pop()
//				stack.push("int32("+ref+"["+idx+"])", intJ)
//			case astore_0, astore_1, astore_2, astore_3:
//				inst.Dest = "a" + localName + strconv.Itoa(int(inst.Op-astore_0))
//				inst.Value, inst.Type = stack.pop()
//				if int(inst.Op-astore_0) >= len(params) {
//					var p [4]string
//					copy(p[:], params)
//					params = p[:]
//				}
//				params[int(inst.Op-astore_0)] = inst.Type
//			case aastore, iastore, castore:
//				val, _ := stack.pop()
//				index, _ := stack.pop()
//				ref, _ := stack.pop()
//				inst.Args = []string{ref, index, val}
//			case istore, dstore, astore, fstore:
//				var prefix string
//				switch inst.Op {
//				case istore:
//					prefix = "i"
//				case dstore:
//					prefix = "d"
//				case astore:
//					prefix = "a"
//				case fstore:
//					prefix = "f"
//				}
//				inst.Dest = prefix + localName + strconv.Itoa(int(inst.operands[0]))
//				inst.Value, inst.Type = stack.pop()
//				if int(inst.operands[0]) >= len(params) {
//					var p = make([]string, int(inst.operands[0])+1)
//					copy(p, params)
//					params = p
//				}
//				params[int(inst.operands[0])] = inst.Type
//			case istore_0, istore_1, istore_2, istore_3:
//				inst.Dest = "i" + localName + strconv.Itoa(int(inst.Op-istore_0))
//				inst.Value, inst.Type = stack.pop()
//			case lstore_0, lstore_1, lstore_2, lstore_3:
//				inst.Dest = "l" + localName + strconv.Itoa(int(inst.Op-lstore_0))
//				inst.Value, inst.Type = stack.pop()
//			case dstore_0, dstore_1, dstore_2, dstore_3:
//				inst.Dest = "d" + localName + strconv.Itoa(int(inst.Op-dstore_0))
//				inst.Value, inst.Type = stack.pop()
//			case invokespecial, invokevirtual, invokestatic:
//				c, n, t := class.GetMethodRef(inst.index())
//				p, _ := translateParams(t)
//				var params []string
//				for j := 0; j < len(p); j++ {
//					v, _ := stack.pop()
//					params = append([]string{v}, params...)
//				}
//				if inst.Op != invokestatic {
//					r, t := stack.pop()                     // the receiver
//					if inst.Op == invokespecial && c != t { // invokespecial calls inits so add the super class
//						r += ".super" // + ValidateName(c)
//					}
//					inst.HasReceiver = true
//					inst.Args = append([]string{r}, params...) // start with receiver
//				} else {
//					inst.Args = params
//				}
//				m := c + "." + n + ":" + t
//				inst.Func = n
//				inst.FDesc = m
//				if strings.HasSuffix(t, ")V") { //  ends in void
//					inst.Dest = "_"
//				} else {
//					v := nextVar()
//					inst.Dest = v
//					inst.Type = getJavaType(t[strings.LastIndex(t, ")")+1:])
//					stack.push(v, inst.Type)
//				}
//			case new_:
//				v := nextVar()
//				name := class.GetClass(inst.index())
//				inst.Type = name
//				inst.Dest = v
//				inst.Args = []string{name}
//				stack.push(v, inst.Type)
//			case dup:
//				s, t := stack.pop()
//				stack.push(s, t)
//				stack.push(s, t)
//			case i2d, l2d:
//				s, _ := stack.pop()
//				stack.push("float64("+s+")", doubleJ)
//			case i2s:
//				s, _ := stack.pop()
//				stack.push("int32(int16("+s+"))", intJ)
//			case i2l:
//				s, _ := stack.pop()
//				stack.push("int64("+s+")", longJ)
//			case d2i, f2i:
//				s, _ := stack.pop()
//				stack.push("int32("+s+")", intJ)
//			case i2c:
//				s, _ := stack.pop()
//				stack.push("int32(uint16("+s+"))", charJ)
//			case d2f:
//				s, _ := stack.pop()
//				stack.push("float32("+s+")", floatJ)
//			case i2b:
//				s, _ := stack.pop()
//				stack.push("int32(int8("+s+"))", intJ)
//			case iadd, dadd, ladd:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"+"+i2, t)
//			case isub, lsub:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"-"+i2, t)
//			case imul, dmul, lmul, fmul:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"*"+i2, t)
//			case idiv, ddiv, ldiv:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"/"+i2, t)
//			case irem, lrem:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"%"+i2, t)
//			case ineg, lneg:
//				i1, t := stack.pop()
//				stack.push("-"+i1, t)
//			case ishl, lshl:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"<<"+i2, t)
//			case ishr, lshr:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+">>"+i2, t)
//			case iushr:
//				i2, _ := stack.pop()
//				i1, t := stack.pop() // go uses logical shift when types are unsigned
//				stack.push(fmt.Sprintf("int32(uint32(%s)>>uint32(%s))", i1, i2), t)
//			case lushr:
//				i2, _ := stack.pop()
//				i1, t := stack.pop() // go uses logical shift when types are unsigned
//				stack.push(fmt.Sprintf("int64(uint64(%s)>>uint64(%s))", i1, i2), t)
//			case iand, land:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"&"+i2, t)
//			case ior, lor:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"|"+i2, t)
//			case ixor, lxor:
//				i2, _ := stack.pop()
//				i1, t := stack.pop()
//				stack.push(i1+"^"+i2, t)
//			case return_:
//			case ireturn, dreturn, areturn, lreturn:
//				inst.Args = make([]string, 1)
//				inst.Args[0], inst.Type = stack.pop()
//			case getstatic:
//				n, f, t := class.GetFieldRef(inst.index())
//				inst.Type = getJavaType(t) // cleans up the type if has L...;
//				inst.Args = []string{ValidateName(n) + "_" + f}
//				stack.push(inst.Args[0], inst.Type)
//			case ldc2_w:
//				var a = [1]string{}
//				inst.Args = (a)[:]
//				inst.Args[0], inst.Type = class.GetConstant(inst.index())
//				stack.push(inst.Args[0], inst.Type)
//			case ldc:
//				var v string
//				v, inst.Type = class.GetConstant(int(inst.operands[0]))
//				switch inst.Type { // do any necessary conversions for this constant
//				case "java/lang/String":
//					inst.Args = []string{fmt.Sprintf("newString(%s)", v)}
//				case intJ, floatJ:
//					inst.Args = []string{v}
//				default:
//					panic("unknown type: " + inst.Type)
//				}
//				stack.push(inst.Args[0], inst.Type)
//			case iconst_m1, iconst_0, iconst_1, iconst_2, iconst_3, iconst_4, iconst_5:
//				stack.push(strconv.Itoa(int(inst.Op)-int(iconst_0)), intJ)
//			case lconst_0, lconst_1:
//				stack.push(strconv.Itoa(int(inst.Op-lconst_0)), longJ)
//			case dconst_0, dconst_1:
//				stack.push(strconv.Itoa(int(inst.Op-dconst_0)), doubleJ)
//			case aconst_null:
//				stack.push("nil", objRefJ)
//			case bipush:
//				stack.push(strconv.Itoa(int(inst.operands[0])), intJ)
//			case sipush:
//				stack.push(strconv.Itoa(inst.index()), intJ)
//			case putfield:
//				v, _ := stack.pop()
//				ref, _ := stack.pop()
//				_, f, t := class.GetFieldRef(inst.index())
//				inst.Type = t
//				inst.Args = []string{ref, f, v}
//			case getfield:
//				v := nextVar()
//				ref, _ := stack.pop()
//				_, f, t := class.GetFieldRef(inst.index())
//				inst.Type = getJavaType(t) // confirm that is a java type
//				inst.Dest = v
//				inst.Args = []string{ref, f}
//				stack.push(v, inst.Type)
//			case pop:
//				ref, _ := stack.pop()
//				inst.Dest = "_"
//				inst.Args = []string{ref}
//			case ifge, ifne, ifgt, ifle, iflt, ifeq:
//				v, _ := stack.pop()
//				inst.Args = []string{v, strconv.Itoa(inst.index() + inst.Loc)}
//			case dcmpg:
//				v := nextVar()
//				a2, _ := stack.pop()
//				a1, _ := stack.pop()
//				inst.Dest = v
//				inst.Type = intJ
//				inst.Args = []string{a1, a2}
//				stack.push(v, inst.Type)
//			case if_icmpge, if_icmplt, if_icmpgt, if_icmple, if_icmpne:
//				v2, _ := stack.pop()
//				v1, _ := stack.pop()
//				inst.Args = []string{v1, v2, strconv.Itoa(inst.index() + inst.Loc)}
//			case anewarray:
//				v := nextVar()
//				size, _ := stack.pop()
//				c := class.GetClass(inst.index())
//				inst.Dest = v
//				inst.Type = "[]" + getGoType(c)
//				inst.Args = []string{size, c}
//				stack.push(v, inst.Type)
//			case newarray:
//				v := nextVar()
//				size, _ := stack.pop()
//				t := arrayTypeCodes(int(inst.operands[0]))
//				inst.Dest = v
//				inst.Type = "[]" + getGoType(t)
//				inst.Args = []string{size, t}
//				stack.push(v, inst.Type)
//			case arraylength:
//				ref, _ := stack.pop()
//				stack.push("int32(len("+ref+"))", inst.Type)
//			case goto_:
//				inst.Args = []string{strconv.Itoa(inst.index() + inst.Loc)}
//			case iinc:
//				//TODO: make sure this is always an int and not some other type
//				inst.Args = []string{"i" + localName + strconv.Itoa(int(inst.operands[0])), strconv.Itoa(int(inst.operands[1]))}
//			case lookupswitch:
//				key, _ := stack.pop()
//				defaultOffset := int(binary.BigEndian.Uint32(inst.operands))
//				npairs := int(binary.BigEndian.Uint32(inst.operands[4:]))
//				inst.Args = []string{key, strconv.Itoa(defaultOffset + inst.Loc)}
//				var i = 8 // starts at 8 bc defaultOffset and npairs
//				for npairs > 0 {
//					match := int(binary.BigEndian.Uint32(inst.operands[i:]))
//					i += 4
//					offset := int(binary.BigEndian.Uint32(inst.operands[i:]))
//					i += 4
//					inst.Args = append(inst.Args, strconv.Itoa(match), strconv.Itoa(offset+inst.Loc))
//					npairs--
//				}
//			default:
//				panic("unknown Op: " + inst.Op.String())
//			}
//			block[i] = inst // update the instruction
//		}
//	}
//	if len(stack) != 0 { // sanity check to make sure nothing went wrong
//		return errors.New(fmt.Sprintf("stack isn't empty: %s", stack))
//	}
//	return nil
//}
