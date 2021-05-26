package e2e

import (
	"fmt"
	"log"
	"os"

	pubnub "github.com/pubnub/go/v5"
	"github.com/pubnub/go/v5/tests/stubs"
	"github.com/stretchr/testify/assert"

	"testing"
)

const historyResponseSuccess = `[[{"timetoken":1111,"message":{"a":11,"b":22}},{"timetoken":2222,"message":{"a":33,"b":44}}],1234,4321]`

func TestHistorySuccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.History().Channel("ch").Execute()

	assert.Nil(err)
}

func TestHistoryCallWithAllParams(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/history/sub-key/%s/channel/ch", config.SubscribeKey),
		Query:              "count=2&end=2&include_token=true&include_meta=false&reverse=true&start=1",
		ResponseBody:       `[[],0,0]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

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
		Path:               fmt.Sprintf("/v2/history/sub-key/%s/channel/ch", config.SubscribeKey),
		Query:              "count=100&include_token=false&include_meta=false&reverse=false",
		ResponseBody:       historyResponseSuccess,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.History().
		Channel("ch").
		Transport(interceptor.Transport).
		Execute()

	assert.Nil(err)
	if res != nil {
		assert.Equal(int64(1234), res.StartTimetoken)
		assert.Equal(int64(4321), res.EndTimetoken)
		assert.Equal(2, len(res.Messages))
		assert.Equal(int64(1111), res.Messages[0].Timetoken)
		assert.Equal(map[string]interface{}{"a": float64(11), "b": float64(22)},
			res.Messages[0].Message)
		assert.Equal(int64(2222), res.Messages[1].Timetoken)
		assert.Equal(map[string]interface{}{"a": float64(33), "b": float64(44)},
			res.Messages[1].Message)
	} else {
		assert.Fail("res nil")
	}
}

func TestHistorySuccessContext(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/history/sub-key/%s/channel/ch", config.SubscribeKey),
		Query:              "count=100&include_token=false&include_meta=false&reverse=false",
		ResponseBody:       historyResponseSuccess,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.HistoryWithContext(backgroundContext).
		Channel("ch").
		Transport(interceptor.Transport).
		Execute()
	if res != nil {
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
	} else {
		assert.Fail("res nil")
	}
}

func TestHistoryEncryptedPNOther(t *testing.T) {
	assert := assert.New(t)

	config.CipherKey = "hello"

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/history/sub-key/%s/channel/ch", config.SubscribeKey),
		Query:              "count=100&include_token=false&include_meta=false&reverse=false",
		ResponseBody:       `[[{"pn_other":"6QoqmS9CnB3W9+I4mhmL7w=="}],14606134331557852,14606134485013970]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "timestamp", "signature"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())
	pn.Config.UseRandomInitializationVector = false

	res, _, err := pn.History().
		Channel("ch").
		Transport(interceptor.Transport).
		Execute()

	assert.Nil(err)
	if res != nil {
		assert.Equal(1, len(res.Messages))

		if msgOther, ok := res.Messages[0].Message.(map[string]interface{}); !ok {
			assert.Fail("!map[string]interface{}")
		} else {
			if msgOther2, ok := msgOther["pn_other"].(map[string]interface{}); !ok {
				assert.Fail("!map[string]interface{} 2")
			} else {
				assert.Equal("hey", msgOther2["text"])
			}
		}
	} else {
		assert.Fail("res nil")
	}

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

/*func TestHistoryPNOtherError(t *testing.T) {
	assert := assert.New(t)

	config.CipherKey = "hello"

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/history/sub-key/sub-c-b9ab9508-43cf-11e8-9967-869954283fb4/channel/ch",
		Query:              "count=100&include_token=false&reverse=false",
		ResponseBody:       `[[{"pn_other":""}],14606134331557852,14606134485013970]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "timestamp", "signature"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.History().
		Channel("ch").
		Execute()

	assert.Equal(res, `"pn_other":""`)
	assert.Contains(err.Error(), "message is empty")

	config.CipherKey = ""
}*/

func TestHistorySuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters: /?#,
	validCharacters := "-._~:[]@!$&'()*+;=`|"

	config.UUID = validCharacters
	//config.AuthKey = validCharacters

	pn := pubnub.NewPubNub(config)

	_, _, err := pn.History().
		Channel(validCharacters).
		Count(100).
		Reverse(true).
		IncludeTimetoken(true).
		Execute()

	assert.Nil(err)
}
