package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/utils"
)

const subscribePath = "/v2/subscribe/%s/%s/0"

type subscribeOpts struct {
	pubnub *PubNub

	Channels         []string
	ChannelGroups    []string
	QueryParam       map[string]string
	Heartbeat        int
	Region           string
	Timetoken        int64
	FilterExpression string
	WithPresence     bool
	State            map[string]interface{}
	stringState      string

	ctx Context
}

type subscribeBuilder struct {
	opts      *subscribeOpts
	operation *SubscribeOperation
}

func newSubscribeBuilder(pubnub *PubNub) *subscribeBuilder {
	builder := subscribeBuilder{
		opts: &subscribeOpts{
			pubnub: pubnub,
		},
		operation: &SubscribeOperation{},
	}

	return &builder
}

// Channels sets the channels to subscribe.
func (b *subscribeBuilder) Channels(channels []string) *subscribeBuilder {
	b.operation.Channels = channels

	return b
}

// ChannelGroups sets the channel groups to subscribe.
func (b *subscribeBuilder) ChannelGroups(groups []string) *subscribeBuilder {
	b.operation.ChannelGroups = groups

	return b
}

// Timetoken sets the timetoken to subscribe. Subscribe will start to fetch the messages from this timetoken onwards.
func (b *subscribeBuilder) Timetoken(tt int64) *subscribeBuilder {
	b.operation.Timetoken = tt

	return b
}

// FilterExpression sets the custom filter expression.
func (b *subscribeBuilder) FilterExpression(expr string) *subscribeBuilder {
	b.operation.FilterExpression = expr

	return b
}

// WithPresence as true subscribes to the presence channels as well.
func (b *subscribeBuilder) WithPresence(pres bool) *subscribeBuilder {
	b.operation.PresenceEnabled = pres

	return b
}

// State sets the state of the channels while subscribing.
func (b *subscribeBuilder) State(state map[string]interface{}) *subscribeBuilder {
	b.operation.State = state
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *subscribeBuilder) QueryParam(queryParam map[string]string) *subscribeBuilder {
	b.operation.QueryParam = queryParam

	return b
}

// Execute runs the Subscribe operation.
func (b *subscribeBuilder) Execute() {
	b.opts.pubnub.subscriptionManager.adaptSubscribe(b.operation)
}

func (o *subscribeOpts) config() Config {
	return *o.pubnub.Config
}

func (o *subscribeOpts) client() *http.Client {
	return o.pubnub.GetSubscribeClient()
}

func (o *subscribeOpts) context() Context {
	return o.ctx
}

func (o *subscribeOpts) validate() error {

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.State != nil {
		state, err := json.Marshal(o.State)
		if err != nil {
			return newValidationError(o, err.Error())
		}

		o.stringState = string(state)
	}
	return nil
}

func (o *subscribeOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	return fmt.Sprintf(subscribePath,
		o.pubnub.Config.SubscribeKey,
		channels,
	), nil
}

func (o *subscribeOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if len(o.ChannelGroups) > 0 {
		channelGroup := utils.JoinChannels(o.ChannelGroups)
		q.Set("channel-group", string(channelGroup))
	}

	if o.Timetoken != 0 {
		q.Set("tt", strconv.FormatInt(o.Timetoken, 10))
	}

	if o.Region != "" {
		q.Set("tr", o.Region)
	}

	if o.FilterExpression != "" {
		q.Set("filter-expr", o.FilterExpression)
	}

	// hb timeout should be at least 4 seconds
	if o.Heartbeat >= 4 {
		q.Set("heartbeat", fmt.Sprintf("%d", o.Heartbeat))
	}

	if o.stringState != "" {
		q.Set("state", o.stringState)
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *subscribeOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *subscribeOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *subscribeOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *subscribeOpts) httpMethod() string {
	return "GET"
}

func (o *subscribeOpts) isAuthRequired() bool {
	return true
}

func (o *subscribeOpts) requestTimeout() int {
	return o.pubnub.Config.SubscribeRequestTimeout
}

func (o *subscribeOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *subscribeOpts) operationType() OperationType {
	return PNSubscribeOperation
}

func (o *subscribeOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
