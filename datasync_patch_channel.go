package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyPatchChannelResponse *PNDataSyncChannelResponse

type patchChannelBuilder struct {
	opts *patchChannelOpts
}

func newPatchChannelBuilder(pubnub *PubNub) *patchChannelBuilder {
	return newPatchChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newPatchChannelOpts(pubnub *PubNub, ctx Context) *patchChannelOpts {
	return &patchChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newPatchChannelBuilderWithContext(pubnub *PubNub, context Context) *patchChannelBuilder {
	return &patchChannelBuilder{opts: newPatchChannelOpts(pubnub, context)}
}

// ID sets the required identifier of the channel to patch.
func (b *patchChannelBuilder) ID(id string) *patchChannelBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *patchChannelBuilder) Operations(operations []PNJSONPatchOperation) *patchChannelBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *patchChannelBuilder) Add(path string, value interface{}) *patchChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *patchChannelBuilder) Remove(path string) *patchChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *patchChannelBuilder) Replace(path string, value interface{}) *patchChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *patchChannelBuilder) Move(from string, path string) *patchChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *patchChannelBuilder) Copy(from string, path string) *patchChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *patchChannelBuilder) Test(path string, value interface{}) *patchChannelBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *patchChannelBuilder) IfMatchETag(eTag string) *patchChannelBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *patchChannelBuilder) QueryParam(queryParam map[string]string) *patchChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the patchChannel request.
func (b *patchChannelBuilder) Transport(tr http.RoundTripper) *patchChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *patchChannelOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the patchChannel request.
func (b *patchChannelBuilder) Execute() (*PNDataSyncChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNPatchDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPatchChannelResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNPatchDataSyncChannelOperation, b.opts.pubnub)
}

type patchChannelOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *patchChannelOpts) validate() error {
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

func (o *patchChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *patchChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *patchChannelOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "PatchChannelSerializationFailed", PNPatchDataSyncChannelOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *patchChannelOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *patchChannelOpts) httpMethod() string {
	return "PATCH"
}

func (o *patchChannelOpts) operationType() OperationType {
	return PNPatchDataSyncChannelOperation
}
