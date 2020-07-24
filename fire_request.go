package pubnub

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"

	"net/http"
	"net/url"
)

type fireOpts struct {
	pubnub *PubNub

	TTL            int
	Channel        string
	Message        interface{}
	Meta           interface{}
	UsePost        bool
	Serialize      bool
	ShouldStore    bool
	DoNotReplicate bool
	Transport      http.RoundTripper
	ctx            Context
	QueryParam     map[string]string
	// nil hacks
	setTTL         bool
	setShouldStore bool
}

type fireBuilder struct {
	opts *fireOpts
}

func newFireBuilder(pubnub *PubNub) *fireBuilder {
	builder := fireBuilder{
		opts: &fireOpts{
			pubnub:    pubnub,
			Serialize: true,
		},
	}

	return &builder
}

func newFireBuilderWithContext(pubnub *PubNub, context Context) *fireBuilder {
	builder := fireBuilder{
		opts: &fireOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// TTL sets the TTL (hours) for the Fire request.
func (b *fireBuilder) TTL(ttl int) *fireBuilder {
	b.opts.TTL = ttl

	return b
}

// Channel sets the Channel for the Fire request.
func (b *fireBuilder) Channel(ch string) *fireBuilder {
	b.opts.Channel = ch

	return b
}

// Message sets the Payload for the Fire request.
func (b *fireBuilder) Message(msg interface{}) *fireBuilder {
	b.opts.Message = msg

	return b
}

// Meta sets the Meta Payload for the Fire request.
func (b *fireBuilder) Meta(meta interface{}) *fireBuilder {
	b.opts.Meta = meta

	return b
}

// UsePost sends the Fire request using HTTP POST.
func (b *fireBuilder) UsePost(post bool) *fireBuilder {
	b.opts.UsePost = post

	return b
}

// Serialize when true (default) serializes the payload before publish.
// Set to false if pre serialized payload is being used.
func (b *fireBuilder) Serialize(serialize bool) *fireBuilder {
	b.opts.Serialize = serialize

	return b
}

// Transport sets the Transport for the Fire request.
func (b *fireBuilder) Transport(tr http.RoundTripper) *fireBuilder {
	b.opts.Transport = tr

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *fireBuilder) QueryParam(queryParam map[string]string) *fireBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Fire request.
func (b *fireBuilder) Execute() (*PublishResponse, StatusResponse, error) {
	b.opts.ShouldStore = false
	b.opts.DoNotReplicate = true
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishResponse, status, err
	}

	return newPublishResponse(rawJSON, status)
}

func (o *fireOpts) config() Config {
	return *o.pubnub.Config
}

func (o *fireOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *fireOpts) context() Context {
	return o.ctx
}

func (o *fireOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	if o.Message == nil {
		return newValidationError(o, StrMissingMessage)
	}

	return nil
}

func (o *fireOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(publishPostPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0"), nil
	}

	var message []byte
	var err error

	if cipherKey := o.pubnub.Config.CipherKey; cipherKey != "" {
		msg := utils.EncryptString(cipherKey, string(message), o.pubnub.Config.UseRandomInitializationVector)

		o.Message = []byte(msg)
	}

	message, err = utils.ValueAsString(o.Message)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(publishGetPath,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.Channel),
		"0",
		utils.URLEncode(string(message))), nil
}

func (o *fireOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Meta != nil {
		meta, err := utils.ValueAsString(o.Meta)
		if err != nil {
			return &url.Values{}, err
		}

		q.Set("meta", string(meta))
	}

	q.Set("store", "0")
	q.Set("norep", "true")

	if o.setTTL {
		if o.TTL > 0 {
			q.Set("ttl", strconv.Itoa(o.TTL))
		}
	}

	q.Set("seqn", strconv.Itoa(o.pubnub.getPublishSequence()))
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *fireOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *fireOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		var msg []byte

		if o.Serialize {
			m, err := utils.ValueAsString(o.Message)
			if err != nil {
				return []byte{}, err
			}
			msg = []byte(m)
		} else {
			if s, ok := o.Message.(string); ok {
				msg = []byte(s)
			} else {
				err := pnerr.NewBuildRequestError("Type error, only string is expected")
				return []byte{}, err
			}
		}

		if cipherKey := o.pubnub.Config.CipherKey; cipherKey != "" {
			enc := utils.EncryptString(cipherKey, string(msg), o.pubnub.Config.UseRandomInitializationVector)
			msg, err := utils.ValueAsString(enc)
			if err != nil {
				return []byte{}, err
			}
			return []byte(msg), nil
		}
		return msg, nil
	}
	return []byte{}, nil
}

func (o *fireOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *fireOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
}

func (o *fireOpts) isAuthRequired() bool {
	return true
}

func (o *fireOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *fireOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *fireOpts) operationType() OperationType {
	return PNFireOperation
}

func (o *fireOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
