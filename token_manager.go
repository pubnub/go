package pubnub

import (
	"fmt"
	"sync"
)

type TokenManager struct {
	sync.RWMutex
	Tokens GrantResourcesWithPermissons
}

func newTokenManager(pubnub *PubNub, ctx Context) *TokenManager {
	g := GrantResourcesWithPermissons{
		Channels: make(map[string]PermissonsWithToken),
		Groups:   make(map[string]PermissonsWithToken),
		Users:    make(map[string]PermissonsWithToken),
		Spaces:   make(map[string]PermissonsWithToken),
	}

	manager := &TokenManager{
		Tokens: g,
	}
	return manager
}

func (m *TokenManager) FindToken() *GrantResourcesWithPermissons {
	g := GrantResourcesWithPermissons{}
	return &g
}

func (m *TokenManager) StoreToken(token string) {
	fmt.Println("--->", token)
	cborObject, err := DecodeCBORToken(token)
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
		g := ParseGrantResources(cborObject.Resources, token)
		m.Lock()
		mergeTokensByResources(m.Tokens.Channels, g.Channels)
		mergeTokensByResources(m.Tokens.Users, g.Users)
		mergeTokensByResources(m.Tokens.Groups, g.Groups)
		mergeTokensByResources(m.Tokens.Spaces, g.Spaces)
		m.Unlock()

		fmt.Println(" --- Tokens ---- ", m.Tokens)

		fmt.Println(" --- Patterns")
		ParseGrantResources(cborObject.Patterns, token)
	}
}

func mergeTokensByResources(m map[string]PermissonsWithToken, r map[string]PermissonsWithToken) {
	for k, v := range r {
		m[k] = v
	}

}
