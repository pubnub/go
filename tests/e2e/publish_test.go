package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

const respSuccess = `[1,"Sent","14981595400555832"]`

// NOTICE: not stubbed publish
func TestPublishSuccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	pn.Config.CipherKey = "enigma"

	res, _, err := pn.Publish().
		Channel("ch").Message("hey").UsePost(true).Serialize(true).Execute()

	assert.Nil(err)
	if res != nil {
		assert.True(14981595400555832 < res.Timestamp)
	}
	pn.Config.CipherKey = ""
}

func TestPublishSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Publish().
		Channel("ch").
		Message("hey").
		ShouldStore(false).
		Execute()

	assert.Nil(err)
}

func TestPublishSuccessSlice(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Publish().
		Channel("ch").
		Message([]string{"hey1", "hey2", "hey3"}).
		ShouldStore(false).
		Execute()

	assert.Nil(err)
}

// !go1.8 returns just "request canceled" error for canceled context
// go1.8 returns "context deadline exceeded" error in such case
func TestPublishContextTimeout(t *testing.T) {
	assert := assert.New(t)
	ms := 50
	timeout := time.Duration(ms) * time.Millisecond
	ctx, cancel := contextWithTimeout(backgroundContext, timeout)
	defer cancel()

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.PublishWithContext(ctx).
		Channel("ch").
		Message("hey").
		Execute()

	if err != nil {
		// 1.6 hack
		if strings.Contains(err.Error(), "request canceled") {
			return
		}

		assert.Contains(err.Error(), "context deadline exceeded")
		return
	}
}

func TestPublishContextCancel(t *testing.T) {
	assert := assert.New(t)
	ms := 500
	timeout := time.Duration(ms) * time.Millisecond
	ctx, cancel := contextWithTimeout(backgroundContext, timeout)

	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.PublishWithContext(ctx).
		Channel("ch").
		Message("hey").
		Execute()

	if err != nil {
		// 1.6 hack
		if strings.Contains(err.Error(), "request canceled") {
			return
		}

		assert.Contains(err.Error(), "context canceled")
		return
	}
}

func XTestPublishTimeout(t *testing.T) {
	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Publish().
		Channel("ch").
		Message("hey").
		UsePost(false).
		Execute()

	assert.Contains(t, err.Error(), "Failed to execute request")

}

func TestPublishMissingPublishKey(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.SubscribeKey = "demo"
	cfg.PublishKey = ""

	pn := pubnub.NewPubNub(cfg)

	_, _, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Contains(err.Error(), "Missing Publish Key")
}

func TestPublishMissingMessage(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.PublishKey = "0a5c823c-c1fd-4c3f-b31a-8a0b545fa463"
	cfg.SubscribeKey = "sub-c-d69e3958-1528-11e7-bc52-02ee2ddab7fe"

	pn := pubnub.NewPubNub(cfg)

	_, _, err := pn.Publish().Channel("ch").Execute()

	assert.Contains(err.Error(), "Missing Message")
}

func TestPublishMissingChannel(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.PublishKey = "0a5c823c-c1fd-4c3f-b31a-8a0b545fa463"
	cfg.SubscribeKey = "sub-c-d69e3958-1528-11e7-bc52-02ee2ddab7fe"

	pn := pubnub.NewPubNub(cfg)

	_, _, err := pn.Publish().Message("hey").Execute()

	assert.Contains(err.Error(), "Missing Channel")
}

func TestPublishServerError(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/publish/%s/%s/0/ch/0/%s", config.PublishKey, config.SubscribeKey, "%22hey%22"),
		Query:              "seqn=1",
		ResponseBody:       "",
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 403,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())

	_, _, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Contains(err.Error(), "403")
}

func TestPublishNetworkError(t *testing.T) {
	assert := assert.New(t)

	cfg := pamConfigCopy()
	cfg.Origin = "foo.bar"
	pn := pubnub.NewPubNub(cfg)

	_, _, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Contains(err.Error(), fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"))

	assert.Contains(err.Error(), "no such host")

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"dial tcp: lookup")
}

// WARNING: not mocked request
func TestPublishSigned(t *testing.T) {
	assert := assert.New(t)

	// Not allowed characters: /?#,
	validCharacters := "-._~:[]@!$&'()*+;=`|"

	config := pamConfigCopy()
	config.UUID = validCharacters

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Publish().Channel("ch").
		Message([]string{"hey", "hey2", "hey3"}).Execute()

	assert.Nil(err)
}

func TestPublishSuperCall(t *testing.T) {
	assert := assert.New(t)

	// Not allowed characters: /?#,
	validCharacters := "-._~:[]@!$&'()*+;=`|"

	config := pamConfigCopy()
	config.UUID = validCharacters

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.Publish().Channel(validCharacters).
		Message([]string{validCharacters, validCharacters,
			validCharacters}).Meta(validCharacters).Execute()

	assert.Nil(err)
}
