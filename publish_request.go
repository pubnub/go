package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/pubnub/go/v5/pnerr"
	"github.com/pubnub/go/v5/utils"

	"net/http"
	"net/url"
)

const publishGetPath = "/publish/%s/%s/0/%s/%s/%s"
const publishPostPath = "/publish/%s/%s/0/%s/%s"

var emptyPublishResponse *PublishResponse

type publishOpts struct {
	pubnub *PubNub

	TTL     int
	Channel string
	Message interface{}
	Meta    interface{}

	UsePost        bool
	ShouldStore    bool
	Serialize      bool
	DoNotReplicate bool
	QueryParam     map[string]string

	Transport http.RoundTripper

	ctx Context

	// nil hacks
	setTTL         bool
	setShouldStore bool
}

// PublishResponse is the response after the execution on Publish and Fire operations.
type PublishResponse struct {
	Timestamp int64
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

	timeString, ok := value[2].(string)
	if !ok {
		return emptyPublishResponse, status, pnerr.NewResponseParsingError(fmt.Sprintf("Error unmarshalling response, %s %v", value[2], value), nil, nil)
	}
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
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
			pubnub:    pubnub,
			Serialize: true,
		},
	}

	return &builder
}

func newPublishBuilderWithContext(pubnub *PubNub, context Context) *publishBuilder {
	builder := publishBuilder{
		opts: &publishOpts{
			pubnub:    pubnub,
			ctx:       context,
			Serialize: true,
		},
	}

	return &builder
}

// TTL sets the TTL (hours) for the Publish request.
func (b *publishBuilder) TTL(ttl int) *publishBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

// Channel sets the Channel for the Publish request.
func (b *publishBuilder) Channel(ch string) *publishBuilder {
	b.opts.Channel = ch

	return b
}

// Message sets the Payload for the Publish request.
func (b *publishBuilder) Message(msg interface{}) *publishBuilder {
	b.opts.Message = msg

	return b
}

// Meta sets the Meta Payload for the Publish request.
func (b *publishBuilder) Meta(meta interface{}) *publishBuilder {
	b.opts.Meta = meta

	return b
}

// UsePost sends the Publish request using HTTP POST.
func (b *publishBuilder) UsePost(post bool) *publishBuilder {
	b.opts.UsePost = post

	return b
}

// ShouldStore if true the messages are stored in History
func (b *publishBuilder) ShouldStore(store bool) *publishBuilder {
	b.opts.ShouldStore = store
	b.opts.setShouldStore = true
	return b
}

// Serialize when true (default) serializes the payload before publish.
// Set to false if pre serialized payload is being used.
func (b *publishBuilder) Serialize(serialize bool) *publishBuilder {
	b.opts.Serialize = serialize

	return b
}

// DoNotReplicate stores the message in one DC.
func (b *publishBuilder) DoNotReplicate(repl bool) *publishBuilder {
	b.opts.DoNotReplicate = repl

	return b
}

// Transport sets the Transport for the Publish request.
func (b *publishBuilder) Transport(tr http.RoundTripper) *publishBuilder {
	b.opts.Transport = tr

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *publishBuilder) QueryParam(queryParam map[string]string) *publishBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Publish request.
func (b *publishBuilder) Execute() (*PublishResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishResponse, status, err
	}

	return newPublishResponse(rawJSON, status)
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

func (o *publishOpts) encryptProcessing(cipherKey string) (string, error) {
	var msg string
	var errJSONMarshal error

	o.pubnub.Config.Log.Println("EncryptString: encrypting", fmt.Sprintf("%s", o.Message))
	if o.pubnub.Config.DisablePNOtherProcessing {
		if msg, errJSONMarshal = utils.SerializeEncryptAndSerialize(o.Message, cipherKey, o.Serialize, o.pubnub.Config.UseRandomInitializationVector); errJSONMarshal != nil {
			o.pubnub.Config.Log.Printf("error in serializing: %v\n", errJSONMarshal)
			return "", errJSONMarshal
		}
	} else {
		//encrypt pn_other only
		o.pubnub.Config.Log.Println("encrypt pn_other only", "reflect.TypeOf(data).Kind()", reflect.TypeOf(o.Message).Kind(), o.Message)
		switch v := o.Message.(type) {
		case map[string]interface{}:

			msgPart, ok := v["pn_other"].(string)

			if ok {
				o.pubnub.Config.Log.Println(ok, msgPart)
				encMsg, errJSONMarshal := utils.SerializeAndEncrypt(msgPart, cipherKey, o.Serialize, o.pubnub.Config.UseRandomInitializationVector)
				if errJSONMarshal != nil {
					o.pubnub.Config.Log.Printf("error in serializing: %v\n", errJSONMarshal)
					return "", errJSONMarshal
				}
				v["pn_other"] = encMsg
				jsonEncBytes, errEnc := json.Marshal(v)
				if errEnc != nil {
					o.pubnub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
					return "", errEnc
				}
				msg = string(jsonEncBytes)
			} else {
				if msg, errJSONMarshal = utils.SerializeEncryptAndSerialize(o.Message, cipherKey, o.Serialize, o.pubnub.Config.UseRandomInitializationVector); errJSONMarshal != nil {
					o.pubnub.Config.Log.Printf("error in serializing: %v\n", errJSONMarshal)
					return "", errJSONMarshal
				}
			}
			break
		default:
			if msg, errJSONMarshal = utils.SerializeEncryptAndSerialize(o.Message, cipherKey, o.Serialize, o.pubnub.Config.UseRandomInitializationVector); errJSONMarshal != nil {
				o.pubnub.Config.Log.Printf("error in serializing: %v\n", errJSONMarshal)
				return "", errJSONMarshal
			}

			break
		}
	}
	return msg, nil
}

func (o *publishOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(publishPostPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0"), nil
	}

	var msg string
	var errJSONMarshal error

	if cipherKey := o.pubnub.Config.CipherKey; cipherKey != "" {
		if msg, errJSONMarshal = o.encryptProcessing(cipherKey); errJSONMarshal != nil {
			return "", errJSONMarshal
		}

		o.pubnub.Config.Log.Println("EncryptString: encrypted", msg)
	} else {
		if o.Serialize {
			jsonEncBytes, errEnc := json.Marshal(o.Message)
			if errEnc != nil {
				o.pubnub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
				return "", errEnc
			}
			msg = string(jsonEncBytes)
			o.pubnub.Config.Log.Println("len(jsonEncBytes)", len(jsonEncBytes))

		} else {
			if serializedMsg, ok := o.Message.(string); ok {
				msg = serializedMsg
			} else {
				return "", pnerr.NewBuildRequestError("buildpath: Message is not JSON serialized.")
			}
		}
	}

	return fmt.Sprintf(publishGetPath,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.Channel),
		"0",
		utils.URLEncode(msg)), nil
}

func (o *publishOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.Meta != nil {
		meta, err := utils.ValueAsString(o.Meta)
		if err != nil {
			return &url.Values{}, err
		}

		q.Set("meta", string(meta))
	}

	if o.setShouldStore {
		if o.ShouldStore {
			q.Set("store", "1")
		} else {
			q.Set("store", "0")
		}
	}

	if o.setTTL {
		if o.TTL > 0 {
			q.Set("ttl", strconv.Itoa(o.TTL))
		}
	}

	seqn := strconv.Itoa(o.pubnub.getPublishSequence())
	o.pubnub.Config.Log.Println("seqn:", seqn)
	q.Set("seqn", seqn)

	SetQueryParam(q, o.QueryParam)

	if o.DoNotReplicate == true {
		q.Set("norep", "true")
	}
	o.pubnub.Config.Log.Println(q)

	return q, nil
}

func (o *publishOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *publishOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		if cipherKey := o.pubnub.Config.CipherKey; cipherKey != "" {
			msg, errJSONMarshal := o.encryptProcessing(cipherKey)
			if errJSONMarshal != nil {
				return []byte{}, errJSONMarshal
			}
			return []byte(msg), nil
		}
		if o.Serialize {
			jsonEncBytes, errEnc := json.Marshal(o.Message)
			if errEnc != nil {
				o.pubnub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
				return []byte{}, errEnc
			}
			return jsonEncBytes, nil
		}
		serializedMsg, ok := o.Message.(string)
		if ok {
			return []byte(serializedMsg), nil
		}
		return []byte{}, pnerr.NewBuildRequestError("buildBody: Message is not JSON serialized.")
	}
	return []byte{}, nil
}

func (o *publishOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *publishOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
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
