package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyPatchRelationshipResponse *PNRelationshipResponse

type patchRelationshipBuilder struct {
	opts *patchRelationshipOpts
}

func newPatchRelationshipBuilder(pubnub *PubNub) *patchRelationshipBuilder {
	return newPatchRelationshipBuilderWithContext(pubnub, pubnub.ctx)
}

func newPatchRelationshipOpts(pubnub *PubNub, ctx Context) *patchRelationshipOpts {
	return &patchRelationshipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newPatchRelationshipBuilderWithContext(pubnub *PubNub, context Context) *patchRelationshipBuilder {
	return &patchRelationshipBuilder{opts: newPatchRelationshipOpts(pubnub, context)}
}

// ID sets the required identifier of the relationship to patch.
func (b *patchRelationshipBuilder) ID(id string) *patchRelationshipBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *patchRelationshipBuilder) Operations(operations []PNJSONPatchOperation) *patchRelationshipBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *patchRelationshipBuilder) Add(path string, value interface{}) *patchRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *patchRelationshipBuilder) Remove(path string) *patchRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *patchRelationshipBuilder) Replace(path string, value interface{}) *patchRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *patchRelationshipBuilder) Move(from string, path string) *patchRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *patchRelationshipBuilder) Copy(from string, path string) *patchRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *patchRelationshipBuilder) Test(path string, value interface{}) *patchRelationshipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *patchRelationshipBuilder) IfMatchETag(eTag string) *patchRelationshipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *patchRelationshipBuilder) QueryParam(queryParam map[string]string) *patchRelationshipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the patchRelationship request.
func (b *patchRelationshipBuilder) Transport(tr http.RoundTripper) *patchRelationshipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *patchRelationshipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the patchRelationship request.
func (b *patchRelationshipBuilder) Execute() (*PNRelationshipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNPatchRelationshipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPatchRelationshipResponse, status, err
	}

	return newPNRelationshipResponse(rawJSON, status, PNPatchRelationshipOperation, b.opts.pubnub)
}

type patchRelationshipOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *patchRelationshipOpts) validate() error {
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

func (o *patchRelationshipOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *patchRelationshipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *patchRelationshipOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "PatchRelationshipSerializationFailed", PNPatchRelationshipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *patchRelationshipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *patchRelationshipOpts) httpMethod() string {
	return "PATCH"
}

func (o *patchRelationshipOpts) operationType() OperationType {
	return PNPatchRelationshipOperation
}
