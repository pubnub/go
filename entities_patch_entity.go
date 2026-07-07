package pubnub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

var emptyPatchEntityResponse *PNEntityResponse

type patchEntityBuilder struct {
	opts *patchEntityOpts
}

func newPatchEntityBuilder(pubnub *PubNub) *patchEntityBuilder {
	return newPatchEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newPatchEntityOpts(pubnub *PubNub, ctx Context) *patchEntityOpts {
	return &patchEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newPatchEntityBuilderWithContext(pubnub *PubNub, context Context) *patchEntityBuilder {
	return &patchEntityBuilder{opts: newPatchEntityOpts(pubnub, context)}
}

// ID sets the required identifier of the entity to patch.
func (b *patchEntityBuilder) ID(id string) *patchEntityBuilder {
	b.opts.ID = id
	return b
}

// Operations sets the full list of RFC 6902 JSON Patch operations to apply.
func (b *patchEntityBuilder) Operations(operations []PNJSONPatchOperation) *patchEntityBuilder {
	b.opts.Operations = operations
	return b
}

// Add appends an "add" JSON Patch operation.
func (b *patchEntityBuilder) Add(path string, value interface{}) *patchEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "add", Path: path, Value: value})
	return b
}

// Remove appends a "remove" JSON Patch operation.
func (b *patchEntityBuilder) Remove(path string) *patchEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "remove", Path: path})
	return b
}

// Replace appends a "replace" JSON Patch operation.
func (b *patchEntityBuilder) Replace(path string, value interface{}) *patchEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "replace", Path: path, Value: value})
	return b
}

// Move appends a "move" JSON Patch operation.
func (b *patchEntityBuilder) Move(from string, path string) *patchEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "move", From: from, Path: path})
	return b
}

// Copy appends a "copy" JSON Patch operation.
func (b *patchEntityBuilder) Copy(from string, path string) *patchEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "copy", From: from, Path: path})
	return b
}

// Test appends a "test" JSON Patch operation.
func (b *patchEntityBuilder) Test(path string, value interface{}) *patchEntityBuilder {
	b.opts.Operations = append(b.opts.Operations, PNJSONPatchOperation{Op: "test", Path: path, Value: value})
	return b
}

// IfMatch sets the ETag for optimistic concurrency control via the If-Match header.
func (b *patchEntityBuilder) IfMatch(eTag string) *patchEntityBuilder {
	b.opts.IfMatch = eTag
	b.opts.setIfMatch = true
	return b
}

// IdempotencyKey sets the idempotency key (UUIDv4). If left unset a random one
// is generated automatically, since the server requires it for PATCH requests.
func (b *patchEntityBuilder) IdempotencyKey(key string) *patchEntityBuilder {
	b.opts.IdempotencyKey = key
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *patchEntityBuilder) QueryParam(queryParam map[string]string) *patchEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the patchEntity request.
func (b *patchEntityBuilder) Transport(tr http.RoundTripper) *patchEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *patchEntityOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID":         o.ID,
		"Operations": fmt.Sprintf("%v", o.Operations),
	}
}

// Execute runs the patchEntity request.
func (b *patchEntityBuilder) Execute() (*PNEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNPatchEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPatchEntityResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNPatchEntityOperation, b.opts.pubnub)
}

type patchEntityOpts struct {
	endpointOpts

	ID             string
	Operations     []PNJSONPatchOperation
	IfMatch        string
	setIfMatch     bool
	IdempotencyKey string
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *patchEntityOpts) validate() error {
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

func (o *patchEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *patchEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *patchEntityOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Operations)
	if errEnc != nil {
		o.pubnub.loggerManager.LogError(errEnc, "PatchEntitySerializationFailed", PNPatchEntityOperation, true)
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *patchEntityOpts) buildHeaders() (map[string]string, error) {
	if o.IdempotencyKey == "" {
		o.IdempotencyKey = uuid.NewString()
	}
	headers := map[string]string{
		"Content-Type":    entityPatchContentType,
		"Idempotency-Key": o.IdempotencyKey,
	}
	if o.setIfMatch {
		headers["If-Match"] = o.IfMatch
	}
	return headers, nil
}

func (o *patchEntityOpts) httpMethod() string {
	return "PATCH"
}

func (o *patchEntityOpts) operationType() OperationType {
	return PNPatchEntityOperation
}
