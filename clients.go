package pubnub

import (
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"time"
)

// NewHTTP1Client creates a new HTTP 1 client with a new transport initialized with connect and read timeout
func NewHTTP1Client(connectTimeout, responseReadTimeout, maxIdleConnsPerHost int) *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			// Covers establishing a new TCP connection
			Timeout: time.Duration(connectTimeout) * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	}

	client := &http.Client{
		Transport: transport,
		// Covers the entire exchange from Dial to reading the body
		Timeout: time.Duration(responseReadTimeout) * time.Second,
	}

	return client
}

// NewHTTP2Client creates a new HTTP 2 client with a new transport initialized with connect and read timeout
func NewHTTP2Client(connectTimeout int, responseReadTimeout int) *http.Client {
	transport := &http2.Transport{}

	client := &http.Client{
		Transport: transport,
		// Covers the entire exchange from Dial to reading the body
		Timeout: time.Duration(responseReadTimeout) * time.Second,
	}

	return client
}
