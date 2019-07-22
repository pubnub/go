package pubnub

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"

// 	"github.com/pubnub/go/pnerr"
// 	"github.com/pubnub/go/utils"
// 	//"reflect"

// 	"net/http"
// 	"net/url"
// )

// var emptyObjectAPIGetUsersResp *ObjectAPIGetUsersResponse

// const objectAPIGetUsersPath = "/v1/objects/%s/users"

// type objectAPIGetUsersBuilder struct {
// 	opts *objectAPIGetUsersOpts
// }

// func newObjectAPIGetUsersBuilder(pubnub *PubNub) *objectAPIGetUsersBuilder {
// 	builder := objectAPIGetUsersBuilder{
// 		opts: &objectAPIGetUsersOpts{
// 			pubnub: pubnub,
// 		},
// 	}

// 	return &builder
// }

// func newObjectAPIGetUsersBuilderWithContext(pubnub *PubNub,
// 	context Context) *objectAPIGetUsersBuilder {
// 	builder := objectAPIGetUsersBuilder{
// 		opts: &objectAPIGetUsersOpts{
// 			pubnub: pubnub,
// 			ctx:    context,
// 		},
// 	}

// 	return &builder
// }

// // Auth sets the Authorization key with permissions to perform the request.
// func (b *objectAPIGetUsersBuilder) Auth(auth string) *objectAPIGetUsersBuilder {
// 	//b.opts.Auth = auth

// 	return b
// }

// // Auth sets the Authorization key with permissions to perform the request.
// func (b *objectAPIGetUsersBuilder) Include(include []string) *objectAPIGetUsersBuilder {
// 	b.opts.Include = include

// 	return b
// }

// // Auth sets the Authorization key with permissions to perform the request.
// func (b *objectAPIGetUsersBuilder) Limit(limit int) *objectAPIGetUsersBuilder {
// 	b.opts.Limit = limit

// 	return b
// }

// // Auth sets the Authorization key with permissions to perform the request.
// func (b *objectAPIGetUsersBuilder) Start(start string) *objectAPIGetUsersBuilder {
// 	b.opts.Start = start

// 	return b
// }

// func (b *objectAPIGetUsersBuilder) End(end string) *objectAPIGetUsersBuilder {
// 	b.opts.End = end

// 	return b
// }

// func (b *objectAPIGetUsersBuilder) Count(count bool) *objectAPIGetUsersBuilder {
// 	b.opts.Count = count

// 	return b
// }

// // QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
// func (b *objectAPIGetUsersBuilder) QueryParam(queryParam map[string]string) *objectAPIGetUsersBuilder {
// 	b.opts.QueryParam = queryParam

// 	return b
// }

// // Transport sets the Transport for the objectAPIGetUsers request.
// func (b *objectAPIGetUsersBuilder) Transport(tr http.RoundTripper) *objectAPIGetUsersBuilder {
// 	b.opts.Transport = tr
// 	return b
// }

// // Execute runs the objectAPIGetUsers request.
// func (b *objectAPIGetUsersBuilder) Execute() (*ObjectAPIGetUsersResponse, StatusResponse, error) {
// 	rawJSON, status, err := executeRequest(b.opts)
// 	if err != nil {
// 		return emptyObjectAPIGetUsersResp, status, err
// 	}

// 	return newObjectAPIGetUsersResponse(rawJSON, b.opts, status)
// }

// type objectAPIGetUsersOpts struct {
// 	pubnub *PubNub

// 	Limit      int
// 	Include    []string
// 	Start      string
// 	End        string
// 	Count      bool
// 	QueryParam map[string]string

// 	Transport http.RoundTripper

// 	ctx Context
// }

// func (o *objectAPIGetUsersOpts) config() Config {
// 	return *o.pubnub.Config
// }

// func (o *objectAPIGetUsersOpts) client() *http.Client {
// 	return o.pubnub.GetClient()
// }

// func (o *objectAPIGetUsersOpts) context() Context {
// 	return o.ctx
// }

// func (o *objectAPIGetUsersOpts) validate() error {
// 	if o.config().SubscribeKey == "" {
// 		return newValidationError(o, StrMissingSubKey)
// 	}

// 	return nil
// }

// func (o *objectAPIGetUsersOpts) buildPath() (string, error) {
// 	return fmt.Sprintf(objectAPIGetUsersPath,
// 		o.pubnub.Config.SubscribeKey), nil
// }

// func (o *objectAPIGetUsersOpts) buildQuery() (*url.Values, error) {

// 	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

// 	if o.Include != nil {
// 		q.Set("include", string(utils.JoinChannels(o.Include)))
// 	}

// 	if o.Auth != "" {
// 		q.Set("auth", o.Auth)
// 	}

// 	if o.Limit != "" {
// 		q.Set("limit", o.Limit)
// 	}

// 	if o.Start != "" {
// 		q.Set("start", o.Start)
// 	}

// 	if o.Count != "" {
// 		q.Set("count", o.Count)
// 	}

// 	if o.End != "" {
// 		q.Set("end", o.End)
// 	}

// 	SetQueryParam(q, o.QueryParam)

// 	return q, nil
// }

// func (o *objectAPIGetUsersOpts) jobQueue() chan *JobQItem {
// 	return o.pubnub.jobQueue
// }

// func (o *objectAPIGetUsersOpts) buildBody() ([]byte, error) {
// 	return []byte{}, nil
// }

// func (o *objectAPIGetUsersOpts) httpMethod() string {
// 	return "GET"
// }

// func (o *objectAPIGetUsersOpts) isAuthRequired() bool {
// 	return true
// }

// func (o *objectAPIGetUsersOpts) requestTimeout() int {
// 	return o.pubnub.Config.NonSubscribeRequestTimeout
// }

// func (o *objectAPIGetUsersOpts) connectTimeout() int {
// 	return o.pubnub.Config.ConnectTimeout
// }

// func (o *objectAPIGetUsersOpts) operationType() OperationType {
// 	return PNGetUsersOperation
// }

// func (o *objectAPIGetUsersOpts) telemetryManager() *TelemetryManager {
// 	return o.pubnub.telemetryManager
// }

// // ObjectAPIGetUsersResponse is the response to objectAPIGetUsers request. It contains a map of type objectAPIGetUsersResponseItem
// type ObjectAPIGetUsersResponse struct {
// 	Id         string                 `json:"id"`
// 	Name       string                 `json:"name"`
// 	ExternalId string                 `json:"externalId"`
// 	ProfileUrl string                 `json:"profileUrl"`
// 	Email      string                 `json:"email"`
// 	Custom     map[string]interface{} `json:"custom"`
// 	Created    string                 `json:"created"`
// 	Updated    string                 `json:"updated"`
// 	ETag       string                 `json:"eTag"`
// }

// type ObjectAPIGetUsersWithData struct {
// 	Status string                      `json:"status"`
// 	Data   []ObjectAPIGetUsersResponse `json:"data"`
// }

// func newObjectAPIGetUsersResponse(jsonBytes []byte, o *objectAPIGetUsersOpts,
// 	status StatusResponse) (*ObjectAPIGetUsersResponse, StatusResponse, error) {

// 	resp := &ObjectAPIGetUsersWithData{}

// 	fmt.Println(string(jsonBytes))

// 	err := json.Unmarshal(jsonBytes, &resp)
// 	if err != nil {
// 		e := pnerr.NewResponseParsingError("Error unmarshalling response",
// 			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

// 		return emptyObjectAPIGetUsersResp, status, e
// 	}

// 	return &resp.Data, status, nil
// }
