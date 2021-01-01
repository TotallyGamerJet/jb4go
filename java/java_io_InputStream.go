package main

import (
	"io"
	"os"
)

func new_java_io_InputStream() *java_lang_Object {
	type fields struct {
		E_input io.Reader
	}
	return &java_lang_Object{
		name:    "java_io_InputStream",
		super:   new_java_lang_Object(),
		fields:  &fields{os.Stdin},
		methods: map[string]interface{}{},
	}
}
