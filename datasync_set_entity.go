package pubnub

import (
	"encoding/json"
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
// entityClass is intentionally omitted because it is immutable.
type setEntityBody struct {
	Data setEntityBodyData `json:"data"`
}

type setEntityBodyData struct {
	EntityClassVersion int                    `json:"entityClassVersion"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the entity to replace.
func (b *setEntityBuilder) ID(id string) *setEntityBuilder {
	b.opts.ID = id
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *setEntityBuilder) EntityClassVersion(version int) *setEntityBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// Status sets the optional entity status (e.g. "active").
func (b *setEntityBuilder) Status(status string) *setEntityBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional user-defined custom properties. This replaces the
// entire payload; omitted fields are removed.
func (b *setEntityBuilder) Payload(payload map[string]interface{}) *setEntityBuilder {
	b.opts.Payload = payload
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
	EntityClassVersion int
	Status             string
	Payload            map[string]interface{}
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
	b := &setEntityBody{
		Data: setEntityBodyData{
			EntityClassVersion: o.EntityClassVersion,
			Status:             o.Status,
			Payload:            o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
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
