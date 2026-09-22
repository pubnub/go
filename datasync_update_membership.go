package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateMembershipResponse *PNDataSyncMembershipResponse

type updateMembershipBuilder struct {
	opts *updateMembershipOpts
}

func newUpdateMembershipBuilder(pubnub *PubNub) *updateMembershipBuilder {
	return newUpdateMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateMembershipOpts(pubnub *PubNub, ctx Context) *updateMembershipOpts {
	return &updateMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateMembershipBuilderWithContext(pubnub *PubNub, context Context) *updateMembershipBuilder {
	return &updateMembershipBuilder{opts: newUpdateMembershipOpts(pubnub, context)}
}

// ID sets the required identifier of the membership to update.
func (b *updateMembershipBuilder) ID(id string) *updateMembershipBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *updateMembershipBuilder) Operations(operations []PNJSONPatchOperation) *updateMembershipBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *updateMembershipBuilder) Add(path string, value interface{}) *updateMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *updateMembershipBuilder) Remove(path string) *updateMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *updateMembershipBuilder) Replace(path string, value interface{}) *updateMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *updateMembershipBuilder) Move(from string, path string) *updateMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *updateMembershipBuilder) Copy(from string, path string) *updateMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *updateMembershipBuilder) Test(path string, value interface{}) *updateMembershipBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateMembershipBuilder) IfMatchETag(eTag string) *updateMembershipBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateMembershipBuilder) QueryParam(queryParam map[string]string) *updateMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateMembership request.
func (b *updateMembershipBuilder) Transport(tr http.RoundTripper) *updateMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateMembershipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the updateMembership request.
func (b *updateMembershipBuilder) Execute() (*PNDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateMembershipResponse, status, err
	}

	return newPNDataSyncMembershipResponse(rawJSON, status, PNUpdateDataSyncMembershipOperation, b.opts.pubnub)
}

type updateMembershipOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *updateMembershipOpts) validate() error {
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

func (o *updateMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateMembershipOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateMembershipSerializationFailed", PNUpdateDataSyncMembershipOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateMembershipOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateMembershipOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateMembershipOpts) operationType() OperationType {
	return PNUpdateDataSyncMembershipOperation
}
