package botify

import "io"

type APIMethod interface {
	Method() string
	Payload() io.Reader
}
