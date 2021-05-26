package pubnub

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/pubnub/go/v5/pnerr"
	"github.com/pubnub/go/v5/utils"
)

var emptySendFileToS3Response *PNSendFileToS3Response

type sendFileToS3Builder struct {
	opts *sendFileToS3Opts
}

func newSendFileToS3Builder(pubnub *PubNub) *sendFileToS3Builder {
	builder := sendFileToS3Builder{
		opts: &sendFileToS3Opts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newSendFileToS3BuilderWithContext(pubnub *PubNub,
	context Context) *sendFileToS3Builder {
	builder := sendFileToS3Builder{
		opts: &sendFileToS3Opts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *sendFileToS3Builder) CipherKey(cipherKey string) *sendFileToS3Builder {
	b.opts.CipherKey = cipherKey

	return b
}

func (b *sendFileToS3Builder) FileUploadRequestData(fileUploadRequestData PNFileUploadRequest) *sendFileToS3Builder {
	b.opts.FileUploadRequestData = fileUploadRequestData

	return b
}

func (b *sendFileToS3Builder) File(f *os.File) *sendFileToS3Builder {
	b.opts.File = f

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *sendFileToS3Builder) QueryParam(queryParam map[string]string) *sendFileToS3Builder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the sendFileToS3 request.
func (b *sendFileToS3Builder) Transport(tr http.RoundTripper) *sendFileToS3Builder {
	b.opts.Transport = tr
	return b
}

// Execute runs the sendFileToS3 request.
func (b *sendFileToS3Builder) Execute() (*PNSendFileToS3Response, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySendFileToS3Response, status, err
	}

	return newPNSendFileToS3Response(rawJSON, b.opts, status)
}

type sendFileToS3Opts struct {
	pubnub *PubNub

	File                  *os.File
	FileUploadRequestData PNFileUploadRequest
	QueryParam            map[string]string
	CipherKey             string
	Transport             http.RoundTripper

	ctx Context
}

func (o *sendFileToS3Opts) config() Config {
	return *o.pubnub.Config
}

func (o *sendFileToS3Opts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *sendFileToS3Opts) context() Context {
	return o.ctx
}

func (o *sendFileToS3Opts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *sendFileToS3Opts) buildPath() (string, error) {
	return o.FileUploadRequestData.URL, nil
}

func (o *sendFileToS3Opts) buildQuery() (*url.Values, error) {
	return &url.Values{}, nil
}

func (o *sendFileToS3Opts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *sendFileToS3Opts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *sendFileToS3Opts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {

	fileInfo, _ := o.File.Stat()
	s := fileInfo.Size()
	buffer := make([]byte, 512)
	_, err := o.File.Read(buffer)
	if err != nil {
		return bytes.Buffer{}, nil, s, err
	}
	o.File.Seek(0, 0)
	contentType := http.DetectContentType(buffer)

	var fileBody bytes.Buffer
	writer := multipart.NewWriter(&fileBody)

	for _, v := range o.FileUploadRequestData.FormFields {
		o.pubnub.Config.Log.Printf("FormFields: Key: %s Value: %s\n", v.Key, v.Value)
		if v.Key == "Content-Type" {
			v.Value = contentType
		}
		_ = writer.WriteField(v.Key, v.Value)
	}

	filePart, errFilePart := writer.CreateFormFile("file", fileInfo.Name())

	if errFilePart != nil {
		o.pubnub.Config.Log.Printf("ERROR: writer CreateFormFile: %s\n", errFilePart.Error())
		return bytes.Buffer{}, writer, s, errFilePart
	}

	if o.CipherKey != "" {
		utils.EncryptFile(o.CipherKey, []byte{}, filePart, o.File)
	} else if o.pubnub.Config.CipherKey != "" {
		utils.EncryptFile(o.pubnub.Config.CipherKey, []byte{}, filePart, o.File)
	} else {
		_, errIOCopy := io.Copy(filePart, o.File)

		if errIOCopy != nil {
			o.pubnub.Config.Log.Printf("ERROR: io Copy error: %s\n", errIOCopy.Error())
			return bytes.Buffer{}, writer, s, errIOCopy
		}
	}

	errWriterClose := writer.Close()
	if errWriterClose != nil {
		o.pubnub.Config.Log.Printf("ERROR: Writer close: %s\n", errWriterClose.Error())
		return bytes.Buffer{}, writer, s, errWriterClose
	}

	return fileBody, writer, s, nil

}

func (o *sendFileToS3Opts) httpMethod() string {
	return "POSTFORM"
}

func (o *sendFileToS3Opts) isAuthRequired() bool {
	return true
}

func (o *sendFileToS3Opts) requestTimeout() int {
	return o.pubnub.Config.FileUploadRequestTimeout
}

func (o *sendFileToS3Opts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *sendFileToS3Opts) operationType() OperationType {
	return PNSendFileToS3Operation
}

func (o *sendFileToS3Opts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// PNSendFileToS3Response is the File Upload API Response for Get Spaces
type PNSendFileToS3Response struct {
}

func newPNSendFileToS3Response(jsonBytes []byte, o *sendFileToS3Opts,
	status StatusResponse) (*PNSendFileToS3Response, StatusResponse, error) {

	resp := &PNSendFileToS3Response{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySendFileToS3Response, status, e
	}
	o.pubnub.Config.Log.Printf("newPNSendFileToS3Response status.StatusCode==> %d", status.StatusCode)

	return resp, status, nil
}
