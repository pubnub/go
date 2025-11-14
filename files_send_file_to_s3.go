package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/pubnub/go/v8/crypto"

	"github.com/pubnub/go/v8/pnerr"
)

var emptySendFileToS3Response *PNSendFileToS3Response

type sendFileToS3Builder struct {
	opts *sendFileToS3Opts
}

func newSendFileToS3Builder(pubnub *PubNub) *sendFileToS3Builder {
	return newSendFileToS3BuilderWithContext(pubnub, pubnub.ctx)
}

func newSendFileToS3Opts(pubnub *PubNub, ctx Context) *sendFileToS3Opts {
	return &sendFileToS3Opts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newSendFileToS3BuilderWithContext(pubnub *PubNub,
	context Context) *sendFileToS3Builder {
	builder := sendFileToS3Builder{
		opts: newSendFileToS3Opts(pubnub, context)}
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
	endpointOpts

	File                  *os.File
	FileUploadRequestData PNFileUploadRequest
	QueryParam            map[string]string
	CipherKey             string
	Transport             http.RoundTripper
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

func (o *sendFileToS3Opts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	if o.File == nil {
		return bytes.Buffer{}, nil, 0, fmt.Errorf("file is nil")
	}

	fileInfo, err := o.File.Stat()
	if err != nil {
		return bytes.Buffer{}, nil, 0, fmt.Errorf("failed to get file info: %v", err)
	}
	s := fileInfo.Size()
	buffer := make([]byte, 512)
	_, err = o.File.Read(buffer)
	if err != nil {
		return bytes.Buffer{}, nil, s, err
	}
	o.File.Seek(0, 0)
	contentType := http.DetectContentType(buffer)

	var fileBody bytes.Buffer
	writer := multipart.NewWriter(&fileBody)

	for _, v := range o.FileUploadRequestData.FormFields {
		o.pubnub.loggerManager.LogSimple(PNLogLevelTrace, fmt.Sprintf("File upload: form field %s=%s", v.Key, v.Value), false)
		if v.Key == "Content-Type" {
			v.Value = contentType
		}
		_ = writer.WriteField(v.Key, v.Value)
	}

	filePart, errFilePart := writer.CreateFormFile("file", fileInfo.Name())

	if errFilePart != nil {
		o.pubnub.loggerManager.LogError(errFilePart, "FileUploadCreateFormFileFailed", PNSendFileToS3Operation, true)
		return bytes.Buffer{}, writer, s, errFilePart
	}

	if o.CipherKey == "" && o.pubnub.getCryptoModule() == nil {
		_, errIOCopy := io.Copy(filePart, o.File)

		if errIOCopy != nil {
			o.pubnub.loggerManager.LogError(errIOCopy, "FileUploadCopyFailed", PNSendFileToS3Operation, true)
			return bytes.Buffer{}, writer, s, errIOCopy
		}
	} else {
		var e error
		cryptoModule := o.pubnub.getCryptoModule()

		if o.CipherKey != "" {
			cryptoModule, e = crypto.NewLegacyCryptoModule(o.CipherKey, true)
			if e != nil {
				o.pubnub.loggerManager.LogError(e, "FileUploadCryptoModuleInitFailed", PNSendFileToS3Operation, true)
				return bytes.Buffer{}, writer, s, e
			}
			o.pubnub.loggerManager.LogSimple(PNLogLevelDebug, `Crypto Module initialized for file upload:
				type: LegacyCryptoModule
				cipherKey: ***
				randomIV: true`, false)
		}

		e = encryptStreamAndCopyTo(cryptoModule, o.File, filePart, o.pubnub.loggerManager)
		if e != nil {
			o.pubnub.loggerManager.LogError(e, "FileUploadEncryptionFailed", PNSendFileToS3Operation, true)
			return bytes.Buffer{}, writer, s, e
		}
	}

	errWriterClose := writer.Close()
	if errWriterClose != nil {
		o.pubnub.loggerManager.LogError(errWriterClose, "FileUploadWriterCloseFailed", PNSendFileToS3Operation, true)
		return bytes.Buffer{}, writer, s, errWriterClose
	}

	return fileBody, writer, s, nil

}

func (o *sendFileToS3Opts) httpMethod() string {
	return "POSTFORM"
}

func (o *sendFileToS3Opts) operationType() OperationType {
	return PNSendFileToS3Operation
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
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySendFileToS3Response, status, e
	}
	o.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("File upload to S3 completed: status=%d", status.StatusCode), false)

	return resp, status, nil
}
