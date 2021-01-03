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
				format := arg1.callMethod("String").(string)
				var args = make([]interface{}, len(arg2))
				for i, v := range arg2 { //TODO: actually implement
					switch v.name {
					case "java_lang_String":
						args[i] = v.callMethod("String").(string)
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
			},
			"println_java_lang_Object_V": func(arg0, arg1 *java_lang_Object) {
				switch arg1.name {
				case "java_lang_String":
					fmt.Println(arg1.callMethod("String").(string))
				case "java_lang_Integer":
					fmt.Println(arg1.getFieldInt("E_val"))
				case "java_lang_Double":
					fmt.Println(arg1.getFieldDouble("E_val"))
				case "java_lang_Character":
					fmt.Println(arg1.getFieldInt("E_val"))
				default:
					fmt.Println(arg1.name)
				}
			},
		},
	}
}
