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

var emptyPNGetUsersResponse *PNGetUsersResponse

const getUsersPath = "/v1/objects/%s/users"

const usersLimit = 100

type getUsersBuilder struct {
	opts *getUsersOpts
}

func newGetUsersBuilder(pubnub *PubNub) *getUsersBuilder {
	builder := getUsersBuilder{
		opts: &getUsersOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = usersLimit

	return &builder
}

func newGetUsersBuilderWithContext(pubnub *PubNub,
	context Context) *getUsersBuilder {
	builder := getUsersBuilder{
		opts: &getUsersOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *getUsersBuilder) Auth(auth string) *getUsersBuilder {
// 	//b.opts.Auth = auth

// 	return b
// }

// Auth sets the Authorization key with permissions to perform the request.
func (b *getUsersBuilder) Include(include []string) *getUsersBuilder {
	b.opts.Include = include

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getUsersBuilder) Limit(limit int) *getUsersBuilder {
	b.opts.Limit = limit

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getUsersBuilder) Start(start string) *getUsersBuilder {
	b.opts.Start = start

	return b
}

func (b *getUsersBuilder) End(end string) *getUsersBuilder {
	b.opts.End = end

	return b
}

func (b *getUsersBuilder) Count(count bool) *getUsersBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUsersBuilder) QueryParam(queryParam map[string]string) *getUsersBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getUsers request.
func (b *getUsersBuilder) Transport(tr http.RoundTripper) *getUsersBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getUsers request.
func (b *getUsersBuilder) Execute() (*PNGetUsersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetUsersResponse, status, err
	}

	return newGetUsersResponse(rawJSON, b.opts, status)
}

type getUsersOpts struct {
	pubnub *PubNub

	Limit      int
	Include    []string
	Start      string
	End        string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getUsersOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getUsersOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getUsersOpts) context() Context {
	return o.ctx
}

func (o *getUsersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getUsersOpts) buildPath() (string, error) {
	return fmt.Sprintf(getUsersPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *getUsersOpts) buildQuery() (*url.Values, error) {

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

func (o *getUsersOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getUsersOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getUsersOpts) httpMethod() string {
	return "GET"
}

func (o *getUsersOpts) isAuthRequired() bool {
	return true
}

func (o *getUsersOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getUsersOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getUsersOpts) operationType() OperationType {
	return PNGetUsersOperation
}

func (o *getUsersOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNGetUsersResponse struct {
	Status     int      `json:"status"`
	Data       []PNUser `json:"data"`
	TotalCount int      `json:"totalCount"`
	Next       string   `json:"next"`
	Prev       string   `json:"prev"`
}

func newGetUsersResponse(jsonBytes []byte, o *getUsersOpts,
	status StatusResponse) (*PNGetUsersResponse, StatusResponse, error) {

	resp := &PNGetUsersResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetUsersResponse, status, e
	}

	return resp, status, nil
}
