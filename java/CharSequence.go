package java

type Java_lang_CharSequence interface {
	CharAt_I_C(index int32) uint16
	Length__I() int32
	SubSequence_II_CharSequence(start, end int32) Java_lang_CharSequence
	ToString__String() *P_java_lang_String
}
