package e2e

import (
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go/v5"
	"github.com/pubnub/go/v5/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestGetStateNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.GetState().
		Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).
		UUID("my-custom-uuid").
		Execute()

	assert.Nil(err)
}

func TestGetStateNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.GetStateWithContext(backgroundContext).
		Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).
		UUID("my-custom-uuid").
		Execute()

	assert.Nil(err)
}

func TestGetStateSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := configCopy()

	pn := pubnub.NewPubNub(config)

	// Not allowed characters: /?#,
	validCharacters := "-._~:[]@!$&'()*+;=`|"

	config.UUID = validCharacters
	config.AuthKey = SPECIAL_CHARACTERS

	_, _, err := pn.GetState().
		Channels([]string{validCharacters, validCharacters, validCharacters}).
		ChannelGroups([]string{validCharacters, validCharacters, validCharacters}).
		UUID(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestGetStateSucess(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/ch/uuid/", config.SubscribeKey) + config.UUID + "/data",
		Query:              "state=%7B%22age%22%3A%2220%22%2C%22name%22%3A%22John%20Doe%22%7D",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"age": "20", "name": "John Doe"}, "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_pres"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/ch/uuid/", config.SubscribeKey) + config.UUID,
		Query:              "",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"age": "20", "name": "John Doe"}, "uuid": "bb45300a-25fb-4b14-8de1-388393274a54", "channel": "ch", "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "channel-group", "l_pres"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	state := make(map[string]interface{})
	state["age"] = "20"
	state["name"] = "John Doe"

	_, _, err := pn.SetState().
		State(state).
		Channels([]string{"ch"}).
		Execute()

	assert.Nil(err)

	res, _, err := pn.GetState().
		Channels([]string{"ch"}).
		Execute()

	assert.Nil(err)
	if s, ok := res.State["ch"].(map[string]interface{}); ok {
		assert.Equal("20", s["age"])
		assert.Equal("John Doe", s["name"])

	} else {
		assert.Fail(fmt.Sprintf("!map[string]interface{} "))
	}
}

func TestGetStateSucessQueryParam(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/ch/uuid/", config.SubscribeKey) + config.UUID + "/data",
		Query:              "state=%7B%22age%22%3A%2220%22%2C%22name%22%3A%22John%20Doe%22%7D&q1=v1&q2=v2",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"age": "20", "name": "John Doe"}, "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_pres"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/channel/ch/uuid/", config.SubscribeKey) + config.UUID,
		Query:              "q1=v1&q2=v2",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"age": "20", "name": "John Doe"}, "uuid": "bb45300a-25fb-4b14-8de1-388393274a54", "channel": "ch", "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "channel-group", "l_pres"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	state := make(map[string]interface{})
	state["age"] = "20"
	state["name"] = "John Doe"

	_, _, err := pn.SetState().
		State(state).
		Channels([]string{"ch"}).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err)

	res, _, err := pn.GetState().
		Channels([]string{"ch"}).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err)
	if s, ok := res.State["ch"].(map[string]interface{}); ok {
		assert.Equal("20", s["age"])
		assert.Equal("John Doe", s["name"])

	} else {
		assert.Fail(fmt.Sprintf("!map[string]interface{} "))
	}
}
