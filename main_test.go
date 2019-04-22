package main

import (
	"testing"
	"unsafe"
)

var vals []int

func init() {
	for i := 0; i < 10000000; i++ {
		vals = append(vals, i)
	}
}

func BenchmarkRange(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		for _, i := range vals {
			bla := i
			bla = bla
		}
	}
}

func BenchmarkUnsafe(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		start := unsafe.Pointer(&vals[0])
		size := unsafe.Sizeof(int(0))
		for i := 0; i < len(vals); i++ {
			bla := *(*int)(unsafe.Pointer(uintptr(start) + size*uintptr(i)))
			bla = bla
		}
	}
}
