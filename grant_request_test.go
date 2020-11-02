package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGrantRequestObjects(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"message":"Success","payload":{"level":"uuid","subscribe_key":"sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306","ttl":1440,"uuids":{"ch1":{"auths":{"pam-key":{"r":0,"w":0,"m":0,"d":0,"g":1,"u":1,"j":1}}},"ch2":{"auths":{"pam-key":{"r":0,"w":0,"m":0,"d":0,"g":1,"u":1,"j":1}}}}},"service":"Access Manager","status":200}`)

	e, _, err := newGrantResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
	assert.Equal(false, e.ManageEnabled)
	assert.Equal(false, e.ReadEnabled)
	assert.Equal(false, e.WriteEnabled)
	assert.Equal(1440, e.TTL)
	assert.Equal(false, e.UUIDs["ch1"].AuthKeys["pam-key"].ManageEnabled)
	assert.Equal(false, e.UUIDs["ch1"].AuthKeys["pam-key"].ReadEnabled)
	assert.Equal(false, e.UUIDs["ch1"].AuthKeys["pam-key"].WriteEnabled)
	assert.Equal(false, e.UUIDs["ch1"].AuthKeys["pam-key"].DeleteEnabled)
	assert.Equal(false, e.UUIDs["ch2"].AuthKeys["pam-key"].ManageEnabled)
	assert.Equal(false, e.UUIDs["ch2"].AuthKeys["pam-key"].ReadEnabled)
	assert.Equal(false, e.UUIDs["ch2"].AuthKeys["pam-key"].WriteEnabled)
	assert.Equal(false, e.UUIDs["ch2"].AuthKeys["pam-key"].DeleteEnabled)
	assert.Equal(true, e.UUIDs["ch1"].AuthKeys["pam-key"].GetEnabled)
	assert.Equal(true, e.UUIDs["ch1"].AuthKeys["pam-key"].UpdateEnabled)
	assert.Equal(true, e.UUIDs["ch1"].AuthKeys["pam-key"].JoinEnabled)
	assert.Equal(true, e.UUIDs["ch2"].AuthKeys["pam-key"].GetEnabled)
	assert.Equal(true, e.UUIDs["ch2"].AuthKeys["pam-key"].UpdateEnabled)
	assert.Equal(true, e.UUIDs["ch2"].AuthKeys["pam-key"].JoinEnabled)

}

func TestNewGrantObjectsPremsBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newGrantBuilder(pubnub)
	o.AuthKeys([]string{"my-auth-key"})
	o.UUIDs([]string{"uuid"})
	o.Get(true)
	o.Update(true)
	o.Join(true)
	o.TTL(5000)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/auth/grant/sub-key/%s", o.opts.pubnub.Config.SubscribeKey),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("auth", "my-auth-key")
	expected.Set("target-uuid", "uuid")
	expected.Set("r", "0")
	expected.Set("w", "0")
	expected.Set("m", "0")
	expected.Set("d", "0")
	expected.Set("g", "1")
	expected.Set("u", "1")
	expected.Set("j", "1")
	expected.Set("ttl", "5000")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid", "timestamp"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestGrantRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := &grantOpts{
		AuthKeys:      []string{"my-auth-key"},
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		Read:          true,
		Write:         true,
		Manage:        true,
		TTL:           5000,
		setTTL:        true,
		pubnub:        pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/auth/grant/sub-key/%s", opts.pubnub.Config.SubscribeKey),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("auth", "my-auth-key")
	expected.Set("channel", "ch")
	expected.Set("channel-group", "cg")
	expected.Set("r", "1")
	expected.Set("w", "1")
	expected.Set("m", "1")
	expected.Set("d", "0")
	expected.Set("ttl", "5000")
	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid", "timestamp"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestGrantRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &grantOpts{
		AuthKeys:      []string{"my-auth-key"},
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		Read:          true,
		Write:         true,
		Manage:        true,
		TTL:           5000,
		setTTL:        true,
		pubnub:        pubnub,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("auth", "my-auth-key")
	expected.Set("channel", "ch")
	expected.Set("channel-group", "cg")
	expected.Set("r", "1")
	expected.Set("w", "1")
	expected.Set("m", "1")
	expected.Set("d", "0")
	expected.Set("ttl", "5000")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid", "timestamp"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)

}

func TestNewGrantBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newGrantBuilder(pubnub)
	o.AuthKeys([]string{"my-auth-key"})
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})
	o.Read(true)
	o.Write(true)
	o.Manage(true)
	o.Delete(true)
	o.TTL(5000)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/auth/grant/sub-key/%s", o.opts.pubnub.Config.SubscribeKey),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("auth", "my-auth-key")
	expected.Set("channel", "ch")
	expected.Set("channel-group", "cg")
	expected.Set("r", "1")
	expected.Set("w", "1")
	expected.Set("m", "1")
	expected.Set("d", "1")
	expected.Set("ttl", "5000")
	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid", "timestamp"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewGrantBuilderDelFalse(t *testing.T) {
	assert := assert.New(t)
	o := newGrantBuilder(pubnub)
	o.AuthKeys([]string{"my-auth-key"})
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})
	o.Read(true)
	o.Write(true)
	o.Manage(true)
	o.Delete(false)
	o.TTL(5000)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/auth/grant/sub-key/%s", o.opts.pubnub.Config.SubscribeKey),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("auth", "my-auth-key")
	expected.Set("channel", "ch")
	expected.Set("channel-group", "cg")
	expected.Set("r", "1")
	expected.Set("w", "1")
	expected.Set("m", "1")
	expected.Set("d", "0")
	expected.Set("ttl", "5000")
	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid", "timestamp"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewGrantBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newGrantBuilderWithContext(pubnub, backgroundContext)
	o.AuthKeys([]string{"my-auth-key"})
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})
	o.Read(true)
	o.Write(true)
	o.Manage(true)
	o.TTL(5000)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/auth/grant/sub-key/%s", o.opts.pubnub.Config.SubscribeKey),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("auth", "my-auth-key")
	expected.Set("channel", "ch")
	expected.Set("channel-group", "cg")
	expected.Set("r", "1")
	expected.Set("w", "1")
	expected.Set("m", "1")
	expected.Set("d", "0")
	expected.Set("ttl", "5000")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")
	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid", "timestamp"}, []string{})

	body, err := o.opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestGrantTokenOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &grantOpts{
		AuthKeys:      []string{"my-auth-key"},
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		Read:          true,
		Write:         true,
		Manage:        true,
		TTL:           5000,
		setTTL:        true,
		pubnub:        pn,
	}

	assert.Equal("pubnub/validation: pubnub: Grant: Missing Subscribe Key", opts.validate().Error())
}

func TestGrantTokenOptsValidateSec(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SecretKey = ""
	opts := &grantOpts{
		AuthKeys:      []string{"my-auth-key"},
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		Read:          true,
		Write:         true,
		Manage:        true,
		TTL:           5000,
		setTTL:        true,
		pubnub:        pn,
	}

	assert.Equal("pubnub/validation: pubnub: Grant: Missing Secret Key", opts.validate().Error())
}

func TestGrantTokenOptsValidatePub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := &grantOpts{
		AuthKeys:      []string{"my-auth-key"},
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
		Read:          true,
		Write:         true,
		Manage:        true,
		TTL:           5000,
		setTTL:        true,
		pubnub:        pn,
	}

	assert.Equal("pubnub/validation: pubnub: Grant: Missing Publish Key", opts.validate().Error())
}

func TestNewGrantResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newGrantResponse(jsonBytes, StatusResponse{})
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestNewGrantResponseManageEnabled(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"message":"Success","payload":{"level":"channel-group+auth","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"r":1,"m":1,"w":1,"channels":{"ch1":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"ch2":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"ch3":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}}},"channel-groups":{"cg1":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"cg2":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0,"ttl":1},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"cg3":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}}}},"service":"Access Manager","status":200}`)

	e, _, err := newGrantResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
	assert.Equal(true, e.ManageEnabled)
	assert.Equal(true, e.ReadEnabled)
	assert.Equal(true, e.WriteEnabled)
	assert.Equal(1440, e.TTL)
	assert.Equal(true, e.ChannelGroups["cg1"].AuthKeys["my-auth-key-1"].ManageEnabled)
	assert.Equal(true, e.ChannelGroups["cg1"].AuthKeys["my-auth-key-1"].ReadEnabled)
	assert.Equal(true, e.ChannelGroups["cg1"].AuthKeys["my-auth-key-1"].WriteEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].ManageEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].ReadEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].WriteEnabled)
	assert.Equal(1440, e.ChannelGroups["cg1"].TTL)
	assert.Equal(true, e.Channels["ch1"].AuthKeys["my-auth-key-1"].ManageEnabled)
	assert.Equal(true, e.Channels["ch1"].AuthKeys["my-auth-key-1"].ReadEnabled)
	assert.Equal(true, e.Channels["ch1"].AuthKeys["my-auth-key-1"].WriteEnabled)

}

func TestNewGrantResponseManageEnabledInv(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"message":"Success","payload":{"level":"channel-group+auth","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":0,"r":0,"m":0,"w":0,"channels":{"ch1":{"auths":{"my-auth-key-1":{"r":0,"w":0,"m":0,"d":1},"my-auth-key-2":{"r":0,"w":0,"m":0,"d":1}}},"ch2":{"auths":{"my-auth-key-1":{"r":0,"w":0,"m":0,"d":1},"my-auth-key-2":{"r":0,"w":0,"m":0,"d":1}}},"ch3":{"auths":{"my-auth-key-1":{"r":0,"w":0,"m":0,"d":1},"my-auth-key-2":{"r":0,"w":0,"m":0,"d":1}}}},"channel-groups":{"cg1":{"auths":{"my-auth-key-1":{"r":0,"w":0,"m":0,"d":1,"ttl":4},"my-auth-key-2":{"r":0,"w":0,"m":0,"d":1}}},"cg2":{"auths":{"my-auth-key-1":{"r":0,"w":0,"m":0,"d":1,"ttl":6},"my-auth-key-2":{"r":0,"w":0,"m":0,"d":1}}},"cg3":{"auths":{"my-auth-key-1":{"r":0,"w":0,"m":0,"d":1},"my-auth-key-2":{"r":0,"w":0,"m":0,"d":1}}}}},"service":"Access Manager","status":200}`)

	e, _, err := newGrantResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
	assert.Equal(false, e.ManageEnabled)
	assert.Equal(false, e.ReadEnabled)
	assert.Equal(false, e.WriteEnabled)
	assert.Equal(0, e.TTL)
	assert.Equal(false, e.ChannelGroups["cg1"].AuthKeys["my-auth-key-1"].ManageEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].AuthKeys["my-auth-key-1"].ReadEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].AuthKeys["my-auth-key-1"].WriteEnabled)
	assert.Equal(0, e.ChannelGroups["cg1"].TTL)
	assert.Equal(false, e.ChannelGroups["cg1"].ManageEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].ReadEnabled)
	assert.Equal(false, e.ChannelGroups["cg1"].WriteEnabled)

	assert.Equal(0, e.ChannelGroups["cg1"].TTL)
	assert.Equal(false, e.Channels["ch1"].AuthKeys["my-auth-key-1"].ManageEnabled)
	assert.Equal(false, e.Channels["ch1"].AuthKeys["my-auth-key-1"].ReadEnabled)
	assert.Equal(false, e.Channels["ch1"].AuthKeys["my-auth-key-1"].WriteEnabled)
}

func TestNewGrantResponseManageEnabledCH(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"message":"Success","payload":{"level":"user","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel":"ch1","auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`)

	_, _, err := newGrantResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
}

func TestNewGrantResponseManageEnabledCHM(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`{"message":"Success","payload":{"level":"user","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel":"ch1","auths":{"my-pam-key":{"r":1,"w":1,"m":1,"d":0}}},"service":"Access Manager","status":200}`)

	_, _, err := newGrantResponse(jsonBytes, StatusResponse{})

	assert.Nil(err)
}

func TestGrantTTL(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	gb := newGrantBuilder(pn)
	gb.TTL(10)
	assert.Equal(10, gb.opts.TTL)
}
