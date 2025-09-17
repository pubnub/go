package pubnub

import (
	"fmt"
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

	o := newHistoryBuilderWithContext(pubnub, backgroundContext)
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
