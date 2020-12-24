package parser

type MethodInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	//attributesCount uint16
	attributes []AttributeInfo
}

func ReadMethodInfo(c *RawClass, p *Parser) (info MethodInfo) {
	info.accessFlags = p.ReadU2()
	info.nameIndex = p.ReadU2()
	info.descriptorIndex = p.ReadU2()
	attributesCount := p.ReadU2()
	info.attributes = make([]AttributeInfo, attributesCount)
	for i := range info.attributes {
		info.attributes[i] = ReadAttributeInfo(c, p)
	}
	return info
}

func (m *MethodInfo) GetName(class RawClass) string {
	return class.GetUtf8(m.nameIndex)
}

// GetDescriptor returns a string representing the params and return type
func (m *MethodInfo) GetDescriptor(class RawClass) string {
	return class.GetUtf8(m.descriptorIndex)
}

// GetCode returns the code, max_stack size, max_local size
func (m *MethodInfo) GetCode() ([]byte, int, int) {
	for _, v := range m.attributes {
		code, ok := v.(codeAttribute)
		if !ok {
			continue
		}
		return code.code, int(code.maxStack), int(code.maxLocals)
	}
	return nil, 0, 0
}

func (m *MethodInfo) IsPublic() bool {
	return m.accessFlags&accPublic != 0
}

func (m *MethodInfo) IsStatic() bool {
	return m.accessFlags&accStatic != 0
}
