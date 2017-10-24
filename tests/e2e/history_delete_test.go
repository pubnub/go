package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestHistoryDeleteNotStubbed(t *testing.T) {
	assert := assert.New(t)

	ch := randomized("h-ch")
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.DeleteMessages().
		Channel(ch).
		Execute()

	assert.Nil(err)
}
