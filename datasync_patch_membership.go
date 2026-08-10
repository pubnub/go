package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyPatchMembershipResponse *PNDataSyncMembershipResponse

type patchMembershipBuilder struct {
	opts *patchMembershipOpts
}

func newPatchMembershipBuilder(pubnub *PubNub) *patchMembershipBuilder {
	return newPatchMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newPatchMembershipOpts(pubnub *PubNub, ctx Context) *patchMembershipOpts {
	return &patchMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newPatchMembershipBuilderWithContext(pubnub *PubNub, context Context) *patchMembershipBuilder {
	return &patchMembershipBuilder{opts: newPatchMembershipOpts(pubnub, context)}
}

// ID sets the required identifier of the membership to patch.
func (b *patchMembershipBuilder) ID(id string) *patchMembershipBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *patchMembershipBuilder) Operations(operations []PNJSONPatchOperation) *patchMembershipBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *patchMembershipBuilder) Add(path string, value interface{}) *patchMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *patchMembershipBuilder) Remove(path string) *patchMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *patchMembershipBuilder) Replace(path string, value interface{}) *patchMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *patchMembershipBuilder) Move(from string, path string) *patchMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *patchMembershipBuilder) Copy(from string, path string) *patchMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *patchMembershipBuilder) Test(path string, value interface{}) *patchMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *patchMembershipBuilder) IfMatchETag(eTag string) *patchMembershipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *patchMembershipBuilder) QueryParam(queryParam map[string]string) *patchMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the patchMembership request.
func (b *patchMembershipBuilder) Transport(tr http.RoundTripper) *patchMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *patchMembershipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the patchMembership request.
func (b *patchMembershipBuilder) Execute() (*PNDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNPatchDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPatchMembershipResponse, status, err
	}

	return newPNDataSyncMembershipResponse(rawJSON, status, PNPatchDataSyncMembershipOperation, b.opts.pubnub)
}

type patchMembershipOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *patchMembershipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingMembershipID)
	}
	if len(o.Operations) == 0 {
		return newValidationError(o, StrMissingPatchOperations)
	}
	return nil
}

func (o *patchMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *patchMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *patchMembershipOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "PatchMembershipSerializationFailed", PNPatchDataSyncMembershipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *patchMembershipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *patchMembershipOpts) httpMethod() string {
	return "PATCH"
}

func (o *patchMembershipOpts) operationType() OperationType {
	return PNPatchDataSyncMembershipOperation
}
