package sysex

import (
	"bytes"
	"errors"
	"fmt"

	libktn "github.com/katana-dev/lib-katana"
)

var (
	ErrUnknownOp   = errors.New("Unknown SysEx operation")
	ErrBadHeader   = fmt.Errorf("Sysex message should start with 0x%x.", sysexStart)
	ErrBadFooter   = fmt.Errorf("Sysex message should end with 0x%x.", sysexEnd)
	ErrBadVendor   = errors.New("Sysex message should either have Roland or Universal Non-realtime vendor.")
	ErrBadUniSub   = errors.New("Sysex Universal Non-realtime subtype not supported.")
	ErrBadUniIdent = errors.New("Sysex Universal Non-realtime should identify as a Katana.")
	ErrBadModel    = errors.New("Sysex Roland message should have Katana model ID.")
	ErrBadRolandOp = fmt.Errorf("Sysex Roland message should have Query (0x%x) or Command (0x%x) operation.", queryFlag, commandFlag)
	ErrBadChecksum = errors.New("Sysex Roland message checksum doesn't match expected value.")
)

//Public constants for building messages.
const (
	OpIdRequest  = 1
	OpIdResponse = 2
	OpQuery      = 3
	OpCommand    = 4

	DevIdAny     = byte(0x7F)
	DevIdDefault = byte(0x00)
)

//Private values to check and serialize things.
var (
	modelId      = []byte{0x00, 0x00, 0x00, 0x33}
	familyCode   = []byte{0x33, 0x03, 0x00, 0x00}
	firmwareV102 = []byte{0x01, 0x00, 0x00, 0x00}
)

const (
	sysexStart = byte(0xF0)
	sysexEnd   = byte(0xF7)

	vendorId = byte(0x41)

	//Universal Non-realtime
	uniNonRt = byte(0x7E)
	uniInfo  = byte(0x06)
	uniIdReq = byte(0x01)
	uniIdRes = byte(0x02)

	//Roland sysex
	queryFlag   = byte(0x11)
	commandFlag = byte(0x12)
)

//Calculate the checksum byte for Sysex messages.
//Passing as many byte slices as you want.
func Checksum(blobs ...[]byte) byte {
	var acc uint32 = 0
	for _, blob := range blobs {
		for _, b := range blob {
			acc = (acc + uint32(b)) & 0x7f
		}
	}
	return byte((0x80 - acc) & 0x7f)
}

//Represents a single sysex message.
type SysexMessage struct {
	DeviceId, Op      byte
	FirmwareVer, Data []byte
	Address           Address
	Size              libktn.Uint28
}

//Creates a SysexMessage from a byte array.
//Be sure to include 0xF0 and 0xF7 header and footers.
func Parse(sysex []byte) (*SysexMessage, error) {
	//Check header.
	if sysex[0] != sysexStart {
		return nil, ErrBadHeader
	}

	//Check footer.
	if sysex[len(sysex)-1] != sysexEnd {
		return nil, ErrBadFooter
	}

	//Copy device ID.
	//TODO: Validate device ID.
	devId := sysex[2]

	//What message spec are we dealing with?
	switch sysex[1] {
	case uniNonRt:
		//Only the info one is supported.
		if sysex[3] != uniInfo {
			return nil, ErrBadUniSub
		}

		//See if it's a request or response.
		switch sysex[4] {
		case uniIdReq:
			//Requests don't have much data to check.
			return &SysexMessage{Op: OpIdRequest, DeviceId: devId}, nil

		case uniIdRes:
			//Responses are matched for Katana signature, firmeware version is not validated.
			if !matchBytes(sysex[5:10], []byte{vendorId}, familyCode) {
				return nil, ErrBadUniIdent
			}

			//Make sure we have our own copy.
			firmVer := make([]byte, 4)
			copy(firmVer, sysex[10:14])
			return &SysexMessage{Op: OpIdResponse, DeviceId: devId, FirmwareVer: firmVer}, nil

		default:
			return nil, ErrBadUniSub
		}

	case vendorId:
		//Check the model is a Katana.
		if !matchBytes(sysex[3:7], modelId) {
			return nil, ErrBadModel
		}

		//Addresses and checksum are the same for query and command, so do this first.
		a := sysex[8:12]
		addr, aerr := MakeAddress(a)
		d := sysex[12 : len(sysex)-2] //-1 for footer, -1 for checksum.
		c := Checksum(a, d)

		switch sysex[7] {
		case queryFlag:
			if aerr != nil {
				return nil, aerr
			}

			//Convert size bytes.
			size, serr := libktn.MakeUint28(d)
			if serr != nil {
				return nil, serr
			}

			//Create message either way.
			m := &SysexMessage{Op: OpQuery, DeviceId: devId, Address: addr, Size: size}

			//See if a checksum warning should be added.
			if sysex[len(sysex)-2] != c {
				return m, ErrBadChecksum
			}

			return m, nil

		case commandFlag:
			if aerr != nil {
				return nil, aerr
			}

			//Copy relevant data.
			data := make([]byte, len(d))
			copy(data, d)

			//Create message either way.
			m := &SysexMessage{Op: OpCommand, DeviceId: devId, Address: addr, Data: data}

			//See if a checksum warning should be added.
			if sysex[len(sysex)-2] != c {
				return m, ErrBadChecksum
			}

			return m, nil

		default:
			return nil, ErrBadRolandOp
		}

	default:
		return nil, ErrBadVendor
	}
}

//Factory for ID request sysex message.
func MakeIdRequest() SysexMessage {
	return SysexMessage{Op: OpIdRequest, DeviceId: DevIdAny}
}

//Factory for a query sysex message.
func MakeQuery(a Address, s libktn.Uint28) SysexMessage {
	return SysexMessage{Op: OpQuery, Address: a, Size: s, DeviceId: DevIdDefault}
}

//Factory for a command sysex message.
func MakeCommand(a Address, din []byte) SysexMessage {
	//Have our own copy of the slice, associated with this message.
	dmsg := make([]byte, len(din))
	copy(dmsg, din)

	return SysexMessage{Op: OpCommand, Address: a, Data: dmsg, DeviceId: DevIdDefault}
}

//Serializes a SysexMessage to bytes, as per Katana MIDI spec.
func (m *SysexMessage) Sysex() ([]byte, error) {
	switch m.Op {
	case OpIdRequest:
		return idRequest(m.DeviceId)
	case OpIdResponse:
		return idResponse(m.DeviceId, m.FirmwareVer)
	case OpQuery:
		return query(m.DeviceId, m.Address, m.Size)
	case OpCommand:
		return command(m.DeviceId, m.Address, m.Data)
	default:
		return nil, ErrUnknownOp
	}
}

//Internal serialize method.
func idRequest(deviceId byte) ([]byte, error) {
	//TODO: validate device ID.
	return []byte{sysexStart, uniNonRt, deviceId, uniInfo, uniIdReq, sysexEnd}, nil
}

//Internal serialize method.
func idResponse(deviceId byte, firmwareVer []byte) ([]byte, error) {
	//TODO: validate device ID.
	if firmwareVer == nil {
		return nil, libktn.RequiredError("FirmwareVer")
	}

	if len(firmwareVer) != 4 {
		return nil, libktn.SliceLengthError{4}
	}

	b := bytes.Buffer{}
	b.Grow(15)
	b.WriteByte(sysexStart)
	b.WriteByte(uniNonRt)
	b.WriteByte(deviceId)
	b.WriteByte(uniInfo)
	b.WriteByte(uniIdRes)
	b.WriteByte(vendorId)
	b.Write(familyCode)
	b.Write(firmwareVer)
	b.WriteByte(sysexEnd)
	return b.Bytes(), nil
}

//Internal serialize method.
func query(deviceId byte, addr Address, size libktn.Uint28) ([]byte, error) {
	//TODO: validate device ID.

	a, aerr := addr.Sysex()
	if aerr != nil {
		return nil, aerr
	}

	s, serr := size.Sysex()
	if serr != nil {
		return nil, serr
	}

	b := bytes.Buffer{}
	b.Grow(18)
	b.WriteByte(sysexStart)
	b.WriteByte(vendorId)
	b.WriteByte(deviceId)
	b.Write(modelId)
	b.WriteByte(queryFlag)
	b.Write(a)
	b.Write(s)
	b.WriteByte(Checksum(a, s))
	b.WriteByte(sysexEnd)
	return b.Bytes(), nil
}

//Internal serialize method.
func command(deviceId byte, addr Address, data []byte) ([]byte, error) {
	//TODO: validate device ID.

	a, aerr := addr.Sysex()
	if aerr != nil {
		return nil, aerr
	}

	b := bytes.Buffer{}
	b.Grow(14 + len(data))
	b.WriteByte(sysexStart)
	b.WriteByte(vendorId)
	b.WriteByte(deviceId)
	b.Write(modelId)
	b.WriteByte(commandFlag)
	b.Write(a)
	b.Write(data)
	b.WriteByte(Checksum(a, data))
	b.WriteByte(sysexEnd)
	return b.Bytes(), nil
}

func matchBytes(ref []byte, comp ...[]byte) bool {
	//Add up lengths.
	l := 0
	for _, b := range comp {
		l += len(b)
	}

	//They can't be equal if different.
	if len(ref) != l {
		return false
	}

	//Move across all slices and return if any byte doesn't match.
	i := 0
	for _, b := range comp {
		for _, x := range b {
			if ref[i] != x {
				return false
			}
			i++
		}
	}

	//Everything matched.
	return true
}
