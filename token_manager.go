package pubnub

import (
	"sync"
)

// TokenManager struct is used to for token manager operations
type TokenManager struct {
	sync.RWMutex
	Token string
}

func newTokenManager(pubnub *PubNub, ctx Context) *TokenManager {
	return &TokenManager{}
}

// CleanUp resets the token manager
func (m *TokenManager) CleanUp() {
	m.Lock()
	m.Token = ""
	m.Unlock()
}

func (m *TokenManager) GetToken() string {
	m.RLock()
	token := m.Token
	m.RUnlock()
	return token
}

// StoreToken Aceepts PAMv3 token format token to store in the token manager
func (m *TokenManager) StoreToken(token string) {
	m.Lock()
	m.Token = token
	m.Unlock()
}
