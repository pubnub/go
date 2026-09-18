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

var emptyRemoveChannelResponse *PNRemoveChannelResponse

// PNRemoveChannelResponse is the (empty) response returned by RemoveChannel. The
// server returns no body on success; the struct exists for API symmetry.
type PNRemoveChannelResponse = PNRemoveEntityResponse

type deleteChannelBuilder struct {
	opts *deleteChannelOpts
}

func newDeleteChannelBuilder(pubnub *PubNub) *deleteChannelBuilder {
	return newDeleteChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newDeleteChannelOpts(pubnub *PubNub, ctx Context) *deleteChannelOpts {
	return &deleteChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newDeleteChannelBuilderWithContext(pubnub *PubNub, context Context) *deleteChannelBuilder {
	return &deleteChannelBuilder{opts: newDeleteChannelOpts(pubnub, context)}
}

// ID sets the required identifier of the channel to remove.
func (b *deleteChannelBuilder) ID(id string) *deleteChannelBuilder {
	b.opts.ID = id
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *deleteChannelBuilder) IfMatchETag(eTag string) *deleteChannelBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteChannelBuilder) QueryParam(queryParam map[string]string) *deleteChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the removeChannel request.
func (b *deleteChannelBuilder) Transport(tr http.RoundTripper) *deleteChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *deleteChannelOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the removeChannel request.
func (b *deleteChannelBuilder) Execute() (*PNRemoveChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNRemoveDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelResponse, status, err
	}

	return newPNRemoveChannelResponse(rawJSON, b.opts, status)
}

type deleteChannelOpts struct {
	endpointOpts

	ID             string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *deleteChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingChannelID)
	}
	return nil
}

func (o *deleteChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *deleteChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteChannelOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *deleteChannelOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteChannelOpts) operationType() OperationType {
	return PNRemoveDataSyncChannelOperation
}

func newPNRemoveChannelResponse(jsonBytes []byte, o *deleteChannelOpts,
	status StatusResponse) (*PNRemoveChannelResponse, StatusResponse, error) {

	resp := &PNRemoveChannelResponse{}

	// The remove endpoint returns an empty body on success; only attempt to
	// unmarshal when the server actually sent a JSON payload.
	if len(bytes.TrimSpace(jsonBytes)) == 0 {
		return resp, status, nil
	}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "RemoveChannelResponseParsingFailed", PNRemoveDataSyncChannelOperation, true)
		return emptyRemoveChannelResponse, status, e
	}

	return resp, status, nil
}
