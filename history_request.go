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

	return newHistoryResponse(rawJson, opts.config().CipherKey)
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

func (o *HistoryOpts) operationType() PNOperationType {
	return PNHistoryOperation
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
				var err error

				switch v := val.(type) {
				case string:
					val, err = unmarshalWithDecrypt(v, cipherKey)
					if err != nil {
						e := pnerr.NewResponseParsingError("Error unmarshalling response",
							ioutil.NopCloser(bytes.NewBufferString(v)), err)

						return emptyHistoryResp, e
					}
					break
				case map[string]interface{}:
					msg, ok := v["pn_other"].(string)
					if !ok {
						e := pnerr.NewResponseParsingError("Decription error: ",
							ioutil.NopCloser(bytes.NewBufferString("message is empty")), nil)

						return emptyHistoryResp, e
					}
					val, err = unmarshalWithDecrypt(msg, cipherKey)
					if err != nil {
						e := pnerr.NewResponseParsingError("Error unmarshalling response",
							ioutil.NopCloser(bytes.NewBufferString(err.Error())), err)

						return emptyHistoryResp, e
					}
					break
				default:
					e := pnerr.NewResponseParsingError("Decription error: ",
						ioutil.NopCloser(bytes.NewBufferString("message is empty")), nil)

					return emptyHistoryResp, e
				}

				msgs[k] = val
			}

			switch v := val.(type) {
			case string:
				items[k].Message = val

				items = append(items, items[k])
				break
			case float64:
				items[k].Message = val

				items = append(items, items[k])
				break
			case map[string]interface{}:
				if v["timetoken"] != nil {
					for key, value := range v {
						if key == "timetoken" {
							tt := value.(float64)
							items[k].Timetoken = int64(tt)
						} else {
							items[k].Message = value
						}
					}
				} else {
					for k, value := range msgs {
						items[k].Message = value
					}
				}
				break
			case []interface{}:
				items[k].Message = v

				items = append(items, items[k])
				break
			default:
				continue
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
