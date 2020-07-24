package pubnub

import (
	"bytes"
	"errors"
	"mime/multipart"
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

func (o *fakeEndpointOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
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

func TestSignatureV2_granttoken(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "POST"
	pubKey := "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	secKey := "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"
	path := "/v3/pam/sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306/grant"
	query := "pnsdk=PubNub-Go%2F4.6.5&timestamp=1583755527&uuid=pn-f1df31f1-6ba9-49e1-ae72-084c15404302"
	body := `{"ttl":10,"permissions":{"resources":{"channels":{},"groups":{},"users":{"u1":7,"u2":7},"spaces":{"s1":31,"s2":31}},"patterns":{"channels":{},"groups":{},"users":{"^u-[0-9a-f]*$":15},"spaces":{"^s-[0-9a-f]*$":27}},"meta":{}}}`
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body, nil)
	assert.Equal("v2.Dd2YCsd-Ds_fFpGyuHIetndf2UrgZ-FTwcLggcPfx98", sigv2)
}

func TestSignatureV2_3(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "POST"
	pubKey := "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	secKey := "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"
	path := "/v3/pam/sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306/grant"
	query := "timestamp=1583479519"
	body := `{"ttl":"3","permissions":{"resources":{"channels":{},"groups":{},"users":{"userid9820":31,"userid1962":31},"spaces":{"spaceid9820":31,"spaceid1962":31}},"patterns":{"channels":{},"groups":{},"users":{},"spaces":{}},"meta":{}}}`
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body, nil)
	assert.Equal("v2.GbQoBigWGAHX9BAFOlr84nW8XO9QgNob9g21KEr1KfA", sigv2)
}

func TestSignatureV2_4(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "POST"
	pubKey := "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	secKey := "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"
	path := "/v3/pam/sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306/grant"
	query := "timestamp=1583486135"
	body := `{"ttl":"3","permissions":{"resources":{"channels":{},"groups":{},"users":{"userid6155":31,"userid7504":31},"spaces":{"spaceid6155":31,"spaceid7504":31}},"patterns":{"channels":{},"groups":{},"users":{},"spaces":{}},"meta":{}}}`
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body, nil)
	assert.Equal("v2.Md5iDqxLLCU1wj3L6wUmsrbW6IrOuhZyPEi1AwdsELs", sigv2)
}

func TestSignatureV2_5(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "POST"
	pubKey := "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	secKey := "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"
	path := "/v3/pam/sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306/grant"
	query := "timestamp=1583486728"
	body := `{"ttl":"3","permissions":{"resources":{"channels":{},"groups":{},"users":{"userid4289":31,"userid4476":31},"spaces":{"spaceid4289":31,"spaceid4476":31}},"patterns":{"channels":{},"groups":{},"users":{},"spaces":{}},"meta":{}}}`
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body, nil)
	assert.Equal("v2._LxrtaeGughFMF7aFMnzqpLAjM9JT4ELlUjVNskXgDs", sigv2)
}

func TestSignatureV2_6(t *testing.T) {
	assert := assert.New(t)
	httpMethod := "POST"
	pubKey := "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	secKey := "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"
	path := "/v1/objects/sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306/users"
	query := "include=custom&pnsdk=PubNub-Go/4.6.5&timestamp=1583493622&uuid=pn-01fb5308-f0ce-480b-b051-6ff98ba22467"
	body := `{"id":"id0","name":"name","externalId":"extid","profileUrl":"purl","email":"email","custom":{"a":"b","c":"d"}}`
	sigv2 := createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body, nil)
	assert.Equal("v2.a--gef4a6Rm3Oe7k2pPOP9IbjRPWi5Ky-RKpIcQIYn0", sigv2)
}
