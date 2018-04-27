package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/utils"
)

const removeChannelsFromPushPath = "/v1/push/sub-key/%s/devices/%s"

var emptyRemoveChannelsFromPushResponse *RemoveChannelsFromPushResponse

type RemoveChannelsFromPushBuilder struct {
	opts *removeChannelsFromPushOpts
}

func newRemoveChannelsFromPushBuilder(pubnub *PubNub) *RemoveChannelsFromPushBuilder {
	builder := RemoveChannelsFromPushBuilder{
		opts: &removeChannelsFromPushOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveChannelsFromPushBuilderWithContext(
	pubnub *PubNub, context Context) *RemoveChannelsFromPushBuilder {
	builder := RemoveChannelsFromPushBuilder{
		opts: &removeChannelsFromPushOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *RemoveChannelsFromPushBuilder) Channels(
	channels []string) *RemoveChannelsFromPushBuilder {
	b.opts.Channels = channels
	return b
}

func (b *RemoveChannelsFromPushBuilder) PushType(
	pushType PNPushType) *RemoveChannelsFromPushBuilder {
	b.opts.PushType = pushType
	return b
}

func (b *RemoveChannelsFromPushBuilder) DeviceIDForPush(
	deviceID string) *RemoveChannelsFromPushBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

func (b *RemoveChannelsFromPushBuilder) Execute() (
	*RemoveChannelsFromPushResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelsFromPushResponse, status, err
	}

	return emptyRemoveChannelsFromPushResponse, status, err
}

type removeChannelsFromPushOpts struct {
	pubnub *PubNub

	Channels []string

	PushType PNPushType

	DeviceIDForPush string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeChannelsFromPushOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeChannelsFromPushOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeChannelsFromPushOpts) context() Context {
	return o.ctx
}

func (o *removeChannelsFromPushOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.DeviceIDForPush == "" {
		return newValidationError(o, StrMissingDeviceID)
	}

	if o.PushType == PNPushTypeNone {
		return newValidationError(o, StrMissingPushType)
	}

	return nil
}

type RemoveChannelsFromPushResponse struct{}

func (o *removeChannelsFromPushOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeChannelsFromPushPath,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.DeviceIDForPush)), nil
}

func (o *removeChannelsFromPushOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)
	q.Set("type", o.PushType.String())
	var channels []string

	for _, v := range o.Channels {
		channels = append(channels, utils.UrlEncode(v))
	}

	q.Set("remove", strings.Join(channels, ","))

	return q, nil
}

func (o *removeChannelsFromPushOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeChannelsFromPushOpts) httpMethod() string {
	return "GET"
}

func (o *removeChannelsFromPushOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelsFromPushOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelsFromPushOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeChannelsFromPushOpts) operationType() OperationType {
	return PNRemoveGroupOperation
}

func (o *removeChannelsFromPushOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
