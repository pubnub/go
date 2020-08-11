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

var emptyPublishFileMessageResponse *PublishFileMessageResponse

const publishFileMessageGetPath = "/v1/files/publish-file/%s/%s/0/%s/%s/%s"
const publishFileMessagePostPath = "/v1/files/publish-file/%s/%s/0/%s/%s"

type publishFileMessageBuilder struct {
	opts *publishFileMessageOpts
}

func newPublishFileMessageBuilder(pubnub *PubNub) *publishFileMessageBuilder {
	builder := publishFileMessageBuilder{
		opts: &publishFileMessageOpts{
			pubnub: pubnub,
		},
	}
	builder.opts.UsePost = false

	return &builder
}

func newPublishFileMessageBuilderWithContext(pubnub *PubNub,
	context Context) *publishFileMessageBuilder {
	builder := publishFileMessageBuilder{
		opts: &publishFileMessageOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// TTL sets the TTL (hours) for the Publish request.
func (b *publishFileMessageBuilder) TTL(ttl int) *publishFileMessageBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

// Meta sets the Meta Payload for the Publish request.
func (b *publishFileMessageBuilder) Meta(meta interface{}) *publishFileMessageBuilder {
	b.opts.Meta = meta

	return b
}

// ShouldStore if true the messages are stored in History
func (b *publishFileMessageBuilder) ShouldStore(store bool) *publishFileMessageBuilder {
	b.opts.ShouldStore = store
	b.opts.setShouldStore = true
	return b
}

// Channel sets the Channel for the PublishFileMessage request.
func (b *publishFileMessageBuilder) Channel(channel string) *publishFileMessageBuilder {
	b.opts.Channel = channel
	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) FileName(name string) *publishFileMessageBuilder {
	b.opts.FileName = name

	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) Message(msg PNPublishFileMessage) *publishFileMessageBuilder {
	b.opts.Message = msg

	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) MessageText(msg string) *publishFileMessageBuilder {
	b.opts.MessageText = msg

	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) FileID(id string) *publishFileMessageBuilder {
	b.opts.FileID = id

	return b
}

// usePost sends the PublishFileMessage request using HTTP POST. Not implemented
func (b *publishFileMessageBuilder) usePost(post bool) *publishFileMessageBuilder {
	b.opts.UsePost = post

	return b
}

// Transport sets the Transport for the objectAPICreateUsers request.
func (b *publishFileMessageBuilder) Transport(tr http.RoundTripper) *publishFileMessageBuilder {
	b.opts.Transport = tr
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *publishFileMessageBuilder) QueryParam(queryParam map[string]string) *publishFileMessageBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the PublishFileMessage request.
func (b *publishFileMessageBuilder) Execute() (*PublishFileMessageResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishFileMessageResponse, status, err
	}

	return newPublishFileMessageResponse(rawJSON, b.opts, status)
}

type publishFileMessageOpts struct {
	pubnub         *PubNub
	Message        interface{}
	Channel        string
	UsePost        bool
	TTL            int
	Meta           interface{}
	ShouldStore    bool
	setTTL         bool
	setShouldStore bool
	MessageText    string
	FileID         string
	FileName       string
	QueryParam     map[string]string
	Transport      http.RoundTripper
	ctx            Context
}

func (o *publishFileMessageOpts) config() Config {
	return *o.pubnub.Config
}

func (o *publishFileMessageOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *publishFileMessageOpts) context() Context {
	return o.ctx
}

func (o *publishFileMessageOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if (o.Message == nil) && (o.FileID == "") {
		return newValidationError(o, StrMissingFileID)
	}

	if (o.Message == nil) && (o.FileName == "") {
		return newValidationError(o, StrMissingFileName)
	}

	if o.Message != nil {
		if filesPayload, okFile := o.Message.(PNPublishFileMessage); okFile {
			if filesPayload.PNFile != nil {
				if filesPayload.PNFile.ID == "" {
					return newValidationError(o, StrMissingFileID)
				}
				if filesPayload.PNFile.Name == "" {
					return newValidationError(o, StrMissingFileName)
				}
			} else {
				return newValidationError(o, StrMissingFileID)
			}
		} else {
			return newValidationError(o, StrMissingMessage)
		}
	}

	return nil
}

func (o *publishFileMessageOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(publishFileMessagePostPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0"), nil
	}

	if o.Message == nil {
		m := &PNPublishMessage{
			Text: o.MessageText,
		}

		file := &PNFileInfoForPublish{
			ID:   o.FileID,
			Name: o.FileName,
		}

		o.Message = PNPublishFileMessage{
			PNFile:    file,
			PNMessage: m,
		}
	}

	if cipherKey := o.pubnub.Config.CipherKey; cipherKey != "" {
		var msg string
		var p *publishBuilder
		if o.context() != nil {
			p = newPublishBuilderWithContext(o.pubnub, o.context())
		} else {
			p = newPublishBuilder(o.pubnub)
		}
		p.opts.Message = o.Message

		msg, errJSONMarshal := p.opts.encryptProcessing(cipherKey)
		if errJSONMarshal != nil {
			return "", errJSONMarshal
		}

		o.pubnub.Config.Log.Println("EncryptString: encrypted", msg)
		return fmt.Sprintf(publishFileMessageGetPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0",
			utils.URLEncode(msg)), nil
	}
	var msg string
	jsonEncBytes, errEnc := json.Marshal(o.Message)
	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
		return "", errEnc
	}
	msg = string(jsonEncBytes)
	return fmt.Sprintf(publishFileMessageGetPath,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.Channel),
		"0",
		utils.URLEncode(msg),
	), nil

}

func (o *publishFileMessageOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *publishFileMessageOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *publishFileMessageOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		jsonEncBytes, errEnc := json.Marshal(o.Message)
		if errEnc != nil {
			o.pubnub.Config.Log.Printf("ERROR: PublishFileMessage error: %s\n", errEnc.Error())
			return []byte{}, errEnc
		}
		return jsonEncBytes, nil
	}
	return []byte{}, nil
}

func (o *publishFileMessageOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *publishFileMessageOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
}

func (o *publishFileMessageOpts) isAuthRequired() bool {
	return true
}

func (o *publishFileMessageOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *publishFileMessageOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *publishFileMessageOpts) operationType() OperationType {
	return PNPublishFileMessageOperation
}

func (o *publishFileMessageOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PublishFileMessageResponse is the response to PublishFileMessage request.
type PublishFileMessageResponse struct {
	Timestamp int64
}

func newPublishFileMessageResponse(jsonBytes []byte, o *publishFileMessageOpts,
	status StatusResponse) (*PublishFileMessageResponse, StatusResponse, error) {

	resp := &PublishFileMessageResponse{}

	var value []interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPublishFileMessageResponse, status, e
	}

	if len(value) > 1 {
		timeString, ok := value[2].(string)
		if !ok {
			return emptyPublishFileMessageResponse, status, pnerr.NewResponseParsingError(fmt.Sprintf("Error unmarshalling response 2, %s %v", value[2], value), nil, nil)
		}
		timestamp, err := strconv.ParseInt(timeString, 10, 64)
		if err != nil {
			return emptyPublishFileMessageResponse, status, err
		}

		return &PublishFileMessageResponse{
			Timestamp: timestamp,
		}, status, nil
	}

	return resp, status, nil
}
