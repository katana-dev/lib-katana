package main

/*
#include "common.h"
*/
import "C"
import (
	"unsafe"

	libktn "github.com/katana-dev/lib-katana"
	"github.com/katana-dev/lib-katana/sysex"
)

//export ktn_sysex_id_request
func ktn_sysex_id_request() unsafe.Pointer {
	m := sysex.MakeIdRequest()
	c, _ := m.Sysex()
	return NewCByteSlice(c)
}

//export ktn_sysex_query
func ktn_sysex_query(region, offset, size C.int) unsafe.Pointer {
	m := sysex.MakeQuery(sysex.Address{Region: libktn.Uint14(region), Offset: libktn.Uint14(offset)}, libktn.Uint28(size))
	c, _ := m.Sysex()
	return NewCByteSlice(c)
}

//export ktn_sysex_command
func ktn_sysex_command(region, offset C.int, arr unsafe.Pointer, size C.int) unsafe.Pointer {
	b := C.GoBytes(arr, size)
	m := sysex.MakeCommand(sysex.Address{Region: libktn.Uint14(region), Offset: libktn.Uint14(offset)}, b)
	c, _ := m.Sysex()
	return NewCByteSlice(c)
}
