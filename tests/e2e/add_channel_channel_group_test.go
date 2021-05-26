package e2e

import (
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go/v5"
	"github.com/pubnub/go/v5/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelToChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{"ch"}).
		ChannelGroup("cg").
		Execute()
	assert.Nil(err)
}

func TestAddChannelToChannelGroupNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelToChannelGroupWithContext(backgroundContext).
		Channels([]string{"ch"}).
		ChannelGroup("cg").
		Execute()
	assert.Nil(err)
}

func TestAddChannelToChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{"ch"}).
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestAddChannelToChannelGroupMissingChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelToChannelGroup().
		ChannelGroup("cg").
		Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestAddChannelToChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters:
	// .,:*#`[]&
	validCharacters := "`[]&|='?;-_~@!$()+"

	channelCharacters := "-_~"

	config.UUID = validCharacters

	pn := pubnub.NewPubNub(config)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{channelCharacters}).
		ChannelGroup(validCharacters).
		QueryParam(queryParam).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToChannelGroupSuccessAdded(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/my-unique-group", config.SubscribeKey),
		Query:              "add=my-channel&q1=v1&q2=v2",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"channel-registry\", \"error\": \"false\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/my-unique-group", config.SubscribeKey),
		Query:              "q1=v1&q2=v2",
		ResponseBody:       "{\"status\": \"200\", \"payload\": {\"channels\": [\"my-channel\"], \"group\": \"my-unique-group\"}, \"service\": \"channel-registry\", \"error\": \"false\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})

	myChannel := "my-channel"
	myGroup := "my-unique-group"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	pn.SetClient(interceptor.GetClient())

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{myChannel}).
		ChannelGroup(myGroup).QueryParam(queryParam).
		Execute()

	assert.Nil(err)

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).QueryParam(queryParam).
		Execute()

	assert.Nil(err)

	assert.Equal(myChannel, res.Channels[0])
	assert.Equal(myGroup, res.ChannelGroup)
}
