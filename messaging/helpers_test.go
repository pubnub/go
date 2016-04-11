package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPnpresToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", addPnpresToString(""))
	assert.Equal("qwer-pnpres", addPnpresToString("qwer"))
	assert.Equal("qwer-pnpres,asdf-pnpres,zxcv-pnpres",
		addPnpresToString("qwer,asdf,zxcv"))
}

func TestSplitItems(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]string{}, splitItems(""))
	assert.Equal([]string{"ch1"}, splitItems("ch1"))
	assert.Equal([]string{"ch1", "ch2"}, splitItems("ch1,ch2"))
}
