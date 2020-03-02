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

var emptyPNCreateSpaceResponse *PNCreateSpaceResponse

const createSpacePath = "/v1/objects/%s/spaces"

type createSpaceBuilder struct {
	opts *createSpaceOpts
}

func newCreateSpaceBuilder(pubnub *PubNub) *createSpaceBuilder {
	builder := createSpaceBuilder{
		opts: &createSpaceOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newCreateSpaceBuilderWithContext(pubnub *PubNub,
	context Context) *createSpaceBuilder {
	builder := createSpaceBuilder{
		opts: &createSpaceOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type createSpaceBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

func (b *createSpaceBuilder) Include(include []PNUserSpaceInclude) *createSpaceBuilder {

	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *createSpaceBuilder) ID(id string) *createSpaceBuilder {
	b.opts.ID = id

	return b
}

func (b *createSpaceBuilder) Name(name string) *createSpaceBuilder {
	b.opts.Name = name

	return b
}

func (b *createSpaceBuilder) Description(description string) *createSpaceBuilder {
	b.opts.Description = description

	return b
}

func (b *createSpaceBuilder) Custom(custom map[string]interface{}) *createSpaceBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *createSpaceBuilder) QueryParam(queryParam map[string]string) *createSpaceBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the createSpace request.
func (b *createSpaceBuilder) Transport(tr http.RoundTripper) *createSpaceBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the createSpace request.
func (b *createSpaceBuilder) Execute() (*PNCreateSpaceResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNCreateSpaceResponse, status, err
	}

	return newPNCreateSpaceResponse(rawJSON, b.opts, status)
}

type createSpaceOpts struct {
	pubnub *PubNub

	Include     []string
	ID          string
	Name        string
	Description string
	Custom      map[string]interface{}
	QueryParam  map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *createSpaceOpts) config() Config {
	return *o.pubnub.Config
}

func (o *createSpaceOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *createSpaceOpts) context() Context {
	return o.ctx
}

func (o *createSpaceOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *createSpaceOpts) buildPath() (string, error) {
	return fmt.Sprintf(createSpacePath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *createSpaceOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		SetArrayTypeQueryParam(q, o.Include, "include")
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *createSpaceOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *createSpaceOpts) buildBody() ([]byte, error) {
	b := &createSpaceBody{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		Custom:      o.Custom,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *createSpaceOpts) httpMethod() string {
	return "POST"
}

func (o *createSpaceOpts) isAuthRequired() bool {
	return true
}

func (o *createSpaceOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *createSpaceOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *createSpaceOpts) operationType() OperationType {
	return PNCreateSpaceOperation
}

func (o *createSpaceOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNCreateSpaceResponse is the Objects API Response for create space
type PNCreateSpaceResponse struct {
	status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNCreateSpaceResponse(jsonBytes []byte, o *createSpaceOpts,
	status StatusResponse) (*PNCreateSpaceResponse, StatusResponse, error) {

	resp := &PNCreateSpaceResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNCreateSpaceResponse, status, e
	}

	return resp, status, nil
}
