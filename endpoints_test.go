package pubnub

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeEndpointOpts struct {
	pubnub *PubNub
}

func (o *fakeEndpointOpts) buildPath() (string, error) {
	return "/my/path", nil
}

func (o *fakeEndpointOpts) buildQuery() (*url.Values, error) {
	q := &url.Values{}

	q.Set("a", "2")

	q.Set("b", "hey")

	return q, nil
}

func (o *fakeEndpointOpts) buildBody() ([]byte, error) {
	return []byte("myBody"), nil
}

func (o *fakeEndpointOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *fakeEndpointOpts) config() Config {
	return *o.pubnub.Config
}

func (o *fakeEndpointOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *fakeEndpointOpts) validate() error {
	return nil
}

func (o *fakeEndpointOpts) context() Context {
	return o.context()
}

func (o *fakeEndpointOpts) httpMethod() string {
	return "GET"
}

func (o *fakeEndpointOpts) operationType() OperationType {
	return PNSubscribeOperation
}

func (o *fakeEndpointOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

func xTestBuildURL(t *testing.T) {
	assert := assert.New(t)

	pnc := NewConfig()
	pn := NewPubNub(pnc)

	e := &fakeEndpointOpts{
		pubnub: pn,
	}

	url, err := buildURL(e)

	assert.Nil(err)
	assert.Equal("https://ps.pndsn.com/my/path?a=2&b=hey", url.RequestURI())
}

func TestSignatureV2(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "POST"
	pubKey := "demo"
	secKey := "wMfbo9G0xVUG8yfTfYw5qIdfJkTd7A"
	path := "/v3/pam/demo/grant"
	query := "PoundsSterling=%C2%A313.37&timestamp=123456789"
	body := `{
  "ttl": 1440,
  "permissions": {
    "resources" : {
      "channels": {
        "inbox-jay": 3
      },
      "groups": {},
      "users": {},
      "spaces": {}
    },
    "patterns" : {
      "channels": {},
      "groups": {},
      "users": {},
      "spaces": {}
    },
    "meta": {
      "user-id": "jay@example.com",
      "contains-unicode": "The 💩 test."
    }
  }
}`
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body, nil)
	assert.Equal("v2.k80LsDMD-sImA8rCBj-ntRKhZ8mSjHY8Ivngt9W3Yc4", sigv2)
}

func TestSignatureV2_1(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "GET"
	pubKey := "pub-c-03f156ea-a2e9-4c35-a733-9535824be897"
	secKey := "sec-c-MmUxNTZjMmYtNzFkNS00OAkzLWE2YjctNmQ4YzE5NWNmZDA3"
	path := "/v1/objects/sub-c-d7da9e59-c997-11e9-a139-dab2c75acd6f/spaces/pandu-ut-sid/users"
	query := "l_obj=0.4545555556&l_pam=1.145&pnsdk=PubNubCSharp4.0.34.0&requestid=19e1dee9-2f87-45d6-97e5-3f4d3f9779a2&timestamp=1568724043&uuid=mytestuuid"
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, "", nil)
	assert.Equal("v2.TcZdUURiXAnxJgN4OLPczxzH4MQO87l-yKfE4fyUHGc", sigv2)
}

func TestSignatureV2_2(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "GET"
	pubKey := "pub-c-38994634-9e06-4967-bc66-2ac2cef65ed9"
	secKey := "sec-c-ZDkzZTBkOTEtNTQxZS00VmQ3LTljMWUtMTNiNGZjNWUwMTVk"
	path := "/v1/objects/sub-c-c9710929-1b7a-11e3-a0c8-02ee2ddab7fe/users"
	query := "count=true&filter=name%3D%3D%27newnamemodified%27&include=custom&pnsdk=NET461CSharp4.5.0.0&timestamp=1577364610&uuid=pn-856feea0-5ba3-4c6d-85a2-a67efa1f4e20"
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, "", nil)
	assert.Equal("v2.w39EWDHc0ibDCxOD2jiJEutjpc_1VJ1SDyGaABEljSs", sigv2)
}
