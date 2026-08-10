package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateMembershipResponse *PNDataSyncMembershipResponse

type updateMembershipBuilder struct {
	opts *updateMembershipOpts
}

func newUpdateMembershipBuilder(pubnub *PubNub) *updateMembershipBuilder {
	return newUpdateMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateMembershipOpts(pubnub *PubNub, ctx Context) *updateMembershipOpts {
	return &updateMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateMembershipBuilderWithContext(pubnub *PubNub, context Context) *updateMembershipBuilder {
	return &updateMembershipBuilder{opts: newUpdateMembershipOpts(pubnub, context)}
}

// updateMembershipBody is the request envelope sent on PUT /memberships/{id}.
// channelId, userId and relationshipClass are intentionally omitted because
// they are immutable after creation.
type updateMembershipBody struct {
	Data updateMembershipBodyData `json:"data"`
}

type updateMembershipBodyData struct {
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Status                   string                 `json:"status,omitempty"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the required identifier of the membership to replace.
func (b *updateMembershipBuilder) ID(id string) *updateMembershipBuilder {
	b.opts.ID = id
	return b
}

// RelationshipClassVersion sets the required schema version of the Membership
// relationship class (>= 1).
func (b *updateMembershipBuilder) RelationshipClassVersion(version int) *updateMembershipBuilder {
	b.opts.RelationshipClassVersion = version
	return b
}

// Status sets the optional membership status (e.g. "active").
func (b *updateMembershipBuilder) Status(status string) *updateMembershipBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom membership fields. Under the default PAM
// projection this replaces the entire payload; omitted fields are removed.
func (b *updateMembershipBuilder) Payload(payload map[string]interface{}) *updateMembershipBuilder {
	b.opts.Payload = payload
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateMembershipBuilder) IfMatchETag(eTag string) *updateMembershipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateMembershipBuilder) QueryParam(queryParam map[string]string) *updateMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateMembership request.
func (b *updateMembershipBuilder) Transport(tr http.RoundTripper) *updateMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateMembershipOpts) GetLogParams() map[string]interface{} {
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

// Execute runs the updateMembership request.
func (b *updateMembershipBuilder) Execute() (*PNDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateMembershipResponse, status, err
	}

	return newPNDataSyncMembershipResponse(rawJSON, status, PNUpdateDataSyncMembershipOperation, b.opts.pubnub)
}

type updateMembershipOpts struct {
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

func (o *updateMembershipOpts) validate() error {
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

func (o *updateMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateMembershipOpts) buildBody() ([]byte, error) {
	b := &updateMembershipBody{
		Data: updateMembershipBodyData{
			RelationshipClassVersion: o.RelationshipClassVersion,
			Status:                   o.Status,
			Payload:                  o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateMembershipSerializationFailed", PNUpdateDataSyncMembershipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateMembershipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": membershipContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateMembershipOpts) httpMethod() string {
	return "PUT"
}

func (o *updateMembershipOpts) operationType() OperationType {
	return PNUpdateDataSyncMembershipOperation
}
