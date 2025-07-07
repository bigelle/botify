package reused

import (
	"bytes"
	"sync"
)

var maxBufCap = 8 * 1024

var bufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 4*1024))
	},
}

func Buf() *bytes.Buffer {
	return  bufPool.Get().(*bytes.Buffer)
}

func PutBuf(b *bytes.Buffer) {
	if b.Cap() > maxBufCap {
		return
	}

	b.Reset()
	bufPool.Put(b)
}
