package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGrantV3ParseResourcePermissions(t *testing.T) {
	assert := assert.New(t)

	m := map[string]ResourcePermissions{
		"channel": ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	r := parseResourcePermissions(m)
	for _, v := range r {
		assert.Equal(int64(31), v)
	}
}

func TestGrantV3ParseResourcePermissions2(t *testing.T) {
	assert := assert.New(t)

	m := map[string]ResourcePermissions{
		"channel": ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	r := parseResourcePermissions(m)
	for _, v := range r {
		assert.Equal(int64(15), v)
	}
}

func TestGrantV3ParseResourcePermissions3(t *testing.T) {
	assert := assert.New(t)

	m := map[string]ResourcePermissions{
		"channel": ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
	}

	r := parseResourcePermissions(m)
	for _, v := range r {
		assert.Equal(int64(7), v)
	}
}

func TestGrantV3ParseResourcePermissions4(t *testing.T) {
	assert := assert.New(t)

	m := map[string]ResourcePermissions{
		"channel": ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
		"channel2": ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  false,
			Manage: true,
			Delete: false,
		},
	}

	r := parseResourcePermissions(m)
	assert.Equal(int64(7), r["channel"])
	assert.Equal(int64(5), r["channel2"])
}

func TestGrantV3(t *testing.T) {
	AssertTestGrantV3(t, true, false)
}

func AssertTestGrantV3(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	o := newGrantBuilder(pn)
	if testContext {
		o = newGrantBuilderWithContext(pn, backgroundContext)
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	ch := map[string]ResourcePermissions{
		"channel": ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
	}

	u := map[string]ResourcePermissions{
		"users": ResourcePermissions{
			Create: false,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	s := map[string]ResourcePermissions{
		"spaces": ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
		},
	}

	cg := map[string]ResourcePermissions{
		"cg": ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
		},
		"cg2": ResourcePermissions{
			Create: true,
			Read:   true,
			Write:  true,
			Manage: false,
			Delete: false,
		},
	}

	o.TTL(100)
	o.Channels(ch)
	o.Users(u)
	o.Spaces(s)
	o.ChannelGroups(cg)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(grantPath, pn.Config.SubscribeKey),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"ttl\":100,\"permissions\":{\"resources\":{\"channels\":{\"channel\":7},\"groups\":{\"cg\":23,\"cg2\":19},\"users\":{\"users\":15},\"spaces\":{\"spaces\":31}},\"patterns\":{\"channels\":null,\"groups\":null,\"users\":null,\"spaces\":null},\"meta\":null}}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}
}
