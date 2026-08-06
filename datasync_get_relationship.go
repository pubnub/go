package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptyGetRelationshipResponse *PNRelationshipResponse

type getRelationshipBuilder struct {
	opts *getRelationshipOpts
}

func newGetRelationshipBuilder(pubnub *PubNub) *getRelationshipBuilder {
	return newGetRelationshipBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetRelationshipOpts(pubnub *PubNub, ctx Context) *getRelationshipOpts {
	return &getRelationshipOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetRelationshipBuilderWithContext(pubnub *PubNub, context Context) *getRelationshipBuilder {
	return &getRelationshipBuilder{opts: newGetRelationshipOpts(pubnub, context)}
}

// ID sets the required identifier of the relationship to read.
func (b *getRelationshipBuilder) ID(id string) *getRelationshipBuilder {
	b.opts.ID = id
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getRelationshipBuilder) QueryParam(queryParam map[string]string) *getRelationshipBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getRelationship request.
func (b *getRelationshipBuilder) Transport(tr http.RoundTripper) *getRelationshipBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getRelationshipOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the getRelationship request.
func (b *getRelationshipBuilder) Execute() (*PNRelationshipResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetRelationshipOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetRelationshipResponse, status, err
	}

	return newPNRelationshipResponse(rawJSON, status, PNGetRelationshipOperation, b.opts.pubnub)
}

type getRelationshipOpts struct {
	endpointOpts

	ID         string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getRelationshipOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingRelationshipID)
	}
	return nil
}

func (o *getRelationshipOpts) buildPath() (string, error) {
	return fmt.Sprintf(relationshipsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getRelationshipOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *getRelationshipOpts) httpMethod() string {
	return "GET"
}

func (o *getRelationshipOpts) operationType() OperationType {
	return PNGetRelationshipOperation
}
