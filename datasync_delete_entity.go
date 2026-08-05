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

var emptyRemoveEntityResponse *PNRemoveEntityResponse

type deleteEntityBuilder struct {
	opts *deleteEntityOpts
}

func newDeleteEntityBuilder(pubnub *PubNub) *deleteEntityBuilder {
	return newDeleteEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newDeleteEntityOpts(pubnub *PubNub, ctx Context) *deleteEntityOpts {
	return &deleteEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newDeleteEntityBuilderWithContext(pubnub *PubNub, context Context) *deleteEntityBuilder {
	return &deleteEntityBuilder{opts: newDeleteEntityOpts(pubnub, context)}
}

// ID sets the required identifier of the entity to remove.
func (b *deleteEntityBuilder) ID(id string) *deleteEntityBuilder {
	b.opts.ID = id
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *deleteEntityBuilder) IfMatchETag(eTag string) *deleteEntityBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteEntityBuilder) QueryParam(queryParam map[string]string) *deleteEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the removeEntity request.
func (b *deleteEntityBuilder) Transport(tr http.RoundTripper) *deleteEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *deleteEntityOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the removeEntity request.
func (b *deleteEntityBuilder) Execute() (*PNRemoveEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNRemoveEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveEntityResponse, status, err
	}

	return newPNRemoveEntityResponse(rawJSON, b.opts, status)
}

type deleteEntityOpts struct {
	endpointOpts

	ID             string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *deleteEntityOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingEntityID)
	}
	return nil
}

func (o *deleteEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *deleteEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteEntityOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *deleteEntityOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteEntityOpts) operationType() OperationType {
	return PNRemoveEntityOperation
}

// PNRemoveEntityResponse is the (empty) response returned by RemoveEntity. The
// server returns no body on success; the struct exists for API symmetry.
type PNRemoveEntityResponse struct {
	Status int `json:"status,omitempty"`
}

func newPNRemoveEntityResponse(jsonBytes []byte, o *deleteEntityOpts,
	status StatusResponse) (*PNRemoveEntityResponse, StatusResponse, error) {

	resp := &PNRemoveEntityResponse{}

	// The remove endpoint returns an empty body on success; only attempt to
	// unmarshal when the server actually sent a JSON payload.
	if len(bytes.TrimSpace(jsonBytes)) == 0 {
		return resp, status, nil
	}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "RemoveEntityResponseParsingFailed", PNRemoveEntityOperation, true)
		return emptyRemoveEntityResponse, status, e
	}

	return resp, status, nil
}
