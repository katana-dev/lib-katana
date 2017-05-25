package patch

import (
	libktn "github.com/katana-dev/lib-katana"
	"github.com/katana-dev/lib-katana/sysex"
)

const (
	sparseCap                = 1040
	offMax     libktn.Uint14 = 2326
	offDiscard libktn.Uint14 = 0xFFFF
	padVal                   = 0
)

type boundary struct{ begin, end, shift libktn.Uint14 }

var bounds = [...]boundary{
	boundary{begin: 0, end: 107, shift: 0},
	boundary{begin: 192, end: 1059, shift: 85},
	boundary{begin: 2064, end: 2107, shift: 1090},
	boundary{begin: 2304, end: 2327, shift: 1287},
}

/*
The sparse patch implementation uses a simple offseting function.
Skipping the 3 largest unused/disabled sections in the patch addresses with the least CPU overhead.
Being just over 1KB per patch it isn't the most compressed but still suited for most situations.

Data loss in this encoding:
 - Parameters that have supported = 0 according to the TSL map may be discarded.

See https://github.com/katana-dev/docs/blob/master/data/tsl-map-1.0.0.csv
*/
type SparsePatch struct {
	data []byte
}

//Creates a new SparsePatch instance.
func NewSparse() Patch {
	return &SparsePatch{data: make([]byte, sparseCap, sparseCap)}
}

func (p *SparsePatch) ApplyMessage(msg *sysex.SysexMessage) WriteStat {
	if msg.Op == sysex.OpCommand && sysex.MutablePatchRegions[msg.Address.Region] {
		s, _ := p.WriteBytes(msg.Address.Offset, msg.Data)
		return s
	}
	return WriteStat{discarded: libktn.Uint14(len(msg.Data))}
}

func (p *SparsePatch) WriteBytes(offset libktn.Uint14, data []byte) (WriteStat, error) {
	var bi, di libktn.Uint14
	stat := WriteStat{}

	for di < libktn.Uint14(len(data)) {
		//Align to next boundary.
		for bounds[bi].end <= offset+di {
			bi++

			//Don't pass last boundary.
			if libktn.Uint14(len(bounds)) == bi {
				stat.discarded += libktn.Uint14(len(data)) - di
				return stat, nil
			}
		}

		//If we overshot we can skip data, or the whole thing.
		if bounds[bi].begin > offset+di {
			//Find out how much we can jump.
			//Cap this not to go beyond the provided array.
			dif := min14(bounds[bi].begin-(offset+di), libktn.Uint14(len(data))-di)

			stat.discarded += dif
			di += dif

			//Don't move past the end of our data.
			if di >= libktn.Uint14(len(data))-1 {
				return stat, nil
			}
		}

		//Copy till whichever ends first. The end of the boundary or the data.
		ds := offset + di - bounds[bi].shift
		de := bounds[bi].end - bounds[bi].shift
		c := libktn.Uint14(copy(p.data[ds:de], data[di:]))
		stat.written += c
		di += c
	}

	return stat, nil
}

func (p *SparsePatch) GetByte(offset libktn.Uint14) (libktn.Uint7, error) {
	if offset > offMax {
		return 0, libktn.ErrOutOfBounds
	}

	o := byteOffset(offset)
	if o == offDiscard {
		return 0, ErrDiscardedOffset
	}

	v, err := libktn.MakeUint7(p.data[o])
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (p *SparsePatch) GetShort(offset libktn.Uint14) (libktn.Uint14, error) {
	if offset > offMax {
		return 0, libktn.ErrOutOfBounds
	}

	omsb := byteOffset(offset)
	olsb := byteOffset(offset + 1)
	if omsb == offDiscard || olsb == offDiscard {
		return 0, ErrDiscardedOffset
	}

	v, err := libktn.MakeUint14([]byte{p.data[omsb], p.data[olsb]})
	if err != nil {
		return 0, err
	}
	return v, nil
}

func byteOffset(offset libktn.Uint14) libktn.Uint14 {
	for _, b := range bounds {
		//Try to move to the correct bound asap.
		if b.end < offset {
			continue
		}

		//Overshot a boundary, meaning it's a discard offset.
		if b.begin < offset {
			return offDiscard
		}

		//We're between the begin and end value of current bound. Do the shift.
		return offset - b.shift
	}

	//Beyond the last boundary is also discarded.
	return offDiscard
}

func min14(a, b libktn.Uint14) libktn.Uint14 {
	if a > b {
		return b
	}
	return a
}
