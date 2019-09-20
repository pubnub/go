package e2e

import (
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.DeleteChannelGroup().
		ChannelGroup("cg").
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelGroupNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.DeleteChannelGroupWithContext(backgroundContext).
		ChannelGroup("cg1").
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.DeleteChannelGroup().
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestRemoveChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.UUID = validCharacters

	pn := pubnub.NewPubNub(config)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	_, _, err := pn.DeleteChannelGroup().
		ChannelGroup(validCharacters).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelGroupSuccessRemoved(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel-remove"
	myGroup := "my-unique-group-remove"

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/my-unique-group-remove", config.SubscribeKey),
		Query:              "add=my-channel-remove",
		ResponseBody:       `{"status": 200, "message": "OK", "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/my-unique-group-remove", config.SubscribeKey),
		Query:              "remove=my-channel-remove&q1=v1&q2=v2",
		ResponseBody:       `{"status": 200, "message": "OK", "service": "channel-registry", "error": false}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/channel-registration/sub-key/%s/channel-group/my-unique-group-remove", config.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status": 200, "payload": {"channels": [], "group": "my-unique-group-remove"}, "service": "channel-registry", "error": false}`,
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

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{myChannel}).
		ChannelGroup(myGroup).
		QueryParam(queryParam).
		Execute()

	assert.Nil(err)

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(0, len(res.Channels))
	assert.Equal(myGroup, res.ChannelGroup)
}
