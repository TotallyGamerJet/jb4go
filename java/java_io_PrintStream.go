package java

import (
	"fmt"
)

type P_java_io_PrintStream struct{}

func (this *P_java_io_PrintStream) P_print_G_V(str *P_java_lang_String) {
	if str == nil {
		str = NewString_string_G("null")
	}
	this._write_String_void(str)
}

func (arg0 *P_java_io_PrintStream) P_println_G_V(arg1 *P_java_lang_String) {
	arg0.P_print_G_V(arg1)
}

func (this *P_java_io_PrintStream) _write_String_void(s *P_java_lang_String) {
	for i := int32(0); i < s.P_length__I(); i++ {
		fmt.Print(string(rune(s.P_charAt_I_C(i))))
	}
}
