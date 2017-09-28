// +build go1.7

package stubs

import "net/http"

func GetRequestCancel(req *http.Request) <-chan error {
	channel := make(chan error)

	func() {
		<-req.Context().Done()
		channel <- req.Context().Err()
	}()

	return channel
}
