package main

import "C"

import (
	"fmt"
	"unsafe"

	"github.com/katana-dev/lib-katana/patch"
	"github.com/katana-dev/lib-katana/sysex"
)

//Creates an interger reference to a new Patch.
//export new_patch
func new_patch() C.int {
	p, _ := patch.New(patch.EncSparse)
	n := trackObj(p)
	return C.int(n)
}

/**
 * Applies a sysex message to a patch.
 *
 * @param int Reference number
 * @param void* Byte array pointer
 * @param int Array length
 * @return char* CString message
 */
//export apply_message_to_patch
func apply_message_to_patch(n C.int, arr unsafe.Pointer, len C.int) *C.char {
	b := C.GoBytes(arr, len)
	m, err := sysex.Parse(b)
	if err != nil {
		return C.CString(err.Error())
	}
	p := getObj(int32(n)).(patch.Patch)
	r := p.ApplyMessage(m)
	return C.CString(fmt.Sprintf("Wrote to ref %v > %+v\n", n, r))
}

func main() {}
