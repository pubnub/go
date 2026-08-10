package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyPatchUserResponse *PNUserResponse

type patchUserBuilder struct {
	opts *patchUserOpts
}

func newPatchUserBuilder(pubnub *PubNub) *patchUserBuilder {
	return newPatchUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newPatchUserOpts(pubnub *PubNub, ctx Context) *patchUserOpts {
	return &patchUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newPatchUserBuilderWithContext(pubnub *PubNub, context Context) *patchUserBuilder {
	return &patchUserBuilder{opts: newPatchUserOpts(pubnub, context)}
}

// ID sets the required identifier of the user to patch.
func (b *patchUserBuilder) ID(id string) *patchUserBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *patchUserBuilder) Operations(operations []PNJSONPatchOperation) *patchUserBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *patchUserBuilder) Add(path string, value interface{}) *patchUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *patchUserBuilder) Remove(path string) *patchUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *patchUserBuilder) Replace(path string, value interface{}) *patchUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *patchUserBuilder) Move(from string, path string) *patchUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *patchUserBuilder) Copy(from string, path string) *patchUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *patchUserBuilder) Test(path string, value interface{}) *patchUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *patchUserBuilder) IfMatchETag(eTag string) *patchUserBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *patchUserBuilder) QueryParam(queryParam map[string]string) *patchUserBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the patchUser request.
func (b *patchUserBuilder) Transport(tr http.RoundTripper) *patchUserBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *patchUserOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the patchUser request.
func (b *patchUserBuilder) Execute() (*PNUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNPatchDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPatchUserResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNPatchDataSyncUserOperation, b.opts.pubnub)
}

type patchUserOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *patchUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingUserID)
	}
	if len(o.Operations) == 0 {
		return newValidationError(o, StrMissingPatchOperations)
	}
	return nil
}

func (o *patchUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *patchUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *patchUserOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "PatchUserSerializationFailed", PNPatchDataSyncUserOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *patchUserOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *patchUserOpts) httpMethod() string {
	return "PATCH"
}

func (o *patchUserOpts) operationType() OperationType {
	return PNPatchDataSyncUserOperation
}
