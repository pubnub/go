package pubnub

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

var (
	fakeResponseState = StatusResponse{}
)

func initHistoryOpts() *historyOpts {
	return &historyOpts{
		Channel:          "ch",
		Start:            int64(100000),
		End:              int64(200000),
		SetStart:         true,
		SetEnd:           true,
		Reverse:          false,
		Count:            3,
		IncludeTimetoken: true,
		pubnub:           pubnub,
	}
}

func TestHistoryRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := &historyOpts{
		Channel: "ch",
		pubnub:  pubnub,
	}

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
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHistoryRequestAllParams(t *testing.T) {
	assert := assert.New(t)

	opts := &historyOpts{
		Channel:          "ch",
		Start:            int64(100000),
		End:              int64(200000),
		SetStart:         true,
		SetEnd:           true,
		Reverse:          false,
		Count:            3,
		IncludeTimetoken: true,
		pubnub:           pubnub,
	}

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
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestHistoryResponseParsingStringMessages(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[["hey-1","hey-two","hey-1","hey-1","hey-1","hey0","hey1","hey2","hey3","hey4","hey5","hey6","hey7","hey8","hey9","hey10","hey0","hey1","hey2","hey3","hey4","hey5","hey6","hey7","hey8","hey9","hey10","hey0","hey1","hey2","hey3","hey4","hey5","hey6","hey7","hey8","hey9","hey10"],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
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

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
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

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
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

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal(float64(int1), messages[0].Message)
}

func TestHistoryResponseParsingSlice(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[[1,2,3],["one","two","three"]],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	// assert.Equal(int64(14991868111600528), resp.EndTimetoken)
	//
	// messages := resp.Messages
	// assert.Equal([]interface{}{float64(1), float64(2), float64(3)},
	// messages[0].Message)
}

func TestHistoryResponseParsingMap(t *testing.T) {
	assert := assert.New(t)
	pnconfig.CipherKey = ""
	jsonString := []byte(`[[{"one":1,"two":2},{"three":3,"four":4}],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	/*for v, k := range messages {
		fmt.Println(v, k)
	}*/
	//fmt.Println(messages[0].Message)
	//fmt.Println(messages[1].Message)
	assert.Equal(map[string]interface{}{"two": float64(2), "one": float64(1)},
		messages[0].Message)
}

func TestHistoryPNOther(t *testing.T) {
	assert := assert.New(t)
	pnconfig.CipherKey = "testCipher"
	//fmt.Println(utils.EncryptString("testCipher", `{"three":3, "four":4}`))
	int64Val := int64(14991775432719844)
	jsonString := []byte(`[[{"pn_other":"ven1bo79fk88nq5EIcnw/N9RmGzLeeWMnsabr1UL3iw="},1,"a",1.1,false,14991775432719844],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	data := messages[0].Message
	switch v := data.(type) {
	case map[string]interface{}:
		testMap := make(map[string]interface{})
		testMap = v
		//pn_other := make(map[string]string)
		pn_other := testMap["pn_other"].(string)
		assert.Equal(pn_other, "{\"three\":3, \"four\":4}")
		//assert.Equal(pn_other["three"], 3)

		break
	default:
		//fmt.Println(reflect.TypeOf(v).Kind(), reflect.TypeOf(data).Kind(), v, data)
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

func TestHistoryResponseParsingSliceInMapWithTimetoken(t *testing.T) {
	assert := assert.New(t)

	jsonString := []byte(`[[{"message":[1,2,3,["one","two","three"]],"timetoken":1111}],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
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

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
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
	pnconfig.CipherKey = "testCipher"
	pubnub = NewPubNub(pnconfig)

	jsonString := []byte(`[["GUI1NhVPOxZap54NuLEaow=="],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)
	fmt.Println(resp)

	messages := resp.Messages
	assert.Equal("hey", messages[0].Message)
	pnconfig.CipherKey = ""
}

func TestHistoryEncryptSlice(t *testing.T) {
	assert := assert.New(t)
	pnconfig.CipherKey = "testCipher"

	jsonString := []byte(`[["gwkdY8qcv60GM/PslArWQsdXrQ6LwJD2HoaEfy0CjMc="],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages

	assert.Equal("[\"hey-1\",\"hey-2\",\"hey-3\"]", messages[0].Message)
	pnconfig.CipherKey = ""
}

func TestHistoryEncryptMap(t *testing.T) {
	assert := assert.New(t)
	pnconfig.CipherKey = "testCipher"

	jsonString := []byte(`[["wIC13nvJcI4vBtWNFVUu0YDiqREr9kavB88xeyWTweDS363Yl84RCWqOHWTol4aY"],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	//log.Println(messages[0].Message)
	assert.Equal(`{"one":1,"two":["hey-1","hey-2"]}`, messages[0].Message)

	pnconfig.CipherKey = ""
}
