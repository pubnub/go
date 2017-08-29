package pubnub

import (
	"fmt"
	"net/url"
	"strings"
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
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/hearbeat",
			strings.Join(opts.Channels, ",")),
		u.EscapedPath(), []int{})
}
