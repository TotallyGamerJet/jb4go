package main

import "fmt"

func new_java_io_PrintStream() *java_lang_Object {
	type fields struct{}
	return &java_lang_Object{
		name:   "java_io_PrintStream",
		super:  new_java_lang_Object(),
		fields: &fields{},
	}
}

func I_java_io_PrintStream_println_G_V(arg0 *java_lang_Object, arg1 *java_lang_Object) {
	fmt.Println(I_java_lang_String_String(arg1))
}

func I_java_io_PrintStream_print_G_V(arg0 *java_lang_Object, arg1 *java_lang_Object) {
	fmt.Print(I_java_lang_String_String(arg1))
}

func I_java_io_PrintStream_println__V(arg0 *java_lang_Object) {
	fmt.Println()
}

func I_java_io_PrintStream_printf_GRjava_lang_Object_java_io_PrintStream(arg0, arg1 *java_lang_Object, arg2 []*java_lang_Object) *java_lang_Object {
	format := I_java_lang_String_String(arg1)
	var args = make([]interface{}, len(arg2))
	for i, v := range arg2 { //TODO: actually implement
		switch v.name {
		case "java_lang_String":
			args[i] = I_java_lang_String_String(v)
		case "java_lang_Integer":
			args[i] = v.getFieldInt("E_val")
		case "java_lang_Double":
			args[i] = v.getFieldDouble("E_val")
		default:
			panic(v.name)
		}
	}
	fmt.Printf(format, args...)
	return arg0
}

func I_java_io_PrintStream_println_java_lang_Object_V(arg0, arg1 *java_lang_Object) {
	switch arg1.name {
	case "java_lang_String":
		fmt.Println(I_java_lang_String_String(arg1))
	case "java_lang_Integer":
		fmt.Println(arg1.getFieldInt("E_val"))
	case "java_lang_Double":
		fmt.Println(arg1.getFieldDouble("E_val"))
	case "java_lang_Character":
		fmt.Println(arg1.getFieldInt("E_val"))
	default:
		fmt.Println(arg1.name)
	}
}
