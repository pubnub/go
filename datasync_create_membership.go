package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v9/pnerr"
)

var emptyCreateMembershipResponse *PNDataSyncMembershipResponse

type createMembershipBuilder struct {
	opts *createMembershipOpts
}

func newCreateMembershipBuilder(pubnub *PubNub) *createMembershipBuilder {
	return newCreateMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newCreateMembershipOpts(pubnub *PubNub, ctx Context) *createMembershipOpts {
	return &createMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newCreateMembershipBuilderWithContext(pubnub *PubNub, context Context) *createMembershipBuilder {
	return &createMembershipBuilder{opts: newCreateMembershipOpts(pubnub, context)}
}

// createMembershipBody is the request envelope sent on POST /memberships.
type createMembershipBody struct {
	Data createMembershipBodyData `json:"data"`
}

type createMembershipBodyData struct {
	ID                       string                 `json:"id,omitempty"`
	ChannelID                string                 `json:"channelId"`
	UserID                   string                 `json:"userId"`
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Status                   string                 `json:"status,omitempty"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
}

// ID sets the optional membership identifier. When omitted the server generates one.
func (b *createMembershipBuilder) ID(id string) *createMembershipBuilder {
	b.opts.ID = id
	return b
}

// ChannelID sets the required linked channel identifier. Immutable after creation.
func (b *createMembershipBuilder) ChannelID(channelID string) *createMembershipBuilder {
	b.opts.ChannelID = channelID
	return b
}

// UserID sets the required linked user identifier. Immutable after creation.
func (b *createMembershipBuilder) UserID(userID string) *createMembershipBuilder {
	b.opts.UserID = userID
	return b
}

// RelationshipClassVersion sets the required schema version of the Membership
// relationship class (>= 1).
func (b *createMembershipBuilder) RelationshipClassVersion(version int) *createMembershipBuilder {
	b.opts.RelationshipClassVersion = version
	return b
}

// Status sets the optional membership status (e.g. "active").
func (b *createMembershipBuilder) Status(status string) *createMembershipBuilder {
	b.opts.Status = status
	return b
}

// Payload sets the optional custom membership fields.
func (b *createMembershipBuilder) Payload(payload map[string]interface{}) *createMembershipBuilder {
	b.opts.Payload = payload
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *createMembershipBuilder) QueryParam(queryParam map[string]string) *createMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the createMembership request.
func (b *createMembershipBuilder) Transport(tr http.RoundTripper) *createMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *createMembershipOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"ChannelID":                o.ChannelID,
		"UserID":                   o.UserID,
		"RelationshipClassVersion": o.RelationshipClassVersion,
	}
	if o.ID != "" {
		params["ID"] = o.ID
	}
	if o.Status != "" {
		params["Status"] = o.Status
	}
	if o.Payload != nil {
		params["Payload"] = fmt.Sprintf("%v", o.Payload)
	}
	return params
}

// Execute runs the createMembership request.
func (b *createMembershipBuilder) Execute() (*PNDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNCreateDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyCreateMembershipResponse, status, err
	}

	return newPNDataSyncMembershipResponse(rawJSON, status, PNCreateDataSyncMembershipOperation, b.opts.pubnub)
}

type createMembershipOpts struct {
	endpointOpts

	ID                       string
	ChannelID                string
	UserID                   string
	RelationshipClassVersion int
	Status                   string
	Payload                  map[string]interface{}
	QueryParam               map[string]string

	Transport http.RoundTripper
}

func (o *createMembershipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ChannelID == "" {
		return newValidationError(o, StrMissingChannelID)
	}
	if o.UserID == "" {
		return newValidationError(o, StrMissingUserID)
	}
	if o.RelationshipClassVersion < 1 {
		return newValidationError(o, StrInvalidRelationshipClassVersion)
	}
	return nil
}

func (o *createMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsPath, o.pubnub.Config.SubscribeKey), nil
}

func (o *createMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *createMembershipOpts) buildBody() ([]byte, error) {
	b := &createMembershipBody{
		Data: createMembershipBodyData{
			ID:                       o.ID,
			ChannelID:                o.ChannelID,
			UserID:                   o.UserID,
			RelationshipClassVersion: o.RelationshipClassVersion,
			Status:                   o.Status,
			Payload:                  o.Payload,
		},
	}

	jsonEncBytes, errEnc := json.Marshal(b)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "CreateMembershipSerializationFailed", PNCreateDataSyncMembershipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *createMembershipOpts) buildHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": membershipContentType,
	}, nil
}

func (o *createMembershipOpts) httpMethod() string {
	return "POST"
}

func (o *createMembershipOpts) operationType() OperationType {
	return PNCreateDataSyncMembershipOperation
}

// newPNDataSyncMembershipResponse unmarshals a single-membership envelope shared
// by the create, get, update and patch operations.
func newPNDataSyncMembershipResponse(jsonBytes []byte, status StatusResponse,
	operation OperationType, pn *PubNub) (*PNDataSyncMembershipResponse, StatusResponse, error) {

	resp := &PNDataSyncMembershipResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		pn.loggerManager.LogError(e, "MembershipResponseParsingFailed", operation, true)
		return nil, status, e
	}

	return resp, status, nil
}
