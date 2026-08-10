package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptyGetMembershipResponse *PNDataSyncMembershipResponse

type getMembershipBuilder struct {
	opts *getMembershipOpts
}

func newGetMembershipBuilder(pubnub *PubNub) *getMembershipBuilder {
	return newGetMembershipBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetMembershipOpts(pubnub *PubNub, ctx Context) *getMembershipOpts {
	return &getMembershipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetMembershipBuilderWithContext(pubnub *PubNub, context Context) *getMembershipBuilder {
	return &getMembershipBuilder{opts: newGetMembershipOpts(pubnub, context)}
}

// ID sets the required identifier of the membership to read.
func (b *getMembershipBuilder) ID(id string) *getMembershipBuilder {
	b.opts.ID = id
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMembershipBuilder) QueryParam(queryParam map[string]string) *getMembershipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getMembership request.
func (b *getMembershipBuilder) Transport(tr http.RoundTripper) *getMembershipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getMembershipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the getMembership request.
func (b *getMembershipBuilder) Execute() (*PNDataSyncMembershipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetDataSyncMembershipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetMembershipResponse, status, err
	}

	return newPNDataSyncMembershipResponse(rawJSON, status, PNGetDataSyncMembershipOperation, b.opts.pubnub)
}

type getMembershipOpts struct {
	endpointOpts

	ID         string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getMembershipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingMembershipID)
	}
	return nil
}

func (o *getMembershipOpts) buildPath() (string, error) {
	return fmt.Sprintf(membershipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getMembershipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *getMembershipOpts) httpMethod() string {
	return "GET"
}

func (o *getMembershipOpts) operationType() OperationType {
	return PNGetDataSyncMembershipOperation
}
