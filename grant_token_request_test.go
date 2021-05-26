package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGrantTokenParseResourcePermissions(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	m := map[string]UserSpacePermissions{
		"channel": UserSpacePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNUsers)
	for _, v := range r {
		assert.Equal(int64(31), v)
	}
}

func TestGrantTokenParseResourcePermissions2(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	m := map[string]UserSpacePermissions{
		"channel": UserSpacePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNUsers)
	for _, v := range r {
		assert.Equal(int64(15), v)
	}
}

func TestGrantTokenParseResourcePermissions3(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	m := map[string]UserSpacePermissions{
		"channel": UserSpacePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNUsers)
	for _, v := range r {
		assert.Equal(int64(7), v)
	}
}

func TestGrantTokenParseResourcePermissions4(t *testing.T) {
	assert := assert.New(t)

	pn := NewPubNub(NewDemoConfig())
	o := newGrantTokenBuilder(pn)

	m := map[string]UserSpacePermissions{
		"channel": UserSpacePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
		"channel2": UserSpacePermissions{
			Create: false,
			Read:   true,
			Write:  false,
			Manage: true,
			Delete: false,
		},
	}

	r := o.opts.parseResourcePermissions(m, PNUsers)
	assert.Equal(int64(7), r["channel"])
	assert.Equal(int64(5), r["channel2"])
}

func TestGrantToken(t *testing.T) {
	AssertTestGrantToken(t, true, false)
}

func AssertTestGrantToken(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGrantTokenBuilder(pn)
	if testContext {
		o = newGrantTokenBuilderWithContext(pn, backgroundContext)
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	// ch := map[string]ChannelPermissions{
	// 	"channel": ChannelPermissions{
	// 		Write:  false,
	// 		Read:   true,
	// 		Delete: false,
	// 	},
	// }

	u := map[string]UserSpacePermissions{
		"users": UserSpacePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	s := map[string]UserSpacePermissions{
		"spaces": UserSpacePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	// cg := map[string]GroupPermissions{
	// 	"cg": GroupPermissions{
	// 		Read:   true,
	// 		Manage: true,
	// 	},
	// 	"cg2": GroupPermissions{
	// 		Read:   true,
	// 		Manage: false,
	// 	},
	// }

	o.TTL(100)
	//o.Channels(ch)
	o.Users(u)
	o.Spaces(s)
	//o.ChannelGroups(cg)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(grantTokenPath, pn.Config.SubscribeKey),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	//expectedBody := "{\"ttl\":100,\"permissions\":{\"resources\":{\"channels\":{\"channel\":1},\"groups\":{\"cg\":5,\"cg2\":1},\"users\":{\"users\":15},\"spaces\":{\"spaces\":31}},\"patterns\":{\"channels\":{},\"groups\":{},\"users\":{},\"spaces\":{}},\"meta\":{}}}" //{\"ttl\":100,\"permissions\":{\"resources\":{\"channels\":{\"channel\":7},\"groups\":{\"cg\":23,\"cg2\":19},\"users\":{\"users\":15},\"spaces\":{\"spaces\":31}},\"patterns\":{\"channels\":null,\"groups\":null,\"users\":null,\"spaces\":null},\"meta\":null}}"
	expectedBody := "{\"ttl\":100,\"permissions\":{\"resources\":{\"channels\":{},\"groups\":{},\"users\":{\"users\":15},\"spaces\":{\"spaces\":31}},\"patterns\":{\"channels\":{},\"groups\":{},\"users\":{},\"spaces\":{}},\"meta\":{}}}"
	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}
}
