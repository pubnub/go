package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/pubnub/go/v7/crypto"
	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

var (
	fakeResponseState = StatusResponse{}
)

func (pn *PubNub) initHistoryOpts() *historyOpts {
	opts := newHistoryOpts(pubnub, pubnub.ctx)
	opts.Channel = "ch"
	opts.Start = int64(100000)
	opts.End = int64(200000)
	opts.setStart = true
	opts.setEnd = true
	opts.Reverse = false
	opts.Count = 3
	opts.IncludeTimetoken = true
	opts.pubnub = pn
	opts.includeCustomMessageType = false
	return opts
}

func TestHistoryRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := newHistoryOpts(pubnub, pubnub.ctx)
	opts.Channel = "ch"

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/history/sub-key/sub_key/channel/%s", opts.Channel),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("count", "100")
	expected.Set("reverse", "false")
	expected.Set("include_token", "false")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "include_meta"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewHistoryBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newHistoryBuilder(pubnub)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/history/sub-key/sub_key/channel/%s", o.opts.Channel),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("count", "100")
	expected.Set("reverse", "false")
	expected.Set("include_token", "false")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "include_meta"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewHistoryBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newHistoryBuilderWithContext(pubnub, pubnub.ctx)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/history/sub-key/sub_key/channel/%s", o.opts.Channel),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("count", "100")
	expected.Set("reverse", "false")
	expected.Set("include_token", "false")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "include_meta"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHistoryRequestAllParams(t *testing.T) {
	assert := assert.New(t)

	opts := newHistoryOpts(pubnub, pubnub.ctx)
	opts.Channel = "ch"
	opts.Start = int64(100000)
	opts.End = int64(200000)
	opts.setStart = true
	opts.setEnd = true
	opts.Reverse = false
	opts.Count = 3
	opts.IncludeTimetoken = true
	opts.includeCustomMessageType = true

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/history/sub-key/sub_key/channel/%s", opts.Channel),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("start", "100000")
	expected.Set("end", "200000")
	expected.Set("reverse", "false")
	expected.Set("count", "3")
	expected.Set("include_token", "true")
	expected.Set("include_custom_message_type", "true")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "include_meta"}, []string{})
}

func HistoryRequestWithMetaCommon(t *testing.T, withMeta bool) {
	assert := assert.New(t)

	opts := newHistoryOpts(pubnub, pubnub.ctx)
	opts.Channel = "ch"
	opts.Start = int64(100000)
	opts.End = int64(200000)
	opts.setStart = true
	opts.setEnd = true
	opts.Reverse = false
	opts.Count = 3
	opts.IncludeTimetoken = true
	opts.WithMeta = withMeta

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/history/sub-key/sub_key/channel/%s", opts.Channel),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("start", "100000")
	expected.Set("end", "200000")
	expected.Set("reverse", "false")
	expected.Set("count", "3")
	expected.Set("include_token", "true")
	expected.Set("include_meta", strconv.FormatBool(withMeta))
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestHistoryRequestWithMetaTrue(t *testing.T) {
	HistoryRequestWithMetaCommon(t, true)
}

func TestHistoryRequestWithMetaFalse(t *testing.T) {
	HistoryRequestWithMetaCommon(t, false)
}

func TestHistoryResponseParsingStringMessages(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[["hey-1","hey-two","hey-1","hey-1","hey-1","hey0","hey1","hey2","hey3","hey4","hey5","hey6","hey7","hey8","hey9","hey10","hey0","hey1","hey2","hey3","hey4","hey5","hey6","hey7","hey8","hey9","hey10","hey0","hey1","hey2","hey3","hey4","hey5","hey6","hey7","hey8","hey9","hey10"],14991775432719844,14991868111600528]`)

	pn := NewPubNubDemo()

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal(messages[0].Message, "hey-1")
	assert.Equal(messages[1].Message, "hey-two")
}

func TestHistoryResponseParsingStringWithTimetoken(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"timetoken":15232761410327866,"message":"hey-1"},{"timetoken":15232761410327866,"message":"hey-2"},{"timetoken":15232761410327866,"message":"hey-3"}],15232761410327866,15232761410327866]`)
	pn := NewPubNubDemo()

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(15232761410327866), resp.StartTimetoken)
	assert.Equal(int64(15232761410327866), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal("hey-1", messages[0].Message)
	assert.Equal(int64(15232761410327866), messages[0].Timetoken)
	assert.Equal("hey-2", messages[1].Message)
	assert.Equal(int64(15232761410327866), messages[1].Timetoken)
	assert.Equal("hey-3", messages[2].Message)
	assert.Equal(int64(15232761410327866), messages[2].Timetoken)
}

func TestHistoryResponseParsingInt(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[1,2,3,4,5,6,7],14991775432719844,14991868111600528]`)
	pn := NewPubNubDemo()

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal(float64(1), messages[0].Message)
}

func TestHistoryResponseParsingInt1(t *testing.T) {
	assert := assert.New(t)
	int1 := int(1)
	jsonString := []byte(fmt.Sprintf("[[%d],14991775432719844,14991868111600528]", int1))

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal(float64(int1), messages[0].Message)
}

func TestHistoryResponseParsingSlice(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[[1,2,3],["one","two","three"]],14991775432719844,14991868111600528]`)
	pn := NewPubNubDemo()

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
}

func TestHistoryResponseParsingMap(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"one":1,"two":2},{"three":3,"four":4}],14991775432719844,14991868111600528]`)
	pn := NewPubNubDemo()

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal(map[string]interface{}{"two": float64(2), "one": float64(1)},
		messages[0].Message)
}

func TestHistoryPNOther(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNubDemo()
	pn.Config.CipherKey = "testCipher"
	pn.Config.UseRandomInitializationVector = false

	int64Val := int64(14991775432719844)
	jsonString := []byte(`[[{"pn_other":"ven1bo79fk88nq5EIcnw/N9RmGzLeeWMnsabr1UL3iw="},1,"a",1.1,false,14991775432719844],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	data := messages[0].Message
	switch v := data.(type) {
	case map[string]interface{}:
		testMap := make(map[string]interface{})
		testMap = v
		if msgOther2, ok := testMap["pn_other"].(map[string]interface{}); !ok {
			assert.Fail("!map[string]interface{} 2")
		} else {
			assert.Equal(float64(3), msgOther2["three"])
			assert.Equal(float64(4), msgOther2["four"])
		}

		break
	default:
		assert.Fail(fmt.Sprintf("%s", reflect.TypeOf(v).Kind()), " expected interafce")
		break
	}
	assert.Equal(float64(1), messages[1].Message)
	assert.Equal("a", messages[2].Message)
	assert.Equal(float64(1.1), messages[3].Message)
	assert.Equal(false, messages[4].Message)
	assert.Equal(float64(int64Val), messages[5].Message)
	pnconfig.CipherKey = ""
}

func TestHistoryPNOtherYay(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNubDemo()
	pn.Config.UseRandomInitializationVector = false
	pn.Config.CipherKey = "enigma"
	int64Val := int64(14991775432719844)
	jsonString := []byte(`[[{"pn_other":"Wi24KS4pcTzvyuGOHubiXg=="},1,"a",1.1,false,14991775432719844],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	data := messages[0].Message
	switch v := data.(type) {
	case map[string]interface{}:
		testMap := make(map[string]interface{})
		testMap = v
		pnOther := testMap["pn_other"].(string)
		assert.Equal("yay!", pnOther)

		break
	default:
		assert.Fail("failed")

		break
	}
	assert.Equal(float64(1), messages[1].Message)
	assert.Equal("a", messages[2].Message)
	assert.Equal(float64(1.1), messages[3].Message)
	assert.Equal(false, messages[4].Message)
	assert.Equal(float64(int64Val), messages[5].Message)
	pnconfig.CipherKey = ""
}

func TestHistoryResponseParsingSliceInMapWithTimetoken(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"message":[1,2,3,["one","two","three"]],"timetoken":1111}],14991775432719844,14991868111600528]`)

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal([]interface{}{float64(1), float64(2), float64(3),
		[]interface{}{"one", "two", "three"}}, messages[0].Message)
}

func TestHistoryResponseParsingMapInSlice(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[[{"one":"two","three":[5,6]}]],14991775432719844,14991868111600528]`)

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal([]interface{}{
		map[string]interface{}{
			"one": "two", "three": []interface{}{float64(5), float64(6)}}},
		messages[0].Message)
}

func TestHistoryEncrypt(t *testing.T) {
	assert := assert.New(t)
	pnconfig := NewDemoConfig()
	pnconfig.CipherKey = "testCipher"
	pnconfig.UseRandomInitializationVector = false
	pubnub := NewPubNub(pnconfig)

	jsonString := []byte(`[["MnwzPGdVgz2osQCIQJviGg=="],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pubnub.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	assert.Equal("hey", messages[0].Message)
	assert.Nil(messages[0].Error)
	pnconfig.CipherKey = ""
}

func TestHistoryEncryptSlice(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNubDemo()
	pn.Config.CipherKey = "testCipher"
	pn.Config.UseRandomInitializationVector = false

	jsonString := []byte(`[["gwkdY8qcv60GM/PslArWQsdXrQ6LwJD2HoaEfy0CjMc="],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	if resp, ok := messages[0].Message.([]interface{}); !ok {
		assert.Fail("!ok []interface{}")
	} else {
		assert.Equal("hey-1", resp[0].(string))
		assert.Equal("hey-2", resp[1].(string))
		assert.Equal("hey-3", resp[2].(string))
	}

	pnconfig.CipherKey = ""
}

func TestHistoryEncryptMap(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNubDemo()
	pn.Config.CipherKey = "testCipher"
	pn.Config.UseRandomInitializationVector = false

	jsonString := []byte(`[["wIC13nvJcI4vBtWNFVUu0YDiqREr9kavB88xeyWTweDS363Yl84RCWqOHWTol4aY"],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	if resp, ok := messages[0].Message.(map[string]interface{}); !ok {
		assert.Fail("!ok map[string]interface {}")
	} else {
		assert.Equal(float64(1), resp["one"].(float64))
		if resp2, ok2 := resp["two"].([]interface{}); !ok2 {
			assert.Fail("!ok2 map[int]interface{}")
		} else {
			assert.Equal("hey-1", resp2[0])
			assert.Equal("hey-2", resp2[1])
		}
	}

	pnconfig.CipherKey = ""
}

func TestHistoryResponseMeta(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"message":"my-message","meta":{"m1":"n1","m2":"n2"}}],15699986472636251,15699986472636251]`)

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(15699986472636251), resp.StartTimetoken)
	assert.Equal(int64(15699986472636251), resp.EndTimetoken)

	messages := resp.Messages
	meta := messages[0].Meta.(map[string]interface{})
	assert.Equal("my-message", messages[0].Message)
	assert.Equal("n1", meta["m1"])
	assert.Equal("n2", meta["m2"])
}

func TestHistoryResponseMetaAndTT(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"message":"my-message","meta":{"m1":"n1","m2":"n2"},"timetoken":15699986472636251}],15699986472636251,15699986472636251]`)

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(15699986472636251), resp.StartTimetoken)
	assert.Equal(int64(15699986472636251), resp.EndTimetoken)

	messages := resp.Messages
	meta := messages[0].Meta.(map[string]interface{})
	assert.Equal("my-message", messages[0].Message)
	assert.Equal(int64(15699986472636251), messages[0].Timetoken)
	assert.Equal("n1", meta["m1"])
	assert.Equal("n2", meta["m2"])
}

func TestHistoryResponseError(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`s`)

	pn := NewPubNubDemo()
	_, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestHistoryResponseStartTTError(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"message":[1,2,3,["one","two","three"]],"timetoken":1111}],"s","a"]`)

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Equal(int64(0), resp.StartTimetoken)
	assert.Equal(int64(0), resp.EndTimetoken)
	assert.Nil(err)

}

func TestHistoryResponseEndTTError(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"message":[1,2,3,["one","two","three"]],"timetoken":1111}],121324,"a"]`)

	pn := NewPubNubDemo()
	resp, _, err := newHistoryResponse(jsonString, pn.initHistoryOpts(), fakeResponseState)
	assert.Equal(int64(121324), resp.StartTimetoken)
	assert.Equal(int64(0), resp.EndTimetoken)
	assert.Nil(err)
}

func TestHistoryCryptoModuleWithEncryptedMessage(t *testing.T) {
	assert := assert.New(t)
	pnconfig := NewDemoConfig()
	pubnub := NewPubNub(pnconfig)
	crypto, init_err := crypto.NewAesCbcCryptoModule("enigma", true)

	assert.Nil(init_err)

	pubnub.Config.CryptoModule = crypto

	// Rust generated cipher text
	jsonString := []byte(`[["UE5FRAFBQ1JIEALf+E65kseYJwTw2J6BUk9MePHiCcBCS+8ykXLkBIOA"],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pubnub.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	assert.Equal("test", messages[0].Message)
	assert.Nil(messages[0].Error)
	pnconfig.CipherKey = ""
}

func TestHistoryCryptoModuleWithNoEncryptedMessage(t *testing.T) {
	assert := assert.New(t)
	pnconfig := NewDemoConfig()
	pubnub := NewPubNub(pnconfig)
	crypto, init_err := crypto.NewAesCbcCryptoModule("enigma", true)

	assert.Nil(init_err)

	pubnub.Config.CryptoModule = crypto

	jsonString := []byte(`[["test"],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, pubnub.initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	assert.Equal("test", messages[0].Message)
	assert.NotNil(messages[0].Error)
	pnconfig.CipherKey = ""
}

// Enhanced History Tests for better coverage
func TestHistoryCountBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Test minimum count
	o := newHistoryBuilder(pn)
	o.Channel("test-channel")
	o.Count(1)

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1", query.Get("count"))

	// Test standard maximum count (implementation caps at 100)
	o.Count(100)
	query, err = o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("count"))

	// Test count over maximum - implementation normalizes to 100
	o.Count(1000)
	query, err = o.opts.buildQuery()
	assert.Nil(err)
	// The implementation caps this at 100
	assert.Equal("100", query.Get("count"))

	// Test zero count - implementation defaults to 100
	o.Count(0)
	query, err = o.opts.buildQuery()
	assert.Nil(err)
	// The implementation defaults zero to 100
	assert.Equal("100", query.Get("count"))
}

func TestHistoryTimetokenBoundaries(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newHistoryBuilder(pn)
	o.Channel("test-channel")

	// Test with minimum valid timetoken
	o.Start(1)
	o.End(2)

	query, err := o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1", query.Get("start"))
	assert.Equal("2", query.Get("end"))

	// Test with maximum reasonable timetoken
	maxTimetoken := int64(9223372036854775807) // Max int64
	o.Start(maxTimetoken - 1000)
	o.End(maxTimetoken)

	query, err = o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("9223372036854774807", query.Get("start"))
	assert.Equal("9223372036854775807", query.Get("end"))

	// Test with zero timetokens
	o.Start(0)
	o.End(0)

	query, err = o.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("0", query.Get("start"))
	assert.Equal("0", query.Get("end"))
}

func TestHistoryErrorResponseParsing(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := pn.initHistoryOpts()

	// Test completely invalid JSON
	invalidJSON := []byte(`not json at all`)
	_, _, err := newHistoryResponse(invalidJSON, opts, fakeResponseState)
	assert.NotNil(err)
	assert.Contains(err.Error(), "parsing")

	// Test malformed JSON structure (expects array, not object)
	malformedJSON := []byte(`{"invalid": "structure"}`)
	_, _, err = newHistoryResponse(malformedJSON, opts, fakeResponseState)
	assert.NotNil(err)
	assert.Contains(err.Error(), "parsing")

	// Test empty object (should error - History expects array format)
	emptyJSON := []byte(`{}`)
	_, _, err = newHistoryResponse(emptyJSON, opts, fakeResponseState)
	assert.NotNil(err)
	assert.Contains(err.Error(), "parsing")

	// Test empty array (should error - History expects at least 3 elements)
	emptyArrayJSON := []byte(`[]`)
	_, _, err = newHistoryResponse(emptyArrayJSON, opts, fakeResponseState)
	assert.NotNil(err)
	assert.Contains(err.Error(), "parsing")

	// Test response with minimum valid structure (empty messages array)
	validMinimalJSON := []byte(`[[], "start", "end"]`)
	r, _, err := newHistoryResponse(validMinimalJSON, opts, fakeResponseState)
	assert.Nil(err)
	assert.NotNil(r)
	assert.NotNil(r.Messages)
	assert.Equal(0, len(r.Messages))
}

func TestHistorySpecialChannelNames(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannels := []string{
		"ch-with-dash",
		"ch_with_underscore",
		"ch.with.dot",
		"ch:with:colon",
		"ch with space",
		"unicode-チャンネル",
	}

	for _, channel := range specialChannels {
		o := newHistoryBuilder(pn)
		o.Channel(channel)

		path, err := o.opts.buildPath()
		assert.Nil(err)
		assert.NotEmpty(path)
		assert.Contains(path, "history")
	}
}

// Validation Tests

func TestHistoryValidateMissingSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Equal("pubnub/validation: pubnub: History: Missing Subscribe Key", opts.validate().Error())
}

func TestHistoryValidateMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = ""

	assert.Equal("pubnub/validation: pubnub: History: Missing Channel", opts.validate().Error())
}

func TestHistoryValidateSuccess(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	assert.Nil(opts.validate())
}

// Systematic Builder Pattern Tests

func TestHistoryBuilderBasic(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHistoryBuilder(pn)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn, builder.opts.pubnub)
}

func TestHistoryBuilderContext(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHistoryBuilderWithContext(pn, pn.ctx)
	assert.NotNil(builder)
	assert.NotNil(builder.opts)
	assert.Equal(pn.ctx, builder.opts.ctx)
}

func TestHistoryBuilderMethodChaining(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{"key": "value"}

	builder := newHistoryBuilder(pn)
	result := builder.Channel("test-channel").
		Start(1000).
		End(2000).
		Count(50).
		Reverse(true).
		IncludeTimetoken(true).
		IncludeMeta(true).
		QueryParam(queryParam)

	// Should return same instance for method chaining
	assert.Equal(builder, result)

	// Verify all values are set correctly
	assert.Equal("test-channel", builder.opts.Channel)
	assert.Equal(int64(1000), builder.opts.Start)
	assert.Equal(int64(2000), builder.opts.End)
	assert.Equal(50, builder.opts.Count)
	assert.Equal(true, builder.opts.Reverse)
	assert.Equal(true, builder.opts.IncludeTimetoken)
	assert.Equal(true, builder.opts.WithMeta)
	assert.Equal(queryParam, builder.opts.QueryParam)
	assert.Equal(true, builder.opts.setStart)
	assert.Equal(true, builder.opts.setEnd)
}

func TestHistoryBuilderSetters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHistoryBuilder(pn)

	// Test Channel setter
	builder.Channel("my-channel")
	assert.Equal("my-channel", builder.opts.Channel)

	// Test Start setter
	builder.Start(12345)
	assert.Equal(int64(12345), builder.opts.Start)
	assert.Equal(true, builder.opts.setStart)

	// Test End setter
	builder.End(67890)
	assert.Equal(int64(67890), builder.opts.End)
	assert.Equal(true, builder.opts.setEnd)

	// Test Count setter
	builder.Count(25)
	assert.Equal(25, builder.opts.Count)

	// Test Reverse setter
	builder.Reverse(true)
	assert.Equal(true, builder.opts.Reverse)

	// Test IncludeTimetoken setter
	builder.IncludeTimetoken(true)
	assert.Equal(true, builder.opts.IncludeTimetoken)

	// Test IncludeMeta setter
	builder.IncludeMeta(true)
	assert.Equal(true, builder.opts.WithMeta)

	// Test QueryParam setter
	queryParam := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	builder.QueryParam(queryParam)
	assert.Equal(queryParam, builder.opts.QueryParam)
}

func TestHistoryBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHistoryBuilder(pn)
	builder.Channel("test-channel")

	queryParam := map[string]string{
		"custom": "param",
		"test":   "value",
	}
	builder.QueryParam(queryParam)

	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("param", query.Get("custom"))
	assert.Equal("value", query.Get("test"))
}

// URL/Path Building Tests

func TestHistoryBuildPath(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = "test-channel"

	path, err := opts.buildPath()
	assert.Nil(err)
	expected := "/v2/history/sub-key/demo/channel/test-channel"
	assert.Equal(expected, path)
}

func TestHistoryBuildPathWithSpecialChars(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = "channel-with-special@chars#and$symbols"

	path, err := opts.buildPath()
	assert.Nil(err)
	// Should URL encode special characters
	assert.Contains(path, "channel-with-special%40chars%23and%24symbols")
}

func TestHistoryBuildQueryConditionalParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)

	// Test without setStart/setEnd flags
	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("", query.Get("start"))
	assert.Equal("", query.Get("end"))
	assert.Equal("100", query.Get("count"))           // Default count
	assert.Equal("false", query.Get("reverse"))       // Default reverse
	assert.Equal("false", query.Get("include_token")) // Default
	assert.Equal("false", query.Get("include_meta"))  // Default

	// Test with setStart/setEnd flags
	opts.Start = 1000
	opts.End = 2000
	opts.setStart = true
	opts.setEnd = true
	opts.Count = 50
	opts.Reverse = true
	opts.IncludeTimetoken = true
	opts.WithMeta = true

	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1000", query.Get("start"))
	assert.Equal("2000", query.Get("end"))
	assert.Equal("50", query.Get("count"))
	assert.Equal("true", query.Get("reverse"))
	assert.Equal("true", query.Get("include_token"))
	assert.Equal("true", query.Get("include_meta"))
}

func TestHistoryBuildQueryCountValidation(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)

	// Test count = 0 (should default to 100)
	opts.Count = 0
	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("count"))

	// Test count within valid range
	opts.Count = 50
	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("50", query.Get("count"))

	// Test count at maximum (100)
	opts.Count = 100
	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("count"))

	// Test count over maximum (should default to 100)
	opts.Count = 150
	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("count"))

	// Test negative count (should default to 100)
	opts.Count = -5
	query, err = opts.buildQuery()
	assert.Nil(err)
	assert.Equal("100", query.Get("count"))
}

func TestHistoryBuildQueryWithComplexParams(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)

	complexParams := map[string]string{
		"filter":         "status=active",
		"sort":           "time,desc",
		"include":        "metadata,custom",
		"special_chars":  "value@with#symbols",
		"unicode":        "测试参数",
		"empty_value":    "",
		"number_string":  "42",
		"boolean_string": "true",
	}
	opts.QueryParam = complexParams

	query, err := opts.buildQuery()
	assert.Nil(err)

	// Verify all parameters are present
	for key, expectedValue := range complexParams {
		actualValue := query.Get(key)
		if key == "special_chars" {
			// Special characters should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should be URL encoded", key)
		} else if key == "unicode" {
			// Unicode should be URL encoded
			assert.Contains(actualValue, "%", "Query parameter %s should contain URL encoded Unicode", key)
		} else if key == "filter" {
			// Filter parameter contains = which gets URL encoded
			assert.Equal("status%3Dactive", actualValue, "Query parameter %s should be URL encoded", key)
		} else if key == "sort" {
			// Sort parameter contains , which gets URL encoded
			assert.Equal("time%2Cdesc", actualValue, "Query parameter %s should be URL encoded", key)
		} else if key == "include" {
			// Include parameter contains , which gets URL encoded
			assert.Equal("metadata%2Ccustom", actualValue, "Query parameter %s should be URL encoded", key)
		} else {
			assert.Equal(expectedValue, actualValue, "Query parameter %s should match", key)
		}
	}
}

// HTTP Method and Operation Tests

func TestHistoryOperationType(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)

	assert.Equal(PNHistoryOperation, opts.operationType())
}

func TestHistoryIsAuthRequired(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)

	assert.True(opts.isAuthRequired())
}

func TestHistoryTimeouts(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)

	assert.Equal(pn.Config.NonSubscribeRequestTimeout, opts.requestTimeout())
	assert.Equal(pn.Config.ConnectTimeout, opts.connectTimeout())
}

// Comprehensive Edge Case Tests

func TestHistoryWithUnicodeChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = "测试频道-русский-チャンネル"

	// Should pass validation
	assert.Nil(opts.validate())

	// Should build path with URL encoding
	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/history/")
	// Unicode should be URL encoded
	assert.Contains(path, "%")
}

func TestHistoryWithExtremeParameters(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHistoryBuilder(pn)
	builder.Channel("test-channel")

	// Test extreme timetoken values
	maxTimetoken := int64(9223372036854775807) // Max int64
	minTimetoken := int64(1)

	builder.Start(minTimetoken)
	builder.End(maxTimetoken)
	builder.Count(1) // Minimum count
	builder.Reverse(true)
	builder.IncludeTimetoken(true)
	builder.IncludeMeta(true)

	assert.Equal(minTimetoken, builder.opts.Start)
	assert.Equal(maxTimetoken, builder.opts.End)
	assert.Equal(1, builder.opts.Count)
	assert.Equal(true, builder.opts.Reverse)
	assert.Equal(true, builder.opts.IncludeTimetoken)
	assert.Equal(true, builder.opts.WithMeta)

	// Test query building with extreme values
	query, err := builder.opts.buildQuery()
	assert.Nil(err)
	assert.Equal("1", query.Get("start"))
	assert.Equal("9223372036854775807", query.Get("end"))
	assert.Equal("1", query.Get("count"))
	assert.Equal("true", query.Get("reverse"))
	assert.Equal("true", query.Get("include_token"))
	assert.Equal("true", query.Get("include_meta"))
}

func TestHistoryWithVeryLongChannelName(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a very long channel name
	longName := ""
	for i := 0; i < 100; i++ {
		longName += fmt.Sprintf("channel_%d_", i)
	}

	opts := newHistoryOpts(pn, pn.ctx)
	opts.Channel = longName

	assert.Nil(opts.validate())

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Contains(path, "/history/")
	assert.Contains(path, "channel_0_")
	assert.Contains(path, "channel_99_")
}

func TestHistoryWithEmptyQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.QueryParam = map[string]string{}

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestHistoryWithNilQueryParam(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newHistoryOpts(pn, pn.ctx)
	opts.QueryParam = nil

	query, err := opts.buildQuery()
	assert.Nil(err)
	assert.NotNil(query)
}

func TestHistoryParameterCombinations(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	testCases := []struct {
		name            string
		count           int
		reverse         bool
		includeToken    bool
		includeMeta     bool
		expectedCount   string
		expectedReverse string
		expectedToken   string
		expectedMeta    string
	}{
		{
			name:            "All false",
			count:           25,
			reverse:         false,
			includeToken:    false,
			includeMeta:     false,
			expectedCount:   "25",
			expectedReverse: "false",
			expectedToken:   "false",
			expectedMeta:    "false",
		},
		{
			name:            "All true",
			count:           75,
			reverse:         true,
			includeToken:    true,
			includeMeta:     true,
			expectedCount:   "75",
			expectedReverse: "true",
			expectedToken:   "true",
			expectedMeta:    "true",
		},
		{
			name:            "Mixed",
			count:           10,
			reverse:         true,
			includeToken:    false,
			includeMeta:     true,
			expectedCount:   "10",
			expectedReverse: "true",
			expectedToken:   "false",
			expectedMeta:    "true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newHistoryBuilder(pn)
			builder.Channel("test-channel")
			builder.Count(tc.count)
			builder.Reverse(tc.reverse)
			builder.IncludeTimetoken(tc.includeToken)
			builder.IncludeMeta(tc.includeMeta)

			query, err := builder.opts.buildQuery()
			assert.Nil(err)
			assert.Equal(tc.expectedCount, query.Get("count"))
			assert.Equal(tc.expectedReverse, query.Get("reverse"))
			assert.Equal(tc.expectedToken, query.Get("include_token"))
			assert.Equal(tc.expectedMeta, query.Get("include_meta"))
		})
	}
}

// Error Scenario Tests

func TestHistoryExecuteError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = "" // Invalid config

	builder := newHistoryBuilder(pn)
	builder.Channel("test-channel")

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Subscribe Key")
}

func TestHistoryExecuteErrorMissingChannel(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	builder := newHistoryBuilder(pn)
	// Don't set Channel, should fail validation

	_, _, err := builder.Execute()
	assert.NotNil(err)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestHistoryPathBuildingEdgeCases(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	specialChannels := []string{
		"channel@with%encoded",
		"channel/with/slashes",
		"channel?with=query&chars",
		"channel#with#hashes",
		"channel with spaces and símböls",
		"测试频道-русский-チャンネル-한국어",
	}

	for _, channel := range specialChannels {
		opts := newHistoryOpts(pn, pn.ctx)
		opts.Channel = channel

		// Should pass validation
		assert.Nil(opts.validate(), "Should validate channel: %s", channel)

		// Should build valid path
		path, err := opts.buildPath()
		assert.Nil(err, "Should build path for channel: %s", channel)
		assert.Contains(path, "/history/", "Should contain history path for: %s", channel)
	}
}

func TestHistoryTransportSetter(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	// Create a mock transport
	transport := &http.Transport{}

	builder := newHistoryBuilder(pn)
	result := builder.Transport(transport)

	// Should return same instance for chaining
	assert.Equal(builder, result)

	// Should set the transport
	assert.Equal(transport, builder.opts.Transport)
}
