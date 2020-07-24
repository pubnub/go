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

	"github.com/pubnub/go/pnerr"
)

var emptyManageMembersResponse *PNManageMembersResponse

const manageMembersPathV2 = "/v2/objects/%s/channels/%s/uuids"

const manageMembersLimitV2 = 100

type manageChannelMembersBuilderV2 struct {
	opts *manageMembersOptsV2
}

func newManageChannelMembersBuilderV2(pubnub *PubNub) *manageChannelMembersBuilderV2 {
	builder := manageChannelMembersBuilderV2{
		opts: &manageMembersOptsV2{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = manageMembersLimitV2

	return &builder
}

func newManageChannelMembersBuilderV2WithContext(pubnub *PubNub,
	context Context) *manageChannelMembersBuilderV2 {
	builder := manageChannelMembersBuilderV2{
		opts: &manageMembersOptsV2{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *manageChannelMembersBuilderV2) Include(include []PNChannelMembersInclude) *manageChannelMembersBuilderV2 {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *manageChannelMembersBuilderV2) Channel(channel string) *manageChannelMembersBuilderV2 {
	b.opts.Channel = channel

	return b
}

func (b *manageChannelMembersBuilderV2) Limit(limit int) *manageChannelMembersBuilderV2 {
	b.opts.Limit = limit

	return b
}

func (b *manageChannelMembersBuilderV2) Start(start string) *manageChannelMembersBuilderV2 {
	b.opts.Start = start

	return b
}

func (b *manageChannelMembersBuilderV2) End(end string) *manageChannelMembersBuilderV2 {
	b.opts.End = end

	return b
}

func (b *manageChannelMembersBuilderV2) Count(count bool) *manageChannelMembersBuilderV2 {
	b.opts.Count = count

	return b
}

func (b *manageChannelMembersBuilderV2) Filter(filter string) *manageChannelMembersBuilderV2 {
	b.opts.Filter = filter

	return b
}

func (b *manageChannelMembersBuilderV2) Sort(sort []string) *manageChannelMembersBuilderV2 {
	b.opts.Sort = sort

	return b
}

func (b *manageChannelMembersBuilderV2) Set(channelMembersInput []PNChannelMembersSet) *manageChannelMembersBuilderV2 {
	b.opts.MembersSet = channelMembersInput

	return b
}

func (b *manageChannelMembersBuilderV2) Remove(channelMembersRemove []PNChannelMembersRemove) *manageChannelMembersBuilderV2 {
	b.opts.MembersRemove = channelMembersRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *manageChannelMembersBuilderV2) QueryParam(queryParam map[string]string) *manageChannelMembersBuilderV2 {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the manageMembers request.
func (b *manageChannelMembersBuilderV2) Transport(tr http.RoundTripper) *manageChannelMembersBuilderV2 {
	b.opts.Transport = tr
	return b
}

// Execute runs the manageMembers request.
func (b *manageChannelMembersBuilderV2) Execute() (*PNManageMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyManageMembersResponse, status, err
	}

	return newPNManageMembersResponse(rawJSON, b.opts, status)
}

type manageMembersOptsV2 struct {
	pubnub        *PubNub
	Channel       string
	Limit         int
	Include       []string
	Start         string
	End           string
	Filter        string
	Sort          []string
	Count         bool
	QueryParam    map[string]string
	MembersRemove []PNChannelMembersRemove
	MembersSet    []PNChannelMembersSet
	Transport     http.RoundTripper

	ctx Context
}

func (o *manageMembersOptsV2) config() Config {
	return *o.pubnub.Config
}

func (o *manageMembersOptsV2) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *manageMembersOptsV2) context() Context {
	return o.ctx
}

func (o *manageMembersOptsV2) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *manageMembersOptsV2) buildPath() (string, error) {
	return fmt.Sprintf(manageMembersPathV2,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *manageMembersOptsV2) buildQuery() (*url.Values, error) {

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

func (o *manageMembersOptsV2) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *manageMembersOptsV2) buildBody() ([]byte, error) {
	b := &PNManageChannelMembersBody{
		Set:    o.MembersSet,
		Remove: o.MembersRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *manageMembersOptsV2) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *manageMembersOptsV2) httpMethod() string {
	return "PATCH"
}

func (o *manageMembersOptsV2) isAuthRequired() bool {
	return true
}

func (o *manageMembersOptsV2) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *manageMembersOptsV2) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *manageMembersOptsV2) operationType() OperationType {
	return PNManageMembersOperation
}

func (o *manageMembersOptsV2) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNManageMembersResponse is the Objects API Response for ManageMembers
type PNManageMembersResponse struct {
	status     int                `json:"status"`
	Data       []PNChannelMembers `json:"data"`
	TotalCount int                `json:"totalCount"`
	Next       string             `json:"next"`
	Prev       string             `json:"prev"`
}

func newPNManageMembersResponse(jsonBytes []byte, o *manageMembersOptsV2,
	status StatusResponse) (*PNManageMembersResponse, StatusResponse, error) {

	resp := &PNManageMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyManageMembersResponse, status, e
	}

	return resp, status, nil
}
