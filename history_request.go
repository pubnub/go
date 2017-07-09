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

const HISTORY_PATH = "/v2/history/sub-key/%s/channel/%s"
const MAX_COUNT = 100

var emptyHistoryResp *HistoryResponse

func HistoryRequest(pn *PubNub, opts *HistoryOpts) (*HistoryResponse, error) {
	opts.pubnub = pn
	rawJson, err := executeRequest(opts)
	if err != nil {
		return emptyHistoryResp, err
	}

	fmt.Println(string(rawJson))

	// TODO: just return the function call after finishing implementation
	r, err := newHistoryResponse(rawJson, opts.config().CipherKey)
	if err != nil {
		return emptyHistoryResp, nil
	}

	return r, nil
}

func HistoryRequestWithContext(ctx Context,
	pn *PubNub, opts *HistoryOpts) (*HistoryResponse, error) {
	opts.pubnub = pn
	opts.ctx = ctx

	_, err := executeRequest(opts)
	if err != nil {
		return emptyHistoryResp, err
	}

	return emptyHistoryResp, nil
}

type HistoryOpts struct {
	pubnub *PubNub

	Channel string

	// Stringified timetoken, default: not set
	Start string

	// Stringified timetoken, default: not set
	End string

	// defualt: 100
	Count int

	// default: false
	Reverse bool

	// default: false
	IncludeTimetoken bool

	Transport http.RoundTripper

	ctx Context
}

func (o *HistoryOpts) config() Config {
	return *o.pubnub.Config
}

func (o *HistoryOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *HistoryOpts) context() Context {
	return o.ctx
}

func (o *HistoryOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if o.Channel == "" {
		return ErrMissingChannel
	}

	return nil
}

func (o *HistoryOpts) buildPath() (string, error) {
	return fmt.Sprintf(HISTORY_PATH,
		o.pubnub.Config.SubscribeKey,
		o.Channel), nil
}

func (o *HistoryOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if o.Start != "" {
		i, err := strconv.Atoi(o.Start)
		if err != nil {
			// TODO: wrap error
			return nil, err
		}

		q.Set("start", strconv.Itoa(i))
	}

	if o.End != "" {
		i, err := strconv.Atoi(o.End)
		if err != nil {
			// TODO: wrap error
			return nil, err
		}

		q.Set("end", strconv.Itoa(i))
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

func (o *HistoryOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *HistoryOpts) httpMethod() string {
	return "GET"
}

func (o *HistoryOpts) isAuthRequired() bool {
	return true
}

func (o *HistoryOpts) requestTimeout() int {
	return o.pubnub.Config.NonSubscribeRequestTimeout
}

func (o *HistoryOpts) connectTimeout() int {
	return o.pubnub.Config.ConnectTimeout
}

type HistoryResponse struct {
	Messages       []HistoryResponseItem
	StartTimetoken int64
	EndTimetoken   int64
}

func newHistoryResponse(jsonBytes []byte, cipherKey string) (*HistoryResponse, error) {
	resp := &HistoryResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyHistoryResp, e
	}

	switch v := value.(type) {
	case []interface{}:
		msgs := v[0].([]interface{})
		items := make([]HistoryResponseItem, len(msgs))

		for k, val := range msgs {
			if cipherKey != "" {
				v, ok := val.(string)
				var err error

				if ok {
					val, err = unmarshalWithDecrypt(v, cipherKey)
					if err != nil {
						e := pnerr.NewResponseParsingError("Error unmarshalling response",
							ioutil.NopCloser(bytes.NewBufferString(v)), err)

						return emptyHistoryResp, e
					}
				}

				msgs[k] = val
			}

			if _, ok := val.(string); ok {
				items[k].Message = val

				items = append(items, items[k])
			}

			if _, ok := val.(float64); ok {
				items[k].Message = val

				items = append(items, items[k])
			}

			if m, ok := val.(map[string]interface{}); ok {
				if m["timetoken"] != nil {
					for k, value := range msgs {
						msg := value.(map[string]interface{})

						timetoken := msg["timetoken"].(float64)

						items[k].Message = msg["message"]
						items[k].Timetoken = int64(timetoken)

						items = append(items, items[k])
					}
				} else {
					for k, value := range msgs {
						items[k].Message = value

						items = append(items, items[k])
					}
				}
			}

			if v, ok := val.([]interface{}); ok {
				items[k].Message = v

				items = append(items, items[k])
			}
		}

		startTimetoken, ok := v[1].(float64)
		if !ok {
			e := pnerr.NewResponseParsingError("Error parsing response",
				ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

			return emptyHistoryResp, e
		}

		endTimetoken, ok := v[2].(float64)
		if !ok {
			e := pnerr.NewResponseParsingError("Error parsing response",
				ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

			return emptyHistoryResp, e
		}

		resp.Messages = items
		resp.StartTimetoken = int64(startTimetoken)
		resp.EndTimetoken = int64(endTimetoken)
		break
	default:
		e := pnerr.NewResponseParsingError("Error parsing response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyHistoryResp, e
	}

	return resp, nil
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
