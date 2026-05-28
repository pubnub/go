package pubnub

import (
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// closeCountingTransport implements http.RoundTripper and exposes a CloseIdleConnections
// hook so tests can verify *http.Client.CloseIdleConnections() reached the underlying transport.
// http.Client.CloseIdleConnections forwards to the transport only when the transport implements
// the unexported io.Closer-like interface { CloseIdleConnections() } — *http.Transport does, and so
// does this stub.
type closeCountingTransport struct {
	closes int64
}

func (t *closeCountingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
	}, nil
}

func (t *closeCountingTransport) CloseIdleConnections() {
	atomic.AddInt64(&t.closes, 1)
}

// TestDestroy_NoPanicWhenClientsNeverInitialised exercises the pre-existing latent defect:
// before this fix, Destroy unconditionally dereferenced pn.client, which is lazily initialised,
// so destroying an instance that never issued a request would panic with a nil-pointer deref.
func TestDestroy_NoPanicWhenClientsNeverInitialised(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.SuppressLeaveEvents = true

	pn := NewPubNub(cfg)

	assert.NotPanics(t, func() { pn.Destroy() })
}

// TestDestroy_NoPanicAfterReconnectInvalidate covers the window introduced by this PR:
// invalidateManagedHTTPClientsAfterSubscribeReconnect clears pn.client mid-life when not pinned,
// so a "create → use → reconnect → destroy" sequence must not panic on the freshly-nilled field.
func TestDestroy_NoPanicAfterReconnectInvalidate(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.UseHTTP2 = true
	cfg.SuppressLeaveEvents = true

	pn := NewPubNub(cfg)
	_ = pn.GetClient()
	_ = pn.GetSubscribeClient()

	pn.invalidateManagedHTTPClientsAfterSubscribeReconnect()

	assert.NotPanics(t, func() { pn.Destroy() })
}

// TestDestroy_ClosesBothManagedClients proves the resource-leak fix: prior code only closed
// pn.client, leaking idle connections held by pn.subscribeClient (most relevant for h2 where a
// single long-lived session per origin holds the only physical connection). UseHTTP2 is forced
// off here so the reconnect-invalidation path is bypassed and CloseIdleConnections is exercised
// exactly once per client by closeManagedHTTPClients itself.
func TestDestroy_ClosesBothManagedClients(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.UseHTTP2 = false
	cfg.SuppressLeaveEvents = true

	pn := NewPubNub(cfg)

	txn := &closeCountingTransport{}
	sub := &closeCountingTransport{}
	pn.SetClient(&http.Client{Transport: txn})
	pn.SetSubscribeClient(&http.Client{Transport: sub})

	pn.Destroy()

	assert.EqualValues(t, 1, atomic.LoadInt64(&txn.closes), "transactional client must have CloseIdleConnections called")
	assert.EqualValues(t, 1, atomic.LoadInt64(&sub.closes), "subscribe client must have CloseIdleConnections called")

	pn.Lock()
	txnPtr := pn.client
	subPtr := pn.subscribeClient
	pn.Unlock()
	assert.Nil(t, txnPtr, "pn.client must be cleared after Destroy")
	assert.Nil(t, subPtr, "pn.subscribeClient must be cleared after Destroy")
}
