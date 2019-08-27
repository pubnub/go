package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sprucehealth/pubnub-go/pnerr"
)

var emptyPNDeleteSpaceResponse *PNDeleteSpaceResponse

const deleteSpacePath = "/v1/objects/%s/spaces/%s"

type deleteSpaceBuilder struct {
	opts *deleteSpaceOpts
}

func newDeleteSpaceBuilder(pubnub *PubNub) *deleteSpaceBuilder {
	builder := deleteSpaceBuilder{
		opts: &deleteSpaceOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newDeleteSpaceBuilderWithContext(pubnub *PubNub,
	context Context) *deleteSpaceBuilder {
	builder := deleteSpaceBuilder{
		opts: &deleteSpaceOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *deleteSpaceBuilder) ID(id string) *deleteSpaceBuilder {
	b.opts.ID = id

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteSpaceBuilder) QueryParam(queryParam map[string]string) *deleteSpaceBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the deleteSpace request.
func (b *deleteSpaceBuilder) Transport(tr http.RoundTripper) *deleteSpaceBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the deleteSpace request.
func (b *deleteSpaceBuilder) Execute() (*PNDeleteSpaceResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNDeleteSpaceResponse, status, err
	}

	return newPNDeleteSpaceResponse(rawJSON, b.opts, status)
}

type deleteSpaceOpts struct {
	pubnub     *PubNub
	ID         string
	QueryParam map[string]string
	Transport  http.RoundTripper

	ctx Context
}

func (o *deleteSpaceOpts) config() Config {
	return *o.pubnub.Config
}

func (o *deleteSpaceOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *deleteSpaceOpts) context() Context {
	return o.ctx
}

func (o *deleteSpaceOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *deleteSpaceOpts) buildPath() (string, error) {
	return fmt.Sprintf(deleteSpacePath,
		o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *deleteSpaceOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *deleteSpaceOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *deleteSpaceOpts) buildBody() ([]byte, error) {
	return []byte{}, nil

}

func (o *deleteSpaceOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteSpaceOpts) isAuthRequired() bool {
	return true
}

func (o *deleteSpaceOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *deleteSpaceOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *deleteSpaceOpts) operationType() OperationType {
	return PNDeleteSpaceOperation
}

func (o *deleteSpaceOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNDeleteSpaceResponse is the Objects API Response for delete space
type PNDeleteSpaceResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newPNDeleteSpaceResponse(jsonBytes []byte, o *deleteSpaceOpts,
	status StatusResponse) (*PNDeleteSpaceResponse, StatusResponse, error) {

	resp := &PNDeleteSpaceResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNDeleteSpaceResponse, status, e
	}

	return resp, status, nil
}
