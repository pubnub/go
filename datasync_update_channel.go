package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateChannelResponse *PNDataSyncChannelResponse

type updateChannelBuilder struct {
	opts *updateChannelOpts
}

func newUpdateChannelBuilder(pubnub *PubNub) *updateChannelBuilder {
	return newUpdateChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateChannelOpts(pubnub *PubNub, ctx Context) *updateChannelOpts {
	return &updateChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateChannelBuilderWithContext(pubnub *PubNub, context Context) *updateChannelBuilder {
	return &updateChannelBuilder{opts: newUpdateChannelOpts(pubnub, context)}
}

// ID sets the required identifier of the channel to update.
func (b *updateChannelBuilder) ID(id string) *updateChannelBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *updateChannelBuilder) Operations(operations []PNJSONPatchOperation) *updateChannelBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *updateChannelBuilder) Add(path string, value interface{}) *updateChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *updateChannelBuilder) Remove(path string) *updateChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *updateChannelBuilder) Replace(path string, value interface{}) *updateChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *updateChannelBuilder) Move(from string, path string) *updateChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *updateChannelBuilder) Copy(from string, path string) *updateChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *updateChannelBuilder) Test(path string, value interface{}) *updateChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateChannelBuilder) IfMatchETag(eTag string) *updateChannelBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateChannelBuilder) QueryParam(queryParam map[string]string) *updateChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateChannel request.
func (b *updateChannelBuilder) Transport(tr http.RoundTripper) *updateChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateChannelOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the updateChannel request.
func (b *updateChannelBuilder) Execute() (*PNDataSyncChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateChannelResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNUpdateDataSyncChannelOperation, b.opts.pubnub)
}

type updateChannelOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *updateChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingChannelID)
	}
	if len(o.Operations) == 0 {
		return newValidationError(o, StrMissingPatchOperations)
	}
	return nil
}

func (o *updateChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateChannelOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateChannelSerializationFailed", PNUpdateDataSyncChannelOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateChannelOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateChannelOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateChannelOpts) operationType() OperationType {
	return PNUpdateDataSyncChannelOperation
}
