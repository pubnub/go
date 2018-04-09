package pubnub

import (
	"fmt"
	"reflect"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
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

	assert.Equal("\"yay!\"", respTest[0].Message)
	assert.Equal("15229448184080121", respTest[0].Timetoken)

	assert.Equal("\"yay!\"", respMyChannel[0].Message)
	assert.Equal("15229448086016618", respMyChannel[0].Timetoken)
	assert.Equal("\"yay!\"", respMyChannel[1].Message)
	assert.Equal("15229448126438499", respMyChannel[1].Timetoken)
	assert.Equal("my-message", respMyChannel[2].Message)
	assert.Equal("15229450607090584", respMyChannel[2].Timetoken)

}

func TestFetchResponseWithCipherInterface(t *testing.T) {
	assert := assert.New(t)
	//fmt.Println(utils.EncryptString("enigma", "{\"not_other\":\"1234\", \"pn_other\":\"yay!\"}"))

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
		assert.Equal(testMap["pn_other"], "\"yay!\"")

		break
	default:
		assert.Fail(string(reflect.TypeOf(data).Kind()), " expected interafce")
		break
	}

	assert.Equal("15229448086016618", respMyChannel[0].Timetoken)
	assert.Equal("{\"not_other\":\"1234\", \"pn_other\":\"yay!\"}", respMyChannel[1].Message)
	assert.Equal("15229448126438499", respMyChannel[1].Timetoken)
	assert.Equal("my-message", respMyChannel[2].Message)
	assert.Equal("15229450607090584", respMyChannel[2].Timetoken)

}

/*func Fetch(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)

	channels := []string{"test1", "test2"}

	res, status, err := pn.Fetch().
		Channels(channels).
		Count(count).
		Reverse(reverse).
		Execute()

	assert.Nil(err)

	if status.Error == nil {
		for channel, messages := range res.Messages {
			fmt.Println("channel", channel)
			for _, messageInt := range messages {
				message := pubnub.FetchResponseItem(messageInt)
				fmt.Println(message.Message)
				fmt.Println(message.Timetoken)
			}
		}
	} else {
		fmt.Println("ParseFetch", err)
		fmt.Println("ParseFetch", status.StatusCode)
	}
}*/
