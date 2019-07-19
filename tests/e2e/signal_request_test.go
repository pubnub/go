package e2e

import (
	//"log"
	//"os"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Signal().
		Channel("ch").
		Message("hey").
		Execute()

	assert.Nil(err)

}

func TestSignal(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Signal().
		Channel("ch").
		Message("hey").
		UsePost(true).
		Execute()

	assert.Nil(err)

}
