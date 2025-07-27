package botify_test

import (
	"testing"

	"github.com/bigelle/botify"
	"github.com/stretchr/testify/assert"
)

//it might be dumb but i have to test it

func TestGetUpdates_ContentType(t *testing.T) {
	var gu botify.GetUpdates

	assert.Equal(t, gu.ContentType(), "application/json")
}
