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

var emptyPNGetChannelMetadataResponse *PNGetChannelMetadataResponse

const getChannelMetadataPath = "/v2/objects/%s/channels/%s"

type getChannelMetadataBuilder struct {
	opts *getChannelMetadataOpts
}

func newGetChannelMetadataBuilder(pubnub *PubNub) *getChannelMetadataBuilder {
	builder := getChannelMetadataBuilder{
		opts: &getChannelMetadataOpts{
			pubnub: pubnub,
		},
	}
	return &builder
}

func newGetChannelMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *getChannelMetadataBuilder {
	builder := getChannelMetadataBuilder{
		opts: &getChannelMetadataOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getChannelMetadataBuilder) Include(include []PNChannelMetadataInclude) *getChannelMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getChannelMetadataBuilder) Channel(channel string) *getChannelMetadataBuilder {
	b.opts.Channel = channel

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getChannelMetadataBuilder) QueryParam(queryParam map[string]string) *getChannelMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getChannelMetadata request.
func (b *getChannelMetadataBuilder) Transport(tr http.RoundTripper) *getChannelMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getChannelMetadata request.
func (b *getChannelMetadataBuilder) Execute() (*PNGetChannelMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetChannelMetadataResponse, status, err
	}

	return newPNGetChannelMetadataResponse(rawJSON, b.opts, status)
}

type getChannelMetadataOpts struct {
	pubnub     *PubNub
	Channel    string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getChannelMetadataOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getChannelMetadataOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getChannelMetadataOpts) context() Context {
	return o.ctx
}

func (o *getChannelMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *getChannelMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getChannelMetadataPath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *getChannelMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getChannelMetadataOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getChannelMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getChannelMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getChannelMetadataOpts) httpMethod() string {
	return "GET"
}

func (o *getChannelMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getChannelMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getChannelMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getChannelMetadataOpts) operationType() OperationType {
	return PNGetChannelMetadataOperation
}

func (o *getChannelMetadataOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetChannelMetadataResponse is the Objects API Response for Get Space
type PNGetChannelMetadataResponse struct {
	status int       `json:"status"`
	Data   PNChannel `json:"data"`
}

func newPNGetChannelMetadataResponse(jsonBytes []byte, o *getChannelMetadataOpts,
	status StatusResponse) (*PNGetChannelMetadataResponse, StatusResponse, error) {

	resp := &PNGetChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
