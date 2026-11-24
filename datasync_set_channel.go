package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptySetChannelResponse *PNDataSyncChannelResponse

type setChannelBuilder struct {
	opts *setChannelOpts
}

func newSetChannelBuilder(pubnub *PubNub) *setChannelBuilder {
	return newSetChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetChannelOpts(pubnub *PubNub, ctx Context) *setChannelOpts {
	return &setChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newSetChannelBuilderWithContext(pubnub *PubNub, context Context) *setChannelBuilder {
	return &setChannelBuilder{opts: newSetChannelOpts(pubnub, context)}
}

// setChannelBody is the request envelope sent on PUT /channels/{id}.
// entityClass and entityClassLevel are omitted unless their setters are called.
type setChannelBody struct {
	Data setChannelBodyData `json:"data"`
}

type setChannelBodyData struct {
	EntityClass        string                 `json:"entityClass,omitempty"`
	EntityClassVersion int                    `json:"entityClassVersion"`
	EntityClassLevel   PNEntityClassLevel     `json:"entityClassLevel,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the channel to create or replace.
func (b *setChannelBuilder) ID(id string) *setChannelBuilder {
	b.opts.ID = id
	return b
}

// EntityClass sets the optional channel class, sent as entityClass. On
// creation, omit it to use the Global Channel class; an explicit class must be
// Channel or a descendant. On replacement, a supplied name must match the
// stored class; omitting it makes no assertion.
func (b *setChannelBuilder) EntityClass(entityClass string) *setChannelBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *setChannelBuilder) EntityClassVersion(version int) *setChannelBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// EntityClassLevel sets the optional class scope (Global or SubKey), sent as
// entityClassLevel. On creation it disambiguates classes that share a name.
// On replacement, a supplied level must match the stored identity; omitting it
// makes no assertion.
func (b *setChannelBuilder) EntityClassLevel(level PNEntityClassLevel) *setChannelBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Status sets the optional channel status (e.g. "active").
func (b *setChannelBuilder) Status(status string) *setChannelBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the custom channel fields. Under the default PAM projection this
// replaces the entire payload; omitted keys are removed. An empty map is sent
// as {}. Not calling Payload, or passing nil, omits the field so the server
// resets it to default.
func (b *setChannelBuilder) Payload(payload map[string]interface{}) *setChannelBuilder {
	b.opts.Payload = payload
	b.opts.setPayload = true
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *setChannelBuilder) IfMatchETag(eTag string) *setChannelBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setChannelBuilder) QueryParam(queryParam map[string]string) *setChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the setChannel request.
func (b *setChannelBuilder) Transport(tr http.RoundTripper) *setChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *setChannelOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"ID":                 o.ID,
		"EntityClassVersion": o.EntityClassVersion,
	}
	if o.EntityClass != "" {
		params["EntityClass"] = o.EntityClass
	}
	if o.EntityClassLevel != "" {
		params["EntityClassLevel"] = o.EntityClassLevel
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the setChannel request.
func (b *setChannelBuilder) Execute() (*PNDataSyncChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNSetDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetChannelResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNSetDataSyncChannelOperation, b.opts.pubnub)
}

type setChannelOpts struct {
	endpointOpts

	ID                 string
	EntityClass        string
	EntityClassVersion int
	EntityClassLevel   PNEntityClassLevel
	Status             string
	Payload            map[string]interface{}
	setPayload         bool
	IfMatchETag        string
	setIfMatchETag     bool
	QueryParam         map[string]string

	Transport http.RoundTripper
}

func (o *setChannelOpts) validate() error {
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

func (o *setChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *setChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *setChannelOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := marshalDataSyncEntitySetBody(o.EntityClassVersion, o.EntityClass, o.EntityClassLevel, o.Status, o.Payload, o.setPayload)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "SetChannelSerializationFailed", PNSetDataSyncChannelOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setChannelOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": channelContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *setChannelOpts) httpMethod() string {
	return "PUT"
}

func (o *setChannelOpts) operationType() OperationType {
	return PNSetDataSyncChannelOperation
}
