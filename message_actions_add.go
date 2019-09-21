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

var emptyPNAddMessageActionsResponse *PNAddMessageActionsResponse

const addMessageActionsPath = "/v1/actions/%s/channel/%s"

type addMessageActionsBuilder struct {
	opts *addMessageActionsOpts
}

func newAddMessageActionsBuilder(pubnub *PubNub) *addMessageActionsBuilder {
	builder := addMessageActionsBuilder{
		opts: &addMessageActionsOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newAddMessageActionsBuilderWithContext(pubnub *PubNub,
	context Context) *addMessageActionsBuilder {
	builder := addMessageActionsBuilder{
		opts: &addMessageActionsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type addMessageActionsBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *addMessageActionsBuilder) Include(include []PNUserSpaceInclude) *addMessageActionsBuilder {

	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *addMessageActionsBuilder) ID(id string) *addMessageActionsBuilder {
	b.opts.ID = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *addMessageActionsBuilder) Name(name string) *addMessageActionsBuilder {
	b.opts.Name = name

	return b
}

func (b *addMessageActionsBuilder) Description(description string) *addMessageActionsBuilder {
	b.opts.Description = description

	return b
}

func (b *addMessageActionsBuilder) Custom(custom map[string]interface{}) *addMessageActionsBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *addMessageActionsBuilder) QueryParam(queryParam map[string]string) *addMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the addMessageActions request.
func (b *addMessageActionsBuilder) Transport(tr http.RoundTripper) *addMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the addMessageActions request.
func (b *addMessageActionsBuilder) Execute() (*PNAddMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNAddMessageActionsResponse, status, err
	}

	return newPNAddMessageActionsResponse(rawJSON, b.opts, status)
}

type addMessageActionsOpts struct {
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

func (o *addMessageActionsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *addMessageActionsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *addMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *addMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *addMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(addMessageActionsPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *addMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *addMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *addMessageActionsOpts) buildBody() ([]byte, error) {
	b := &addMessageActionsBody{
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

func (o *addMessageActionsOpts) httpMethod() string {
	return "POST"
}

func (o *addMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *addMessageActionsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *addMessageActionsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *addMessageActionsOpts) operationType() OperationType {
	return PNAddMessageActionsOperation
}

func (o *addMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNAddMessageActionsResponse is the Objects API Response
type PNAddMessageActionsResponse struct {
	status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNAddMessageActionsResponse(jsonBytes []byte, o *addMessageActionsOpts,
	status StatusResponse) (*PNAddMessageActionsResponse, StatusResponse, error) {

	resp := &PNAddMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNAddMessageActionsResponse, status, e
	}

	return resp, status, nil
}
