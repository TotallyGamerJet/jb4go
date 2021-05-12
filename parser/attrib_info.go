package parser

import "fmt"

type AttributeInfo struct {
	NameIndex uint16
	Info      []byte
}

func ReadAttributeInfo(p *Parser) (info AttributeInfo, err error) {
	//attribute_info {
	//    u2 attribute_name_index;
	//    u4 attribute_length;
	//    u1 info[attribute_length];
	//}
	wrap := func(s string, e error) error {
		return fmt.Errorf("failed to read %s: %w", s, e)
	}
	if info.NameIndex, err = p.ReadU2(); err != nil {
		return info, wrap("name index", err)
	}
	length, err := p.ReadU4()
	if err != nil {
		return info, wrap("attrib length", err)
	}
	var n int
	info.Info = make([]byte, length)
	n, err = p.Read(info.Info)
	if err != nil {
		return info, wrap("attrib info", err)
	}
	if n != len(info.Info) {
		return info, fmt.Errorf("didn't read all of attrib info")
	}
	return info, nil
}

func ReadAllAttributeInfos(p *Parser) (Attributes []AttributeInfo, err error) {
	wrap := func(s string, e error) error {
		return fmt.Errorf("failed to read %s: %w", s, e)
	}
	attribCount, err := p.ReadU2()
	if err != nil {
		return Attributes, wrap("attribute count", err)
	}
	Attributes = make([]AttributeInfo, attribCount)
	for i := range Attributes {
		Attributes[i], err = ReadAttributeInfo(p)
		if err != nil {
			return Attributes, wrap(fmt.Sprintf("attrib (%d)", i), err)
		}
	}
	return Attributes, nil
}
