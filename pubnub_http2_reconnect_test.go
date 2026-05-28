package pubnub

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidateManagedHTTPClientsAfterSubscribeReconnect_DiscardManagedHTTP2Clients(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.UseHTTP2 = true

	pn := NewPubNub(cfg)
	a := pn.GetSubscribeClient()
	b := pn.GetClient()

	pn.invalidateManagedHTTPClientsAfterSubscribeReconnect()

	a2 := pn.GetSubscribeClient()
	b2 := pn.GetClient()

	assert.False(t, a == a2 && b == b2)
	assert.False(t, a == a2)
	assert.False(t, b == b2)
	assert.False(t, pn.txnHTTPClientPinned)
	assert.False(t, pn.subscribeClientPinned)
}

func TestInvalidateManagedHTTPClientsAfterSubscribeReconnect_SkipsWhenUseHTTP2Disabled(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.UseHTTP2 = false

	pn := NewPubNub(cfg)
	ref := pn.GetSubscribeClient()

	pn.invalidateManagedHTTPClientsAfterSubscribeReconnect()

	assert.True(t, ref == pn.GetSubscribeClient())
}

func TestInvalidateManagedHTTPClientsAfterSubscribeReconnect_PinnedClientsRetainPointer(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.UseHTTP2 = true

	pn := NewPubNub(cfg)
	txn := &http.Client{}
	sub := &http.Client{}
	pn.SetClient(txn)
	pn.SetSubscribeClient(sub)

	pn.invalidateManagedHTTPClientsAfterSubscribeReconnect()

	assert.True(t, txn == pn.GetClient())
	assert.True(t, sub == pn.GetSubscribeClient())
	assert.True(t, pn.txnHTTPClientPinned)
	assert.True(t, pn.subscribeClientPinned)
}
