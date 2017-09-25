// +build go1.7

package stubs

import "net/http"

func GetRequestCancelChannel(req *http.Request) chan error {
	cancel := make(chan error)

	go func() {
		select {
		case <-req.Context().Done():
			cancel <- req.Context().Err()
			return
		}
	}()

	return cancel
}
