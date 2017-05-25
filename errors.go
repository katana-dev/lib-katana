package libktn

import (
	"errors"
	"fmt"
)

var (
	ErrOutOfBounds error = errors.New("Value is out of bounds")
)

type RequiredError string

func (e RequiredError) Error() string {
	return fmt.Sprintf("%s is a required field")
}

type SliceLengthError []int

func (e SliceLengthError) Error() string {
	switch len(e) {
	case 1:
		return fmt.Sprintf("Invalid slice length, expecting %d", e[0])
	case 2:
		return fmt.Sprintf("Invalid slice length, expecting between %d and %d", e[0], e[1])
	default:
		// Did you think I was going to produce more errors?
		return "Invalid slice length"
	}
}
