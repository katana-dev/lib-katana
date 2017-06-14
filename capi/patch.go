package main

/*

typedef struct {
	void* data;
	int size;
} KtnSizedArray;

*/
import "C"

import (
	"fmt"
	"unsafe"

	libktn "github.com/katana-dev/lib-katana"
	"github.com/katana-dev/lib-katana/patch"
	"github.com/katana-dev/lib-katana/sysex"
)

//Creates an interger reference to a new Patch.
//export ktn_new_patch
func ktn_new_patch() C.int {
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
//export ktn_apply_message_to_patch
func ktn_apply_message_to_patch(n C.int, arr unsafe.Pointer, len C.int) *C.char {
	b := C.GoBytes(arr, len)
	m, err := sysex.Parse(b)
	if err != nil {
		return C.CString(err.Error())
	}
	p := getObj(int32(n)).(patch.Patch)
	r := p.ApplyMessage(m)
	return C.CString(fmt.Sprintf("Wrote to ref %v > %+v\n", n, r))
}

//export ktn_get_patch_byte_lossy
func ktn_get_patch_byte_lossy(n C.int, off C.ushort) C.uchar {
	p := getObj(int32(n)).(patch.Patch)
	o := libktn.Uint14(off)
	v, err := p.GetByte(o)
	if err != nil {
		return C.uchar(0)
	}
	return C.uchar(v)
}

//export ktn_get_patch_short_lossy
func ktn_get_patch_short_lossy(n C.int, off C.ushort) C.ushort {
	p := getObj(int32(n)).(patch.Patch)
	o := libktn.Uint14(off)
	v, err := p.GetShort(o)
	if err != nil {
		return C.ushort(0)
	}
	return C.ushort(v)
}

//export ktn_get_patch_fx_chain
func ktn_get_patch_fx_chain(n C.int) C.KtnSizedArray {
	p := getObj(int32(n)).(patch.Patch)
	v := p.GetFxChain()
	c := make([]byte, len(v))
	for i, b := range v {
		c[i] = byte(b)
	}
	return C.KtnSizedArray{data: C.CBytes(c), size: C.int(len(c))}
}

func main() {}
