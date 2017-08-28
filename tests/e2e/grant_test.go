package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestGrantSucccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{"pam-key"}).Channels([]string{"ch1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
}

func TestGrantMultipleMixed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{"my-auth-key-1", "my-auth-key-2"}).
		Channels([]string{"ch1", "ch2", "ch3"}).
		Groups([]string{"cg1", "cg2", "cg3"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
}

func TestGrantSingleChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).
		Channels([]string{"ch1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.True(res.Channels["ch1"].WriteEnabled)
	assert.True(res.Channels["ch1"].ReadEnabled)
	assert.False(res.Channels["ch1"].ManageEnabled)
}

func TestGrantSingleChannelWithAuth(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).Manage(false).
		AuthKeys([]string{"my-pam-key"}).
		Channels([]string{"ch1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.Channels["ch1"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).
		Channels([]string{"ch1", "ch2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.Channels["ch1"].WriteEnabled)
	assert.True(res.Channels["ch1"].ReadEnabled)
	assert.False(res.Channels["ch1"].ManageEnabled)

	assert.True(res.Channels["ch2"].WriteEnabled)
	assert.True(res.Channels["ch2"].ReadEnabled)
	assert.False(res.Channels["ch2"].ManageEnabled)
}

func TestGrantMultipleChannelsWithAuth(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).
		AuthKeys([]string{"my-pam-key"}).
		Channels([]string{"ch1", "ch2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.Channels["ch1"].AuthKeys["my-pam-key"].ManageEnabled)

	assert.True(res.Channels["ch2"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.Channels["ch2"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.Channels["ch2"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantSingleGroup(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).
		Groups([]string{"cg1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].ManageEnabled)
}

func TestGrantSingleGroupWithAuth(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Groups([]string{"cg1"}).
		AuthKeys([]string{"my-pam-key"}).
		Write(true).
		Read(true).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantMultipleGroups(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).
		Groups([]string{"cg1", "cg2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].ManageEnabled)

	assert.True(res.ChannelGroups["cg2"].WriteEnabled)
	assert.True(res.ChannelGroups["cg2"].ReadEnabled)
	assert.False(res.ChannelGroups["cg2"].ManageEnabled)
}

func TestGrantMultipleGroupsWithAuth(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	res, err := pn.Grant().
		Read(true).Write(true).
		AuthKeys([]string{"my-pam-key"}).
		Groups([]string{"cg1", "cg2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ManageEnabled)

	assert.True(res.ChannelGroups["cg2"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.ChannelGroups["cg2"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.ChannelGroups["cg2"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(config)

	_, err := pn.Grant().
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{SPECIAL_CHARACTERS}).
		Channels([]string{SPECIAL_CHANNEL}).
		Groups([]string{SPECIAL_CHANNEL}).
		Execute()

	assert.Nil(err)
}
