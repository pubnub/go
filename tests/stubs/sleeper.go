package stubs

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Sleeper struct {
	Timeout int
}

// timeout - timeout in milliseconds to sleep
func NewSleeperClient(timeout int) *http.Client {
	return &http.Client{
		Transport: &Sleeper{
			Timeout: timeout,
		},
	}
}

func (s *Sleeper) RoundTrip(req *http.Request) (*http.Response,
	error) {
	// function helper for 1.6
	cancel := GetRequestCancelChannel(req)

	select {
	case <-time.After(time.Duration(s.Timeout) * time.Millisecond):
		body := ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
			"%d ms passed", s.Timeout)))
		return &http.Response{
			Status:           "530 RoundTrip Timeout",
			StatusCode:       530,
			Proto:            "HTTP/1.0",
			ProtoMajor:       1,
			ProtoMinor:       0,
			Request:          req,
			TransferEncoding: nil,
			Body:             body,
			Close:            true,
			ContentLength:    0,
		}, nil
	case err := <-cancel:
		return nil, err
	}

	return nil, errors.New("sleeper unexpected case")
}
