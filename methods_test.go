package botify_test

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/bigelle/botify"
	"github.com/stretchr/testify/assert"
)

func TestSetWebhook_WritePayload(t *testing.T) {
	swh := botify.SetWebhook{
		URL:         "https://example.com/webhook",
		Certificate: botify.InputFileLocal{Name: "cert", Data: strings.NewReader("some ultra secret certificate")},
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	ct, err := swh.WritePayload(buf)

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.NotEmpty(t, ct) {
		t.FailNow()
	}

	_, params, err := mime.ParseMediaType(ct)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	r := multipart.NewReader(buf, params["boundary"])

	for {
		part, err := r.NextPart()
		if err == io.EOF {
			break
		}
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		switch part.FormName() {
		case "url":
			pr := bytes.NewBuffer(nil)
			io.Copy(pr, part)

			assert.Equal(t, swh.URL, pr.String())
		}
	}
}
