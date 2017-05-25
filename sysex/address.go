package sysex

import (
	libktn "github.com/katana-dev/lib-katana"
)

const (
	CH1Region   = 2049  //10 01
	CH2Region   = 2050  //10 02
	CH3Region   = 2051  //10 03
	CH4Region   = 2052  //10 04
	PanelRegion = 12288 //60 00
)

var (
	MutablePatchRegions = map[libktn.Uint14]bool{
		CH1Region:   true,
		CH2Region:   true,
		CH3Region:   true,
		CH4Region:   true,
		PanelRegion: true,
	}
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
