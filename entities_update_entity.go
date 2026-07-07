package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateEntityResponse *PNEntityResponse

type updateEntityBuilder struct {
	opts *updateEntityOpts
}

func newUpdateEntityBuilder(pubnub *PubNub) *updateEntityBuilder {
	return newUpdateEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateEntityOpts(pubnub *PubNub, ctx Context) *updateEntityOpts {
	return &updateEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateEntityBuilderWithContext(pubnub *PubNub, context Context) *updateEntityBuilder {
	return &updateEntityBuilder{opts: newUpdateEntityOpts(pubnub, context)}
}

// updateEntityBody is the request envelope sent on PUT /entities/{id}.
// entityClass is intentionally omitted because it is immutable.
type updateEntityBody struct {
	Data updateEntityBodyData `json:"data"`
}

type updateEntityBodyData struct {
	EntityClassVersion int                    `json:"entityClassVersion"`
	Status             string                 `json:"status,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the entity to replace.
func (b *updateEntityBuilder) ID(id string) *updateEntityBuilder {
	b.opts.ID = id
	return b
}

// EntityClassVersion sets the required schema version of the entity class (>= 1).
func (b *updateEntityBuilder) EntityClassVersion(version int) *updateEntityBuilder {
	b.opts.EntityClassVersion = version
	return b
}

// Status sets the optional entity status (e.g. "active").
func (b *updateEntityBuilder) Status(status string) *updateEntityBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional user-defined custom properties. This replaces the
// entire payload; omitted fields are removed.
func (b *updateEntityBuilder) Payload(payload map[string]interface{}) *updateEntityBuilder {
	b.opts.Payload = payload
	return b
}

// IfMatch sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateEntityBuilder) IfMatch(eTag string) *updateEntityBuilder {
	b.opts.IfMatch = eTag
	b.opts.setIfMatch = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateEntityBuilder) QueryParam(queryParam map[string]string) *updateEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateEntity request.
func (b *updateEntityBuilder) Transport(tr http.RoundTripper) *updateEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateEntityOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the updateEntity request.
func (b *updateEntityBuilder) Execute() (*PNEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateEntityResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNUpdateEntityOperation, b.opts.pubnub)
}

type updateEntityOpts struct {
	endpointOpts

	ID                 string
	EntityClassVersion int
	Status             string
	Payload            map[string]interface{}
	IfMatch            string
	setIfMatch         bool
	QueryParam         map[string]string

	Transport http.RoundTripper
}

func (o *updateEntityOpts) validate() error {
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

func (o *updateEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateEntityOpts) buildBody() ([]byte, error) {
	b := &updateEntityBody{
		Data: updateEntityBodyData{
			EntityClassVersion: o.EntityClassVersion,
			Status:             o.Status,
			Payload:            o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateEntitySerializationFailed", PNUpdateEntityOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateEntityOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityContentType,
	}
	if o.setIfMatch {
		headers["If-Match"] = o.IfMatch
	}
	return headers, nil
}

func (o *updateEntityOpts) httpMethod() string {
	return "PUT"
}

func (o *updateEntityOpts) operationType() OperationType {
	return PNUpdateEntityOperation
}
