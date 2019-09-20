package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenManager(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	pn.Config.SecretKey = ""

	tm := newTokenManager(pn, nil)

	//{map[] map[] map[testuser_16669:{{true true true true true} 31  1568805412 3}] map[testspace_15011:{{true true true true true} 31 p0F2AkF0Gl2CEiRDdHRsA0NyZXOkRGNoYW6gQ2dycKBDdXNyoW50ZXN0dXNlcl8xNjY2ORgfQ3NwY6FvdGVzdHNwYWNlXzE1MDExGB9DcGF0pERjaGFuoENncnCgQ3VzcqBDc3BjoERtZXRhoENzaWdYIMqDoIOYPP9ULfXKLDK3eoGQ-C8nJxPTWFCDAc-Flxu7 1568805412 3}] map[] map[] map[^.*:{{true true true true true} 31 p0F2AkF0Gl2CEiVDdHRsA0NyZXOkRGNoYW6gQ2dycKBDdXNyoENzcGOgQ3BhdKREY2hhbqBDZ3JwoEN1c3KhY14uKhgfQ3NwY6FjXi4qGB9EbWV0YaBDc2lnWCDfqMStM0r1GgghNjt1MPeSaA0ADTw6aGsuQgMT3jYylg== 1568805413 3}] map[^.*:{{true true true true true} 31  1568805413 3}]}

	t1 := "p0F2AkF0Gl2AX-JDdHRsCkNyZXOkRGNoYW6gQ2dycKBDdXNyoWl1LTMzNTIwNTUPQ3NwY6Fpcy0xNzA3OTgzGB9DcGF0pERjaGFuoENncnCgQ3VzcqBDc3BjoERtZXRhoENzaWdYINqGs2EyEMHPZrp6znVqTBzXNBAD_31hUH3JuUSWE2A6"
	t2 := "p0F2AkF0Gl2AaMlDdHRsCkNyZXOkRGNoYW6gQ2dycKBDdXNyoWl1LTE5NzQxMDcPQ3NwY6Fpcy0yMzExMDExGB9DcGF0pERjaGFuoENncnCgQ3VzcqBDc3BjoERtZXRhoENzaWdYIO1ti19DLbEKK-s_COJPlM1xtZCpP8K4sV51nvRPTIxf"
	t3 := "p0F2AkF0Gl2CEiRDdHRsA0NyZXOkRGNoYW6gQ2dycKBDdXNyoW50ZXN0dXNlcl8xNjY2ORgfQ3NwY6FvdGVzdHNwYWNlXzE1MDExGB9DcGF0pERjaGFuoENncnCgQ3VzcqBDc3BjoERtZXRhoENzaWdYIMqDoIOYPP9ULfXKLDK3eoGQ-C8nJxPTWFCDAc-Flxu7"
	t4 := "p0F2AkF0Gl2CEiVDdHRsA0NyZXOkRGNoYW6gQ2dycKBDdXNyoENzcGOgQ3BhdKREY2hhbqBDZ3JwoEN1c3KhY14uKhgfQ3NwY6FjXi4qGB9EbWV0YaBDc2lnWCDfqMStM0r1GgghNjt1MPeSaA0ADTw6aGsuQgMT3jYylg=="

	tm.StoreTokens([]string{t1, t2, t3, t4})

	g := tm.GetAllTokens()

	assert.Equal(0, len(g.Channels))
	assert.Equal(0, len(g.Groups))
	assert.Equal(0, len(g.ChannelsPattern))
	assert.Equal(0, len(g.GroupsPattern))
	assert.Equal(t3, g.Users["testuser_16669"].Token)
	assert.Equal(int64(31), g.Users["testuser_16669"].BitMaskPerms)
	assert.Equal(3, g.Users["testuser_16669"].TTL)
	assert.Equal(int64(1568805412), g.Users["testuser_16669"].Timestamp)
	assert.Equal(true, g.Users["testuser_16669"].Permissions.Read)
	assert.Equal(true, g.Users["testuser_16669"].Permissions.Write)
	assert.Equal(true, g.Users["testuser_16669"].Permissions.Delete)
	assert.Equal(true, g.Users["testuser_16669"].Permissions.Create)
	assert.Equal(true, g.Users["testuser_16669"].Permissions.Manage)

	assert.Equal(t3, g.Spaces["testspace_15011"].Token)
	assert.Equal(true, g.Spaces["testspace_15011"].Permissions.Read)
	assert.Equal(true, g.Spaces["testspace_15011"].Permissions.Write)
	assert.Equal(true, g.Spaces["testspace_15011"].Permissions.Delete)
	assert.Equal(true, g.Spaces["testspace_15011"].Permissions.Create)
	assert.Equal(true, g.Spaces["testspace_15011"].Permissions.Manage)

	assert.Equal(t2, g.Users["u-1974107"].Token)
	assert.Equal(t1, g.Spaces["s-1707983"].Token)

	assert.Equal(t4, g.UsersPattern["^.*"].Token)
	assert.Equal(t4, g.SpacesPattern["^.*"].Token)

	g2 := tm.GetTokensByResource(PNUsers)
	assert.Equal(t3, g2.Users["testuser_16669"].Token)
	assert.Equal(int64(31), g2.Users["testuser_16669"].BitMaskPerms)
	assert.Equal(3, g2.Users["testuser_16669"].TTL)
	assert.Equal(int64(1568805412), g2.Users["testuser_16669"].Timestamp)
	assert.Equal(true, g2.Users["testuser_16669"].Permissions.Read)
	assert.Equal(true, g2.Users["testuser_16669"].Permissions.Write)
	assert.Equal(true, g2.Users["testuser_16669"].Permissions.Delete)
	assert.Equal(true, g2.Users["testuser_16669"].Permissions.Create)
	assert.Equal(true, g2.Users["testuser_16669"].Permissions.Manage)
	assert.Equal(t2, g2.Users["u-1974107"].Token)
	assert.Equal(t4, g2.UsersPattern["^.*"].Token)

	g3 := tm.GetTokensByResource(PNSpaces)
	assert.Equal(t3, g3.Spaces["testspace_15011"].Token)
	assert.Equal(true, g3.Spaces["testspace_15011"].Permissions.Read)
	assert.Equal(true, g3.Spaces["testspace_15011"].Permissions.Write)
	assert.Equal(true, g3.Spaces["testspace_15011"].Permissions.Delete)
	assert.Equal(true, g3.Spaces["testspace_15011"].Permissions.Create)
	assert.Equal(true, g3.Spaces["testspace_15011"].Permissions.Manage)
	assert.Equal(t1, g3.Spaces["s-1707983"].Token)
	assert.Equal(t4, g3.SpacesPattern["^.*"].Token)

	g4 := tm.GetToken("testspace_15011", PNSpaces)
	assert.Equal(t3, g4)

	g5 := tm.GetToken("testuser_16669", PNUsers)
	assert.Equal(t3, g5)

	g6 := tm.GetToken("^.*", PNSpaces)
	assert.Equal(t4, g6)

	g7 := tm.GetToken("^.*", PNUsers)
	assert.Equal(t4, g7)

	g8 := tm.GetToken("NONEXISTENT", PNSpaces)
	assert.Equal(t4, g8)

	g9 := tm.GetToken("NONEXISTENT", PNUsers)
	assert.Equal(t4, g9)

}
