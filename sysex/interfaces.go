package sysex

type Serializable interface {
	Sysex() ([]byte, error)
}
