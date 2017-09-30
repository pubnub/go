package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg").
		Execute()

	assert.Nil(err)
}

func TestAddChannelChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{"ch"}).
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestAddChannelChannelGroupMissingChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelChannelGroup().
		Group("cg").
		Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestAddChannelChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := configCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{validCharacters}).
		Group(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestAddChannelChannelGroupSuccessAdded(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel-group/my-unique-group",
		Query:              "add=my-channel",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"channel-registry\", \"error\": \"false\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel-group/my-unique-group",
		Query:              "",
		ResponseBody:       "{\"status\": \"200\", \"payload\": {\"channels\": [\"my-channel\"], \"group\": \"my-unique-group\"}, \"service\": \"channel-registry\", \"error\": \"false\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	myChannel := "my-channel"
	myGroup := "my-unique-group"

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{myChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	res, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(myChannel, res.Channels[0])
	assert.Equal(myGroup, res.Group)
}
