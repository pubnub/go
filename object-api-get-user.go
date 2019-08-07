package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
	//"reflect"
)

var emptyPNFetchUserResponse *PNFetchUserResponse

const fetchUserPath = "/v1/objects/%s/users/%s"

type getUserBuilder struct {
	opts *fetchUserOpts
}

func newGetUserBuilder(pubnub *PubNub) *getUserBuilder {
	builder := getUserBuilder{
		opts: &fetchUserOpts{
			pubnub: pubnub,
		},
	}
	return &builder
}

func newGetUserBuilderWithContext(pubnub *PubNub,
	context Context) *getUserBuilder {
	builder := getUserBuilder{
		opts: &fetchUserOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getUserBuilder) Include(include []string) *getUserBuilder {
	b.opts.Include = include

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getUserBuilder) Id(id string) *getUserBuilder {
	b.opts.Id = id

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUserBuilder) QueryParam(queryParam map[string]string) *getUserBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the fetchUser request.
func (b *getUserBuilder) Transport(tr http.RoundTripper) *getUserBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the fetchUser request.
func (b *getUserBuilder) Execute() (*PNFetchUserResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNFetchUserResponse, status, err
	}

	return newPNFetchUserResponse(rawJSON, b.opts, status)
}

type fetchUserOpts struct {
	pubnub     *PubNub
	Id         string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *fetchUserOpts) config() Config {
	return *o.pubnub.Config
}

func (o *fetchUserOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *fetchUserOpts) context() Context {
	return o.ctx
}

func (o *fetchUserOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *fetchUserOpts) buildPath() (string, error) {
	return fmt.Sprintf(fetchUserPath,
		o.pubnub.Config.SubscribeKey, o.Id), nil
}

func (o *fetchUserOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *fetchUserOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *fetchUserOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *fetchUserOpts) httpMethod() string {
	return "GET"
}

func (o *fetchUserOpts) isAuthRequired() bool {
	return true
}

func (o *fetchUserOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *fetchUserOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *fetchUserOpts) operationType() OperationType {
	return PNFetchUserOperation
}

func (o *fetchUserOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNFetchUserResponse struct {
	Status int    `json:"status"`
	Data   PNUser `json:"data"`
}

func newPNFetchUserResponse(jsonBytes []byte, o *fetchUserOpts,
	status StatusResponse) (*PNFetchUserResponse, StatusResponse, error) {

	resp := &PNFetchUserResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNFetchUserResponse, status, e
	}

	return resp, status, nil
}
