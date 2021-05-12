package main

import (
	"github.com/pkg/errors"
	"github.com/totallygamerjet/jb4go/gen"
	"github.com/totallygamerjet/jb4go/parser"
	"github.com/totallygamerjet/jb4go/transformer"
	"log"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	f, err := os.Open(args[1]) //TODO: handle .jar files
	if err != nil {
		return err
	}
	p := parser.Parse(f)
	raw, err := parser.ReadClass(p)
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
	o, err := os.OpenFile(gFile.FileName, os.O_CREATE|os.O_RDWR, 0775)
	if err != nil {
		return err
	}
	err = gen.Generate(gFile, o)
	if err != nil {
		return errors.Wrap(err, "Failed to generate: ")
	}
	return nil
}
