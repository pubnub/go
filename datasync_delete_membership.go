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

var emptyRemoveMembershipResponse *PNRemoveDataSyncMembershipResponse

// PNRemoveDataSyncMembershipResponse is the (empty) response returned by
// RemoveMembership. The server returns no body on success; the struct exists
// for API symmetry. Named to avoid colliding with Objects RemoveMemberships.
type PNRemoveDataSyncMembershipResponse struct {
	Status int `json:"status,omitempty"`
}

type deleteMembershipBuilder struct {
	opts *deleteMembershipOpts
}

func newDeleteMembershipBuilder(pubnub *PubNub) *deleteMembershipBuilder {
	return newDeleteMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newDeleteMembershipOpts(pubnub *PubNub, ctx Context) *deleteMembershipOpts {
	return &deleteMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newDeleteMembershipBuilderWithContext(pubnub *PubNub, context Context) *deleteMembershipBuilder {
	return &deleteMembershipBuilder{opts: newDeleteMembershipOpts(pubnub, context)}
}

// ID sets the required identifier of the membership to remove.
func (b *deleteMembershipBuilder) ID(id string) *deleteMembershipBuilder {
	b.opts.ID = id
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *deleteMembershipBuilder) IfMatchETag(eTag string) *deleteMembershipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteMembershipBuilder) QueryParam(queryParam map[string]string) *deleteMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the removeMembership request.
func (b *deleteMembershipBuilder) Transport(tr http.RoundTripper) *deleteMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *deleteMembershipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the removeMembership request.
func (b *deleteMembershipBuilder) Execute() (*PNRemoveDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNRemoveDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveMembershipResponse, status, err
	}

	return newPNRemoveDataSyncMembershipResponse(rawJSON, b.opts, status)
}

type deleteMembershipOpts struct {
	endpointOpts

	ID             string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *deleteMembershipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingMembershipID)
	}
	return nil
}

func (o *deleteMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *deleteMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteMembershipOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *deleteMembershipOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteMembershipOpts) operationType() OperationType {
	return PNRemoveDataSyncMembershipOperation
}

func newPNRemoveDataSyncMembershipResponse(jsonBytes []byte, o *deleteMembershipOpts,
	status StatusResponse) (*PNRemoveDataSyncMembershipResponse, StatusResponse, error) {

	resp := &PNRemoveDataSyncMembershipResponse{}

	// The remove endpoint returns an empty body on success; only attempt to
	// unmarshal when the server actually sent a JSON payload.
	if len(bytes.TrimSpace(jsonBytes)) == 0 {
		return resp, status, nil
	}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "RemoveMembershipResponseParsingFailed", PNRemoveDataSyncMembershipOperation, true)
		return emptyRemoveMembershipResponse, status, e
	}

	return resp, status, nil
}
