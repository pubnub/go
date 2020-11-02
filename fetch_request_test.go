package pubnub

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessFetchGet(t *testing.T, expectedString string, channels []string) {
	assert := assert.New(t)

	opts := &fetchOpts{
		Channels: channels,
		Reverse:  false,
		pubnub:   pubnub,
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

func TestFetchMetaAndActions(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, true, false, false)
}

func TestFetchActions(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, false, true, false, false)
}

func TestFetchMeta(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, false, false, false)
}

func TestFetchMetaAndActionsWithUUID(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, true, true, false)
}

func TestFetchActionsWithUUID(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, false, true, true, false)
}

func TestFetchMetaWithUUID(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, false, true, false)
}

func TestFetchMetaAndActionsWithMT(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, true, false, true)
}

func TestFetchActionsWithMT(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, false, true, false, true)
}

func TestFetchMetaWithMT(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, false, false, true)
}

func TestFetchMetaAndActionsWithUUIDAndMT(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, true, true, true)
}

func TestFetchActionsWithUUIDAndMT(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, false, true, true, true)
}

func TestFetchMetaWithUUIDAndMT(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGetMetaAndActions(t, "test1,test2", channels, true, false, true, true)
}

func AssertSuccessFetchGetMetaAndActions(t *testing.T, expectedString string, channels []string, withMeta, withMessageActions, withUUID, withMessageType bool) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &fetchOpts{
		Channels:   channels,
		Reverse:    false,
		pubnub:     pubnub,
		QueryParam: queryParam,
	}
	opts.WithMeta = withMeta
	opts.WithMessageActions = withMessageActions
	if withUUID {
		opts.WithUUID = withUUID
	}
	if withMessageType {
		opts.WithMessageType = withMessageType
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	if withMessageActions {

		h.AssertPathsEqual(t,
			fmt.Sprintf("/v3/history-with-actions/sub-key/sub_key/channel/%s", expectedString),
			u.EscapedPath(), []int{})
	} else {
		h.AssertPathsEqual(t,
			fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", expectedString),
			u.EscapedPath(), []int{})
	}

	u1, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u1.Get("q1"))
	assert.Equal("v2", u1.Get("q2"))
	assert.Equal(strconv.FormatBool(withMeta), u1.Get("include_meta"))
	if withMessageType {
		assert.Equal(strconv.FormatBool(withMessageType), u1.Get("include_message_type"))
	}
	if withUUID {
		assert.Equal(strconv.FormatBool(withUUID), u1.Get("include_uuid"))
	}

}

func TestFetchMessageActionValidation(t *testing.T) {
	assert := assert.New(t)

	channels := []string{"test1", "test2"}
	o := newFetchBuilder(pubnub)
	o.Channels(channels)
	o.Reverse(false)
	o.IncludeMessageActions(true)

	_, _, e := o.Execute()
	assert.Equal("pubnub/validation: pubnub: Fetch Messages: Only one channel is supported when WithMessageActions is true", e.Error())

}

func AssertSuccessFetchGetQueryParam(t *testing.T, expectedString string, channels []string) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &fetchOpts{
		Channels:   channels,
		Reverse:    false,
		pubnub:     pubnub,
		QueryParam: queryParam,
	}

	u, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

}

func TestSuccessFetchQueryOneChannel(t *testing.T) {
	opts := &fetchOpts{
		Channels: []string{"ch"},
		Reverse:  false,
		pubnub:   pubnub,
	}

	query, _ := opts.buildQuery()

	assert.Equal(t, "100", query.Get("max"))
	assert.Equal(t, "false", query.Get("reverse"))

}

func TestSuccessFetchQueryOneChannelWithMessageActions(t *testing.T) {
	opts := &fetchOpts{
		Channels:           []string{"ch"},
		Reverse:            false,
		pubnub:             pubnub,
		WithMessageActions: true,
	}

	query, _ := opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))
	assert.Equal(t, "false", query.Get("reverse"))

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
		assert.Fail(fmt.Sprintf("%s", reflect.TypeOf(data).Kind()), " expected interafce")
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
		assert.Fail(fmt.Sprintf("%s", reflect.TypeOf(data).Kind()), " expected interafce")
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

func TestFetchResponseMetaCipher(t *testing.T) {
	FetchResponseMetaCommon(t, true)
}

func TestFetchResponseMeta(t *testing.T) {
	FetchResponseMetaCommon(t, false)
}

func FetchResponseMetaCommon(t *testing.T, withCipher bool) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "channels": {"my-channel": [{"message": "my-message", "timetoken": "15699986472636251", "meta": {"m1": "n1", "m2": "n2"}}]}, "error_message": "", "error": false}`)
	cipher := ""
	if withCipher {
		cipher = "enigma"
		jsonString = []byte(`{"status": 200, "channels": {"my-channel": [{"message": "6f+dRox3OKNiBHdGRT5HpA==", "timetoken": "15699986472636251", "meta": {"m1": "n1", "m2": "n2"}}]}, "error_message": "", "error": false}`)
	}

	resp, _, err := newFetchResponse(jsonString, initFetchOpts(cipher), fakeResponseState)
	assert.Nil(err)
	if resp != nil {
		messages := resp.Messages
		m0 := messages["my-channel"]
		if m0 != nil {
			assert.Equal("my-message", m0[0].Message)
			assert.Equal("15699986472636251", m0[0].Timetoken)
			meta := m0[0].Meta.(map[string]interface{})
			assert.Equal("n1", meta["m1"])
			assert.Equal("n2", meta["m2"])
		} else {
			assert.Fail("m0 nil")
		}
	} else {
		assert.Fail("res nil")
	}
}

func TestFetchResponseMessageTypeAndUUID(t *testing.T) {
	FetchResponseCommonForMessageTypeAndUUID(t, false)
}

func FetchResponseCommonForMessageTypeAndUUID(t *testing.T, withCipher bool) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"my-channel":[{"message_type": 4, "message": "my-message", "timetoken": "15959610984115342", "meta": "", "uuid": "db9c5e39-7c95-40f5-8d71-125765b6f561"}]}}`)

	resp, _, err := newFetchResponse(jsonString, initFetchOpts(""), fakeResponseState)
	assert.Nil(err)
	if resp != nil {
		messages := resp.Messages
		m0 := messages["my-channel"]
		if m0 != nil {
			assert.Equal("my-message", m0[0].Message)
			assert.Equal("15959610984115342", m0[0].Timetoken)
			assert.Equal(4, m0[0].MessageType)
			assert.Equal("db9c5e39-7c95-40f5-8d71-125765b6f561", m0[0].UUID)
		} else {
			assert.Fail("m0 nil")
		}
	} else {
		assert.Fail("res nil")
	}
}

func TestFetchResponseWithoutMessageTypeAndUUID(t *testing.T) {
	FetchResponseCommonForWithoutMessageTypeAndUUID(t, false)
}

func FetchResponseCommonForWithoutMessageTypeAndUUID(t *testing.T, withCipher bool) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "error": false, "error_message": "", "channels": {"my-channel":[{"message": "my-message", "timetoken": "15959610984115342", "meta": ""}]}}`)

	resp, _, err := newFetchResponse(jsonString, initFetchOpts(""), fakeResponseState)
	assert.Nil(err)
	if resp != nil {
		messages := resp.Messages
		m0 := messages["my-channel"]
		if m0 != nil {
			assert.Equal("my-message", m0[0].Message)
			assert.Equal("15959610984115342", m0[0].Timetoken)
		} else {
			assert.Fail("m0 nil")
		}
	} else {
		assert.Fail("res nil")
	}
}

func TestFetchResponseMetaAndActionsCipher(t *testing.T) {
	FetchResponseMetaAndActionsCommon(t, true)
}

func TestFetchResponseMetaAndActions(t *testing.T) {
	FetchResponseMetaAndActionsCommon(t, false)
}

func FetchResponseMetaAndActionsCommon(t *testing.T, withCipher bool) {
	assert := assert.New(t)

	jsonString := []byte(`{"status": 200, "channels": {"my-channel": [{"message": "my-message", "timetoken": "15699986472636251", "meta": {"m1": "n1", "m2": "n2"}, "actions": {"reaction2": {"smiley_face": [{"uuid": "pn-f3d10ae1-0437-4366-b509-0b5abd797a02", "actionTimetoken": "15700177371680470"}]}, "reaction": {"smiley_face": [{"uuid": "pn-f3d10ae1-0437-4366-b509-0b5abd797a02", "actionTimetoken": "15700177592799750"}, {"uuid": "pn-0e6345ab-529e-4fce-be3e-6bd041296661", "actionTimetoken": "15700010213930810"}], "frown_face": [{"uuid": "pn-f3d10ae1-0437-4366-b509-0b5abd797a02", "actionTimetoken": "15700177482326900"}]}}}]}, "error_message": "", "error": false}`)
	cipher := ""
	if withCipher {
		cipher = "enigma"
		jsonString = []byte(`{"status": 200, "channels": {"my-channel": [{"message": "6f+dRox3OKNiBHdGRT5HpA==", "timetoken": "15699986472636251", "meta": {"m1": "n1", "m2": "n2"}, "actions": {"reaction2": {"smiley_face": [{"uuid": "pn-f3d10ae1-0437-4366-b509-0b5abd797a02", "actionTimetoken": "15700177371680470"}]}, "reaction": {"smiley_face": [{"uuid": "pn-f3d10ae1-0437-4366-b509-0b5abd797a02", "actionTimetoken": "15700177592799750"}, {"uuid": "pn-0e6345ab-529e-4fce-be3e-6bd041296661", "actionTimetoken": "15700010213930810"}], "frown_face": [{"uuid": "pn-f3d10ae1-0437-4366-b509-0b5abd797a02", "actionTimetoken": "15700177482326900"}]}}}]}, "error_message": "", "error": false}`)
	}

	resp, _, err := newFetchResponse(jsonString, initFetchOpts(cipher), fakeResponseState)
	assert.Nil(err)

	if resp != nil {
		messages := resp.Messages
		m0 := messages["my-channel"]
		if m0 != nil {
			assert.Equal("my-message", m0[0].Message)
			assert.Equal("15699986472636251", m0[0].Timetoken)
			meta := m0[0].Meta.(map[string]interface{})
			assert.Equal("n1", meta["m1"])
			assert.Equal("n2", meta["m2"])
		} else {
			assert.Fail("m0 nil")
		}

		actions := m0[0].MessageActions

		if len(actions) > 0 {
			a0 := actions["reaction2"]
			r00 := a0.ActionsTypeValues["smiley_face"]
			if r00 != nil {
				assert.Equal("pn-f3d10ae1-0437-4366-b509-0b5abd797a02", r00[0].UUID)
				assert.Equal("15700177371680470", r00[0].ActionTimetoken)
			} else {
				assert.Fail("r0 nil")
			}

			a1 := actions["reaction"]
			r10 := a1.ActionsTypeValues["smiley_face"]
			if r10 != nil {
				assert.Equal("pn-f3d10ae1-0437-4366-b509-0b5abd797a02", r10[0].UUID)
				assert.Equal("15700177592799750", r10[0].ActionTimetoken)
				assert.Equal("pn-0e6345ab-529e-4fce-be3e-6bd041296661", r10[1].UUID)
				assert.Equal("15700010213930810", r10[1].ActionTimetoken)
			} else {
				assert.Fail("r0 nil")
			}
			r11 := a1.ActionsTypeValues["frown_face"]
			if r11 != nil {
				assert.Equal("pn-f3d10ae1-0437-4366-b509-0b5abd797a02", r11[0].UUID)
				assert.Equal("15700177482326900", r11[0].ActionTimetoken)
			} else {
				assert.Fail("r0 nil")
			}
		} else {
			assert.Fail("actions nil")
		}
	} else {
		assert.Fail("res nil")
	}
}

func TestFireValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &fetchOpts{
		Reverse: false,
		pubnub:  pn,
	}

	assert.Equal("pubnub/validation: pubnub: Fetch Messages: Missing Subscribe Key", opts.validate().Error())
}

func TestFireValidateCH(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &fetchOpts{
		Reverse: false,
		pubnub:  pn,
	}
	assert.Equal("pubnub/validation: pubnub: Fetch Messages: Missing Channel", opts.validate().Error())
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
