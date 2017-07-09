package e2e

import (
	"log"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestHistorySuccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	res, err := pn.History(&pubnub.HistoryOpts{
		Channel: "ch",
	})
	assert.Nil(err)

	log.Println(res)

	assert.True(14981595400555832 < res.StartTimetoken)
}
