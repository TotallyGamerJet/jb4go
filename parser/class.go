package parser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
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
	var wide bool
	for i := range c.constantPool {
		if wide { // skip over 8-byte constants
			wide = false
			continue
		}
		c.constantPool[i], wide = ReadCPInfo(p)
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
			return c.GetUtf8(src.sourceFileIndex)
		}
	}
	return c.GetName() + ".java" // if its not found use the classes' name plus ext
}

// takes the index in the constant pool starting at 1 and returns the item at that location
func (c *RawClass) GetCP(index int) CPInfo {
	return c.constantPool[index-1]
}

func (c *RawClass) GetUtf8(index uint16) string {
	return c.GetCP(int(index)).(utf8Info).str
}

// returns the name of this raw class
func (c *RawClass) GetName() string {
	return c.GetClass(int(c.thisClass))
}

// returns the name of the super object of this raw class
func (c *RawClass) GetSuperName() string {
	return c.GetClass(int(c.superClass))
}

func (c *RawClass) IsPublic() bool {
	return c.accessFlags&accPublic != 0
}

func (c *RawClass) GetClass(index int) string {
	info := c.GetCP(index).(classInfo)
	return c.GetUtf8(info.nameIndex)
}

func (c *RawClass) GetNameAndType(index int) (string, string) {
	info := c.GetCP(index).(nameAndTypeInfo)
	return c.GetUtf8(info.nameIndex), c.GetUtf8(info.descriptorIndex)
}

// returns the class name, the name of the field and the type.
// CLASSNAME, FIELD_NAME, TYPE
func (c *RawClass) GetFieldRef(index int) (string, string, string) {
	info := c.GetCP(index).(fieldRefInfo)
	n, t := c.GetNameAndType(int(info.nameAndTypeIndex))
	return c.GetClass(int(info.classIndex)), n, t
}

func (c *RawClass) GetMethodRef(index int) (string, string, string) {
	info := c.GetCP(index).(methodRefInfo)
	n, t := c.GetNameAndType(int(info.nameAndTypeIndex))
	return c.GetClass(int(info.classIndex)), n, t
}

// GetConstant returns the constant at the index as a string and a string representation
// of the type
func (c *RawClass) GetConstant(index int) (string, string) {
	switch info := c.GetCP(index); v := info.(type) {
	case stringInfo:
		return "`" + c.GetUtf8(v.stringIndex) + "`", "java/lang/String"
	case doubleInfo:
		bits := (int64(v.high) << 32) + int64(v.low)
		var double = math.Float64frombits(uint64(bits))
		return strconv.FormatFloat(double, 'E', -1, 64), "double"
	case integerInfo:
		return strconv.Itoa(int(v.bytes)), "int"
	case longInfo:
		bits := (int64(v.high) << 32) + int64(v.low)
		return strconv.Itoa(int(bits)), "long"
	default:
		panic(fmt.Sprintf("unknown constant: %d", v.Tag()))
	}
}
