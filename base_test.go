//    Copyright 2016 Levi Gross

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package circular

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestBasicBuffer(t *testing.T) {
	myBuf, err := NewBuffer(128)
	if err != nil {
		t.Fatal("Unable to create buffer", err)
	}
	if !myBuf.Empty() {
		t.Error("My empty buffer is not empty", myBuf.Size())
	}

	if myBuf.Full() {
		t.Error("Empty buffer is full", myBuf.Size())
	}

	if myBuf.Size() != 0 {
		t.Error("Buf size is not zero", myBuf.Size())
	}
}

func TestBufferOverCap(t *testing.T) {
	myBuf, err := NewBuffer(128)
	if err != nil {
		t.Fatal("Unable to create buffer", err)
	}
	for i := 0; i != 1000; i++ {
		myBuf.Push(unsafe.Pointer(&i))
	}
	for i := 0; i != 1000; i++ {
		if i%2 == 0 {
			myBuf.Push(unsafe.Pointer(&i))
		}
		myBuf.Pop()
	}
	if !myBuf.Empty() && myBuf.Size() == 0 {
		t.Error("Buffer not empty or size is more than 0", myBuf.Empty(), myBuf.Size())
	}
}

func TestBufferOps(t *testing.T) {
	myBuf, err := NewBuffer(128)
	if err != nil {
		t.Fatal("Unable to create buffer", err)
	}
	for i := 0; i != 128; i++ {
		myInt := i
		myBuf.Push(unsafe.Pointer(&myInt))
	}
	if !myBuf.Full() {
		t.Error("Buffer is full but it doesn't think it is", myBuf.Size())
	}

	for i := 0; i != 128; i++ {
		derVal := *(*int)(myBuf.Pop())
		if i != derVal {
			t.Error("Was expecting", i, "got", derVal)
		}
	}

	if !myBuf.Empty() {
		t.Error("Buffer isn't empty", myBuf.Size())
	}

}

type foo struct {
	count       int
	stringCount string
	derBytes    []byte
}

func TestConcurrentReadWrite(t *testing.T) {
	doneChan := make(chan struct{})
	myBuf, _ := NewBuffer(128)
	go func() {
		for {
			select {
			case <-doneChan:
				return
			default:
				myInt := 5436
				myBuf.Push(unsafe.Pointer(&myInt))
			}
		}
	}()
	go func() {
		for {
			anInt := 294
			select {
			case <-doneChan:
				return
			default:
				myBuf.Push(unsafe.Pointer(&anInt))
				myBuf.Pop()
			}
		}

	}()
	select {
	case <-time.After(time.Second):
		close(doneChan)
	}
}

func TestBufferCustomStruct(t *testing.T) {
	vals := make([]foo, 100)
	for i := range vals {
		vals[i].count = i
		vals[i].stringCount = fmt.Sprint(i)
		vals[i].derBytes = []byte(vals[i].stringCount + vals[i].stringCount)
	}
	myBuf, err := NewBuffer(128)
	if err != nil {
		t.Fatal("Unable to create buffer", err)
	}
	for i := range vals {
		myBuf.Push(unsafe.Pointer(&vals[i]))
	}

	if myBuf.Size() != 100 {
		t.Error("We size should be 100", myBuf.Size())
	}

	for i := range vals {
		derFoo := translateFoo(myBuf.Pop())
		if derFoo.stringCount != fmt.Sprint(i) {
			t.Error("Was expecting ", i, "got", derFoo.stringCount)
		}
		if derFoo.count != i {
			t.Error("Was expecting ", i, "got", derFoo.count)
		}
		if bytes.Compare(derFoo.derBytes,
			[]byte(vals[i].stringCount+vals[i].stringCount)) != 0 {
			t.Error("Was expecting",
				[]byte(vals[i].stringCount+vals[i].stringCount),
				"got", derFoo.derBytes)
		}
	}

}

func translateFoo(p unsafe.Pointer) foo {
	return *(*foo)(p)
}
