package e2e

import (
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelChannelGroup().
		Channels([]string{"ch"}).
		Group("cg").
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelChannelGroupMissingChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelChannelGroup().
		Group("cg").
		Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestRemoveChannelChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelChannelGroup().
		Channels([]string{"ch"}).
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestRemoveChannelChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := configCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.RemoveChannelChannelGroup().
		Channels([]string{validCharacters}).
		Group(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelChannelGroupSuccess(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel"
	myAnotherChannel := "my-another-channel"
	myGroup := "my-unique-group"

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.AddChannelChannelGroup().
		Channels([]string{myChannel, myAnotherChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	res, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(2, len(res.Channels))
	assert.Equal(myChannel, res.Channels[1])
	assert.Equal(myAnotherChannel, res.Channels[0])
	assert.Equal(myGroup, res.Group)

	_, _, err = pn.RemoveChannelChannelGroup().
		Channels([]string{myChannel, myAnotherChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	// await for remove channels
	<-time.After(1 * time.Second)

	res, _, err = pn.ListAllChannelsChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(0, len(res.Channels))
	assert.Equal(myGroup, res.Group)
}
