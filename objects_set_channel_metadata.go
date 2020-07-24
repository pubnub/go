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

var emptyPNSetChannelMetadataResponse *PNSetChannelMetadataResponse

const setChannelMetadataPath = "/v2/objects/%s/channels/%s"

type setChannelMetadataBuilder struct {
	opts *setChannelMetadataOpts
}

func newSetChannelMetadataBuilder(pubnub *PubNub) *setChannelMetadataBuilder {
	builder := setChannelMetadataBuilder{
		opts: &setChannelMetadataOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newSetChannelMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *setChannelMetadataBuilder {
	builder := setChannelMetadataBuilder{
		opts: &setChannelMetadataOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// SetChannelMetadataBody is the input to update space
type SetChannelMetadataBody struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

func (b *setChannelMetadataBuilder) Include(include []PNChannelMetadataInclude) *setChannelMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *setChannelMetadataBuilder) Channel(channel string) *setChannelMetadataBuilder {
	b.opts.Channel = channel

	return b
}

func (b *setChannelMetadataBuilder) Name(name string) *setChannelMetadataBuilder {
	b.opts.Name = name

	return b
}

func (b *setChannelMetadataBuilder) Description(description string) *setChannelMetadataBuilder {
	b.opts.Description = description

	return b
}

func (b *setChannelMetadataBuilder) Custom(custom map[string]interface{}) *setChannelMetadataBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setChannelMetadataBuilder) QueryParam(queryParam map[string]string) *setChannelMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the setChannelMetadata request.
func (b *setChannelMetadataBuilder) Transport(tr http.RoundTripper) *setChannelMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the setChannelMetadata request.
func (b *setChannelMetadataBuilder) Execute() (*PNSetChannelMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNSetChannelMetadataResponse, status, err
	}

	return newPNSetChannelMetadataResponse(rawJSON, b.opts, status)
}

type setChannelMetadataOpts struct {
	pubnub      *PubNub
	Include     []string
	Channel     string
	Name        string
	Description string
	Custom      map[string]interface{}
	QueryParam  map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *setChannelMetadataOpts) config() Config {
	return *o.pubnub.Config
}

func (o *setChannelMetadataOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *setChannelMetadataOpts) context() Context {
	return o.ctx
}

func (o *setChannelMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *setChannelMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(setChannelMetadataPath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *setChannelMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *setChannelMetadataOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *setChannelMetadataOpts) buildBody() ([]byte, error) {
	b := &SetChannelMetadataBody{
		Name:        o.Name,
		Description: o.Description,
		Custom:      o.Custom,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *setChannelMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *setChannelMetadataOpts) httpMethod() string {
	return "PATCH"
}

func (o *setChannelMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *setChannelMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *setChannelMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *setChannelMetadataOpts) operationType() OperationType {
	return PNSetChannelMetadataOperation
}

func (o *setChannelMetadataOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNSetChannelMetadataResponse is the Objects API Response for Update Space
type PNSetChannelMetadataResponse struct {
	status int       `json:"status"`
	Data   PNChannel `json:"data"`
}

func newPNSetChannelMetadataResponse(jsonBytes []byte, o *setChannelMetadataOpts,
	status StatusResponse) (*PNSetChannelMetadataResponse, StatusResponse, error) {

	resp := &PNSetChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNSetChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
