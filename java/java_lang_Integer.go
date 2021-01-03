package main

func fn_java_lang_Integer() map[string]interface{} {
	return map[string]interface{}{}
}

func new_java_lang_Integer() *java_lang_Object {
	type fields struct {
		E_val int32
	}
	return &java_lang_Object{
		name:    "java_lang_Integer",
		super:   new_java_lang_Object(), // TODO: extend from java.lang.Number
		fields:  &fields{},
		methods: fn_java_lang_Integer,
	}
}

func java_lang_Integer_valueOf_I_java_lang_Integer(arg0 int32) *java_lang_Object {
	d := new_java_lang_Integer()
	d.setField("E_val", arg0)
	return d
}
