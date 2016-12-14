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
	"sync/atomic"
	"unsafe"
)

// Buffer is our circular buffer
type Buffer struct {
	read, write uint32
	data        []unsafe.Pointer
}

// NewBuffer allocates a new buffer
func NewBuffer(size uint32) *Buffer {
	b := &Buffer{data: make([]unsafe.Pointer, size)}
	return b
}

// Size is the size of the buffer
func (b Buffer) Size() uint32 {
	return atomic.LoadUint32(&b.write) - atomic.LoadUint32(&b.read)
}

// Empty will tell you if the buffer is empty
func (b Buffer) Empty() bool {
	return atomic.LoadUint32(&b.write) == atomic.LoadUint32(&b.read)
}

// Full returns true if the buffer is "full"
func (b Buffer) Full() bool {
	return b.Size() == uint32(len(b.data))
}

func (b Buffer) mask(val uint32) uint32 {
	return val % uint32(len(b.data))
}

// Push places an item onto the ring buffer
func (b *Buffer) Push(object unsafe.Pointer) {
	atomic.StorePointer(&b.data[b.mask(atomic.LoadUint32(&b.write))], object)
	atomic.AddUint32(&b.write, 1)
}

// Pop returns the next item on the ring buffer
func (b *Buffer) Pop() unsafe.Pointer {
	if b.Empty() {
		return nil
	}
	val := atomic.LoadPointer(&b.data[b.mask(atomic.LoadUint32(&b.read))])
	atomic.AddUint32(&b.read, 1)
	return val
}
