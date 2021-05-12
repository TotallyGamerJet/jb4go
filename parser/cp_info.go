package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

const (
	TagUtf8               = 1
	TagInteger            = 3
	TagFloat              = 4
	TagLong               = 5
	TagDouble             = 6
	TagClass              = 7
	TagString             = 8
	TagFieldref           = 9
	TagMethodref          = 10
	TagInterfaceMethodref = 11
	TagNameAndType        = 12
	TagMethodHandle       = 15
	TagMethodType         = 16
	TagInvokeDynamic      = 18
)

type CPInfo struct {
	Tag  byte
	Info []byte
}

func ReadCPInfo(p *Parser) (info CPInfo, err error) {
	if info.Tag, err = p.ReadU1(); err != nil {
		return info, fmt.Errorf("failed to read cp info: %w", err)
	}
	switch info.Tag {
	case TagUtf8:
		length, err := p.ReadU2()
		if err != nil {
			return info, fmt.Errorf("failed to read UTF8 length: %w", err)
		}
		info.Info = make([]byte, length)
	case TagClass, TagString, TagMethodType:
		info.Info = make([]byte, 2)
	case TagMethodHandle:
		info.Info = make([]byte, 3)
	case TagInteger, TagFloat, TagFieldref, TagMethodref,
		TagInterfaceMethodref, TagNameAndType, TagInvokeDynamic:
		info.Info = make([]byte, 4)
	case TagLong, TagDouble:
		info.Info = make([]byte, 8)
	default:
		return info, fmt.Errorf("unknown cp info tag %d", info.Tag)
	}
	var n int
	n, err = p.Read(info.Info)
	if err != nil {
		return info, fmt.Errorf("failed to read cp type %d: %w", info.Tag, err)
	}
	if n != len(info.Info) {
		return info, fmt.Errorf("didn't read all of cp info type %d", info.Tag)
	}
	return info, nil
}

func (info CPInfo) check(t byte) {
	if info.Tag != t {
		panic(fmt.Sprintf("cp info type %d is not type %d", info.Tag, t))
	}
}

//CONSTANT_InvokeDynamic_info {
//    u1 tag;
//    u2 bootstrap_method_attr_index;
//    u2 name_and_type_index;
//}
func (info CPInfo) InvokeDynamic() (uint16, uint16) {
	info.check(TagInvokeDynamic)
	be := binary.BigEndian
	return be.Uint16(info.Info[:2]), be.Uint16(info.Info[2:])
}

//CONSTANT_MethodType_info {
//    u1 tag;
//    u2 descriptor_index;
//}
func (info CPInfo) MethodType() uint16 {
	info.check(TagMethodType)
	return binary.BigEndian.Uint16(info.Info)
}

//CONSTANT_MethodHandle_info {
//    u1 tag;
//    u1 reference_kind;
//    u2 reference_index;
//}
func (info CPInfo) MethodHandle() (uint8, uint16) {
	info.check(TagMethodHandle)
	return info.Info[0], binary.BigEndian.Uint16(info.Info[1:])
}

//CONSTANT_Utf8_info {
//    u1 tag;
//    u2 length; // implicit
//    u1 bytes[length];
//}
func (info CPInfo) UTF8() string {
	info.check(TagUtf8)
	length := len(info.Info)
	buf := bytes.Buffer{}
	buf.Grow(length)
	for i := 0; i < length; {
		x := int(info.Info[i]) & 0xFF
		switch x >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			/* 0xxxxxxx*/
			i++
			_ = buf.WriteByte(byte(x))
		case 12, 13:
			/* 110x xxxx   10xx xxxx*/
			i += 2
			y := int(info.Info[i-1])
			_ = buf.WriteByte(byte(((x & 0x1f) << 6) + (y & 0x3f)))
		case 14:
			/* 1110 xxxx  10xx xxxx  10xx xxxx */
			i += 3
			y := int(info.Info[i-2])
			z := int(info.Info[i-1])
			_ = buf.WriteByte(byte(((x & 0xf) << 12) + ((y & 0x3f) << 6) + (z & 0x3f)))
		default:
			panic("malformed input")
		}
	}
	return buf.String()
}

//CONSTANT_NameAndType_info {
//    u1 tag;
//    u2 name_index;
//    u2 descriptor_index;
//}
func (info CPInfo) NameAndType() (uint16, uint16) {
	info.check(TagNameAndType)
	be := binary.BigEndian
	return be.Uint16(info.Info[:2]), be.Uint16(info.Info[2:])
}

//CONSTANT_Long_info {
//    u1 tag;
//    u4 high_bytes;
//    u4 low_bytes;
//}
func (info CPInfo) Long() int64 {
	info.check(TagLong)
	return int64(binary.BigEndian.Uint64(info.Info))
}

//CONSTANT_Double_info {
//    u1 tag;
//    u4 high_bytes;
//    u4 low_bytes;
//}
func (info CPInfo) Double() float64 {
	info.check(TagDouble)
	bits := binary.BigEndian.Uint64(info.Info)
	return math.Float64frombits(bits)
}

//CONSTANT_Integer_info {
//    u1 tag;
//    u4 bytes;
//}
func (info CPInfo) Integer() int32 {
	info.check(TagInteger)
	return int32(binary.BigEndian.Uint32(info.Info))
}

//CONSTANT_Float_info {
//    u1 tag;
//    u4 bytes;
//}
func (info CPInfo) Float() float32 {
	info.check(TagFloat)
	bits := binary.BigEndian.Uint32(info.Info)
	return math.Float32frombits(bits) // TODO: ensure this is right
}

//CONSTANT_String_info {
//    u1 tag;
//    u2 string_index;
//}
func (info CPInfo) String() uint16 {
	info.check(TagString)
	return binary.BigEndian.Uint16(info.Info)
}

// Class returns the class name index. This method panics if
// CPInfo is not a class type
//CONSTANT_Class_info {
//    u1 tag;
//    u2 name_index;
//}
func (info CPInfo) Class() uint16 {
	info.check(TagClass)
	return binary.BigEndian.Uint16(info.Info)
}

//CONSTANT_Fieldref_info {
//    u1 tag;
//    u2 class_index;
//    u2 name_and_type_index;
//}
//
func (info CPInfo) FieldRef() (uint16, uint16) {
	info.check(TagFieldref)
	be := binary.BigEndian
	return be.Uint16(info.Info[:2]), be.Uint16(info.Info[2:])
}

//CONSTANT_Methodref_info {
//    u1 tag;
//    u2 class_index;
//    u2 name_and_type_index;
//}
//
func (info CPInfo) MethodRef() (uint16, uint16) {
	info.check(TagMethodref)
	be := binary.BigEndian
	return be.Uint16(info.Info[:2]), be.Uint16(info.Info[2:])
}

//CONSTANT_InterfaceMethodref_info {
//    u1 tag;
//    u2 class_index;
//    u2 name_and_type_index;
//}
func (info CPInfo) InterfaceMethodRef() (uint16, uint16) {
	info.check(TagInterfaceMethodref)
	be := binary.BigEndian
	return be.Uint16(info.Info[:2]), be.Uint16(info.Info[2:])
}
