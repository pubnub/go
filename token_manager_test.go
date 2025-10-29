package pubnub

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenManager(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	tm := newTokenManager(pn, pn.ctx)

	assert.NotNil(tm)
	assert.Empty(tm.Token)
}

func TestTokenManagerStoreToken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	testToken := "test-token-123"
	tm.StoreToken(testToken)

	assert.Equal(testToken, tm.Token)
}

func TestTokenManagerGetToken(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	testToken := "test-token-456"
	tm.StoreToken(testToken)

	retrievedToken := tm.GetToken()
	assert.Equal(testToken, retrievedToken)
}

func TestTokenManagerCleanUp(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	// Store a token first
	testToken := "test-token-789"
	tm.StoreToken(testToken)
	assert.Equal(testToken, tm.GetToken())

	// Clean up
	tm.CleanUp()
	assert.Empty(tm.GetToken())
}

func TestTokenManagerConcurrency(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	const numGoroutines = 100
	const testToken = "concurrent-test-token"

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Test concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			tm.StoreToken(testToken)
		}(i)
	}

	wg.Wait()

	// Verify final state
	assert.Equal(testToken, tm.GetToken())
}

func TestTokenManagerConcurrentReadWrite(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	const numReaders = 50
	const numWriters = 50
	const testToken = "read-write-test-token"

	var wg sync.WaitGroup

	// Start concurrent readers
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				token := tm.GetToken()
				// Token should be either empty or valid token
				assert.True(token == "" || token == testToken)
			}
		}()
	}

	// Start concurrent writers
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func() {
			defer wg.Done()
			tm.StoreToken(testToken)
		}()
	}

	wg.Wait()

	// Final token should be the test token
	assert.Equal(testToken, tm.GetToken())
}

func TestTokenManagerMultipleCleanUps(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	// Store token
	tm.StoreToken("token1")
	assert.Equal("token1", tm.GetToken())

	// Multiple cleanups should be safe
	tm.CleanUp()
	tm.CleanUp()
	tm.CleanUp()

	assert.Empty(tm.GetToken())

	// Should be able to store after cleanup
	tm.StoreToken("token2")
	assert.Equal("token2", tm.GetToken())
}

func TestTokenManagerEmptyTokenStorage(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	tm := newTokenManager(pn, pn.ctx)

	// Store empty token
	tm.StoreToken("")
	assert.Empty(tm.GetToken())

	// Store token, then empty
	tm.StoreToken("test-token")
	assert.Equal("test-token", tm.GetToken())

	tm.StoreToken("")
	assert.Empty(tm.GetToken())
}
