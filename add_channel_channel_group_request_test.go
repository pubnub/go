package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestAddChannelRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelOpts{
		Channels:     []string{"ch1", "ch2", "ch3"},
		ChannelGroup: "cg",
		pubnub:       pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("add", "ch1,ch2,ch3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestAddChannelRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelOpts{
		Channels:     []string{"ch1", "ch2", "ch3"},
		ChannelGroup: "cg",
		pubnub:       pubnub,
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	opts.QueryParam = queryParam

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")
	expected.Set("add", "ch1,ch2,ch3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewAddChannelToChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newAddChannelToChannelGroupBuilder(pubnub)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})
	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("add", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewAddChannelToChannelGroupBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newAddChannelToChannelGroupBuilderWithContext(pubnub, backgroundContext)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})
	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("add", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestAddChannelOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &addChannelOpts{
		Channels:     []string{"ch1", "ch2", "ch3"},
		ChannelGroup: "cg",
		pubnub:       pn,
	}
	assert.Equal("pubnub/validation: pubnub: Add Channel To Channel Group: Missing Subscribe Key", opts.validate().Error())
}
