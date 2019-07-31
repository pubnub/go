package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
	//"reflect"

	"net/http"
	"net/url"
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
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
// func (b *createSpaceBuilder) Auth(auth string) *createSpaceBuilder {
// 	b.opts.Auth = auth

// 	return b
// }

// Auth sets the Authorization key with permissions to perform the request.
func (b *createSpaceBuilder) Include(include []string) *createSpaceBuilder {
	b.opts.Include = include

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *createSpaceBuilder) Id(id string) *createSpaceBuilder {
	b.opts.Id = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
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

	Auth        string
	Include     []string
	Id          string
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
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	if o.Auth != "" {
		q.Set("auth", o.Auth)
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *createSpaceOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *createSpaceOpts) buildBody() ([]byte, error) {
	b := &createSpaceBody{
		Id:          o.Id,
		Name:        o.Name,
		Description: o.Description,
		Custom:      o.Custom,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	fmt.Println(fmt.Sprintf("%v %s", b, string(jsonEncBytes)))
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

type PNCreateSpaceResponse struct {
	Status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNCreateSpaceResponse(jsonBytes []byte, o *createSpaceOpts,
	status StatusResponse) (*PNCreateSpaceResponse, StatusResponse, error) {

	resp := &PNCreateSpaceResponse{}

	fmt.Println(string(jsonBytes))

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("error", err)
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNCreateSpaceResponse, status, e
	}

	return resp, status, nil
}
