package sysex

import (
	"testing"
	"github.com/stvp/assert"

	libktn "github.com/katana-dev/lib-katana"
)

func TestAddress(t *testing.T) {
	var (
		valid = [][]byte{
			[]byte{0x00, 0x00, 0x00, 0x00},
			[]byte{0x42, 0x42, 0x42, 0x42},
			[]byte{0x7F, 0x00, 0x00, 0x00},
			[]byte{0x00, 0x7F, 0x00, 0x00},
			[]byte{0x00, 0x00, 0x7F, 0x00},
			[]byte{0x00, 0x00, 0x00, 0x7F},
			[]byte{0x7F, 0x7F, 0x7F, 0x7F},
		}
		oob = [][]byte{
			[]byte{0x00, 0x00, 0x00, 0xFF},
			[]byte{0x00, 0x00, 0xFF, 0x00},
			[]byte{0x00, 0xFF, 0x00, 0x00},
			[]byte{0xFF, 0x00, 0x00, 0x00},
		}
	)

	var (
		a Address
		b []byte
		e error
	)

	for _, in := range valid {
		a, e = MakeAddress(in)
		assert.Nil(t, e)
		b, e = a.Sysex()
		assert.Nil(t, e)
		assert.Equal(t, in, b)
	}

	for _, in := range oob {
		a, e = MakeAddress(in)
		assert.Equal(t, Address{}, a)
		assert.Equal(t, libktn.ErrOutOfBounds, e)
	}

	a, e = MakeAddress([]byte{1, 2, 3})
	assert.Equal(t, Address{}, a)
	assert.Equal(t, libktn.SliceLengthError{4}, e)
}
