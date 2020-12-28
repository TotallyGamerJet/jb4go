package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/totallygamerjet/jb4go/gen"
	"github.com/totallygamerjet/jb4go/parser"
	"github.com/totallygamerjet/jb4go/transformer"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	f, err := os.Open(os.Args[1]) //TODO: handle .jar files
	if err != nil {
		return err
	}
	raw, err := parser.Parse(f)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse java file: ")
	}
	class, err := transformer.Simplify(raw)
	if err != nil {
		return errors.Wrap(err, "Couldn't simplify raw class: ")
	}
	gFile, err := transformer.Translate(class)
	if err != nil {
		return err
	}
	fmt.Println(gFile)
	o, err := os.OpenFile(gFile.FileName, os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	err = gen.Generate(gFile, o)
	if err != nil {
		return errors.Wrap(err, "Failed to generate: ")
	}
	return nil
}
