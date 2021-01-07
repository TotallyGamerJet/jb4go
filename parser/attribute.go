package parser

import "fmt"

type AttributeInfo interface{}

func ReadAttributeInfo(c *RawClass, p *Parser) (info AttributeInfo) {
	b := basic{
		attributeNameIndex: p.ReadU2(),
		attributeLength:    p.ReadU4(),
	}
	tag, ok := c.GetCP(int(b.attributeNameIndex)).(utf8Info)
	if !ok {
		panic("improper attribute, must point to utf8")
	}
	switch tag.str {
	case "Code":
		var code codeAttribute
		code.basic = b
		code.maxStack = p.ReadU2()
		code.maxLocals = p.ReadU2()
		code.codeLength = p.ReadU4()
		code.code = make([]byte, code.codeLength)
		p.Read(code.code)
		code.exceptionTableLength = p.ReadU2()
		code.exceptionTable = make([]exception, code.exceptionTableLength)
		for i := range code.exceptionTable {
			code.exceptionTable[i] = readException(p)
		}
		code.attributesCount = p.ReadU2()
		code.attributes = make([]AttributeInfo, code.attributesCount)
		for i := range code.attributes {
			code.attributes[i] = ReadAttributeInfo(c, p)
		}
		info = code
	case "LineNumberTable":
		size := p.ReadU2()
		table := make([]lineNumberTableEntry, size)
		for i := range table {
			table[i] = lineNumberTableEntry{startPc: p.ReadU2(), lineNumber: p.ReadU2()}
		}
		info = lineNumberTable{b, size, table}
	case "InnerClasses":
		var c innerClassesAttribute
		c.basic = b
		c.numberOfClasses = p.ReadU2()
		c.classes = make([]classAttrib, c.numberOfClasses)
		for i := range c.classes {
			c.classes[i] = classAttrib{
				innerClassInfoIndex:   p.ReadU2(),
				outerClassInfoIndex:   p.ReadU2(),
				innerNameIndex:        p.ReadU2(),
				innerClassAccessFlags: p.ReadU2(),
			}
		}
		info = c
	case "BootstrapMethods":
		var boot bootstrapMethodsAttribute
		boot.basic = b
		boot.numBootstrapMethods = p.ReadU2()
		boot.bootstrapMethods = make([]bootstrapMethod, boot.numBootstrapMethods)
		for i := range boot.bootstrapMethods {
			m := bootstrapMethod{
				bootstrapMethodRef:    p.ReadU2(),
				numBootstrapArguments: p.ReadU2(),
			}
			m.bootstrapArguments = make([]uint16, m.numBootstrapArguments)
			for i := range m.bootstrapArguments {
				m.bootstrapArguments[i] = p.ReadU2()
			}
			boot.bootstrapMethods[i] = m
		}
		info = boot
	case "SourceFile":
		info = sourceFileAttribute{
			basic:           b,
			sourceFileIndex: p.ReadU2(),
		}
	case "ConstantValue":
		info = constantValueAttribute{
			basic:              b,
			constantValueIndex: p.ReadU2(),
		}
	default: // should ignore unknown attributes
		fmt.Printf("Unknown attribute: %s\n", tag.str)
		info = b
		bytes := make([]byte, b.attributeLength)
		_ = p.Read(bytes)
	}
	return info
}

type basic struct {
	attributeNameIndex uint16
	attributeLength    uint32
}

type lineNumberTable struct {
	basic
	lineNumberTableLength uint16
	lineNumberTable       []lineNumberTableEntry
}

type lineNumberTableEntry struct {
	startPc    uint16
	lineNumber uint16
}

type codeAttribute struct {
	basic
	maxStack             uint16
	maxLocals            uint16
	codeLength           uint32
	code                 []byte
	exceptionTableLength uint16
	exceptionTable       []exception
	attributesCount      uint16
	attributes           []AttributeInfo
}

type exception struct {
	startPc   uint16
	endPc     uint16
	handlerPc uint16
	catchType uint16
}

func readException(p *Parser) exception {
	return exception{
		startPc:   p.ReadU2(),
		endPc:     p.ReadU2(),
		handlerPc: p.ReadU2(),
		catchType: p.ReadU2(),
	}
}

type innerClassesAttribute struct {
	basic
	numberOfClasses uint16
	classes         []classAttrib
}

type classAttrib struct {
	innerClassInfoIndex   uint16
	outerClassInfoIndex   uint16
	innerNameIndex        uint16
	innerClassAccessFlags uint16
}

type bootstrapMethodsAttribute struct {
	basic
	numBootstrapMethods uint16
	bootstrapMethods    []bootstrapMethod
}

type bootstrapMethod struct {
	bootstrapMethodRef    uint16
	numBootstrapArguments uint16
	bootstrapArguments    []uint16
}

type sourceFileAttribute struct {
	basic
	sourceFileIndex uint16
}

type constantValueAttribute struct {
	basic
	constantValueIndex uint16
}
