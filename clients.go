package pubnub

import (
	"net"
	"net/http"
	"time"
)

func NewHttpClient(connectTimeout int, nonSubscribeTimeout int) *http.Client {
	transport := &http.Transport{
		// MaxIdleConns: 30,
		Dial: (&net.Dialer{
			Timeout: time.Duration(connectTimeout) * time.Second,
		}).Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(nonSubscribeTimeout) * time.Second,
	}

	return client
}
