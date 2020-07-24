package pubnub

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestHeartbeatRequestBasic(t *testing.T) {
	assert := assert.New(t)

	state := make(map[string]interface{})
	state["one"] = []string{"qwerty"}
	state["two"] = 2

	opts := &heartbeatOpts{
		pubnub:        pubnub,
		State:         state,
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/heartbeat",
			strings.Join(opts.Channels, ",")),
		u.EscapedPath(), []int{})

	u2, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("cg", u2.Get("channel-group"))
	assert.Equal(`{"one":["qwerty"],"two":2}`, u2.Get("state"))
}

func TestNewHeartbeatBuilder(t *testing.T) {
	assert := assert.New(t)

	state := make(map[string]interface{})
	state["one"] = []string{"qwerty"}
	state["two"] = 2

	o := newHeartbeatBuilder(pubnub)
	o.State(state)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/heartbeat",
			strings.Join(o.opts.Channels, ",")),
		u.EscapedPath(), []int{})

	u2, err := o.opts.buildQuery()
	assert.Nil(err)

	assert.Equal("cg", u2.Get("channel-group"))
	assert.Equal(`{"one":["qwerty"],"two":2}`, u2.Get("state"))
}

func TestNewHeartbeatBuilderContext(t *testing.T) {
	assert := assert.New(t)

	state := make(map[string]interface{})
	state["one"] = []string{"qwerty"}
	state["two"] = 2

	o := newHeartbeatBuilderWithContext(pubnub, backgroundContext)
	o.State(state)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/heartbeat",
			strings.Join(o.opts.Channels, ",")),
		u.EscapedPath(), []int{})

	u2, err := o.opts.buildQuery()
	assert.Nil(err)

	assert.Equal("cg", u2.Get("channel-group"))
	assert.Equal(`{"one":["qwerty"],"two":2}`, u2.Get("state"))
}

func TestHeartbeatValidateChAndCg(t *testing.T) {
	assert := assert.New(t)

	opts := &heartbeatOpts{
		pubnub: pubnub,
	}
	err := opts.validate()
	assert.Equal("pubnub/validation: pubnub: Heartbeat: Missing Channel or Channel Group", err.Error())
}

func TestHeartbeatValidateSubKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &heartbeatOpts{
		pubnub: pn,
	}
	err := opts.validate()
	assert.Equal("pubnub/validation: pubnub: Heartbeat: Missing Subscribe Key", err.Error())
}
