package pubnub

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// recordingRoundTripper captures the outgoing request method and returns a
// canned response, allowing executeRequest to be exercised without a network.
type recordingRoundTripper struct {
	lastMethod string
	statusCode int
	body       string
}

func (rt *recordingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.lastMethod = req.Method
	return &http.Response{
		StatusCode: rt.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(rt.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// TestExecuteRequestPUTMethod verifies the PUT branch added to executeRequest
// actually issues an HTTP PUT (UpdateEntity uses PUT for full replacement).
func TestExecuteRequestPUTMethod(t *testing.T) {
	assert := assert.New(t)

	rt := &recordingRoundTripper{statusCode: 200, body: `{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":2,"eTag":"Cxyz1..."}}`}

	pn := NewPubNub(NewDemoConfig())
	pn.SetClient(&http.Client{Transport: rt})

	res, status, err := pn.DataSync.UpdateEntity().
		ID("entity-abc").
		EntityClassVersion(2).
		Status("active").
		Execute()

	assert.Nil(err)
	assert.Equal("PUT", rt.lastMethod)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
	assert.Equal("Cxyz1...", res.Data.ETag)
}

// TestExecuteRequestAccepts201 verifies that a 201 Created response (returned by
// CreateEntity) is treated as success and not as a server error.
func TestExecuteRequestAccepts201(t *testing.T) {
	assert := assert.New(t)

	rt := &recordingRoundTripper{statusCode: 201, body: `{"status":201,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"eTag":"BfklQ..."}}`}

	pn := NewPubNub(NewDemoConfig())
	pn.SetClient(&http.Client{Transport: rt})

	res, status, err := pn.DataSync.CreateEntity().
		EntityClass("vehicle").
		EntityClassVersion(1).
		Execute()

	assert.Nil(err)
	assert.Equal("POST", rt.lastMethod)
	assert.Equal(201, status.StatusCode)
	assert.NotNil(res)
	assert.Equal("entity-abc", res.Data.ID)
	assert.Equal("BfklQ...", res.Data.ETag)
}

// TestExecuteRequestDeleteEmptyBody verifies DELETE succeeds with an empty 200 body.
func TestExecuteRequestDeleteEmptyBody(t *testing.T) {
	assert := assert.New(t)

	rt := &recordingRoundTripper{statusCode: 200, body: ""}

	pn := NewPubNub(NewDemoConfig())
	pn.SetClient(&http.Client{Transport: rt})

	res, status, err := pn.DataSync.RemoveEntity().ID("entity-abc").Execute()

	assert.Nil(err)
	assert.Equal("DELETE", rt.lastMethod)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}

// TestExecuteRequestCreateRelationship verifies CreateRelationship issues POST
// with the relationship media type and accepts a 201 Created response.
func TestExecuteRequestCreateRelationship(t *testing.T) {
	assert := assert.New(t)

	rt := &recordingRoundTripper{statusCode: 201, body: `{"status":201,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"eTag":"1"}}`}

	pn := NewPubNub(NewDemoConfig())
	pn.SetClient(&http.Client{Transport: rt})

	res, status, err := pn.DataSync.CreateRelationship().
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Execute()

	assert.Nil(err)
	assert.Equal("POST", rt.lastMethod)
	assert.Equal(201, status.StatusCode)
	assert.NotNil(res)
	assert.Equal("r-123", res.Data.ID)
	assert.Equal("u123", res.Data.EntityAID)
}
