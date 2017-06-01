package pubnub

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pubnub/go/pnerr"
)

func executeRequest(opts endpointOpts) (interface{}, error) {
	err := opts.validate()
	if err != nil {
		return nil, err
	}

	url := buildUrl(opts)
	config := opts.config()

	log.Println("pubnub: >>> %s", url)

	cnTimeout := config.ConnectionTimeout
	nonSubTimeout := config.NonSubscribeRequestTimeout

	// TODO: fetch client reference from options
	client := NewHttpClient(cnTimeout, nonSubTimeout)

	// TODO: can be POST
	req, err := http.NewRequest("GET", url, nil)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

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

func executeRequestWithContext(ctx context.Context,
	opts endpointOpts) (interface{}, error) {

	err := opts.validate()
	if err != nil {
		return nil, err
	}

	url := buildUrl(opts)
	config := opts.config()

	cnTimeout := config.ConnectionTimeout
	nonSubTimeout := config.NonSubscribeRequestTimeout

	client := NewHttpClient(cnTimeout, nonSubTimeout)

	// TODO: can be POST
	req, err := http.NewRequest("GET", url, nil)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req.WithContext(ctx))
	// Host lookup failed
	if err != nil {
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

func parseResponse(resp *http.Response) (interface{}, error) {
	if resp.StatusCode != 200 {
		// Errors like 400, 403, 500
		e := pnerr.NewServerError(resp.StatusCode, resp.Body)

		log.Println(e.Error())

		return nil, e
	}

	log.Println("pubnub: OK >>>", resp.Status, resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error reading response body", resp.Body, err)
		log.Println(e)

		return nil, e
	}

	var value []byte

	err = json.Unmarshal(body, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response", resp.Body, err)
		log.Println(e)

		return nil, e
	}

	return value, nil
}
