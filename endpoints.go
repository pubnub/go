package pubnub

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

type endpointOpts interface {
	config() Config
	client() *http.Client
	context() Context
	validate() error

	buildPath() (string, error)
	buildQuery() (*url.Values, error)
	// or bytes[]?
	buildBody() ([]byte, error)

	httpMethod() string
	operationType() OperationType
	telemetryManager() *TelemetryManager
}

func defaultQuery(uuid string, telemetryManager *TelemetryManager) *url.Values {
	v := &url.Values{}

	v.Set("pnsdk", "PubNub-Go/"+Version)
	v.Set("uuid", uuid)

	for queryName, queryParam := range telemetryManager.OperationLatency() {
		v.Set(queryName, queryParam)
	}

	return v
}

func buildUrl(o endpointOpts) (*url.URL, error) {
	var stringifiedQuery string
	var signature string

	path, err := o.buildPath()
	if err != nil {
		return &url.URL{}, err
	}

	query, err := o.buildQuery()
	if err != nil {
		return &url.URL{}, err
	}

	if o.config().FilterExpression != "" {
		query.Set("filter-expr", o.config().FilterExpression)
	}

	//if v := query.Get("auth"); v != "" {
	if v := o.config().AuthKey; v != "" {
		query.Set("auth", v)
	}

	if o.config().SecretKey != "" {
		timestamp := time.Now().Unix()
		query.Set("timestamp", strconv.Itoa(int(timestamp)))

		signedInput := o.config().SubscribeKey + "\n" + o.config().PublishKey + "\n"

		if o.operationType() == PNAccessManagerGrant ||
			o.operationType() == PNAccessManagerRevoke {
			signedInput += "grant\n"
		} else {
			signedInput += fmt.Sprintf("%s\n", path)
		}

		signedInput += utils.PreparePamParams(query)

		signature = utils.GetHmacSha256(o.config().SecretKey, signedInput)
	}

	if o.operationType() == PNPublishOperation {
		v := query.Get("meta")
		if v != "" {
			query.Set("meta", utils.UrlEncode(v))
		}
	}

	if o.operationType() == PNSetStateOperation {
		v := query.Get("state")
		query.Set("state", utils.UrlEncode(v))
	}

	if v := query.Get("uuid"); v != "" {
		query.Set("uuid", utils.UrlEncode(v))
	}

	i := 0
	for k, v := range *query {
		if i == len(*query)-1 {
			stringifiedQuery += fmt.Sprintf("%s=%s", k, v[0])
		} else {
			stringifiedQuery += fmt.Sprintf("%s=%s&", k, v[0])
		}

		i++
	}

	if signature != "" {
		stringifiedQuery += fmt.Sprintf("&signature=%s", signature)
	}

	path = fmt.Sprintf("//%s%s", o.config().Origin, path)

	retUrl := &url.URL{
		Opaque:   path,
		Scheme:   "https",
		Host:     o.config().Origin,
		RawQuery: stringifiedQuery,
	}

	return retUrl, nil
}

func newValidationError(o endpointOpts, msg string) error {
	return pnerr.NewValidationError(string(o.operationType()), msg)
}
