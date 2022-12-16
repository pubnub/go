package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertRevokeToken(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRevokeTokenBuilder(pn)
	if testContext {
		o = newRevokeTokenBuilderWithContext(pn, backgroundContext)
	}

	token := "token"
	o.QueryParam(queryParam)
	o.Token(token)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(revokeTokenPath, pn.Config.SubscribeKey, token),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestRevokeToken(t *testing.T) {
	AssertRevokeToken(t, true, false)
}

func TestRevokeTokenContext(t *testing.T) {
	AssertRevokeToken(t, true, true)
}

func TestRevokeTokenResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)
	jsonBytes := []byte(`s`)

	_, _, err := newPNRevokeTokenResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRevokeTokenResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	opts := newRevokeTokenOpts(pn, pn.ctx)
	jsonBytes := []byte(`{"status":200}`)

	_, s, err := newPNRevokeTokenResponse(jsonBytes, opts, StatusResponse{StatusCode: 200})
	assert.Equal(200, s.StatusCode)

	assert.Nil(err)
}
