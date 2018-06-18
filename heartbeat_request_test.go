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
}
