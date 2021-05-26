package pubnub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"

	"github.com/pubnub/go/v5/pnerr"
	"github.com/pubnub/go/v5/utils"

	"net/http"
	"net/url"
)

const historyPath = "/v2/history/sub-key/%s/channel/%s"
const maxCount = 100

var emptyHistoryResp *HistoryResponse

type historyBuilder struct {
	opts *historyOpts
}

func newHistoryBuilder(pubnub *PubNub) *historyBuilder {
	builder := historyBuilder{
		opts: &historyOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newHistoryBuilderWithContext(pubnub *PubNub,
	context Context) *historyBuilder {
	builder := historyBuilder{
		opts: &historyOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channel sets the Channel for the History request.
func (b *historyBuilder) Channel(ch string) *historyBuilder {
	b.opts.Channel = ch
	return b
}

// Start sets the Start Timetoken for the History request.
func (b *historyBuilder) Start(start int64) *historyBuilder {
	b.opts.Start = start
	b.opts.setStart = true
	return b
}

// End sets the End Timetoken for the History request.
func (b *historyBuilder) End(end int64) *historyBuilder {
	b.opts.End = end
	b.opts.setEnd = true
	return b
}

// Count sets the number of items to return in the History request.
func (b *historyBuilder) Count(count int) *historyBuilder {
	b.opts.Count = count
	return b
}

// Reverse sets the order of messages in the History request.
func (b *historyBuilder) Reverse(r bool) *historyBuilder {
	b.opts.Reverse = r
	return b
}

// IncludeTimetoken tells the server to send the timetoken associated with each history item.
func (b *historyBuilder) IncludeTimetoken(i bool) *historyBuilder {
	b.opts.IncludeTimetoken = i
	return b
}

// IncludeMeta fetches the meta data associated with the message
func (b *historyBuilder) IncludeMeta(withMeta bool) *historyBuilder {
	b.opts.WithMeta = withMeta
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *historyBuilder) QueryParam(queryParam map[string]string) *historyBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the History request.
func (b *historyBuilder) Transport(tr http.RoundTripper) *historyBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the History request.
func (b *historyBuilder) Execute() (*HistoryResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyHistoryResp, status, err
	}

	return newHistoryResponse(rawJSON, b.opts, status)
}

type historyOpts struct {
	pubnub *PubNub

	Channel string

	Start      int64
	End        int64
	QueryParam map[string]string
	WithMeta   bool

	// default: 100
	Count int

	// default: false
	Reverse bool

	// default: false
	IncludeTimetoken bool

	// nil hacks
	setStart bool
	setEnd   bool

	Transport http.RoundTripper

	ctx Context
}

func (o *historyOpts) config() Config {
	return *o.pubnub.Config
}

func (o *historyOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *historyOpts) context() Context {
	return o.ctx
}

func (o *historyOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *historyOpts) buildPath() (string, error) {
	return fmt.Sprintf(historyPath,
		o.pubnub.Config.SubscribeKey,
		utils.URLEncode(o.Channel)), nil
}

func (o *historyOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	if o.setStart {
		q.Set("start", strconv.FormatInt(o.Start, 10))
	}

	if o.setEnd {
		q.Set("end", strconv.FormatInt(o.End, 10))
	}

	if o.Count > 0 && o.Count <= maxCount {
		q.Set("count", strconv.Itoa(o.Count))
	} else {
		q.Set("count", "100")
	}

	q.Set("reverse", strconv.FormatBool(o.Reverse))
	q.Set("include_token", strconv.FormatBool(o.IncludeTimetoken))
	q.Set("include_meta", strconv.FormatBool(o.WithMeta))

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *historyOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *historyOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *historyOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *historyOpts) httpMethod() string {
	return "GET"
}

func (o *historyOpts) isAuthRequired() bool {
	return true
}

func (o *historyOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *historyOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *historyOpts) operationType() OperationType {
	return PNHistoryOperation
}

func (o *historyOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// HistoryResponse is used to store the response from the History request.
type HistoryResponse struct {
	Messages       []HistoryResponseItem
	StartTimetoken int64
	EndTimetoken   int64
}

// HistoryResponseItem is used to store the Message and the associated timetoken from the History request.
type HistoryResponseItem struct {
	Message   interface{}
	Meta      interface{}
	Timetoken int64
}

func logAndCreateNewResponseParsingError(o *historyOpts, err error, jsonBody string, message string) *pnerr.ResponseParsingError {
	o.pubnub.Config.Log.Println(err.Error())
	e := pnerr.NewResponseParsingError(message,
		ioutil.NopCloser(bytes.NewBufferString(jsonBody)), err)
	return e
}

func getHistoryItemsWithoutTimetoken(historyResponseRaw []byte, o *historyOpts, err1 error, jsonBytes []byte) ([]HistoryResponseItem, *pnerr.ResponseParsingError) {
	var historyResponseItems []interface{}
	err0 := json.Unmarshal(historyResponseRaw, &historyResponseItems)
	if err0 != nil {
		e := logAndCreateNewResponseParsingError(o, fmt.Errorf("%e, %e, %s", err0, err1, string(jsonBytes)), string(jsonBytes), "Error unmarshalling response")

		return nil, e
	}

	items := make([]HistoryResponseItem, len(historyResponseItems))

	for i, v := range historyResponseItems {
		o.pubnub.Config.Log.Println(v)
		items[i].Message, _ = parseCipherInterface(v, o.pubnub.Config)
	}
	return items, nil
}

func getHistoryItemsWithTimetoken(historyResponseItems []HistoryResponseItem, o *historyOpts, historyResponseRaw []byte, jsonBytes []byte) ([]HistoryResponseItem, *pnerr.ResponseParsingError) {
	items := make([]HistoryResponseItem, len(historyResponseItems))

	b := false

	for i, v := range historyResponseItems {
		if v.Message != nil {
			o.pubnub.Config.Log.Println(v.Message)
			items[i].Message, _ = parseCipherInterface(v.Message, o.pubnub.Config)

			o.pubnub.Config.Log.Println(v.Timetoken)
			items[i].Timetoken = v.Timetoken

			o.pubnub.Config.Log.Println(v.Meta)
			items[i].Meta = v.Meta
		} else {
			b = true
			break
		}
	}
	if b {
		items, e := getHistoryItemsWithoutTimetoken(historyResponseRaw, o, nil, jsonBytes)
		return items, e
	}

	return items, nil
}

func newHistoryResponse(jsonBytes []byte, o *historyOpts,
	status StatusResponse) (*HistoryResponse, StatusResponse, error) {

	resp := &HistoryResponse{}

	var historyResponseRaw []json.RawMessage

	err := json.Unmarshal(jsonBytes, &historyResponseRaw)
	if err != nil {
		e := logAndCreateNewResponseParsingError(o, err, string(jsonBytes), "Error unmarshalling response")

		return emptyHistoryResp, status, e
	}

	if historyResponseRaw != nil && len(historyResponseRaw) > 2 {
		o.pubnub.Config.Log.Println("M", string(historyResponseRaw[0]))
		o.pubnub.Config.Log.Println("T1", string(historyResponseRaw[1]))
		o.pubnub.Config.Log.Println("T2", string(historyResponseRaw[2]))

		var historyResponseItems []HistoryResponseItem
		var items []HistoryResponseItem

		err1 := json.Unmarshal(historyResponseRaw[0], &historyResponseItems)
		var e *pnerr.ResponseParsingError
		if err1 != nil {
			o.pubnub.Config.Log.Println(err1.Error())

			items, e = getHistoryItemsWithoutTimetoken(historyResponseRaw[0], o, err1, jsonBytes)
			if e != nil {
				return emptyHistoryResp, status, e
			}
		} else {
			items, e = getHistoryItemsWithTimetoken(historyResponseItems, o, historyResponseRaw[0], jsonBytes)
			if e != nil {
				return emptyHistoryResp, status, e
			}
		}
		if items != nil {
			resp.Messages = items
			o.pubnub.Config.Log.Printf("returning []interface, %v\n", items)
		} else {
			o.pubnub.Config.Log.Println("items nil")
		}

		startTimetoken, err := strconv.ParseInt(string(historyResponseRaw[1]), 10, 64)
		if err == nil {
			resp.StartTimetoken = startTimetoken
		}

		endTimetoken, err := strconv.ParseInt(string(historyResponseRaw[2]), 10, 64)
		if err == nil {
			resp.EndTimetoken = endTimetoken
		}
	} else if historyResponseRaw != nil && len(historyResponseRaw) > 0 {
		e := logAndCreateNewResponseParsingError(o, err, string(jsonBytes), "Error unmarshalling response")

		return emptyHistoryResp, status, e
	} else {
		e := logAndCreateNewResponseParsingError(o, err, string(jsonBytes), "Error unmarshalling response")

		return emptyHistoryResp, status, e
	}

	return resp, status, nil
}
