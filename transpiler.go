package main

import (
	"encoding/json"
	"fmt"
	"github.com/totallygamerjet/jb4go/gen"
	"github.com/totallygamerjet/jb4go/parser"
	"github.com/totallygamerjet/jb4go/transformer"
	"io/ioutil"
	"log"
)

//https://www.mirkosertic.de/blog/2017/06/compiling-bytecode-to-javascript/
//https://tomassetti.me/how-to-write-a-transpiler/
func main() { //"Employee.class")
	b, err := ioutil.ReadFile("/Users/jarrettkuklis/eclipse-workspace/csis 312/src/gcd/gcdClass.class")
	if err != nil {
		log.Fatal("Couldn't open file: ", err)
	}
	raw, err := parser.ParseBytecode(b)
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
	//fmt.Println(class)
	err = gen.Transpile(class)
	if err != nil {
		log.Fatal("Failed to generator: ", err)
	}
}
