package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenUuid(t *testing.T) {
	assert := assert.New(t)

	uuid, err := GenUuid()
	assert.Nil(err)
	assert.Len(uuid, 32)
}
