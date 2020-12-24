package main

import (
	"encoding/json"
	"fmt"
	"github.com/totallygamerjet/jb4go/gen"
	"github.com/totallygamerjet/jb4go/parser"
	"github.com/totallygamerjet/jb4go/transformer"
	"log"
	"os"
)

//https://www.mirkosertic.de/blog/2017/06/compiling-bytecode-to-javascript/
//https://tomassetti.me/how-to-write-a-transpiler/
func main() {
	f, err := os.Open("./examples/gcdClass.class") //TODO: handle .jar files
	if err != nil {
		log.Fatal(err)
	}
	raw, err := parser.Parse(f)
	if err != nil {
		log.Fatal("Couldn't parse java file: ", err)
	}
	class, err := transformer.Simplify(raw)
	if err != nil {
		log.Fatal("Couldn't simplify raw class: ", err)
	}
	j, err := json.MarshalIndent(class, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
	gFile, err := transformer.Translate(class)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gFile)
	err = gen.Generate(gFile, os.Stdout)
	if err != nil {
		log.Fatal("Failed to generate: ", err)
	}
}
