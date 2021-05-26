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

var emptyFetchResp *FetchResponse

const fetchPath = "/v3/history/sub-key/%s/channel/%s"
const historyWithMessageActionsPath = "/v3/history-with-actions/sub-key/%s/channel/%s"

const maxCountFetch = 100
const maxCountFetchMoreThanOneChannel = 25
const maxCountHistoryWithMessageActions = 25

type fetchBuilder struct {
	opts *fetchOpts
}

func newFetchBuilder(pubnub *PubNub) *fetchBuilder {
	builder := fetchBuilder{
		opts: &fetchOpts{
			pubnub:          pubnub,
			WithUUID:        true,
			WithMessageType: true,
		},
	}

	return &builder
}

func newFetchBuilderWithContext(pubnub *PubNub,
	context Context) *fetchBuilder {
	builder := fetchBuilder{
		opts: &fetchOpts{
			pubnub:          pubnub,
			ctx:             context,
			WithUUID:        true,
			WithMessageType: true,
		},
	}

	return &builder
}

// Channels sets the Channels for the Fetch request.
func (b *fetchBuilder) Channels(channels []string) *fetchBuilder {
	b.opts.Channels = channels
	return b
}

// Start sets the Start Timetoken for the Fetch request.
func (b *fetchBuilder) Start(start int64) *fetchBuilder {
	b.opts.Start = start
	b.opts.setStart = true
	return b
}

// End sets the End Timetoken for the Fetch request.
func (b *fetchBuilder) End(end int64) *fetchBuilder {
	b.opts.End = end
	b.opts.setEnd = true
	return b
}

// Count sets the number of items to return in the Fetch request.
func (b *fetchBuilder) Count(count int) *fetchBuilder {
	b.opts.Count = count
	return b
}

// Reverse sets the order of messages in the Fetch request.
func (b *fetchBuilder) Reverse(r bool) *fetchBuilder {
	b.opts.Reverse = r
	return b
}

// IncludeMeta fetches the meta data associated with the message
func (b *fetchBuilder) IncludeMeta(withMeta bool) *fetchBuilder {
	b.opts.WithMeta = withMeta
	return b
}

// IncludeMessageActions fetches the actions associated with the message
func (b *fetchBuilder) IncludeMessageActions(withMessageActions bool) *fetchBuilder {
	b.opts.WithMessageActions = withMessageActions
	return b
}

// IncludeUUID fetches the UUID associated with the message
func (b *fetchBuilder) IncludeUUID(withUUID bool) *fetchBuilder {
	b.opts.WithUUID = withUUID
	return b
}

// IncludeMessageType fetches the Message Type associated with the message
func (b *fetchBuilder) IncludeMessageType(withMessageType bool) *fetchBuilder {
	b.opts.WithMessageType = withMessageType
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *fetchBuilder) QueryParam(queryParam map[string]string) *fetchBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the Fetch request.
func (b *fetchBuilder) Transport(tr http.RoundTripper) *fetchBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the Fetch request.
func (b *fetchBuilder) Execute() (*FetchResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyFetchResp, status, err
	}

	return newFetchResponse(rawJSON, b.opts, status)
}

type fetchOpts struct {
	pubnub *PubNub

	Channels []string

	Start              int64
	End                int64
	WithMessageActions bool
	WithMeta           bool
	WithUUID           bool
	WithMessageType    bool

	// default: 100
	Count int

	// default: false
	Reverse bool

	QueryParam map[string]string

	// nil hacks
	setStart bool
	setEnd   bool

	Transport http.RoundTripper

	ctx Context
}

func (o *fetchOpts) config() Config {
	return *o.pubnub.Config
}

func (o *fetchOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *fetchOpts) context() Context {
	return o.ctx
}

func (o *fetchOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) <= 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.WithMessageActions && len(o.Channels) > 1 {
		return newValidationError(o, "Only one channel is supported when WithMessageActions is true")
	}

	return nil
}

func (o *fetchOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	if o.WithMessageActions {
		return fmt.Sprintf(historyWithMessageActionsPath,
			o.pubnub.Config.SubscribeKey,
			channels), nil
	}
	return fmt.Sprintf(fetchPath,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *fetchOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.setStart {
		q.Set("start", strconv.FormatInt(o.Start, 10))
	}

	if o.setEnd {
		q.Set("end", strconv.FormatInt(o.End, 10))
	}

	maxCount := maxCountFetch

	if o.WithMessageActions {
		maxCount = maxCountHistoryWithMessageActions
	}

	if len(o.Channels) > 1 {
		maxCount = maxCountFetchMoreThanOneChannel
	}

	if o.Count > 0 && o.Count <= maxCount {
		q.Set("max", strconv.Itoa(o.Count))
	} else {
		q.Set("max", strconv.Itoa(maxCount))
	}

	q.Set("reverse", strconv.FormatBool(o.Reverse))
	q.Set("include_meta", strconv.FormatBool(o.WithMeta))
	q.Set("include_message_type", strconv.FormatBool(o.WithMessageType))
	q.Set("include_uuid", strconv.FormatBool(o.WithUUID))

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *fetchOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *fetchOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *fetchOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *fetchOpts) httpMethod() string {
	return "GET"
}

func (o *fetchOpts) isAuthRequired() bool {
	return true
}

func (o *fetchOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *fetchOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *fetchOpts) operationType() OperationType {
	return PNFetchMessagesOperation
}

func (o *fetchOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

func (o *fetchOpts) parseMessageActions(actions interface{}) map[string]PNHistoryMessageActionsTypeMap {
	o.pubnub.Config.Log.Println(actions)
	resp := make(map[string]PNHistoryMessageActionsTypeMap)

	if actions != nil {
		actionsMap := actions.(map[string]interface{})

		for actionType, action := range actionsMap {

			o.pubnub.Config.Log.Println("action:", action)
			o.pubnub.Config.Log.Println("actionType:", actionType) //reaction2

			actionMap := action.(map[string]interface{})

			if actionMap != nil {
				messageActionsTypeMap := PNHistoryMessageActionsTypeMap{}
				messageActionsTypeMap.ActionsTypeValues = make(map[string][]PNHistoryMessageActionTypeVal, len(actionMap))
				for actionVal, val := range actionMap {
					o.pubnub.Config.Log.Println("actionVal:", actionVal) // smiley_face
					o.pubnub.Config.Log.Println("val:", val)

					actionValInt := val.([]interface{})
					if actionValInt != nil {
						params := make([]PNHistoryMessageActionTypeVal, len(actionValInt))
						pCount := 0
						for _, actionParam := range actionValInt {

							pv := PNHistoryMessageActionTypeVal{}
							for actionParamName, actionParamVal := range actionParam.(map[string]interface{}) {
								o.pubnub.Config.Log.Println("actionParamName", actionParamName)
								o.pubnub.Config.Log.Println("actionParamVal", actionParamVal)
								switch actionParamName {
								case "uuid":
									pv.UUID = actionParamVal.(string)
								case "actionTimetoken":
									pv.ActionTimetoken = actionParamVal.(string)
								}
							}
							params[pCount] = pv
							pCount++
						}
						messageActionsTypeMap.ActionsTypeValues[actionVal] = params
					}
				}
				resp[actionType] = messageActionsTypeMap
			}

		}
	}

	return resp
}

//{"status": 200, "error": false, "error_message": "", "channels": {"ch1":[{"message_type": "", "message": {"text": "hey"}, "timetoken": "15959610984115342", "meta": "", "uuid": "db9c5e39-7c95-40f5-8d71-125765b6f561"}]}}
func (o *fetchOpts) fetchMessages(channels map[string]interface{}) map[string][]FetchResponseItem {
	messages := make(map[string][]FetchResponseItem, len(channels))

	for channel, histResponseSliceMap := range channels {
		if histResponseMap, ok2 := histResponseSliceMap.([]interface{}); ok2 {
			o.pubnub.Config.Log.Printf("Channel:%s, count:%d\n", channel, len(histResponseMap))
			items := make([]FetchResponseItem, len(histResponseMap))
			count := 0

			for _, val := range histResponseMap {
				if histResponse, ok3 := val.(map[string]interface{}); ok3 {
					msg, _ := parseCipherInterface(histResponse["message"], o.pubnub.Config)

					histItem := FetchResponseItem{
						Message:   msg,
						Timetoken: histResponse["timetoken"].(string),
						Meta:      histResponse["meta"],
					}
					if d, ok := histResponse["message_type"]; ok {
						switch v := d.(type) {
						case float64:
							histItem.MessageType = int(v)
						case string:
							t, err := strconv.ParseInt(v, 10, 64)
							if err == nil {
								histItem.MessageType = int(t)
							} else {
								o.pubnub.Config.Log.Printf("MessageType conversion error.")
							}
						default:
							o.pubnub.Config.Log.Printf("histResponse message_type type %vv", d)
							if v != nil {
								o.pubnub.Config.Log.Printf("histResponse message_type type %vv", reflect.TypeOf(v).Kind())
							} else {
								o.pubnub.Config.Log.Printf("histResponse message_type nil")
							}
						}
					}
					if d, ok := histResponse["uuid"]; ok {
						histItem.UUID = d.(string)
					}
					histItem.MessageActions = o.parseMessageActions(histResponse["actions"])
					if filesPayload, okFile := msg.(map[string]interface{}); okFile {
						f, m := ParseFileInfo(filesPayload)

						if f.Name != "" && f.ID != "" {
							histItem.File = f
							histItem.Message = m
						}
					}

					items[count] = histItem
					o.pubnub.Config.Log.Printf("Channel:%s, count:%d %d\n", channel, count, len(items))
					count++
				} else {
					o.pubnub.Config.Log.Printf("histResponse not a map %v\n", histResponse)
					continue
				}
			}
			messages[channel] = items
			o.pubnub.Config.Log.Printf("Channel:%s, count:%d\n", channel, len(messages[channel]))
		} else {
			o.pubnub.Config.Log.Printf("histResponseSliceMap not an []interface %v\n", histResponseSliceMap)
			continue
		}
	}
	return messages
}

func newFetchResponse(jsonBytes []byte, o *fetchOpts,
	status StatusResponse) (*FetchResponse, StatusResponse, error) {

	resp := &FetchResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyFetchResp, status, e
	}

	if result, ok := value.(map[string]interface{}); ok {
		o.pubnub.Config.Log.Println(result["channels"])
		if channels, ok1 := result["channels"].(map[string]interface{}); ok1 {
			if channels != nil {
				resp.Messages = o.fetchMessages(channels)
			} else {
				o.pubnub.Config.Log.Printf("type assertion to map failed %v\n", result)
			}
		}
	} else {
		o.pubnub.Config.Log.Printf("type assertion to map failed %v\n", value)
	}

	return resp, status, nil
}

// FetchResponse is the response to Fetch request. It contains a map of type FetchResponseItem
type FetchResponse struct {
	Messages map[string][]FetchResponseItem
}

// FetchResponseItem contains the message and the associated timetoken.
type FetchResponseItem struct {
	Message        interface{}                               `json:"message"`
	Meta           interface{}                               `json:"meta"`
	MessageActions map[string]PNHistoryMessageActionsTypeMap `json:"actions"`
	File           PNFileDetails                             `json:"file"`
	Timetoken      string                                    `json:"timetoken"`
	UUID           string                                    `json:"uuid"`
	MessageType    int                                       `json:"message_type"`
}

// PNHistoryMessageActionsTypeMap is the struct used in the Fetch request that includes Message Actions
type PNHistoryMessageActionsTypeMap struct {
	ActionsTypeValues map[string][]PNHistoryMessageActionTypeVal `json:"-"`
}

// PNHistoryMessageActionTypeVal is the struct used in the Fetch request that includes Message Actions
type PNHistoryMessageActionTypeVal struct {
	UUID            string `json:"uuid"`
	ActionTimetoken string `json:"actionTimetoken"`
}
