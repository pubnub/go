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
	"github.com/pubnub/go/utils"
	//"reflect"
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

// Auth sets the Authorization key with permissions to perform the request.
// func (b *getMembersBuilder) Auth(auth string) *getMembersBuilder {
// 	//b.opts.Auth = auth

// 	return b
// }

func (b *getMembersBuilder) SpaceId(id string) *getMembersBuilder {
	b.opts.Id = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMembersBuilder) Include(include []string) *getMembersBuilder {
	b.opts.Include = include

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMembersBuilder) Limit(limit int) *getMembersBuilder {
	b.opts.Limit = limit

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMembersBuilder) Start(start string) *getMembersBuilder {
	b.opts.Start = start

	return b
}

func (b *getMembersBuilder) End(end string) *getMembersBuilder {
	b.opts.End = end

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
	Id         string
	Limit      int
	Include    []string
	Start      string
	End        string
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
		o.pubnub.Config.SubscribeKey, o.Id), nil
}

func (o *getMembersOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	// if o.Auth != "" {
	// 	q.Set("auth", o.Auth)
	// }

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

type PNGetMembersResponse struct {
	Status     int                 `json:"status"`
	Data       []PNSpaceMembership `json:"data"`
	TotalCount int                 `json:"totalCount"`
	Next       string              `json:"next"`
	Prev       string              `json:"prev"`
}

func newPNGetMembersResponse(jsonBytes []byte, o *getMembersOpts,
	status StatusResponse) (*PNGetMembersResponse, StatusResponse, error) {

	resp := &PNGetMembersResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetMembersResponse, status, e
	}

	return resp, status, nil
}
