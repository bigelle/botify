package reused

import (
	"sync"
)

// or should it be bigger? idk
var bufMaxCap = 16 * 1024

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 4*1024)
		return &b
	},
}

func Buf() *[]byte {
	return bufPool.Get().(*[]byte)
}

func PutBuf(b *[]byte) {
	if cap(*b) > bufMaxCap {
		return
	}

	*b = (*b)[:0]
	bufPool.Put(b)
}
