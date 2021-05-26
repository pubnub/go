package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/v5/pnerr"
)

var emptyGetChannelMembersResponse *PNGetChannelMembersResponse

const getChannelMembersPathV2 = "/v2/objects/%s/channels/%s/uuids"

const membersLimitV2 = 100

type getChannelMembersBuilderV2 struct {
	opts *getChannelMembersOptsV2
}

func newGetChannelMembersBuilderV2(pubnub *PubNub) *getChannelMembersBuilderV2 {
	builder := getChannelMembersBuilderV2{
		opts: &getChannelMembersOptsV2{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = membersLimitV2

	return &builder
}

func newGetChannelMembersBuilderV2WithContext(pubnub *PubNub,
	context Context) *getChannelMembersBuilderV2 {
	builder := getChannelMembersBuilderV2{
		opts: &getChannelMembersOptsV2{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getChannelMembersBuilderV2) Channel(channel string) *getChannelMembersBuilderV2 {
	b.opts.Channel = channel

	return b
}

func (b *getChannelMembersBuilderV2) Include(include []PNChannelMembersInclude) *getChannelMembersBuilderV2 {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getChannelMembersBuilderV2) Limit(limit int) *getChannelMembersBuilderV2 {
	b.opts.Limit = limit

	return b
}

func (b *getChannelMembersBuilderV2) Start(start string) *getChannelMembersBuilderV2 {
	b.opts.Start = start

	return b
}

func (b *getChannelMembersBuilderV2) End(end string) *getChannelMembersBuilderV2 {
	b.opts.End = end

	return b
}

func (b *getChannelMembersBuilderV2) Filter(filter string) *getChannelMembersBuilderV2 {
	b.opts.Filter = filter

	return b
}

func (b *getChannelMembersBuilderV2) Sort(sort []string) *getChannelMembersBuilderV2 {
	b.opts.Sort = sort

	return b
}

func (b *getChannelMembersBuilderV2) Count(count bool) *getChannelMembersBuilderV2 {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getChannelMembersBuilderV2) QueryParam(queryParam map[string]string) *getChannelMembersBuilderV2 {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getChannelMembers request.
func (b *getChannelMembersBuilderV2) Transport(tr http.RoundTripper) *getChannelMembersBuilderV2 {
	b.opts.Transport = tr
	return b
}

// Execute runs the getChannelMembers request.
func (b *getChannelMembersBuilderV2) Execute() (*PNGetChannelMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetChannelMembersResponse, status, err
	}

	return newPNGetChannelMembersResponse(rawJSON, b.opts, status)
}

type getChannelMembersOptsV2 struct {
	pubnub     *PubNub
	Channel    string
	Limit      int
	Include    []string
	Start      string
	End        string
	Filter     string
	Sort       []string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getChannelMembersOptsV2) config() Config {
	return *o.pubnub.Config
}

func (o *getChannelMembersOptsV2) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getChannelMembersOptsV2) context() Context {
	return o.ctx
}

func (o *getChannelMembersOptsV2) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *getChannelMembersOptsV2) buildPath() (string, error) {
	return fmt.Sprintf(getChannelMembersPathV2,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *getChannelMembersOptsV2) buildQuery() (*url.Values, error) {

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

func (o *getChannelMembersOptsV2) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getChannelMembersOptsV2) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getChannelMembersOptsV2) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getChannelMembersOptsV2) httpMethod() string {
	return "GET"
}

func (o *getChannelMembersOptsV2) isAuthRequired() bool {
	return true
}

func (o *getChannelMembersOptsV2) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getChannelMembersOptsV2) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getChannelMembersOptsV2) operationType() OperationType {
	return PNGetChannelMembersOperation
}

func (o *getChannelMembersOptsV2) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetChannelMembersResponse is the Objects API Response for Get Members
type PNGetChannelMembersResponse struct {
	status     int                `json:"status"`
	Data       []PNChannelMembers `json:"data"`
	TotalCount int                `json:"totalCount"`
	Next       string             `json:"next"`
	Prev       string             `json:"prev"`
}

func newPNGetChannelMembersResponse(jsonBytes []byte, o *getChannelMembersOptsV2,
	status StatusResponse) (*PNGetChannelMembersResponse, StatusResponse, error) {

	resp := &PNGetChannelMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetChannelMembersResponse, status, e
	}

	return resp, status, nil
}
