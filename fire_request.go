package pubnub

import (
	"fmt"
	"strconv"

	"github.com/pubnub/go/v8/pnerr"
	"github.com/pubnub/go/v8/utils"

	"net/http"
	"net/url"
)

func newFireOpts(pubnub *PubNub, ctx Context) *fireOpts {
	return &fireOpts{
		endpointOpts: endpointOpts{
			pubnub: pubnub,
			ctx:    ctx,
		},
		Serialize: true,
	}
}

type fireOpts struct {
	endpointOpts

	TTL            int
	Channel        string
	Message        interface{}
	Meta           interface{}
	UsePost        bool
	Serialize      bool
	ShouldStore    bool
	DoNotReplicate bool
	Transport      http.RoundTripper
	QueryParam     map[string]string
	// nil hacks
	setTTL         bool
	setShouldStore bool
}

type fireBuilder struct {
	opts *fireOpts
}

func newFireBuilder(pubnub *PubNub) *fireBuilder {
	return newFireBuilderWithContext(pubnub, pubnub.ctx)
}

func newFireBuilderWithContext(pubnub *PubNub, context Context) *fireBuilder {
	return &fireBuilder{
		opts: newFireOpts(
			pubnub,
			context,
		),
	}
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

// GetLogParams returns the user-provided parameters for logging
func (o *fireOpts) GetLogParams() map[string]interface{} {
	params := map[string]interface{}{
		"Channel":   o.Channel,
		"UsePost":   o.UsePost,
		"Serialize": o.Serialize,
	}
	if o.setTTL {
		params["TTL"] = o.TTL
	}
	if o.Meta != nil {
		params["Meta"] = fmt.Sprintf("%v", o.Meta)
	}
	// Truncate message for logging
	if o.Message != nil {
		msgStr := fmt.Sprintf("%v", o.Message)
		if len(msgStr) > 100 {
			params["Message"] = msgStr[:100] + "... (truncated)"
		} else {
			params["Message"] = msgStr
		}
	}
	return params
}

// Execute runs the Fire request.
func (b *fireBuilder) Execute() (*PublishResponse, StatusResponse, error) {
	b.opts.ShouldStore = false
	b.opts.DoNotReplicate = true
	b.opts.pubnub.loggerManager.LogUserInput(PNLogLevelDebug, PNFireOperation, b.opts.GetLogParams(), true)

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishResponse, status, err
	}

	return newPublishResponse(rawJSON, status, b.opts.pubnub.loggerManager)
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

	if o.pubnub.getCryptoModule() != nil {
		var msg string
		if msg, err = serializeEncryptAndSerialize(o.pubnub.getCryptoModule(), o.Message, o.Serialize, o.pubnub.loggerManager); err != nil {
			o.pubnub.loggerManager.LogError(err, "FireMessageSerializationFailed", PNFireOperation, true)
			return "", err
		}
		message = []byte(msg)
	} else {
		message, err = utils.ValueAsString(o.Message)
		if err != nil {
			return "", err
		}
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

func (o *fireOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		var msg []byte

		if o.Serialize {
			m, err := utils.ValueAsString(o.Message)
			if err != nil {
				return []byte{}, err
			}
			msg = m
		} else {
			if s, ok := o.Message.(string); ok {
				msg = []byte(s)
			} else {
				err := pnerr.NewBuildRequestError("Type error, only string is expected")
				return []byte{}, err
			}
		}

		if o.pubnub.getCryptoModule() != nil {
			enc, err := encryptString(o.pubnub.getCryptoModule(), string(msg))
			if err != nil {
				return []byte{}, err
			}
			m, err := utils.ValueAsString(enc)
			if err != nil {
				return []byte{}, err
			}
			return m, nil
		}
		return msg, nil
	}
	return []byte{}, nil
}

func (o *fireOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
}

func (o *fireOpts) operationType() OperationType {
	return PNFireOperation
}
