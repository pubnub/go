package pubnub

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pubnub/go/pnerr"
)

type StatusResponse struct {
	Error error

	Category  StatusCategory
	Operation OperationType

	StatusCode int

	TlsEnabled bool

	Uuid             string
	AuthKey          string
	Origin           string
	OriginalResponse string

	AffectedChannels      []string
	AffectedChannelGroups []string
}

type ResponseInfo struct {
	Operation OperationType

	StatusCode int

	TlsEnabled bool

	Origin  string
	Uuid    string
	AuthKey string

	OriginalResponse *http.Response
}

func executeRequest(opts endpointOpts) ([]byte, StatusResponse, error) {
	err := opts.validate()
	if err != nil {
		log.Println(err)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	url, err := buildUrl(opts)

	if err != nil {
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}
	log.Println("pubnub: >>>", url)
	log.Println(opts.httpMethod())

	var req *http.Request

	if opts.httpMethod() == "POST" {
		b, err := opts.buildBody()
		if err != nil {
			return nil,
				createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
				err
		}

		body := bytes.NewReader(b)
		req, err = newRequest("POST", url, body)
	} else if opts.httpMethod() == "DELETE" {
		req, err = newRequest("DELETE", url, nil)
	} else {
		req, err = newRequest("GET", url, nil)
	}

	if err != nil {
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	ctx := opts.context()
	if ctx != nil {
		// with !go1.7 you can't assign context directly to a request,
		// the request.cancel is mapped to the ctx.Done() channel instead
		// go1.7 can assign context to an executed request
		req = setRequestContext(req, ctx)
	}

	client := opts.client()
	res, err := client.Do(req)
	// Host lookup failed
	if err != nil {
		log.Println(err.Error())
		e := pnerr.NewConnectionError("Failed to execute request", err)

		log.Println(e.Error())

		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, e),
			e
	}

	val, status, err := parseResponse(res)
	// Already wrapped error
	if err != nil {
		return nil, status, err
	}

	responseInfo := ResponseInfo{
		StatusCode:       res.StatusCode,
		OriginalResponse: res,
		Operation:        opts.operationType(),
		Origin:           url.Host,
	}

	if url.Scheme == "https" {
		responseInfo.TlsEnabled = true
	}

	if uuid, ok := url.Query()["uuid"]; ok {
		responseInfo.Uuid = uuid[0]
	}

	if auth, ok := url.Query()["auth"]; ok {
		responseInfo.AuthKey = auth[0]
	}

	status = createStatus(PNUnknownCategory, string(val), responseInfo, nil)

	return val, status, nil
}

func newRequest(method string, u *url.URL, body io.Reader) (*http.Request,
	error) {

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}

	req := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       rc,
		Host:       u.Host,
	}

	return req, nil
}

func parseResponse(resp *http.Response) ([]byte, StatusResponse, error) {
	status := StatusResponse{}

	if resp.StatusCode != 200 {
		// Errors like 400, 403, 500
		log.Println(resp.Body)

		e := pnerr.NewServerError(resp.StatusCode, resp.Body)

		log.Println(e.Error())

		if resp.StatusCode == 408 {
			status = createStatus(PNTimeoutCategory, "", ResponseInfo{}, e)

			return nil, status, e
		}

		if resp.StatusCode == 400 {
			status = createStatus(PNBadRequestCategory, "", ResponseInfo{}, e)

			return nil, status, e
		}

		status = createStatus(PNUnknownCategory, "", ResponseInfo{}, e)

		return nil, status, e
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error reading response body", resp.Body, err)
		log.Println(e)

		return nil, status, e
	}

	log.Println("pubnub: <<<", resp.Status, string(body))

	return body, status, nil
}

func createStatus(category StatusCategory, response string,
	responseInfo ResponseInfo, err error) StatusResponse {
	resp := StatusResponse{}

	if response != "" {
		resp.OriginalResponse = response
	}

	if err != nil {
		resp.Error = err
	}

	resp.StatusCode = responseInfo.StatusCode
	resp.TlsEnabled = responseInfo.TlsEnabled
	resp.Origin = responseInfo.Origin
	resp.Uuid = responseInfo.Uuid
	resp.AuthKey = responseInfo.AuthKey
	resp.Operation = responseInfo.Operation
	resp.Category = category
	resp.AffectedChannels = []string{}
	resp.AffectedChannelGroups = []string{}

	return resp
}
