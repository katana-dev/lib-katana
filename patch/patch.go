package patch

import (
	"errors"

	libktn "github.com/katana-dev/lib-katana"
	"github.com/katana-dev/lib-katana/sysex"
)

const (
	EncSparse uint16 = 0
)

const (
	offFxChain = 928 //07 20
	lenFxChain = 20
)

var (
	ErrUnknownEncoding = errors.New("Unknown encoding flags")
	ErrDiscardedOffset = errors.New("Patch encoding discards this offset")
)

type WriteStat struct{ written, discarded libktn.Uint14 }

type Patch interface {
	GetFxChain() []libktn.Uint7
	GetByte(libktn.Uint14) (libktn.Uint7, error)
	GetShort(libktn.Uint14) (libktn.Uint14, error)
	WriteBytes(libktn.Uint14, []byte) (WriteStat, error)
	ApplyMessage(*sysex.SysexMessage) WriteStat
}

func New(enc uint16) (Patch, error) {
	switch enc {
	case EncSparse:
		return NewSparse(), nil
	default:
		return nil, ErrUnknownEncoding
	}
}
