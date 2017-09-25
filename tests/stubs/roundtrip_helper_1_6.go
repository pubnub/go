// +build !go1.7

package stubs

import (
	"errors"
	"net/http"
)

func GetRequestCancelChannel(req *http.Request) chan error {
	cancel := make(chan error)

	go func() {
		select {
		case <-req.Cancel:
			cancel <- errors.New("request canceled")
			return
		}
	}()

	return cancel
}
