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

	// Not allowed characters:
	// .,:*
	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config := pamConfigCopy()

	// Not allowed characters: /
	config.Uuid = validCharacters

	config.AuthKey = validCharacters

	pn := pubnub.NewPubNub(config)
	state := make(map[string]interface{})
	state["qwerty"] = validCharacters
	state["a"] = "b"

	// Not allowed characters:
	// ?#[]@!$&'()+;=`|
	groupCharacters := "-_~"

	_, _, err := pn.SetState().
		State(state).
		Channels([]string{validCharacters}).
		ChannelGroups([]string{groupCharacters}).
		Execute()

	assert.Nil(err)
}
