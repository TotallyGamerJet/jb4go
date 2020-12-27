package gen

import (
	"github.com/totallygamerjet/jb4go/transformer"
	"io"
	"text/template"
)

func Generate(g transformer.GoFile, w io.Writer) error {
	t, err := template.ParseFiles("./gen/go.gotmpl")
	if err != nil {
		return err
	}
	return t.Execute(w, g)
}
