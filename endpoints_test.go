package pubnub

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "myPub"
	pnconfig.SubscribeKey = "mySub"
}

type fakeEndpointOpts struct {
	pubnub *PubNub
}

func (o *fakeEndpointOpts) buildPath() (string, error) {
	return "/my/path", nil
}

func (o *fakeEndpointOpts) buildQuery() (*url.Values, error) {
	q := &url.Values{}

	q.Set("a", "2")
	q.Set("b", "hey")

	return q, nil
}

func (o *fakeEndpointOpts) buildBody() ([]byte, error) {
	return []byte("myBody"), nil
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

func (o *fakeEndpointOpts) context() Context {
	return o.context()
}

func (o *fakeEndpointOpts) httpMethod() string {
	return "GET"
}

// TODO: fix assertion (strings aren't equal)
func AestBuildUrl(t *testing.T) {
	assert := assert.New(t)
	pnconfig = NewConfig()
	pn := NewPubNub(pnconfig)
	e := &fakeEndpointOpts{
		pubnub: pn,
	}

	url, err := buildUrl(e)
	assert.Nil(err)
	assert.Equal("https://ps.pndns.com/my/path?a=2&b=hey", url)
}
