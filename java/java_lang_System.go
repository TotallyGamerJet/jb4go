package java

import "os"

var java_lang_System_out = &java_io_PrintStream{}
var java_lang_System_in = &java_io_InputStream{input: os.Stdin}

type java_lang_System = struct {
	*java_lang_Object
}
