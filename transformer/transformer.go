package transformer

import (
	"github.com/totallygamerjet/jb4go/parser"
)

// JClass is an easy to use representation of a RawClass.
// It provides everything you need to transpile java byte code
type JClass struct {
	SrcFileName string    `json:"filename"`
	Name        string    `json:"name"`
	SuperName   string    `json:"super"`
	IsPublic    bool      `json:"public"`
	Fields      []JField  `json:"fields"`
	Methods     []JMethod `json:"methods"`
	// TODO: attributes
}

type JField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsPublic bool   `json:"public"`
	IsStatic bool   `json:"static"`
}

type JMethod struct {
	Name      string        `json:"name"`
	IsPublic  bool          `json:"public"`
	IsStatic  bool          `json:"static"`
	Params    []nameAndType `json:"params"`
	Code      []basicBlock  `json:"code"`
	MaxStack  int           `json:"maxStack"`
	MaxLocals int           `json:"maxLocals"`
	Return    string        `json:"return"`
}

func Simplify(raw parser.RawClass) (c JClass, err error) {
	c.SrcFileName = raw.GetFileName()
	c.Name = raw.GetName()
	c.SuperName = raw.GetSuperName()
	c.IsPublic = raw.IsPublic()
	c.Fields = make([]JField, len(raw.Fields))
	for i, f := range raw.Fields {
		var field JField
		field.Name = f.GetName(raw)
		field.Type = getJavaType(f.GetType(raw))
		field.IsPublic = f.IsPublic()
		field.IsStatic = f.IsStatic()
		c.Fields[i] = field
	}
	c.Methods = make([]JMethod, len(raw.Methods))
	for i, info := range raw.Methods {
		var m JMethod
		m.Name = info.GetName(raw)
		m.IsPublic = info.IsPublic()
		m.IsStatic = info.IsStatic()
		m.Params, m.Return = translateParams(info.GetDescriptor(raw))
		var code []byte
		code, m.MaxStack, m.MaxLocals = info.GetCode()
		instrs := readInstructions(code)
		blocks := createBasicBlocks(instrs)
		var params = make([]string, len(m.Params))
		for i, v := range m.Params {
			params[i] = v.type_
		}
		if !m.IsStatic {
			params = append([]string{ValidateName(c.Name)}, params...)
		}
		createIntermediate(blocks, raw, params)
		m.Code = blocks
		c.Methods[i] = m
	}
	return c, nil
}
