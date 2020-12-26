package java

type java_lang_CharSequence interface {
	charAt_I_C(index int32) uint16
	length__I() int32
	subSequence_II_java_lang_CharSequence(start, end int32) java_lang_CharSequence
	toString__G() *java_lang_String
}
