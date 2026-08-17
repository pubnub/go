package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateRelationshipResponse *PNRelationshipResponse

type updateRelationshipBuilder struct {
	opts *updateRelationshipOpts
}

func newUpdateRelationshipBuilder(pubnub *PubNub) *updateRelationshipBuilder {
	return newUpdateRelationshipBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateRelationshipOpts(pubnub *PubNub, ctx Context) *updateRelationshipOpts {
	return &updateRelationshipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateRelationshipBuilderWithContext(pubnub *PubNub, context Context) *updateRelationshipBuilder {
	return &updateRelationshipBuilder{opts: newUpdateRelationshipOpts(pubnub, context)}
}

// ID sets the required identifier of the relationship to update.
func (b *updateRelationshipBuilder) ID(id string) *updateRelationshipBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *updateRelationshipBuilder) Operations(operations []PNJSONPatchOperation) *updateRelationshipBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *updateRelationshipBuilder) Add(path string, value interface{}) *updateRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *updateRelationshipBuilder) Remove(path string) *updateRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *updateRelationshipBuilder) Replace(path string, value interface{}) *updateRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *updateRelationshipBuilder) Move(from string, path string) *updateRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *updateRelationshipBuilder) Copy(from string, path string) *updateRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *updateRelationshipBuilder) Test(path string, value interface{}) *updateRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateRelationshipBuilder) IfMatchETag(eTag string) *updateRelationshipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateRelationshipBuilder) QueryParam(queryParam map[string]string) *updateRelationshipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateRelationship request.
func (b *updateRelationshipBuilder) Transport(tr http.RoundTripper) *updateRelationshipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateRelationshipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the updateRelationship request.
func (b *updateRelationshipBuilder) Execute() (*PNRelationshipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateRelationshipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateRelationshipResponse, status, err
	}

	return newPNRelationshipResponse(rawJSON, status, PNUpdateRelationshipOperation, b.opts.pubnub)
}

type updateRelationshipOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *updateRelationshipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingRelationshipID)
	}
	if len(o.Operations) == 0 {
		return newValidationError(o, StrMissingPatchOperations)
	}
	return nil
}

func (o *updateRelationshipOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateRelationshipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateRelationshipOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateRelationshipSerializationFailed", PNUpdateRelationshipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateRelationshipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateRelationshipOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateRelationshipOpts) operationType() OperationType {
	return PNUpdateRelationshipOperation
}
