package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/v8/utils"
)

const removeChannelsFromPushPath = "/v1/push/sub-key/%s/devices/%s"
const removeChannelsFromPushPathAPNS2 = "/v2/push/sub-key/%s/devices-apns2/%s"

var emptyRemoveChannelsFromPushResponse *RemoveChannelsFromPushResponse

type removeChannelsFromPushBuilder struct {
	opts *removeChannelsFromPushOpts
}

func newRemoveChannelsFromPushBuilder(pubnub *PubNub) *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilderWithContext(pubnub, pubnub.ctx)
}

func newRemoveChannelsFromPushBuilderWithContext(pubnub *PubNub, context Context) *removeChannelsFromPushBuilder {
	builder := removeChannelsFromPushBuilder{
		opts: newRemoveChannelsFromPushOpts(pubnub, context)}
	return &builder
}

// Channels sets the channels to remove from Push Notifications
func (b *removeChannelsFromPushBuilder) Channels(channels []string) *removeChannelsFromPushBuilder {
	b.opts.Channels = channels
	return b
}

// PushType sets the PushType for the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) PushType(pushType PNPushType) *removeChannelsFromPushBuilder {
	b.opts.PushType = pushType
	return b
}

// DeviceIDForPush sets the DeviceIDForPush for the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) DeviceIDForPush(deviceID string) *removeChannelsFromPushBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

// Topic sets the topic of for APNS2 Push Notifcataions
func (b *removeChannelsFromPushBuilder) Topic(topic string) *removeChannelsFromPushBuilder {
	b.opts.Topic = topic
	return b
}

// Environment sets the environment of for APNS2 Push Notifcataions
func (b *removeChannelsFromPushBuilder) Environment(env PNPushEnvironment) *removeChannelsFromPushBuilder {
	b.opts.Environment = env
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelsFromPushBuilder) QueryParam(queryParam map[string]string) *removeChannelsFromPushBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the RemoveChannelsFromPush request.
func (b *removeChannelsFromPushBuilder) Transport(tr http.RoundTripper) *removeChannelsFromPushBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) Execute() (*RemoveChannelsFromPushResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelsFromPushResponse, status, err
	}

	return emptyRemoveChannelsFromPushResponse, status, err
}

func newRemoveChannelsFromPushOpts(pubnub *PubNub, ctx Context) *removeChannelsFromPushOpts {
	return &removeChannelsFromPushOpts{endpointOpts: endpointOpts{
		pubnub: pubnub,
		ctx:    ctx,
	}}
}

type removeChannelsFromPushOpts struct {
	endpointOpts

	Channels        []string
	QueryParam      map[string]string
	PushType        PNPushType
	DeviceIDForPush string
	Topic           string
	Environment     PNPushEnvironment

	Transport http.RoundTripper
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

	if o.PushType == PNPushTypeAPNS2 && (o.Topic == "") {
		return newValidationError(o, StrMissingPushTopic)
	}

	return nil
}

// RemoveChannelsFromPushResponse is the struct returned when the Execute function of RemovePushNotificationsFromChannels is called.
type RemoveChannelsFromPushResponse struct{}

func (o *removeChannelsFromPushOpts) buildPath() (string, error) {
	if o.PushType == PNPushTypeAPNS2 {
		return fmt.Sprintf(removeChannelsFromPushPathAPNS2,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.DeviceIDForPush)), nil

	}
	return fmt.Sprintf(removeChannelsFromPushPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.DeviceIDForPush)), nil
}

func (o *removeChannelsFromPushOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	q.Set("type", o.PushType.String())

	q.Set("remove", strings.Join(o.Channels, ","))
	SetPushEnvironment(q, o.Environment)
	SetPushTopic(q, o.Topic)
	SetQueryParam(q, o.QueryParam)
	return q, nil
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

func (o *removeChannelsFromPushOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeChannelsFromPushOpts) httpMethod() string {
	return "GET"
}

func (o *removeChannelsFromPushOpts) operationType() OperationType {
	return PNRemovePushNotificationsFromChannelsOperation
}
