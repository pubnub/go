package pubnub

import (
	"net/url"
	"sync"
)

// TokenManager struct is used to for token manager operations
type TokenManager struct {
	sync.RWMutex
	Tokens              GrantResourcesWithPermissions
	pubnub              *PubNub
	lastUpdateTimestamp int
	lastQueryTimestamp  int
}

func newTokenManager(pubnub *PubNub, ctx Context) *TokenManager {

	g := GrantResourcesWithPermissions{
		Channels:        make(map[string]ChannelPermissionsWithToken),
		Groups:          make(map[string]GroupPermissionsWithToken),
		Users:           make(map[string]UserSpacePermissionsWithToken),
		Spaces:          make(map[string]UserSpacePermissionsWithToken),
		ChannelsPattern: make(map[string]ChannelPermissionsWithToken),
		GroupsPattern:   make(map[string]GroupPermissionsWithToken),
		UsersPattern:    make(map[string]UserSpacePermissionsWithToken),
		SpacesPattern:   make(map[string]UserSpacePermissionsWithToken),
	}

	manager := &TokenManager{
		Tokens: g,
	}
	manager.pubnub = pubnub

	return manager
}

// CleanUp resets the token manager
func (m *TokenManager) CleanUp() {
	m.Tokens = GrantResourcesWithPermissions{}
}

// SetAuthParan sets the auth param in the requests by retrieving the corresponding tokens from the token manager
func (m *TokenManager) SetAuthParan(q *url.Values, resourceID string, resourceType PNResourceType) {
	authParam := "auth"
	m.RLock()
	token := m.GetToken(resourceID, resourceType)
	if token != "" {
		switch resourceType {
		case PNChannels:
			q.Set(authParam, token)
		case PNGroups:
			q.Set(authParam, token)
		case PNUsers:
			q.Set(authParam, token)
		case PNSpaces:
			q.Set(authParam, token)
		}
	}
	m.RUnlock()
}

// GetAllTokens retrieves all the tokens from the token manager
func (m *TokenManager) GetAllTokens() GrantResourcesWithPermissions {
	m.RLock()
	t := m.Tokens
	m.RUnlock()
	return t
}

// GetTokensByResource retrieves the tokens by PNResourceType from the token manager
func (m *TokenManager) GetTokensByResource(resourceType PNResourceType) GrantResourcesWithPermissions {
	g := GrantResourcesWithPermissions{
		Channels:        make(map[string]ChannelPermissionsWithToken),
		Groups:          make(map[string]GroupPermissionsWithToken),
		Users:           make(map[string]UserSpacePermissionsWithToken),
		Spaces:          make(map[string]UserSpacePermissionsWithToken),
		ChannelsPattern: make(map[string]ChannelPermissionsWithToken),
		GroupsPattern:   make(map[string]GroupPermissionsWithToken),
		UsersPattern:    make(map[string]UserSpacePermissionsWithToken),
		SpacesPattern:   make(map[string]UserSpacePermissionsWithToken),
	}
	m.RLock()
	switch resourceType {
	case PNChannels:
		for k, v := range m.Tokens.Channels {
			g.Channels[k] = v
		}

		for k, v := range m.Tokens.ChannelsPattern {
			g.ChannelsPattern[k] = v
		}
	case PNGroups:
		for k, v := range m.Tokens.Groups {
			g.Groups[k] = v
		}

		for k, v := range m.Tokens.GroupsPattern {
			g.GroupsPattern[k] = v
		}
	case PNUsers:
		for k, v := range m.Tokens.Users {
			g.Users[k] = v
		}

		for k, v := range m.Tokens.UsersPattern {
			g.UsersPattern[k] = v
		}
	case PNSpaces:
		for k, v := range m.Tokens.Spaces {
			g.Spaces[k] = v
		}

		for k, v := range m.Tokens.SpacesPattern {
			g.SpacesPattern[k] = v
		}
	}
	m.RUnlock()
	return g
}

// GetToken first match for direct ids, if no match found use the first token from pattern match ignoring the regex (by design).
func (m *TokenManager) GetToken(resourceID string, resourceType PNResourceType) string {
	m.RLock()
	switch resourceType {
	case PNChannels:
		if d, ok := m.Tokens.Channels[resourceID]; ok {
			return d.Token
		}

		for _, v := range m.Tokens.ChannelsPattern {
			return v.Token
		}
	case PNGroups:
		if d, ok := m.Tokens.Groups[resourceID]; ok {
			return d.Token
		}

		for _, v := range m.Tokens.GroupsPattern {
			return v.Token
		}
	case PNUsers:
		if d, ok := m.Tokens.Users[resourceID]; ok {
			return d.Token
		}

		for _, v := range m.Tokens.UsersPattern {
			return v.Token
		}
	case PNSpaces:
		if d, ok := m.Tokens.Spaces[resourceID]; ok {
			return d.Token
		}

		for _, v := range m.Tokens.SpacesPattern {
			return v.Token
		}
	}
	m.RUnlock()
	return ""
}

func mergeTokensByResource(m interface{}, resource interface{}, resourceType PNResourceType) {
	switch resourceType {
	case PNChannels:
		c := resource.(map[string]ChannelPermissionsWithToken)
		d := m.(map[string]ChannelPermissionsWithToken)
		for k, v := range c {
			d[k] = v
		}
	case PNGroups:
		c := resource.(map[string]GroupPermissionsWithToken)
		d := m.(map[string]GroupPermissionsWithToken)
		for k, v := range c {
			d[k] = v
		}
	default:
		//case PNUsers:
		//case PNSpaces:
		c := resource.(map[string]UserSpacePermissionsWithToken)
		d := m.(map[string]UserSpacePermissionsWithToken)
		for k, v := range c {
			d[k] = v
		}
	}
}

// StoreTokens Aceepts PAMv3 token format tokens to store in the token manager
func (m *TokenManager) StoreTokens(token []string) {
	for _, k := range token {
		m.StoreToken(k)
	}
}

// StoreToken Aceepts PAMv3 token format token to store in the token manager
func (m *TokenManager) StoreToken(token string) {

	if m.pubnub.Config.StoreTokensOnGrant && m.pubnub.Config.SecretKey == "" {
		m.pubnub.Config.Log.Println("token: ", token)
		cborObject, err := GetPermissions(token)
		if err == nil {

			res := ParseGrantResources(cborObject.Resources, token, cborObject.Timestamp, cborObject.TTL)
			m.Lock()
			mergeTokensByResource(m.Tokens.Channels, res.Channels, PNChannels)
			mergeTokensByResource(m.Tokens.Users, res.Users, PNUsers)
			mergeTokensByResource(m.Tokens.Groups, res.Groups, PNGroups)
			mergeTokensByResource(m.Tokens.Spaces, res.Spaces, PNSpaces)

			//clear all Users/Spaces pattern maps (by design, store last token only for patterns)
			pat := ParseGrantResources(cborObject.Patterns, token, cborObject.Timestamp, cborObject.TTL)
			if len(pat.Users) > 0 {
				m.Tokens.UsersPattern = make(map[string]UserSpacePermissionsWithToken)
				m.pubnub.Config.Log.Println("Clearing UsersPattern from Token Manager")
			}
			if len(pat.Spaces) > 0 {
				m.Tokens.SpacesPattern = make(map[string]UserSpacePermissionsWithToken)
				m.pubnub.Config.Log.Println("Clearing SpacesPattern from Token Manager")
			}

			mergeTokensByResource(m.Tokens.ChannelsPattern, pat.Channels, PNChannels)
			mergeTokensByResource(m.Tokens.UsersPattern, pat.Users, PNUsers)
			mergeTokensByResource(m.Tokens.GroupsPattern, pat.Groups, PNGroups)
			mergeTokensByResource(m.Tokens.SpacesPattern, pat.Spaces, PNSpaces)

			m.pubnub.Config.Log.Println("Tokens: ", m.Tokens)

			m.Unlock()
		} else {
			m.pubnub.Config.Log.Println("Not storing tokens as StoreTokensOnGrant is false and SecretKey is set ")
		}
	} else {

	}
}
