package e2e

import (
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

var pnconfig *pubnub.Config

const RESP_SUCCESS = `[1,"Sent","14981595400555832"]`

func init() {
	pnconfig = pubnub.NewConfig()
	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"
	pnconfig.ConnectTimeout = 2
	pnconfig.NonSubscribeRequestTimeout = 2
}

// NOTICE: not stubbed publish
func TestPublishSuccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	pn.Config.CipherKey = "enigma"

	res, err := pn.Publish().
		Channel("ch").Message("hey").UsePost(true).Serialize(true).Execute()

	assert.Nil(err)
	assert.True(14981595400555832 < res.Timestamp)
}

func TestPublishSuccess(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/publish/pub_key/sub_key/0/ch/0/%22hey%22",
		Query:              "seqn=1&store=0",
		ResponseBody:       RESP_SUCCESS,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pnconfig)
	pn.SetClient(interceptor.GetClient())

	_, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Nil(err)
}

func TestPublishSuccessSlice(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/publish/pub_key/sub_key/0/ch/0/%5B%22hey1%22,%22hey2%22,%22hey3%22%5D",
		Query:              "seqn=1&store=0",
		ResponseBody:       RESP_SUCCESS,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pnconfig)
	pn.SetClient(interceptor.GetClient())

	_, err := pn.Publish().Channel("ch").Message([]string{"hey", "hey2", "hey3"}).
		Execute()

	assert.Nil(err)
}

// !go1.8 returns just "request canceled" error for canceled context
// go1.8 returns "context deadline exceeded" error in such case
func TestPublishContextTimeout(t *testing.T) {
	assert := assert.New(t)
	ms := 500
	timeout := time.Duration(ms) * time.Millisecond
	ctx, cancel := contextWithTimeout(backgroundContext, timeout)
	defer cancel()

	pn := pubnub.NewPubNub(pnconfig)
	pn.SetClient(stubs.NewSleeperClient(ms + 3000))

	res, err := pn.PublishWithContext(ctx).Channel("ch").Message("hey").Execute()

	if err == nil {
		assert.Fail("Received success instead of context deadline: %v", res)
		return
	}

	assert.Contains(err.Error(), fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"))

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		ERR_CONTEXT_DEADLINE)
}

// TODO: replace with transport listener
func TestPublishContextCancel(t *testing.T) {
	assert := assert.New(t)
	ms := 500
	timeout := time.Duration(ms) * time.Millisecond
	ctx, cancel := contextWithTimeout(backgroundContext, timeout)

	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel()
	}()

	pn := pubnub.NewPubNub(pnconfig)
	pn.SetClient(stubs.NewSleeperClient(ms + 3000))

	res, err := pn.PublishWithContext(ctx).Channel("ch").Message("hey").Execute()

	if err == nil {
		assert.Fail("Received success instead of context deadline: %v", res)
		return
	}

	assert.Contains(err.Error(), fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"))

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		ERR_CONTEXT_CANCELLED)
}

// TODO: fix this test after timeouts refactoring
func ATestPublishTimeout(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pnconfig)
	pn.SetClient(stubs.NewSleeperClient(
		pnconfig.NonSubscribeRequestTimeout*1000 + 1000))

	_, err := pn.Publish().Channel("ch").Message("hey").
		UsePost(false).Execute()

	assert.Equal(fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"), err.Error())

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"timeout awaiting response headers")
}

func TestPublishMissingPublishKey(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.SubscribeKey = "demo"
	cfg.PublishKey = ""

	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Contains(err.Error(), "pubnub: Missing Publish Key")
}

func TestPublishMissingMessage(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.PublishKey = "0a5c823c-c1fd-4c3f-b31a-8a0b545fa463"
	cfg.SubscribeKey = "sub-c-d69e3958-1528-11e7-bc52-02ee2ddab7fe"

	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish().Channel("ch").Execute()

	assert.Contains(err.Error(), "pubnub: Missing Message")
}

func TestPublishMissingChannel(t *testing.T) {
	assert := assert.New(t)

	cfg := pubnub.NewConfig()
	cfg.PublishKey = "0a5c823c-c1fd-4c3f-b31a-8a0b545fa463"
	cfg.SubscribeKey = "sub-c-d69e3958-1528-11e7-bc52-02ee2ddab7fe"

	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish().Message("hey").Execute()

	assert.Contains(err.Error(), "pubnub: Missing Channel")
}

// Grant permissions added from another sdk
func xTestPublishServerError(t *testing.T) {
	assert := assert.New(t)

	cfg := pamConfigCopy()
	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Contains(err.Error(), fmt.Sprintf(serverErrorTemplate, 403))
}

func TestPublishNetworkError(t *testing.T) {
	assert := assert.New(t)

	cfg := pamConfigCopy()
	cfg.Origin = "foo.bar"
	pn := pubnub.NewPubNub(cfg)

	_, err := pn.Publish().Channel("ch").Message("hey").Execute()

	assert.Contains(err.Error(), fmt.Sprintf(connectionErrorTemplate,
		"Failed to execute request"))

	assert.Contains(err.Error(), "no such host")

	assert.Contains(err.(*pnerr.ConnectionError).OrigError.Error(),
		"dial tcp: lookup")
}

// WARNING: not mocked request
func TestPublishSigned(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()
	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHANNEL

	pn := pubnub.NewPubNub(config)

	_, err := pn.Publish().Channel("ch").
		Message([]string{"hey", "hey2", "hey3"}).Execute()

	assert.Nil(err)
}

func TestPublishSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()
	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHANNEL

	pn := pubnub.NewPubNub(config)

	_, err := pn.Publish().Channel(SPECIAL_CHANNEL).
		Message([]string{SPECIAL_CHARACTERS, SPECIAL_CHARACTERS,
			SPECIAL_CHARACTERS}).Meta(SPECIAL_CHARACTERS).Execute()

	assert.Nil(err)
}
