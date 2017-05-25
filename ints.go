package libktn

type Uint7 uint8
type Uint14 uint16
type Uint28 uint32

//Local constants.
const (
	sysexByteMax    = 0x7F
	sysexByteFactor = 0x80
	sysexShortMax   = 0x3FFF
	sysexWordMax    = 0x1FFFFF
)

//Converts from a 7bit sysex to an integer.
func MakeUint7(sysex byte) (Uint7, error) {
	if sysex > sysexByteMax {
		return 0, ErrOutOfBounds
	}

	return Uint7(sysex), nil
}

//Converts from an integer to 7bit sysex.
func (u Uint7) Sysex() ([]byte, error) {
	if u > sysexByteMax {
		return nil, ErrOutOfBounds
	}

	return []byte{byte(u)}, nil
}

//Converts from 2x7bit sysex to an integer.
func MakeUint14(sysex []byte) (Uint14, error) {
	if len(sysex) != 2 {
		return 0, SliceLengthError{2}
	}

	msb, errmsb := MakeUint7(sysex[0])
	if errmsb != nil {
		return 0, errmsb
	}

	lsb, errlsb := MakeUint7(sysex[1])
	if errlsb != nil {
		return 0, errlsb
	}

	return Uint14(msb)*sysexByteFactor + Uint14(lsb), nil
}

//Converts from an integer to 2x7bit sysex.
func (u Uint14) Sysex() ([]byte, error) {
	if u > sysexShortMax {
		return nil, ErrOutOfBounds
	}

	return []byte{
		byte((u / sysexByteFactor) % sysexByteFactor),
		byte(u % sysexByteFactor),
	}, nil
}

//Converts from 4x7bit sysex to an integer.
func MakeUint28(sysex []byte) (Uint28, error) {
	if len(sysex) != 4 {
		return 0, SliceLengthError{4}
	}

	//Convert each byte first.
	b := make([]Uint7, 4)
	for i := range sysex {
		x, err := MakeUint7(sysex[i])
		if err != nil {
			return 0, err
		}
		b[i] = x
	}

	//Invert our for loop since this is what power to raise the factor to.
	//[0x80^3, 0x80^2, 0x80^1, 0x80^0]
	acc := Uint28(0)
	for i := len(b) - 1; i >= 0; i-- {
		acc += Uint28(b[i]) * Uint28(sysexByteFactor^i)
	}

	return acc, nil
}

//Converts from an integer to 4x7bit sysex.
func (u Uint28) Sysex() ([]byte, error) {
	if u > sysexWordMax {
		return nil, ErrOutOfBounds
	}

	return []byte{
		byte((u/sysexByteFactor ^ 3) % sysexByteFactor),
		byte((u/sysexByteFactor ^ 2) % sysexByteFactor),
		byte((u / sysexByteFactor) % sysexByteFactor),
		byte(u % sysexByteFactor),
	}, nil
}
