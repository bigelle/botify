package botify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type APIMethod interface {
	ContentType() string
	Method() string
	Payload() (io.Reader, error)
}

type GetUpdates struct {
	Offset         int       `json:"offset"`
	Limit          int       `json:"limit"`
	Timeout        int       `json:"timeout"`
	AllowedUpdates *[]string `json:"allowed_updates"`
}

func (m *GetUpdates) ContentType() string {
	return "application/json"
}

func (m *GetUpdates) Method() string {
	return "getUpdates"
}

func (m *GetUpdates) Payload() (io.Reader, error) {
	buf := &bytes.Buffer{}
	
	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		return nil, fmt.Errorf("encoding getUpdates payload: %w", err)
	}

	return buf, nil
}
