package pubnub

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/pubnub/go/v5/pnerr"
)

// StatusResponse is used to store the usable properties in the response of an request.
type StatusResponse struct {
	Error                 error
	Category              StatusCategory
	Operation             OperationType
	StatusCode            int
	TLSEnabled            bool
	UUID                  string
	AuthKey               string
	Origin                string
	OriginalResponse      string
	Request               string
	AffectedChannels      []string
	AffectedChannelGroups []string
	AdditionalData        interface{}
}

// ResponseInfo is used to store the properties in the response of an request.
type ResponseInfo struct {
	Operation        OperationType
	StatusCode       int
	TLSEnabled       bool
	Origin           string
	UUID             string
	AuthKey          string
	OriginalResponse *http.Response
}

func fillJobQ(req *http.Request, client *http.Client, opts endpointOpts, j chan *JobQResponse) {
	jqi := &JobQItem{
		Req:         req,
		Client:      client,
		JobResponse: j,
	}
	opts.jobQueue() <- jqi
}

func addToJobQ(req *http.Request, client *http.Client, opts endpointOpts, j chan *JobQResponse, ctx Context) {
	if ctx != nil {
		select {
		case <-ctx.Done():
			return
		default:
			fillJobQ(req, client, opts, j)
		}
	} else {
		fillJobQ(req, client, opts, j)
	}
}

func buildBody(opts endpointOpts, url *url.URL) (io.Reader, error) {

	b, err := opts.buildBody()
	if err != nil {
		opts.config().Log.Println("PNUnknownCategory", err, url)
		return nil, err
	}
	opts.config().Log.Println("BODY", string(b))

	return bytes.NewReader(b), nil
}

func executeRequest(opts endpointOpts) ([]byte, StatusResponse, error) {
	err := opts.validate()

	if err != nil {
		opts.config().Log.Println("PNUnknownCategory", err)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	url, err := buildURL(opts)

	if err != nil {
		opts.config().Log.Println("PNUnknownCategory", err)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	opts.config().Log.Println(fmt.Sprintf("url:%s\nmethod:%s", url, opts.httpMethod()))

	var req *http.Request

	if opts.httpMethod() == "POST" {
		body, err := buildBody(opts, url)
		if err != nil {
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}
		req, err = newRequest("POST", url, body, opts.config().UseHTTP2)
		req.Header.Set("Content-Type", "application/json")
	} else if opts.httpMethod() == "POSTFORM" {

		body, w, _, err := opts.buildBodyMultipartFileUpload()
		if err != nil {
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}

		req, err = newRequestForMultipartWriter("POST", url.RequestURI(), &body, w, opts.config().UseHTTP2)
		if err != nil {
			opts.config().Log.Println("POST ERROR : ", err)
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}

		req.Header.Set("Content-Type", w.FormDataContentType())
	} else if opts.httpMethod() == "DELETE" {
		req, err = newRequest("DELETE", url, nil, opts.config().UseHTTP2)
	} else if opts.httpMethod() == "PATCH" {
		body, err := buildBody(opts, url)
		if err != nil {
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}

		req, err = newRequest("PATCH", url, body, opts.config().UseHTTP2)
	} else {
		req, err = newRequest("GET", url, nil, opts.config().UseHTTP2)
	}

	if err != nil {
		opts.config().Log.Println("PNUnknownCategory", err, url)
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

	startTimestamp := time.Now()

	var res *http.Response
	runRequestWorker := false

	switch opts.operationType() {
	case PNPublishOperation, PNAccessManagerGrant:
		runRequestWorker = true
	}

	if runRequestWorker && opts.config().MaxWorkers > 0 {
		j := make(chan *JobQResponse)
		go addToJobQ(req, client, opts, j, ctx)
		jr := <-j
		close(j)
		res = jr.Resp
		err = jr.Error
	} else {
		res, err = client.Do(req)
	}

	// Host lookup failed
	if err != nil {
		opts.config().Log.Println("err.Error()", err.Error())
		e := pnerr.NewConnectionError("Failed to execute request", err)

		opts.config().Log.Println("PNUnknownCategory", e.Error(), url)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, e),
			e
	}

	val, status, err := parseResponse(res, opts)
	// Already wrapped error
	if err != nil {
		opts.config().Log.Println("res.StatusCode, status, err.Error()", res.StatusCode, status, err.Error())
		return nil, status, err
	}

	elapsedTime := time.Since(startTimestamp)

	manager := opts.telemetryManager()
	manager.StoreLatency(elapsedTime.Seconds(), opts.operationType())

	responseInfo := ResponseInfo{
		StatusCode:       res.StatusCode,
		OriginalResponse: res,
		Operation:        opts.operationType(),
		Origin:           url.Host,
	}

	if url.Scheme == "https" {
		responseInfo.TLSEnabled = true
	}

	if uuid, ok := url.Query()["uuid"]; ok {
		responseInfo.UUID = uuid[0]
	}

	if auth, ok := url.Query()["auth"]; ok {
		responseInfo.AuthKey = auth[0]
	}

	if opts.httpMethod() != "POSTFORM" {
		opts.config().Log.Println("PNUnknownCategory", string(val), responseInfo)
	}
	status = createStatus(PNUnknownCategory, string(val), responseInfo, nil)

	return val, status, nil
}

func newRequestForMultipartWriter(method string, URL string, body io.Reader, writer *multipart.Writer, useHTTP2 bool) (*http.Request, error) {
	req, err := http.NewRequest(method, URL, body)
	if useHTTP2 {
		req.Proto = "HTTP/2.0"
		req.ProtoMajor = 2
		req.ProtoMinor = 0
	} else {
		req.Proto = "HTTP/1.1"
		req.ProtoMajor = 1
		req.ProtoMinor = 1
	}
	return req, err
}

func newRequest(method string, u *url.URL, body io.Reader, useHTTP2 bool) (*http.Request, error) {
	var rc io.ReadCloser
	var ok bool

	rc, ok = body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}

	if useHTTP2 {
		req := &http.Request{
			Method:     method,
			URL:        u,
			Proto:      "HTTP/2.0",
			ProtoMajor: 2,
			ProtoMinor: 0,
			Header:     make(http.Header),
			Body:       rc,
			Host:       u.Host,
		}
		return req, nil
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

func parseResponse(resp *http.Response, opts endpointOpts) ([]byte, StatusResponse, error) {
	status := StatusResponse{}

	if (resp.StatusCode != 200) && (resp.StatusCode != 204) {
		// Errors like 400, 403, 500
		e := pnerr.NewServerError(resp.StatusCode, resp.Body)

		opts.config().Log.Println(e.Error())

		if resp.StatusCode == 408 {
			opts.config().Log.Println("PNTimeoutCategory: resp.StatusCode, resp.Body, resp.Request.URL", resp.StatusCode, resp.Body, resp.Request.URL)
			status = createStatus(PNTimeoutCategory, "", ResponseInfo{StatusCode: resp.StatusCode}, e)

			return nil, status, e
		}

		if resp.StatusCode == 400 {
			opts.config().Log.Println("PNBadRequestCategory: resp.StatusCode, resp.Body, resp.Request.URL", resp.StatusCode, resp.Body, resp.Request.URL)
			status = createStatus(PNBadRequestCategory, "", ResponseInfo{StatusCode: resp.StatusCode}, e)

			return nil, status, e
		}
		opts.config().Log.Println("PNUnknownCategory: resp.StatusCode, resp.Body, resp.Request.URL", resp.StatusCode, resp.Body, resp.Request.URL)
		status = createStatus(PNUnknownCategory, "", ResponseInfo{StatusCode: resp.StatusCode, Operation: opts.operationType()}, e)

		return nil, status, e
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error reading response body", resp.Body, err)
		opts.config().Log.Println("Read All error: resp.Body, resp.Request.URL, e", resp.StatusCode, resp.Body, resp.Request.URL, e)

		return nil, status, e
	}

	opts.config().Log.Println("200 OK: resp.StatusCode, resp.Status, resp.Body, resp.Request.URL, string(body)", resp.StatusCode, resp.Status, resp.Body, resp.Request.URL, string(body))

	//opts.config().Log.Println("200 OK: resp.StatusCode, resp.Status, resp.Request.URL", resp.StatusCode, resp.Status, resp.Request.URL)
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
	resp.TLSEnabled = responseInfo.TLSEnabled
	resp.Origin = responseInfo.Origin
	resp.UUID = responseInfo.UUID
	resp.AuthKey = responseInfo.AuthKey
	resp.Operation = responseInfo.Operation
	resp.Category = category
	resp.AffectedChannels = []string{}
	resp.AffectedChannelGroups = []string{}

	return resp
}
