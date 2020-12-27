package main

import (
	"fmt"
	"github.com/totallygamerjet/jb4go/gen"
	"github.com/totallygamerjet/jb4go/parser"
	"github.com/totallygamerjet/jb4go/transformer"
	"log"
	"os"
)

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
	gFile, err := transformer.Translate(class)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gFile)
	o, err := os.OpenFile(gFile.FileName, os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = gen.Generate(gFile, o)
	if err != nil {
		log.Fatal("Failed to generate: ", err)
	}
}
