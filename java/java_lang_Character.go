package main

func new_java_lang_Character() *java_lang_Object {
	type fields struct {
		E_val uint16
	}
	return &java_lang_Object{
		name:    "java_lang_Character",
		super:   new_java_lang_Object(), // TODO: extend from java.lang.Number
		fields:  &fields{},
		methods: fn_java_lang_Character,
	}
}

func fn_java_lang_Character() map[string]interface{} {
	return map[string]interface{}{}
}

func java_lang_Character_valueOf_C_java_lang_Character(arg0 int32) *java_lang_Object {
	d := new_java_lang_Character()
	d.setField("E_val", uint16(arg0))
	return d
}
