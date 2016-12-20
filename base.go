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
	"fmt"
	"sync/atomic"
	"unsafe"
)

// Buffer is our circular buffer
type Buffer struct {
	read, write uint64
	data        []unsafe.Pointer
}

// NewBuffer allocates a new buffer. This number needs to be a power of two
// or the buffer won't allocate.
func NewBuffer(size uint64) (*Buffer, error) {
	if size&(size-1) != 0 {
		return nil, fmt.Errorf("%d is not a power of two", size)
	}
	b := &Buffer{data: make([]unsafe.Pointer, size)}
	return b, nil
}

// Size is the size of the buffer
func (b Buffer) Size() uint64 {
	return atomic.LoadUint64(&b.write) - atomic.LoadUint64(&b.read)
}

// Empty will tell you if the buffer is empty
func (b Buffer) Empty() bool {
	return atomic.LoadUint64(&b.write) == atomic.LoadUint64(&b.read)
}

// Full returns true if the buffer is "full"
func (b Buffer) Full() bool {
	return b.Size() == uint64(len(b.data))
}

func (b Buffer) mask(val uint64) uint64 {
	return val & uint64(len(b.data)-1)
}

// Push places an item onto the ring buffer
func (b *Buffer) Push(object unsafe.Pointer) {
	for atomic.CompareAndSwapUint64(&b.write, atomic.LoadUint64(&b.write), atomic.LoadUint64(&b.write)+1) {
		atomic.StorePointer(&b.data[b.mask(atomic.LoadUint64(&b.write)-1)], object)
		break
	}
}

// Pop returns the next item on the ring buffer
func (b *Buffer) Pop() unsafe.Pointer {
	if atomic.LoadUint64(&b.write) == atomic.LoadUint64(&b.read) {
		return nil
	}
	var val unsafe.Pointer
	for atomic.CompareAndSwapPointer(&val, val, b.data[b.mask(atomic.LoadUint64(&b.read))]) {
		atomic.AddUint64(&b.read, 1)
		break
	}

	return val
}
