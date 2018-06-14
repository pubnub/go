package e2e

import (
	//"fmt"
	//"log"
	//"os"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelToChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{"ch"}).
		ChannelGroup("cg").
		Execute()
	//fmt.Println(err.Error())
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
	//config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Not allowed characters:
	// .,:*#`[]&
	validCharacters := "`[]&|='?;-_~@!$()+"

	channelCharacters := "-_~"

	config.UUID = validCharacters
	//config.AuthKey = validCharacters

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{channelCharacters}).
		ChannelGroup(validCharacters).
		Execute()
	//fmt.Println(err.Error())
	assert.Nil(err)
}

func TestAddChannelToChannelGroupSuccessAdded(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64/channel-group/my-unique-group",
		Query:              "add=my-channel",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"service\": \"channel-registry\", \"error\": \"false\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v1/channel-registration/sub-key/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64/channel-group/my-unique-group",
		Query:              "",
		ResponseBody:       "{\"status\": \"200\", \"payload\": {\"channels\": [\"my-channel\"], \"group\": \"my-unique-group\"}, \"service\": \"channel-registry\", \"error\": \"false\"}",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "l_cg"},
		ResponseStatusCode: 200,
	})

	myChannel := "my-channel"
	myGroup := "my-unique-group"

	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	pn.SetClient(interceptor.GetClient())

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{myChannel}).
		ChannelGroup(myGroup).
		Execute()
	//fmt.Println(err.Error())
	assert.Nil(err)

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(myChannel, res.Channels[0])
	assert.Equal(myGroup, res.ChannelGroup)
}
