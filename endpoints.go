package pubnub

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pubnub/go/pnerr"
	"github.com/pubnub/go/utils"
)

type endpointOpts interface {
	jobQueue() chan *JobQItem
	config() Config
	client() *http.Client
	context() Context
	validate() error
	buildPath() (string, error)
	buildQuery() (*url.Values, error)
	buildBody() ([]byte, error)
	httpMethod() string
	operationType() OperationType
	telemetryManager() *TelemetryManager
}

func SetQueryParam(q *url.Values, queryParam map[string]string) {
	if queryParam != nil {
		for key, value := range queryParam {
			q.Set(key, value)
		}
	}
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

func buildURL(o endpointOpts) (*url.URL, error) {
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

	if v := o.config().AuthKey; v != "" && query.Get("auth") == "" {
		query.Set("auth", v)
	}

	if o.config().SecretKey != "" {
		timestamp := time.Now().Unix()
		query.Set("timestamp", strconv.Itoa(int(timestamp)))

		if (!o.config().UsePAMV3) || ((o.operationType() == PNPublishOperation) && (o.httpMethod() == "POST")) {
			signedInput := o.config().SubscribeKey + "\n" + o.config().PublishKey + "\n"

			signedInput += fmt.Sprintf("%s\n", path)

			signedInput += utils.PreparePamParams(query)
			o.config().Log.Println("signedInput:", signedInput)

			signature = utils.GetHmacSha256(o.config().SecretKey, signedInput)
		} else {
			signature = createSignatureV2(o, path, query)
		}
	}

	if o.operationType() == PNPublishOperation {
		v := query.Get("meta")
		if v != "" {
			query.Set("meta", utils.URLEncode(v))
		}
	}

	if o.operationType() == PNSetStateOperation {
		v := query.Get("state")
		query.Set("state", utils.URLEncode(v))
	}

	if v := query.Get("uuid"); v != "" {
		query.Set("uuid", utils.URLEncode(v))
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

	secure := ""
	if o.config().Secure {
		secure = "s"
	}

	retURL := &url.URL{
		Opaque:   path,
		Scheme:   fmt.Sprintf("http%s", secure),
		Host:     o.config().Origin,
		RawQuery: stringifiedQuery,
	}

	return retURL, nil
}

func createSignatureV2(o endpointOpts, path string, query *url.Values) string {
	bodyString := ""
	b, err := o.buildBody()
	if err == nil {
		bodyString = string(b)
	}

	sig := createSignatureV2FromStrings(
		o.httpMethod(),
		o.config().PublishKey,
		o.config().SecretKey,
		fmt.Sprintf("%s", path),
		utils.PreparePamParams(query),
		bodyString,
		o.config().Log,
	)

	o.config().Log.Println("signaturev2:", sig)
	return sig
}

func createSignatureV2FromStrings(httpMethod, pubKey, secKey, path, query, body string, l *log.Logger) string {
	signedInputV2 := httpMethod + "\n"
	signedInputV2 += pubKey + "\n"
	signedInputV2 += path + "\n"
	signedInputV2 += query + "\n"
	signedInputV2 += body
	if l != nil {
		l.Println("signedInputV2:", signedInputV2)
	}

	encoded := utils.GetHmacSha256(secKey, signedInputV2)
	encoded = strings.TrimRight(encoded, "=")
	signatureV2 := "v2." + encoded
	return signatureV2
}

func newValidationError(o endpointOpts, msg string) error {
	return pnerr.NewValidationError(string(o.operationType()), msg)
}
