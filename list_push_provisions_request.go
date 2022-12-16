package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v7/pnerr"
	"github.com/pubnub/go/v7/utils"
)

const listChannelsOfPushPath = "/v1/push/sub-key/%s/devices/%s"
const listChannelsOfPushPathAPNS2 = "/v2/push/sub-key/%s/devices-apns2/%s"

var emptyListPushProvisionsRequestResponse *ListPushProvisionsRequestResponse

type listPushProvisionsRequestBuilder struct {
	opts *listPushProvisionsRequestOpts
}

func newListPushProvisionsRequestBuilder(pubnub *PubNub) *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilderWithContext(pubnub, pubnub.ctx)
}

func newListPushProvisionsRequestBuilderWithContext(
	pubnub *PubNub, context Context) *listPushProvisionsRequestBuilder {
	builder := listPushProvisionsRequestBuilder{
		opts: newListPushProvisionsRequestOpts(
			pubnub,
			context,
		),
	}

	return &builder
}

// PushType sets the PushType for the List Push Provisions request.
func (b *listPushProvisionsRequestBuilder) PushType(pushType PNPushType) *listPushProvisionsRequestBuilder {
	b.opts.PushType = pushType
	return b
}

// DeviceIDForPush sets the device id for List Push Provisions request.
func (b *listPushProvisionsRequestBuilder) DeviceIDForPush(deviceID string) *listPushProvisionsRequestBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

// Topic sets the topic of for APNS2 Push Notifcataions
func (b *listPushProvisionsRequestBuilder) Topic(topic string) *listPushProvisionsRequestBuilder {
	b.opts.Topic = topic
	return b
}

// Environment sets the environment of for APNS2 Push Notifcataions
func (b *listPushProvisionsRequestBuilder) Environment(env PNPushEnvironment) *listPushProvisionsRequestBuilder {
	b.opts.Environment = env
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *listPushProvisionsRequestBuilder) QueryParam(queryParam map[string]string) *listPushProvisionsRequestBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the List Push Provisions request.
func (b *listPushProvisionsRequestBuilder) Execute() (
	*ListPushProvisionsRequestResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyListPushProvisionsRequestResponse, status, err
	}

	return newListPushProvisionsRequestResponse(rawJSON, status)
}

func newListPushProvisionsRequestResponse(jsonBytes []byte, status StatusResponse) (
	*ListPushProvisionsRequestResponse, StatusResponse, error) {
	resp := &ListPushProvisionsRequestResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyListPushProvisionsRequestResponse, status, e
	}

	if parsedValue, ok := value.([]interface{}); ok {
		a := make([]string, len(parsedValue))
		for i, v := range parsedValue {
			a[i] = v.(string)
		}
		resp.Channels = a
	}

	return resp, status, nil
}

func newListPushProvisionsRequestOpts(pubnub *PubNub, ctx Context) *listPushProvisionsRequestOpts {
	return &listPushProvisionsRequestOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
	}
}

type listPushProvisionsRequestOpts struct {
	endpointOpts

	PushType PNPushType

	DeviceIDForPush string
	QueryParam      map[string]string
	Transport       http.RoundTripper
	Topic           string
	Environment     PNPushEnvironment
}

func (o *listPushProvisionsRequestOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
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

// ListPushProvisionsRequestResponse is the struct returned when the Execute function of ListPushProvisions is called.
type ListPushProvisionsRequestResponse struct {
	Channels []string
}

func (o *listPushProvisionsRequestOpts) buildPath() (string, error) {
	if o.PushType == PNPushTypeAPNS2 {
		return fmt.Sprintf(listChannelsOfPushPathAPNS2,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.DeviceIDForPush)), nil
	}

	return fmt.Sprintf(listChannelsOfPushPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.DeviceIDForPush)), nil
}

func (o *listPushProvisionsRequestOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)
	q.Set("type", o.PushType.String())
	SetPushEnvironment(q, o.Environment)
	SetPushTopic(q, o.Topic)

	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *listPushProvisionsRequestOpts) isAuthRequired() bool {
	return true
}

func (o *listPushProvisionsRequestOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *listPushProvisionsRequestOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *listPushProvisionsRequestOpts) operationType() OperationType {
	return PNRemoveGroupOperation
}
