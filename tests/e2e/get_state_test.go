package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestGetStateNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.GetState().
		Channels([]string{"ch"}).
		ChannelGroups([]string{"cg"}).
		Execute()

	assert.Nil(err)
}

func TestGetStateSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := configCopy()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/presence/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel/-.%2C_%7E%3A%5B%5D%40%21%24%26%27%28%29%2A%2B%3B%3D%60%7C/uuid/-.%2C_~%3A%5B%5D%40%21%24%26%27%28%29%2A%2B%3B%3D%60%7C",
		Query:              "channel-group=-.%2C_%7E%3A%5B%5D%40%21%24%26%27%28%29%2A%2B%3B%3D%60%7C",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"channels": {}, "total_channels": 0, "total_occupancy": 0}, "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	// Not allowed characters: /?#
	validCharacters := "-.,_~:[]@!$&'()*+;=`|"

	config.Uuid = validCharacters
	config.AuthKey = SPECIAL_CHARACTERS

	_, _, err := pn.GetState().
		Channels([]string{validCharacters}).
		ChannelGroups([]string{validCharacters}).
		Execute()

	assert.Nil(err)
}

func TestGetStateSucess(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/presence/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel/ch/uuid/" + config.Uuid + "/data",
		Query:              "state=%7B%22age%22%3A%2220%22%2C%22name%22%3A%22John%20Doe%22%7D",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"age": "20", "name": "John Doe"}, "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/presence/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel/ch/uuid/" + config.Uuid,
		Query:              "",
		ResponseBody:       `{"status": 200, "message": "OK", "payload": {"age": "20", "name": "John Doe"}, "uuid": "bb45300a-25fb-4b14-8de1-388393274a54", "channel": "ch", "service": "Presence"}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "channel-group"},
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
	age, _ := res.State["age"].(string)
	name, _ := res.State["name"].(string)

	assert.Equal("20", age)
	assert.Equal("John Doe", name)
}
