package parser

type FieldInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []AttributeInfo
}

func ReadFieldInfo(c *RawClass, p *Parser) (info FieldInfo) {
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

func (info *FieldInfo) GetName(class RawClass) string {
	return class.GetUtf8(info.nameIndex)
}

func (info *FieldInfo) GetType(class RawClass) string {
	return class.GetUtf8(info.descriptorIndex)
}

func (info *FieldInfo) IsPublic() bool {
	return info.accessFlags&accPublic != 0
}

func (info *FieldInfo) IsStatic() bool {
	return info.accessFlags&accStatic != 0
}

func (info *FieldInfo) GetConstantValue(class RawClass) string {
	for _, v := range info.attributes {
		if c, ok := v.(constantValueAttribute); ok {
			val, _ := class.GetConstant(int(c.constantValueIndex))
			return val
		}
	}
	panic("doesn't have constant value attribute")
}
