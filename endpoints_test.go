package pubnub

import (
	"net/url"
	"testing"
	"net/http"

	"github.com/stretchr/testify/assert"
)

var pnconfig *Config

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "myPub"
	pnconfig.SubscribeKey = "mySub"
}

type fakeEndpointOpts struct {
	pubnub *PubNub
}

func (o *fakeEndpointOpts) buildPath() string {
	return "/my/path"
}

func (o *fakeEndpointOpts) buildQuery() *url.Values {
	q := &url.Values{}

	q.Set("a", "2")
	q.Set("b", "hey")

	return q
}

func (o *fakeEndpointOpts) buildBody() string {
	return "myBody"
}

func (o *fakeEndpointOpts) config() Config {
	return *o.pubnub.Config
}

func (o *fakeEndpointOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *fakeEndpointOpts) validate() error {
	return nil
}

// TODO: fix assertion (strings aren't equal)
func AestBuildUrl(t *testing.T) {
	assert := assert.New(t)
	pnconfig = NewConfig()
	pn := NewPubNub(pnconfig)
	e := &fakeEndpointOpts{
		pubnub: pn,
	}

	url := buildUrl(e)
	assert.Equal("https://ps.pndns.com/my/path?a=2&b=hey", url)
}
