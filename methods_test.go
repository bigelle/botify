package botify

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// it might be dumb but i have to test it
func TestGetUpdates_ContentType(t *testing.T) {
	var gu GetUpdates

	assert.Equal(t, gu.ContentType(), "application/json")
}

func Test_multipartPayload(t *testing.T) {
	r, ct, err := multipartPayload([]multipartField{
		// text, string
		{
			Type:  "text",
			Name:  "text",
			Value: "text",
		},
		// text, string pointer
		{
			Type:  "text",
			Name:  "string pointer",
			Value: func(i string) *string { return &i }("text"),
		},
		// text, int
		{
			Type:  "text",
			Name:  "int",
			Value: 42,
		},
		// text, int pointer
		{
			Type:  "text",
			Name:  "int pointer",
			Value: func(i int) *int { return &i }(42),
		},
		// text, bool
		{
			Type:  "text",
			Name:  "bool",
			Value: true,
		},
		// text, bool pointer
		{
			Type:  "text",
			Name:  "bool pointer",
			Value: func(i bool) *bool { return &i }(true),
		},
		// text, struct
		{
			Type:  "text",
			Name:  "struct",
			Value: struct{ Foo string }{Foo: "text"},
		},
		// text, struct pointer
		{
			Type:  "text",
			Name:  "struct pointer",
			Value: &struct{ Foo string }{Foo: "text"},
		},
		// text, slice
		{
			Type:  "text",
			Name:  "slice",
			Value: []string{"foobar", "more text"},
		},
		// text, struct pointer
		{
			Type:  "text",
			Name:  "slice pointer",
			Value: &[]string{"foobar", "more text"},
		},
		// file, reader
		{
			Type:     "file",
			Name:     "file",
			FileName: "foo.txt",
			Value:    bytes.NewReader([]byte("goooaaaalllll")),
		},
	})

	if assert.NoError(t, err) {
		assert.NotEmpty(t, r)
		assert.NotEmpty(t, ct)

		// i need a better way to check if it was printed as it should've
		b, _ := io.ReadAll(r)
		t.Log(string(b))
	}
}
