package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	pubnub = NewPubNub(pnconfig)
}

func TestLeaveRequestSingleChannel(t *testing.T) {
	assert := assert.New(t)

	opts := &LeaveOpts{
		Channels: []string{"ch"},
		pubnub:   pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/leave", opts.Channels[0]),
		u.EscapedPath(), []int{})
}

func TestLeaveRequestMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	opts := &LeaveOpts{
		Channels: []string{"ch1", "ch2", "ch3"},
		pubnub:   pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/leave",
		u.EscapedPath(), []int{})
}

func TestLeaveRequestSingleChannelGroup(t *testing.T) {
	assert := assert.New(t)

	opts := &LeaveOpts{
		ChannelGroups: []string{"cg"},
		pubnub:        pubnub,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestLeaveRequestMultipleChannelGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &LeaveOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestLeaveRequestChannelsAndGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &LeaveOpts{
		Channels:      []string{"ch1", "ch2", "ch3"},
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/leave",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}
