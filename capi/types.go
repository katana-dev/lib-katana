package main

/*
#include <stdlib.h>
#include "common.h"

void* mallocByteSlice(void* data, int size){
	ByteSlice* b = malloc(sizeof(ByteSlice));
	if(b == NULL)
		return b;

	b->data = data;
	b->size = size;
	return b;
}
*/
import "C"

import (
	"unsafe"
)

func NewCByteSlice(c []byte) unsafe.Pointer {
	return C.mallocByteSlice(C.CBytes(c), C.int(len(c)))
}
