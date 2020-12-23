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
	Name      string   `json:"name"`
	IsPublic  bool     `json:"public"`
	IsStatic  bool     `json:"static"`
	Params    []string `json:"params"`
	Code      []byte   `json:"code"`
	MaxStack  int      `json:"maxStack"`
	MaxLocals int      `json:"maxLocals"`
	Return    string   `json:"return"`
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
		field.Type = f.GetType(raw)
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
		m.Code, m.MaxStack, m.MaxLocals = info.GetCode()
		instrs := readInstructions(m.Code)
		blocks := createBasicBlocks(instrs)
		createIntermediate(blocks, raw)
		c.Methods[i] = m
	}
	return c, nil
}

func translateParams(str string) (params []string, ret string) {
	var isName = false // true if reading in class isName
	var temp string
	for _, v := range str {
		if isName {
			if v != ';' {
				temp += string(v)
			} else {
				// end of class name
				params = append(params, temp)
				temp = ""
				isName = false
			}
			continue
		}
		switch v {
		case '(', ')': // ignore
		case 'B':
			params = append(params, "byte")
		case 'C':
			params = append(params, "char")
		case 'D':
			params = append(params, "double")
		case 'F':
			params = append(params, "float")
		case 'I':
			params = append(params, "int")
		case 'J':
			params = append(params, "long")
		case 'L':
			isName = true
			//temp = "L"
		case 'S':
			params = append(params, "short")
		case 'Z':
			params = append(params, "boolean")
		case 'V':
			params = append(params, "void")
		case '[':
			params = append(params, "[")
		}
	}
	if len(params) > 0 {
		ret = params[len(params)-1]
		params = params[:len(params)-1]
	}
	return params, ret
}
