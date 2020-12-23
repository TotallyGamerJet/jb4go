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

type stack []string

func (s *stack) push(v string) {
	*s = append(*s, v)
}

func (s *stack) pop() string {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

// creates an intermediate form of the code
func createIntermediate(blocks []basicBlock, class parser.RawClass) string {
	stack := make(stack, 0)
	nextVar := getUniqueCounter("v")
	var inter = []string{prelude}
	for _, block := range blocks {
		inter = append(inter, "case "+strconv.Itoa(block[0].loc)+":")
		for _, inst := range block {
			switch inst.opcode {
			case nop: //ignore
			case aload_0, aload_1, aload_2, aload_3:
				v := nextVar()
				inter = append(inter, "var "+v+" = "+localName+strconv.Itoa(int(inst.opcode-aload_0)))
				stack.push(v)
			case iload:
				v := nextVar()
				inter = append(inter, "var "+v+" = "+localName+strconv.Itoa(int(inst.operands[0])))
				stack.push(v)
			case iload_0, iload_1, iload_2, iload_3:
				v := nextVar()
				inter = append(inter, "var "+v+" = "+localName+strconv.Itoa(int(inst.opcode-iload_0)))
				stack.push(v)
			case astore_0, astore_1, astore_2, astore_3:
				inter = append(inter, localName+strconv.Itoa(int(inst.opcode-astore_0))+" = "+stack.pop())
			case istore:
				inter = append(inter, localName+strconv.Itoa(int(inst.operands[0]))+" = "+stack.pop())
			case istore_0, istore_1, istore_2, istore_3:
				inter = append(inter, localName+strconv.Itoa(int(inst.opcode-istore_0))+" = "+stack.pop())
			case invokespecial, invokevirtual:
				c, n, t := class.GetMethodRef(inst.index())
				p, _ := translateParams(t)
				var params string
				for j := 0; j < len(p); j++ {
					params += stack.pop() + ","
				}
				r := stack.pop() // the receiver
				m := c + "." + n + ":" + t
				if strings.HasSuffix(t, ")V") { //  ends in void
					inter = append(inter, "    _ = "+r+"."+n+"("+params+") //["+m+"]")
				} else {
					v := nextVar()
					inter = append(inter, "var "+v+" = "+r+"."+n+"("+params+") //["+m+"]")
					stack.push(v)
				}
			case invokestatic:
				c, n, t := class.GetMethodRef(inst.index())
				p, _ := translateParams(t)
				var params string
				for j := 0; j < len(p); j++ {
					params += stack.pop() + ","
				}
				m := c + "." + n + ":" + t
				if strings.HasSuffix(t, ")V") { //  ends in void
					inter = append(inter, "    _ = "+c+"."+n+"("+params+") //["+m+"]")
				} else {
					v := nextVar()
					inter = append(inter, "var "+v+" = "+c+"."+n+"("+params+") //["+m+"]")
					stack.push(v)
				}
			case new_:
				v := nextVar()
				name := class.GetClass(inst.index())
				inter = append(inter, "var "+v+" = "+fmt.Sprintf("new(%s)", name))
				stack.push(v)
			case dup:
				v := nextVar()
				s := stack.pop()
				stack.push(s)
				stack.push(v)
				inter = append(inter, "var "+v+" = "+s)
			case irem:
				v := nextVar()
				i2 := stack.pop()
				i1 := stack.pop()
				inter = append(inter, fmt.Sprintf("var %s = %s %% %s", v, i1, i2))
				stack.push(v)
			case return_:
				inter = append(inter, "return")
			case ireturn, areturn:
				inter = append(inter, "return "+stack.pop())
			case getstatic:
				v := nextVar()
				n, f, t := class.GetFieldRef(inst.index())
				inter = append(inter, "var "+v+" = "+fmt.Sprintf("%s.%s:%s", n, f, t))
				stack.push(v)
			case ldc:
				v := nextVar()
				c := class.GetConstant(int(inst.operands[0]))
				inter = append(inter, "var "+v+" = "+c)
				stack.push(v)
			case putfield:
				val := stack.pop()
				ref := stack.pop()
				n, f, t := class.GetFieldRef(inst.index())
				inter = append(inter, fmt.Sprintf("%s.%s = %s //%s:%s", ref, f, val, n, t))
			case getfield:
				v := nextVar()
				ref := stack.pop()
				n, f, t := class.GetFieldRef(inst.index())
				inter = append(inter, fmt.Sprintf("var %s = %s.%s //%s:%s", v, ref, f, n, t))
				stack.push(v)
			case ifge:
				v := stack.pop()
				inter = append(inter, fmt.Sprintf("if %s >= 0 { currentLabel = %d; continue controlFlowLoop }", v, inst.index()+inst.loc))
			case ifne:
				v := stack.pop()
				inter = append(inter, fmt.Sprintf("if %s != 0 { currentLabel = %d; continue controlFlowLoop }", v, inst.index()+inst.loc))
			case goto_:
				inter = append(inter, fmt.Sprintf("currentLabel = %d; continue controlFlowLoop", inst.index()+inst.loc))
			default:
				panic("unknown opcode: " + inst.opcode.String())
			}
			//fmt.Println("intermediate:")
			//for _, v := range inter {
			//	fmt.Printf("\t%s\n", v)
			//}
		}
		inter = append(inter, "fallthrough") // continue to the next block if it didn't already jump
	}
	inter = append(inter, envoi)
	if len(stack) != 0 { // sanity check to make sure nothing went wrong
		panic("stack isn't empty")
	}
	b := strings.Builder{}
	for _, s := range inter {
		b.WriteString(s + "\n")
	}
	fmt.Println(b.String())
	return b.String()
}
