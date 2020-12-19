package java

type P_java_lang_String struct {
	Java_lang_Object
	value []uint16
	hash  int32
}

func NewString_string_String_C(str string) (this *P_java_lang_String) {
	this = &P_java_lang_String{}
	this.value = make([]uint16, len(str))
	for i := 0; i < len(str); i++ {
		this.value[i] = uint16(str[i])
	}
	return this
}

func NewString_RCII_String_C(value []uint16, offset, count int32) (this *P_java_lang_String) {
	this = &P_java_lang_String{}
	if offset < 0 {
		//TODO: throw new StringIndexOutOfBoundsException(offset);
	}
	if count < 0 {
		//TODO: throw new StringIndexOutOfBoundsException(count);
	}
	// Note: offset or count might be near -1>>>1.
	if offset > int32(len(value))-count {
		//TODO: throw new StringIndexOutOfBoundsException(offset + count);
	}
	this.value = Arrays_CopyOfRange_RCII_RC(value, offset, offset+count)
	return this
}

func (this *P_java_lang_String) ToString__String() *P_java_lang_String {
	return this
}

func (this *P_java_lang_String) SubSequence_II_CharSequence(start, end int32) Java_lang_CharSequence {
	return this.Substring_II_String(start, end)
}

func (this *P_java_lang_String) Substring_II_String(beginIndex, endIndex int32) *P_java_lang_String {
	if beginIndex < 0 {
		//TODO: throw new StringIndexOutOfBoundsException(beginIndex);
	}
	if endIndex > int32(len(this.value)) {
		//TODO: throw new StringIndexOutOfBoundsException(endIndex);
	}
	subLen := endIndex - beginIndex
	if subLen < 0 {
		//TODO: throw new StringIndexOutOfBoundsException(subLen);
	}
	if (beginIndex == 0) && (endIndex == int32(len(this.value))) {
		return this
	}
	return NewString_RCII_String_C(this.value, beginIndex, subLen)
}

func (this *P_java_lang_String) Length__I() int32 {
	return int32(len(this.value))
}

func (this *P_java_lang_String) CharAt_I_C(index int32) uint16 {
	if (index < 0) || (index >= int32(len(this.value))) {
		//TODO: throw new*StringIndexOutOfBoundsException(index);
	}
	return this.value[index]
}

func (this *P_java_lang_String) IndexOf_StringI_I(str *P_java_lang_String, fromIndex int32) int32 {
	return _String_IndexOf_CIIRCIII_I(this.value, 0, int32(len(this.value)), str.value, 0, int32(len(str.value)), fromIndex)
}

func (this *P_java_lang_String) IndexOf_String_I(str *P_java_lang_String) int32 {
	return this.IndexOf_StringI_I(str, 0)
}

func (this *P_java_lang_String) Contains_CharSequence_bool(s Java_lang_CharSequence) bool {
	return this.IndexOf_String_I(s.ToString__String()) > -1
}

func _String_IndexOf_CIIRCIII_I(source []uint16, sourceOffset, sourceCount int32, target []uint16, targetOffset, targetCount, fromIndex int32) int32 {
	if fromIndex >= sourceCount {
		if targetCount == 0 {
			return sourceCount
		} else {
			return -1
		}
	}
	if fromIndex < 0 {
		fromIndex = 0
	}
	if targetCount == 0 {
		return fromIndex
	}

	first := target[targetOffset]
	max := sourceOffset + (sourceCount - targetCount)

	for i := sourceOffset + fromIndex; i <= max; i++ {
		// Look for first character.
		if source[i] != first {
			for { //while (++i <= max && source[i] != first);
				i++
				if i > max || source[i] == first {
					break
				}
			}
		}

		// Found first character, now look at the rest of v2
		if i <= max {
			j := i + 1
			end := j + targetCount - 1
			for k := targetOffset + 1; j < end && source[j] == target[k]; j, k = j+1, k+1 {
			}

			if j == end {
				// Found whole string.
				return i - sourceOffset
			}
		}
	}
	return -1
}
