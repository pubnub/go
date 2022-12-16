package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
)

var emptyGetFileURLResponse *PNGetFileURLResponse

const getFileURLPath = "/v1/files/%s/channels/%s/files/%s/%s"

type getFileURLBuilder struct {
	opts *getFileURLOpts
}

func newGetFileURLBuilder(pubnub *PubNub) *getFileURLBuilder {
	return newGetFileURLBuilderWithContext(pubnub, pubnub.ctx)
}

func newGetFileURLBuilderWithContext(pubnub *PubNub, context Context) *getFileURLBuilder {
	builder := getFileURLBuilder{
		opts: newGetFileURLOpts(
			pubnub,
			context,
		),
	}

	return &builder
}

func (b *getFileURLBuilder) Channel(channel string) *getFileURLBuilder {
	b.opts.Channel = channel

	return b
}

func (b *getFileURLBuilder) ID(id string) *getFileURLBuilder {
	b.opts.ID = id

	return b
}

func (b *getFileURLBuilder) Name(name string) *getFileURLBuilder {
	b.opts.Name = name

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getFileURLBuilder) QueryParam(queryParam map[string]string) *getFileURLBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getFileURL request.
func (b *getFileURLBuilder) Transport(tr http.RoundTripper) *getFileURLBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getFileURL request.
func (b *getFileURLBuilder) Execute() (*PNGetFileURLResponse, StatusResponse, error) {
	u, _ := buildURL(b.opts)

	resp := &PNGetFileURLResponse{
		URL: u.RequestURI(),
	}
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
	return resp, stat, nil
}

func newGetFileURLOpts(pubnub *PubNub, ctx Context) *getFileURLOpts {
	return &getFileURLOpts{endpointOpts: endpointOpts{
		pubnub: pubnub,
		ctx:    ctx,
	}}
}

type getFileURLOpts struct {
	endpointOpts

	Channel    string
	ID         string
	Name       string
	QueryParam map[string]string

	Transport http.RoundTripper
}

func (o *getFileURLOpts) validate() error {
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

func (o *getFileURLOpts) buildPath() (string, error) {
	return fmt.Sprintf(getFileURLPath,
		o.pubnub.Config.SubscribeKey, o.Channel, o.ID, o.Name), nil
}

func (o *getFileURLOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getFileURLOpts) operationType() OperationType {
	return PNGetFileURLOperation
}

// PNGetFileURLResponse is the File Upload API Response for Get Spaces
type PNGetFileURLResponse struct {
	URL string `json:"location"`
}
