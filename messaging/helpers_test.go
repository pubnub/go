package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPnpresToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", addPnpresToString(""))
	assert.Equal("qwer-pnpres", addPnpresToString("qwer"))
	assert.Equal("qwer-pnpres,asdf-pnpres,zxcv-pnpres",
		addPnpresToString("qwer,asdf,zxcv"))
}
