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

const removeChannelFromChannelGroupPath = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyRemoveChannelFromChannelGroupResponse *RemoveChannelFromChannelGroupResponse

type removeChannelFromChannelGroupBuilder struct {
	opts *removeChannelOpts
}

func newRemoveChannelFromChannelGroupBuilder(
	pubnub *PubNub) *removeChannelFromChannelGroupBuilder {
	builder := removeChannelFromChannelGroupBuilder{
		opts: &removeChannelOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveChannelFromChannelGroupBuilderWithContext(
	pubnub *PubNub, context Context) *removeChannelFromChannelGroupBuilder {
	builder := removeChannelFromChannelGroupBuilder{
		opts: &removeChannelOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channels sets the chnanels to remove from the channel group
func (b *removeChannelFromChannelGroupBuilder) Channels(
	ch []string) *removeChannelFromChannelGroupBuilder {
	b.opts.Channels = ch
	return b
}

// ChannelGroup sets the ChannelGroup to remove the channels
func (b *removeChannelFromChannelGroupBuilder) ChannelGroup(
	cg string) *removeChannelFromChannelGroupBuilder {
	b.opts.ChannelGroup = cg
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelFromChannelGroupBuilder) QueryParam(queryParam map[string]string) *removeChannelFromChannelGroupBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs RemoveChannelFromChannelGroup request
func (b *removeChannelFromChannelGroupBuilder) Execute() (
	*RemoveChannelFromChannelGroupResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelFromChannelGroupResponse, status, err
	}

	return newRemoveChannelFromChannelGroupResponse(rawJSON, status)
}

type removeChannelOpts struct {
	pubnub *PubNub

	Channels     []string
	QueryParam   map[string]string
	ChannelGroup string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeChannelOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeChannelOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeChannelOpts) context() Context {
	return o.ctx
}

func (o *removeChannelOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.ChannelGroup == "" {
		return newValidationError(o, StrMissingChannelGroup)
	}

	return nil
}

func (o *removeChannelOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeChannelFromChannelGroupPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.ChannelGroup)), nil
}

func (o *removeChannelOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	var channels []string

	for _, ch := range o.Channels {
		channels = append(channels, utils.URLEncode(ch))
	}

	q.Set("remove", strings.Join(channels, ","))
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *removeChannelOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *removeChannelOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeChannelOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeChannelOpts) httpMethod() string {
	return "GET"
}

func (o *removeChannelOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeChannelOpts) operationType() OperationType {
	return PNRemoveChannelFromChannelGroupOperation
}

func (o *removeChannelOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// RemoveChannelFromChannelGroupResponse is the struct returned when the Execute function of RemoveChannelFromChannelGroup is called.
type RemoveChannelFromChannelGroupResponse struct {
}

func newRemoveChannelFromChannelGroupResponse(jsonBytes []byte,
	status StatusResponse) (*RemoveChannelFromChannelGroupResponse,
	StatusResponse, error) {
	return emptyRemoveChannelFromChannelGroupResponse, status, nil
}
