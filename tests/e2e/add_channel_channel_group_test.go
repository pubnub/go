package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
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
	myChannel := "my-channel"
	myGroup := "my-unique-group"

	pn := pubnub.NewPubNub(configCopy())

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

	_, _, err = pn.RemoveChannelChannelGroup().
		Channels([]string{myChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)
}
