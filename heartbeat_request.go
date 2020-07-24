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
	"strings"

	"github.com/pubnub/go/utils"
)

const heartbeatPath = "/v2/presence/sub-key/%s/channel/%s/heartbeat"

type heartbeatBuilder struct {
	opts *heartbeatOpts
}

func newHeartbeatBuilder(pubnub *PubNub) *heartbeatBuilder {
	builder := heartbeatBuilder{
		opts: &heartbeatOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newHeartbeatBuilderWithContext(pubnub *PubNub,
	context Context) *heartbeatBuilder {
	builder := heartbeatBuilder{
		opts: &heartbeatOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *heartbeatBuilder) QueryParam(queryParam map[string]string) *heartbeatBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// State sets the state for the Heartbeat request.
func (b *heartbeatBuilder) State(state interface{}) *heartbeatBuilder {
	b.opts.State = state

	return b
}

// Channels sets the Channels for the Heartbeat request.
func (b *heartbeatBuilder) Channels(ch []string) *heartbeatBuilder {
	b.opts.Channels = ch

	return b
}

// ChannelGroups sets the ChannelGroups for the Heartbeat request.
func (b *heartbeatBuilder) ChannelGroups(cg []string) *heartbeatBuilder {
	b.opts.ChannelGroups = cg

	return b
}

// Execute runs the Heartbeat request
func (b *heartbeatBuilder) Execute() (interface{}, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return "", status, err
	}

	var value interface{}

	err = json.Unmarshal(rawJSON, &value)
	if err != nil {
		return nil, status, err
	}

	return value, status, nil
}

type heartbeatOpts struct {
	pubnub *PubNub

	State interface{}

	Channels      []string
	ChannelGroups []string
	QueryParam    map[string]string

	ctx Context
}

func (o *heartbeatOpts) config() Config {
	return *o.pubnub.Config
}

func (o *heartbeatOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *heartbeatOpts) context() Context {
	return o.ctx
}

func (o *heartbeatOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	return nil
}

func (o *heartbeatOpts) buildPath() (string, error) {
	channels := string(utils.JoinChannels(o.Channels))

	return fmt.Sprintf(heartbeatPath,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *heartbeatOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	q.Set("heartbeat", strconv.Itoa(o.pubnub.Config.PresenceTimeout))

	if len(o.ChannelGroups) > 0 {
		q.Set("channel-group", strings.Join(o.ChannelGroups, ","))
	}

	if o.State != nil {
		state, err := utils.ValueAsString(o.State)
		if err != nil {
			return &url.Values{}, err
		}

		if string(state) != "{}" {
			q.Set("state", string(state))
		}
	}
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *heartbeatOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *heartbeatOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *heartbeatOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *heartbeatOpts) httpMethod() string {
	return "GET"
}

func (o *heartbeatOpts) isAuthRequired() bool {
	return true
}

func (o *heartbeatOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *heartbeatOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *heartbeatOpts) operationType() OperationType {
	return PNHeartBeatOperation
}

func (o *heartbeatOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
