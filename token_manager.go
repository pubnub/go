package pubnub

import (
	"fmt"
	"sync"
)

type TokenManager struct {
	sync.RWMutex
	Tokens GrantResourcesWithPermissions
	pubnub *PubNub
}

// Match tokens for subscribe or other calls
// Check ttl expiration

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

func (m *TokenManager) GetTokens(channels, groups, users, spaces []string) *GrantResourcesWithPermissions {
	g := GrantResourcesWithPermissions{
		Channels: make(map[string]ChannelPermissionsWithToken),
		Groups:   make(map[string]GroupPermissionsWithToken),
		Users:    make(map[string]UserSpacePermissionsWithToken),
		Spaces:   make(map[string]UserSpacePermissionsWithToken),
		// ChannelsPattern: make(map[string]ChannelPermissionsWithToken),
		// GroupsPattern:   make(map[string]GroupPermissionsWithToken),
		// UsersPattern:    make(map[string]UserSpacePermissionsWithToken),
		// SpacesPattern:   make(map[string]UserSpacePermissionsWithToken),
	}
	//findTokenInTokensChannels(channels, g.Channels, m.Tokens.Channels)
	findTokenInTokens(channels, g.Channels, m.Tokens.Channels, PNChannels)
	findTokenInTokens(groups, g.Groups, m.Tokens.Groups, PNGroups)
	findTokenInTokens(users, g.Users, m.Tokens.Users, PNUsers)
	findTokenInTokens(spaces, g.Spaces, m.Tokens.Spaces, PNSpaces)
	// findTokenInTokens(channels, g.Channels, m.Tokens.Channels, PNChannels)
	// findTokenInTokens(groups, g.Groups, m.Tokens.Groups, PNGroups)
	// findTokenInTokens(users, g.Users, m.Tokens.Users, PNUsers)
	// findTokenInTokens(spaces, g.Spaces, m.Tokens.Spaces, PNSpaces)

	return &g
}

func matchTokensForSubscribe(g *GrantResourcesWithPermissions) {

}

func findTokenInTokens(r []string, resource, merge interface{}, resourceType PNResourceType) {
	switch resourceType {
	case PNChannels:
		a := resource.(map[string]ChannelPermissionsWithToken)
		e := merge.(map[string]ChannelPermissionsWithToken)
		for _, k := range r {
			if d, ok := e[k]; ok {
				a[k] = d
			}
		}
	case PNGroups:
		a := resource.(map[string]GroupPermissionsWithToken)
		e := merge.(map[string]GroupPermissionsWithToken)
		for _, k := range r {
			if d, ok := e[k]; ok {
				a[k] = d
			}
		}
	default:
		//case PNUsers:
		//case PNSpaces:
		a := resource.(map[string]UserSpacePermissionsWithToken)
		e := merge.(map[string]UserSpacePermissionsWithToken)
		for _, k := range r {
			if d, ok := e[k]; ok {
				a[k] = d
			}
		}
	}

}

// func findTokenInTokensChannels(r []string, a, m map[string]ChannelPermissionsWithToken) {
// 	for _, k := range r {
// 		if d, ok := m[k]; ok {
// 			a[k] = d
// 		}
// 	}
// }

// func findTokenInTokensGroups(r []string, a, m map[string]GroupPermissionsWithToken) {
// 	for _, k := range r {
// 		if d, ok := m[k]; ok {
// 			a[k] = d
// 		}
// 	}
// }

// func findTokenInTokensUserSpace(r []string, a, m map[string]UserSpacePermissionsWithToken) {
// 	for _, k := range r {
// 		if d, ok := m[k]; ok {
// 			a[k] = d
// 		}
// 	}
// }

// func mergeTokensByChannels(m map[string]ChannelPermissionsWithToken, r map[string]ChannelPermissionsWithToken) {
// 	for k, v := range r {
// 		m[k] = v
// 	}

// }
// func mergeTokensByGroups(m map[string]GroupPermissionsWithToken, r map[string]GroupPermissionsWithToken) {
// 	for k, v := range r {
// 		m[k] = v
// 	}

// }
// func mergeTokensByUserSpace(m map[string]UserSpacePermissionsWithToken, r map[string]UserSpacePermissionsWithToken) {
// 	for k, v := range r {
// 		m[k] = v
// 	}

// }

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

func (m *TokenManager) StoreToken(token string) {
	if m.pubnub.Config.StoreTokensOnGrant {
		fmt.Println("--->", token)
		cborObject, err := GetPermissions(token)
		if err == nil {
			fmt.Printf("\nCBOR decode Token---> %#v", cborObject)
			fmt.Println("")
			fmt.Println("Sig: ", string(cborObject.Signature))
			fmt.Println("Version: ", cborObject.Version)
			fmt.Println("Timetoken: ", cborObject.Timestamp)
			fmt.Println("TTL: ", cborObject.TTL)
			fmt.Println(fmt.Sprintf("Meta: %#v", cborObject.Meta))
			fmt.Println("")
			fmt.Println(" --- Resources")
			res := ParseGrantResources(cborObject.Resources, token, cborObject.Timestamp)
			m.Lock()
			mergeTokensByResource(m.Tokens.Channels, res.Channels, PNChannels)
			mergeTokensByResource(m.Tokens.Users, res.Users, PNUsers)
			mergeTokensByResource(m.Tokens.Groups, res.Groups, PNGroups)
			mergeTokensByResource(m.Tokens.Spaces, res.Spaces, PNSpaces)

			// mergeTokensByChannels(m.Tokens.Channels, g.Channels)
			// mergeTokensByUserSpace(m.Tokens.Users, g.Users)
			// mergeTokensByGroups(m.Tokens.Groups, g.Groups)
			// mergeTokensByUserSpace(m.Tokens.Spaces, g.Spaces)

			fmt.Println(" --- Tokens ---- ", m.Tokens)

			fmt.Println(" --- Patterns")
			pat := ParseGrantResources(cborObject.Patterns, token, cborObject.Timestamp)
			mergeTokensByResource(m.Tokens.ChannelsPattern, pat.Channels, PNChannels)
			mergeTokensByResource(m.Tokens.UsersPattern, pat.Users, PNUsers)
			mergeTokensByResource(m.Tokens.GroupsPattern, pat.Groups, PNGroups)
			mergeTokensByResource(m.Tokens.SpacesPattern, pat.Spaces, PNSpaces)

			fmt.Println(" --- Tokens ---- ", m.Tokens)

			m.Unlock()
		}
	}
}

// func (m *TokenManager) StoreToken(token string) {
// 	if m.pubnub.Config.StoreTokensOnGrant {
// 		fmt.Println("--->", token)
// 		cborObject, err := GetPermissions(token)
// 		if err == nil {
// 			fmt.Printf("\nCBOR decode Token---> %#v", cborObject)
// 			fmt.Println("")
// 			fmt.Println("Sig: ", string(cborObject.Signature))
// 			fmt.Println("Version: ", cborObject.Version)
// 			fmt.Println("Timetoken: ", cborObject.Timetoken)
// 			fmt.Println("TTL: ", cborObject.TTL)
// 			fmt.Println(fmt.Sprintf("Meta: %#v", cborObject.Meta))
// 			fmt.Println("")
// 			fmt.Println(" --- Resources")
// 			g := ParseGrantResources(cborObject.Resources, token)
// 			m.Lock()
// 			m.Tokens.Channels = mergeTokensByResources(m.Tokens.Channels, g.Channels)
// 			m.Tokens.Users = mergeTokensByResources(m.Tokens.Users, g.Users)
// 			m.Tokens.Groups = mergeTokensByResources(m.Tokens.Groups, g.Groups)
// 			m.Tokens.Spaces = mergeTokensByResources(m.Tokens.Spaces, g.Spaces)
// 			m.Unlock()

// 			fmt.Println(" --- Tokens ---- ", m.Tokens)

// 			fmt.Println(" --- Patterns")
// 			ParseGrantResources(cborObject.Patterns, token)
// 		}
// 	}
// }

// func mergeTokensByResources(m map[string]PermissionsWithToken, r map[string]PermissionsWithToken) map[string]PermissionsWithToken {
// 	for k, v := range r {
// 		m[k] = v
// 	}
// 	return m
// }
