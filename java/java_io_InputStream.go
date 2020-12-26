package java

import (
	"io"
)

type java_io_InputStream struct {
	*java_lang_Object
	input io.Reader
}
