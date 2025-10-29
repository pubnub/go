package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	h "github.com/pubnub/go/v8/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessFetchGet(t *testing.T, expectedString string, channels []string) {
	assert := assert.New(t)

	opts := newFetchOpts(pubnub, pubnub.ctx, fetchOpts{
		Channels: channels,
	})
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
	AssertSuccessFetchGetQueryParam(t, "%22test%22?max=25", channels)
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

	opts := newFetchOpts(pubnub, pubnub.ctx, fetchOpts{
		Channels:   channels,
		QueryParam: queryParam,
	})
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

	opts := newFetchOpts(pubnub, pubnub.ctx, fetchOpts{
		Channels:   channels,
		QueryParam: queryParam,
	})

	u, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

}

func TestSuccessFetchQueryOneChannel(t *testing.T) {
	opts := newFetchOpts(pubnub, pubnub.ctx, fetchOpts{
		Channels: []string{"ch"},
	})

	query, _ := opts.buildQuery()

	assert.Equal(t, "100", query.Get("max"))

}

func TestSuccessFetchQueryOneChannelWithMessageActions(t *testing.T) {
	opts := newFetchOpts(pubnub, pubnub.ctx, fetchOpts{
		Channels:           []string{"ch"},
		WithMessageActions: true,
	})

	query, _ := opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))

}

func AssertSuccessFetchQuery(t *testing.T, expectedString string, channels []string) {
	opts := newFetchOpts(pubnub, pubnub.ctx, fetchOpts{
		Channels: channels,
	})

	query, _ := opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))

}

func AssertNewFetchBuilder(t *testing.T, expectedString string, channels []string) {
	o := newFetchBuilder(pubnub)
	o.Channels(channels)

	query, _ := o.opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))

}

func TestNewFetchBuilder(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertNewFetchBuilder(t, "%22test%22?max=25", channels)
}

func AssertNewFetchBuilderContext(t *testing.T, expectedString string, channels []string) {
	o := newFetchBuilderWithContext(pubnub, pubnub.ctx)
	o.Channels(channels)

	query, _ := o.opts.buildQuery()

	assert.Equal(t, "25", query.Get("max"))

}

func TestNewFetchBuilderContext(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertNewFetchBuilderContext(t, "%22test%22?max=25", channels)
}

func TestFetchPath(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchGet(t, "test1,test2", channels)
}

func TestFetchQuery(t *testing.T) {
	channels := []string{"test1", "test2"}
	AssertSuccessFetchQuery(t, "%22test%22?max=25", channels)
}

func initFetchOpts(cipher string) *fetchOpts {
	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = cipher
	pn.Config.UseRandomInitializationVector = false
	return newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"test1,test2"},
	})
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
	pn.Config.UseRandomInitializationVector = false
	f := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"test1,test2"},
	})

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
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{})

	assert.Equal("pubnub/validation: pubnub: Fetch Messages: Missing Subscribe Key", opts.validate().Error())
}

func TestFireValidateCH(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{})
	assert.Equal("pubnub/validation: pubnub: Fetch Messages: Missing Channel", opts.validate().Error())
}

func TestNewFetchResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{})
	jsonBytes := []byte(`s`)

	_, _, err := newFetchResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

// Comprehensive Validation Tests

func TestFetchValidateWithMessageActionsMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels:           []string{"channel1", "channel2"},
		WithMessageActions: true,
	})

	err := opts.validate()
	assert.NotNil(err)
	assert.Equal("pubnub/validation: pubnub: Fetch Messages: Only one channel is supported when WithMessageActions is true", err.Error())
}

func TestFetchValidateWithMessageActionsSingleChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels:           []string{"channel1"},
		WithMessageActions: true,
	})

	assert.Nil(opts.validate())
}

func TestFetchValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"channel1", "channel2"},
	})

	assert.Nil(opts.validate())
}

func TestFetchValidateSuccessComplexParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels:        []string{"channel1", "channel2"},
		Start:           123456789,
		End:             987654321,
		Count:           25,
		WithMeta:        true,
		WithUUID:        true,
		WithMessageType: true,
	})
	opts.setStart = true
	opts.setEnd = true

	assert.Nil(opts.validate())
}

// HTTP Method and Operation Tests

func TestFetchHTTPMethod(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"test-channel"},
	})

	// Fetch should use GET method (default when httpMethod() not defined)
	// We can verify this by checking that buildBody() returns empty
	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body) // GET requests have empty body
}

func TestFetchOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"test-channel"},
	})

	assert.Equal(PNFetchMessagesOperation, opts.operationType())
}

// Comprehensive Builder Pattern Tests (11 setters!)

func TestFetchBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newFetchBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestFetchBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newFetchBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestFetchBuilderSettersIndividual(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newFetchBuilder(pn)

	// Test Channels setter
	channels := []string{"channel1", "channel2"}
	builder.Channels(channels)
	assert.Equal(channels, builder.opts.Channels)

	// Test Start setter
	builder.Start(123456789)
	assert.Equal(int64(123456789), builder.opts.Start)
	assert.True(builder.opts.setStart)

	// Test End setter
	builder.End(987654321)
	assert.Equal(int64(987654321), builder.opts.End)
	assert.True(builder.opts.setEnd)

	// Test Count setter
	builder.Count(50)
	assert.Equal(50, builder.opts.Count)

	// Test IncludeMeta setter
	builder.IncludeMeta(true)
	assert.True(builder.opts.WithMeta)

	// Test IncludeMessageActions setter
	builder.IncludeMessageActions(true)
	assert.True(builder.opts.WithMessageActions)

	// Test IncludeUUID setter
	builder.IncludeUUID(false)
	assert.False(builder.opts.WithUUID)

	// Test IncludeMessageType setter
	builder.IncludeMessageType(false)
	assert.False(builder.opts.WithMessageType)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestFetchBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	channels := []string{"channel1", "channel2"}
	queryParam := map[string]string{"key": "value"}

	builder := newFetchBuilder(pn)
	result := builder.Channels(channels).
		Start(123456789).
		End(987654321).
		Count(25).
		IncludeMeta(true).
		IncludeMessageActions(false).
		IncludeUUID(false).
		IncludeMessageType(false).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(int64(123456789), builder.opts.Start)
	assert.Equal(int64(987654321), builder.opts.End)
	assert.Equal(25, builder.opts.Count)
	assert.True(builder.opts.WithMeta)
	assert.False(builder.opts.WithMessageActions)
	assert.False(builder.opts.WithUUID)
	assert.False(builder.opts.WithMessageType)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.True(builder.opts.setStart)
	assert.True(builder.opts.setEnd)
}

func TestFetchBuilderTransport(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newFetchBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}

// Complex Query Building Tests (Dynamic Max Count Logic)

func TestFetchBuildQuerySingleChannelMaxCount(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"channel1"},
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("max"))            // maxCountFetch = 100 for single channel
	assert.Equal("false", query.Get("include_meta")) // Default is false, not true
	assert.Equal("true", query.Get("include_message_type"))
	assert.Equal("true", query.Get("include_uuid"))
}

func TestFetchBuildQueryMultipleChannelsMaxCount(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"channel1", "channel2"},
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("25", query.Get("max")) // maxCountFetchMoreThanOneChannel = 25
}

func TestFetchBuildQueryWithMessageActionsMaxCount(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels:           []string{"channel1"},
		WithMessageActions: true,
	})

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("25", query.Get("max")) // maxCountHistoryWithMessageActions = 25
}

func TestFetchBuildQueryCountBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		channels    []string
		withActions bool
		count       int
		expectedMax string
	}{
		{
			name:        "Single channel, count within limit",
			channels:    []string{"ch1"},
			withActions: false,
			count:       50,
			expectedMax: "50",
		},
		{
			name:        "Single channel, count over limit",
			channels:    []string{"ch1"},
			withActions: false,
			count:       150,
			expectedMax: "100", // Capped at maxCountFetch
		},
		{
			name:        "Multiple channels, count within limit",
			channels:    []string{"ch1", "ch2"},
			withActions: false,
			count:       20,
			expectedMax: "20",
		},
		{
			name:        "Multiple channels, count over limit",
			channels:    []string{"ch1", "ch2"},
			withActions: false,
			count:       50,
			expectedMax: "25", // Capped at maxCountFetchMoreThanOneChannel
		},
		{
			name:        "With actions, count within limit",
			channels:    []string{"ch1"},
			withActions: true,
			count:       20,
			expectedMax: "20",
		},
		{
			name:        "With actions, count over limit",
			channels:    []string{"ch1"},
			withActions: true,
			count:       50,
			expectedMax: "25", // Capped at maxCountHistoryWithMessageActions
		},
		{
			name:        "Zero count uses default",
			channels:    []string{"ch1"},
			withActions: false,
			count:       0,
			expectedMax: "100", // Uses default maxCountFetch
		},
		{
			name:        "Negative count uses default",
			channels:    []string{"ch1"},
			withActions: false,
			count:       -5,
			expectedMax: "100", // Uses default maxCountFetch
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newFetchOpts(pn, pn.ctx, fetchOpts{
				Channels:           tc.channels,
				WithMessageActions: tc.withActions,
				Count:              tc.count,
			})

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectedMax, query.Get("max"))
		})
	}
}

func TestFetchBuildQueryBooleanFlagCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name               string
		withMeta           bool
		withMessageActions bool
		withUUID           bool
		withMessageType    bool
	}{
		{
			name:               "All flags false",
			withMeta:           false,
			withMessageActions: false,
			withUUID:           false,
			withMessageType:    false,
		},
		{
			name:               "All flags true",
			withMeta:           true,
			withMessageActions: false, // Cannot be true with multiple channels
			withUUID:           true,
			withMessageType:    true,
		},
		{
			name:               "Mixed flags 1",
			withMeta:           false,
			withMessageActions: false,
			withUUID:           true,
			withMessageType:    false,
		},
		{
			name:               "Mixed flags 2",
			withMeta:           true,
			withMessageActions: false,
			withUUID:           false,
			withMessageType:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newFetchOpts(pn, pn.ctx, fetchOpts{
				Channels: []string{"channel1", "channel2"},
			})
			// Set the flags after creation to override defaults
			opts.WithMeta = tc.withMeta
			opts.WithUUID = tc.withUUID
			opts.WithMessageType = tc.withMessageType

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(strconv.FormatBool(tc.withMeta), query.Get("include_meta"))
			assert.Equal(strconv.FormatBool(tc.withUUID), query.Get("include_uuid"))
			assert.Equal(strconv.FormatBool(tc.withMessageType), query.Get("include_message_type"))
		})
	}
}

func TestFetchBuildQueryTimetokenHandling(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		start       int64
		end         int64
		setStart    bool
		setEnd      bool
		expectStart string
		expectEnd   string
	}{
		{
			name:        "No timetokens",
			start:       0,
			end:         0,
			setStart:    false,
			setEnd:      false,
			expectStart: "",
			expectEnd:   "",
		},
		{
			name:        "Only start",
			start:       123456789,
			end:         0,
			setStart:    true,
			setEnd:      false,
			expectStart: "123456789",
			expectEnd:   "",
		},
		{
			name:        "Only end",
			start:       0,
			end:         987654321,
			setStart:    false,
			setEnd:      true,
			expectStart: "",
			expectEnd:   "987654321",
		},
		{
			name:        "Both timetokens",
			start:       123456789,
			end:         987654321,
			setStart:    true,
			setEnd:      true,
			expectStart: "123456789",
			expectEnd:   "987654321",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newFetchOpts(pn, pn.ctx, fetchOpts{
				Channels: []string{"channel1"},
				Start:    tc.start,
				End:      tc.end,
			})
			opts.setStart = tc.setStart
			opts.setEnd = tc.setEnd

			query, err := opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectStart, query.Get("start"))
			assert.Equal(tc.expectEnd, query.Get("end"))
		})
	}
}

func TestFetchBuildQueryWithCustomParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	customParams := map[string]string{
		"custom":         "value",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}

	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels:   []string{"channel1"},
		QueryParam: customParams,
	})

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all custom parameters are present
	for key, expectedValue := range customParams {
		actualValue := query.Get(key)
		if key == "special_chars" {
			// Special characters should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should be URL encoded", key)
		} else if key == "unicode" {
			// Unicode should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should contain URL encoded Unicode", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}
}

// Conditional Path Building Tests

func TestFetchBuildPathWithoutMessageActions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: []string{"channel1", "channel2"},
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v3/history/sub-key/demo/channel/channel1,channel2"
	assert.Equal(expected, path)
}

func TestFetchBuildPathWithMessageActions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels:           []string{"channel1"},
		WithMessageActions: true,
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v3/history-with-actions/sub-key/demo/channel/channel1"
	assert.Equal(expected, path)
}

func TestFetchBuildPathChannelJoining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name             string
		channels         []string
		withActions      bool
		expectedPath     string
		expectedBasePath string
	}{
		{
			name:             "Single channel",
			channels:         []string{"test-channel"},
			withActions:      false,
			expectedPath:     "/v3/history/sub-key/demo/channel/test-channel",
			expectedBasePath: "/v3/history/",
		},
		{
			name:             "Two channels",
			channels:         []string{"channel1", "channel2"},
			withActions:      false,
			expectedPath:     "/v3/history/sub-key/demo/channel/channel1,channel2",
			expectedBasePath: "/v3/history/",
		},
		{
			name:             "Three channels",
			channels:         []string{"ch1", "ch2", "ch3"},
			withActions:      false,
			expectedPath:     "/v3/history/sub-key/demo/channel/ch1,ch2,ch3",
			expectedBasePath: "/v3/history/",
		},
		{
			name:             "Single channel with actions",
			channels:         []string{"action-channel"},
			withActions:      true,
			expectedPath:     "/v3/history-with-actions/sub-key/demo/channel/action-channel",
			expectedBasePath: "/v3/history-with-actions/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newFetchOpts(pn, pn.ctx, fetchOpts{
				Channels:           tc.channels,
				WithMessageActions: tc.withActions,
			})

			path, err := opts.buildPath()
			assert.Nil(err)
			assert.Equal(tc.expectedPath, path)
			assert.Contains(path, tc.expectedBasePath)
		})
	}
}

func TestFetchBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannels := []string{
		"channel@with%encoded",
		"channel/with/slashes",
		"channel?with=query&chars",
		"channel#with#hashes",
		"channel with spaces and símböls",
		"测试频道-русский-チャンネル",
	}

	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: specialChannels,
	})

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v3/history/sub-key/demo/channel/")
	// Should contain comma-separated channels
	assert.Contains(path, ",")
}

// Comprehensive Edge Case Tests

func TestFetchWithUnicodeChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	channels := []string{"测试频道", "русский-канал", "チャンネル", "한국어채널"}

	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: channels,
	})

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v3/history/sub-key/demo/channel/")
}

func TestFetchWithManyChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create 50 channels
	channels := make([]string, 50)
	for i := 0; i < 50; i++ {
		channels[i] = fmt.Sprintf("channel_%d", i)
	}

	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: channels,
	})

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "channel_0")
	assert.Contains(path, "channel_49")
	assert.Contains(path, ",") // Should have comma separators

	// Should use multiple channel max count
	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("25", query.Get("max"))
}

func TestFetchWithExtremeTimetokens(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newFetchBuilder(pn)
	builder.Channels([]string{"test-channel"})

	// Test extreme timetoken values
	maxTimetoken := int64(9223372036854775807) // Max int64
	minTimetoken := int64(1)

	builder.Start(minTimetoken)
	builder.End(maxTimetoken)

	assert.Equal(minTimetoken, builder.opts.Start)
	assert.Equal(maxTimetoken, builder.opts.End)

	// Test query building with extreme values
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1", query.Get("start"))
	assert.Equal("9223372036854775807", query.Get("end"))
}

func TestFetchWithVeryLongChannelNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create very long channel names
	longChannels := make([]string, 3)
	for i := 0; i < 3; i++ {
		longName := ""
		for j := 0; j < 100; j++ {
			longName += fmt.Sprintf("segment_%d_%d_", i, j)
		}
		longChannels[i] = longName
	}

	opts := newFetchOpts(pn, pn.ctx, fetchOpts{
		Channels: longChannels,
	})

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/v3/history/")
	assert.Contains(path, "segment_0_0_")
	assert.Contains(path, "segment_2_99_")
}

func TestFetchParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name     string
		channels []string
		start    int64
		end      int64
		count    int
		meta     bool
		actions  bool
		uuid     bool
		msgType  bool
	}{
		{
			name:     "Minimal parameters",
			channels: []string{"ch1"},
		},
		{
			name:     "All parameters enabled (single channel)",
			channels: []string{"ch1"},
			start:    123456789,
			end:      987654321,
			count:    25,
			meta:     true,
			actions:  true,
			uuid:     true,
			msgType:  true,
		},
		{
			name:     "All parameters enabled (multiple channels, no actions)",
			channels: []string{"ch1", "ch2", "ch3"},
			start:    123456789,
			end:      987654321,
			count:    25,
			meta:     true,
			actions:  false, // Cannot be true with multiple channels
			uuid:     true,
			msgType:  true,
		},
		{
			name:     "Mixed parameters",
			channels: []string{"ch1", "ch2"},
			start:    555666777,
			count:    15,
			meta:     true,
			uuid:     false,
			msgType:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newFetchBuilder(pn)
			builder.Channels(tc.channels)

			if tc.start != 0 {
				builder.Start(tc.start)
			}
			if tc.end != 0 {
				builder.End(tc.end)
			}
			if tc.count != 0 {
				builder.Count(tc.count)
			}

			builder.IncludeMeta(tc.meta)
			builder.IncludeMessageActions(tc.actions)
			builder.IncludeUUID(tc.uuid)
			builder.IncludeMessageType(tc.msgType)

			// Should pass validation
			assert.Nil(builder.opts.validate())

			// Should build valid path
			path, err := builder.opts.buildPath()
			assert.Nil(err)
			if tc.actions {
				assert.Contains(path, "/v3/history-with-actions/")
			} else {
				assert.Contains(path, "/v3/history/")
			}

			// Should build valid query
			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.NotNil(query)
		})
	}
}

// Error Scenario Tests

func TestFetchExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newFetchBuilder(pn)
	builder.Channels([]string{"test-channel"})

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestFetchExecuteErrorWithMessageActionsMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newFetchBuilder(pn)
	builder.Channels([]string{"channel1", "channel2"})
	builder.IncludeMessageActions(true)

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Only one channel is supported when WithMessageActions is true")
}

func TestFetchPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	edgeCases := []struct {
		name        string
		channels    []string
		withActions bool
	}{
		{
			name:        "Empty channel name",
			channels:    []string{""},
			withActions: false,
		},
		{
			name:        "Channel with only spaces",
			channels:    []string{"   "},
			withActions: false,
		},
		{
			name:        "Very special characters",
			channels:    []string{"!@#$%^&*()_+-=[]{}|;':\",./<>?"},
			withActions: false,
		},
		{
			name:        "Mix of normal and special",
			channels:    []string{"normal", "special@#$%", "unicode测试"},
			withActions: false,
		},
		{
			name:        "Single special channel with actions",
			channels:    []string{"action@channel#with$symbols"},
			withActions: true,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newFetchOpts(pn, pn.ctx, fetchOpts{
				Channels:           tc.channels,
				WithMessageActions: tc.withActions,
			})

			// Should pass validation
			assert.Nil(opts.validate(), "Should validate for case: %s", tc.name)

			// Should build valid path
			path, err := opts.buildPath()
			assert.Nil(err, "Should build path for case: %s", tc.name)

			if tc.withActions {
				assert.Contains(path, "/v3/history-with-actions/", "Should contain actions path for: %s", tc.name)
			} else {
				assert.Contains(path, "/v3/history/", "Should contain history path for: %s", tc.name)
			}
		})
	}
}

func TestFetchQueryBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name        string
		queryParam  map[string]string
		expectError bool
	}{
		{
			name:        "Empty query params",
			queryParam:  map[string]string{},
			expectError: false,
		},
		{
			name:        "Nil query params",
			queryParam:  nil,
			expectError: false,
		},
		{
			name: "Very large query params",
			queryParam: map[string]string{
				"param1": strings.Repeat("a", 1000),
				"param2": strings.Repeat("b", 1000),
				"param3": strings.Repeat("c", 1000),
			},
			expectError: false,
		},
		{
			name: "Special character query params",
			queryParam: map[string]string{
				"special@key":   "special@value",
				"unicode测试":     "unicode值",
				"with spaces":   "also spaces",
				"equals=key":    "equals=value",
				"ampersand&key": "ampersand&value",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := newFetchOpts(pn, pn.ctx, fetchOpts{
				Channels:   []string{"test-channel"},
				QueryParam: tc.queryParam,
			})

			query, err := opts.buildQuery()
			if tc.expectError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(query)
			}
		})
	}
}

func TestFetchBuilderCompleteness(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test that all parameters can be set and retrieved
	builder := newFetchBuilder(pn)

	channels := []string{"complete-test-channel-1", "complete-test-channel-2"}
	queryParam := map[string]string{
		"custom1": "value1",
		"custom2": "value2",
	}

	// Set all possible parameters
	builder.Channels(channels).
		Start(111111111).
		End(999999999).
		Count(15).
		IncludeMeta(true).
		IncludeMessageActions(false). // false because multiple channels
		IncludeUUID(false).
		IncludeMessageType(false).
		QueryParam(queryParam)

	// Verify all values are set
	assert.Equal(channels, builder.opts.Channels)
	assert.Equal(int64(111111111), builder.opts.Start)
	assert.Equal(int64(999999999), builder.opts.End)
	assert.Equal(15, builder.opts.Count)
	assert.True(builder.opts.WithMeta)
	assert.False(builder.opts.WithMessageActions)
	assert.False(builder.opts.WithUUID)
	assert.False(builder.opts.WithMessageType)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.True(builder.opts.setStart)
	assert.True(builder.opts.setEnd)

	// Should pass validation
	assert.Nil(builder.opts.validate())

	// Should build correct path
	path, err := builder.opts.buildPath()
	assert.Nil(err)
	expectedPath := "/v3/history/sub-key/demo/channel/complete-test-channel-1,complete-test-channel-2"
	assert.Equal(expectedPath, path)

	// Should build query with custom params
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("value1", query.Get("custom1"))
	assert.Equal("value2", query.Get("custom2"))
	assert.Equal("15", query.Get("max"))
	assert.Equal("true", query.Get("include_meta"))
	assert.Equal("false", query.Get("include_uuid"))
	assert.Equal("false", query.Get("include_message_type"))
}
