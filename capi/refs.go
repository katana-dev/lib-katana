package main

import (
	"C"
	"fmt"
	"sync"
)

//A single reference object.
type ref struct {
	obj interface{}
	cnt int32
}

//A singleton tracker for references.
var t struct {
	sync.Mutex
	next int32
	nums map[interface{}]int32
	objs map[int32]ref
}

//Init tracking features.
func initTracker() {
	if t.next != 0 {
		return
	}

	t.Lock()
	defer t.Unlock()
	t.next = 1 //Start at non-zero.
	t.nums = make(map[interface{}]int32)
	t.objs = make(map[int32]ref)
}

//Gets an object without manipulating the counter.
func getObj(n int32) interface{} {
	initTracker()

	t.Lock()
	r, ok := t.objs[n]
	t.Unlock()
	if !ok {
		panic(fmt.Sprintf("getObj unknown reference number: %d", n))
	}
	return r.obj
}

//Track a given Go object which is outbound.
func trackObj(o interface{}) int32 {
	initTracker()

	t.Lock()
	defer t.Unlock()

	//Try deduplicate the reference.
	n := t.nums[o]
	if n != 0 {
		r := t.objs[n]
		t.objs[n] = ref{r.obj, r.cnt + 1}
	} else {
		//Push a new reference.
		t.next++
		if t.next < 0 {
			panic("trackObj reference number overflow")
		}
		n = t.next
		t.nums[o] = n
		t.objs[n] = ref{o, 1}
	}

	return int32(n)
}

//Decrements the reference counter for given reference number, eventually deleting on zero.
func releaseRef(n int32) {
	initTracker()

	t.Lock()
	defer t.Unlock()

	r, ok := t.objs[n]
	if !ok {
		panic(fmt.Sprintf("releaseRef unknown reference number: %d", n))
	}

	//When this is the last reference, remove the entry from the maps.
	if r.cnt <= 1 {
		delete(t.objs, n)
		delete(t.nums, r.obj)
	} else {
		//Otherwise decrement.
		r := t.objs[n]
		t.objs[n] = ref{r.obj, r.cnt - 1}
	}
}

//Releases a reference generated before so it may be garbage collected.
//export release_ref
func release_ref(n C.int) {
	releaseRef(int32(n))
}
