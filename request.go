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

func executeRequest(opts endpointOpts) ([]byte, error) {
	err := opts.validate()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	url, err := buildUrl(opts)
	if err != nil {
		return nil, err
	}
	log.Println("pubnub: >>>", url)
	log.Println(opts.httpMethod())

	var req *http.Request

	if opts.httpMethod() == "POST" {
		b, err := opts.buildBody()
		if err != nil {
			return nil, err
		}

		body := bytes.NewReader(b)
		req, err = newRequest("POST", url, body)
	} else {
		req, err = newRequest("GET", url, nil)
	}

	if err != nil {
		return nil, err
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

		return nil, e
	}

	val, err := parseResponse(res)
	// Already wrapped error
	if err != nil {
		return nil, err
	}

	return val, nil
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

func parseResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != 200 {
		// Errors like 400, 403, 500
		log.Println(resp.Body)

		e := pnerr.NewServerError(resp.StatusCode, resp.Body)

		log.Println(e.Error())

		return nil, e
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error reading response body", resp.Body, err)
		log.Println(e)

		return nil, e
	}

	log.Println("pubnub: <<<", resp.Status, string(body))

	return body, nil
}
