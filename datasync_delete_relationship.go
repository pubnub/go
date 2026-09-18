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

var emptyRemoveRelationshipResponse *PNRemoveRelationshipResponse

type deleteRelationshipBuilder struct {
	opts *deleteRelationshipOpts
}

func newDeleteRelationshipBuilder(pubnub *PubNub) *deleteRelationshipBuilder {
	return newDeleteRelationshipBuilderWithContext(pubnub, pubnub.ctx)
}

func newDeleteRelationshipOpts(pubnub *PubNub, ctx Context) *deleteRelationshipOpts {
	return &deleteRelationshipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newDeleteRelationshipBuilderWithContext(pubnub *PubNub, context Context) *deleteRelationshipBuilder {
	return &deleteRelationshipBuilder{opts: newDeleteRelationshipOpts(pubnub, context)}
}

// ID sets the required identifier of the relationship to remove.
func (b *deleteRelationshipBuilder) ID(id string) *deleteRelationshipBuilder {
	b.opts.ID = id
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *deleteRelationshipBuilder) IfMatchETag(eTag string) *deleteRelationshipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteRelationshipBuilder) QueryParam(queryParam map[string]string) *deleteRelationshipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the removeRelationship request.
func (b *deleteRelationshipBuilder) Transport(tr http.RoundTripper) *deleteRelationshipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *deleteRelationshipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the removeRelationship request.
func (b *deleteRelationshipBuilder) Execute() (*PNRemoveRelationshipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNRemoveRelationshipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveRelationshipResponse, status, err
	}

	return newPNRemoveRelationshipResponse(rawJSON, b.opts, status)
}

type deleteRelationshipOpts struct {
	endpointOpts

	ID             string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *deleteRelationshipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingRelationshipID)
	}
	return nil
}

func (o *deleteRelationshipOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *deleteRelationshipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteRelationshipOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *deleteRelationshipOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteRelationshipOpts) operationType() OperationType {
	return PNRemoveRelationshipOperation
}

// PNRemoveRelationshipResponse is the (empty) response returned by
// RemoveRelationship. The server returns no body on success; the struct exists
// for API symmetry.
type PNRemoveRelationshipResponse struct {
	Status int `json:"status,omitempty"`
}

func newPNRemoveRelationshipResponse(jsonBytes []byte, o *deleteRelationshipOpts,
	status StatusResponse) (*PNRemoveRelationshipResponse, StatusResponse, error) {

	resp := &PNRemoveRelationshipResponse{}

	// The remove endpoint returns an empty body on success; only attempt to
	// unmarshal when the server actually sent a JSON payload.
	if len(bytes.TrimSpace(jsonBytes)) == 0 {
		return resp, status, nil
	}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "RemoveRelationshipResponseParsingFailed", PNRemoveRelationshipOperation, true)
		return emptyRemoveRelationshipResponse, status, e
	}

	return resp, status, nil
}
