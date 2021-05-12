package parser

import (
	"fmt"
)

type ClassFile struct {
	Magic        uint32
	MinorVersion uint16
	MajorVersion uint16
	ConstantPool []CPInfo
	AccessFlags  uint16
	ThisClass    uint16
	SuperClass   uint16
	Interfaces   []uint16
	Fields       []FieldInfo
	Methods      []MethodInfo
	Attributes   []AttributeInfo
}

func ReadClass(p *Parser) (cf ClassFile, err error) {
	wrap := func(s string, e error) error {
		return fmt.Errorf("failed to read %s: %w", s, e)
	}
	if cf.Magic, err = p.ReadU4(); err != nil {
		return cf, wrap("magic", err)
	}
	if cf.MinorVersion, err = p.ReadU2(); err != nil {
		return cf, wrap("minor version", err)
	}
	if cf.MajorVersion, err = p.ReadU2(); err != nil {
		return cf, wrap("major version", err)
	}
	cpCount, err := p.ReadU2()
	if err != nil {
		return cf, wrap("cp count", err)
	}
	cf.ConstantPool = make([]CPInfo, cpCount)
	for i := 1; i < len(cf.ConstantPool); i++ {
		switch cf.ConstantPool[i-1].Tag { // there is an empty space after these two
		case TagDouble, TagLong:
			continue
		}
		cf.ConstantPool[i], err = ReadCPInfo(p)
		if err != nil {
			return cf, wrap(fmt.Sprintf("cp index (%d)", i), err)
		}
	}
	if cf.AccessFlags, err = p.ReadU2(); err != nil {
		return cf, wrap("access flags", err)
	}
	if cf.ThisClass, err = p.ReadU2(); err != nil {
		return cf, wrap("this class", err)
	}
	if cf.SuperClass, err = p.ReadU2(); err != nil {
		return cf, wrap("super class", err)
	}
	interfaceCount, err := p.ReadU2()
	if err != nil {
		return cf, wrap("interface count", err)
	}
	cf.Interfaces = make([]uint16, interfaceCount)
	for i := range cf.Interfaces {
		cf.Interfaces[i], err = p.ReadU2()
		if err != nil {
			return cf, wrap(fmt.Sprintf("interface index (%d)", i), err)
		}
	}
	fieldCount, err := p.ReadU2()
	if err != nil {
		return cf, wrap("field count", err)
	}
	cf.Fields = make([]FieldInfo, fieldCount)
	for i := range cf.Fields {
		cf.Fields[i], err = ReadFieldInfo(p)
		if err != nil {
			return cf, wrap(fmt.Sprintf("field index (%d)", i), err)
		}
	}
	metCount, err := p.ReadU2()
	if err != nil {
		return cf, wrap("met count", err)
	}
	cf.Methods = make([]MethodInfo, metCount)
	for i := range cf.Methods {
		cf.Methods[i], err = ReadMethodInfo(p)
		if err != nil {
			return cf, wrap(fmt.Sprintf("method index (%d)", i), err)
		}
	}
	if cf.Attributes, err = ReadAllAttributeInfos(p); err != nil {
		return cf, wrap("classfile", err)
	}
	return cf, nil
}

func (c ClassFile) Name() string {
	idx := c.ConstantPool[c.ThisClass].Class()
	return c.ConstantPool[idx].UTF8()
}

func (c ClassFile) SuperName() string {
	idx := c.ConstantPool[c.SuperClass].Class()
	return c.ConstantPool[idx].UTF8()
}

func (c ClassFile) IsPublic() bool {
	return c.AccessFlags&AccPublic != 0
}

func (c ClassFile) IsFinal() bool {
	return c.AccessFlags&AccFinal != 0
}

func (c ClassFile) IsSuper() bool {
	return c.AccessFlags&AccSuper != 0
}

func (c ClassFile) IsInterface() bool {
	return c.AccessFlags&AccInterface != 0
}

func (c ClassFile) IsAbstract() bool {
	return c.AccessFlags&AccAbstract != 0
}

func (c ClassFile) IsSynthetic() bool {
	return c.AccessFlags&AccSynthetic != 0
}

func (c ClassFile) IsAnnotation() bool {
	return c.AccessFlags&AccAnnotation != 0
}

func (c ClassFile) IsEnum() bool {
	return c.AccessFlags&AccEnum != 0
}
