package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/pubnub/go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertGetChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	incl := []PNChannelMetadataInclude{
		PNChannelMetadataIncludeCustom,
	}
	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newGetChannelMetadataBuilder(pn)
	if testContext {
		o = newGetChannelMetadataBuilderWithContext(pn, backgroundContext)
	}

	o.Include(incl)
	o.Channel("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestGetChannelMetadata(t *testing.T) {
	AssertGetChannelMetadata(t, true, false)
}

func TestGetChannelMetadataContext(t *testing.T) {
	AssertGetChannelMetadata(t, true, true)
}

func TestGetChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getChannelMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNGetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestGetChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &getChannelMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","description":"desc","custom":{"a":"b"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T13:26:08.341297Z","eTag":"Aee9zsKNndXlHw"}}`)

	r, _, err := newPNGetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("desc", r.Data.Description)
	//assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Created)
	assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Updated)
	assert.Equal("Aee9zsKNndXlHw", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])

	assert.Nil(err)
}
