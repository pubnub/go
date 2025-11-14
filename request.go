package pubnub

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/pubnub/go/v8/pnerr"
)

// httpHeaderToMap converts http.Header to map[string]string for logging
func httpHeaderToMap(httpHeaders http.Header) map[string]string {
	headers := make(map[string]string)
	for k, v := range httpHeaders {
		if len(v) > 0 {
			headers[k] = v[0] // Take first value
		}
	}
	return headers
}

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

func fillJobQ(req *http.Request, client *http.Client, opts endpoint, j chan *JobQResponse) {
	jqi := &JobQItem{
		Req:         req,
		Client:      client,
		JobResponse: j,
	}
	opts.jobQueue() <- jqi
}

func addToJobQ(req *http.Request, client *http.Client, opts endpoint, j chan *JobQResponse, ctx Context) {
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

func buildBody(opts endpoint, url *url.URL) ([]byte, io.Reader, error) {

	b, err := opts.buildBody()
	if err != nil {
		opts.getPubNub().loggerManager.LogError(err, "BuildBodyFailed", opts.operationType(), true)
		return nil, nil, err
	}

	// Return both the bytes (for logging) and a reader (for the request)
	return b, bytes.NewReader(b), nil
}

func executeRequest(opts endpoint) ([]byte, StatusResponse, error) {
	var err error

	err = opts.validate()

	if err != nil {
		opts.getPubNub().loggerManager.LogError(err, "ValidationFailed", opts.operationType(), true)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	url, err := buildURL(opts)

	if err != nil {
		opts.getPubNub().loggerManager.LogError(err, "BuildURLFailed", opts.operationType(), true)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	opts.getPubNub().loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Preparing request: method=%s, url=%s", opts.httpMethod(), url.String()), false)

	var req *http.Request
	var requestBodyBytes []byte

	if opts.httpMethod() == "POST" {
		var body io.Reader
		requestBodyBytes, body, err = buildBody(opts, url)
		if err != nil {
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}
		req, err = newRequest("POST", url, body, opts.config().UseHTTP2)
		req.Header.Set("Content-Type", "application/json")
	} else if opts.httpMethod() == "POSTFORM" {

		body, w, _, err := opts.buildBodyMultipartFileUpload()
		if err != nil {
			opts.getPubNub().loggerManager.LogError(err, "BuildMultipartBodyFailed", opts.operationType(), true)
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}

		req, err = newRequestForMultipartWriter("POST", url.RequestURI(), &body, w, opts.config().UseHTTP2)
		if err != nil {
			opts.getPubNub().loggerManager.LogError(err, "CreateMultipartRequestFailed", opts.operationType(), true)
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}

		req.Header.Set("Content-Type", w.FormDataContentType())
	} else if opts.httpMethod() == "DELETE" {
		req, err = newRequest("DELETE", url, nil, opts.config().UseHTTP2)
	} else if opts.httpMethod() == "PATCH" {
		var body io.Reader
		requestBodyBytes, body, err = buildBody(opts, url)
		if err != nil {
			return nil, createStatus(PNUnknownCategory, "", ResponseInfo{}, err), err
		}

		req, err = newRequest("PATCH", url, body, opts.config().UseHTTP2)
	} else {
		req, err = newRequest("GET", url, nil, opts.config().UseHTTP2)
	}

	if err != nil {
		opts.getPubNub().loggerManager.LogError(err, "CreateRequestFailed", opts.operationType(), true)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}

	// Apply custom headers from endpoint
	headers, err := opts.buildHeaders()
	if err != nil {
		opts.getPubNub().loggerManager.LogError(err, "BuildHeadersFailed", opts.operationType(), true)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, err),
			err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if len(headers) > 0 {
		opts.getPubNub().loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Custom headers applied: %v", headers), false)
	}

	ctx := opts.context()
	if ctx != nil {
		// with !go1.7 you can't assign context directly to a request,
		// the request.cancel is mapped to the ctx.Done() channel instead
		// go1.7 can assign context to an executed request
		req = setRequestContext(req, ctx)
	}

	client := opts.client()

	// Log the outgoing network request
	requestHeaders := httpHeaderToMap(req.Header)
	requestBody := ""
	if opts.httpMethod() == "POST" || opts.httpMethod() == "PATCH" {
		if opts.httpMethod() == "POSTFORM" {
			requestBody = "[Multipart form data]"
		} else if len(requestBodyBytes) > 0 {
			// Truncate body for logging readability
			bodyStr := string(requestBodyBytes)
			if len(bodyStr) > 1000 {
				requestBody = bodyStr[:1000] + fmt.Sprintf("... (truncated, total: %d bytes)", len(requestBodyBytes))
			} else {
				requestBody = bodyStr
			}
		}
	}
	opts.getPubNub().loggerManager.LogNetworkRequest(PNLogLevelDebug, req.Method, req.URL.String(), requestHeaders, requestBody, true)

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
		opts.getPubNub().loggerManager.LogSimple(PNLogLevelError, fmt.Sprintf("HTTP request failed: %s", err.Error()), false)
		e := pnerr.NewConnectionError("Failed to execute request", err)

		opts.getPubNub().loggerManager.LogError(e, "NetworkRequestFailed", opts.operationType(), true)
		return nil,
			createStatus(PNUnknownCategory, "", ResponseInfo{}, e),
			e
	}

	val, status, err := parseResponse(res, opts)
	// Already wrapped error
	if err != nil {
		// Error logging already handled in parseResponse
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

	// Log successful response
	responseBody := string(val)
	if opts.httpMethod() == "POSTFORM" {
		responseBody = "[Multipart response - omitted]"
	} else if len(responseBody) > 1000 {
		responseBody = responseBody[:1000] + fmt.Sprintf("... (truncated, total: %d bytes)", len(val))
	}
	opts.getPubNub().loggerManager.LogNetworkResponse(PNLogLevelDebug, res.StatusCode, req.URL.String(), responseBody, true)

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
		rc = io.NopCloser(body)
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

func parseResponse(resp *http.Response, opts endpoint) ([]byte, StatusResponse, error) {
	status := StatusResponse{}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if (resp.StatusCode != 200) && (resp.StatusCode != 204) {
		// Errors like 400, 403, 500
		// Read body once for both error creation and logging
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)

		// Create error with the body bytes
		e := pnerr.NewServerError(resp.StatusCode, io.NopCloser(bytes.NewReader(bodyBytes)))

		// Log error response (truncate for readability)
		logLevel := PNLogLevelError
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			logLevel = PNLogLevelWarn // Client errors
		}
		logBodyStr := bodyStr
		if len(logBodyStr) > 1000 {
			logBodyStr = logBodyStr[:1000] + fmt.Sprintf("... (truncated, total: %d bytes)", len(bodyBytes))
		}
		opts.getPubNub().loggerManager.LogNetworkResponse(logLevel, resp.StatusCode, resp.Request.URL.String(), logBodyStr, true)

		if resp.StatusCode == 408 {
			opts.getPubNub().loggerManager.LogError(e, "RequestTimeout", opts.operationType(), true)
			status = createStatus(PNTimeoutCategory, "", ResponseInfo{StatusCode: resp.StatusCode}, e)
			return nil, status, e
		}

		if resp.StatusCode == 400 {
			opts.getPubNub().loggerManager.LogError(e, "BadRequest", opts.operationType(), true)
			status = createStatus(PNBadRequestCategory, "", ResponseInfo{StatusCode: resp.StatusCode}, e)
			return nil, status, e
		}

		if resp.StatusCode == 412 {
			opts.getPubNub().loggerManager.LogError(e, "PreconditionFailed", opts.operationType(), true)
			status = createStatus(PNPreconditionFailedCategory, "", ResponseInfo{StatusCode: resp.StatusCode, Operation: opts.operationType()}, e)
			return nil, status, e
		}

		opts.getPubNub().loggerManager.LogError(e, fmt.Sprintf("HTTPError%d", resp.StatusCode), opts.operationType(), true)
		status = createStatus(PNUnknownCategory, "", ResponseInfo{StatusCode: resp.StatusCode, Operation: opts.operationType()}, e)
		return nil, status, e
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error reading response body", resp.Body, err)
		opts.getPubNub().loggerManager.LogError(e, "ReadResponseBodyFailed", opts.operationType(), true)
		return nil, status, e
	}

	opts.getPubNub().loggerManager.LogSimple(PNLogLevelDebug, fmt.Sprintf("Response parsed successfully: %d bytes", len(body)), false)

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
