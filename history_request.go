package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
	"io/ioutil"
	"reflect"
	"strconv"

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

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *historyOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *historyOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
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

// parseInterface umarshals the response data, marshals the data again in a
// different format and returns the json string. It also unescapes the data.
//
// parameters:
// vv: interface array to parse and extract data from.
// o : historyOpts
//
// returns []HistoryResponseItem.
func parseInterface(vv []interface{}, o *historyOpts) []HistoryResponseItem {
	o.pubnub.Config.Log.Println(vv)
	items := make([]HistoryResponseItem, len(vv))
	for i := range vv {

		val := vv[i]
		o.pubnub.Config.Log.Println("reflect.TypeOf(val).Kind()", reflect.TypeOf(val).Kind())
		switch v := val.(type) {
		case map[string]interface{}:
			o.pubnub.Config.Log.Println("Map", v)
			if v["timetoken"] != nil {
				o.pubnub.Config.Log.Println("timetoken:", v["timetoken"])
				if f, ok := v["timetoken"].(float64); ok {
					s := fmt.Sprintf("%.0f", f)
					o.pubnub.Config.Log.Println("s:", s)

					if tt, err := strconv.ParseInt(s, 10, 64); err == nil {
						o.pubnub.Config.Log.Println("tt:", tt)
						items[i].Timetoken = int64(tt)
					} else {
						o.pubnub.Config.Log.Println(f, s, err)
					}
				} else {
					o.pubnub.Config.Log.Println("v[timetoken].(int64)", ok, items[i].Timetoken)
				}
				items[i].Message, _ = parseCipherInterface(v["message"], o.pubnub.Config)
			} else {
				o.pubnub.Config.Log.Println("value", v)
				items[i].Message, _ = parseCipherInterface(v, o.pubnub.Config)
				o.pubnub.Config.Log.Println("items[i]", items[i])
			}
			break
		default:
			o.pubnub.Config.Log.Println(v)
			items[i].Message, _ = parseCipherInterface(v, o.pubnub.Config)
			break
		}
	}
	return items
}

func newHistoryResponse(jsonBytes []byte, o *historyOpts,
	status StatusResponse) (*HistoryResponse, StatusResponse, error) {

	resp := &HistoryResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyHistoryResp, status, e
	}

	switch v := value.(type) {
	case []interface{}:
		startTimetoken, ok := v[1].(float64)
		if !ok {
			e := pnerr.NewResponseParsingError("Error parsing response",
				ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

			return emptyHistoryResp, status, e
		}

		endTimetoken, ok := v[2].(float64)
		if !ok {
			e := pnerr.NewResponseParsingError("Error parsing response",
				ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

			return emptyHistoryResp, status, e
		}

		msgs := v[0].([]interface{})
		o.pubnub.Config.Log.Println(msgs)

		items := parseInterface(msgs, o)
		if items != nil {
			resp.Messages = items
			o.pubnub.Config.Log.Printf("returning []interface, %v\n", items)
		} else {
			o.pubnub.Config.Log.Println("items nil")
		}

		resp.StartTimetoken = int64(startTimetoken)
		resp.EndTimetoken = int64(endTimetoken)
		break
	default:
		e := pnerr.NewResponseParsingError("Error parsing response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyHistoryResp, status, e
	}

	return resp, status, nil
}

// HistoryResponseItem is used to store the Message and the associated timetoken from the History request.
type HistoryResponseItem struct {
	Message   interface{}
	Timetoken int64
}
