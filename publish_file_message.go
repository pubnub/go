package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pubnub/go/v7/pnerr"
	"github.com/pubnub/go/v7/utils"

	"net/http"
	"net/url"
	"strconv"
)

var emptyPublishFileMessageResponse *PublishFileMessageResponse

const publishFileMessageGetPath = "/v1/files/publish-file/%s/%s/0/%s/%s/%s"
const publishFileMessagePostPath = "/v1/files/publish-file/%s/%s/0/%s/%s"

type publishFileMessageBuilder struct {
	opts *publishFileMessageOpts
}

func newPublishFileMessageBuilder(pubnub *PubNub) *publishFileMessageBuilder {
	return newPublishFileMessageBuilderWithContext(pubnub, pubnub.ctx)
}

func newPublishFileMessageOpts(pubnub *PubNub, ctx Context) *publishFileMessageOpts {
	return &publishFileMessageOpts{endpointOpts: endpointOpts{pubnub: pubnub, ctx: ctx}}
}
func newPublishFileMessageBuilderWithContext(pubnub *PubNub,
	context Context) *publishFileMessageBuilder {
	builder := publishFileMessageBuilder{
		opts: newPublishFileMessageOpts(pubnub, context)}
	return &builder
}

// TTL sets the TTL (hours) for the Publish request.
func (b *publishFileMessageBuilder) TTL(ttl int) *publishFileMessageBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

// Meta sets the Meta Payload for the Publish request.
func (b *publishFileMessageBuilder) Meta(meta interface{}) *publishFileMessageBuilder {
	b.opts.Meta = meta

	return b
}

// ShouldStore if true the messages are stored in History
func (b *publishFileMessageBuilder) ShouldStore(store bool) *publishFileMessageBuilder {
	b.opts.ShouldStore = store
	b.opts.setShouldStore = true
	return b
}

// Channel sets the Channel for the PublishFileMessage request.
func (b *publishFileMessageBuilder) Channel(channel string) *publishFileMessageBuilder {
	b.opts.Channel = channel
	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) FileName(name string) *publishFileMessageBuilder {
	b.opts.FileName = name

	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) Message(msg PNPublishFileMessage) *publishFileMessageBuilder {
	b.opts.Message = msg
	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) MessageText(msg string) *publishFileMessageBuilder {
	b.opts.MessageText = msg

	return b
}

// Message sets the Payload for the PublishFileMessage request.
func (b *publishFileMessageBuilder) FileID(id string) *publishFileMessageBuilder {
	b.opts.FileID = id

	return b
}

// usePost sends the PublishFileMessage request using HTTP POST. Not implemented
func (b *publishFileMessageBuilder) usePost(post bool) *publishFileMessageBuilder {
	b.opts.UsePost = post

	return b
}

// Transport sets the Transport for the objectAPICreateUsers request.
func (b *publishFileMessageBuilder) Transport(tr http.RoundTripper) *publishFileMessageBuilder {
	b.opts.Transport = tr
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *publishFileMessageBuilder) QueryParam(queryParam map[string]string) *publishFileMessageBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// UseRawText sets whether the message text should be sent as raw text instead of being wrapped in a JSON "text" field.
// When true, the message will be sent directly without the {"text": "message"} wrapper.
// Defaults to false for backward compatibility.
// Example:
//
//	UseRawText(false): {"message": {"text": "Hello"}, "file": {"id": "123", "name": "file.txt"}}
//	UseRawText(true):  {"message": "Hello", "file": {"id": "123", "name": "file.txt"}}
func (b *publishFileMessageBuilder) UseRawText(useRawText bool) *publishFileMessageBuilder {
	b.opts.UseRawText = useRawText

	return b
}

// Execute runs the PublishFileMessage request.
func (b *publishFileMessageBuilder) Execute() (*PublishFileMessageResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyPublishFileMessageResponse, status, err
	}

	return newPublishFileMessageResponse(rawJSON, b.opts, status)
}

type publishFileMessageOpts struct {
	endpointOpts
	Message        interface{}
	Channel        string
	UsePost        bool
	TTL            int
	Meta           interface{}
	ShouldStore    bool
	setTTL         bool
	setShouldStore bool
	MessageText    string
	FileID         string
	FileName       string
	QueryParam     map[string]string
	Transport      http.RoundTripper
	UseRawText     bool
}

func (o *publishFileMessageOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if (o.Message == nil) && (o.FileID == "") {
		return newValidationError(o, StrMissingFileID)
	}

	if (o.Message == nil) && (o.FileName == "") {
		return newValidationError(o, StrMissingFileName)
	}

	if o.Message != nil {
		if filesPayload, okFile := o.Message.(PNPublishFileMessage); okFile {
			if filesPayload.PNFile != nil {
				if filesPayload.PNFile.ID == "" {
					return newValidationError(o, StrMissingFileID)
				}
				if filesPayload.PNFile.Name == "" {
					return newValidationError(o, StrMissingFileName)
				}
			} else {
				return newValidationError(o, StrMissingFileID)
			}
		} else if filesPayloadRaw, okFileRaw := o.Message.(PNPublishFileMessageRaw); okFileRaw {
			if filesPayloadRaw.PNFile != nil {
				if filesPayloadRaw.PNFile.ID == "" {
					return newValidationError(o, StrMissingFileID)
				}
				if filesPayloadRaw.PNFile.Name == "" {
					return newValidationError(o, StrMissingFileName)
				}
			} else {
				return newValidationError(o, StrMissingFileID)
			}
			// Set UseRawText to true when a raw message is passed
			o.UseRawText = true
		} else {
			return newValidationError(o, StrMissingMessage)
		}
	}

	return nil
}

// buildRawTextMessage creates the appropriate message structure for raw text mode
func (o *publishFileMessageOpts) buildRawTextMessage() interface{} {
	if filesPayload, ok := o.Message.(PNPublishFileMessage); ok && filesPayload.PNMessage != nil {
		return map[string]interface{}{
			"message": filesPayload.PNMessage.Text,
			"file": map[string]interface{}{
				"id":   filesPayload.PNFile.ID,
				"name": filesPayload.PNFile.Name,
			},
		}
	}
	if filesPayloadRaw, ok := o.Message.(PNPublishFileMessageRaw); ok && filesPayloadRaw.PNMessage != nil {
		return map[string]interface{}{
			"message": filesPayloadRaw.PNMessage.Text,
			"file": map[string]interface{}{
				"id":   filesPayloadRaw.PNFile.ID,
				"name": filesPayloadRaw.PNFile.Name,
			},
		}
	}
	// Fallback: construct message from individual fields (MessageText, FileID, FileName)
	return map[string]interface{}{
		"message": o.MessageText,
		"file": map[string]interface{}{
			"id":   o.FileID,
			"name": o.FileName,
		},
	}
}

func (o *publishFileMessageOpts) buildPath() (string, error) {
	if o.UsePost {
		return fmt.Sprintf(publishFileMessagePostPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0"), nil
	}

	if o.Message == nil {
		file := &PNFileInfoForPublish{
			ID:   o.FileID,
			Name: o.FileName,
		}

		o.Message = PNPublishFileMessage{
			PNFile: file,
			PNMessage: &PNPublishMessage{
				Text: o.MessageText,
			},
		}
	}

	var messageToProcess interface{}
	if o.UseRawText {
		messageToProcess = o.buildRawTextMessage()
	} else {
		messageToProcess = o.Message
	}

	if o.pubnub.getCryptoModule() != nil {
		var msg string
		var p *publishBuilder
		if o.context() != nil {
			p = newPublishBuilderWithContext(o.pubnub, o.context())
		} else {
			p = newPublishBuilder(o.pubnub)
		}

		p.opts.Message = messageToProcess
		msg, errJSONMarshal := p.opts.encryptProcessing()
		if errJSONMarshal != nil {
			return "", errJSONMarshal
		}

		o.pubnub.Config.Log.Println("EncryptString: encrypted", msg)
		return fmt.Sprintf(publishFileMessageGetPath,
			o.pubnub.Config.PublishKey,
			o.pubnub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0",
			utils.URLEncode(msg)), nil
	}

	jsonEncBytes, errEnc := json.Marshal(messageToProcess)
	if errEnc != nil {
		o.pubnub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
		return "", errEnc
	}
	msg := string(jsonEncBytes)
	return fmt.Sprintf(publishFileMessageGetPath,
		o.pubnub.Config.PublishKey,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.Channel),
		"0",
		utils.URLEncode(msg),
	), nil

}

func (o *publishFileMessageOpts) buildQuery() (*url.Values, error) {
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

	return q, nil
}

func (o *publishFileMessageOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		jsonEncBytes, errEnc := json.Marshal(o.Message)
		if errEnc != nil {
			o.pubnub.Config.Log.Printf("ERROR: PublishFileMessage error: %s\n", errEnc.Error())
			return []byte{}, errEnc
		}
		return jsonEncBytes, nil
	}
	return []byte{}, nil
}

func (o *publishFileMessageOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
}

func (o *publishFileMessageOpts) isAuthRequired() bool {
	return true
}

func (o *publishFileMessageOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *publishFileMessageOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *publishFileMessageOpts) operationType() OperationType {
	return PNPublishFileMessageOperation
}

// PublishFileMessageResponse is the response to PublishFileMessage request.
type PublishFileMessageResponse struct {
	Timestamp int64
}

func newPublishFileMessageResponse(jsonBytes []byte, o *publishFileMessageOpts,
	status StatusResponse) (*PublishFileMessageResponse, StatusResponse, error) {

	resp := &PublishFileMessageResponse{}

	var value []interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			io.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyPublishFileMessageResponse, status, e
	}

	if len(value) > 1 {
		timeString, ok := value[2].(string)
		if !ok {
			return emptyPublishFileMessageResponse, status, pnerr.NewResponseParsingError(fmt.Sprintf("Error unmarshalling response 2, %s %v", value[2], value), nil, nil)
		}
		timestamp, err := strconv.ParseInt(timeString, 10, 64)
		if err != nil {
			return emptyPublishFileMessageResponse, status, err
		}

		return &PublishFileMessageResponse{
			Timestamp: timestamp,
		}, status, nil
	}

	return resp, status, nil
}
