package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v7/pnerr"
)

var emptyPNRemoveMessageActionsResponse *PNRemoveMessageActionsResponse

const removeMessageActionsPath = "/v1/message-actions/%s/channel/%s/message/%s/action/%s"

type removeMessageActionsBuilder struct {
	opts *removeMessageActionsOpts
}

func newRemoveMessageActionsBuilder(pubnub *PubNub) *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilderWithContext(pubnub, pubnub.ctx)
}

func newRemoveMessageActionsOpts(pubnub *PubNub, ctx Context) *removeMessageActionsOpts {
	return &removeMessageActionsOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newRemoveMessageActionsBuilderWithContext(pubnub *PubNub,
	context Context) *removeMessageActionsBuilder {
	builder := removeMessageActionsBuilder{
		opts: newRemoveMessageActionsOpts(pubnub, context)}
	return &builder
}

func (b *removeMessageActionsBuilder) Channel(channel string) *removeMessageActionsBuilder {
	b.opts.Channel = channel

	return b
}

func (b *removeMessageActionsBuilder) MessageTimetoken(timetoken string) *removeMessageActionsBuilder {
	b.opts.MessageTimetoken = timetoken

	return b
}

func (b *removeMessageActionsBuilder) ActionTimetoken(timetoken string) *removeMessageActionsBuilder {
	b.opts.ActionTimetoken = timetoken

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeMessageActionsBuilder) QueryParam(queryParam map[string]string) *removeMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeMessageActions request.
func (b *removeMessageActionsBuilder) Transport(tr http.RoundTripper) *removeMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeMessageActions request.
func (b *removeMessageActionsBuilder) Execute() (*PNRemoveMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNRemoveMessageActionsResponse, status, err
	}

	return newPNRemoveMessageActionsResponse(rawJSON, b.opts, status)
}

type removeMessageActionsOpts struct {
	endpointOpts

	Channel          string
	MessageTimetoken string
	ActionTimetoken  string
	Custom           map[string]interface{}
	QueryParam       map[string]string

	Transport http.RoundTripper
}

func (o *removeMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeMessageActionsPath,
		o.pubnub.Config.SubscribeKey, o.Channel, o.MessageTimetoken, o.ActionTimetoken), nil
}

func (o *removeMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeMessageActionsOpts) httpMethod() string {
	return "DELETE"
}

func (o *removeMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *removeMessageActionsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeMessageActionsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeMessageActionsOpts) operationType() OperationType {
	return PNRemoveMessageActionsOperation
}

// PNRemoveMessageActionsResponse is the Objects API Response for create space
type PNRemoveMessageActionsResponse struct {
	status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newPNRemoveMessageActionsResponse(jsonBytes []byte, o *removeMessageActionsOpts,
	status StatusResponse) (*PNRemoveMessageActionsResponse, StatusResponse, error) {

	resp := &PNRemoveMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNRemoveMessageActionsResponse, status, e
	}

	return resp, status, nil
}
