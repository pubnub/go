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
)

var emptyGetSpacesResponse *PNGetSpacesResponse

const getSpacesPath = "/v1/objects/%s/spaces"

const spaceLimit = 100

type getSpacesBuilder struct {
	opts *getSpacesOpts
}

func newGetSpacesBuilder(pubnub *PubNub) *getSpacesBuilder {
	builder := getSpacesBuilder{
		opts: &getSpacesOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.Limit = spaceLimit

	return &builder
}

func newGetSpacesBuilderWithContext(pubnub *PubNub,
	context Context) *getSpacesBuilder {
	builder := getSpacesBuilder{
		opts: &getSpacesOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getSpacesBuilder) Include(include []PNUserSpaceInclude) *getSpacesBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getSpacesBuilder) Limit(limit int) *getSpacesBuilder {
	b.opts.Limit = limit

	return b
}

func (b *getSpacesBuilder) Start(start string) *getSpacesBuilder {
	b.opts.Start = start

	return b
}

func (b *getSpacesBuilder) End(end string) *getSpacesBuilder {
	b.opts.End = end

	return b
}

func (b *getSpacesBuilder) Filter(filter string) *getSpacesBuilder {
	b.opts.Filter = filter

	return b
}

func (b *getSpacesBuilder) Count(count bool) *getSpacesBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getSpacesBuilder) QueryParam(queryParam map[string]string) *getSpacesBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getSpaces request.
func (b *getSpacesBuilder) Transport(tr http.RoundTripper) *getSpacesBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getSpaces request.
func (b *getSpacesBuilder) Execute() (*PNGetSpacesResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetSpacesResponse, status, err
	}

	return newPNGetSpacesResponse(rawJSON, b.opts, status)
}

type getSpacesOpts struct {
	pubnub *PubNub

	Limit      int
	Include    []string
	Start      string
	End        string
	Filter     string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getSpacesOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getSpacesOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getSpacesOpts) context() Context {
	return o.ctx
}

func (o *getSpacesOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getSpacesOpts) buildPath() (string, error) {
	return fmt.Sprintf(getSpacesPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *getSpacesOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
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
		q.Set("filter", utils.URLEncode(o.Filter))
	}

	o.pubnub.tokenManager.SetAuthParan(q, "", PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getSpacesOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getSpacesOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getSpacesOpts) httpMethod() string {
	return "GET"
}

func (o *getSpacesOpts) isAuthRequired() bool {
	return true
}

func (o *getSpacesOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getSpacesOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getSpacesOpts) operationType() OperationType {
	return PNGetSpacesOperation
}

func (o *getSpacesOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetSpacesResponse is the Objects API Response for Get Spaces
type PNGetSpacesResponse struct {
	status     int       `json:"status"`
	Data       []PNSpace `json:"data"`
	TotalCount int       `json:"totalCount"`
	Next       string    `json:"next"`
	Prev       string    `json:"prev"`
}

func newPNGetSpacesResponse(jsonBytes []byte, o *getSpacesOpts,
	status StatusResponse) (*PNGetSpacesResponse, StatusResponse, error) {

	resp := &PNGetSpacesResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetSpacesResponse, status, e
	}

	return resp, status, nil
}
