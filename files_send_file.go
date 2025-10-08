package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pubnub/go/v7/pnerr"
)

var emptySendFileResponse *PNSendFileResponse

const sendFilePath = "/v1/files/%s/channels/%s/generate-upload-url"

type sendFileBuilder struct {
	opts *sendFileOpts
}

func newSendFileBuilder(pubnub *PubNub) *sendFileBuilder {
	return newSendFileBuilderWithContext(pubnub, pubnub.ctx)
}

func newSendFileOpts(pubnub *PubNub, ctx Context) *sendFileOpts {
	return &sendFileOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newSendFileBuilderWithContext(pubnub *PubNub,
	context Context) *sendFileBuilder {
	builder := sendFileBuilder{
		opts: newSendFileOpts(pubnub, context)}
	return &builder
}

// TTL sets the TTL (hours) for the Publish request.
func (b *sendFileBuilder) TTL(ttl int) *sendFileBuilder {
	b.opts.TTL = ttl

	return b
}

// Meta sets the Meta Payload for the Publish request.
func (b *sendFileBuilder) Meta(meta interface{}) *sendFileBuilder {
	b.opts.Meta = meta

	return b
}

// ShouldStore if true the messages are stored in History
func (b *sendFileBuilder) ShouldStore(store bool) *sendFileBuilder {
	b.opts.ShouldStore = store
	return b
}

func (b *sendFileBuilder) CipherKey(cipher string) *sendFileBuilder {
	b.opts.CipherKey = cipher

	return b
}

func (b *sendFileBuilder) Channel(channel string) *sendFileBuilder {
	b.opts.Channel = channel

	return b
}

func (b *sendFileBuilder) Name(name string) *sendFileBuilder {
	b.opts.Name = name

	return b
}

// Message sets the message content for the file upload.
// Accepts any JSON-serializable type: string, map[string]interface{}, []interface{}, number, bool, etc.
// Examples:
//   - String: .Message("Hello")
//   - JSON object: .Message(map[string]interface{}{"type": "document", "priority": "high"})
//   - Array: .Message([]string{"item1", "item2"})
func (b *sendFileBuilder) Message(message interface{}) *sendFileBuilder {
	b.opts.Message = message

	return b
}

func (b *sendFileBuilder) File(f *os.File) *sendFileBuilder {
	b.opts.File = f

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *sendFileBuilder) QueryParam(queryParam map[string]string) *sendFileBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the sendFile request.
func (b *sendFileBuilder) Transport(tr http.RoundTripper) *sendFileBuilder {
	b.opts.Transport = tr

	return b
}

// CustomMessageType sets the User-specified message type string - limited by 3-50 case-sensitive alphanumeric characters
// with only `-` and `_` special characters allowed.
func (b *sendFileBuilder) CustomMessageType(messageType string) *sendFileBuilder {
	b.opts.CustomMessageType = messageType

	return b
}

// UseRawMessage sets whether the message should be sent as raw content instead of being wrapped in a JSON "text" field.
// When true, the message will be sent directly without the {"text": ...} wrapper.
// When false (default), the message is wrapped in a "text" field for backward compatibility.
// Works with any JSON-serializable type: strings, objects, arrays, numbers, booleans, etc.
// Examples:
//
//	UseRawMessage(false): {"message": {"text": "Hello"}, "file": {"id": "123", "name": "file.txt"}}
//	UseRawMessage(true):  {"message": "Hello", "file": {"id": "123", "name": "file.txt"}}
//	UseRawMessage(true):  {"message": {"type": "doc", "priority": "high"}, "file": {...}}
func (b *sendFileBuilder) UseRawMessage(useRawMessage bool) *sendFileBuilder {
	b.opts.UseRawMessage = useRawMessage

	return b
}

// Execute runs the sendFile request.
func (b *sendFileBuilder) Execute() (*PNSendFileResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySendFileResponse, status, err
	}

	return newPNSendFileResponse(rawJSON, b.opts, status)
}

type sendFileOpts struct {
	endpointOpts

	Channel           string
	Name              string
	Message           interface{}
	File              *os.File
	CipherKey         string
	TTL               int
	Meta              interface{}
	ShouldStore       bool
	QueryParam        map[string]string
	CustomMessageType string
	UseRawMessage     bool

	Transport http.RoundTripper
}

func (o *sendFileOpts) isCustomMessageTypeCorrect() bool {
	return isCustomMessageTypeValid(o.CustomMessageType)
}

func (o *sendFileOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	if o.Name == "" {
		return newValidationError(o, StrMissingFileName)
	}

	if o.File == nil {
		return newValidationError(o, "file is required")
	}

	if !o.isCustomMessageTypeCorrect() {
		return newValidationError(o, StrInvalidCustomMessageType)
	}

	return nil
}

func (o *sendFileOpts) buildPath() (string, error) {
	return fmt.Sprintf(sendFilePath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *sendFileOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	if len(o.CustomMessageType) > 0 {
		q.Set("custom_message_type", o.CustomMessageType)
	}

	return q, nil
}

// PNSendFileBody is used to create the body of the request
type PNSendFileBody struct {
	Name string `json:"name"`
}

func (o *sendFileOpts) buildBody() ([]byte, error) {
	b := &PNSendFileBody{
		Name: o.Name,
	}
	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *sendFileOpts) httpMethod() string {
	return "POST"
}

func (o *sendFileOpts) isAuthRequired() bool {
	return true
}

func (o *sendFileOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *sendFileOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *sendFileOpts) operationType() OperationType {
	return PNSendFileOperation
}

// PNSendFileResponseForS3 is the File Upload API Response for SendFile.
type PNSendFileResponseForS3 struct {
	Status            int                 `json:"status"`
	Data              PNFileData          `json:"data"`
	FileUploadRequest PNFileUploadRequest `json:"file_upload_request"`
}

// PNSendFileResponse is the type used to store the response info of Send File.
type PNSendFileResponse struct {
	Timestamp int64
	Status    int        `json:"status"`
	Data      PNFileData `json:"data"`
}

// TODO Add retry on publish failure
func newPNSendFileResponse(jsonBytes []byte, o *sendFileOpts,
	status StatusResponse) (*PNSendFileResponse, StatusResponse, error) {

	respForS3 := &PNSendFileResponseForS3{}

	err := json.Unmarshal(jsonBytes, &respForS3)
	if err != nil {
		e := pnerr.NewResponseParsingError("error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)
		return emptySendFileResponse, status, e
	}
	var s *sendFileToS3Builder
	if o.context() != nil {
		s = newSendFileToS3BuilderWithContext(o.pubnub, o.context())
	} else {
		s = newSendFileToS3Builder(o.pubnub)
	}
	_, s3ResponseStatus, errS3Response := s.File(o.File).CipherKey(o.CipherKey).FileUploadRequestData(respForS3.FileUploadRequest).Execute()
	if s3ResponseStatus.StatusCode != 204 {
		o.pubnub.Config.Log.Printf("s3ResponseStatus: %d", s3ResponseStatus.StatusCode)
		return emptySendFileResponse, s3ResponseStatus, errS3Response
	}

	file := &PNFileInfoForPublish{
		ID:   respForS3.Data.ID,
		Name: o.Name,
	}

	m := &PNPublishMessage{
		Text: o.Message,
	}
	message := PNPublishFileMessage{
		PNFile:    file,
		PNMessage: m,
	}

	sent := false
	tryCount := 0
	var timestamp int64
	maxCount := o.config().FileMessagePublishRetryLimit
	for !sent && tryCount < maxCount {
		tryCount++
		pubFileMessageResponse, pubFileResponseStatus, errPubFileResponse := o.pubnub.PublishFileMessage().TTL(o.TTL).Meta(o.Meta).ShouldStore(o.ShouldStore).Channel(o.Channel).Message(message).UseRawMessage(o.UseRawMessage).Execute()
		if errPubFileResponse != nil {
			if tryCount >= maxCount {
				pubFileResponseStatus.AdditionalData = file
				return emptySendFileResponse, pubFileResponseStatus, errPubFileResponse
			}
			continue
		} else {
			timestamp = pubFileMessageResponse.Timestamp
			sent = true
			break
		}
	}
	resp := &PNSendFileResponse{}
	d := PNFileData{}
	d.ID = respForS3.Data.ID
	resp.Data = d
	resp.Timestamp = timestamp

	return resp, status, nil
}
