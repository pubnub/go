package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	res, _, err := pn.Time().Execute()

	assert.Nil(err)

	assert.True(int64(15059085932399340) < res.Timetoken)
}
