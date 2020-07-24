package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"

	"net/http"
	"net/url"
	"strconv"
)

var emptySignalResponse *SignalResponse

const signalGetPath = "/signal/%s/%s/0/%s/%s/%s"
const signalPostPath = "/signal/%s/%s/0/%s/%s"

type signalBuilder struct {
	opts *signalOpts
}

func newSignalBuilder(pubnub *PubNub) *signalBuilder {
	builder := signalBuilder{
		opts: &signalOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.UsePost = false

	return &builder
}

func newSignalBuilderWithContext(pubnub *PubNub,
	context Context) *signalBuilder {
	builder := signalBuilder{
		opts: &signalOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channel sets the Channel for the Signal request.
func (b *signalBuilder) Channel(channel string) *signalBuilder {
	b.opts.Channel = channel
	return b
}

// Message sets the Payload for the Signal request.
func (b *signalBuilder) Message(msg interface{}) *signalBuilder {
	b.opts.Message = msg

	return b
}

// usePost sends the Signal request using HTTP POST. Not implemented
func (b *signalBuilder) usePost(post bool) *signalBuilder {
	b.opts.UsePost = post

	return b
}

// Transport sets the Transport for the objectAPICreateUsers request.
func (b *signalBuilder) Transport(tr http.RoundTripper) *signalBuilder {
	b.opts.Transport = tr
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *signalBuilder) QueryParam(queryParam map[string]string) *signalBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Signal request.
func (b *signalBuilder) Execute() (*SignalResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySignalResponse, status, err
	}

	return newSignalResponse(rawJSON, b.opts, status)
}

type signalOpts struct {
	pubnub     *PubNub
	Message    interface{}
	Channel    string
	UsePost    bool
	QueryParam map[string]string
	Transport  http.RoundTripper
	ctx        Context
}

func (o *signalOpts) config() Config {
	return *o.pubnub.Config
}

func (o *signalOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *signalOpts) context() Context {
	return o.ctx
}

func (o *signalOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	return nil
}

func (o *signalOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(signalPostPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0"), nil
	}

	var msg string
	jsonEncBytes, errEnc := json.Marshal(o.Message)
	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
		return "", errEnc
	}
	msg = string(jsonEncBytes)
	return fmt.Sprintf(signalGetPath,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.Channel),
		"0",
		utils.URLEncode(msg),
	), nil
}

func (o *signalOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *signalOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *signalOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		jsonEncBytes, errEnc := json.Marshal(o.Message)
		if errEnc != nil {
			o.pubnub.Config.Log.Printf("ERROR: Signal error: %s\n", errEnc.Error())
			return []byte{}, errEnc
		}
		return jsonEncBytes, nil
	}
	return []byte{}, nil
}

func (o *signalOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *signalOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
}

func (o *signalOpts) isAuthRequired() bool {
	return true
}

func (o *signalOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *signalOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *signalOpts) operationType() OperationType {
	return PNSignalOperation
}

func (o *signalOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// SignalResponse is the response to Signal request.
type SignalResponse struct {
	Timestamp int64
}

func newSignalResponse(jsonBytes []byte, o *signalOpts,
	status StatusResponse) (*SignalResponse, StatusResponse, error) {

	resp := &SignalResponse{}

	var value []interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySignalResponse, status, e
	}

	if len(value) > 1 {
		timeString, ok := value[2].(string)
		if !ok {
			return emptySignalResponse, status, pnerr.NewResponseParsingError(fmt.Sprintf("Error unmarshalling response 2, %s %v", value[2], value), nil, nil)
		}
		timestamp, err := strconv.ParseInt(timeString, 10, 64)
		if err != nil {
			return emptySignalResponse, status, err
		}

		return &SignalResponse{
			Timestamp: timestamp,
		}, status, nil
	}

	return resp, status, nil
}
