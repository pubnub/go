package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateChannelResponse *PNDataSyncChannelResponse

type updateChannelBuilder struct {
	opts *updateChannelOpts
}

func newUpdateChannelBuilder(pubnub *PubNub) *updateChannelBuilder {
	return newUpdateChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateChannelOpts(pubnub *PubNub, ctx Context) *updateChannelOpts {
	return &updateChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateChannelBuilderWithContext(pubnub *PubNub, context Context) *updateChannelBuilder {
	return &updateChannelBuilder{opts: newUpdateChannelOpts(pubnub, context)}
}

// updateChannelBody is the request envelope sent on PUT /channels/{id}.
// entityClass is intentionally omitted because it is immutable.
type updateChannelBody struct {
	Data updateChannelBodyData `json:"data"`
}

type updateChannelBodyData struct {
	EntityClassVersion int                    `json:"entityClassVersion"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the channel to replace.
func (b *updateChannelBuilder) ID(id string) *updateChannelBuilder {
	b.opts.ID = id
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *updateChannelBuilder) EntityClassVersion(version int) *updateChannelBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// Status sets the optional channel status (e.g. "active").
func (b *updateChannelBuilder) Status(status string) *updateChannelBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom channel fields. Under the default PAM
// projection this replaces the entire payload; omitted fields are removed.
func (b *updateChannelBuilder) Payload(payload map[string]interface{}) *updateChannelBuilder {
	b.opts.Payload = payload
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateChannelBuilder) IfMatchETag(eTag string) *updateChannelBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateChannelBuilder) QueryParam(queryParam map[string]string) *updateChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateChannel request.
func (b *updateChannelBuilder) Transport(tr http.RoundTripper) *updateChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateChannelOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"ID":                 o.ID,
		"EntityClassVersion": o.EntityClassVersion,
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the updateChannel request.
func (b *updateChannelBuilder) Execute() (*PNDataSyncChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateChannelResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNUpdateDataSyncChannelOperation, b.opts.pubnub)
}

type updateChannelOpts struct {
	endpointOpts

	ID                 string
	EntityClassVersion int
	Status             string
	Payload            map[string]interface{}
	IfMatchETag        string
	setIfMatchETag     bool
	QueryParam         map[string]string

	Transport http.RoundTripper
}

func (o *updateChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingChannelID)
	}
	if o.EntityClassVersion < 1 {
		return newValidationError(o, StrInvalidEntityClassVersion)
	}
	return nil
}

func (o *updateChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateChannelOpts) buildBody() ([]byte, error) {
	b := &updateChannelBody{
		Data: updateChannelBodyData{
			EntityClassVersion: o.EntityClassVersion,
			Status:             o.Status,
			Payload:            o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateChannelSerializationFailed", PNUpdateDataSyncChannelOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateChannelOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": channelContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateChannelOpts) httpMethod() string {
	return "PUT"
}

func (o *updateChannelOpts) operationType() OperationType {
	return PNUpdateDataSyncChannelOperation
}
