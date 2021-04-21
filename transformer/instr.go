package transformer

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
)

type instruction struct {
	Loc      int // the instruction number
	Op       opcode
	operands []byte // raw additional bytes

	Type        string   // the java type of the new variable
	Dest        string   // the variable or field name
	Value       string   // if the inst is a const this is set
	Args        []string // the first arg is the receiver if it's a method
	Func        string   // Function name if is a function
	FDesc       string   // describes the function params and return
	HasReceiver bool     // does this method have a receiver
}

func (i instruction) index() int {
	return int(int16(i.operands[0])<<8 | int16(i.operands[1]))
}

func (i instruction) String() (s string) {
	if i.Type == "" && i.Func != "" && len(i.Dest) <= 0 {
		return "Op:" + i.Op.String() + " operands:" + fmt.Sprint(i.operands)
	}
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

type basicBlock []instruction

func (b basicBlock) getLineStart() int {
	return b[0].Loc
}

func (b basicBlock) getLast() instruction {
	var last instruction
	for i := len(b) - 1; i >= 0; i-- {
		last = b[i]
		if last.Op != nop {
			break
		}
	}
	return last
}

func getCFG(blocks []basicBlock, l2b map[int]int) map[int][]int {
	successors := make(map[int][]int) // makes names of blocks to a slice of names of blocks
	for idx, block := range blocks {
		var last = block.getLast()
		if last.Op == goto_ || last.Op == goto_w {
			successors[idx] = []int{l2b[getBrOffset(last)+last.Loc]}
		} else if hasBranch(last) {
			successors[idx] = []int{l2b[getBrOffset(last)+last.Loc], idx + 1}
		} else if isTerminator(last) {
			successors[idx] = []int{}
		} else { // fallthrough to next block
			if idx+1 >= len(blocks) { // no more blocks
				successors[idx] = []int{}
			} else {
				successors[idx] = []int{idx + 1}
			}
		}
	}
	return successors
}

func createBasicBlocks(instrs []instruction) ([]basicBlock, map[int]int) {
	add, get := func() (func(int), func() []int) {
		// if a line number is in this slice it means that it is the first instruction of a block
		//added 0 because a block starts at the beginning
		var startOfBlock = []int{0}
		var exists = make(map[int]struct{})
		return func(i int) {
				if _, ok := exists[i]; ok {
					return
				}
				startOfBlock = append(startOfBlock, i)
				exists[i] = struct{}{}
			}, func() []int {
				return startOfBlock
			}
	}()
	for _, v := range instrs {
		if v.Op == lookupswitch {
			defaultOffset := int(binary.BigEndian.Uint32(v.operands))
			add(defaultOffset + v.Loc)
			npairs := int(binary.BigEndian.Uint32(v.operands[4:]))
			var i = 8
			for npairs > 0 {
				i += 4 // skip over match
				offset := int(binary.BigEndian.Uint32(v.operands[i:]))
				add(offset + v.Loc)
				i += 4 // move to next int32
				npairs--
			}
			add(v.Loc + len(v.operands) + 1) // basic block at the next instruction
		} else if isTerminator(v) {
			add(v.Loc + len(v.operands) + 1)                     // basic block at the next instruction
			if hasBranch(v) || v.Op == goto_ || v.Op == goto_w { // a branch creates a new block at the branch location
				add(v.Loc + getBrOffset(v))
			}
		}
	}
	add(len(instrs)) // the last inst
	startOfBlock := get()
	// some instructions may jump backwards and therefore their starts are out of order
	sort.Ints(startOfBlock)
	blocks := make([]basicBlock, len(startOfBlock)-1) // number of blocks is one less size
	l2b := make(map[int]int)                          // match line numbers to block ids
	// loop through all but the last bc we add 1 to each one
	for i, v := range startOfBlock[:len(startOfBlock)-1] {
		// the next block is from this current start block # and the next one
		blocks[i] = instrs[v:startOfBlock[i+1]]
		l2b[v] = i // line v is the start of block i
	}
	return blocks, l2b
}

// readInstructions takes in a slice of bytes and breaks them up into a slice of useful instructions
// if an instruction has operands there will be nop instructions replaced in the slice so that
// jumps are easy to resolve
func readInstructions(b []byte) (instrs []instruction) {
	for i := 0; i < len(b); i++ {
		var instr = instruction{
			Loc: i,
			Op:  opcode(b[i]),
		}
		switch instr.Op {
		case wide, tableswitch:
			panic("not implemented") //TODO: implement
		}
		var operN = 0 // # of operands to read in
		switch instr.Op {
		case lookupswitch:
			var padding = 0 // padding gets placed at the end so that nops are added in the right place
			// zero to three bytes must act as padding, such that defaultbyte1 begins at an address that is a multiple of four bytes
			for i%4 != 0 {
				i++
				padding++
			}
			instr.operands = make([]byte, 4+4) // enough space for default and npairs
			copy(instr.operands, b[i:i+8])     // copies default and npairs into operands
			i += 8
			var npairs = binary.BigEndian.Uint32(instr.operands[4:8])           // number of int32 pairs
			var pairs = int(npairs) * (4 + 4)                                   // size of int32 pairs
			n := len(instr.operands)                                            // the location to start filling in pairs
			instr.operands = append(instr.operands, make([]byte, pairs)...)     // make enough room for pairs
			copy(instr.operands[n:], b[i:])                                     // copy in the pairs
			i += pairs - 1                                                      // move it to the last instruction of lookupswitch; for loop will increase by 1
			instr.operands = append(instr.operands, make([]byte, padding-1)...) // add padding
		case goto_w, jsr_w, invokedynamic, invokeinterface:
			// has 4 operands
			operN++
			fallthrough
		case multianewarray: // has three operands
			operN++
			fallthrough
		case goto_, if_acmpeq, if_acmpne, if_icmpeq, if_icmpge, if_icmpgt,
			if_icmple, if_icmplt, if_icmpne, ifeq, ifge, ifgt, ifle, iflt,
			ifne, ifnonnull, ifnull, jsr, sipush, iinc, anewarray, checkcast,
			getfield, getstatic, instanceof, invokespecial, invokestatic,
			invokevirtual, ldc_w, ldc2_w, new_, putfield, putstatic:
			// has two bytes
			operN++
			fallthrough
		case newarray, bipush, aload, astore, dload, dstore, fload, fstore,
			iload, istore, ldc, lload, lstore, ret:
			//has one operand
			operN++
			fallthrough
		default: // default is one byte so do nothing
			if operN > 0 {
				instr.operands = make([]byte, operN)
				copy(instr.operands, b[i+1:])
				i += operN
			}
		}
		instrs = append(instrs, instr)
		// add a nop instruction so that branches will go to the proper place
		instrs = append(instrs, make([]instruction, len(instr.operands))...)
	}
	return instrs
}

// getBrOffset returns the offset of a branch instruction.
// If a branch instruction is not sent the function behavior is undefined.
func getBrOffset(instr instruction) int {
	switch instr.Op {
	case goto_w, jsr_w:
		panic("not implemented") // TODO: implement
	}
	return int(int16(instr.operands[0])<<8 | int16(instr.operands[1]))
}

func hasBranch(instr instruction) bool {
	switch instr.Op {
	case jsr, jsr_w, if_acmpeq, if_acmpne, if_icmpeq, if_icmpge, if_icmpgt, if_icmple,
		if_icmplt, if_icmpne, ifeq, ifge, ifgt, ifle, iflt, ifne, ifnonnull, ifnull: //, goto_, goto_w:
		return true
	}
	return false
}

func isTerminator(instr instruction) bool {
	switch instr.Op {
	case jsr, jsr_w, if_acmpeq, if_acmpne, if_icmpeq, if_icmpge, if_icmpgt, if_icmple,
		if_icmplt, if_icmpne, ifeq, ifge, ifgt, ifle, iflt, ifne, ifnonnull, ifnull, goto_, goto_w,
		//invokeinterface, invokedynamic, invokevirtual, invokespecial, invokestatic,
		areturn, return_, ret, dreturn, freturn, ireturn, lreturn:
		return true
	}
	return false
}

//go:generate stringer -type=opcode
type opcode byte

const (
	nop             opcode = 0x00 //0000 0000		[No change]	perform no operation
	aconst_null     opcode = 0x01 //0000 0001		→ null	push a null reference onto the stack
	iconst_m1       opcode = 0x02 //0000 0010		→ -1	load the int value −1 onto the stack
	iconst_0        opcode = 0x03 //0000 0011		→ 0	load the int value 0 onto the stack
	iconst_1        opcode = 0x04 //0000 0100		→ 1	load the int value 1 onto the stack
	iconst_2        opcode = 0x05 //0000 0101		→ 2	load the int value 2 onto the stack
	iconst_3        opcode = 0x06 //0000 0110		→ 3	load the int value 3 onto the stack
	iconst_4        opcode = 0x07 //0000 0111		→ 4	load the int value 4 onto the stack
	iconst_5        opcode = 0x08 //0000 1000		→ 5	load the int value 5 onto the stack
	lconst_0        opcode = 0x09 //0000 1001		→ 0L	push 0L (the number zero with type long) onto the stack
	lconst_1        opcode = 0xa  //0000 1010		→ 1L	push 1L (the number one with type long) onto the stack
	fconst_0        opcode = 0xb  //0000 1011		→ 0.0f	push 0.0f on the stack
	fconst_1        opcode = 0xc  //0000 1100		→ 1.0f	push 1.0f on the stack
	fconst_2        opcode = 0xd  //0000 1101		→ 2.0f	push 2.0f on the stack
	dconst_0        opcode = 0xe  //0000 1110		→ 0.0	push the constant 0.0 (a double) onto the stack
	dconst_1        opcode = 0x0f //0000 1111		→ 1.0	push the constant 1.0 (a double) onto the stack
	bipush          opcode = 0x10 //0001 0000	1: byte	→ value	push a byte onto the stack as an integer value
	sipush          opcode = 0x11 //0001 0001	2: byte1, byte2	→ value	push a short onto the stack as an integer value
	ldc             opcode = 0x12 //0001 0010	1: index	→ value	push a constant #index from a constant pool (String, int, float, Class, java.lang.invoke.MethodType, java.lang.invoke.MethodHandle, or a dynamically-computed constant) onto the stack
	ldc_w           opcode = 0x13 //0001 0011	2: indexbyte1, indexbyte2	→ value	push a constant #index from a constant pool (String, int, float, Class, java.lang.invoke.MethodType, java.lang.invoke.MethodHandle, or a dynamically-computed constant) onto the stack (wide index is constructed as indexbyte1 << 8 | indexbyte2)
	ldc2_w          opcode = 0x14 //0001 0100	2: indexbyte1, indexbyte2	→ value	push a constant #index from a constant pool (double, long, or a dynamically-computed constant) onto the stack (wide index is constructed as indexbyte1 << 8 | indexbyte2)
	iload           opcode = 0x15 //0001 0101	1: index	→ value	load an int value from a local variable #index
	lload           opcode = 0x16 //0001 0110	1: index	→ value	load a long value from a local variable #index
	fload           opcode = 0x17 //0001 0111	1: index	→ value	load a float value from a local variable #index
	dload           opcode = 0x18 //0001 1000	1: index	→ value	load a double value from a local variable #index
	aload           opcode = 0x19 //0001 1001	1: index	→ objectref	load a reference onto the stack from a local variable #index
	iload_0         opcode = 0x1a //0001 1010		→ value	load an int value from local variable 0
	iload_1         opcode = 0x1b //0001 1011		→ value	load an int value from local variable 1
	iload_2         opcode = 0x1c //0001 1100		→ value	load an int value from local variable 2
	iload_3         opcode = 0x1d //0001 1101		→ value	load an int value from local variable 3
	lload_0         opcode = 0x1e //0001 1110		→ value	load a long value from a local variable 0
	lload_1         opcode = 0x1f //0001 1111		→ value	load a long value from a local variable 1
	lload_2         opcode = 0x20 //0010 0000		→ value	load a long value from a local variable 2
	lload_3         opcode = 0x21 //0010 0001		→ value	load a long value from a local variable 3
	fload_0         opcode = 0x22 //0010 0010		→ value	load a float value from local variable 0
	fload_1         opcode = 0x23 //0010 0011		→ value	load a float value from local variable 1
	fload_2         opcode = 0x24 //0010 0100		→ value	load a float value from local variable 2
	fload_3         opcode = 0x25 //0010 0101		→ value	load a float value from local variable 3
	dload_0         opcode = 0x26 //0010 0110		→ value	load a double from local variable 0
	dload_1         opcode = 0x27 //0010 0111		→ value	load a double from local variable 1
	dload_2         opcode = 0x28 //0010 1000		→ value	load a double from local variable 2
	dload_3         opcode = 0x29 //0010 1001		→ value	load a double from local variable 3
	aload_0         opcode = 0x2a //0010 1010		→ objectref	load a reference onto the stack from local variable 0
	aload_1         opcode = 0x2b //0010 1011		→ objectref	load a reference onto the stack from local variable 1
	aload_2         opcode = 0x2c //0010 1100		→ objectref	load a reference onto the stack from local variable 2
	aload_3         opcode = 0x2d //0010 1101		→ objectref	load a reference onto the stack from local variable 3
	iaload          opcode = 0x2e //0010 1110		arrayref, index → value	load an int from an array
	laload          opcode = 0x2f //0010 1111		arrayref, index → value	load a long from an array
	faload          opcode = 0x30 //0011 0000		arrayref, index → value	load a float from an array
	daload          opcode = 0x31 //0011 0001		arrayref, index → value	load a double from an array
	aaload          opcode = 0x32 //0011 0010		arrayref, index → value	load onto the stack a reference from an array
	baload          opcode = 0x33 //0011 0011		arrayref, index → value	load a byte or Boolean value from an array
	caload          opcode = 0x34 //0011 0100		arrayref, index → value	load a char from an array
	saload          opcode = 0x35 //0011 0101		arrayref, index → value	load short from array
	istore          opcode = 0x36 //0011 0110	1: index	value →	store int value into variable #index
	lstore          opcode = 0x37 //0011 0111	1: index	value →	store a long value in a local variable #index
	fstore          opcode = 0x38 //0011 1000	1: index	value →	store a float value into a local variable #index
	dstore          opcode = 0x39 //0011 1001	1: index	value →	store a double value into a local variable #index
	astore          opcode = 0x3a //0011 1010	1: index	objectref →	store a reference into a local variable #index
	istore_0        opcode = 0x3b //0011 1011		value →	store int value into variable 0
	istore_1        opcode = 0x3c //0011 1100		value →	store int value into variable 1
	istore_2        opcode = 0x3d //0011 1101		value →	store int value into variable 2
	istore_3        opcode = 0x3e //0011 1110		value →	store int value into variable 3
	lstore_0        opcode = 0x3f //0011 1111		value →	store a long value in a local variable 0
	lstore_1        opcode = 0x40 //0100 0000		value →	store a long value in a local variable 1
	lstore_2        opcode = 0x41 //0100 0001		value →	store a long value in a local variable 2
	lstore_3        opcode = 0x42 //0100 0010		value →	store a long value in a local variable 3
	fstore_0        opcode = 0x43 //0100 0011		value →	store a float value into local variable 0
	fstore_1        opcode = 0x44 //0100 0100		value →	store a float value into local variable 1
	fstore_2        opcode = 0x45 //0100 0101		value →	store a float value into local variable 2
	fstore_3        opcode = 0x46 //0100 0110		value →	store a float value into local variable 3
	dstore_0        opcode = 0x47 //0100 0111		value →	store a double into local variable 0
	dstore_1        opcode = 0x48 //0100 1000		value →	store a double into local variable 1
	dstore_2        opcode = 0x49 //0100 1001		value →	store a double into local variable 2
	dstore_3        opcode = 0x4a //0100 1010		value →	store a double into local variable 3
	astore_0        opcode = 0x4b //0100 1011		objectref →	store a reference into local variable 0
	astore_1        opcode = 0x4c //0100 1100		objectref →	store a reference into local variable 1
	astore_2        opcode = 0x4d //0100 1101		objectref →	store a reference into local variable 2
	astore_3        opcode = 0x4e //0100 1110		objectref →	store a reference into local variable 3
	iastore         opcode = 0x4f //0100 1111		arrayref, index, value →	store an int into an array
	lastore         opcode = 0x50 //0101 0000		arrayref, index, value →	store a long to an array
	fastore         opcode = 0x51 //0101 0001		arrayref, index, value →	store a float in an array
	dastore         opcode = 0x52 //0101 0010		arrayref, index, value →	store a double into an array
	aastore         opcode = 0x53 //0101 0011		arrayref, index, value →	store a reference in an array
	bastore         opcode = 0x54 //0101 0100		arrayref, index, value →	store a byte or Boolean value into an array
	castore         opcode = 0x55 //0101 0101		arrayref, index, value →	store a char into an array
	sastore         opcode = 0x56 //0101 0110		arrayref, index, value →	store short to array
	pop             opcode = 0x57 //0101 0111		value →	discard the top value on the stack
	pop2            opcode = 0x58 //0101 1000		{value2, value1} →	discard the top two values on the stack (or one value, if it is a double or long)
	dup             opcode = 0x59 //0101 1001		value → value, value	duplicate the value on top of the stack
	dup_x1          opcode = 0x5a //0101 1010		value2, value1 → value1, value2, value1	insert a copy of the top value into the stack two values from the top. value1 and value2 must not be of the type double or long.
	dup_x2          opcode = 0x5b //0101 1011		value3, value2, value1 → value1, value3, value2, value1	insert a copy of the top value into the stack two (if value2 is double or long it takes up the entry of value3, too) or three values (if value2 is neither double nor long) from the top
	dup2            opcode = 0x5c //0101 1100		{value2, value1} → {value2, value1}, {value2, value1}	duplicate top two stack words (two values, if value1 is not double nor long; a single value, if value1 is double or long)
	dup2_x1         opcode = 0x5d //0101 1101		value3, {value2, value1} → {value2, value1}, value3, {value2, value1}	duplicate two words and insert beneath third word (see explanation above)
	dup2_x2         opcode = 0x5e //0101 1110		{value4, value3}, {value2, value1} → {value2, value1}, {value4, value3}, {value2, value1}	duplicate two words and insert beneath fourth word
	swap            opcode = 0x5f //0101 1111		value2, value1 → value1, value2	swaps two top words on the stack (note that value1 and value2 must not be double or long)
	iadd            opcode = 0x60 //0110 0000		value1, value2 → result	add two ints
	ladd            opcode = 0x61 //0110 0001		value1, value2 → result	add two longs
	fadd            opcode = 0x62 //0110 0010		value1, value2 → result	add two floats
	dadd            opcode = 0x63 //0110 0011		value1, value2 → result	add two doubles
	isub            opcode = 0x64 //0110 0100		value1, value2 → result	int subtract
	lsub            opcode = 0x65 //0110 0101		value1, value2 → result	subtract two longs
	fsub            opcode = 0x66 //0110 0110		value1, value2 → result	subtract two floats
	dsub            opcode = 0x67 //0110 0111		value1, value2 → result	subtract a double from another
	imul            opcode = 0x68 //0110 1000		value1, value2 → result	multiply two integers
	lmul            opcode = 0x69 //0110 1001		value1, value2 → result	multiply two longs
	fmul            opcode = 0x6a //0110 1010		value1, value2 → result	multiply two floats
	dmul            opcode = 0x6b //0110 1011		value1, value2 → result	multiply two doubles
	idiv            opcode = 0x6c //0110 1100		value1, value2 → result	divide two integers
	ldiv            opcode = 0x6d //0110 1101		value1, value2 → result	divide two longs
	fdiv            opcode = 0x6e //0110 1110		value1, value2 → result	divide two floats
	ddiv            opcode = 0x6f //0110 1111		value1, value2 → result	divide two doubles
	irem            opcode = 0x70 //0111 0000		value1, value2 → result	logical int remainder
	lrem            opcode = 0x71 //0111 0001		value1, value2 → result	remainder of division of two longs
	frem            opcode = 0x72 //0111 0010		value1, value2 → result	get the remainder from a division between two floats
	drem            opcode = 0x73 //0111 0011		value1, value2 → result	get the remainder from a division between two doubles
	ineg            opcode = 0x74 //0111 0100		value → result	negate int
	lneg            opcode = 0x75 //0111 0101		value → result	negate a long
	fneg            opcode = 0x76 //0111 0110		value → result	negate a float
	dneg            opcode = 0x77 //0111 0111		value → result	negate a double
	ishl            opcode = 0x78 //0111 1000		value1, value2 → result	int shift left
	lshl            opcode = 0x79 //0111 1001		value1, value2 → result	bitwise shift left of a long value1 by int value2 positions
	ishr            opcode = 0x7a //0111 1010		value1, value2 → result	int arithmetic shift right
	lshr            opcode = 0x7b //0111 1011		value1, value2 → result	bitwise shift right of a long value1 by int value2 positions
	iushr           opcode = 0x7c //0111 1100		value1, value2 → result	int logical shift right
	lushr           opcode = 0x7d //0111 1101		value1, value2 → result	bitwise shift right of a long value1 by int value2 positions, unsigned
	iand            opcode = 0x7e //0111 1110		value1, value2 → result	perform a bitwise AND on two integers
	land            opcode = 0x7f //0111 1111		value1, value2 → result	bitwise AND of two longs
	ior             opcode = 0x80 //1000 0000		value1, value2 → result	bitwise int OR
	lor             opcode = 0x81 //1000 0001		value1, value2 → result	bitwise OR of two longs
	ixor            opcode = 0x82 //1000 0010		value1, value2 → result	int xor
	lxor            opcode = 0x83 //1000 0011		value1, value2 → result	bitwise XOR of two longs
	iinc            opcode = 0x84 //1000 0100	2: index, const	[No change]	increment local variable #index by signed byte const
	i2l             opcode = 0x85 //1000 0101		value → result	convert an int into a long
	i2f             opcode = 0x86 //1000 0110		value → result	convert an int into a float
	i2d             opcode = 0x87 //1000 0111		value → result	convert an int into a double
	l2i             opcode = 0x88 //1000 1000		value → result	convert a long to a int
	l2f             opcode = 0x89 //1000 1001		value → result	convert a long to a float
	l2d             opcode = 0x8a //1000 1010		value → result	convert a long to a double
	f2i             opcode = 0x8b //1000 1011		value → result	convert a float to an int
	f2l             opcode = 0x8c //1000 1100		value → result	convert a float to a long
	f2d             opcode = 0x8d //1000 1101		value → result	convert a float to a double
	d2i             opcode = 0x8e //1000 1110		value → result	convert a double to an int
	d2l             opcode = 0x8f //1000 1111		value → result	convert a double to a long
	d2f             opcode = 0x90 //1001 0000		value → result	convert a double to a float
	i2b             opcode = 0x91 //1001 0001		value → result	convert an int into a byte
	i2c             opcode = 0x92 //1001 0010		value → result	convert an int into a character
	i2s             opcode = 0x93 //1001 0011		value → result	convert an int into a short
	lcmp            opcode = 0x94 //1001 0100		value1, value2 → result	push 0 if the two longs are the same, 1 if value1 is greater than value2, -1 otherwise
	fcmpl           opcode = 0x95 //1001 0101		value1, value2 → result	compare two floats, -1 on NaN
	fcmpg           opcode = 0x96 //1001 0110		value1, value2 → result	compare two floats, 1 on NaN
	dcmpl           opcode = 0x97 //1001 0111		value1, value2 → result	compare two doubles, -1 on NaN
	dcmpg           opcode = 0x98 //1001 1000		value1, value2 → result	compare two doubles, 1 on NaN
	ifeq            opcode = 0x99 //1001 1001	2: branchbyte1, branchbyte2	value →	if value is 0, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	ifne            opcode = 0x9a //1001 1010	2: branchbyte1, branchbyte2	value →	if value is not 0, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	iflt            opcode = 0x9b //1001 1011	2: branchbyte1, branchbyte2	value →	if value is less than 0, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	ifge            opcode = 0x9c //1001 1100	2: branchbyte1, branchbyte2	value →	if value is greater than or equal to 0, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	ifgt            opcode = 0x9d //1001 1101	2: branchbyte1, branchbyte2	value →	if value is greater than 0, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	ifle            opcode = 0x9e //1001 1110	2: branchbyte1, branchbyte2	value →	if value is less than or equal to 0, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_icmpeq       opcode = 0x9f //1001 1111	2: branchbyte1, branchbyte2	value1, value2 →	if ints are equal, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_icmpne       opcode = 0xa0 //1010 0000	2: branchbyte1, branchbyte2	value1, value2 →	if ints are not equal, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_icmplt       opcode = 0xa1 //1010 0001	2: branchbyte1, branchbyte2	value1, value2 →	if value1 is less than value2, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_icmpge       opcode = 0xa2 //1010 0010	2: branchbyte1, branchbyte2	value1, value2 →	if value1 is greater than or equal to value2, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_icmpgt       opcode = 0xa3 //1010 0011	2: branchbyte1, branchbyte2	value1, value2 →	if value1 is greater than value2, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_icmple       opcode = 0xa4 //1010 0100	2: branchbyte1, branchbyte2	value1, value2 →	if value1 is less than or equal to value2, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_acmpeq       opcode = 0xa5 //1010 0101	2: branchbyte1, branchbyte2	value1, value2 →	if references are equal, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	if_acmpne       opcode = 0xa6 //1010 0110	2: branchbyte1, branchbyte2	value1, value2 →	if references are not equal, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	goto_           opcode = 0xa7 //1010 0111	2: branchbyte1, branchbyte2	[no change]	goes to another instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	jsr             opcode = 0xa8 //1010 1000	2: branchbyte1, branchbyte2	→ address	jump to subroutine at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2) and place the return address on the stack
	ret             opcode = 0xa9 //1010 1001	1: index	[No change]	continue execution from address taken from a local variable #index (the asymmetry with jsr is intentional)
	tableswitch     opcode = 0xaa //1010 1010	16+: [0–3 bytes padding], defaultbyte1, defaultbyte2, defaultbyte3, defaultbyte4, lowbyte1, lowbyte2, lowbyte3, lowbyte4, highbyte1, highbyte2, highbyte3, highbyte4, jump offsets...	index →	continue execution from an address in the table at offset index
	lookupswitch    opcode = 0xab //1010 1011	8+: <0–3 bytes padding>, defaultbyte1, defaultbyte2, defaultbyte3, defaultbyte4, npairs1, npairs2, npairs3, npairs4, match-offset pairs...	key →	a target address is looked up from a table using a key and execution continues from the instruction at that address
	ireturn         opcode = 0xac //1010 1100		value → [empty]	return an integer from a method
	lreturn         opcode = 0xad //1010 1101		value → [empty]	return a long value
	freturn         opcode = 0xae //1010 1110		value → [empty]	return a float
	dreturn         opcode = 0xaf //1010 1111		value → [empty]	return a double from a method
	areturn         opcode = 0xb0 //1011 0000		objectref → [empty]	return a reference from a method
	return_         opcode = 0xb1 //1011 0001		→ [empty]	return void from method
	getstatic       opcode = 0xb2 //1011 0010	2: indexbyte1, indexbyte2	→ value	get a static field value of a class, where the field is identified by field reference in the constant pool index (indexbyte1 << 8 | indexbyte2)
	putstatic       opcode = 0xb3 //1011 0011	2: indexbyte1, indexbyte2	value →	set static field to value in a class, where the field is identified by a field reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	getfield        opcode = 0xb4 //1011 0100	2: indexbyte1, indexbyte2	objectref → value	get a field value of an object objectref, where the field is identified by field reference in the constant pool index (indexbyte1 << 8 | indexbyte2)
	putfield        opcode = 0xb5 //1011 0101	2: indexbyte1, indexbyte2	objectref, value →	set field to value in an object objectref, where the field is identified by a field reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	invokevirtual   opcode = 0xb6 //1011 0110	2: indexbyte1, indexbyte2	objectref, [arg1, arg2, ...] → result	invoke virtual method on object objectref and puts the result on the stack (might be void); the method is identified by method reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	invokespecial   opcode = 0xb7 //1011 0111	2: indexbyte1, indexbyte2	objectref, [arg1, arg2, ...] → result	invoke instance method on object objectref and puts the result on the stack (might be void); the method is identified by method reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	invokestatic    opcode = 0xb8 //1011 1000	2: indexbyte1, indexbyte2	[arg1, arg2, ...] → result	invoke a static method and puts the result on the stack (might be void); the method is identified by method reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	invokeinterface opcode = 0xb9 //1011 1001	4: indexbyte1, indexbyte2, count, 0	objectref, [arg1, arg2, ...] → result	invokes an interface method on object objectref and puts the result on the stack (might be void); the interface method is identified by method reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	invokedynamic   opcode = 0xba //1011 1010	4: indexbyte1, indexbyte2, 0, 0	[arg1, [arg2 ...]] → result	invokes a dynamic method and puts the result on the stack (might be void); the method is identified by method reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	new_            opcode = 0xbb //1011 1011	2: indexbyte1, indexbyte2	→ objectref	create new_ object of type identified by class reference in constant pool index (indexbyte1 << 8 | indexbyte2)
	newarray        opcode = 0xbc //1011 1100	1: atype	count → arrayref	create new array with count elements of primitive type identified by atype
	anewarray       opcode = 0xbd //1011 1101	2: indexbyte1, indexbyte2	count → arrayref	create a new array of references of length count and component type identified by the class reference index (indexbyte1 << 8 | indexbyte2) in the constant pool
	arraylength     opcode = 0xbe //1011 1110		arrayref → length	get the length of an array
	athrow          opcode = 0xbf //1011 1111		objectref → [empty], objectref	throws an error or exception (notice that the rest of the stack is cleared, leaving only a reference to the Throwable)
	checkcast       opcode = 0xc0 //1100 0000	2: indexbyte1, indexbyte2	objectref → objectref	checks whether an objectref is of a certain type, the class reference of which is in the constant pool at index (indexbyte1 << 8 | indexbyte2)
	instanceof      opcode = 0xc1 //1100 0001	2: indexbyte1, indexbyte2	objectref → result	determines if an object objectref is of a given type, identified by class reference index in constant pool (indexbyte1 << 8 | indexbyte2)
	monitorenter    opcode = 0xc2 //1100 0010		objectref →	enter monitor for object ("grab the lock" – start of synchronized() section)
	monitorexit     opcode = 0xc3 //1100 0011		objectref →	exit monitor for object ("release the lock" – end of synchronized() section)
	wide            opcode = 0xc4 //1100 0100	3/5: opcode, indexbyte1, indexbyte2 or iinc, indexbyte1, indexbyte2, countbyte1, countbyte2	[same as for corresponding instructions]	execute opcode, where opcode is either iload, fload, aload, lload, dload, istore, fstore, astore, lstore, dstore, or ret, but assume the index is 16 bit; or execute iinc, where the index is 16 bits and the constant to increment by is a signed 16 bit short
	multianewarray  opcode = 0xc5 //1100 0101	3: indexbyte1, indexbyte2, dimensions	count1, [count2,...] → arrayref	create a new array of dimensions dimensions of type identified by class reference in constant pool index (indexbyte1 << 8 | indexbyte2); the sizes of each dimension is identified by count1, [count2, etc.]
	ifnull          opcode = 0xc6 //1100 0110	2: branchbyte1, branchbyte2	value →	if value is null, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	ifnonnull       opcode = 0xc7 //1100 0111	2: branchbyte1, branchbyte2	value →	if value is not null, branch to instruction at branchoffset (signed short constructed from unsigned bytes branchbyte1 << 8 | branchbyte2)
	goto_w          opcode = 0xc8 //1100 1000	4: branchbyte1, branchbyte2, branchbyte3, branchbyte4	[no change]	goes to another instruction at branchoffset (signed int constructed from unsigned bytes branchbyte1 << 24 | branchbyte2 << 16 | branchbyte3 << 8 | branchbyte4)
	jsr_w           opcode = 0xc9 //1100 1001	4: branchbyte1, branchbyte2, branchbyte3, branchbyte4	→ address	jump to subroutine at branchoffset (signed int constructed from unsigned bytes branchbyte1 << 24 | branchbyte2 << 16 | branchbyte3 << 8 | branchbyte4) and place the return address on the stack
	breakpoint      opcode = 0xca //1100 1010			reserved for breakpoints in Java debuggers; should not appear in any class file
	//(no name)	cb-fd				these values are currently unassigned for opcodes and are reserved for future use
	_       opcode = 0xCB
	_       opcode = 0xCC
	_       opcode = 0xCD
	_       opcode = 0xCE
	_       opcode = 0xCF
	_       opcode = 0xD0
	_       opcode = 0xD1
	_       opcode = 0xD2
	_       opcode = 0xD3
	_       opcode = 0xD4
	_       opcode = 0xD5
	_       opcode = 0xD6
	_       opcode = 0xD7
	_       opcode = 0xD8
	_       opcode = 0xD9
	_       opcode = 0xDA
	_       opcode = 0xDB
	_       opcode = 0xDC
	_       opcode = 0xDD
	_       opcode = 0xDE
	_       opcode = 0xDF
	_       opcode = 0xE0
	_       opcode = 0xE1
	_       opcode = 0xE2
	_       opcode = 0xE3
	_       opcode = 0xE4
	_       opcode = 0xE5
	_       opcode = 0xE6
	_       opcode = 0xE7
	_       opcode = 0xE8
	_       opcode = 0xE9
	_       opcode = 0xEA
	_       opcode = 0xEB
	_       opcode = 0xEC
	_       opcode = 0xED
	_       opcode = 0xEE
	_       opcode = 0xEF
	_       opcode = 0xF0
	_       opcode = 0xF1
	_       opcode = 0xF2
	_       opcode = 0xF3
	_       opcode = 0xF4
	_       opcode = 0xF5
	_       opcode = 0xF6
	_       opcode = 0xF7
	_       opcode = 0xF8
	_       opcode = 0xF9
	_       opcode = 0xFA
	_       opcode = 0xFB
	_       opcode = 0xFC
	_       opcode = 0xFD
	impdep1 opcode = 0xfe //1111 1110			reserved for implementation-dependent operations within debuggers; should not appear in any class file
	impdep2 opcode = 0xff //1111 1111			reserved for implementation-dependent operations within debuggers; should not appear in any class file
)
