package pubnub

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// defaultHTTP2ClientMaxIdleConnsPerHost matches NewConfigWithUserId's Default MaxIdleConnsPerHost
// used when callers use NewHTTP2Client without PubNub's Config.
const defaultHTTP2ClientMaxIdleConnsPerHost = 30

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

// NewHTTP2Client returns an [*http.Client] whose transport prefers HTTP/2 on TLS via ALPN
// and falls back to HTTP/1.1 when the origin does not support HTTP/2. MaxIdleConnsPerHost defaults
// to match NewConfigWithUserId unless the client is built through PubNub's GetClient/GetSubscribeClient,
// which use Config.MaxIdleConnsPerHost.
func NewHTTP2Client(connectTimeout int, responseReadTimeout int) *http.Client {
	return newHTTP2Client(connectTimeout, responseReadTimeout, defaultHTTP2ClientMaxIdleConnsPerHost)
}

func newHTTP2Client(connectTimeout, responseReadTimeout, maxIdleConnsPerHost int) *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(connectTimeout) * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		// Match net/http.DefaultTransport so TLS connections negotiate HTTP/2 via ALPN
		// when the peer supports it.
		ForceAttemptHTTP2: true,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		// Should not occur on a freshly built Transport; degrade to HTTP/1-only.
		return NewHTTP1Client(connectTimeout, responseReadTimeout, maxIdleConnsPerHost)
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(responseReadTimeout) * time.Second,
	}
}

// qosMinRequestTimeout is the shortest lifetime a QoS publish or signal may use.
// The server can hold the response while it propagates the message.
const qosMinRequestTimeout = 30 * time.Second

// clientForQoS returns the HTTP client for one QoS publish or signal.
//
// The base client is never modified. Timeout 0 means no client-level deadline
// and is kept. A positive timeout already at or above qosMinRequestTimeout is
// kept, including the same client pointer. A shorter positive timeout is raised
// on a new client that shares the base transport, redirect policy, and cookie jar,
// so proxy and HTTP protocol settings stay in effect.
func clientForQoS(base *http.Client) *http.Client {
	if base == nil {
		return &http.Client{Timeout: qosMinRequestTimeout}
	}
	if base.Timeout <= 0 || base.Timeout >= qosMinRequestTimeout {
		return base
	}
	return &http.Client{
		Transport:     base.Transport,
		CheckRedirect: base.CheckRedirect,
		Jar:           base.Jar,
		Timeout:       qosMinRequestTimeout,
	}
}
