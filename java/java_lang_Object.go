package main

import (
	"reflect"
	"sync"
)

type java_lang_Object struct {
	name       string            // the java name of this class
	super      *java_lang_Object // pointer to this object's super class
	sync.Mutex                   // Used to synchronize
	fields     interface{}       // pointer to actual fields
}

func I_java_lang_Object_init__V(arg0 *java_lang_Object) {

}

func new_java_lang_Object() *java_lang_Object {
	return &java_lang_Object{
		name: "java_lang_Object",
	}
}

// getFieldDouble is a convenience method for casting the result from getField to float64
func (arg0 *java_lang_Object) getFieldDouble(field string) float64 {
	return arg0.getField(field).(float64)
}

// getFieldInt is a convenience method for casting the result from getField to int32
func (arg0 *java_lang_Object) getFieldInt(field string) int32 {
	return arg0.getField(field).(int32)
}

// getFieldObject is a convenience method for casting the result from getField to *java_lang_Object
func (arg0 *java_lang_Object) getFieldObject(field string) *java_lang_Object {
	return arg0.getField(field).(*java_lang_Object)
}

// getField returns the field with a given name
func (arg0 *java_lang_Object) getField(field string) interface{} {
	return reflect.ValueOf(arg0.fields).Elem().FieldByName(field).Interface()
}

// setField takes in the field name and the value to assign it to
func (arg0 *java_lang_Object) setField(field string, val interface{}) {
	reflect.ValueOf(arg0.fields).Elem().FieldByName(field).Set(reflect.ValueOf(val))
}

func (arg0 *java_lang_Object) instanceof(obj *java_lang_Object) bool {
	if arg0.super != nil {
		if arg0.super.name == obj.name {
			return true
		}
		return arg0.super.instanceof(obj)
	}
	return false
}
