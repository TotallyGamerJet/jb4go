package parser

import "fmt"

type MethodInfo struct {
	AccessFlags uint16
	NameIndex   uint16
	DescIndex   uint16
	Attributes  []AttributeInfo
}

func ReadMethodInfo(p *Parser) (info MethodInfo, err error) {
	//method_info {
	//    u2             access_flags;
	//    u2             name_index;
	//    u2             descriptor_index;
	//    u2             attributes_count;
	//    attribute_info attributes[attributes_count];
	//}
	wrap := func(s string, e error) error {
		return fmt.Errorf("failed to read %s: %w", s, e)
	}
	if info.AccessFlags, err = p.ReadU2(); err != nil {
		return info, wrap("access flag", err)
	}
	if info.NameIndex, err = p.ReadU2(); err != nil {
		return info, wrap("name index", err)
	}
	if info.DescIndex, err = p.ReadU2(); err != nil {
		return info, wrap("desc index", err)
	}
	if info.Attributes, err = ReadAllAttributeInfos(p); err != nil {
		return info, wrap("method", err)
	}
	return info, nil
}

func (m MethodInfo) IsPublic() bool {
	return m.AccessFlags&AccPublic != 0
}

func (m MethodInfo) IsPrivate() bool {
	return m.AccessFlags&AccPrivate != 0
}

func (m MethodInfo) IsProtected() bool {
	return m.AccessFlags&AccProtected != 0
}

func (m MethodInfo) IsStatic() bool {
	return m.AccessFlags&AccStatic != 0
}

func (m MethodInfo) IsFinal() bool {
	return m.AccessFlags&AccFinal != 0
}

func (m MethodInfo) IsSynchronized() bool {
	return m.AccessFlags&AccSynchronized != 0
}

func (m MethodInfo) IsBridge() bool {
	return m.AccessFlags&AccBridge != 0
}

func (m MethodInfo) IsVarArgs() bool {
	return m.AccessFlags&AccVarArgs != 0
}

func (m MethodInfo) IsNative() bool {
	return m.AccessFlags&AccNative != 0
}

func (m MethodInfo) IsAbstract() bool {
	return m.AccessFlags&AccAbstract != 0
}

func (m MethodInfo) IsStrict() bool {
	return m.AccessFlags&AccStrict != 0
}

func (m MethodInfo) IsSynthetic() bool {
	return m.AccessFlags&AccSynthetic != 0
}
