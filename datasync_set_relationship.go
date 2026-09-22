package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptySetRelationshipResponse *PNRelationshipResponse

type setRelationshipBuilder struct {
	opts *setRelationshipOpts
}

func newSetRelationshipBuilder(pubnub *PubNub) *setRelationshipBuilder {
	return newSetRelationshipBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetRelationshipOpts(pubnub *PubNub, ctx Context) *setRelationshipOpts {
	return &setRelationshipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newSetRelationshipBuilderWithContext(pubnub *PubNub, context Context) *setRelationshipBuilder {
	return &setRelationshipBuilder{opts: newSetRelationshipOpts(pubnub, context)}
}

// setRelationshipBody is the request envelope sent on PUT /relationships/{id}.
// entityAId, entityBId and relationshipClass are intentionally omitted because
// they are immutable after creation (same pattern as SetEntity).
type setRelationshipBody struct {
	Data setRelationshipBodyData `json:"data"`
}

type setRelationshipBodyData struct {
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Status                   string                 `json:"status,omitempty"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the relationship to replace.
func (b *setRelationshipBuilder) ID(id string) *setRelationshipBuilder {
	b.opts.ID = id
	return b
}

// RelationshipClassVersion sets the required schema version of the relationship class (>= 1).
func (b *setRelationshipBuilder) RelationshipClassVersion(version int) *setRelationshipBuilder {
	b.opts.RelationshipClassVersion = version
	return b
}

// Status sets the optional relationship status (e.g. "active").
func (b *setRelationshipBuilder) Status(status string) *setRelationshipBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the user-defined custom properties. Under the default PAM
// projection this replaces the entire payload; omitted keys are removed. An
// empty map is sent as {}. Not calling Payload, or passing nil, omits the
// field so the server resets it to default.
func (b *setRelationshipBuilder) Payload(payload map[string]interface{}) *setRelationshipBuilder {
	b.opts.Payload = payload
	b.opts.setPayload = true
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *setRelationshipBuilder) IfMatchETag(eTag string) *setRelationshipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setRelationshipBuilder) QueryParam(queryParam map[string]string) *setRelationshipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the setRelationship request.
func (b *setRelationshipBuilder) Transport(tr http.RoundTripper) *setRelationshipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *setRelationshipOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"ID":                       o.ID,
		"RelationshipClassVersion": o.RelationshipClassVersion,
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the setRelationship request.
func (b *setRelationshipBuilder) Execute() (*PNRelationshipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNSetRelationshipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetRelationshipResponse, status, err
	}

	return newPNRelationshipResponse(rawJSON, status, PNSetRelationshipOperation, b.opts.pubnub)
}

type setRelationshipOpts struct {
	endpointOpts

	ID                       string
	RelationshipClassVersion int
	Status                   string
	Payload                  map[string]interface{}
	setPayload               bool
	IfMatchETag              string
	setIfMatchETag           bool
	QueryParam               map[string]string

	Transport http.RoundTripper
}

func (o *setRelationshipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingRelationshipID)
	}
	if o.RelationshipClassVersion < 1 {
		return newValidationError(o, StrInvalidRelationshipClassVersion)
	}
	return nil
}

func (o *setRelationshipOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *setRelationshipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *setRelationshipOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := marshalDataSyncSetBody("relationshipClassVersion", o.RelationshipClassVersion, o.Status, o.Payload, o.setPayload)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "SetRelationshipSerializationFailed", PNSetRelationshipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setRelationshipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": relationshipContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *setRelationshipOpts) httpMethod() string {
	return "PUT"
}

func (o *setRelationshipOpts) operationType() OperationType {
	return PNSetRelationshipOperation
}
