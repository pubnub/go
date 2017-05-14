package pubnub

import (
	"context"
	"net/http"
	"time"
)

func executeRequest(ctx context.Context, e Endpoint,
	eOk chan interface{}, eErr chan error) (interface{}, error) {

	ok := make(chan interface{})
	err := make(chan error)

	eCtx, _ := context.WithTimeout(ctx,
		time.Duration(e.PubNub().PNConfig.NonSubscribeRequestTimeout)*time.Second)

	go func() {
		url := buildUrl(e)
		req, er := http.NewRequest("GET", url, nil)
		// TODO: seems here should be a non default client
		resp, er := http.DefaultClient.Do(req.WithContext(eCtx))
		if err != nil {
			err <- er
		} else {
			ok <- resp
		}
	}()

	select {
	case resp := <-ok:
		if eOk != nil {
			eOk <- resp
		}

		return resp, nil
	case er := <-err:
		if eErr != nil {
			eErr <- er
		}
		return nil, er
	case <-eCtx.Done():
		return nil, nil
	}
}
