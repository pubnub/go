package pubnub

import (
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeSingleChannel(t *testing.T) {
	assert := assert.New(t)
	opts := newSubscribeOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/ch/0", u.EscapedPath(), []int{})
}

func TestSubscribeMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	opts := newSubscribeOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch-1", "ch-2", "ch-3"}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/ch-1,ch-2,ch-3/0", u.EscapedPath(), []int{})
}

func TestSubscribeChannelGroups(t *testing.T) {
	assert := assert.New(t)
	opts := newSubscribeOpts(pubnub, pubnub.ctx)
	opts.ChannelGroups = []string{"cg-1", "cg-2", "cg-3"}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/,/0", u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg-1,cg-2,cg-3")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid"}, []string{})
}

func TestSubscribeMixedParams(t *testing.T) {
	assert := assert.New(t)

	opts := newSubscribeOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.Region = "us-east-1"
	opts.Timetoken = 123
	opts.FilterExpression = "abc"

	path, err := opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/ch/0", u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("tr", "us-east-1")
	expected.Set("filter-expr", "abc")
	expected.Set("tt", "123")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid"}, []string{})
}

func TestSubscribeMixedQueryParams(t *testing.T) {
	assert := assert.New(t)

	opts := newSubscribeOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.Region = "us-east-1"
	opts.Timetoken = 123
	opts.FilterExpression = "abc"
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("tr", "us-east-1")
	expected.Set("filter-expr", "abc")
	expected.Set("tt", "123")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid"}, []string{})
}

func TestSubscribeValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newSubscribeOpts(pn, pn.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.Region = "us-east-1"
	opts.Timetoken = 123
	opts.FilterExpression = "abc"

	assert.Equal("pubnub/validation: pubnub: Subscribe: Missing Subscribe Key", opts.validate().Error())
}

func TestSubscribeValidatePublishKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := newSubscribeOpts(pubnub, pubnub.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.Region = "us-east-1"
	opts.Timetoken = 123
	opts.FilterExpression = "abc"

	assert.Nil(opts.validate())
}

func TestSubscribeValidateCHAndCG(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSubscribeOpts(pn, pn.ctx)
	opts.Region = "us-east-1"
	opts.Timetoken = 123
	opts.FilterExpression = "abc"

	assert.Equal("pubnub/validation: pubnub: Subscribe: Missing Channel", opts.validate().Error())
}

func TestSubscribeValidateState(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newSubscribeOpts(pn, pn.ctx)
	opts.Channels = []string{"ch"}
	opts.ChannelGroups = []string{"cg"}
	opts.Region = "us-east-1"
	opts.Timetoken = 123
	opts.FilterExpression = "abc"
	opts.State = map[string]interface{}{"a": "a"}

	assert.Nil(opts.validate())
}

func TestSubscribeDuplicateChannelsInSameCall(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Subscribe with duplicate channels in the same call
	subscribeOp := &SubscribeOperation{
		Channels: []string{"test-channel", "test-channel", "test-channel"},
	}

	pn.subscriptionManager.adaptSubscribe(subscribeOp)

	// Get the channels from state manager
	channels := pn.subscriptionManager.stateManager.prepareChannelList(false)

	// Verify that the channel appears only once
	assert.Equal(1, len(channels), "Channel should appear only once")
	assert.Equal("test-channel", channels[0], "Channel name should be 'test-channel'")
}

func TestSubscribeDuplicateChannelsMultipleCalls(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Subscribe to the same channel twice in separate operations
	subscribeOp1 := &SubscribeOperation{
		Channels: []string{"test-channel"},
	}
	pn.subscriptionManager.adaptSubscribe(subscribeOp1)

	subscribeOp2 := &SubscribeOperation{
		Channels: []string{"test-channel"},
	}
	pn.subscriptionManager.adaptSubscribe(subscribeOp2)

	// Get the channels from state manager
	channels := pn.subscriptionManager.stateManager.prepareChannelList(false)

	// Verify that the channel appears only once
	assert.Equal(1, len(channels), "Channel should appear only once after multiple subscribe calls")
	assert.Equal("test-channel", channels[0], "Channel name should be 'test-channel'")
}

func TestSubscribeDuplicateChannelsInRequestPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Subscribe with duplicate channels
	subscribeOp := &SubscribeOperation{
		Channels: []string{"ch-1", "ch-2", "ch-1", "ch-3", "ch-2"},
	}

	pn.subscriptionManager.adaptSubscribe(subscribeOp)

	// Get the channels from state manager
	channels := pn.subscriptionManager.stateManager.prepareChannelList(false)

	// Verify that we only have 3 unique channels
	assert.Equal(3, len(channels), "Should have 3 unique channels")

	// Create opts and build path to verify the actual request path
	opts := newSubscribeOpts(pn, pn.ctx)
	opts.Channels = channels

	path, err := opts.buildPath()
	assert.Nil(err)

	// The path should contain each channel only once
	u := &url.URL{
		Path: path,
	}

	// Count occurrences of each channel in the path
	pathStr := u.EscapedPath()
	assert.Contains(pathStr, "ch-1")
	assert.Contains(pathStr, "ch-2")
	assert.Contains(pathStr, "ch-3")

	// Verify each channel appears exactly once in the path
	// The path format is: /v2/subscribe/sub_key/ch-1,ch-2,ch-3/0
	// We check that there are no duplicate channel names in the comma-separated list
	channelsPart := ""
	if strings.Contains(pathStr, "/v2/subscribe/") {
		parts := strings.Split(pathStr, "/")
		if len(parts) >= 5 {
			channelsPart = parts[4] // The channel list is the 5th part (index 4)
		}
	}

	// Split by comma and verify uniqueness
	channelList := strings.Split(channelsPart, ",")
	uniqueChannels := make(map[string]bool)
	for _, ch := range channelList {
		if ch != "" {
			assert.False(uniqueChannels[ch], "Channel %s should not appear more than once in the path", ch)
			uniqueChannels[ch] = true
		}
	}
	assert.Equal(3, len(uniqueChannels), "Should have exactly 3 unique channels in the request path")
}
