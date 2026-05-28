package pubnub

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// roundTripCountTransport counts RoundTrip invocations for subscribe client reuse checks.
type roundTripCountTransport struct {
	base  http.RoundTripper
	count *int64
}

func (t *roundTripCountTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt64(t.count, 1)
	return t.base.RoundTrip(req)
}

func TestGetSubscribeClientSingletonPointer(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub-key"
	cfg.UseHTTP2 = true
	pn := NewPubNub(cfg)

	c1 := pn.GetSubscribeClient()
	c2 := pn.GetSubscribeClient()
	assert.True(t, c1 == c2, "expected same *http.Client pointer from GetSubscribeClient")

	opts := newSubscribeOpts(pn, backgroundContext)
	assert.True(t, opts.client() == c1, "subscribe opts must use GetSubscribeClient")
}

func TestSubscribeMultipleExecuteRequestReusesClient(t *testing.T) {
	for _, useHTTP2 := range []bool{false, true} {
		t.Run(fmt.Sprintf("UseHTTP2_%v", useHTTP2), func(t *testing.T) {
			const body = `{"t":{"t":"16999999999999999","r":1}}`

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.URL.Path, "/v2/subscribe/") {
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(body))
			}))
			defer srv.Close()

			u, err := url.Parse(srv.URL)
			require.NoError(t, err)

			cfg := NewConfigWithUserId(UserId(GenerateUUID()))
			cfg.SubscribeKey = "test-sub-key"
			cfg.Origin = u.Host
			cfg.Secure = false
			cfg.UseHTTP2 = useHTTP2

			pn := NewPubNub(cfg)
			base := pn.GetSubscribeClient()

			var roundTrips int64
			pn.SetSubscribeClient(&http.Client{
				Transport: &roundTripCountTransport{base: base.Transport, count: &roundTrips},
				Timeout:   base.Timeout,
			})

			opts := newSubscribeOpts(pn, backgroundContext)
			opts.Channels = []string{"reuse-ch"}

			ref := pn.GetSubscribeClient()
			_, _, err = executeRequest(opts)
			require.NoError(t, err)
			_, _, err = executeRequest(opts)
			require.NoError(t, err)

			assert.EqualValues(t, 2, atomic.LoadInt64(&roundTrips))
			assert.True(t, pn.GetSubscribeClient() == ref, "subscribe client pointer must stay stable")
		})
	}
}

func TestSubscribeUsesHTTP2ProtoOverTLSWhenConfigured(t *testing.T) {
	const body = `{"t":{"t":"16999999999999999","r":1}}`

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v2/subscribe/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	srv.EnableHTTP2 = true
	srv.StartTLS()
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "test-sub-key"
	cfg.Origin = u.Host
	cfg.Secure = true
	cfg.UseHTTP2 = true

	pn := NewPubNub(cfg)
	base := pn.GetSubscribeClient()
	tr, ok := base.Transport.(*http.Transport)
	require.True(t, ok)
	tr.TLSClientConfig = trustTestServerCert(tr.TLSClientConfig, srv.Certificate())

	var negotiatedProto string
	counting := &roundTripCountTransport{
		base: &protoCaptureTransport{
			base: tr,
			dest: &negotiatedProto,
		},
		count: new(int64),
	}
	pn.SetSubscribeClient(&http.Client{
		Transport: counting,
		Timeout:   base.Timeout,
	})

	opts := newSubscribeOpts(pn, backgroundContext)
	opts.Channels = []string{"h2-ch"}

	_, _, err = executeRequest(opts)
	require.NoError(t, err)

	assert.EqualValues(t, 1, atomic.LoadInt64(counting.count))
	assert.Equal(t, "HTTP/2.0", negotiatedProto)
}
