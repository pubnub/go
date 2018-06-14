package pubnub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pubnub/go/utils"
)

const removeAllPushChannelsForDevicePath = "/v1/push/sub-key/%s/devices/%s/remove"

var emptyRemoveAllPushChannelsForDeviceResponse *RemoveAllPushChannelsForDeviceResponse

type RemoveAllPushChannelsForDeviceBuilder struct {
	opts *removeAllPushChannelsForDeviceOpts
}

func newRemoveAllPushChannelsForDeviceBuilder(pubnub *PubNub) *RemoveAllPushChannelsForDeviceBuilder {
	builder := RemoveAllPushChannelsForDeviceBuilder{
		opts: &removeAllPushChannelsForDeviceOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveAllPushChannelsForDeviceBuilderWithContext(
	pubnub *PubNub, context Context) *RemoveAllPushChannelsForDeviceBuilder {
	builder := RemoveAllPushChannelsForDeviceBuilder{
		opts: &removeAllPushChannelsForDeviceOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

//
func (b *RemoveAllPushChannelsForDeviceBuilder) PushType(
	pushType PNPushType) *RemoveAllPushChannelsForDeviceBuilder {
	b.opts.PushType = pushType
	return b
}

func (b *RemoveAllPushChannelsForDeviceBuilder) DeviceIDForPush(
	deviceID string) *RemoveAllPushChannelsForDeviceBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

func (b *RemoveAllPushChannelsForDeviceBuilder) Execute() (
	*RemoveAllPushChannelsForDeviceResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveAllPushChannelsForDeviceResponse, status, err
	}

	return emptyRemoveAllPushChannelsForDeviceResponse, status, err
}

type removeAllPushChannelsForDeviceOpts struct {
	pubnub *PubNub

	PushType PNPushType

	DeviceIDForPush string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeAllPushChannelsForDeviceOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeAllPushChannelsForDeviceOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeAllPushChannelsForDeviceOpts) context() Context {
	return o.ctx
}

func (o *removeAllPushChannelsForDeviceOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.DeviceIDForPush == "" {
		return newValidationError(o, StrMissingDeviceID)
	}

	if o.PushType == PNPushTypeNone {
		return newValidationError(o, StrMissingPushType)
	}

	return nil
}

type RemoveAllPushChannelsForDeviceResponse struct{}

func (o *removeAllPushChannelsForDeviceOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeAllPushChannelsForDevicePath,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.DeviceIDForPush)), nil
}

func (o *removeAllPushChannelsForDeviceOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	q.Set("type", o.PushType.String())

	return q, nil
}

func (o *removeAllPushChannelsForDeviceOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeAllPushChannelsForDeviceOpts) httpMethod() string {
	return "GET"
}

func (o *removeAllPushChannelsForDeviceOpts) isAuthRequired() bool {
	return true
}

func (o *removeAllPushChannelsForDeviceOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeAllPushChannelsForDeviceOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeAllPushChannelsForDeviceOpts) operationType() OperationType {
	return PNRemoveGroupOperation
}

func (o *removeAllPushChannelsForDeviceOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
