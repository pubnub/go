package pubnub

import (
	"fmt"
	"github.com/pubnub/go/utils"
	"net/http"
	"net/url"
	"strings"
)

const ADD_CHANNELS_TO_PUSH = "/v1/push/sub-key/%s/devices/%s"

var emptyAddChannelsToPushResponse *AddChannelsToPushResponse

type addChannelsToPushBuilder struct {
	opts *addChannelsToPushOpts
}

func newAddChannelsToPushBuilder(pubnub *PubNub) *addChannelsToPushBuilder {
	builder := addChannelsToPushBuilder{
		opts: &addChannelsToPushOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newAddChannelsToPushBuilderWithContext(
	pubnub *PubNub, context Context) *addChannelsToPushBuilder {
	builder := addChannelsToPushBuilder{
		opts: &addChannelsToPushOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *addChannelsToPushBuilder) Channels(
	channels []string) *addChannelsToPushBuilder {
	b.opts.Channels = channels
	return b
}

func (b *addChannelsToPushBuilder) PushType(
	pushType PNPushType) *addChannelsToPushBuilder {
	b.opts.PushType = pushType
	return b
}

func (b *addChannelsToPushBuilder) DeviceIDForPush(
	deviceID string) *addChannelsToPushBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

func (b *addChannelsToPushBuilder) Execute() (
	*AddChannelsToPushResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAddChannelsToPushResponse, status, err
	}

	return emptyAddChannelsToPushResponse, status, nil
}

type addChannelsToPushOpts struct {
	pubnub *PubNub

	Channels []string

	PushType PNPushType

	DeviceIDForPush string

	Transport http.RoundTripper

	ctx Context
}

func (o *addChannelsToPushOpts) config() Config {
	return *o.pubnub.Config
}

func (o *addChannelsToPushOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *addChannelsToPushOpts) context() Context {
	return o.ctx
}

func (o *addChannelsToPushOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.DeviceIDForPush == "" {
		return newValidationError(o, StrMissingDeviceID)
	}

	if len(o.Channels) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.PushType == PNPushTypeNone {
		return newValidationError(o, StrMissingPushType)
	}

	return nil
}

type AddChannelsToPushResponse struct{}

func (o *addChannelsToPushOpts) buildPath() (string, error) {
	return fmt.Sprintf(ADD_CHANNELS_TO_PUSH,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.DeviceIDForPush)), nil
}

func (o *addChannelsToPushOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	var channels []string

	for _, v := range o.Channels {
		channels = append(channels, utils.UrlEncode(v))
	}

	q.Set("add", strings.Join(channels, ","))
	q.Set("type", o.PushType.String())

	return q, nil
}

func (o *addChannelsToPushOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *addChannelsToPushOpts) httpMethod() string {
	return "GET"
}

func (o *addChannelsToPushOpts) isAuthRequired() bool {
	return true
}

func (o *addChannelsToPushOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *addChannelsToPushOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *addChannelsToPushOpts) operationType() OperationType {
	return PNRemoveGroupOperation
}

func (o *addChannelsToPushOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
