package pubnub

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pubnub/go/v5/pnerr"
	"github.com/pubnub/go/v5/utils"
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
	buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error)
	httpMethod() string
	operationType() OperationType
	telemetryManager() *TelemetryManager
}

// SetQueryParam appends the query params map to the query string
func SetQueryParam(q *url.Values, queryParam map[string]string) {
	if queryParam != nil {
		for key, value := range queryParam {
			q.Set(key, utils.URLEncode(value))
		}
	}
}

// SetArrayTypeQueryParam appends to the query string the key val pair
func SetArrayTypeQueryParam(q *url.Values, val []string, key string) {
	for _, value := range val {
		q.Add(key, utils.URLEncode(value))
	}
}

// SetQueryParamAsCommaSepString appends to the query string the comma separated string.
func SetQueryParamAsCommaSepString(q *url.Values, val []string, key string) {
	q.Set(key, strings.Join(val, ","))
}

// SetPushEnvironment appends the push environment to the query string
func SetPushEnvironment(q *url.Values, env PNPushEnvironment) {
	if string(env) != "" {
		q.Set("environment", string(env))
	} else {
		q.Set("environment", string(PNPushEnvironmentDevelopment))
	}
}

// SetPushTopic appends the topic to the query string
func SetPushTopic(q *url.Values, topic string) {
	if topic != "" {
		q.Set("topic", utils.URLEncode(topic))
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

	if v := query.Get("filter-expr"); v != "" {
		query.Set("filter-expr", utils.URLEncode(v))
	}

	if v := query.Get("filter"); v != "" {
		query.Set("filter", utils.URLEncode(v))
	}
	if v := query.Get("include"); v != "" {
		query.Set("include", utils.URLEncode(v))
	}
	if v := query.Get("sort"); v != "" {
		query.Set("sort", utils.URLEncode(v))
	}

	i := 0
	for k, v := range *query {
		for j, value := range v {
			if (i == len(*query)-1) && (j == len(v)-1) {
				stringifiedQuery += fmt.Sprintf("%s=%s", k, value)
			} else {
				stringifiedQuery += fmt.Sprintf("%s=%s&", k, value)
			}
			j++
		}

		i++
	}

	if signature != "" {
		stringifiedQuery += fmt.Sprintf("&signature=%s", signature)
	}

	secure := ""
	if o.config().Secure {
		secure = "s"
	}

	scheme := fmt.Sprintf("http%s", secure)

	host := o.config().Origin

	if o.httpMethod() != "POSTFORM" {
		path = fmt.Sprintf("//%s%s", o.config().Origin, path)
	} else {
		p := strings.Split(path, "://")
		scheme = p[0]
		p2 := strings.Split(p[1], "?")
		path = fmt.Sprintf("//%s", p2[0])
		h := strings.Split(p[1], "/")
		host = h[0]
		stringifiedQuery = ""
	}

	retURL := &url.URL{
		Opaque:   path,
		Scheme:   scheme,
		Host:     host,
		RawQuery: stringifiedQuery,
	}

	return retURL, nil
}

func createSignatureV2(o endpointOpts, path string, query *url.Values) string {
	bodyString := ""
	b, err := o.buildBody()
	if err == nil {
		bodyString = string(b)
	} else {
		o.config().Log.Println("buildBody error", err.Error())
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
	return pnerr.NewValidationError(o.operationType().String(), msg)
}
