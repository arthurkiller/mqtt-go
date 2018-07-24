package packets

import "sync"

var _bufPool sync.Pool

// NewLeakyBuf creates a leaky buffer which can hold at most n buffer, each
// with bufSize bytes.
func NewBufferPool(bufSize int) {
	_bufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, bufSize)
		},
	}
}

// Get returns a buffer from the leaky buffer or create a new buffer.
func Getbuf() (b []byte) {
	return _bufPool.Get().([]byte)
}

// Put add the buffer into the free buffer pool for reuse. Panic if the buffer
// size is not the same with the leaky buffer's. This is intended to expose
// error usage of leaky buffer.
func Putbuf(b []byte) {
	_bufPool.Put(b)
}
