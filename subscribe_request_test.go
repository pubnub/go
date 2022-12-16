package pubnub

import (
	"net/url"
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
