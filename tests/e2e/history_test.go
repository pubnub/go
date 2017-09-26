package e2e

import (
	"log"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

const HISTORY_RESP_SUCCESS = `[[{"timetoken":1111,"message":{"a":11,"b":22}},{"timetoken":2222,"message":{"a":33,"b":44}}],1234,4321]`

func TestHistorySuccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.History().Channel("ch").Execute()

	assert.Nil(err)
}

func TestHistoryCallWithAllParams(t *testing.T) {
	assert := assert.New(t)
	pn := pubnub.NewPubNub(configCopy())

	res, _, err := pn.History().
		Channel("ch").
		Count(2).
		IncludeTimetoken(true).
		Reverse(true).
		Start(int64(1)).
		End(int64(2)).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
}

func TestHistorySuccess(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/history/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel/ch",
		Query:              "count=100&include_token=false&reverse=false",
		ResponseBody:       HISTORY_RESP_SUCCESS,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.History().
		Channel("ch").
		Transport(interceptor.Transport).
		Execute()

	log.Println(res)

	assert.Nil(err)
	assert.Equal(int64(1234), res.StartTimetoken)
	assert.Equal(int64(4321), res.EndTimetoken)
	assert.Equal(2, len(res.Messages))
	assert.Equal(int64(1111), res.Messages[0].Timetoken)
	assert.Equal(map[string]interface{}{"a": float64(11), "b": float64(22)},
		res.Messages[0].Message)
	assert.Equal(int64(2222), res.Messages[1].Timetoken)
	assert.Equal(map[string]interface{}{"a": float64(33), "b": float64(44)},
		res.Messages[1].Message)
}

func TestHistoryEncryptedPNOther(t *testing.T) {
	assert := assert.New(t)

	config.CipherKey = "hello"

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/history/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel/ch",
		Query:              "count=100&include_token=false&reverse=false",
		ResponseBody:       `[[{"pn_other":"6QoqmS9CnB3W9+I4mhmL7w=="}],14606134331557852,14606134485013970]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "timestamp", "signature"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.History().
		Channel("ch").
		Transport(interceptor.Transport).
		Execute()

	assert.Nil(err)
	assert.Equal(1, len(res.Messages))
	assert.Equal(map[string]interface{}{"text": "hey"}, res.Messages[0].Message)

	config.CipherKey = ""
}

func TestHistoryMissingChannel(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	res, _, err := pn.History().
		Channel("").
		Execute()

	assert.Nil(res)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestHistoryPNOtherError(t *testing.T) {
	assert := assert.New(t)

	config.CipherKey = "hello"

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/history/sub-key/sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f/channel/ch",
		Query:              "count=100&include_token=false&reverse=false",
		ResponseBody:       `[[{"pn_other":""}],14606134331557852,14606134485013970]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "timestamp", "signature"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.History().
		Channel("ch").
		Transport(interceptor.Transport).
		Execute()

	assert.Nil(res)
	assert.Contains(err.Error(), "message is empty")

	config.CipherKey = ""
}

func TestHistorySuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	config.Uuid = SPECIAL_CHARACTERS
	config.AuthKey = SPECIAL_CHARACTERS

	pn := pubnub.NewPubNub(pamConfigCopy())

	_, _, err := pn.History().
		Channel(SPECIAL_CHANNEL).
		Count(100).
		Reverse(true).
		IncludeTimetoken(true).
		Execute()

	assert.Nil(err)
}
