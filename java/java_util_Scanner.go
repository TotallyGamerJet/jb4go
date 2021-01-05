package main

import (
	"bufio"
	"io"
	"strconv"
)

func new_java_util_Scanner() *java_lang_Object {
	type fields struct {
		E_r *bufio.Reader
	}
	return &java_lang_Object{
		name:   "java_util_Scanner",
		super:  new_java_lang_Object(),
		fields: &fields{},
	}
}

func I_java_util_Scanner_init_java_io_InputStream_V(arg0 *java_lang_Object, arg1 *java_lang_Object) {
	arg0.setField("E_r", bufio.NewReader(arg1.getField("E_input").(io.Reader)))
}

func I_java_util_Scanner_nextInt__I(arg0 *java_lang_Object) int32 {
	l, _, err := arg0.getField("E_r").(*bufio.Reader).ReadLine() //TODO: leave new line on the stream
	if err != nil {
		panic(err)
	}
	i, err := strconv.Atoi(string(l))
	if err != nil {
		panic(err)
	}
	return int32(i)
}

func I_java_util_Scanner_nextLine__G(arg0 *java_lang_Object) *java_lang_Object {
	l, _, err := arg0.getField("E_r").(*bufio.Reader).ReadLine() //TODO: leave new line on the stream
	if err != nil {
		panic(err)
	}
	return newString(string(l))
}

func I_java_util_Scanner_nextDouble__D(arg0 *java_lang_Object) float64 {
	l, _, err := arg0.getField("E_r").(*bufio.Reader).ReadLine()
	if err != nil {
		panic(err)
	}
	f, err := strconv.ParseFloat(string(l), 64)
	if err != nil {
		panic(err)
	}
	return f
}
