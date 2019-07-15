package packets

import "sync"

var _bufPool sync.Pool

type leakBuffer struct {
	b []byte
}

// NewBufferPool creates a leaky buffer which can hold at most n buffer, each
// with bufSize bytes.
func NewBufferPool(bufSize int) {
	_bufPool = sync.Pool{
		New: func() interface{} {
			return &leakBuffer{make([]byte, bufSize)}
		},
	}
}

// Getbuf returns a buffer from the leaky buffer or create a new buffer.
func Getbuf() (b *leakBuffer) {
	return _bufPool.Get().(*leakBuffer)
}

// Putbuf add the buffer into the free buffer pool for reuse. Panic if the buffer
// size is not the same with the leaky buffer's. This is intended to expose
// error usage of leaky buffer.
func Putbuf(b *leakBuffer) {
	_bufPool.Put(b)
}
