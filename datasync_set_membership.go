package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptySetMembershipResponse *PNDataSyncMembershipResponse

type setMembershipBuilder struct {
	opts *setMembershipOpts
}

func newSetMembershipBuilder(pubnub *PubNub) *setMembershipBuilder {
	return newSetMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetMembershipOpts(pubnub *PubNub, ctx Context) *setMembershipOpts {
	return &setMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newSetMembershipBuilderWithContext(pubnub *PubNub, context Context) *setMembershipBuilder {
	return &setMembershipBuilder{opts: newSetMembershipOpts(pubnub, context)}
}

// setMembershipBody is the request envelope sent on PUT /memberships/{id}.
// channelId, userId and relationshipClass are intentionally omitted because
// they are immutable after creation.
type setMembershipBody struct {
	Data setMembershipBodyData `json:"data"`
}

type setMembershipBodyData struct {
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Status                   string                 `json:"status,omitempty"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the membership to replace.
func (b *setMembershipBuilder) ID(id string) *setMembershipBuilder {
	b.opts.ID = id
	return b
}

// RelationshipClassVersion sets the required schema version of the Membership
// relationship class (>= 1).
func (b *setMembershipBuilder) RelationshipClassVersion(version int) *setMembershipBuilder {
	b.opts.RelationshipClassVersion = version
	return b
}

// Status sets the optional membership status (e.g. "active").
func (b *setMembershipBuilder) Status(status string) *setMembershipBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom membership fields. Under the default PAM
// projection this replaces the entire payload; omitted fields are removed.
func (b *setMembershipBuilder) Payload(payload map[string]interface{}) *setMembershipBuilder {
	b.opts.Payload = payload
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *setMembershipBuilder) IfMatchETag(eTag string) *setMembershipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setMembershipBuilder) QueryParam(queryParam map[string]string) *setMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the setMembership request.
func (b *setMembershipBuilder) Transport(tr http.RoundTripper) *setMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *setMembershipOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the setMembership request.
func (b *setMembershipBuilder) Execute() (*PNDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNSetDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetMembershipResponse, status, err
	}

	return newPNDataSyncMembershipResponse(rawJSON, status, PNSetDataSyncMembershipOperation, b.opts.pubnub)
}

type setMembershipOpts struct {
	endpointOpts

	ID                       string
	RelationshipClassVersion int
	Status                   string
	Payload                  map[string]interface{}
	IfMatchETag              string
	setIfMatchETag           bool
	QueryParam               map[string]string

	Transport http.RoundTripper
}

func (o *setMembershipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingMembershipID)
	}
	if o.RelationshipClassVersion < 1 {
		return newValidationError(o, StrInvalidRelationshipClassVersion)
	}
	return nil
}

func (o *setMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *setMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *setMembershipOpts) buildBody() ([]byte, error) {
	b := &setMembershipBody{
		Data: setMembershipBodyData{
			RelationshipClassVersion: o.RelationshipClassVersion,
			Status:                   o.Status,
			Payload:                  o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "SetMembershipSerializationFailed", PNSetDataSyncMembershipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setMembershipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": membershipContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *setMembershipOpts) httpMethod() string {
	return "PUT"
}

func (o *setMembershipOpts) operationType() OperationType {
	return PNSetDataSyncMembershipOperation
}
