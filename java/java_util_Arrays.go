package java

func P_java_util_Arrays_copyOfRange_RCII_RC(original []uint16, from, to int32) []uint16 {
	newLength := to - from
	if newLength < 0 {
		//TODO: throw new IllegalArgumentException(from + " > " + to);
	}
	cpy := make([]uint16, newLength)
	n := copy(cpy, original[from:to])
	if int32(n) != newLength {
		panic("improper")
	}
	/*char[] copy = new char[newLength];
	System.arraycopy(original, from, copy, 0, Math.min(original.length - from, newLength));*/
	return cpy
}
