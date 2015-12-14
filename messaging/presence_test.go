package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnlyPresence(t *testing.T) {
	assert := assert.New(t)

	assert.True(hasNonPresenceChannels("qwer"))
	assert.True(hasNonPresenceChannels("qwer,asdf"))
	assert.True(hasNonPresenceChannels("qwer,asdf-pnpres"))
	assert.False(hasNonPresenceChannels("qwer-pnpres"))
	assert.False(hasNonPresenceChannels("qwer-pnpres,asdf-pnpres"))
}
