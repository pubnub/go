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
			eErrCh <- err
			return
		}

		res, err := client.Do(req)
		if err != nil {
			log.Println("pubnub: Error response")
			errCh <- err
		}

		// TODO: Do not parse request if context is cancelled/deadline has reached
		okCh <- res
	}()

	select {
	case res := <-okCh:
		// TODO: move this logic into a separate func
		log.Println("> ok2")
		if res.StatusCode == 200 {
			// TODO: move 1st case above
			log.Println("> Status", res.Status)
			if eOkCh != nil {
				eOkCh <- res
			}
			return res, nil
		} else {
			myerr := errors.New(fmt.Sprintf("Response error: %s", res.Status))
			eErrCh <- myerr
			return nil, myerr
		}
	case er := <-errCh:
		log.Println("> err")
		if eErrCh != nil {
			eErrCh <- er
		}
		return nil, er
		// TODO: do not return nil for sync call
	case <-eCtx.Done():
		log.Println("> ctx done")
		return nil, nil
	}
}
