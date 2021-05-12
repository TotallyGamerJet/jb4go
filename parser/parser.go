package parser

import (
	"encoding/binary"
	"io"
)

type Parser struct {
	io.Reader
}

func (p *Parser) ReadU1() (v byte, err error) {
	var b = make([]byte, 1)
	_, err = p.Read(b)
	return b[0], err
}

func (p *Parser) ReadU2() (v uint16, err error) {
	var b = make([]byte, 2)
	_, err = p.Read(b)
	return binary.BigEndian.Uint16(b), err
}

func (p *Parser) ReadU4() (v uint32, err error) {
	var b = make([]byte, 4)
	_, err = p.Read(b)
	return binary.BigEndian.Uint32(b), err
}

func Parse(r io.Reader) *Parser {
	return &Parser{r}
}
