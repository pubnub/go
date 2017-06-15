package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pnconfig *Config

func init() {
	pnconfig = NewConfig()

	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"
	pnconfig.Secure = false
	pnconfig.ConnectionTimeout = 2
	pnconfig.NonSubscribeRequestTimeout = 2
}

func TestBuildUrl(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(pnconfig)

	params := &PublishOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pn,
	}

	url := buildUrl(params)

	assert.Equal(url,
		"http://ps.pndsn.com/publish/my_pub_key/my_sub_key/0/ch/0/hey?pnsdk=4&seqn=1&store=0&uuid=TODO-setup-uuid")
}
