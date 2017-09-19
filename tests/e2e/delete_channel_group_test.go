package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, err := pn.DeleteChannelGroup().
		ChannelGroup("cg").
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, err := pn.DeleteChannelGroup().
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestRemoveChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := configCopy()

	// Not allowed characters:
	// .,:*

	validCharacters := "-_~?#[]@!$&'()+;=`|"

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, err := pn.DeleteChannelGroup().
		ChannelGroup(validCharacters).
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelGroupSuccessRemoved(t *testing.T) {
	assert := assert.New(t)
	myChannel := "my-channel"
	myGroup := "my-unique-group"

	pn := pubnub.NewPubNub(configCopy())

	_, err := pn.AddChannelChannelGroup().
		Channels([]string{myChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	_, err = pn.RemoveChannelChannelGroup().
		Channels([]string{myChannel}).
		Group(myGroup).
		Execute()

	assert.Nil(err)

	res, err := pn.ListAllChannelsChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	assert.Equal(0, len(res.Channels))
	assert.Equal(myGroup, res.Group)
}
