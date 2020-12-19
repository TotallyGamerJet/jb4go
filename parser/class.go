package parser

import (
	"errors"
)

// Holds the raw information to a java bytecode class
type RawClass struct {
	magic        uint32
	minorVersion uint16
	majorVersion uint16
	constantPool []CPInfo
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	Fields       []FieldInfo
	Methods      []MethodInfo
	attributes   []AttributeInfo
}

func ReadClass(p *Parser) (RawClass, error) {
	var c RawClass
	if c.magic = p.ReadU4(); c.magic != 0xCAFEBABE {
		return c, errors.New("improper file")
	}
	c.minorVersion = p.ReadU2()
	c.majorVersion = p.ReadU2()
	if c.majorVersion != 52 && c.minorVersion != 0 {
		return c, errors.New("improper version")
	}
	constantPoolCount := p.ReadU2()
	c.constantPool = make([]CPInfo, constantPoolCount-1)
	for i := range c.constantPool {
		c.constantPool[i] = ReadCPInfo(p)
	}
	c.accessFlags = p.ReadU2()
	c.thisClass = p.ReadU2()
	c.superClass = p.ReadU2()
	interfacesCount := p.ReadU2()
	c.interfaces = make([]uint16, interfacesCount)
	for i := range c.interfaces {
		c.interfaces[i] = p.ReadU2()
	}
	fieldsCount := p.ReadU2()
	c.Fields = make([]FieldInfo, fieldsCount)
	for i := range c.Fields {
		c.Fields[i] = ReadFieldInfo(&c, p)
	}
	methodsCount := p.ReadU2()
	c.Methods = make([]MethodInfo, methodsCount)
	for i := range c.Methods {
		c.Methods[i] = ReadMethodInfo(&c, p)
	}
	attributesCount := p.ReadU2()
	c.attributes = make([]AttributeInfo, attributesCount)
	for i := range c.attributes {
		c.attributes[i] = ReadAttributeInfo(&c, p)
	}
	return c, nil
}

func (c *RawClass) GetFileName() string {
	for _, v := range c.attributes {
		if src, ok := v.(sourceFileAttribute); !ok {
			continue
		} else {
			return c.nameFromIndex(src.sourceFileIndex)
		}
	}
	return c.GetName() + ".java" // if its not found use the classes' name plus ext
}

// takes the index in the constant pool starting at 1 and returns the item at that location
func (c *RawClass) GetCP(index int) CPInfo {
	return c.constantPool[index-1]
}

func (c *RawClass) nameFromIndex(index uint16) string {
	return c.GetCP(int(index)).(utf8Info).str
}

// returns the name of this raw class
func (c *RawClass) GetName() string {
	info := c.GetCP(int(c.thisClass)).(classInfo)
	return c.nameFromIndex(info.nameIndex)
}

// returns the name of the super object of this raw class
func (c *RawClass) GetSuperName() string {
	info := c.GetCP(int(c.superClass)).(classInfo)
	return c.nameFromIndex(info.nameIndex)
}

func (c *RawClass) IsPublic() bool {
	return c.accessFlags&accPublic != 0
}
