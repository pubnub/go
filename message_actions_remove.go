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

var emptyPNRemoveMessageActionsResponse *PNRemoveMessageActionsResponse

const removeMessageActionsPath = "/v1/actions/%s/channel/%s/message/%s/action/%s"

type removeMessageActionsBuilder struct {
	opts *removeMessageActionsOpts
}

func newRemoveMessageActionsBuilder(pubnub *PubNub) *removeMessageActionsBuilder {
	builder := removeMessageActionsBuilder{
		opts: &removeMessageActionsOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newRemoveMessageActionsBuilderWithContext(pubnub *PubNub,
	context Context) *removeMessageActionsBuilder {
	builder := removeMessageActionsBuilder{
		opts: &removeMessageActionsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type removeMessageActionsBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *removeMessageActionsBuilder) Include(include []PNUserSpaceInclude) *removeMessageActionsBuilder {

	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *removeMessageActionsBuilder) ID(id string) *removeMessageActionsBuilder {
	b.opts.ID = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *removeMessageActionsBuilder) Name(name string) *removeMessageActionsBuilder {
	b.opts.Name = name

	return b
}

func (b *removeMessageActionsBuilder) Description(description string) *removeMessageActionsBuilder {
	b.opts.Description = description

	return b
}

func (b *removeMessageActionsBuilder) Custom(custom map[string]interface{}) *removeMessageActionsBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeMessageActionsBuilder) QueryParam(queryParam map[string]string) *removeMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeMessageActions request.
func (b *removeMessageActionsBuilder) Transport(tr http.RoundTripper) *removeMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeMessageActions request.
func (b *removeMessageActionsBuilder) Execute() (*PNRemoveMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNRemoveMessageActionsResponse, status, err
	}

	return newPNRemoveMessageActionsResponse(rawJSON, b.opts, status)
}

type removeMessageActionsOpts struct {
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

func (o *removeMessageActionsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *removeMessageActionsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *removeMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *removeMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeMessageActionsPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *removeMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *removeMessageActionsOpts) buildBody() ([]byte, error) {
	b := &removeMessageActionsBody{
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

func (o *removeMessageActionsOpts) httpMethod() string {
	return "POST"
}

func (o *removeMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *removeMessageActionsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *removeMessageActionsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *removeMessageActionsOpts) operationType() OperationType {
	return PNRemoveMessageActionsOperation
}

func (o *removeMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNRemoveMessageActionsResponse is the Objects API Response for create space
type PNRemoveMessageActionsResponse struct {
	status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNRemoveMessageActionsResponse(jsonBytes []byte, o *removeMessageActionsOpts,
	status StatusResponse) (*PNRemoveMessageActionsResponse, StatusResponse, error) {

	resp := &PNRemoveMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNRemoveMessageActionsResponse, status, e
	}

	return resp, status, nil
}
