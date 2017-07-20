package pubnub

import (
	"net"
	"net/http"
	"time"
)

func NewHttpClient(connectTimeout int, responseReadTimeout int) *http.Client {
	transport := &http.Transport{
		// MaxIdleConns: 30,
		Dial: (&net.Dialer{
			// Covers establishing a new TCP connection
			Timeout: time.Duration(connectTimeout) * time.Second,
		}).Dial,
	}

	client := &http.Client{
		Transport: transport,
		// Covers the entire exchange from Dial to reading the body
		Timeout: time.Duration(responseReadTimeout) * time.Second,
	}

	return client
}
