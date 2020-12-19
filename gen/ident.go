package gen

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
		return "A"
	case "void":
		return "V"
	case "Ljava/lang/String":
		return "G"
	default:
		if strings.HasPrefix(str, "L") {
			return strings.ReplaceAll(str[1:], "/", "_")
		}
		panic("unknown ident")
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
			name := validateName(strings.ReplaceAll(javaType[1:], "/", "_"), true)
			f.Op("*").Qual("github.com/totallygamerjet/java", name)
			return
		}
		panic("not implemented: " + javaType)
	}
}

// takes in a name and its visibility and converts it to a valid name
func validateName(name string, public bool) string {
	if public {
		name = "P_" + name
	} else { // private method
		name = "_" + name
	}
	if jen.IsReservedWord(name) { // add _ if name is reversed
		name += "_"
	}
	return name
}
