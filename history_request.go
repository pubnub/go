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

const HISTORY_PATH = "/v2/history/sub-key/%s/channel/%s"
const MAX_COUNT = 100

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

func (b *historyBuilder) Channel(ch string) *historyBuilder {
	b.opts.Channel = ch
	return b
}

func (b *historyBuilder) Start(start int64) *historyBuilder {
	b.opts.Start = start
	b.opts.SetStart = true
	return b
}

func (b *historyBuilder) End(end int64) *historyBuilder {
	b.opts.End = end
	b.opts.SetEnd = true
	return b
}

func (b *historyBuilder) Count(count int) *historyBuilder {
	b.opts.Count = count
	return b
}

func (b *historyBuilder) Reverse(r bool) *historyBuilder {
	b.opts.Reverse = r
	return b
}

func (b *historyBuilder) IncludeTimetoken(i bool) *historyBuilder {
	b.opts.IncludeTimetoken = i
	return b
}

func (b *historyBuilder) Transport(tr http.RoundTripper) *historyBuilder {
	b.opts.Transport = tr
	return b
}

func (b *historyBuilder) Execute() (*HistoryResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyHistoryResp, status, err
	}

	return newHistoryResponse(rawJson, b.opts, status)
}

type historyOpts struct {
	pubnub *PubNub

	Channel string

	Start int64
	End   int64

	// defualt: 100
	Count int

	// default: false
	Reverse bool

	// default: false
	IncludeTimetoken bool

	// nil hacks
	SetStart bool
	SetEnd   bool

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
	return fmt.Sprintf(HISTORY_PATH,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.Channel)), nil
}

func (o *historyOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	if o.SetStart {
		q.Set("start", strconv.FormatInt(o.Start, 10))
	}

	if o.SetEnd {
		q.Set("end", strconv.FormatInt(o.End, 10))
	}

	if o.Count > 0 && o.Count <= MAX_COUNT {
		q.Set("count", strconv.Itoa(o.Count))
	} else {
		q.Set("count", "100")
	}

	q.Set("reverse", strconv.FormatBool(o.Reverse))
	q.Set("include_token", strconv.FormatBool(o.IncludeTimetoken))

	return q, nil
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

type HistoryResponse struct {
	Messages       []HistoryResponseItem
	StartTimetoken int64
	EndTimetoken   int64
}

// parseCipherInterface handles the decryption in case a cipher key is used
// in case of error it returns data as is.
//
// parameters
// data: the data to decrypt as interface.
// cipherKey: cipher key to use to decrypt.
//
// returns the decrypted data as interface.
func parseCipherInterface(data interface{}, cipherKey string) interface{} {
	if cipherKey != "" {
		var intf interface{}
		decrypted, errDecryption := utils.DecryptString(cipherKey, data.(string))
		if errDecryption != nil {
			intf = data
		} else {
			intf = decrypted
		}
		return intf
	} else {
		return data
	}
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
	cipherKey := o.pubnub.Config.CipherKey
	items := make([]HistoryResponseItem, len(vv))
	for i, _ := range vv {

		val := vv[i]
		o.pubnub.Config.Log.Println("reflect.TypeOf(val).Kind()", reflect.TypeOf(val).Kind())
		switch v := val.(type) {
		case map[string]interface{}:
			if v["timetoken"] != nil {
				for key, value := range v {
					if key == "timetoken" {
						tt := value.(float64)
						items[i].Timetoken = int64(tt)
					} else {
						items[i].Message = parseCipherInterface(value, cipherKey)
					}
				}
			} else {
				for _, value := range vv {
					items[i].Message = parseCipherInterface(value, cipherKey)
				}
			}
			break
		default:
			items[i].Message = parseCipherInterface(v, cipherKey)
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

		items := parseInterface(msgs, o)
		if items != nil {
			resp.Messages = items
			o.pubnub.Config.Log.Println("returning []interface, %s", items)
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

type HistoryResponseItem struct {
	Message   interface{}
	Timetoken int64
}

func unmarshalWithDecrypt(val string, cipherKey string) (interface{}, error) {
	v, err := utils.DecryptString(cipherKey, val)
	if err != nil {
		return nil, err
	}

	value := v.(string)

	var result interface{}
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
