package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v8/pnerr"
)

var emptyPNSetChannelMetadataResponse *PNSetChannelMetadataResponse

const setChannelMetadataPath = "/v2/objects/%s/channels/%s"

type setChannelMetadataBuilder struct {
	opts *setChannelMetadataOpts
}

func newSetChannelMetadataBuilder(pubnub *PubNub) *setChannelMetadataBuilder {
	return newSetChannelMetadataBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetChannelMetadataOpts(pubnub *PubNub, ctx Context) *setChannelMetadataOpts {
	return &setChannelMetadataOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newSetChannelMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *setChannelMetadataBuilder {
	builder := setChannelMetadataBuilder{
		opts: newSetChannelMetadataOpts(pubnub, context)}
	return &builder
}

// SetChannelMetadataBody is the input to update space
type SetChannelMetadataBody struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Type        string                 `json:"type,omitempty"`
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

func (b *setChannelMetadataBuilder) Status(status string) *setChannelMetadataBuilder {
	b.opts.Status = status

	return b
}

// Called channelType instead of type because type is a reserved word in Go
func (b *setChannelMetadataBuilder) Type(channelType string) *setChannelMetadataBuilder {
	b.opts.Type = channelType

	return b
}

// IfMatchETag sets the ETag value for conditional update via If-Match header
func (b *setChannelMetadataBuilder) IfMatchETag(eTag string) *setChannelMetadataBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true

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
	endpointOpts
	Include        []string
	Channel        string
	Name           string
	Description    string
	Custom         map[string]interface{}
	Status         string
	Type           string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
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

func (o *setChannelMetadataOpts) buildBody() ([]byte, error) {
	b := &SetChannelMetadataBody{
		Name:        o.Name,
		Description: o.Description,
		Custom:      o.Custom,
		Status:      o.Status,
		Type:        o.Type,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

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

// buildHeaders adds the If-Match header for conditional updates when IfMatchETag is set.
// The If-Match header enables optimistic concurrency control using ETags.
func (o *setChannelMetadataOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

// PNSetChannelMetadataResponse is the Objects API Response for Update Space
type PNSetChannelMetadataResponse struct {
	Status int       `json:"status"`
	Data   PNChannel `json:"data"`
}

func newPNSetChannelMetadataResponse(jsonBytes []byte, o *setChannelMetadataOpts,
	status StatusResponse) (*PNSetChannelMetadataResponse, StatusResponse, error) {

	resp := &PNSetChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNSetChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
