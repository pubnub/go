package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/v7/pnerr"
)

var emptySetChannelMembersResponse *PNSetChannelMembersResponse

const setChannelMembersPath = "/v2/objects/%s/channels/%s/uuids"

const setChannelMembersLimit = 100

type setChannelMembersBuilder struct {
	opts *setChannelMembersOpts
}

func newSetChannelMembersBuilder(pubnub *PubNub) *setChannelMembersBuilder {
	return newSetChannelMembersBuilderWithContext(pubnub, pubnub.ctx)
}

func newSetChannelMembersBuilderWithContext(pubnub *PubNub,
	context Context) *setChannelMembersBuilder {
	builder := setChannelMembersBuilder{
		opts: newSetChannelMembersOpts(
			pubnub,
			context,
		),
	}
	builder.opts.Limit = setChannelMembersLimit

	return &builder
}

func (b *setChannelMembersBuilder) Include(include []PNChannelMembersInclude) *setChannelMembersBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *setChannelMembersBuilder) Channel(channel string) *setChannelMembersBuilder {
	b.opts.Channel = channel

	return b
}

func (b *setChannelMembersBuilder) Limit(limit int) *setChannelMembersBuilder {
	b.opts.Limit = limit

	return b
}

func (b *setChannelMembersBuilder) Start(start string) *setChannelMembersBuilder {
	b.opts.Start = start

	return b
}

func (b *setChannelMembersBuilder) End(end string) *setChannelMembersBuilder {
	b.opts.End = end

	return b
}

func (b *setChannelMembersBuilder) Count(count bool) *setChannelMembersBuilder {
	b.opts.Count = count

	return b
}

func (b *setChannelMembersBuilder) Filter(filter string) *setChannelMembersBuilder {
	b.opts.Filter = filter

	return b
}

func (b *setChannelMembersBuilder) Sort(sort []string) *setChannelMembersBuilder {
	b.opts.Sort = sort

	return b
}

func (b *setChannelMembersBuilder) Set(channelMembersSet []PNChannelMembersSet) *setChannelMembersBuilder {
	b.opts.ChannelMembersSet = channelMembersSet

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setChannelMembersBuilder) QueryParam(queryParam map[string]string) *setChannelMembersBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the setChannelMembers request.
func (b *setChannelMembersBuilder) Transport(tr http.RoundTripper) *setChannelMembersBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the setChannelMembers request.
func (b *setChannelMembersBuilder) Execute() (*PNSetChannelMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetChannelMembersResponse, status, err
	}

	return newPNSetChannelMembersResponse(rawJSON, b.opts, status)
}

func newSetChannelMembersOpts(pubnub *PubNub, ctx Context) *setChannelMembersOpts {
	return &setChannelMembersOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
		Limit: setChannelMembersLimit}
}

type setChannelMembersOpts struct {
	endpointOpts
	Channel           string
	Limit             int
	Include           []string
	Start             string
	End               string
	Filter            string
	Sort              []string
	Count             bool
	QueryParam        map[string]string
	ChannelMembersSet []PNChannelMembersSet
	Transport         http.RoundTripper
}

func (o *setChannelMembersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *setChannelMembersOpts) buildPath() (string, error) {
	return fmt.Sprintf(setChannelMembersPath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *setChannelMembersOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}
	q.Set("limit", strconv.Itoa(o.Limit))

	if o.Start != "" {
		q.Set("start", o.Start)
	}

	if o.Count {
		q.Set("count", "1")
	} else {
		q.Set("count", "0")
	}

	if o.End != "" {
		q.Set("end", o.End)
	}
	if o.Filter != "" {
		q.Set("filter", o.Filter)
	}

	if o.Sort != nil {
		SetQueryParamAsCommaSepString(q, o.Sort, "sort")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

// PNChannelMembersSetChangeset is the Objects API input to add, remove or update membership
type PNChannelMembersSetChangeset struct {
	Set []PNChannelMembersSet `json:"set"`
}

func (o *setChannelMembersOpts) buildBody() ([]byte, error) {
	b := &PNChannelMembersSetChangeset{
		Set: o.ChannelMembersSet,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setChannelMembersOpts) httpMethod() string {
	return "PATCH"
}

func (o *setChannelMembersOpts) isAuthRequired() bool {
	return true
}

func (o *setChannelMembersOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *setChannelMembersOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *setChannelMembersOpts) operationType() OperationType {
	return PNSetChannelMembersOperation
}

// PNSetChannelMembersResponse is the Objects API Response for SetChannelMembers
type PNSetChannelMembersResponse struct {
	status     int                `json:"status"`
	Data       []PNChannelMembers `json:"data"`
	TotalCount int                `json:"totalCount"`
	Next       string             `json:"next"`
	Prev       string             `json:"prev"`
}

func newPNSetChannelMembersResponse(jsonBytes []byte, o *setChannelMembersOpts,
	status StatusResponse) (*PNSetChannelMembersResponse, StatusResponse, error) {

	resp := &PNSetChannelMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySetChannelMembersResponse, status, e
	}

	return resp, status, nil
}
