package circular

import (
	"testing"
	"unsafe"
)

var myInt = 1

func BenchmarkPush(b *testing.B) {
	myBuf := NewBuffer(128)
	for i := 0; i < b.N; i++ {
		myBuf.Push(unsafe.Pointer(&myInt))
	}
}

func BenchmarkPop(b *testing.B) {
	myBuf := NewBuffer(128)
	for i := 0; i < b.N; i++ {
		myBuf.Push(unsafe.Pointer(&myInt))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myBuf.Pop()
	}
}
