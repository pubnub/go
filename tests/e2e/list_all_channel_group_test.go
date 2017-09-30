package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestListAllChannelGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup("cg").
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

	config := configCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestListAllChannelGroupSuccess(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel"
	myGroup := randomized("my-group")

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
