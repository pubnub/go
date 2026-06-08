package stubs

import "net/http"

func GetRequestCancelChannel(req *http.Request) <-chan error {
	channel := make(chan error, 1)

	go func() {
		<-req.Context().Done()
		channel <- req.Context().Err()
	}()

	return channel
}
