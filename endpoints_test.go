package pubnub

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pnconfig *Config

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "myPub"
	pnconfig.SubscribeKey = "mySub"
}

type fakeEndpoint struct {
	pubnub *PubNub
}

func (e *fakeEndpoint) buildPath() string {
	return "/my/path"
}

func (e *fakeEndpoint) buildQuery() *url.Values {
	q := &url.Values{}

	q.Set("a", "2")
	q.Set("b", "hey")

	return q
}

func (e *fakeEndpoint) buildBody() string {
	return "myBody"
}

func (e *fakeEndpoint) PubNub() *PubNub {
	return e.pubnub
}

// TODO: fix assertion (strings aren't equal)
func AestBuildUrl(t *testing.T) {
	assert := assert.New(t)
	pnconfig = NewConfig()
	pn := NewPubNub(pnconfig)
	e := &fakeEndpoint{
		pubnub: pn,
	}

	url := buildUrl(e)
	assert.Equal("https://ps.pndns.com/my/path?a=2&b=hey", url)
}
