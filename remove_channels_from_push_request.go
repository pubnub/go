package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sprucehealth/pubnub-go/utils"
)

const removeChannelsFromPushPath = "/v1/push/sub-key/%s/devices/%s"

var emptyRemoveChannelsFromPushResponse *RemoveChannelsFromPushResponse

type removeChannelsFromPushBuilder struct {
	opts *removeChannelsFromPushOpts
}

func newRemoveChannelsFromPushBuilder(pubnub *PubNub) *removeChannelsFromPushBuilder {
	builder := removeChannelsFromPushBuilder{
		opts: &removeChannelsFromPushOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveChannelsFromPushBuilderWithContext(
	pubnub *PubNub, context Context) *removeChannelsFromPushBuilder {
	builder := removeChannelsFromPushBuilder{
		opts: &removeChannelsFromPushOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channels sets the channels to remove from Push Notifications
func (b *removeChannelsFromPushBuilder) Channels(
	channels []string) *removeChannelsFromPushBuilder {
	b.opts.Channels = channels
	return b
}

// PushType sets the PushType for the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) PushType(
	pushType PNPushType) *removeChannelsFromPushBuilder {
	b.opts.PushType = pushType
	return b
}

// DeviceIDForPush sets the DeviceIDForPush for the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) DeviceIDForPush(
	deviceID string) *removeChannelsFromPushBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelsFromPushBuilder) QueryParam(queryParam map[string]string) *removeChannelsFromPushBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) Execute() (
	*RemoveChannelsFromPushResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelsFromPushResponse, status, err
	}

	return emptyRemoveChannelsFromPushResponse, status, err
}

type removeChannelsFromPushOpts struct {
	pubnub *PubNub

	Channels        []string
	QueryParam      map[string]string
	PushType        PNPushType
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

// RemoveChannelsFromPushResponse is the struct returned when the Execute function of RemovePushNotificationsFromChannels is called.
type RemoveChannelsFromPushResponse struct{}

func (o *removeChannelsFromPushOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeChannelsFromPushPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.DeviceIDForPush)), nil
}

func (o *removeChannelsFromPushOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	q.Set("type", o.PushType.String())
	var channels []string

	for _, v := range o.Channels {
		channels = append(channels, utils.URLEncode(v))
	}

	q.Set("remove", strings.Join(channels, ","))
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *removeChannelsFromPushOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
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
