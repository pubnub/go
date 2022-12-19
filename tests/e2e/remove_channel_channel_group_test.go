package e2e

import (
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelFromChannelGroupNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"ch"}).
		ChannelGroup(randomized("cg")).
		Execute()

	assert.Nil(err)
}

func TestRemoveChannelFromChannelGroupMissingChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroup().
		ChannelGroup(randomized("cg")).
		Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestRemoveChannelFromChannelGroupMissingChannelContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroupWithContext(backgroundContext).
		ChannelGroup(randomized("cg")).
		Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestRemoveChannelFromChannelGroupMissingGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{randomized("cg")}).
		Execute()

	assert.Contains(err.Error(), "Missing Channel Group")
}

func TestRemoveChannelFromChannelGroupSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters:
	// .,:*
	validCharacters := "?#[]@!$&'()+;=`|"

	config.SetUserId(pubnub.UserId(validCharacters))

	pn := pubnub.NewPubNub(config)

	// Not allowed characters:
	// ?#[]@!$&'()+;=`|
	groupCharacters := "-_~"

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{groupCharacters}).
		ChannelGroup(groupCharacters).
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
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	// await for adding channels
	time.Sleep(2 * time.Second)

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)
	if res != nil {
		assert.Equal(2, len(res.Channels))
		if len(res.Channels) > 1 {
			assert.Equal(myChannel, res.Channels[1])
			assert.Equal(myAnotherChannel, res.Channels[0])
		} else {
			assert.Fail("len(res.Channels) <= 1")
		}
		assert.Equal(myGroup, res.ChannelGroup)
	}
	_, _, err = pn.RemoveChannelFromChannelGroup().
		Channels([]string{myChannel, myAnotherChannel}).
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)

	// await for remove channels
	<-time.After(1 * time.Second)

	res, _, err = pn.ListChannelsInChannelGroup().
		ChannelGroup(myGroup).
		Execute()

	assert.Nil(err)
	if res != nil {
		assert.Equal(0, len(res.Channels))
		assert.Equal(myGroup, res.ChannelGroup)
	}
}
