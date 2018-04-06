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

var emptyFetchResp *FetchResponse

const FETCH_PATH = "/v3/history/sub-key/%s/channel/%s"
const MAX_COUNT_FETCH = 25

type fetchBuilder struct {
	opts *fetchOpts
}

func newFetchBuilder(pubnub *PubNub) *fetchBuilder {
	builder := fetchBuilder{
		opts: &fetchOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newFetchBuilderWithContext(pubnub *PubNub,
	context Context) *fetchBuilder {
	builder := fetchBuilder{
		opts: &fetchOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *fetchBuilder) Channels(channels []string) *fetchBuilder {
	b.opts.Channels = channels
	return b
}

func (b *fetchBuilder) Start(start int64) *fetchBuilder {
	b.opts.Start = start
	b.opts.SetStart = true
	return b
}

func (b *fetchBuilder) End(end int64) *fetchBuilder {
	b.opts.End = end
	b.opts.SetEnd = true
	return b
}

func (b *fetchBuilder) Count(count int) *fetchBuilder {
	b.opts.Count = count
	return b
}

func (b *fetchBuilder) Reverse(r bool) *fetchBuilder {
	b.opts.Reverse = r
	return b
}

/*func (b *fetchBuilder) IncludeTimetoken(i bool) *fetchBuilder {
	b.opts.IncludeTimetoken = i
	return b
}*/

func (b *fetchBuilder) Transport(tr http.RoundTripper) *fetchBuilder {
	b.opts.Transport = tr
	return b
}

func (b *fetchBuilder) Execute() (*FetchResponse, StatusResponse, error) {
	rawJson, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyFetchResp, status, err
	}

	return newFetchResponse(rawJson, b.opts, status)
}

type fetchOpts struct {
	pubnub *PubNub

	Channels []string

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

	return nil
}

func (o *fetchOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	return fmt.Sprintf(FETCH_PATH,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *fetchOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	if o.SetStart {
		q.Set("start", strconv.FormatInt(o.Start, 10))
	}

	if o.SetEnd {
		q.Set("end", strconv.FormatInt(o.End, 10))
	}

	if o.Count > 0 && o.Count <= MAX_COUNT_FETCH {
		q.Set("max", strconv.Itoa(o.Count))
	} else {
		q.Set("max", strconv.Itoa(MAX_COUNT_FETCH))
	}

	q.Set("reverse", strconv.FormatBool(o.Reverse))
	//q.Set("include_token", strconv.FormatBool(o.IncludeTimetoken))

	return q, nil
}

func (o *fetchOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
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

type FetchResponse struct {
	Messages map[string][]FetchResponseItem
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

	result := value.(map[string]interface{})

	if result != nil {
		o.pubnub.Config.Log.Println(result["channels"])
		channels := result["channels"].(map[string]interface{})

		if channels != nil {
			messages := make(map[string][]FetchResponseItem, len(channels))

			for channel, histResponseSliceMap := range channels {
				histResponseMap := histResponseSliceMap.([]interface{})
				o.pubnub.Config.Log.Printf("Channel:%s, count:%d", channel, len(histResponseMap))
				items := make([]FetchResponseItem, len(histResponseMap))
				count := 0

				for _, val := range histResponseMap {
					histResponse := val.(map[string]interface{})

					histItem := FetchResponseItem{
						Message:   parseCipherInterface(histResponse["message"], o.pubnub.Config.CipherKey),
						Timetoken: histResponse["timetoken"].(string),
					}
					items[count] = histItem
					o.pubnub.Config.Log.Printf("Channel:%s, count:%d %d", channel, count, len(items))
					count++
				}
				messages[channel] = items
				o.pubnub.Config.Log.Printf("Channel:%s, count:%d", channel, len(messages[channel]))
			}
			resp.Messages = messages
		}
	}

	return resp, status, nil
}

type FetchResponseItem struct {
	Message   interface{}
	Timetoken string
}
