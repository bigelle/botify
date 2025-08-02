package reused

import (
	"bytes"
	"sync"

	"github.com/go-playground/validator/v10"
)

var maxBufCap = 8 * 1024

var bufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 4*1024))
	},
}

func Buf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func PutBuf(b *bytes.Buffer) {
	if b.Cap() > maxBufCap {
		return
	}

	b.Reset()
	bufPool.Put(b)
}

var (
	valid *validator.Validate
	once  sync.Once
)

func Validator() *validator.Validate {
	once.Do(func() {
		valid = validator.New(validator.WithRequiredStructEnabled())
	})
	return valid
}
