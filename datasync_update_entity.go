package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateEntityResponse *PNEntityResponse

type updateEntityBuilder struct {
	opts *updateEntityOpts
}

func newUpdateEntityBuilder(pubnub *PubNub) *updateEntityBuilder {
	return newUpdateEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateEntityOpts(pubnub *PubNub, ctx Context) *updateEntityOpts {
	return &updateEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateEntityBuilderWithContext(pubnub *PubNub, context Context) *updateEntityBuilder {
	return &updateEntityBuilder{opts: newUpdateEntityOpts(pubnub, context)}
}

// ID sets the required identifier of the entity to update.
func (b *updateEntityBuilder) ID(id string) *updateEntityBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *updateEntityBuilder) Operations(operations []PNJSONPatchOperation) *updateEntityBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *updateEntityBuilder) Add(path string, value interface{}) *updateEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *updateEntityBuilder) Remove(path string) *updateEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *updateEntityBuilder) Replace(path string, value interface{}) *updateEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *updateEntityBuilder) Move(from string, path string) *updateEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *updateEntityBuilder) Copy(from string, path string) *updateEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *updateEntityBuilder) Test(path string, value interface{}) *updateEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateEntityBuilder) IfMatchETag(eTag string) *updateEntityBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateEntityBuilder) QueryParam(queryParam map[string]string) *updateEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateEntity request.
func (b *updateEntityBuilder) Transport(tr http.RoundTripper) *updateEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateEntityOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the updateEntity request.
func (b *updateEntityBuilder) Execute() (*PNEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateEntityResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNUpdateEntityOperation, b.opts.pubnub)
}

type updateEntityOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *updateEntityOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingEntityID)
	}
	if len(o.Operations) == 0 {
		return newValidationError(o, StrMissingPatchOperations)
	}
	return nil
}

func (o *updateEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateEntityOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateEntitySerializationFailed", PNUpdateEntityOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateEntityOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateEntityOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateEntityOpts) operationType() OperationType {
	return PNUpdateEntityOperation
}
