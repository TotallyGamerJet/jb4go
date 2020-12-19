package parser

import (
	"bytes"
	"fmt"
)

type cpTag byte

const (
	_       cpTag = iota
	tagUtf8       //CONSTANT_Utf8	1
	_
	tagInteger            //CONSTANT_Integer	3
	tagFloat              //CONSTANT_Float	4
	tagLong               //CONSTANT_Long	5
	tagDouble             //CONSTANT_Double	6
	tagClass              //CONSTANT_Class	7
	tagString             //CONSTANT_String	8
	tagFieldref           //CONSTANT_Fieldref	9
	tagMethodref          //CONSTANT_Methodref	10
	tagInterfaceMethodref //CONSTANT_InterfaceMethodref	11
	tagNameAndType        //CONSTANT_NameAndType	12
	_
	_
	tagMethodHandle //CONSTANT_MethodHandle	15
	tagMethodType   //CONSTANT_MethodType	16
	_
	tagInvokeDynamic //CONSTANT_InvokeDynamic	18
)

type CPInfo interface {
	Tag() cpTag
}

func ReadCPInfo(p *Parser) (info CPInfo) {
	tag := p.ReadU1()
	switch cpTag(tag) {
	case tagUtf8:
		length := p.ReadU2()
		b := make([]byte, length)
		_ = p.Read(b)
		info = utf8Info{
			str: readUTF8(b),
		}
	//case tagInteger:
	//case tagFloat:
	//case tagLong:
	//case tagDouble:
	case tagClass:
		info = classInfo{
			nameIndex: p.ReadU2(),
		}
	case tagString:
		info = stringInfo{
			stringIndex: p.ReadU2(),
		}
	case tagFieldref:
		info = fieldRefInfo{
			classIndex:       p.ReadU2(),
			nameAndTypeIndex: p.ReadU2(),
		}
	case tagMethodref:
		info = methodRefInfo{
			classIndex:       p.ReadU2(),
			nameAndTypeIndex: p.ReadU2(),
		}
	//case tagInterfaceMethodref:
	case tagNameAndType:
		info = nameAndTypeInfo{
			nameIndex:       p.ReadU2(),
			descriptorIndex: p.ReadU2(),
		}
	case tagMethodHandle:
		info = methodHandleInfo{
			referenceKind:  p.ReadU1(),
			referenceIndex: p.ReadU2(),
		}
	//case tagMethodType:
	case tagInvokeDynamic:
		info = invokeDynamicInfo{
			bootstrapMethodAttrIndex: p.ReadU2(),
			nameAndTypeIndex:         p.ReadU2(),
		}
	default:
		panic(fmt.Sprintf("unknown constant pool tag: %d", tag))
	}
	return info
}

type methodRefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (i methodRefInfo) Tag() cpTag {
	return tagMethodref
}

type classInfo struct {
	nameIndex uint16
}

func (i classInfo) Tag() cpTag {
	return tagClass
}

type fieldRefInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (i fieldRefInfo) Tag() cpTag {
	return tagFieldref
}

type stringInfo struct {
	stringIndex uint16
}

func (i stringInfo) Tag() cpTag {
	return tagString
}

type invokeDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

func (i invokeDynamicInfo) Tag() cpTag {
	return tagInvokeDynamic
}

type utf8Info struct {
	str string
}

func (i utf8Info) Tag() cpTag {
	return tagUtf8
}

type nameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

func (i nameAndTypeInfo) Tag() cpTag {
	return tagNameAndType
}

type methodHandleInfo struct {
	referenceKind  byte
	referenceIndex uint16
}

func (i methodHandleInfo) Tag() cpTag {
	return tagMethodHandle
}

// converts from modified utf8 to unicode
func readUTF8(in []byte) string {
	length := len(in)
	buf := bytes.Buffer{}
	buf.Grow(length)
	var i = 0 // start after the length
	for i < len(in) {
		x := int(in[i]) & 0xFF
		switch x >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			/* 0xxxxxxx*/
			i++
			_ = buf.WriteByte(byte(x))
		case 12, 13:
			/* 110x xxxx   10xx xxxx*/
			i += 2
			y := int(in[i-1])
			_ = buf.WriteByte(byte(((x & 0x1f) << 6) + (y & 0x3f)))
		case 14:
			/* 1110 xxxx  10xx xxxx  10xx xxxx */
			i += 3
			y := int(in[i-2])
			z := int(in[i-1])
			_ = buf.WriteByte(byte(((x & 0xf) << 12) + ((y & 0x3f) << 6) + (z & 0x3f)))
		default:
			panic("malformed input")
		}
	}
	return buf.String()
}
