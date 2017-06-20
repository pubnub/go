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

func TestUuid(t *testing.T) {
	assert := assert.New(t)

	assert.Len(Uuid(), 36)
}
