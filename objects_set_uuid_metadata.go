package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v7/pnerr"
)

var emptyPNSetUUIDMetadataResponse *PNSetUUIDMetadataResponse

const setUUIDMetadataPath = "/v2/objects/%s/uuids/%s"

type setUUIDMetadataBuilder struct {
	opts *setUUIDMetadataOpts
}

func newSetUUIDMetadataBuilder(pubnub *PubNub) *setUUIDMetadataBuilder {
	return newSetUUIDMetadataBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetUUIDMetadataOpts(pubnub *PubNub, ctx Context) *setUUIDMetadataOpts {
	return &setUUIDMetadataOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newSetUUIDMetadataBuilderWithContext(pubnub *PubNub,
	context Context) *setUUIDMetadataBuilder {
	return &setUUIDMetadataBuilder{
		opts: newSetUUIDMetadataOpts(pubnub, context)}
}

// SetUUIDMetadataBody is the input to update user
type SetUUIDMetadataBody struct {
	Name       string                 `json:"name,omitempty"`
	ExternalID string                 `json:"externalId,omitempty"`
	ProfileURL string                 `json:"profileUrl,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Custom     map[string]interface{} `json:"custom,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Type       string                 `json:"type,omitempty"`
}

func (b *setUUIDMetadataBuilder) UUID(uuid string) *setUUIDMetadataBuilder {
	b.opts.UUID = uuid

	return b
}

func (b *setUUIDMetadataBuilder) Include(include []PNUUIDMetadataInclude) *setUUIDMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *setUUIDMetadataBuilder) Name(name string) *setUUIDMetadataBuilder {
	b.opts.Name = name

	return b
}

func (b *setUUIDMetadataBuilder) ExternalID(externalID string) *setUUIDMetadataBuilder {
	b.opts.ExternalID = externalID

	return b
}

func (b *setUUIDMetadataBuilder) ProfileURL(profileURL string) *setUUIDMetadataBuilder {
	b.opts.ProfileURL = profileURL

	return b
}

func (b *setUUIDMetadataBuilder) Email(email string) *setUUIDMetadataBuilder {
	b.opts.Email = email

	return b
}

func (b *setUUIDMetadataBuilder) Custom(custom map[string]interface{}) *setUUIDMetadataBuilder {
	b.opts.Custom = custom

	return b
}

func (b *setUUIDMetadataBuilder) Status(status string) *setUUIDMetadataBuilder {
	b.opts.Status = status

	return b
}

// Called uuidType instead of type because type is a reserved word in Go
func (b *setUUIDMetadataBuilder) Type(uuidType string) *setUUIDMetadataBuilder {
	b.opts.Type = uuidType

	return b
}

// IfMatchETag sets the ETag value for conditional update via If-Match header
func (b *setUUIDMetadataBuilder) IfMatchETag(eTag string) *setUUIDMetadataBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *setUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the setUUIDMetadata request.
func (b *setUUIDMetadataBuilder) Transport(tr http.RoundTripper) *setUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the setUUIDMetadata request.
func (b *setUUIDMetadataBuilder) Execute() (*PNSetUUIDMetadataResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.pubnub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNSetUUIDMetadataResponse, status, err
	}

	return newPNSetUUIDMetadataResponse(rawJSON, b.opts, status)
}

type setUUIDMetadataOpts struct {
	endpointOpts
	Include        []string
	UUID           string
	Name           string
	ExternalID     string
	ProfileURL     string
	Email          string
	Custom         map[string]interface{}
	Status         string
	Type           string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *setUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *setUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(setUUIDMetadataPath,
		o.pubnub.Config.SubscribeKey, o.UUID), nil
}

func (o *setUUIDMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *setUUIDMetadataOpts) buildBody() ([]byte, error) {
	b := &SetUUIDMetadataBody{
		Name:       o.Name,
		ExternalID: o.ExternalID,
		ProfileURL: o.ProfileURL,
		Email:      o.Email,
		Custom:     o.Custom,
		Status:     o.Status,
		Type:       o.Type,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *setUUIDMetadataOpts) httpMethod() string {
	return "PATCH"
}

func (o *setUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *setUUIDMetadataOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *setUUIDMetadataOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *setUUIDMetadataOpts) operationType() OperationType {
	return PNSetUUIDMetadataOperation
}

// buildHeaders adds the If-Match header for conditional updates when IfMatchETag is set.
// The If-Match header enables optimistic concurrency control using ETags.
func (o *setUUIDMetadataOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

// PNSetUUIDMetadataResponse is the Objects API Response for Update user
type PNSetUUIDMetadataResponse struct {
	Status int    `json:"status"`
	Data   PNUUID `json:"data"`
}

func newPNSetUUIDMetadataResponse(jsonBytes []byte, o *setUUIDMetadataOpts,
	status StatusResponse) (*PNSetUUIDMetadataResponse, StatusResponse, error) {

	resp := &PNSetUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNSetUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
