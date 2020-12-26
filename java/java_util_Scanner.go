package java

import (
	"bufio"
	"strconv"
)

type java_util_Scanner struct {
	*java_lang_Object
	r *bufio.Reader
}

func (arg0 *java_util_Scanner) init_java_io_InputStream_V(arg1 *java_io_InputStream) {
	arg0.r = bufio.NewReader(arg1.input)
}

func (arg0 *java_util_Scanner) nextInt__I() int32 {
	l, _, err := arg0.r.ReadLine()
	if err != nil {
		panic(err)
	}
	i, err := strconv.Atoi(string(l))
	if err != nil {
		panic(err)
	}
	return int32(i)
}
