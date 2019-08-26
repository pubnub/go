package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

var emptyPNUpdateSpaceResponse *PNUpdateSpaceResponse

const updateSpacePath = "/v1/objects/%s/spaces/%s"

type updateSpaceBuilder struct {
	opts *updateSpaceOpts
}

func newUpdateSpaceBuilder(pubnub *PubNub) *updateSpaceBuilder {
	builder := updateSpaceBuilder{
		opts: &updateSpaceOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newUpdateSpaceBuilderWithContext(pubnub *PubNub,
	context Context) *updateSpaceBuilder {
	builder := updateSpaceBuilder{
		opts: &updateSpaceOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type UpdateSpaceBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceBuilder) Include(include []PNUserSpaceInclude) *updateSpaceBuilder {
	b.opts.Include = utils.EnumArrayToStringArray(fmt.Sprint(include))

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceBuilder) ID(id string) *updateSpaceBuilder {
	b.opts.ID = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *updateSpaceBuilder) Name(name string) *updateSpaceBuilder {
	b.opts.Name = name

	return b
}

func (b *updateSpaceBuilder) Description(description string) *updateSpaceBuilder {
	b.opts.Description = description

	return b
}

func (b *updateSpaceBuilder) Custom(custom map[string]interface{}) *updateSpaceBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *updateSpaceBuilder) QueryParam(queryParam map[string]string) *updateSpaceBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the updateSpace request.
func (b *updateSpaceBuilder) Transport(tr http.RoundTripper) *updateSpaceBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the updateSpace request.
func (b *updateSpaceBuilder) Execute() (*PNUpdateSpaceResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNUpdateSpaceResponse, status, err
	}

	return newPNUpdateSpaceResponse(rawJSON, b.opts, status)
}

type updateSpaceOpts struct {
	pubnub      *PubNub
	Include     []string
	ID          string
	Name        string
	Description string
	Custom      map[string]interface{}
	QueryParam  map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *updateSpaceOpts) config() Config {
	return *o.pubnub.Config
}

func (o *updateSpaceOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *updateSpaceOpts) context() Context {
	return o.ctx
}

func (o *updateSpaceOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *updateSpaceOpts) buildPath() (string, error) {
	return fmt.Sprintf(updateSpacePath,
		o.pubnub.Config.SubscribeKey, o.ID), nil
}

func (o *updateSpaceOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *updateSpaceOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *updateSpaceOpts) buildBody() ([]byte, error) {
	b := &UpdateSpaceBody{
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

func (o *updateSpaceOpts) httpMethod() string {
	return "PATCH"
}

func (o *updateSpaceOpts) isAuthRequired() bool {
	return true
}

func (o *updateSpaceOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *updateSpaceOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *updateSpaceOpts) operationType() OperationType {
	return PNUpdateSpaceOperation
}

func (o *updateSpaceOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type PNUpdateSpaceResponse struct {
	Status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNUpdateSpaceResponse(jsonBytes []byte, o *updateSpaceOpts,
	status StatusResponse) (*PNUpdateSpaceResponse, StatusResponse, error) {

	resp := &PNUpdateSpaceResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNUpdateSpaceResponse, status, e
	}

	return resp, status, nil
}
