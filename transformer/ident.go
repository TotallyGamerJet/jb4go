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
			name := ValidateName(strings.ReplaceAll(javaType[1:], "/", "_"), true)
			f.Op("*").Qual("github.com/totallygamerjet/java", name)
			return
		}
		panic("not implemented: " + javaType)
	}
}

// takes in a name and its visibility and converts it to a valid name
func ValidateName(name string, public bool) string {
	if strings.HasPrefix(name, "L") && strings.HasSuffix(name, ";") {
		name = name[1 : len(name)-1] // trim off the l and ;
	}
	name = strings.ReplaceAll(name, "/", "_")
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
		return "*" + ValidateName(jType, true) // is this right?
	}
}
