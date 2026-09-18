package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v9/pnerr"
)

var emptyRemoveUserResponse *PNRemoveUserResponse

// PNRemoveUserResponse is the (empty) response returned by RemoveUser. The
// server returns no body on success; the struct exists for API symmetry.
type PNRemoveUserResponse = PNRemoveEntityResponse

type deleteUserBuilder struct {
	opts *deleteUserOpts
}

func newDeleteUserBuilder(pubnub *PubNub) *deleteUserBuilder {
	return newDeleteUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newDeleteUserOpts(pubnub *PubNub, ctx Context) *deleteUserOpts {
	return &deleteUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newDeleteUserBuilderWithContext(pubnub *PubNub, context Context) *deleteUserBuilder {
	return &deleteUserBuilder{opts: newDeleteUserOpts(pubnub, context)}
}

// ID sets the required identifier of the user to remove.
func (b *deleteUserBuilder) ID(id string) *deleteUserBuilder {
	b.opts.ID = id
	return b
}

// IfMatchETag sets the ETag for optimistic concurrency control via the If-Match header.
func (b *deleteUserBuilder) IfMatchETag(eTag string) *deleteUserBuilder {
	b.opts.IfMatchETag = eTag
	b.opts.setIfMatchETag = true
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteUserBuilder) QueryParam(queryParam map[string]string) *deleteUserBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the removeUser request.
func (b *deleteUserBuilder) Transport(tr http.RoundTripper) *deleteUserBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *deleteUserOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the removeUser request.
func (b *deleteUserBuilder) Execute() (*PNRemoveUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNRemoveDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveUserResponse, status, err
	}

	return newPNRemoveUserResponse(rawJSON, b.opts, status)
}

type deleteUserOpts struct {
	endpointOpts

	ID             string
	IfMatchETag    string
	setIfMatchETag bool
	QueryParam     map[string]string

	Transport http.RoundTripper
}

func (o *deleteUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingUserID)
	}
	return nil
}

func (o *deleteUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *deleteUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteUserOpts) buildHeaders() (map[string]string, error) {
	headers := make(map[string]string)
	if o.setIfMatchETag {
		headers["If-Match"] = o.IfMatchETag
	}
	return headers, nil
}

func (o *deleteUserOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteUserOpts) operationType() OperationType {
	return PNRemoveDataSyncUserOperation
}

func newPNRemoveUserResponse(jsonBytes []byte, o *deleteUserOpts,
	status StatusResponse) (*PNRemoveUserResponse, StatusResponse, error) {

	resp := &PNRemoveUserResponse{}

	// The remove endpoint returns an empty body on success; only attempt to
	// unmarshal when the server actually sent a JSON payload.
	if len(bytes.TrimSpace(jsonBytes)) == 0 {
		return resp, status, nil
	}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		o.pubnub.loggerManager.LogError(e, "RemoveUserResponseParsingFailed", PNRemoveDataSyncUserOperation, true)
		return emptyRemoveUserResponse, status, e
	}

	return resp, status, nil
}
