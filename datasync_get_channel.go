package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptyGetChannelResponse *PNDataSyncChannelResponse

type getChannelBuilder struct {
	opts *getChannelOpts
}

func newGetChannelBuilder(pubnub *PubNub) *getChannelBuilder {
	return newGetChannelBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetChannelOpts(pubnub *PubNub, ctx Context) *getChannelOpts {
	return &getChannelOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}

func newGetChannelBuilderWithContext(pubnub *PubNub, context Context) *getChannelBuilder {
	return &getChannelBuilder{opts: newGetChannelOpts(pubnub, context)}
}

// ID sets the required identifier of the channel to read.
func (b *getChannelBuilder) ID(id string) *getChannelBuilder {
	b.opts.ID = id
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getChannelBuilder) QueryParam(queryParam map[string]string) *getChannelBuilder {
	b.opts.QueryParam = queryParam
	return b
}

// Transport sets the Transport for the getChannel request.
func (b *getChannelBuilder) Transport(tr http.RoundTripper) *getChannelBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *getChannelOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"ID": o.ID,
	}
}

// Execute runs the getChannel request.
func (b *getChannelBuilder) Execute() (*PNDataSyncChannelResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNGetDataSyncChannelOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetChannelResponse, status, err
	}

	return newPNEntityResponse(rawJSON, status, PNGetDataSyncChannelOperation, b.opts.pubnub)
}

type getChannelOpts struct {
	endpointOpts

	ID         string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.ID == "" {
		return newValidationError(o, StrMissingChannelID)
	}
	return nil
}

func (o *getChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(channelsIDPath, o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *getChannelOpts) httpMethod() string {
	return "GET"
}

func (o *getChannelOpts) operationType() OperationType {
	return PNGetDataSyncChannelOperation
}
