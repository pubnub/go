package pubnub

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pubnub/go/v8/crypto"
)

var emptyDownloadFileResponse *PNDownloadFileResponse

const downloadFilePath = "/v1/files/%s/channels/%s/files/%s/%s"

const downloadFileLimit = 100

type downloadFileBuilder struct {
	opts *downloadFileOpts
}

func newDownloadFileBuilder(pubnub *PubNub) *downloadFileBuilder {
	return newDownloadFileBuilderWithContext(pubnub, pubnub.ctx)
}

func newDownloadFileOpts(pubnub *PubNub, ctx Context) *downloadFileOpts {
	return &downloadFileOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newDownloadFileBuilderWithContext(pubnub *PubNub,
	context Context) *downloadFileBuilder {
	builder := downloadFileBuilder{
		opts: newDownloadFileOpts(pubnub, context)}
	return &builder
}

func (b *downloadFileBuilder) Channel(channel string) *downloadFileBuilder {
	b.opts.Channel = channel

	return b
}

func (b *downloadFileBuilder) CipherKey(cipherKey string) *downloadFileBuilder {
	b.opts.CipherKey = cipherKey

	return b
}

func (b *downloadFileBuilder) ID(id string) *downloadFileBuilder {
	b.opts.ID = id

	return b
}

func (b *downloadFileBuilder) Name(name string) *downloadFileBuilder {
	b.opts.Name = name

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *downloadFileBuilder) QueryParam(queryParam map[string]string) *downloadFileBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the downloadFile request.
func (b *downloadFileBuilder) Transport(tr http.RoundTripper) *downloadFileBuilder {
	b.opts.Transport = tr
	return b
}

// GetLogParams returns the user-provided parameters for logging
func (o *downloadFileOpts) GetLogParams() map[string]interface{} {
	return map[string]interface{}{
		"Channel": o.Channel,
		"ID":      o.ID,
		"Name":    o.Name,
	}
}

func (b *downloadFileBuilder) Execute() (*PNDownloadFileResponse, StatusResponse, error) {
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNDownloadFileOperation, b.opts.GetLogParams(), true)

	u, _ := buildURL(b.opts)
	stat := StatusResponse{
		AffectedChannels: []string{b.opts.Channel},
		AuthKey:          b.opts.config().AuthKey,
		Category:         PNUnknownCategory,
		Operation:        PNGetFileURLOperation,
		StatusCode:       200,
		TLSEnabled:       b.opts.config().Secure,
		Origin:           b.opts.config().Origin,
		UUID:             b.opts.config().UUID,
	}
	b.opts.pubnub.loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Downloading file: URI=%s", u.RequestURI()), false)
	resp, err := b.opts.client().Get(u.RequestURI())
	if err != nil {
		b.opts.pubnub.loggerManager.LogError(err, "FileDownloadRequestFailed", PNDownloadFileOperation, true)
		return nil, stat, err
	}
	if resp.StatusCode != 200 {
		stat.StatusCode = resp.StatusCode
		return nil, stat, err
	}

	var respDL *PNDownloadFileResponse
	if b.opts.CipherKey == "" && b.opts.pubnub.getCryptoModule() == nil {
		respDL = &PNDownloadFileResponse{
			File: resp.Body,
		}
	} else {
		var e error
		cryptoModule := b.opts.pubnub.getCryptoModule()
		if b.opts.CipherKey != "" {
			cryptoModule, e = crypto.NewLegacyCryptoModule(b.opts.CipherKey, true)
			if e != nil {
				b.opts.pubnub.loggerManager.LogError(e, "FileDownloadCryptoModuleInitFailed", PNDownloadFileOperation, true)
				return nil, stat, e
			}
			b.opts.pubnub.loggerManager.LogSimple(PNLogLevelDebug, "Crypto Module initialized for file download: type=LegacyCryptoModule, randomIV=true", false)
		}

		r, e := cryptoModule.DecryptStream(resp.Body)
		if e != nil {
			return nil, stat, e
		}
		respDL = &PNDownloadFileResponse{
			File: r,
		}

	}

	return respDL, stat, nil
}

type downloadFileOpts struct {
	endpointOpts
	Channel    string
	CipherKey  string
	ID         string
	Name       string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *downloadFileOpts) validate() error {
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

func (o *downloadFileOpts) buildPath() (string, error) {
	return fmt.Sprintf(downloadFilePath,
		o.pubnub.Config.SubscribeKey, o.Channel, o.ID, o.Name), nil
}

func (o *downloadFileOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *downloadFileOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *downloadFileOpts) httpMethod() string {
	return "GET"
}

func (o *downloadFileOpts) isAuthRequired() bool {
	return true
}

func (o *downloadFileOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *downloadFileOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *downloadFileOpts) operationType() OperationType {
	return PNDownloadFileOperation
}

// PNDownloadFileResponse is the File Upload API Response for Get Spaces
type PNDownloadFileResponse struct {
	Status int       `json:"status"`
	File   io.Reader `json:"data"`
}

func newPNDownloadFileResponse(jsonBytes []byte, o *downloadFileOpts,
	status StatusResponse) (*PNDownloadFileResponse, StatusResponse, error) {

	resp := &PNDownloadFileResponse{}

	return resp, status, nil
}
