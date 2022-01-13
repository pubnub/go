package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringAsString(t *testing.T) {
	assert := assert.New(t)

	str, err := ValueAsString("blah")

	assert.Nil(err)
	assert.Equal([]byte("\"blah\""), str)
}

func TestUUID(t *testing.T) {
	assert := assert.New(t)

	assert.Len(UUID(), 39)
}

func TestURLEncode(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("%5B%22hey1%22%2C%20%22hey2%22%2C%20%22hey3%5D",
		URLEncode(`["hey1", "hey2", "hey3]`))
}
