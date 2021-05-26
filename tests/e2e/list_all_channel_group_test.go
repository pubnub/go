package e2e

import (
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v5"
	"github.com/pubnub/go/v5/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestListAllChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	_, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup("cg").
		Execute()

	assert.Nil(err)
}

func TestListAllChannelGroupNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	_, _, err := pn.ListChannelsInChannelGroupWithContext(backgroundContext).
		ChannelGroup("cg1").
		Execute()

	assert.Nil(err)
}

func TestListAllChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.DeleteChannelGroup().
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestListAllChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.UUID = validCharacters

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestListAllChannelGroupSuccess(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel"
	myGroup := randomized("my-group")

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/%s", config.SubscribeKey, myGroup),
		Query:              "add=my-channel",
		ResponseBody:       `{"status": 200, "message": "OK", "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/%s", config.SubscribeKey, myGroup),
		Query:              "",
		ResponseBody:       `{"status": 200, "payload": {"channels": ["my-channel"], "group": "` + myGroup + `"}, "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/%s", config.SubscribeKey, myGroup),
		Query:              "remove=my-channel",
		ResponseBody:       `{"status": 200, "message": "OK", "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{myChannel}).
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	// await for adding channel
	time.Sleep(2 * time.Second)

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(myChannel, res.Channels[0])
	assert.Equal(myGroup, res.ChannelGroup)

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{myChannel}).
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)
}
