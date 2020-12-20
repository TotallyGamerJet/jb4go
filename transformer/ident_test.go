package transformer

import "fmt"

func ExampleTranslateIdent() {
	fmt.Println(TranslateIdent("byte"))
	fmt.Println(TranslateIdent("int"))
	fmt.Println(TranslateIdent("Ljava/lang/String"))

	//Output:
	// B
	// I
	// java_lang_String
}
