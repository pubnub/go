package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

var emptyPNAddMessageActionsResponse *PNAddMessageActionsResponse

const addMessageActionsPath = "/v1/message-actions/%s/channel/%s/message/%s"

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

// MessageAction struct is used to create a Message Action
type MessageAction struct {
	ActionType  string `json:"type"`
	ActionValue string `json:"value"`
}

func (b *addMessageActionsBuilder) Channel(channel string) *addMessageActionsBuilder {
	b.opts.Channel = channel

	return b
}

func (b *addMessageActionsBuilder) MessageTimetoken(timetoken string) *addMessageActionsBuilder {
	b.opts.MessageTimetoken = timetoken

	return b
}

func (b *addMessageActionsBuilder) Action(action MessageAction) *addMessageActionsBuilder {
	b.opts.Action = action

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

	Channel          string
	MessageTimetoken string
	Action           MessageAction
	QueryParam       map[string]string

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
		o.pubnub.Config.SubscribeKey, o.Channel, o.MessageTimetoken), nil
}

func (o *addMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *addMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *addMessageActionsOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Action)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *addMessageActionsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
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

// PNMessageActionsResponse Message Actions response.
type PNMessageActionsResponse struct {
	ActionType       string `json:"type"`
	ActionValue      string `json:"value"`
	ActionTimetoken  string `json:"actionTimetoken"`
	MessageTimetoken string `json:"messageTimetoken"`
	UUID             string `json:"uuid"`
}

// PNAddMessageActionsResponse is the Add Message Actions API Response
type PNAddMessageActionsResponse struct {
	status int                      `json:"status"`
	Data   PNMessageActionsResponse `json:"data"`
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
