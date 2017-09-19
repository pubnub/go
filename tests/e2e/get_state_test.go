package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestGetStateNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, err := pn.GetState().
		Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).
		Execute()

	assert.Nil(err)
}

func TestGetStateSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := configCopy()

	pn := pubnub.NewPubNub(config)

	// Not allowed characters: /
	validCharacters := "-.,_~:?#[]@!$&'()*+;=`|"

	config.Uuid = validCharacters
	config.AuthKey = SPECIAL_CHARACTERS

	_, err := pn.GetState().
		Channels([]string{validCharacters}).
		ChannelGroups([]string{validCharacters}).
		Execute()

	assert.Nil(err)
}

func TestGetStateSucess(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	state := make(map[string]interface{})
	state["age"] = "20"
	state["name"] = "John Doe"

	_, err := pn.SetState().
		State(state).
		Channels([]string{"ch"}).
		Execute()

	assert.Nil(err)

	res, err := pn.GetState().
		Channels([]string{"ch"}).
		Execute()

	assert.Nil(err)
	age, _ := res.State["age"].(string)
	name, _ := res.State["name"].(string)

	assert.Equal("20", age)
	assert.Equal("John Doe", name)
}
