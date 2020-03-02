package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

var emptyPNGetUserResponse *PNGetUserResponse

const getUserPath = "/v1/objects/%s/users/%s"

type getUserBuilder struct {
	opts *getUserOpts
}

func newGetUserBuilder(pubnub *PubNub) *getUserBuilder {
	builder := getUserBuilder{
		opts: &getUserOpts{
			pubnub: pubnub,
		},
	}
	return &builder
}

func newGetUserBuilderWithContext(pubnub *PubNub,
	context Context) *getUserBuilder {
	builder := getUserBuilder{
		opts: &getUserOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getUserBuilder) Include(include []PNUserSpaceInclude) *getUserBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getUserBuilder) ID(id string) *getUserBuilder {
	b.opts.ID = id

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUserBuilder) QueryParam(queryParam map[string]string) *getUserBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getUser request.
func (b *getUserBuilder) Transport(tr http.RoundTripper) *getUserBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getUser request.
func (b *getUserBuilder) Execute() (*PNGetUserResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetUserResponse, status, err
	}

	return newPNGetUserResponse(rawJSON, b.opts, status)
}

type getUserOpts struct {
	pubnub     *PubNub
	ID         string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getUserOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getUserOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getUserOpts) context() Context {
	return o.ctx
}

func (o *getUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(getUserPath,
		o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getUserOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetArrayTypeQueryParam(q, o.Include, "include")
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNUsers)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getUserOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getUserOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getUserOpts) httpMethod() string {
	return "GET"
}

func (o *getUserOpts) isAuthRequired() bool {
	return true
}

func (o *getUserOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getUserOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getUserOpts) operationType() OperationType {
	return PNGetUserOperation
}

func (o *getUserOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetUserResponse is the Objects API Response for Get User
type PNGetUserResponse struct {
	status int    `json:"status"`
	Data   PNUser `json:"data"`
}

func newPNGetUserResponse(jsonBytes []byte, o *getUserOpts,
	status StatusResponse) (*PNGetUserResponse, StatusResponse, error) {

	resp := &PNGetUserResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetUserResponse, status, e
	}

	return resp, status, nil
}
