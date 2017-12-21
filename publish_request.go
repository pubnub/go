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

const PUBLISH_GET_PATH = "/publish/%s/%s/0/%s/%s/%s"
const PUBLISH_POST_PATH = "/publish/%s/%s/0/%s/%s"

var emptyPublishResponse *PublishResponse

type publishOpts struct {
	pubnub *PubNub

	Ttl     int
	Channel string
	Message interface{}
	Meta    interface{}

	UsePost        bool
	ShouldStore    bool
	Serialize      bool
	DoNotReplicate bool

	Transport http.RoundTripper

	ctx Context

	// nil hacks
	SetTtl         bool
	SetShouldStore bool
}

type PublishResponse struct {
	Timestamp int
}

type publishBuilder struct {
	opts *publishOpts
}

func newPublishResponse(jsonBytes []byte, status StatusResponse) (
	*PublishResponse, StatusResponse, error) {
	var value []interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPublishResponse, status, e
	}

	timeString := value[2].(string)
	timestamp, err := strconv.Atoi(timeString)
	if err != nil {
		return emptyPublishResponse, status, err
	}

	return &PublishResponse{
		Timestamp: timestamp,
	}, status, nil
}

func newPublishBuilder(pubnub *PubNub) *publishBuilder {
	builder := publishBuilder{
		opts: &publishOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newPublishBuilderWithContext(pubnub *PubNub, context Context) *publishBuilder {
	builder := publishBuilder{
		opts: &publishOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *publishBuilder) Ttl(ttl int) *publishBuilder {
	b.opts.Ttl = ttl

	return b
}

func (b *publishBuilder) Channel(ch string) *publishBuilder {
	b.opts.Channel = ch

	return b
}

func (b *publishBuilder) Message(msg interface{}) *publishBuilder {
	b.opts.Message = msg

	return b
}

func (b *publishBuilder) Meta(meta interface{}) *publishBuilder {
	b.opts.Meta = meta

	return b
}

func (b *publishBuilder) UsePost(post bool) *publishBuilder {
	b.opts.UsePost = post

	return b
}

func (b *publishBuilder) ShouldStore(store bool) *publishBuilder {
	b.opts.ShouldStore = store
	b.opts.SetShouldStore = true

	return b
}

func (b *publishBuilder) Serialize(serialize bool) *publishBuilder {
	b.opts.Serialize = serialize

	return b
}

func (b *publishBuilder) DoNotReplicate(repl bool) *publishBuilder {
	b.opts.DoNotReplicate = repl

	return b
}

func (b *publishBuilder) Transport(tr http.RoundTripper) *publishBuilder {
	b.opts.Transport = tr

	return b
}

func (b *publishBuilder) Execute() (*PublishResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishResponse, status, err
	}

	return newPublishResponse(rawJson, status)
}

func (o *publishOpts) config() Config {
	return *o.pubnub.Config
}

func (o *publishOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *publishOpts) context() Context {
	return o.ctx
}

func (o *publishOpts) validate() error {
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

func (o *publishOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(PUBLISH_POST_PATH,
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

	return fmt.Sprintf(PUBLISH_GET_PATH,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.Channel),
		"0",
		utils.UrlEncode(string(message))), nil
}

func (o *publishOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	if o.Meta != nil {
		meta, err := utils.ValueAsString(o.Meta)
		if err != nil {
			return &url.Values{}, err
		}

		q.Set("meta", string(meta))
	}

	if o.SetShouldStore {
		if o.ShouldStore {
			q.Set("store", "1")
		} else {
			q.Set("store", "0")
		}
	}

	if o.SetTtl {
		if o.Ttl > 0 {
			q.Set("ttl", strconv.Itoa(o.Ttl))
		}
	}

	q.Set("seqn", strconv.Itoa(<-o.pubnub.publishSequence))

	if o.DoNotReplicate == true {
		q.Set("norep", "true")
	}

	return q, nil
}

func (o *publishOpts) buildBody() ([]byte, error) {
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

func (o *publishOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	} else {
		return "GET"
	}
}

func (o *publishOpts) isAuthRequired() bool {
	return true
}

func (o *publishOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *publishOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *publishOpts) operationType() OperationType {
	return PNPublishOperation
}

func (o *publishOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}
