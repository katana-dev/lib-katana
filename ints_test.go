package libktn

import (
	"testing"
	"github.com/stvp/assert"
)

func TestUint7(t *testing.T) {
	var (
		valid = [][]byte{
			[]byte{0x00},
			[]byte{0x42},
			[]byte{0x7F},
		}
		oob = [][]byte{
			[]byte{0xFF},
		}
	)

	var (
		v Uint7
		s []byte
		e error
	)

	for _, in := range valid {
		v, e = MakeUint7(in[0])
		assert.Nil(t, e)
		s, e = v.Sysex()
		assert.Nil(t, e)
		assert.Equal(t, in, s)
	}

	for _, in := range oob {
		v, e = MakeUint7(in[0])
		assert.Equal(t, ErrOutOfBounds, e)
		assert.Equal(t, Uint7(0), v)

		s, e = Uint7(in[0]).Sysex()
		assert.Equal(t, ErrOutOfBounds, e)
		assert.Nil(t, s)
	}
}

func TestMakeUint14(t *testing.T) {
	var (
		valid = [][]byte{
			[]byte{0x00, 0x00},
			[]byte{0x00, 0x42},
			[]byte{0x42, 0x42},
			[]byte{0x7F, 0x7F},
		}
		oob = [][]byte{
			[]byte{0x7F, 0xFF},
			[]byte{0xFF, 0x7F},
		}
	)

	var (
		v Uint14
		s []byte
		e error
	)

	for _, in := range valid {
		v, e = MakeUint14(in)
		assert.Nil(t, e)
		s, e = v.Sysex()
		assert.Nil(t, e)
		assert.Equal(t, in, s)
	}

	for _, in := range oob {
		v, e = MakeUint14(in)
		assert.Equal(t, ErrOutOfBounds, e)
		assert.Equal(t, Uint14(0), v)

		s, e = Uint14(0xFFFF).Sysex()
		assert.Equal(t, ErrOutOfBounds, e)
		assert.Nil(t, s)
	}

	v, e = MakeUint14([]byte{1, 2, 3})
	assert.Equal(t, SliceLengthError{2}, e)
	assert.Equal(t, Uint14(0), v)
}
