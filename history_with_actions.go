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

var emptyPNHistoryWithActionsResponse *PNHistoryWithActionsResponse

const historyWithActionsPath = "/v3/history-with-actions/%s/channel/%s"

type historyWithActionsBuilder struct {
	opts *historyWithActionsOpts
}

func newHistoryWithActionsBuilder(pubnub *PubNub) *historyWithActionsBuilder {
	builder := historyWithActionsBuilder{
		opts: &historyWithActionsOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newHistoryWithActionsBuilderWithContext(pubnub *PubNub,
	context Context) *historyWithActionsBuilder {
	builder := historyWithActionsBuilder{
		opts: &historyWithActionsOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

type historyWithActionsBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Custom      map[string]interface{} `json:"custom"`
}

func (b *historyWithActionsBuilder) Include(include []PNUserSpaceInclude) *historyWithActionsBuilder {

	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *historyWithActionsBuilder) ID(id string) *historyWithActionsBuilder {
	b.opts.ID = id

	return b
}

func (b *historyWithActionsBuilder) Name(name string) *historyWithActionsBuilder {
	b.opts.Name = name

	return b
}

func (b *historyWithActionsBuilder) Description(description string) *historyWithActionsBuilder {
	b.opts.Description = description

	return b
}

func (b *historyWithActionsBuilder) Custom(custom map[string]interface{}) *historyWithActionsBuilder {
	b.opts.Custom = custom

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *historyWithActionsBuilder) QueryParam(queryParam map[string]string) *historyWithActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the historyWithActions request.
func (b *historyWithActionsBuilder) Transport(tr http.RoundTripper) *historyWithActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the historyWithActions request.
func (b *historyWithActionsBuilder) Execute() (*PNHistoryWithActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPNHistoryWithActionsResponse, status, err
	}

	return newPNHistoryWithActionsResponse(rawJSON, b.opts, status)
}

type historyWithActionsOpts struct {
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

func (o *historyWithActionsOpts) config() Config {
	return *o.pubnub.Config
}

func (o *historyWithActionsOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *historyWithActionsOpts) context() Context {
	return o.ctx
}

func (o *historyWithActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *historyWithActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(historyWithActionsPath,
		o.pubnub.Config.SubscribeKey), nil
}

func (o *historyWithActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Include != nil {
		q.Set("include", string(utils.JoinChannels(o.Include)))
	}
	o.pubnub.tokenManager.SetAuthParan(q, o.ID, PNSpaces)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *historyWithActionsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *historyWithActionsOpts) buildBody() ([]byte, error) {
	b := &historyWithActionsBody{
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

func (o *historyWithActionsOpts) httpMethod() string {
	return "POST"
}

func (o *historyWithActionsOpts) isAuthRequired() bool {
	return true
}

func (o *historyWithActionsOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *historyWithActionsOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *historyWithActionsOpts) operationType() OperationType {
	return PNHistoryWithActionsOperation
}

func (o *historyWithActionsOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNHistoryWithActionsResponse is the Objects API Response for create space
type PNHistoryWithActionsResponse struct {
	status int     `json:"status"`
	Data   PNSpace `json:"data"`
}

func newPNHistoryWithActionsResponse(jsonBytes []byte, o *historyWithActionsOpts,
	status StatusResponse) (*PNHistoryWithActionsResponse, StatusResponse, error) {

	resp := &PNHistoryWithActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPNHistoryWithActionsResponse, status, e
	}

	return resp, status, nil
}
