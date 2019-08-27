package pubnub

import (
	"fmt"
	"reflect"
	"testing"

	h "github.com/sprucehealth/pubnub-go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessFetchGet(t *testing.T, expectedString string, channels []string) {
	assert := assert.New(t)

	opts := &fetchOpts{
		Channels:         channels,
		Reverse:          false,
		IncludeTimetoken: true,
		pubnub:           pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", expectedString),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)
}

func TestFetchQueryParam(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetQueryParam(t, "%22test%22?max=25&reverse=false", channels)
}

func AssertSuccessFetchGetQueryParam(t *testing.T, expectedString string, channels []string) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &fetchOpts{
		Channels:         channels,
		Reverse:          false,
		IncludeTimetoken: true,
		pubnub:           pubnub,
		QueryParam:       queryParam,
	}

	u, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

}

func AssertSuccessFetchQuery(t *testing.T, expectedString string, channels []string) {
	opts := &fetchOpts{
		Channels: channels,
		Reverse:  false,
		pubnub:   pubnub,
	}

	query, _ := opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))
	assert.Equal(t, "false", query.Get("reverse"))

}

func AssertNewFetchBuilder(t *testing.T, expectedString string, channels []string) {
	o := newFetchBuilder(pubnub)
	o.Channels(channels)
	o.Reverse(false)

	query, _ := o.opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))
	assert.Equal(t, "false", query.Get("reverse"))

}

func TestNewFetchBuilder(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertNewFetchBuilder(t, "%22test%22?max=25&reverse=false", channels)
}

func AssertNewFetchBuilderContext(t *testing.T, expectedString string, channels []string) {
	o := newFetchBuilderWithContext(pubnub, backgroundContext)
	o.Channels(channels)
	o.Reverse(false)

	query, _ := o.opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))
	assert.Equal(t, "false", query.Get("reverse"))

}

func TestNewFetchBuilderContext(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertNewFetchBuilderContext(t, "%22test%22?max=25&reverse=false", channels)
}

func TestFetchPath(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGet(t, "test1,test2", channels)
}

func TestFetchQuery(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchQuery(t, "%22test%22?max=25&reverse=false", channels)
}

func initFetchOpts(cipher string) *fetchOpts {
	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = cipher
	return &fetchOpts{
		Channels: []string{"test1,test2"},
		Reverse:  false,
		pubnub:   pn,
	}
}

func TestFetchResponseWithoutCipher(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"test":[{"message":"nyQDWnNPc1ryr5RgzVCKWw==","timetoken":"15229448184080121"}],"my-channel":[{"message":"nyQDWnNPc1ryr5RgzVCKWw==","timetoken":"15229448086016618"},{"message":"nyQDWnNPc1ryr5RgzVCKWw==","timetoken":"15229448126438499"},{"message":"my-message","timetoken":"15229450607090584"}]}}`)

	resp, _, err := newFetchResponse(jsonString, initFetchOpts(""), fakeResponseState)
	assert.Nil(err)

	respTest := resp.Messages["test"]
	respMyChannel := resp.Messages["my-channel"]

	assert.Equal("nyQDWnNPc1ryr5RgzVCKWw==", respTest[0].Message)
	assert.Equal("15229448184080121", respTest[0].Timetoken)

	assert.Equal("nyQDWnNPc1ryr5RgzVCKWw==", respMyChannel[0].Message)
	assert.Equal("15229448086016618", respMyChannel[0].Timetoken)
	assert.Equal("nyQDWnNPc1ryr5RgzVCKWw==", respMyChannel[1].Message)
	assert.Equal("15229448126438499", respMyChannel[1].Timetoken)
	assert.Equal("my-message", respMyChannel[2].Message)
	assert.Equal("15229450607090584", respMyChannel[2].Timetoken)

}

func TestFetchResponseWithCipher(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"test":[{"message":"Wi24KS4pcTzvyuGOHubiXg==","timetoken":"15229448184080121"}],"my-channel":[{"message":"Wi24KS4pcTzvyuGOHubiXg==","timetoken":"15229448086016618"},{"message":"Wi24KS4pcTzvyuGOHubiXg==","timetoken":"15229448126438499"},{"message":"my-message","timetoken":"15229450607090584"}]}}`)

	resp, _, err := newFetchResponse(jsonString, initFetchOpts("enigma"), fakeResponseState)
	assert.Nil(err)

	respTest := resp.Messages["test"]
	respMyChannel := resp.Messages["my-channel"]

	assert.Equal("yay!", respTest[0].Message)
	assert.Equal("15229448184080121", respTest[0].Timetoken)

	assert.Equal("yay!", respMyChannel[0].Message)
	assert.Equal("15229448086016618", respMyChannel[0].Timetoken)
	assert.Equal("yay!", respMyChannel[1].Message)
	assert.Equal("15229448126438499", respMyChannel[1].Timetoken)
	assert.Equal("my-message", respMyChannel[2].Message)
	assert.Equal("15229450607090584", respMyChannel[2].Timetoken)

}

func TestFetchResponseWithCipherInterface(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"test":[{"message":"{\"not_other\":\"1234\", \"pn_other\":\"yay!\"}","timetoken":"15229448184080121"}],"my-channel":[{"message":{"not_other":"1234", "pn_other":"Wi24KS4pcTzvyuGOHubiXg=="},"timetoken":"15229448086016618"},{"message":"bCC/kQbGdScQ0teYcawUsunrJLoUdp3Mwb/24ifa87QDBWv5aTkXkkjVMMXizEDb","timetoken":"15229448126438499"},{"message":"my-message","timetoken":"15229450607090584"}]}}`)

	resp, _, err := newFetchResponse(jsonString, initFetchOpts("enigma"), fakeResponseState)
	assert.Nil(err)

	respTest := resp.Messages["test"]
	respMyChannel := resp.Messages["my-channel"]

	assert.Equal("{\"not_other\":\"1234\", \"pn_other\":\"yay!\"}", respTest[0].Message)
	assert.Equal("15229448184080121", respTest[0].Timetoken)

	data := respMyChannel[0].Message
	switch v := data.(type) {
	case map[string]interface{}:
		testMap := make(map[string]interface{})
		testMap = v
		assert.Equal(testMap["not_other"], "1234")
		assert.Equal(testMap["pn_other"], "yay!")

		break
	default:
		assert.Fail(string(reflect.TypeOf(data).Kind()), " expected interafce")
		break
	}

	assert.Equal("15229448086016618", respMyChannel[0].Timetoken)
	if testMap, ok := respMyChannel[1].Message.(map[string]interface{}); !ok {
		assert.Fail("respMyChannel[1].Message ! map[string]interface{}")
	} else {
		assert.Equal("1234", testMap["not_other"])
		assert.Equal("yay!", testMap["pn_other"])
	}
	assert.Equal("15229448126438499", respMyChannel[1].Timetoken)
	assert.Equal("my-message", respMyChannel[2].Message)
	assert.Equal("15229450607090584", respMyChannel[2].Timetoken)

}

func TestFetchResponseWithCipherInterfacePNOtherDisabled(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"
	pn.Config.DisablePNOtherProcessing = true
	f := &fetchOpts{
		Channels: []string{"test1,test2"},
		Reverse:  false,
		pubnub:   pn,
	}

	jsonString := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"test":[{"message":"{\"not_other\":\"1234\", \"pn_other\":\"yay!\"}","timetoken":"15229448184080121"}],"my-channel":[{"message":{"not_other":"1234", "pn_other":"Wi24KS4pcTzvyuGOHubiXg=="},"timetoken":"15229448086016618"},{"message":"bCC/kQbGdScQ0teYcawUsunrJLoUdp3Mwb/24ifa87QDBWv5aTkXkkjVMMXizEDb","timetoken":"15229448126438499"},{"message":"my-message","timetoken":"15229450607090584"}]}}`)

	resp, _, err := newFetchResponse(jsonString, f, fakeResponseState)
	assert.Nil(err)

	respTest := resp.Messages["test"]
	respMyChannel := resp.Messages["my-channel"]

	assert.Equal("{\"not_other\":\"1234\", \"pn_other\":\"yay!\"}", respTest[0].Message)
	assert.Equal("15229448184080121", respTest[0].Timetoken)

	data := respMyChannel[0].Message
	switch v := data.(type) {
	case map[string]interface{}:
		testMap := make(map[string]interface{})
		testMap = v
		assert.Equal(testMap["not_other"], "1234")
		assert.Equal(testMap["pn_other"], "Wi24KS4pcTzvyuGOHubiXg==")

		break
	default:
		assert.Fail(string(reflect.TypeOf(data).Kind()), " expected interafce")
		break
	}

	assert.Equal("15229448086016618", respMyChannel[0].Timetoken)
	if testMap, ok := respMyChannel[1].Message.(map[string]interface{}); !ok {
		assert.Fail("respMyChannel[1].Message ! map[string]interface{}")
	} else {
		assert.Equal("1234", testMap["not_other"])
		assert.Equal("yay!", testMap["pn_other"])
	}
	assert.Equal("15229448126438499", respMyChannel[1].Timetoken)
	assert.Equal("my-message", respMyChannel[2].Message)
	assert.Equal("15229450607090584", respMyChannel[2].Timetoken)
	pn.Config.CipherKey = ""

}

func TestFireValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &fetchOpts{
		Reverse: false,
		pubnub:  pn,
	}

	assert.Equal("pubnub/validation: pubnub: \x06: Missing Subscribe Key", opts.validate().Error())
}

func TestFireValidateCH(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &fetchOpts{
		Reverse: false,
		pubnub:  pn,
	}
	assert.Equal("pubnub/validation: pubnub: \x06: Missing Channel", opts.validate().Error())
}

func TestNewFetchResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &fetchOpts{
		Reverse: false,
		pubnub:  pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newFetchResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}
