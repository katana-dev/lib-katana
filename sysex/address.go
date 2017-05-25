package sysex

import (
	libktn "github.com/katana-dev/lib-katana"
)

//Represents an address in sysex memory.
type Address struct {
	Region libktn.Uint14
	Offset libktn.Uint14
}

//Converts from 4x7bit sysex address to an Address struct.
func MakeAddress(sysex []byte) (Address, error) {
	if len(sysex) != 4 {
		return Address{}, libktn.SliceLengthError{4}
	}

	region, errRegion := libktn.MakeUint14(sysex[:2])
	if errRegion != nil {
		return Address{}, errRegion
	}

	offset, errOffset := libktn.MakeUint14(sysex[2:])
	if errOffset != nil {
		return Address{}, errOffset
	}

	return Address{Region: region, Offset: offset}, nil
}

//Converts from an Address struct to 4x7bit sysex address.
func (a *Address) Sysex() ([]byte, error) {
	region, errRegion := a.Region.Sysex()
	if errRegion != nil {
		return nil, errRegion
	}

	offset, errOffset := a.Offset.Sysex()
	if errOffset != nil {
		return nil, errOffset
	}

	return append(region, offset...), nil
}
