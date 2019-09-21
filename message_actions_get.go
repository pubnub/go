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

var emptyPNGetMessageActionsResponse *PNGetMessageActionsResponse

const getMessageActionsPath = "/v1/actions/%s/channel/%s/message/%s"

type getMessageActionsBuilder struct {
	opts *getMessageActionsOpts
}

func newGetMessageActionsBuilder(pubnub *PubNub) *getMessageActionsBuilder {
	builder := getMessageActionsBuilder{
		opts: &getMessageActionsOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newGetMessageActionsBuilderWithContext(pubnub *PubNub,
	context Context) *getMessageActionsBuilder {
	builder := getMessageActionsBuilder{
		opts: &getMessageActionsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type getMessageActionsBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMessageActionsBuilder) Include(include []PNUserSpaceInclude) *getMessageActionsBuilder {

	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMessageActionsBuilder) ID(id string) *getMessageActionsBuilder {
	b.opts.ID = id

	return b
}

// Auth sets the Authorization key with permissions to perform the request.
func (b *getMessageActionsBuilder) Name(name string) *getMessageActionsBuilder {
	b.opts.Name = name

	return b
}

func (b *getMessageActionsBuilder) Description(description string) *getMessageActionsBuilder {
	b.opts.Description = description

	return b
}

func (b *getMessageActionsBuilder) Custom(custom map[string]interface{}) *getMessageActionsBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMessageActionsBuilder) QueryParam(queryParam map[string]string) *getMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getMessageActions request.
func (b *getMessageActionsBuilder) Transport(tr http.RoundTripper) *getMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getMessageActions request.
func (b *getMessageActionsBuilder) Execute() (*PNGetMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNGetMessageActionsResponse, status, err
	}

	return newPNGetMessageActionsResponse(rawJSON, b.opts, status)
}

type getMessageActionsOpts struct {
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

func (o *getMessageActionsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *getMessageActionsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *getMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *getMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(getMessageActionsPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *getMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *getMessageActionsOpts) buildBody() ([]byte, error) {
	b := &getMessageActionsBody{
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

func (o *getMessageActionsOpts) httpMethod() string {
	return "POST"
}

func (o *getMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *getMessageActionsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *getMessageActionsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *getMessageActionsOpts) operationType() OperationType {
	return PNGetMessageActionsOperation
}

func (o *getMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNGetMessageActionsResponse is the Objects API Response
type PNGetMessageActionsResponse struct {
	status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNGetMessageActionsResponse(jsonBytes []byte, o *getMessageActionsOpts,
	status StatusResponse) (*PNGetMessageActionsResponse, StatusResponse, error) {

	resp := &PNGetMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNGetMessageActionsResponse, status, e
	}

	return resp, status, nil
}
