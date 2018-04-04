package pubnub

import (
	"fmt"
	"log"
	"net/url"
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

	jsonString := []byte(`[[{"timetoken":1111,"message":"hey-1"},{"timetoken":2222,"message":"hey-2"},{"timetoken":3333,"message":"hey-3"}],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal("hey-1", messages[0].Message)
	assert.Equal(int64(1111), messages[0].Timetoken)
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

	jsonString := []byte(`[[{"one":1,"two":2},{"three":3,"four":4}],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	assert.Equal(int64(14991775432719844), resp.StartTimetoken)
	assert.Equal(int64(14991868111600528), resp.EndTimetoken)

	messages := resp.Messages
	assert.Equal(map[string]interface{}{"three": float64(3), "four": float64(4)},
		messages[0].Message)
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
	pnconfig.CipherKey = "enigma"
	pubnub = NewPubNub(pnconfig)

	jsonString := []byte(`[["GUI1NhVPOxZap54NuLEaow=="],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	assert.Equal("hey", messages[0].Message)
}

func TestHistoryEncryptSlice(t *testing.T) {
	assert := assert.New(t)
	pnconfig.CipherKey = "testCipher"

	jsonString := []byte(`[["gwkdY8qcv60GM/PslArWQsdXrQ6LwJD2HoaEfy0CjMc="],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages

	assert.Equal("[\"hey-1\",\"hey-2\",\"hey-3\"]", messages[0].Message)
}

func TestHistoryEncryptMap(t *testing.T) {
	assert := assert.New(t)
	pnconfig.CipherKey = "testCipher"

	jsonString := []byte(`[["wIC13nvJcI4vBtWNFVUu0YDiqREr9kavB88xeyWTweDS363Yl84RCWqOHWTol4aY"],14991775432719844,14991868111600528]`)

	resp, _, err := newHistoryResponse(jsonString, initHistoryOpts(), fakeResponseState)
	assert.Nil(err)

	messages := resp.Messages
	log.Println(messages[0].Message)
	assert.Equal(`{"one":1,"two":["hey-1","hey-2"]}`, messages[0].Message)

	pnconfig.CipherKey = ""
}
