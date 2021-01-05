package main

func new_java_lang_Double() *java_lang_Object {
	type fields struct {
		E_val float64
	}
	return &java_lang_Object{
		name:   "java_lang_Double",
		super:  new_java_lang_Object(), // TODO: extend from java.lang.Number
		fields: &fields{},
	}
}

func java_lang_Double_valueOf_D_java_lang_Double(arg0 float64) *java_lang_Object {
	d := new_java_lang_Double()
	d.setField("E_val", arg0)
	return d
}
