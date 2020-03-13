package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/pnerr"
)

var emptyGetMembersResponse *PNGetMembersResponse

const getMembersPath = "/v1/objects/%s/spaces/%s/users"

const membersLimit = 100

type getMembersBuilder struct {
	opts *getMembersOpts
}

func newGetMembersBuilder(pubnub *PubNub) *getMembersBuilder {
	builder := getMembersBuilder{
		opts: &getMembersOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = membersLimit

	return &builder
}

func newGetMembersBuilderWithContext(pubnub *PubNub,
	context Context) *getMembersBuilder {
	builder := getMembersBuilder{
		opts: &getMembersOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getMembersBuilder) SpaceID(id string) *getMembersBuilder {
	b.opts.ID = id

	return b
}

func (b *getMembersBuilder) Include(include []PNMembersInclude) *getMembersBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getMembersBuilder) Limit(limit int) *getMembersBuilder {
	b.opts.Limit = limit

	return b
}

func (b *getMembersBuilder) Start(start string) *getMembersBuilder {
	b.opts.Start = start

	return b
}

func (b *getMembersBuilder) End(end string) *getMembersBuilder {
	b.opts.End = end

	return b
}

func (b *getMembersBuilder) Filter(filter string) *getMembersBuilder {
	b.opts.Filter = filter

	return b
}

func (b *getMembersBuilder) Sort(sort []string) *getMembersBuilder {
	b.opts.Sort = sort

	return b
}

func (b *getMembersBuilder) Count(count bool) *getMembersBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMembersBuilder) QueryParam(queryParam map[string]string) *getMembersBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getMembers request.
func (b *getMembersBuilder) Transport(tr http.RoundTripper) *getMembersBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getMembers request.
func (b *getMembersBuilder) Execute() (*PNGetMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetMembersResponse, status, err
	}

	return newPNGetMembersResponse(rawJSON, b.opts, status)
}

type getMembersOpts struct {
	pubnub     *PubNub
	ID         string
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

func (o *getMembersOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getMembersOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getMembersOpts) context() Context {
	return o.ctx
}

func (o *getMembersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getMembersOpts) buildPath() (string, error) {
	return fmt.Sprintf(getMembersPath,
		o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getMembersOpts) buildQuery() (*url.Values, error) {

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

	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getMembersOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getMembersOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getMembersOpts) httpMethod() string {
	return "GET"
}

func (o *getMembersOpts) isAuthRequired() bool {
	return true
}

func (o *getMembersOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getMembersOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getMembersOpts) operationType() OperationType {
	return PNGetMembersOperation
}

func (o *getMembersOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetMembersResponse is the Objects API Response for Get Members
type PNGetMembersResponse struct {
	status     int         `json:"status"`
	Data       []PNMembers `json:"data"`
	TotalCount int         `json:"totalCount"`
	Next       string      `json:"next"`
	Prev       string      `json:"prev"`
}

func newPNGetMembersResponse(jsonBytes []byte, o *getMembersOpts,
	status StatusResponse) (*PNGetMembersResponse, StatusResponse, error) {

	resp := &PNGetMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetMembersResponse, status, e
	}

	return resp, status, nil
}
