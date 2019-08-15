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

var emptyGetSpaceMembershipsResponse *PNGetSpaceMembershipsResponse

const getSpaceMembershipsPath = "/v1/objects/%s/users/%s/spaces"

const spaceMembershipLimit = 100

type getSpaceMembershipsBuilder struct {
	opts *getSpaceMembershipsOpts
}

func newGetSpaceMembershipsBuilder(pubnub *PubNub) *getSpaceMembershipsBuilder {
	builder := getSpaceMembershipsBuilder{
		opts: &getSpaceMembershipsOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceMembershipLimit

	return &builder
}

func newGetSpaceMembershipsBuilderWithContext(pubnub *PubNub,
	context Context) *getSpaceMembershipsBuilder {
	builder := getSpaceMembershipsBuilder{
		opts: &getSpaceMembershipsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *getSpaceMembershipsBuilder) Auth(auth string) *getSpaceMembershipsBuilder {
// 	//b.opts.Auth = auth

// 	return b
// }

func (b *getSpaceMembershipsBuilder) UserId(id string) *getSpaceMembershipsBuilder {
	b.opts.Id = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getSpaceMembershipsBuilder) Include(include []PNSpaceMembershipsIncude) *getSpaceMembershipsBuilder {
	b.opts.Include = utils.EnumArrayToStringArray(fmt.Sprint(include))

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getSpaceMembershipsBuilder) Limit(limit int) *getSpaceMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getSpaceMembershipsBuilder) Start(start string) *getSpaceMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *getSpaceMembershipsBuilder) End(end string) *getSpaceMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *getSpaceMembershipsBuilder) Count(count bool) *getSpaceMembershipsBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getSpaceMembershipsBuilder) QueryParam(queryParam map[string]string) *getSpaceMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getSpaceMemberships request.
func (b *getSpaceMembershipsBuilder) Transport(tr http.RoundTripper) *getSpaceMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getSpaceMemberships request.
func (b *getSpaceMembershipsBuilder) Execute() (*PNGetSpaceMembershipsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetSpaceMembershipsResponse, status, err
	}

	return newPNGetSpaceMembershipsResponse(rawJSON, b.opts, status)
}

type getSpaceMembershipsOpts struct {
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

func (o *getSpaceMembershipsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getSpaceMembershipsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getSpaceMembershipsOpts) context() Context {
	return o.ctx
}

func (o *getSpaceMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getSpaceMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(getSpaceMembershipsPath,
		o.pubnub.Config.SubscribeKey, o.Id), nil
}

func (o *getSpaceMembershipsOpts) buildQuery() (*url.Values, error) {

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

func (o *getSpaceMembershipsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getSpaceMembershipsOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getSpaceMembershipsOpts) httpMethod() string {
	return "GET"
}

func (o *getSpaceMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *getSpaceMembershipsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getSpaceMembershipsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getSpaceMembershipsOpts) operationType() OperationType {
	return PNGetSpaceMembershipsOperation
}

func (o *getSpaceMembershipsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNGetSpaceMembershipsResponse struct {
	Status     int                 `json:"status"`
	Data       []PNSpaceMembership `json:"data"`
	TotalCount int                 `json:"totalCount"`
	Next       string              `json:"next"`
	Prev       string              `json:"prev"`
}

func newPNGetSpaceMembershipsResponse(jsonBytes []byte, o *getSpaceMembershipsOpts,
	status StatusResponse) (*PNGetSpaceMembershipsResponse, StatusResponse, error) {

	resp := &PNGetSpaceMembershipsResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetSpaceMembershipsResponse, status, e
	}

	return resp, status, nil
}
