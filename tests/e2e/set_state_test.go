package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestSetStateSucessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	state := make(map[string]interface{})
	state["age"] = "20"

	_, _, err := pn.SetState().State(state).Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).Execute()

	assert.Nil(err)
}

func TestSetStateSuperCall(t *testing.T) {
	assert := assert.New(t)

	setStateCharacters := "-.,_~:?#[]@!$&'()*+;=`|"

	config := pamConfigCopy()

	// Not allowed characters: /
	config.Uuid = setStateCharacters

	config.AuthKey = SPECIAL_CHANNEL

	pn := pubnub.NewPubNub(config)
	state := make(map[string]interface{})
	state["qwerty"] = SPECIAL_CHARACTERS

	_, _, err := pn.SetState().State(state).Channels([]string{setStateCharacters}).
		ChannelGroups([]string{setStateCharacters}).
		State(state).
		Execute()

	assert.Nil(err)
}
