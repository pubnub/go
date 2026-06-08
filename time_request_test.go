package pubnub

import (
	"net/http"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v9/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// protoCaptureTransport records the negotiated protocol from the underlying RoundTripper
// (e.g. "HTTP/2.0" vs "HTTP/1.1") for integration checks.
type protoCaptureTransport struct {
	base http.RoundTripper
	dest *string
}

func (p *protoCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := p.base.RoundTrip(req)
	if resp != nil && p.dest != nil {
		*p.dest = resp.Proto
	}
	return resp, err
}

// TestTimeRequestUsesHTTP2OnH2PubNubAPIOrigin exercises the PubNub HTTP stack against the
// h2.pubnubapi.com origin, which terminates HTTP/2. Requires outbound HTTPS access.
func TestTimeRequestUsesHTTP2OnH2PubNubAPIOrigin(t *testing.T) {
	assert := assert.New(t)

	const h2Origin = "h2.pubnubapi.com"

	config := NewConfigWithUserId(UserId(GenerateUUID()))
	config.Origin = h2Origin
	config.UseHTTP2 = true

	pn := NewPubNub(config)

	base := pn.GetClient()
	var negotiatedProto string
	pn.SetClient(&http.Client{
		Transport: &protoCaptureTransport{base: base.Transport, dest: &negotiatedProto},
		Timeout:   base.Timeout,
	})

	_, s, err := pn.Time().Execute()

	assert.Nil(err)
	assert.Equal(200, s.StatusCode)
	assert.Equal(h2Origin, s.Origin)
	assert.Equal("HTTP/2.0", negotiatedProto,
		"h2.pubnubapi.com should negotiate HTTP/2; got %s (TLS / connectivity)", negotiatedProto)
}

// TestTimeRequestUsesHTTP11WhenUseHTTP2FalseOnH2PubNubAPIOrigin checks that UseHTTP2=false selects
// the HTTP/1-only client so TLS negotiates HTTP/1.1 even against the h2-named origin (which must
// still offer http/1.1 ALPN). Requires outbound HTTPS access.
func TestTimeRequestUsesHTTP11WhenUseHTTP2FalseOnH2PubNubAPIOrigin(t *testing.T) {
	assert := assert.New(t)

	const h2Origin = "h2.pubnubapi.com"

	config := NewConfigWithUserId(UserId(GenerateUUID()))
	config.Origin = h2Origin
	config.UseHTTP2 = false

	pn := NewPubNub(config)

	base := pn.GetClient()
	var negotiatedProto string
	pn.SetClient(&http.Client{
		Transport: &protoCaptureTransport{base: base.Transport, dest: &negotiatedProto},
		Timeout:   base.Timeout,
	})

	_, s, err := pn.Time().Execute()

	assert.Nil(err)
	assert.Equal(200, s.StatusCode)
	assert.Equal(h2Origin, s.Origin)
	assert.Equal("HTTP/1.1", negotiatedProto,
		"with UseHTTP2=false expect HTTP/1.1 ALPN against h2.pubnubapi.com; got %s", negotiatedProto)
}

func TestTimeRequestHTTP2(t *testing.T) {
	assert := assert.New(t)

	config := NewConfigWithUserId(UserId(GenerateUUID()))
	config.Origin = "ssp.pubnub.com"
	config.UseHTTP2 = true

	pn := NewPubNub(config)

	_, s, err := pn.Time().Execute()

	assert.Nil(err)
	assert.Equal(200, s.StatusCode)
}

func TestNewTimeResponseUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newTimeResponse(jsonBytes, fakeResponseState)
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())

	opts := &timeOpts{}
	a, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal(a, []byte{})
}

func TestNewTimeResponseQueryParam(t *testing.T) {
	assert := assert.New(t)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	config := NewConfigWithUserId(UserId(GenerateUUID()))
	pn := NewPubNub(config)

	opts := &timeOpts{}
	opts.pubnub = pn
	opts.QueryParam = queryParam

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	path, err := opts.buildPath()
	u := &url.URL{
		Path: path,
	}
	assert.Nil(err)

	query, err := opts.buildQuery()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		"/time/0",
		u.EscapedPath(), []int{})

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	a, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal(a, []byte{})
}

func TestNewTimeBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newTimeBuilder(pubnub)
	_, err := o.opts.buildBody()
	assert.Nil(err)
}

func TestNewTimeBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newTimeBuilderWithContext(pubnub, backgroundContext)
	_, err := o.opts.buildBody()
	assert.Nil(err)
}
