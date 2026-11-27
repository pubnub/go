package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptySetEntityResponse *PNEntityResponse

type setEntityBuilder struct {
	opts *setEntityOpts
}

func newSetEntityBuilder(pubnub *PubNub) *setEntityBuilder {
	return newSetEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetEntityOpts(pubnub *PubNub, ctx Context) *setEntityOpts {
	return &setEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newSetEntityBuilderWithContext(pubnub *PubNub, context Context) *setEntityBuilder {
	return &setEntityBuilder{opts: newSetEntityOpts(pubnub, context)}
}

// setEntityBody is the request envelope sent on PUT /entities/{id}.
// entityClass and entityClassLevel are omitted unless their setters are called.
type setEntityBody struct {
	Data setEntityBodyData `json:"data"`
}

type setEntityBodyData struct {
	EntityClass        string                 `json:"entityClass,omitempty"`
	EntityClassVersion int                    `json:"entityClassVersion"`
	EntityClassLevel   PNEntityClassLevel     `json:"entityClassLevel,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the entity to create or replace.
func (b *setEntityBuilder) ID(id string) *setEntityBuilder {
	b.opts.ID = id
	return b
}

// EntityClass sets the optional entity class name, sent as entityClass.
// The server requires it when this PUT creates a missing entity. On
// replacement, a supplied name must match the stored class; omitting it
// makes no assertion.
func (b *setEntityBuilder) EntityClass(entityClass string) *setEntityBuilder {
	b.opts.EntityClass = entityClass
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *setEntityBuilder) EntityClassVersion(version int) *setEntityBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// EntityClassLevel sets the optional class scope (Global or SubKey), sent as
// entityClassLevel. On creation it disambiguates classes that share a name.
// On replacement, a supplied level must match the stored identity; omitting it
// makes no assertion.
func (b *setEntityBuilder) EntityClassLevel(level PNEntityClassLevel) *setEntityBuilder {
	b.opts.EntityClassLevel = level
	return b
}

// Status sets the optional entity status (e.g. "active").
func (b *setEntityBuilder) Status(status string) *setEntityBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the user-defined custom properties. This replaces the entire
// payload; omitted keys are removed. An empty map is sent as {}. Not calling
// Payload, or passing nil, omits the field so the server resets it to default.
func (b *setEntityBuilder) Payload(payload map[string]interface{}) *setEntityBuilder {
	b.opts.Payload = payload
	b.opts.setPayload = true
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *setEntityBuilder) IfMatchETag(eTag string) *setEntityBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setEntityBuilder) QueryParam(queryParam map[string]string) *setEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the setEntity request.
func (b *setEntityBuilder) Transport(tr http.RoundTripper) *setEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *setEntityOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the setEntity request.
func (b *setEntityBuilder) Execute() (*PNEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNSetEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetEntityResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNSetEntityOperation, b.opts.pubnub)
}

type setEntityOpts struct {
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

func (o *setEntityOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingEntityID)
	}
	if o.EntityClassVersion < 1 {
		return newValidationError(o, StrInvalidEntityClassVersion)
	}
	return nil
}

func (o *setEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *setEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *setEntityOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := marshalDataSyncEntitySetBody(o.EntityClassVersion, o.EntityClass, o.EntityClassLevel, o.Status, o.Payload, o.setPayload)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "SetEntitySerializationFailed", PNSetEntityOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setEntityOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *setEntityOpts) httpMethod() string {
	return "PUT"
}

func (o *setEntityOpts) operationType() OperationType {
	return PNSetEntityOperation
}
