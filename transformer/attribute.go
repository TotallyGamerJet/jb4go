package transformer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/totallygamerjet/jb4go/parser"
)

const (
	AttribConstantValue                  = "ConstantValue"
	AttribCode                           = "Code"
	AttribStackMapTable                  = "StackMapTable"
	AttribExceptions                     = "Exceptions"
	AttribInnerClasses                   = "InnerClasses"
	AttribEnclosingMethod                = "EnclosingMethod"
	AttribSynthetic                      = "Synthetic"
	AttribSignature                      = "Signature"
	AttribSourceFile                     = "SourceFile"
	AttribSourceDebugExtension           = "SourceDebugExtension"
	AttribLineNumberTable                = "LineNumberTable"
	AttribLocalVarTable                  = "LocalVariableTable"
	AttribLocalVarTypeTable              = "LocalVariableTypeTable"
	AttribDeprecated                     = "Deprecated"
	AttribRuntimeVisibleAnnotations      = "RuntimeVisibleAnnotations"
	AttribRuntimeInvisibleAnnotations    = "RuntimeInvisibleAnnotations"
	AttribRuntimeVisibleParAnnotations   = "RuntimeVisibleParameterAnnotations"
	AttribRuntimeInvisibleParAnnotations = "RuntimeInvisibleParameterAnnotations"
	AttribAnnotationDefault              = "AnnotationDefault"
	AttribBootstrapMethods               = "BootstrapMethods"
)

type Attribute struct {
	Name string
	Info []byte
}

func ResolveAttributeInfo(cp []parser.CPInfo, info parser.AttributeInfo) Attribute {
	return Attribute{
		Name: cp[info.NameIndex].UTF8(),
		Info: info.Info,
	}
}

func (a Attribute) check(t string) {
	if a.Name != t {
		panic(fmt.Sprintf("attribute %s is not type %s", a.Name, t))
	}
}

//ConstantValue_attribute {
//    u2 attribute_name_index;
//    u4 attribute_length;
//    u2 constantvalue_index;
//}
func (a Attribute) ConstantValue() uint16 {
	a.check(AttribConstantValue)
	return binary.BigEndian.Uint16(a.Info)
}

type CodeAttribute struct {
	MaxStack       uint16
	MaxLocals      uint16
	Code           []byte
	ExceptionTable []ExceptionTableEntry
	Attributes     []parser.AttributeInfo
}

type ExceptionTableEntry struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

func (a Attribute) Code() (code CodeAttribute) {
	a.check(AttribCode)
	p := parser.Parse(bytes.NewReader(a.Info))
	// TODO: error checking
	code.MaxStack, _ = p.ReadU2()
	code.MaxLocals, _ = p.ReadU2()
	length, _ := p.ReadU4()
	code.Code = make([]byte, length)
	_, _ = p.Read(code.Code)
	exLength, _ := p.ReadU2()
	code.ExceptionTable = make([]ExceptionTableEntry, exLength)
	for i := range code.ExceptionTable {
		//u2 start_pc;
		//        u2 end_pc;
		//        u2 handler_pc;
		//        u2 catch_type;
		ex := ExceptionTableEntry{}
		ex.StartPc, _ = p.ReadU2()
		ex.EndPc, _ = p.ReadU2()
		ex.HandlerPc, _ = p.ReadU2()
		ex.CatchType, _ = p.ReadU2()
		code.ExceptionTable[i] = ex
	}
	code.Attributes, _ = parser.ReadAllAttributeInfos(p)
	return code
}
