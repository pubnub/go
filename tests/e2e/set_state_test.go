package e2e

import (
	"fmt"
	"reflect"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestNewSetStateResponse(t *testing.T) {
	assert := assert.New(t)

	pubnub.Config.Uuid = "my-custom-uuid"

	jsonBytes := []byte(`{"status": 200, "message": "OK", "payload": {"k": "v"}, "uuid": "my-custom-uuid", "service": "Presence"}`)

	res, _, err := newSetStateResponse(jsonBytes, fakeResponseState)
	assert.Nil(err)
	if s, ok := res.State.(map[string]interface{}); ok {
		assert.Equal("v", s["k"])
	} else {
		assert.Fail("!map[string]interface{}")
	}
}

func TestSetStateSucessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	state := make(map[string]interface{})
	state["age"] = "20"

	setStateRes, _, err := pn.SetState().State(state).Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).Execute()

	assert.Nil(err)
	if s, ok := setStateRes.State.(map[string]interface{}); ok {
		assert.Equal("20", s["age"])
	} else {
		assert.Fail(fmt.Sprintf("!map[string]interface{} %v %v", reflect.TypeOf(setStateRes.State).Kind(), reflect.TypeOf(setStateRes.State)))
	}

	assert.Equal("OK", setStateRes.Message)

	getStateRes, _, err := pn.GetState().
		Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).
		Execute()

	assert.Nil(err)
	if s, ok := getStateRes.State["ch"].(map[string]interface{}); ok {
		assert.Equal("20", s["age"])
	} else {
		assert.Fail(fmt.Sprintf("!map[string]interface{} %v %v", reflect.TypeOf(getStateRes.State["ch"]).Kind(), reflect.TypeOf(setStateRes.State)))
	}
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
