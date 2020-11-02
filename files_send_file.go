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
	"os"

	"github.com/pubnub/go/pnerr"
)

var emptySendFileResponse *PNSendFileResponse

const sendFilePath = "/v1/files/%s/channels/%s/generate-upload-url"

type sendFileBuilder struct {
	opts *sendFileOpts
}

func newSendFileBuilder(pubnub *PubNub) *sendFileBuilder {
	builder := sendFileBuilder{
		opts: &sendFileOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newSendFileBuilderWithContext(pubnub *PubNub,
	context Context) *sendFileBuilder {
	builder := sendFileBuilder{
		opts: &sendFileOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

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

func (b *sendFileBuilder) Message(message string) *sendFileBuilder {
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

// Execute runs the sendFile request.
func (b *sendFileBuilder) Execute() (*PNSendFileResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySendFileResponse, status, err
	}

	return newPNSendFileResponse(rawJSON, b.opts, status)
}

type sendFileOpts struct {
	pubnub *PubNub

	Channel     string
	Name        string
	Message     string
	File        *os.File
	CipherKey   string
	TTL         int
	Meta        interface{}
	ShouldStore bool
	QueryParam  map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *sendFileOpts) config() Config {
	return *o.pubnub.Config
}

func (o *sendFileOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *sendFileOpts) context() Context {
	return o.ctx
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
	return nil
}

func (o *sendFileOpts) buildPath() (string, error) {
	return fmt.Sprintf(sendFilePath,
		o.pubnub.Config.SubscribeKey, o.Channel), nil
}

func (o *sendFileOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *sendFileOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
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

func (o *sendFileOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
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

func (o *sendFileOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNSendFileResponseForS3 is the File Upload API Response for SendFile.
type PNSendFileResponseForS3 struct {
	status            int                 `json:"status"`
	Data              PNFileData          `json:"data"`
	FileUploadRequest PNFileUploadRequest `json:"file_upload_request"`
}

// PNSendFileResponse is the type used to store the response info of Send File.
type PNSendFileResponse struct {
	Timestamp int64
	status    int        `json:"status"`
	Data      PNFileData `json:"data"`
}

// TODO Add retry on publish failure
func newPNSendFileResponse(jsonBytes []byte, o *sendFileOpts,
	status StatusResponse) (*PNSendFileResponse, StatusResponse, error) {

	respForS3 := &PNSendFileResponseForS3{}

	err := json.Unmarshal(jsonBytes, &respForS3)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)
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

	m := &PNPublishMessage{
		Text: o.Message,
	}

	file := &PNFileInfoForPublish{
		ID:   respForS3.Data.ID,
		Name: o.Name,
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
		pubFileMessageResponse, pubFileResponseStatus, errPubFileResponse := o.pubnub.PublishFileMessage().TTL(o.TTL).Meta(o.Meta).ShouldStore(o.ShouldStore).Channel(o.Channel).Message(message).Execute()
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
