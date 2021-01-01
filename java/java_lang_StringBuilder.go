package main

import (
	"strconv"
	"strings"
)

func new_java_lang_StringBuilder() *java_lang_Object {
	type fields struct {
		E_b *strings.Builder
	}
	return &java_lang_Object{
		name:   "java_lang_StringBuilder",
		super:  new_java_lang_Object(),
		fields: &fields{},
		methods: map[string]interface{}{
			"init__V": func(arg0 *java_lang_Object) {
				arg0.setField("E_b", &strings.Builder{})
			},
			"append_G_java_lang_StringBuilder": func(arg0 *java_lang_Object, arg1 *java_lang_Object) *java_lang_Object {
				arg0.getField("E_b").(*strings.Builder).WriteString(arg1.callMethod("String").(string))
				return arg0
			},
			"append_I_java_lang_StringBuilder": func(arg0 *java_lang_Object, arg1 int32) *java_lang_Object {
				arg0.getField("E_b").(*strings.Builder).WriteString(strconv.Itoa(int(arg1)))
				return arg0
			},
			"toString__G": func(arg0 *java_lang_Object) *java_lang_Object {
				return newString(arg0.getField("E_b").(*strings.Builder).String())
			},
		},
	}
}
