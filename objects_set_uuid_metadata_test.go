package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/pubnub/go/utils"
	"github.com/stretchr/testify/assert"
)

func AssertSetUUIDMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	incl := []PNUUIDMetadataInclude{
		PNUUIDMetadataIncludeCustom,
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

	o := newSetUUIDMetadataBuilder(pn)
	if testContext {
		o = newSetUUIDMetadataBuilderWithContext(pn, backgroundContext)
	}

	o.Include(incl)
	o.UUID("id0")
	o.Name("name")
	o.ExternalID("exturl")
	o.ProfileURL("prourl")
	o.Email("email")
	o.Custom(custom)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/uuids/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\",\"externalId\":\"exturl\",\"profileUrl\":\"prourl\",\"email\":\"email\",\"custom\":{\"a\":\"b\",\"c\":\"d\"}}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestSetUUIDMetadata(t *testing.T) {
	AssertSetUUIDMetadata(t, true, false)
}

func TestSetUUIDMetadataContext(t *testing.T) {
	AssertSetUUIDMetadata(t, true, true)
}

func TestSetUUIDMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &setUUIDMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newPNSetUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetUUIDMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := &setUUIDMetadataOpts{
		pubnub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:19.140324Z","updated":"2019-08-20T13:26:19.140324Z","eTag":"AbyT4v2p6K7fpQE"}}`)

	r, _, err := newPNSetUUIDMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("extid", r.Data.ExternalID)
	assert.Equal("purl", r.Data.ProfileURL)
	assert.Equal("email", r.Data.Email)
	// assert.Equal("2019-08-20T13:26:19.140324Z", r.Data.Created)
	assert.Equal("2019-08-20T13:26:19.140324Z", r.Data.Updated)
	assert.Equal("AbyT4v2p6K7fpQE", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])
	assert.Equal("d", r.Data.Custom["c"])

	assert.Nil(err)
}
