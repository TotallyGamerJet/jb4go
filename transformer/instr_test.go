package transformer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_readInstructions(t *testing.T) {
	b := []byte{byte(iload_2), byte(ifge), 0x00, 0x27}
	i := readInstructions(b)
	assert.Equal(t, []instruction{
		{Loc: 0, Op: iload_2, operands: nil},
		{Loc: 1, Op: ifge, operands: []byte{0x00, 0x27}}, {}, {},
	}, i)
}
