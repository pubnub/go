package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptyGetEntityResponse *PNEntityResponse

type getEntityBuilder struct {
	opts *getEntityOpts
}

func newGetEntityBuilder(pubnub *PubNub) *getEntityBuilder {
	return newGetEntityBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetEntityOpts(pubnub *PubNub, ctx Context) *getEntityOpts {
	return &getEntityOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetEntityBuilderWithContext(pubnub *PubNub, context Context) *getEntityBuilder {
	return &getEntityBuilder{opts: newGetEntityOpts(pubnub, context)}
}

// ID sets the required identifier of the entity to read.
func (b *getEntityBuilder) ID(id string) *getEntityBuilder {
	b.opts.ID = id
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getEntityBuilder) QueryParam(queryParam map[string]string) *getEntityBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getEntity request.
func (b *getEntityBuilder) Transport(tr http.RoundTripper) *getEntityBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getEntityOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the getEntity request.
func (b *getEntityBuilder) Execute() (*PNEntityResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetEntityOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetEntityResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNGetEntityOperation, b.opts.pubnub)
}

type getEntityOpts struct {
	endpointOpts

	ID         string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getEntityOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingEntityID)
	}
	return nil
}

func (o *getEntityOpts) buildPath() (string, error) {
	return fmt.Sprintf(entitiesIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getEntityOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *getEntityOpts) httpMethod() string {
	return "GET"
}

func (o *getEntityOpts) operationType() OperationType {
	return PNGetEntityOperation
}
