package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

var emptyPNRemoveUUIDMetadataResponse *PNRemoveUUIDMetadataResponse

const removeUUIDMetadataPath = "/v2/objects/%s/uuids/%s"

type removeUUIDMetadataBuilder struct {
	opts *removeUUIDMetadataOpts
}

func newRemoveUUIDMetadataBuilder(pubnub *PubNub) *removeUUIDMetadataBuilder {
	builder := removeUUIDMetadataBuilder{
		opts: &removeUUIDMetadataOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveUUIDMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *removeUUIDMetadataBuilder {
	builder := removeUUIDMetadataBuilder{
		opts: &removeUUIDMetadataOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *removeUUIDMetadataBuilder) UUID(uuid string) *removeUUIDMetadataBuilder {
	b.opts.UUID = uuid

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *removeUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeUUIDMetadata request.
func (b *removeUUIDMetadataBuilder) Transport(tr http.RoundTripper) *removeUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeUUIDMetadata request.
func (b *removeUUIDMetadataBuilder) Execute() (*PNRemoveUUIDMetadataResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNRemoveUUIDMetadataResponse, status, err
	}

	return newPNRemoveUUIDMetadataResponse(rawJSON, b.opts, status)
}

type removeUUIDMetadataOpts struct {
	pubnub     *PubNub
	UUID       string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeUUIDMetadataOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeUUIDMetadataOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeUUIDMetadataOpts) context() Context {
	return o.ctx
}

func (o *removeUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeUUIDMetadataPath,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *removeUUIDMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeUUIDMetadataOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *removeUUIDMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil

}

func (o *removeUUIDMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeUUIDMetadataOpts) httpMethod() string {
	return "DELETE"
}

func (o *removeUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *removeUUIDMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeUUIDMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeUUIDMetadataOpts) operationType() OperationType {
	return PNRemoveUUIDMetadataOperation
}

func (o *removeUUIDMetadataOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNRemoveUUIDMetadataResponse is the Objects API Response for delete user
type PNRemoveUUIDMetadataResponse struct {
	status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newPNRemoveUUIDMetadataResponse(jsonBytes []byte, o *removeUUIDMetadataOpts,
	status StatusResponse) (*PNRemoveUUIDMetadataResponse, StatusResponse, error) {

	resp := &PNRemoveUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNRemoveUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
