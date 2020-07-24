package pubnub

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/pubnub/go/utils"
)

const addChannelsToPushPath = "/v1/push/sub-key/%s/devices/%s"
const addChannelsToPushPathAPNS2 = "/v2/push/sub-key/%s/devices-apns2/%s"

var emptyAddPushNotificationsOnChannelsResponse *AddPushNotificationsOnChannelsResponse

// addPushNotificationsOnChannelsBuilder provides a builder to add Push Notifications on channels
type addPushNotificationsOnChannelsBuilder struct {
	opts *addChannelsToPushOpts
}

func newAddPushNotificationsOnChannelsBuilder(pubnub *PubNub) *addPushNotificationsOnChannelsBuilder {
	builder := addPushNotificationsOnChannelsBuilder{
		opts: &addChannelsToPushOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newAddPushNotificationsOnChannelsBuilderWithContext(
	pubnub *PubNub, context Context) *addPushNotificationsOnChannelsBuilder {
	builder := addPushNotificationsOnChannelsBuilder{
		opts: &addChannelsToPushOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channels sets the channels to enable Push Notifications
func (b *addPushNotificationsOnChannelsBuilder) Channels(channels []string) *addPushNotificationsOnChannelsBuilder {
	b.opts.Channels = channels

	return b
}

// PushType set the type of Push: GCM, APNS, MPNS
func (b *addPushNotificationsOnChannelsBuilder) PushType(pushType PNPushType) *addPushNotificationsOnChannelsBuilder {
	b.opts.PushType = pushType
	return b
}

// DeviceIDForPush sets the device of for Push Notifcataions
func (b *addPushNotificationsOnChannelsBuilder) DeviceIDForPush(deviceID string) *addPushNotificationsOnChannelsBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

// Topic sets the topic of for APNS2 Push Notifcataions
func (b *addPushNotificationsOnChannelsBuilder) Topic(topic string) *addPushNotificationsOnChannelsBuilder {
	b.opts.Topic = topic
	return b
}

// Environment sets the environment of for APNS2 Push Notifcataions
func (b *addPushNotificationsOnChannelsBuilder) Environment(env PNPushEnvironment) *addPushNotificationsOnChannelsBuilder {
	b.opts.Environment = env
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *addPushNotificationsOnChannelsBuilder) QueryParam(queryParam map[string]string) *addPushNotificationsOnChannelsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs add Push Notifications on channels request
func (b *addPushNotificationsOnChannelsBuilder) Execute() (*AddPushNotificationsOnChannelsResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAddPushNotificationsOnChannelsResponse, status, err
	}

	return emptyAddPushNotificationsOnChannelsResponse, status, nil
}

type addChannelsToPushOpts struct {
	pubnub          *PubNub
	Channels        []string
	PushType        PNPushType
	DeviceIDForPush string
	QueryParam      map[string]string
	Transport       http.RoundTripper
	ctx             Context
	Topic           string
	Environment     PNPushEnvironment
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

	if o.PushType == PNPushTypeAPNS2 && (o.Topic == "") {
		return newValidationError(o, StrMissingPushTopic)
	}

	return nil
}

// AddPushNotificationsOnChannelsResponse is response structure for AddPushNotificationsOnChannelsBuilder
type AddPushNotificationsOnChannelsResponse struct{}

func (o *addChannelsToPushOpts) buildPath() (string, error) {
	if o.PushType == PNPushTypeAPNS2 {
		return fmt.Sprintf(addChannelsToPushPathAPNS2,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.DeviceIDForPush)), nil
	}

	return fmt.Sprintf(addChannelsToPushPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.DeviceIDForPush)), nil
}

func (o *addChannelsToPushOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	var channels []string

	for _, v := range o.Channels {
		channels = append(channels, utils.URLEncode(v))
	}

	q.Set("add", strings.Join(channels, ","))
	q.Set("type", o.PushType.String())
	SetPushEnvironment(q, o.Environment)
	SetPushTopic(q, o.Topic)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *addChannelsToPushOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *addChannelsToPushOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *addChannelsToPushOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
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
