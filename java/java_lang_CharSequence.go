package java

type P_java_lang_CharSequence interface {
	P_charAt_I_C(index int32) uint16
	P_length__I() int32
	P_subSequence_II_java_lang_CharSequence(start, end int32) P_java_lang_CharSequence
	P_toString__G() *P_java_lang_String
}
