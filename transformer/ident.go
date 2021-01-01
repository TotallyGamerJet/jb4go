package transformer

import (
	"github.com/dave/jennifer/jen"
	"strings"
)

// TranslateIdent takes a java type and converts it to a single character
// identifier or if its an Object name it converts it to a valid Go name.
func TranslateIdent(str string) string {
	switch str {
	case "byte":
		return "B"
	case "char":
		return "C"
	case "double":
		return "D"
	case "float":
		return "F"
	case "int":
		return "I"
	case "long":
		return "J"
	case "short":
		return "S"
	case "boolean":
		return "Z"
	case "[":
		return "R"
	case "void":
		return "V"
	case "java/lang/String":
		return "G"
	default:
		return strings.ReplaceAll(str, "/", "_")
	}
}

// toGoType takes a java type as a string and add to the file the proper go type
func toGoType(javaType string, f *jen.Statement) {
	switch javaType {
	case "byte":
		f.Int8()
	case "char":
		f.Uint16()
	case "double":
		f.Float64()
	case "float":
		f.Float32()
	case "int":
		f.Int32()
	case "long":
		f.Int64()
	case "short":
		f.Int16()
	case "boolean":
		f.Bool()
	case "[":
		panic("not implemented")
	case "void": // do nothing
	default:
		if strings.HasPrefix(javaType, "L") {
			name := ValidateName(strings.ReplaceAll(javaType[1:], "/", "_"))
			f.Op("*").Qual("github.com/totallygamerjet/java", name)
			return
		}
		panic("not implemented: " + javaType)
	}
}

// takes in a name and its visibility and converts it to a valid name
func ValidateName(name string) string {
	if strings.HasPrefix(name, "L") && strings.HasSuffix(name, ";") {
		name = name[1 : len(name)-1] // trim off the l and ;
	}
	name = strings.ReplaceAll(name, "/", "_")
	//if public {
	//	name = "P_" + name
	//} else { // private method
	//	name = "_" + name
	//}
	//if jen.IsReservedWord(name) { // add _ if name is reversed
	//	name += "_"
	//}
	return name
}

func getGoType(jType string) string {
	switch jType {
	case "char":
		return "uint16"
	case "short":
		return "int16"
	case "byte":
		return "int8"
	case "int":
		return "int32"
	case "long":
		return "int64"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "boolean":
		return "bool"
	default:
		if jType[0] == '[' { // ignore arrays
			return jType
		}
		return "*java_lang_Object" //"*" + ValidateName(jType) // is this right?
	}
}

// takes a shortened java type (ex. I for int) and returns the java type
func getJavaType(str string) string {
	switch str {
	case "B":
		return "byte"
	case "C":
		return "char"
	case "D":
		return "double"
	case "F":
		return "float"
	case "I":
		return "int"
	case "J":
		return "long"
	case "S":
		return "short"
	case "Z":
		return "boolean"
	case "void":
		return ""
	default:
		var out string
		if strings.HasPrefix(str, "[") {
			out += "[]"
			str = str[1:] // consume [
		}
		if strings.HasPrefix(str, "L") && strings.HasSuffix(str, ";") {
			out += str[1 : len(str)-1] // remove the beginning and end
		}
		return out
	}
}

type nameAndType struct {
	type_   string
	isArray bool
}

func translateParams(str string) (params []nameAndType, ret string) {
	var isName = false // true if reading in class isName
	var temp string
	var isArray = false
	for _, v := range str {
		if isName {
			if v != ';' {
				temp += string(v)
			} else {
				// end of class name
				params = append(params, struct {
					type_   string
					isArray bool
				}{type_: temp, isArray: isArray})
				temp = ""
				isName = false
			}
			continue
		}
		switch v {
		case '(', ')':
			continue
		case 'B':
			temp += "byte"
		case 'C':
			temp += "char"
		case 'D':
			temp += "double"
		case 'F':
			temp += "float"
		case 'I':
			temp += "int"
		case 'J':
			temp += "long"
		case 'L':
			isName = true
			continue
		case 'S':
			temp += "short"
		case 'Z':
			temp += "boolean"
		case 'V':
			temp += "void"
		case '[':
			isArray = true
			continue //get the type of this array
		}
		params = append(params, nameAndType{type_: temp, isArray: isArray})
		temp = ""
	}
	if len(params) > 0 {
		ret = params[len(params)-1].type_
		params = params[:len(params)-1]
	}
	return params, ret
}
