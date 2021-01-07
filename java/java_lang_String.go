package main

func new_java_lang_String() *java_lang_Object {
	type fields struct {
		E_str string
	}
	return &java_lang_Object{
		name:    "java_lang_String",
		super:   new_java_lang_Object(),
		fields:  &fields{},
		methods: fn_java_lang_String,
	}
}

func fn_java_lang_String() map[string]interface{} {
	return map[string]interface{}{
		"String": func(arg0 *java_lang_Object) string {
			return arg0.getField("E_str").(string)
		},
	}
}

func newString(str string) *java_lang_Object {
	o := new_java_lang_String()
	o.setField("E_str", str)
	return o
}

func java_lang_String_format_GRjava_lang_Object_G(format *java_lang_Object, objs []*java_lang_Object) *java_lang_Object {
	return newString("(TODO)")
}
