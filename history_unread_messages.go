package pubnub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
	"reflect"
	"strings"

	"net/http"
	"net/url"
)

var emptyHistoryWithMessagesResp *HistoryWithMessagesResponse

const historyWithMessagesPath = "/v3/history/sub-key/%s/channels-with-messages/%s"

type historyWithMessagesBuilder struct {
	opts *historyWithMessagesOpts
}

func newHistoryWithMessagesBuilder(pubnub *PubNub) *historyWithMessagesBuilder {
	builder := historyWithMessagesBuilder{
		opts: &historyWithMessagesOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newHistoryWithMessagesBuilderWithContext(pubnub *PubNub,
	context Context) *historyWithMessagesBuilder {
	builder := historyWithMessagesBuilder{
		opts: &historyWithMessagesOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channels sets the Channels for the HistoryWithMessages request.
func (b *historyWithMessagesBuilder) Channels(channels []string) *historyWithMessagesBuilder {
	b.opts.Channels = channels
	return b
}

// Timetoken sets the number of items to return in the HistoryWithMessages request.
func (b *historyWithMessagesBuilder) Timetoken(timetoken string) *historyWithMessagesBuilder {
	b.opts.Timetoken = timetoken
	return b
}

// ChannelTimetokens sets the order of messages in the HistoryWithMessages request.
func (b *historyWithMessagesBuilder) ChannelTimetokens(channelTimetokens []string) *historyWithMessagesBuilder {
	b.opts.ChannelTimetokens = channelTimetokens
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *historyWithMessagesBuilder) QueryParam(queryParam map[string]string) *historyWithMessagesBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the HistoryWithMessages request.
func (b *historyWithMessagesBuilder) Transport(tr http.RoundTripper) *historyWithMessagesBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the HistoryWithMessages request.
func (b *historyWithMessagesBuilder) Execute() (*HistoryWithMessagesResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyHistoryWithMessagesResp, status, err
	}

	return newHistoryWithMessagesResponse(rawJSON, b.opts, status)
}

type historyWithMessagesOpts struct {
	pubnub *PubNub

	Channels          []string
	Timetoken         string
	ChannelTimetokens []string

	QueryParam map[string]string

	// nil hacks
	Transport http.RoundTripper

	ctx Context
}

func (o *historyWithMessagesOpts) config() Config {
	return *o.pubnub.Config
}

func (o *historyWithMessagesOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *historyWithMessagesOpts) context() Context {
	return o.ctx
}

func (o *historyWithMessagesOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) <= 0 {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *historyWithMessagesOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	return fmt.Sprintf(historyWithMessagesPath,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *historyWithMessagesOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

	q.Set("timetoken", o.Timetoken)
	q.Set("channelTimetokens", strings.Join(o.ChannelTimetokens, ","))
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *historyWithMessagesOpts) jobQueue() chan *JobQItem {
	return o.pubnub.jobQueue
}

func (o *historyWithMessagesOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *historyWithMessagesOpts) httpMethod() string {
	return "GET"
}

func (o *historyWithMessagesOpts) isAuthRequired() bool {
	return true
}

func (o *historyWithMessagesOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *historyWithMessagesOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *historyWithMessagesOpts) operationType() OperationType {
	return PNHistoryWithMessagesOperation
}

func (o *historyWithMessagesOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

// HistoryWithMessagesResponse is the response to HistoryWithMessages request. It contains a map of type HistoryWithMessagesResponseItem
type HistoryWithMessagesResponse struct {
	Channels map[string]int
}

//http://ps.pndsn.com/v3/history/sub-key/demo/channels-with-messages/my-channel,my-channel1?timestamp=1549982652&pnsdk=PubNub-Go/4.1.6&uuid=pn-82f145ea-adc3-4917-a11d-76a957347a82&timetoken=15499825804610610&channelTimetokens=15499825804610610,15499925804610615&auth=akey&signature=pVDVge_suepcOlSMllpsXg_jpOjtEpW7B3HHFaViI4s=
//{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}
func newHistoryWithMessagesResponse(jsonBytes []byte, o *historyWithMessagesOpts,
	status StatusResponse) (*HistoryWithMessagesResponse, StatusResponse, error) {

	resp := &HistoryWithMessagesResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyHistoryWithMessagesResp, status, e
	}

	if result, ok := value.(map[string]interface{}); ok {
		o.pubnub.Config.Log.Println(result["channels"])
		if channels, ok1 := result["channels"].(map[string]interface{}); ok1 {
			if channels != nil {
				resp.Channels = make(map[string]int)
				for ch, v := range channels {
					resp.Channels[ch] = int(v.(float64))
				}
			} else {
				o.pubnub.Config.Log.Printf("type assertion to map failed %v\n", result)
			}
		} else {
			o.pubnub.Config.Log.Println("Assertion failed", reflect.TypeOf(result["channels"]))
		}
	} else {
		o.pubnub.Config.Log.Printf("type assertion to map failed %v\n", value)
	}

	return resp, status, nil
}
