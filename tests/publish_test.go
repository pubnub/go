package pntests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/pnerr"
	"github.com/stretchr/testify/assert"
)

var pnconfig *pubnub.Config

func init() {
	pnconfig = pubnub.NewConfig()
	pnconfig.PublishKey = "my_pub_key"
	pnconfig.SubscribeKey = "my_sub_key"
	pnconfig.SecretKey = "my_secret_key"
	pnconfig.Origin = "localhost:3000"
	pnconfig.Secure = false
	pnconfig.ConnectionTimeout = 2
	pnconfig.NonSubscribeRequestTimeout = 2
}

func TestPublishContextTimeoutSync(t *testing.T) {
	assert := assert.New(t)
	ms := 500
	timeout := time.Duration(ms) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	shutdown := make(chan bool)
	go servePublish(pnconfig.NonSubscribeRequestTimeout+1, shutdown)

	pn := pubnub.NewPubNub(pnconfig)

	res, err := pn.PublishWithContext(ctx, &pubnub.PublishOpts{
		Channel: "ch",
		Message: "hey",
	})

	if err == nil {
		assert.Fail("Received success instead of context deadline: %v", res)
		return
	}

	assert.Equal(fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"), err.Error())

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"context deadline exceeded")

	shutdown <- true
}

func TestPublishContextCancel(t *testing.T) {
	assert := assert.New(t)
	ms := 500
	timeout := time.Duration(ms) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	shutdown := make(chan bool)
	go servePublish(pnconfig.NonSubscribeRequestTimeout+1, shutdown)
	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel()
	}()

	pn := pubnub.NewPubNub(pnconfig)

	res, err := pn.PublishWithContext(ctx, &pubnub.PublishOpts{
		Channel: "ch",
		Message: "hey",
	})

	if err == nil {
		assert.Fail("Received success instead of context deadline: %v", res)
		return
	}

	assert.Equal(fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"), err.Error())

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"context canceled")

	shutdown <- true
}

func TestRequestTimeoutSync(t *testing.T) {
	assert := assert.New(t)
	shutdown := make(chan bool)

	go servePublish(2, shutdown)

	pn := pubnub.NewPubNub(pnconfig)

	params := &pubnub.PublishOpts{
		Channel: "ch1",
		Message: "hey",
		UsePost: false,
	}

	_, err := pn.Publish(params)

	assert.Equal(fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"), err.Error())
	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"timeout awaiting response headers")

	shutdown <- true
}

func TestPublishMissingPublishKey(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.SubscribeKey = "demo"
	cfg.PublishKey = ""

	pn := pubnub.NewPubNub(cfg)

	params := &pubnub.PublishOpts{
		Channel: "ch",
		Message: "hey",
	}

	_, err := pn.Publish(params)

	assert.Contains(err.Error(), "pubnub: Missing Publish Key")
}

func TestPublishMissingMessage(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.PublishKey = "0a5c823c-c1fd-4c3f-b31a-8a0b545fa463"
	cfg.SubscribeKey = "sub-c-d69e3958-1528-11e7-bc52-02ee2ddab7fe"

	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish(&pubnub.PublishOpts{
		Channel: "hey",
	})

	assert.Contains(err.Error(), "pubnub: Missing Message")
}

func TestPublishMissingChannel(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.PublishKey = "0a5c823c-c1fd-4c3f-b31a-8a0b545fa463"
	cfg.SubscribeKey = "sub-c-d69e3958-1528-11e7-bc52-02ee2ddab7fe"

	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish(&pubnub.PublishOpts{
		Message: "hey",
	})

	assert.Contains(err.Error(), "pubnub: Missing Channel")
}

func TestPublishServerError(t *testing.T) {
	assert := assert.New(t)

	cfg := pamConfigCopy()
	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish(&pubnub.PublishOpts{
		Channel: "ch",
		Message: "hey",
	})

	assert.Equal(fmt.Sprintf(serverErrorTemplate, 403), err.Error())
}

func TestPublishNetworkError(t *testing.T) {
	assert := assert.New(t)

	cfg := pamConfigCopy()
	cfg.Origin = "foo.bar"
	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish(&pubnub.PublishOpts{
		Channel: "ch",
		Message: "hey",
	})

	assert.Equal(fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"), err.Error())

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"dial tcp: lookup")
}

func TestNewRequestErrorHost(t *testing.T) {
	assert := assert.New(t)

	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://aaaaaa.com/", nil)

	_, err := client.Do(r)

	assert.Contains(err.Error(), "no such host")
}
