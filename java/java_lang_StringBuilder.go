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
	}
}

func I_java_lang_StringBuilder_init__V(arg0 *java_lang_Object) {
	arg0.setField("E_b", &strings.Builder{})
}

func I_java_lang_StringBuilder_append_G_java_lang_StringBuilder(arg0 *java_lang_Object, arg1 *java_lang_Object) *java_lang_Object {
	arg0.getField("E_b").(*strings.Builder).WriteString(I_java_lang_String_String(arg1))
	return arg0
}

func I_java_lang_StringBuilder_append_I_java_lang_StringBuilder(arg0 *java_lang_Object, arg1 int32) *java_lang_Object {
	arg0.getField("E_b").(*strings.Builder).WriteString(strconv.Itoa(int(arg1)))
	return arg0
}

func I_java_lang_StringBuilder_append_D_java_lang_StringBuilder(arg0 *java_lang_Object, arg1 float64) *java_lang_Object {
	arg0.getField("E_b").(*strings.Builder).WriteString(strconv.FormatFloat(arg1, 'f', -1, 64))
	return arg0
}

func I_java_lang_StringBuilder_toString__G(arg0 *java_lang_Object) *java_lang_Object {
	return newString(arg0.getField("E_b").(*strings.Builder).String())
}
