package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pubnub/go/utils"
)

const historyDeletePath = "/v3/history/sub-key/%s/channel/%s"

var emptyHistoryDeleteResp *HistoryDeleteResponse

type historyDeleteBuilder struct {
	opts *historyDeleteOpts
}

func newHistoryDeleteBuilder(pubnub *PubNub) *historyDeleteBuilder {
	builder := historyDeleteBuilder{
		opts: &historyDeleteOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newHistoryDeleteBuilderWithContext(pubnub *PubNub,
	context Context) *historyDeleteBuilder {
	builder := historyDeleteBuilder{
		opts: &historyDeleteOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

func (b *historyDeleteBuilder) Channel(ch string) *historyDeleteBuilder {
	b.opts.Channel = ch
	return b
}

func (b *historyDeleteBuilder) Start(start int64) *historyDeleteBuilder {
	b.opts.Start = start
	b.opts.SetStart = true
	return b
}

func (b *historyDeleteBuilder) End(end int64) *historyDeleteBuilder {
	b.opts.End = end
	b.opts.SetEnd = true
	return b
}

func (b *historyDeleteBuilder) Transport(tr http.RoundTripper) *historyDeleteBuilder {
	b.opts.Transport = tr
	return b
}

func (b *historyDeleteBuilder) Execute() (*HistoryDeleteResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyHistoryDeleteResp, status, err
	}

	return emptyHistoryDeleteResp, status, nil
}

type historyDeleteOpts struct {
	pubnub *PubNub

	Channel string

	Start int64
	End   int64

	SetStart bool
	SetEnd   bool

	Transport http.RoundTripper

	ctx Context
}

func (o *historyDeleteOpts) config() Config {
	return *o.pubnub.Config
}

func (o *historyDeleteOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *historyDeleteOpts) context() Context {
	return o.ctx
}

func (o *historyDeleteOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().SecretKey == "" {
		return newValidationError(o, StrMissingSecretKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *historyDeleteOpts) buildPath() (string, error) {
	return fmt.Sprintf(historyDeletePath,
		o.pubnub.Config.SubscribeKey,
		utils.UrlEncode(o.Channel)), nil
}

func (o *historyDeleteOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid, o.pubnub.telemetryManager)

	if o.SetStart {
		q.Set("start", strconv.FormatInt(o.Start, 10))
	}

	if o.SetEnd {
		q.Set("end", strconv.FormatInt(o.End, 10))
	}

	return q, nil
}

func (o *historyDeleteOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *historyDeleteOpts) httpMethod() string {
	return "DELETE"
}

func (o *historyDeleteOpts) isAuthRequired() bool {
	return true
}

func (o *historyDeleteOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *historyDeleteOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

func (o *historyDeleteOpts) operationType() OperationType {
	return PNDeleteMessagesOperation
}

func (o *historyDeleteOpts) telemetryManager() *TelemetryManager {
	return o.pubnub.telemetryManager
}

type HistoryDeleteResponse struct {
}
