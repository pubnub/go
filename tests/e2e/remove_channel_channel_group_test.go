package e2e

import (
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelFromChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"ch"}).
		Group("cg").
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelFromChannelGroupMissingChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Group("cg").
		Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestRemoveChannelFromChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"ch"}).
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestRemoveChannelFromChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{validCharacters}).
		Group(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelFromChannelGroupSuccess(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel"
	myAnotherChannel := "my-another-channel"
	myGroup := randomized("my-group")

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelToChannelGroup().
		Channels([]string{myChannel, myAnotherChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	// await for adding channels
	time.Sleep(2 * time.Second)

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(2, len(res.Channels))
	assert.Equal(myChannel, res.Channels[1])
	assert.Equal(myAnotherChannel, res.Channels[0])
	assert.Equal(myGroup, res.Group)

	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{myChannel, myAnotherChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	// await for remove channels
	<-time.After(1 * time.Second)

	res, _, err = pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(0, len(res.Channels))
	assert.Equal(myGroup, res.Group)
}
