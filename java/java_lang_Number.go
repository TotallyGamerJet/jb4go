package main

var (
	java_lang_Number_serialVersionUID int64 = -8742448824652078965
)

func new_java_lang_Number() *java_lang_Object {
	type fields struct {
	}
	return &java_lang_Object{
		name:    "java_lang_Number",
		super:   new_java_lang_Object(),
		fields:  &fields{},
		methods: fn_java_lang_Number,
	}
}
func fn_java_lang_Number() map[string]interface{} {
	return map[string]interface{}{
		"init__V": func(aarg0 *java_lang_Object) {
			aarg0.super.callMethod("init__V")
			return
		},
		"intValue__I":    nil, // abstract method - implement in subclasses
		"longValue__J":   nil, // abstract method - implement in subclasses
		"floatValue__F":  nil, // abstract method - implement in subclasses
		"doubleValue__D": nil, // abstract method - implement in subclasses
		"byteValue__B": func(aarg0 *java_lang_Object) int32 {
			var v0 int32
			v0 = aarg0.callMethodInt("intValue__I")
			return int32(int8(v0))
		},
		"shortValue__S": func(aarg0 *java_lang_Object) int32 {
			var v0 int32
			v0 = aarg0.callMethodInt("intValue__I")
			return int32(int16(v0))
		},
	}
}
