package java

import "strconv"

type java_text_DecimalFormat struct {
	*java_lang_Object
	format string
}

func (arg0 java_text_DecimalFormat) init_G_V(arg1 *java_lang_String) {
	arg0.format = arg1.String()
}

func (arg0 java_text_DecimalFormat) format_D_G(arg1 float64) *java_lang_String {
	//TODO: actually handle the formatting properly
	return New_string_G(strconv.FormatFloat(arg1, 'f', 1, 64))
}
