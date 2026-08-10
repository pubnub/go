package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptyGetUserResponse *PNUserResponse

type getUserBuilder struct {
	opts *getUserOpts
}

func newGetUserBuilder(pubnub *PubNub) *getUserBuilder {
	return newGetUserBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetUserOpts(pubnub *PubNub, ctx Context) *getUserOpts {
	return &getUserOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetUserBuilderWithContext(pubnub *PubNub, context Context) *getUserBuilder {
	return &getUserBuilder{opts: newGetUserOpts(pubnub, context)}
}

// ID sets the required identifier of the user to read.
func (b *getUserBuilder) ID(id string) *getUserBuilder {
	b.opts.ID = id
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUserBuilder) QueryParam(queryParam map[string]string) *getUserBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getUser request.
func (b *getUserBuilder) Transport(tr http.RoundTripper) *getUserBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getUserOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the getUser request.
func (b *getUserBuilder) Execute() (*PNUserResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetDataSyncUserOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetUserResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNGetDataSyncUserOperation, b.opts.pubnub)
}

type getUserOpts struct {
	endpointOpts

	ID         string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingUserID)
	}
	return nil
}

func (o *getUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(usersIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getUserOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *getUserOpts) httpMethod() string {
	return "GET"
}

func (o *getUserOpts) operationType() OperationType {
	return PNGetDataSyncUserOperation
}
