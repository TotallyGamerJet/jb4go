package parser

import "fmt"

type FieldInfo struct {
	AccessFlags uint16
	NameIndex   uint16
	DescIndex   uint16
	Attributes  []AttributeInfo
}

func ReadFieldInfo(p *Parser) (info FieldInfo, err error) {
	//field_info {
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
		return info, wrap("access flags", err)
	}
	if info.NameIndex, err = p.ReadU2(); err != nil {
		return info, wrap("name index", err)
	}
	if info.DescIndex, err = p.ReadU2(); err != nil {
		return info, wrap("descriptor index", err)
	}
	if info.Attributes, err = ReadAllAttributeInfos(p); err != nil {
		return info, wrap("field info", err)
	}
	return info, nil
}

func (f FieldInfo) IsPublic() bool {
	return f.AccessFlags&AccPublic != 0
}

func (f FieldInfo) IsPrivate() bool {
	return f.AccessFlags&AccPrivate != 0
}

func (f FieldInfo) IsProtected() bool {
	return f.AccessFlags&AccProtected != 0
}

func (f FieldInfo) IsStatic() bool {
	return f.AccessFlags&AccStatic != 0
}

func (f FieldInfo) IsFinal() bool {
	return f.AccessFlags&AccFinal != 0
}

func (f FieldInfo) IsVolatile() bool {
	return f.AccessFlags&AccVolatile != 0
}

func (f FieldInfo) IsTransient() bool {
	return f.AccessFlags&AccTransient != 0
}

func (f FieldInfo) IsSynthetic() bool {
	return f.AccessFlags&AccSynthetic != 0
}

func (f FieldInfo) IsEnum() bool {
	return f.AccessFlags&AccEnum != 0
}
