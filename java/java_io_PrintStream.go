package main

import "fmt"

func new_java_io_PrintStream() *java_lang_Object {
	type fields struct{}
	return &java_lang_Object{
		name:   "java_io_PrintStream",
		super:  new_java_lang_Object(),
		fields: &fields{},
		methods: map[string]interface{}{
			"println_G_V": func(arg0 *java_lang_Object, arg1 *java_lang_Object) {
				fmt.Println(arg1.callMethod("String").(string))
			},
			"print_G_V": func(arg0 *java_lang_Object, arg1 *java_lang_Object) {
				fmt.Print(arg1.callMethod("String").(string))
			},
			"println__V": func(arg0 *java_lang_Object) {
				fmt.Println()
			},
			"printf_GRjava_lang_Object_java_io_PrintStream": func(arg0, arg1 *java_lang_Object, arg2 []*java_lang_Object) *java_lang_Object {
				fmt.Println(arg1.callMethod("String").(string))
				for _, v := range arg2 { //TODO: actually implement
					switch v.name {
					case "java_lang_String":
						fmt.Print(v.callMethod("String").(string))
					default:
						fmt.Print(v.name, " ")
					}
				}
				fmt.Println()
				return arg0
			},
		},
	}
}
