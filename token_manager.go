package pubnub

import (
	"fmt"
	"sync"
)

type TokenManager struct {
	sync.RWMutex
	Tokens GrantResourcesWithPermissons
	pubnub *PubNub
}

// Match tokens for subscribe or other calls
// Check ttl expiration

func newTokenManager(pubnub *PubNub, ctx Context) *TokenManager {

	g := GrantResourcesWithPermissons{
		Channels: make(map[string]ChannelPermissonsWithToken),
		Groups:   make(map[string]GroupPermissonsWithToken),
		Users:    make(map[string]UserSpacePermissonsWithToken),
		Spaces:   make(map[string]UserSpacePermissonsWithToken),
	}

	manager := &TokenManager{
		Tokens: g,
	}
	manager.pubnub = pubnub

	return manager
}

func (m *TokenManager) GetTokens(channels, groups, users, spaces []string) *GrantResourcesWithPermissons {
	g := GrantResourcesWithPermissons{
		Channels: make(map[string]ChannelPermissonsWithToken),
		Groups:   make(map[string]GroupPermissonsWithToken),
		Users:    make(map[string]UserSpacePermissonsWithToken),
		Spaces:   make(map[string]UserSpacePermissonsWithToken),
	}
	findTokenInTokens(channels, g.Channels, m.Tokens.Channels, PNChannels)
	//findTokenInTokensChannels(channels, g.Channels, m.Tokens.Channels)
	findTokenInTokensGroups(groups, g.Groups, m.Tokens.Groups)
	findTokenInTokensUserSpace(users, g.Users, m.Tokens.Users)
	findTokenInTokensUserSpace(spaces, g.Spaces, m.Tokens.Spaces)

	return &g
}

func matchTokensForSubscribe(g *GrantResourcesWithPermissons) {

}

func findTokenInTokens(r []string, resource, merge interface{}, resourceType PNResourceType) {
	switch resourceType {
	case PNChannels:
		a := resource.(map[string]ChannelPermissions)
		e := merge.(map[string]ChannelPermissions)
		for _, k := range r {
			if d, ok := e[k]; ok {
				a[k] = d
			}
		}
	case PNGroups:
		a := resource.(map[string]GroupPermissions)
		e := merge.(map[string]GroupPermissions)
		for _, k := range r {
			if d, ok := e[k]; ok {
				a[k] = d
			}
		}
	default:
		//case PNUsers:
		//case PNSpaces:
		a := resource.(map[string]UserSpacePermissions)
		e := merge.(map[string]UserSpacePermissions)
		for _, k := range r {
			if d, ok := e[k]; ok {
				a[k] = d
			}
		}
	}

}

func findTokenInTokensChannels(r []string, a, m map[string]ChannelPermissonsWithToken) {
	for _, k := range r {
		if d, ok := m[k]; ok {
			a[k] = d
		}
	}
}

func findTokenInTokensGroups(r []string, a, m map[string]GroupPermissonsWithToken) {
	for _, k := range r {
		if d, ok := m[k]; ok {
			a[k] = d
		}
	}
}

func findTokenInTokensUserSpace(r []string, a, m map[string]UserSpacePermissonsWithToken) {
	for _, k := range r {
		if d, ok := m[k]; ok {
			a[k] = d
		}
	}
}

// func mergeTokensByChannels(m map[string]ChannelPermissonsWithToken, r map[string]ChannelPermissonsWithToken) {
// 	for k, v := range r {
// 		m[k] = v
// 	}

// }
// func mergeTokensByGroups(m map[string]GroupPermissonsWithToken, r map[string]GroupPermissonsWithToken) {
// 	for k, v := range r {
// 		m[k] = v
// 	}

// }
// func mergeTokensByUserSpace(m map[string]UserSpacePermissonsWithToken, r map[string]UserSpacePermissonsWithToken) {
// 	for k, v := range r {
// 		m[k] = v
// 	}

// }

func mergeTokensByResource(m interface{}, resource interface{}, resourceType PNResourceType) {
	switch resourceType {
	case PNChannels:
		c := resource.(map[string]ChannelPermissonsWithToken)
		d := m.(map[string]ChannelPermissonsWithToken)
		for k, v := range c {
			d[k] = v
		}
	case PNGroups:
		c := resource.(map[string]GroupPermissonsWithToken)
		d := m.(map[string]GroupPermissonsWithToken)
		for k, v := range c {
			d[k] = v
		}
	default:
		//case PNUsers:
		//case PNSpaces:
		c := resource.(map[string]UserSpacePermissonsWithToken)
		d := m.(map[string]UserSpacePermissonsWithToken)
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
			fmt.Println("Timetoken: ", cborObject.Timetoken)
			fmt.Println("TTL: ", cborObject.TTL)
			fmt.Println(fmt.Sprintf("Meta: %#v", cborObject.Meta))
			fmt.Println("")
			fmt.Println(" --- Resources")
			g := ParseGrantResources(cborObject.Resources, token, cborObject.Timetoken)
			m.Lock()
			mergeTokensByResource(m.Tokens.Channels, g.Channels, PNChannels)
			mergeTokensByResource(m.Tokens.Users, g.Users, PNUsers)
			mergeTokensByResource(m.Tokens.Groups, g.Groups, PNGroups)
			mergeTokensByResource(m.Tokens.Spaces, g.Spaces, PNSpaces)

			// mergeTokensByChannels(m.Tokens.Channels, g.Channels)
			// mergeTokensByUserSpace(m.Tokens.Users, g.Users)
			// mergeTokensByGroups(m.Tokens.Groups, g.Groups)
			// mergeTokensByUserSpace(m.Tokens.Spaces, g.Spaces)
			m.Unlock()

			fmt.Println(" --- Tokens ---- ", m.Tokens)

			fmt.Println(" --- Patterns")
			ParseGrantResources(cborObject.Patterns, token, cborObject.Timetoken)
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

// func mergeTokensByResources(m map[string]PermissonsWithToken, r map[string]PermissonsWithToken) map[string]PermissonsWithToken {
// 	for k, v := range r {
// 		m[k] = v
// 	}
// 	return m
// }
