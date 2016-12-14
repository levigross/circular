# Circular Buffer
A lock free circular buffer in Go

[![Build Status](https://travis-ci.org/levigross/circular.svg?branch=master)](https://travis-ci.org/levigross/circular) [![GoDoc](https://godoc.org/github.com/levigross/circular?status.svg)](https://godoc.org/github.com/levigross/circular)

This package is written to be a lock free circular buffer. It stores the items as unsafe.Pointer 
objects. Usage is quite simple...

## How to download

`go get -u github.com/levigross/circular`

### Usage

```go
import "github.com/levigross/circular"


// This creates a buffer that is 100 elements large 
myBuffer := circular.Buffer(100)

for i := 0; i != 100; i++ {
		myInt := i // copy the int because we are storing unsafe.Pointers
		myBuf.Push(unsafe.Pointer(&myInt))
}

// Getting an item requires you to cast the object from an unsafe.Pointer
myint := *(*int)(myBuf.Pop())

```




