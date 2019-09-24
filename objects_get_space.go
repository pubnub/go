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
)

var emptyPNGetSpaceResponse *PNGetSpaceResponse

const getSpacePath = "/v1/objects/%s/spaces/%s"

type getSpaceBuilder struct {
	opts *getSpaceOpts
}

func newGetSpaceBuilder(pubnub *PubNub) *getSpaceBuilder {
	builder := getSpaceBuilder{
		opts: &getSpaceOpts{
			pubnub: pubnub,
		},
	}
	return &builder
}

func newGetSpaceBuilderWithContext(pubnub *PubNub,
	context Context) *getSpaceBuilder {
	builder := getSpaceBuilder{
		opts: &getSpaceOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *getSpaceBuilder) Include(include []PNUserSpaceInclude) *getSpaceBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getSpaceBuilder) ID(id string) *getSpaceBuilder {
	b.opts.ID = id

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getSpaceBuilder) QueryParam(queryParam map[string]string) *getSpaceBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getSpace request.
func (b *getSpaceBuilder) Transport(tr http.RoundTripper) *getSpaceBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getSpace request.
func (b *getSpaceBuilder) Execute() (*PNGetSpaceResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetSpaceResponse, status, err
	}

	return newPNGetSpaceResponse(rawJSON, b.opts, status)
}

type getSpaceOpts struct {
	pubnub     *PubNub
	ID         string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getSpaceOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getSpaceOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getSpaceOpts) context() Context {
	return o.ctx
}

func (o *getSpaceOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getSpaceOpts) buildPath() (string, error) {
	return fmt.Sprintf(getSpacePath,
		o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *getSpaceOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getSpaceOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getSpaceOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getSpaceOpts) httpMethod() string {
	return "GET"
}

func (o *getSpaceOpts) isAuthRequired() bool {
	return true
}

func (o *getSpaceOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getSpaceOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getSpaceOpts) operationType() OperationType {
	return PNGetSpaceOperation
}

func (o *getSpaceOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetSpaceResponse is the Objects API Response for Get Space
type PNGetSpaceResponse struct {
	status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNGetSpaceResponse(jsonBytes []byte, o *getSpaceOpts,
	status StatusResponse) (*PNGetSpaceResponse, StatusResponse, error) {

	resp := &PNGetSpaceResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetSpaceResponse, status, e
	}

	return resp, status, nil
}
