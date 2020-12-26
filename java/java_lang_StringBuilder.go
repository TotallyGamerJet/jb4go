package java

import (
	"strconv"
	"strings"
)

type java_lang_StringBuilder struct {
	*java_lang_Object
	b strings.Builder
}

func (arg0 *java_lang_StringBuilder) init__V() {
	arg0.b = strings.Builder{}
}

func (arg0 *java_lang_StringBuilder) append_G_java_lang_StringBuilder(arg1 *java_lang_String) *java_lang_StringBuilder {
	arg0.b.WriteString(arg1.String())
	return arg0
}

func (arg0 *java_lang_StringBuilder) append_I_java_lang_StringBuilder(arg1 int32) *java_lang_StringBuilder {
	arg0.b.WriteString(strconv.Itoa(int(arg1)))
	return arg0
}

func (arg0 *java_lang_StringBuilder) toString__G() *java_lang_String {
	return New_string_G(arg0.b.String())
}
