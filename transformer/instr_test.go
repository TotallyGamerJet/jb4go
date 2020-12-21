package transformer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_readInstructions(t *testing.T) {
	b := []byte{byte(iload_2), byte(ifge), 0x00, 0x27}
	i := readInstructions(b)
	assert.Equal(t, []instruction{
		{iload_2, nil},
		{opcode: ifge, operands: []byte{0x00, 0x27}}, {}, {},
	}, i)
}
