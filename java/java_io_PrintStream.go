package java

import (
	"fmt"
)

type java_io_PrintStream struct{}

func (this *java_io_PrintStream) print_G_V(str *java_lang_String) {
	if str == nil {
		str = New_string_G("null")
	}
	this._write_String_void(str)
}

func (arg0 *java_io_PrintStream) println_G_V(arg1 *java_lang_String) {
	arg0.print_G_V(arg1)
	arg0.print_G_V(New_string_G("\n"))
}

func (this *java_io_PrintStream) _write_String_void(s *java_lang_String) {
	for i := int32(0); i < s.length__I(); i++ {
		fmt.Print(string(rune(s.charAt_I_C(i))))
	}
}
