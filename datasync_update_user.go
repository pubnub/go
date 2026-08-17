package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var emptyUpdateUserResponse *PNUserResponse

type updateUserBuilder struct {
	opts *updateUserOpts
}

func newUpdateUserBuilder(pubnub *PubNub) *updateUserBuilder {
	return newUpdateUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newUpdateUserOpts(pubnub *PubNub, ctx Context) *updateUserOpts {
	return &updateUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newUpdateUserBuilderWithContext(pubnub *PubNub, context Context) *updateUserBuilder {
	return &updateUserBuilder{opts: newUpdateUserOpts(pubnub, context)}
}

// ID sets the required identifier of the user to update.
func (b *updateUserBuilder) ID(id string) *updateUserBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *updateUserBuilder) Operations(operations []PNJSONPatchOperation) *updateUserBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *updateUserBuilder) Add(path string, value interface{}) *updateUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *updateUserBuilder) Remove(path string) *updateUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *updateUserBuilder) Replace(path string, value interface{}) *updateUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *updateUserBuilder) Move(from string, path string) *updateUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *updateUserBuilder) Copy(from string, path string) *updateUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *updateUserBuilder) Test(path string, value interface{}) *updateUserBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *updateUserBuilder) IfMatchETag(eTag string) *updateUserBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateUserBuilder) QueryParam(queryParam map[string]string) *updateUserBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the updateUser request.
func (b *updateUserBuilder) Transport(tr http.RoundTripper) *updateUserBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *updateUserOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the updateUser request.
func (b *updateUserBuilder) Execute() (*PNUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNUpdateDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyUpdateUserResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNUpdateDataSyncUserOperation, b.opts.pubnub)
}

type updateUserOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *updateUserOpts) validate() error {
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

func (o *updateUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *updateUserOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "UpdateUserSerializationFailed", PNUpdateDataSyncUserOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *updateUserOpts) buildHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": entityPatchContentType,
	}
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *updateUserOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateUserOpts) operationType() OperationType {
	return PNUpdateDataSyncUserOperation
}
