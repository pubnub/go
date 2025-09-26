package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v7/pnerr"
)

var emptyPNGetUUIDMetadataResponse *PNGetUUIDMetadataResponse

const getUUIDMetadataPath = "/v2/objects/%s/uuids/%s"

type getUUIDMetadataBuilder struct {
	opts *getUUIDMetadataOpts
}

func newGetUUIDMetadataBuilder(pubnub *PubNub) *getUUIDMetadataBuilder {
	return newGetUUIDMetadataBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetUUIDMetadataOpts(pubnub *PubNub, ctx Context) *getUUIDMetadataOpts {
	return &getUUIDMetadataOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newGetUUIDMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *getUUIDMetadataBuilder {
	builder := getUUIDMetadataBuilder{
		opts: newGetUUIDMetadataOpts(pubnub, context)}
	return &builder
}

func (b *getUUIDMetadataBuilder) Include(include []PNUUIDMetadataInclude) *getUUIDMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getUUIDMetadataBuilder) UUID(uuid string) *getUUIDMetadataBuilder {
	b.opts.UUID = uuid

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *getUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getUUIDMetadata request.
func (b *getUUIDMetadataBuilder) Transport(tr http.RoundTripper) *getUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getUUIDMetadata request.
func (b *getUUIDMetadataBuilder) Execute() (*PNGetUUIDMetadataResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetUUIDMetadataResponse, status, err
	}

	return newPNGetUUIDMetadataResponse(rawJSON, b.opts, status)
}

type getUUIDMetadataOpts struct {
	endpointOpts
	UUID       string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getUUIDMetadataPath,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *getUUIDMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getUUIDMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getUUIDMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getUUIDMetadataOpts) httpMethod() string {
	return "GET"
}

func (o *getUUIDMetadataOpts) operationType() OperationType {
	return PNGetUUIDMetadataOperation
}

// PNGetUUIDMetadataResponse is the Objects API Response for Get User
type PNGetUUIDMetadataResponse struct {
	Status int    `json:"status"`
	Data   PNUUID `json:"data"`
}

func newPNGetUUIDMetadataResponse(jsonBytes []byte, o *getUUIDMetadataOpts,
	status StatusResponse) (*PNGetUUIDMetadataResponse, StatusResponse, error) {

	resp := &PNGetUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
