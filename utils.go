package sb

import (
	"bytes"
	"sync"
)

var pool = &sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

func put(buffer *bytes.Buffer) {
	buffer.Reset()
	pool.Put(buffer)
}
