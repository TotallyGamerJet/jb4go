// Package gen generates go files
package gen

import (
	_ "embed"
	"github.com/totallygamerjet/jb4go/transformer"
	"io"
	"text/template"
)

//go:embed go.gotmpl
var gen string

func Generate(g transformer.GoFile, w io.Writer) error {
	t, err := template.New("jb2go").Parse(gen)
	if err != nil {
		return err
	}
	return t.Execute(w, g)
}
