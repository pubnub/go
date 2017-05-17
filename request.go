package pubnub

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func executeRequest(ctx context.Context, e Endpoint,
	eOkCh chan interface{}, eErrCh chan error) (interface{}, error) {

	okCh := make(chan *http.Response)
	errCh := make(chan error)

	eCtx, _ := context.WithCancel(context.Background())

	go func() {
		cnTimeout := e.PubNub().Config.ConnectionTimeout
		nonSubTimeout := e.PubNub().Config.NonSubscribeRequestTimeout
		url := buildUrl(e)

		client := NewHttpClient(cnTimeout, nonSubTimeout)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("pubnub: request: Error response", err)
			eErrCh <- err
			return
		}

		res, err := client.Do(req)
		if err != nil {
			log.Println("pubnub: request: Error response", err)
			errCh <- err
		}

		okCh <- res
	}()

	select {
	case res := <-okCh:
		resp, err := parseResponse(res)
		if resp != nil {
			if eOkCh != nil {
				eOkCh <- resp
			}
			return resp, nil
		} else {
			if eErrCh != nil {
				eErrCh <- err
			}
			return nil, err
		}
	case er := <-errCh:
		if eErrCh != nil {
			eErrCh <- er
		}
		return nil, er

	case <-eCtx.Done():
		log.Println("pubnub: request: context done")
		return nil, nil
	}
}

func parseResponse(resp *http.Response) (interface{}, error) {
	if resp.StatusCode == 200 {
		log.Println("pubnub: OK >>>", resp.Status, resp.Body)
		return resp, nil
	} else {
		myerr := errors.New(fmt.Sprintf("Response error: %s", resp.Status))
		log.Println("pubnub: ERROR >>>", resp.Status, resp.Body)
		return nil, myerr
	}
}
