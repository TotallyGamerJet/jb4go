package transformer

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"strings"
)

const (
	booleanJ = "boolean"
	doubleJ  = "double"
	floatJ   = "float"
	charJ    = "char"
	byteJ    = "byte"
	shortJ   = "short"
	intJ     = "int"
	longJ    = "long"
	voidJ    = "void"
)

// TranslateIdent takes a java type and converts it to a single character
// identifier or if its an Object name it converts it to a valid Go name.
func TranslateIdent(str string) string {
	switch str {
	case byteJ:
		return "B"
	case charJ:
		return "C"
	case doubleJ:
		return "D"
	case floatJ:
		return "F"
	case intJ:
		return "I"
	case longJ:
		return "J"
	case shortJ:
		return "S"
	case booleanJ:
		return "Z"
	case "[":
		return "R"
	case voidJ:
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
	case byteJ:
		f.Int8()
	case charJ:
		f.Uint16()
	case doubleJ:
		f.Float64()
	case floatJ:
		f.Float32()
	case intJ:
		f.Int32()
	case longJ:
		f.Int64()
	case shortJ:
		f.Int16()
	case booleanJ:
		f.Bool()
	case "[":
		panic("not implemented")
	case voidJ: // do nothing
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
	case charJ:
		return "int32" //jvm doesn't distinguish between char and ints "uint16"
	case shortJ:
		return "int32" //jvm doesn't distinguish between shorts and ints "int16"
	case byteJ:
		return "int8"
	case intJ:
		return "int32"
	case longJ:
		return "int64"
	case floatJ:
		return "float32"
	case doubleJ:
		return "float64"
	case booleanJ:
		return "int32" //jvm doesn't distinguish between bools and ints "bool"
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
		return byteJ
	case "C":
		return charJ
	case "D":
		return doubleJ
	case "F":
		return floatJ
	case "I":
		return intJ
	case "J":
		return longJ
	case "S":
		return shortJ
	case "Z":
		return booleanJ
	case voidJ:
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
				params = append(params, nameAndType{type_: temp, isArray: isArray})
				temp = ""
				isName = false
			}
			continue
		}
		switch v {
		case '(', ')':
			continue
		case 'B':
			temp += byteJ
		case 'C':
			temp += charJ
		case 'D':
			temp += doubleJ
		case 'F':
			temp += floatJ
		case 'I':
			temp += intJ
		case 'J':
			temp += longJ
		case 'L':
			isName = true
			continue
		case 'S':
			temp += shortJ
		case 'Z':
			temp += booleanJ
		case 'V':
			temp += voidJ
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

func arrayTypeCodes(code int) string {
	switch code {
	case 4:
		return booleanJ
	case 5:
		return charJ
	case 6:
		return floatJ
	case 7:
		return doubleJ
	case 8:
		return byteJ
	case 9:
		return shortJ
	case 10:
		return intJ
	case 11:
		return longJ
	default:
		panic(fmt.Sprintf("unknown code: %d", code))
	}
}
