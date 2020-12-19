package java

import (
	"fmt"
)

type Java_io_PrintStream struct{}

func (this Java_io_PrintStream) Print_String_void(str *P_java_lang_String) {
	if str == nil {
		str = NewString_string_String_C("null")
	}
	this.write_String_void(str)
}

func (this Java_io_PrintStream) write_String_void(s *P_java_lang_String) {
	for i := int32(0); i < s.Length__I(); i++ {
		fmt.Print(string(rune(s.CharAt_I_C(i))))
	}
}
