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

func PublishRequest(pn *PubNub, opts *PublishOpts) (PublishResponse, error) {
	opts.pubnub = pn
	resp, err := executeRequest(opts)
	if err != nil {
		return PublishResponse{}, err
	}
	var value []interface{}

	err = json.Unmarshal(resp, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(resp))), err)

		return PublishResponse{}, e
	}

	timeString := value[2].(string)
	timestamp, err := strconv.Atoi(timeString)
	if err != nil {
		return PublishResponse{}, err
	}

	return PublishResponse{
		Timestamp: timestamp,
	}, nil
}

func PublishRequestWithContext(ctx Context,
	pn *PubNub, opts *PublishOpts) (PublishResponse, error) {
	opts.pubnub = pn
	opts.ctx = ctx

	_, err := executeRequest(opts)
	if err != nil {
		return PublishResponse{}, err
	}

	return PublishResponse{
		Timestamp: 123,
	}, nil
}

type PublishOpts struct {
	pubnub *PubNub

	Ttl     int
	Channel string
	Message interface{}
	Meta    interface{}

	UsePost        bool
	DoNotStore     bool
	Serialize      bool
	DoNotReplicate bool

	Transport http.RoundTripper

	ctx Context
}

type PublishResponse struct {
	Timestamp int
}

func (o *PublishOpts) config() Config {
	return *o.pubnub.Config
}

func (o *PublishOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *PublishOpts) context() Context {
	return o.ctx
}

func (o *PublishOpts) validate() error {
	if o.config().PublishKey == "" {
		return ErrMissingPubKey
	}

	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if o.Channel == "" {
		return ErrMissingChannel
	}

	if o.Message == nil {
		return ErrMissingMessage
	}

	return nil
}

func (o *PublishOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(PUBLISH_POST_PATH,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			o.Channel,
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

	stringifiedMessage := utils.UrlEncode(string(message))

	return fmt.Sprintf(PUBLISH_GET_PATH,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		o.Channel,
		"0",
		stringifiedMessage), nil
}

func (o *PublishOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if o.Meta != nil {
		// TODO: serialize
		meta, _ := utils.ValueAsString(o.Meta)
		q.Set("meta", string(meta))
	}

	if o.DoNotStore {
		q.Set("store", "1")
	} else {
		q.Set("store", "0")
	}

	// TODO: 0 value?
	if o.Ttl > 0 {
		q.Set("ttl", strconv.Itoa(o.Ttl))
	}

	q.Set("seqn", strconv.Itoa(<-o.pubnub.publishSequence))

	if o.DoNotReplicate == true {
		q.Set("norep", "true")
	}

	return q, nil
}

func (o *PublishOpts) buildBody() ([]byte, error) {
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

func (o *PublishOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	} else {
		return "GET"
	}
}

func (o *PublishOpts) isAuthRequired() bool {
	return true
}

func (o *PublishOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *PublishOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *PublishOpts) operationType() PNOperationType {
	return PNPublishOperation
}
