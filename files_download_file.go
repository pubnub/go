package pubnub

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/v7/utils"
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

func (b *downloadFileBuilder) Execute() (*PNDownloadFileResponse, StatusResponse, error) {
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
	b.opts.pubnub.Config.Log.Printf("u.RequestURI(): %s", u.RequestURI())
	resp, err := b.opts.client().Get(u.RequestURI())
	if err != nil {
		b.opts.pubnub.Config.Log.Printf("err %s", err)
		return nil, stat, err
	}
	if resp.StatusCode != 200 {
		stat.StatusCode = resp.StatusCode
		return nil, stat, err
	}
	contentLenEnc, err := strconv.ParseInt(string(resp.Header.Get("Content-Length")), 10, 64)
	if err != nil {
		b.opts.pubnub.Config.Log.Printf("err in parsing content length %s", err)
		return nil, stat, err
	}

	var respDL *PNDownloadFileResponse
	if b.opts.CipherKey != "" {
		r, w := io.Pipe()
		utils.DecryptFile(b.opts.CipherKey, contentLenEnc, resp.Body, w)
		respDL = &PNDownloadFileResponse{
			File: r,
		}

	} else if b.opts.pubnub.Config.CipherKey != "" {
		r, w := io.Pipe()
		utils.DecryptFile(b.opts.pubnub.Config.CipherKey, contentLenEnc, resp.Body, w)
		respDL = &PNDownloadFileResponse{
			File: r,
		}

	} else {
		respDL = &PNDownloadFileResponse{
			File: resp.Body,
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

func (o *downloadFileOpts) operationType() OperationType {
	return PNDownloadFileOperation
}

// PNDownloadFileResponse is the File Upload API Response for Get Spaces
type PNDownloadFileResponse struct {
	status int       `json:"status"`
	File   io.Reader `json:"data"`
}

func newPNDownloadFileResponse(jsonBytes []byte, o *downloadFileOpts,
	status StatusResponse) (*PNDownloadFileResponse, StatusResponse, error) {

	resp := &PNDownloadFileResponse{}

	return resp, status, nil
}
