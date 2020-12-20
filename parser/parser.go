package parser

import (
	"encoding/binary"
	errors2 "github.com/pkg/errors"
	"io"
	"io/ioutil"
)

type Parser struct {
	data []byte
	pos  int
}

func (p *Parser) Read(b []byte) (n int) {
	n = copy(b, p.data[p.pos:])
	p.pos += n
	return n
}

func (p *Parser) ReadU1() (v byte) {
	v = p.data[p.pos]
	p.pos++
	return v
}

func (p *Parser) ReadU2() (v uint16) {
	v = binary.BigEndian.Uint16(p.data[p.pos:])
	p.pos += 2
	return v
}

func (p *Parser) ReadU4() (v uint32) {
	v = binary.BigEndian.Uint32(p.data[p.pos:])
	p.pos += 4
	return v
}

func ParseBytecode(java []byte) (RawClass, error) {
	p := &Parser{java, 0}
	return ReadClass(p)
}

func Parse(r io.Reader) (RawClass, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return RawClass{}, errors2.Wrap(err, "Unable to read Reader")
	}
	return ParseBytecode(b)
}
