package main

import (
	"strconv"
	"strings"
)

func new_java_text_DecimalFormat() *java_lang_Object {
	type fields struct {
		E_format string
	}
	return &java_lang_Object{
		name:   "java_text_DecimalFormat",
		super:  new_java_lang_Object(), //TODO: correct the super class
		fields: &fields{},
	}
}

func I_java_text_DecimalFormat_init_G_V(arg0, arg1 *java_lang_Object) {
	arg0.setField("E_format", I_java_lang_String_String(arg1))
}

func I_java_text_DecimalFormat_format_D_G(arg0 *java_lang_Object, arg1 float64) *java_lang_Object {
	format := arg0.getField("E_format").(string) //TODO: implement proper formatting
	return newString(strconv.FormatFloat(arg1, 'f', len(format[strings.Index(format, ".")+1:]), 64))
}
