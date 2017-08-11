package sysex

import (
	"testing"

	"github.com/stvp/assert"

	libktn "github.com/katana-dev/lib-katana"
)

//TODO: Test Sysex methods

func TestParseIdRequest(t *testing.T) {
	var (
		valid = map[[6]byte]SysexMessage{
			[6]byte{0xF0, 0x7E, 0x7F, 0x06, 0x01, 0xF7}: SysexMessage{Op: OpIdRequest, DeviceId: 0x7F},
			[6]byte{0xF0, 0x7E, 0x03, 0x06, 0x01, 0xF7}: SysexMessage{Op: OpIdRequest, DeviceId: 0x03},
		}
	)

	var (
		m *SysexMessage
		e error
	)

	for in, exp := range valid {
		m, e = Parse(in[:])
		assert.Nil(t, e)
		assert.Equal(t, exp, *m)
	}
}

func TestMakeIdRequest(t *testing.T) {
	m := MakeIdRequest()
	assert.Equal(t, SysexMessage{Op: OpIdRequest, DeviceId: 0x7F}, m)
}

func TestParseIdResponse(t *testing.T) {
	var (
		valid = map[[15]byte]SysexMessage{
			[15]byte{0xF0, 0x7E, 0x7F, 0x06, 0x02, 0x41, 0x33, 0x03, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xF7}: SysexMessage{Op: OpIdResponse, DeviceId: 0x7F, FirmwareVer: []byte{0x01, 0x00, 0x00, 0x00}},
			[15]byte{0xF0, 0x7E, 0x04, 0x06, 0x02, 0x41, 0x33, 0x03, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0xF7}: SysexMessage{Op: OpIdResponse, DeviceId: 0x04, FirmwareVer: []byte{0x01, 0x02, 0x03, 0x04}},
		}
	)

	var (
		m *SysexMessage
		e error
	)

	for in, exp := range valid {
		m, e = Parse(in[:])
		assert.Nil(t, e)
		assert.Equal(t, exp, *m)
	}
}

func TestParseQuery(t *testing.T) {
	var (
		valid = map[[18]byte]SysexMessage{
			[18]byte{0xF0, 0x41, 0x00, 0x00, 0x00, 0x00, 0x33, 0x11, 0x60, 0x00, 0x00, 0x53, 0x00, 0x00, 0x00, 0x01, 0x4C, 0xF7}: SysexMessage{Op: OpQuery, DeviceId: 0x00, Address: Address{Region: PanelRegion, Offset: 0x53}, Size: 1},
			[18]byte{0xF0, 0x41, 0x02, 0x00, 0x00, 0x00, 0x33, 0x11, 0x10, 0x02, 0x00, 0x06, 0x00, 0x00, 0x00, 0x06, 0x62, 0xF7}: SysexMessage{Op: OpQuery, DeviceId: 0x02, Address: Address{Region: CH2Region, Offset: 0x06}, Size: 6},
		}
	)

	var (
		m *SysexMessage
		e error
	)

	for in, exp := range valid {
		m, e = Parse(in[:])
		assert.Nil(t, e)
		assert.Equal(t, exp, *m)
	}
}

func TestMakeQuery(t *testing.T) {
	a := Address{Region: PanelRegion, Offset: 0x42}
	s := libktn.Uint28(1337)
	m := MakeQuery(a, s)
	assert.Equal(t, SysexMessage{Op: OpQuery, Address: a, Size: s, DeviceId: 0x00}, m)
}

func TestParseCommand(t *testing.T) {
	var (
		valid = map[[18]byte]SysexMessage{
			[18]byte{0xF0, 0x41, 0x00, 0x00, 0x00, 0x00, 0x33, 0x12, 0x60, 0x00, 0x00, 0x00, 0x4B, 0x41, 0x54, 0x41, 0x7F, 0xF7}: SysexMessage{Op: OpCommand, DeviceId: 0x00, Address: Address{Region: PanelRegion, Offset: 0x00}, Data: []byte{0x4B, 0x41, 0x54, 0x41}},
			[18]byte{0xF0, 0x41, 0x02, 0x00, 0x00, 0x00, 0x33, 0x12, 0x10, 0x02, 0x00, 0x06, 0x20, 0x20, 0x20, 0x20, 0x68, 0xF7}: SysexMessage{Op: OpCommand, DeviceId: 0x02, Address: Address{Region: CH2Region, Offset: 0x06}, Data: []byte{0x20, 0x20, 0x20, 0x20}},
		}
	)

	var (
		m *SysexMessage
		e error
	)

	for in, exp := range valid {
		m, e = Parse(in[:])
		assert.Nil(t, e)
		assert.Equal(t, exp, *m)
	}
}

func TestMakeCommand(t *testing.T) {
	a := Address{Region: PanelRegion, Offset: 0x42}
	d := []byte{0, 1, 2, 3, 4}
	m := MakeCommand(a, d)
	assert.Equal(t, SysexMessage{Op: OpCommand, Address: a, Data: d, DeviceId: 0x00}, m)
}
