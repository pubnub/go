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

	"github.com/pubnub/go/v5/pnerr"
)

var emptyPNRemoveChannelMetadataResponse *PNRemoveChannelMetadataResponse

const removeChannelMetadataPath = "/v2/objects/%s/channels/%s"

type removeChannelMetadataBuilder struct {
	opts *removeChannelMetadataOpts
}

func newRemoveChannelMetadataBuilder(pubnub *PubNub) *removeChannelMetadataBuilder {
	builder := removeChannelMetadataBuilder{
		opts: &removeChannelMetadataOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveChannelMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *removeChannelMetadataBuilder {
	builder := removeChannelMetadataBuilder{
		opts: &removeChannelMetadataOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *removeChannelMetadataBuilder) Channel(channel string) *removeChannelMetadataBuilder {
	b.opts.Channel = channel

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelMetadataBuilder) QueryParam(queryParam map[string]string) *removeChannelMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeChannelMetadata request.
func (b *removeChannelMetadataBuilder) Transport(tr http.RoundTripper) *removeChannelMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeChannelMetadata request.
func (b *removeChannelMetadataBuilder) Execute() (*PNRemoveChannelMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNRemoveChannelMetadataResponse, status, err
	}

	return newPNRemoveChannelMetadataResponse(rawJSON, b.opts, status)
}

type removeChannelMetadataOpts struct {
	pubnub     *PubNub
	Channel    string
	QueryParam map[string]string
	Transport  http.RoundTripper

	ctx Context
}

func (o *removeChannelMetadataOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeChannelMetadataOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeChannelMetadataOpts) context() Context {
	return o.ctx
}

func (o *removeChannelMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *removeChannelMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeChannelMetadataPath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *removeChannelMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeChannelMetadataOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *removeChannelMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil

}

func (o *removeChannelMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeChannelMetadataOpts) httpMethod() string {
	return "DELETE"
}

func (o *removeChannelMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeChannelMetadataOpts) operationType() OperationType {
	return PNRemoveChannelMetadataOperation
}

func (o *removeChannelMetadataOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNRemoveChannelMetadataResponse is the Objects API Response for delete space
type PNRemoveChannelMetadataResponse struct {
	status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newPNRemoveChannelMetadataResponse(jsonBytes []byte, o *removeChannelMetadataOpts,
	status StatusResponse) (*PNRemoveChannelMetadataResponse, StatusResponse, error) {

	resp := &PNRemoveChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNRemoveChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
