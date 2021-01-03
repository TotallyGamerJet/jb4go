package main

import (
	"fmt"
	"reflect"
	"sync"
)

type java_lang_Object struct {
	name       string                        // the java name of this class
	super      *java_lang_Object             // pointer to this object's super class
	sync.Mutex                               // Used to synchronize
	fields     interface{}                   // pointer to actual fields
	methods    func() map[string]interface{} // maps method names to func implementations
}

func fn_java_lang_Object() map[string]interface{} {
	return map[string]interface{}{
		"init__V": func(arg0 *java_lang_Object) {

		},
	}
}

func new_java_lang_Object() *java_lang_Object {
	return &java_lang_Object{
		name:    "java_lang_Object",
		methods: fn_java_lang_Object,
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

// callMethodObject is a convenience method that casts the result of callMethod to *java_lang_Object
func (arg0 *java_lang_Object) callMethodObject(methodName string, params ...interface{}) *java_lang_Object {
	if len(params) > 0 {
		return arg0.callMethod(methodName, params...).(*java_lang_Object)
	}
	return arg0.callMethod(methodName).(*java_lang_Object)
}

func (arg0 *java_lang_Object) callMethodDouble(methodName string, params ...interface{}) float64 {
	if len(params) > 0 {
		return arg0.callMethod(methodName, params...).(float64)
	}
	return arg0.callMethod(methodName).(float64)
}

func (arg0 *java_lang_Object) callMethodInt(methodName string, params ...interface{}) int32 {
	if len(params) > 0 {
		return arg0.callMethod(methodName, params...).(int32)
	}
	return arg0.callMethod(methodName).(int32)
}

// call method calls the method by a given name with any parameters and returns any results or nil if none
func (arg0 *java_lang_Object) callMethod(methodName string, params ...interface{}) interface{} {
	params = append([]interface{}{arg0}, params...) // always start with receiver
	var p = make([]reflect.Value, len(params))
	for i, v := range params {
		p[i] = reflect.ValueOf(v)
	}
	m, ok := arg0.methods()[methodName]
	if !ok {
		panic(fmt.Sprintf("%s has no method: %s", arg0.name, methodName))
	}
	rets := reflect.ValueOf(m).Call(p)
	if len(rets) > 0 {
		return rets[0].Interface()
	}
	return nil
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
