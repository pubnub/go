package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/sprucehealth/pubnub-go/pnerr"
	"github.com/sprucehealth/pubnub-go/utils"

	"net/http"
	"net/url"
)

var emptyFetchResp *FetchResponse

const fetchPath = "/v3/history/sub-key/%s/channel/%s"
const maxCountFetch = 25

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

	Start int64
	End   int64

	// default: 100
	Count int

	// default: false
	Reverse bool

	// default: false
	IncludeTimetoken bool
	QueryParam       map[string]string

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

	return nil
}

func (o *fetchOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

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

	if o.Count > 0 && o.Count <= maxCountFetch {
		q.Set("max", strconv.Itoa(o.Count))
	} else {
		q.Set("max", strconv.Itoa(maxCountFetch))
	}

	q.Set("reverse", strconv.FormatBool(o.Reverse))
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *fetchOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
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

// FetchResponse is the response to Fetch request. It contains a map of type FetchResponseItem
type FetchResponse struct {
	Messages map[string][]FetchResponseItem
}

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

// FetchResponseItem contains the message and the associated timetoken.
type FetchResponseItem struct {
	Message   interface{}
	Timetoken string
}
