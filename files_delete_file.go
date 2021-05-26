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

	"github.com/pubnub/go/v5/pnerr"
)

var emptyDeleteFileResponse *PNDeleteFileResponse

const deleteFilePath = "/v1/files/%s/channels/%s/files/%s/%s"

type deleteFileBuilder struct {
	opts *deleteFileOpts
}

func newDeleteFileBuilder(pubnub *PubNub) *deleteFileBuilder {
	builder := deleteFileBuilder{
		opts: &deleteFileOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newDeleteFileBuilderWithContext(pubnub *PubNub,
	context Context) *deleteFileBuilder {
	builder := deleteFileBuilder{
		opts: &deleteFileOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *deleteFileBuilder) Channel(channel string) *deleteFileBuilder {
	b.opts.Channel = channel

	return b
}

func (b *deleteFileBuilder) ID(id string) *deleteFileBuilder {
	b.opts.ID = id

	return b
}

func (b *deleteFileBuilder) Name(name string) *deleteFileBuilder {
	b.opts.Name = name

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteFileBuilder) QueryParam(queryParam map[string]string) *deleteFileBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the deleteFile request.
func (b *deleteFileBuilder) Transport(tr http.RoundTripper) *deleteFileBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the deleteFile request.
func (b *deleteFileBuilder) Execute() (*PNDeleteFileResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyDeleteFileResponse, status, err
	}

	return newPNDeleteFileResponse(rawJSON, b.opts, status)
}

type deleteFileOpts struct {
	pubnub *PubNub

	Channel    string
	ID         string
	Name       string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *deleteFileOpts) config() Config {
	return *o.pubnub.Config
}

func (o *deleteFileOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *deleteFileOpts) context() Context {
	return o.ctx
}

func (o *deleteFileOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	if o.Name == "" {
		return newValidationError(o, StrMissingFileName)
	}

	if o.ID == "" {
		return newValidationError(o, StrMissingFileID)
	}

	return nil
}

func (o *deleteFileOpts) buildPath() (string, error) {
	return fmt.Sprintf(deleteFilePath,
		o.pubnub.Config.SubscribeKey, o.Channel, o.ID, o.Name), nil
}

func (o *deleteFileOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *deleteFileOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *deleteFileOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *deleteFileOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *deleteFileOpts) httpMethod() string {
	return "DELETE"
}

func (o *deleteFileOpts) isAuthRequired() bool {
	return true
}

func (o *deleteFileOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *deleteFileOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *deleteFileOpts) operationType() OperationType {
	return PNDeleteFileOperation
}

func (o *deleteFileOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNDeleteFileResponse is the File Upload API Response for Delete file operation
type PNDeleteFileResponse struct {
	status int `json:"status"`
}

func newPNDeleteFileResponse(jsonBytes []byte, o *deleteFileOpts,
	status StatusResponse) (*PNDeleteFileResponse, StatusResponse, error) {

	resp := &PNDeleteFileResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyDeleteFileResponse, status, e
	}

	return resp, status, nil
}
