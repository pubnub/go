package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"

	"net/http"
	"net/url"
)

type fireOpts struct {
	pubnub *PubNub

	Ttl            int
	Channel        string
	Message        interface{}
	Meta           interface{}
	UsePost        bool
	Serialize      bool
	ShouldStore    bool
	DoNotReplicate bool
	Transport      http.RoundTripper
	ctx            Context
	// nil hacks
	setTtl         bool
	setShouldStore bool
}

type fireBuilder struct {
	opts *fireOpts
}

func newFireResponse(jsonBytes []byte, status StatusResponse) (
	*PublishResponse, StatusResponse, error) {
	var value []interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPublishResponse, status, e
	}

	if timeString, ok := value[2].(string); !ok {
		return emptyPublishResponse, status, pnerr.NewResponseParsingError(fmt.Sprintf("Error unmarshalling response, %s %v", value[2], value), nil, nil)
	} else {
		timestamp, err := strconv.Atoi(timeString)
		if err != nil {
			return emptyPublishResponse, status, err
		}

		return &PublishResponse{
			Timestamp: timestamp,
		}, status, nil
	}
}

func newFireBuilder(pubnub *PubNub) *fireBuilder {
	builder := fireBuilder{
		opts: &fireOpts{
			pubnub: pubnub,
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

func (b *fireBuilder) Ttl(ttl int) *fireBuilder {
	b.opts.Ttl = ttl

	return b
}

func (b *fireBuilder) Channel(ch string) *fireBuilder {
	b.opts.Channel = ch

	return b
}

func (b *fireBuilder) Message(msg interface{}) *fireBuilder {
	b.opts.Message = msg

	return b
}

func (b *fireBuilder) Meta(meta interface{}) *fireBuilder {
	b.opts.Meta = meta

	return b
}

func (b *fireBuilder) UsePost(post bool) *fireBuilder {
	b.opts.UsePost = post

	return b
}

func (b *fireBuilder) Serialize(serialize bool) *fireBuilder {
	b.opts.Serialize = serialize

	return b
}

func (b *fireBuilder) Transport(tr http.RoundTripper) *fireBuilder {
	b.opts.Transport = tr

	return b
}

func (b *fireBuilder) Execute() (*PublishResponse, StatusResponse, error) {
	b.opts.ShouldStore = false
	b.opts.DoNotReplicate = true
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishResponse, status, err
	}

	return newPublishResponse(rawJson, status)
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
			utils.UrlEncode(o.Channel),
			"0"), nil
	}

	var message []byte
	var err error

	if cipherKey := o.pubnub.Config.CipherKey; cipherKey != "" {
		msg := utils.EncryptString(cipherKey, string(message))

		o.Message = []byte(msg)
	}

	message, err = utils.ValueAsString(o.Message)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(publishGetPath,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.Channel),
		"0",
		utils.UrlEncode(string(message))), nil
}

func (o *fireOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	if o.Meta != nil {
		meta, err := utils.ValueAsString(o.Meta)
		if err != nil {
			return &url.Values{}, err
		}

		q.Set("meta", string(meta))
	}

	q.Set("store", "0")
	q.Set("norep", "true")

	if o.setTtl {
		if o.Ttl > 0 {
			q.Set("ttl", strconv.Itoa(o.Ttl))
		}
	}

	q.Set("seqn", strconv.Itoa(o.pubnub.getPublishSequence()))

	return q, nil
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
			enc := utils.EncryptString(cipherKey, string(msg))
			msg, err := utils.ValueAsString(enc)
			if err != nil {
				return []byte{}, err
			}
			return []byte(msg), nil
		} else {
			return msg, nil
		}
	} else {
		return []byte{}, nil
	}
}

func (o *fireOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	} else {
		return "GET"
	}
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
