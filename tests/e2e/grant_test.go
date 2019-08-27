package e2e

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	pubnub "github.com/sprucehealth/pubnub-go"
	"github.com/sprucehealth/pubnub-go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestGrantParseLogsForAuthKey(t *testing.T) {

	assert := assert.New(t)
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.SecretKey = "sec-key"
	pn.Config.AuthKey = "myAuthKey"

	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	pn.Grant().
		Read(true).Write(true).Manage(true).
		Channels([]string{"ch1", "ch2"}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-tic.C:
		tic.Stop()
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	//fmt.Printf("Captured: %s", out)

	s := fmt.Sprintf("%s", out)
	// //https://ps.pndsn.com/v2/auth/grant/sub-key/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64?w=1&m=1&channel=ch1,ch2&timestamp=1535719943&auth=myAuthKey&pnsdk=PubNub-Go/4.1.3&uuid=pn-621c7b2a-f87c-4362-bd1e-6c6762dfc667&r=1&signature=PntTQe-zBfJa6AvN4bu4u0txG_TOoksHGod7OnijmwM=
	// expected := fmt.Sprintf("https://%s/v2/auth/grant/sub-key/%s?&uuid=%sw=1&m=1&channel=ch1,ch2",
	// 	pn.Config.Origin,
	// 	pn.Config.SubscribeKey,
	// )

	// assert.Contains(s, expected)

	//auth=myAuthKey&pnsdk=PubNub-Go/4.1.3
	expected2 := fmt.Sprintf("auth=%s",
		pn.Config.AuthKey)

	assert.Contains(s, expected2)

}

func TestGrantParseLogsForMultipleAuthKeys(t *testing.T) {

	assert := assert.New(t)
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	pn := pubnub.NewPubNub(configCopy())
	pn.Config.SecretKey = "sec-key"
	pn.Config.AuthKey = "myAuthKey"

	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	pn.Grant().
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{"authkey1", "authkey2"}).
		Channels([]string{"ch1", "ch2"}).
		Execute()

	tic := time.NewTicker(time.Duration(timeout) * time.Second)
	select {
	case <-tic.C:
		tic.Stop()
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	//fmt.Printf("Captured: %s", out)

	s := fmt.Sprintf("%s", out)

	//https://ps.pndsn.com/v2/auth/grant/sub-key/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64?m=1&auth=authkey1,authkey2&channel=ch1,ch2&timestamp=1535719219&pnsdk=PubNub-Go/4.1.3&uuid=pn-a83164fe-7ecf-42ab-ba14-d2d8e6eabd7a&r=1&w=1&signature=0SkyfvohAq8_0phVi0YhCL4c2ZRSPBVwCwQ9fANvPmM=
	assert.Contains(s, "auth=authkey1,authkey2")
}

func TestGrantSucccessNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	pn.Config.UUID = "asd,|//&aqwe"

	res, _, err := pn.Grant().
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{"pam-key"}).Channels([]string{"ch1", "ch2"}).
		Execute()

	assert.Nil(err)
	log.Println(res)
	assert.NotNil(res)

	assert.True(res.Channels["ch2"].AuthKeys["pam-key"].WriteEnabled)
	assert.True(res.Channels["ch2"].AuthKeys["pam-key"].ReadEnabled)
	assert.True(res.Channels["ch2"].AuthKeys["pam-key"].ManageEnabled)
	assert.True(!res.Channels["ch2"].AuthKeys["pam-key"].DeleteEnabled)

}

func TestGrantSucccessAppLevelFalse(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	pn.Config.UUID = "asd,|//&aqwe"

	res, _, err := pn.Grant().
		Read(false).Write(false).Manage(false).Delete(false).
		Execute()

	assert.Nil(err)
	log.Println(res)
	assert.NotNil(res)

	assert.True(!res.WriteEnabled)
	assert.True(!res.ReadEnabled)
	assert.True(!res.ManageEnabled)
	assert.True(!res.DeleteEnabled)

}

func TestGrantSucccessAppLevelMixed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	pn.Config.UUID = "asd,|//&aqwe"

	res, _, err := pn.Grant().
		Read(false).Write(true).Manage(false).Delete(true).
		Execute()

	assert.Nil(err)
	log.Println(res)
	assert.NotNil(res)

	assert.True(res.WriteEnabled)
	assert.True(!res.ReadEnabled)
	assert.True(!res.ManageEnabled)
	assert.True(res.DeleteEnabled)

}

func TestGrantSucccessAppLevelMixed2(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	pn.Config.UUID = "asd,|//&aqwe"

	res, _, err := pn.Grant().
		Read(true).Write(false).Manage(true).Delete(false).
		Execute()

	assert.Nil(err)
	log.Println(res)
	assert.NotNil(res)

	assert.True(!res.WriteEnabled)
	assert.True(res.ReadEnabled)
	assert.True(res.ManageEnabled)
	assert.True(!res.DeleteEnabled)

}

func TestGrantSucccessNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())

	pn.Config.UUID = "asd,|//&aqwe"

	res, _, err := pn.GrantWithContext(backgroundContext).
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{"pam-key"}).Channels([]string{"ch1", "ch2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
}

func TestGrantMultipleMixed(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "auth=my-auth-key-1%2Cmy-auth-key-2&channel=ch1%2Cch2%2Cch3&channel-group=cg1%2Ccg2%2Ccg3&r=1&m=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel-group+auth","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channels":{"ch1":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"ch2":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"ch3":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}}},"channel-groups":{"cg1":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"cg2":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}},"cg3":{"auths":{"my-auth-key-1":{"r":1,"w":1,"m":1,"d":0},"my-auth-key-2":{"r":1,"w":1,"m":1,"d":0}}}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "timestamp", "signature"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	res, _, err := pn.Grant().
		Read(true).Write(true).Manage(true).
		AuthKeys([]string{"my-auth-key-1", "my-auth-key-2"}).
		Channels([]string{"ch1", "ch2", "ch3"}).
		ChannelGroups([]string{"cg1", "cg2", "cg3"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
}

func TestGrantSingleChannel(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "channel=ch1&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channels":{"ch1":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())
	pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	res, _, err := pn.Grant().
		Read(true).Write(true).
		Channels([]string{"ch1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.True(res.Channels["ch1"].WriteEnabled)
	assert.True(res.Channels["ch1"].ReadEnabled)
	assert.False(res.Channels["ch1"].ManageEnabled)
}

func TestGrantSingleChannelWithAuth(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "auth=my-pam-key&channel=ch1&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"user","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel":"ch1","auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		Read(true).Write(true).Manage(false).
		AuthKeys([]string{"my-pam-key"}).
		Channels([]string{"ch1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)
	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.Channels["ch1"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "channel=ch1%2Cch2&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channels":{"ch1":{"r":1,"w":1,"m":0,"d":0},"ch2":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		Read(true).Write(true).
		Channels([]string{"ch1", "ch2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.Channels["ch1"].WriteEnabled)
	assert.True(res.Channels["ch1"].ReadEnabled)
	assert.False(res.Channels["ch1"].ManageEnabled)

	assert.True(res.Channels["ch2"].WriteEnabled)
	assert.True(res.Channels["ch2"].ReadEnabled)
	assert.False(res.Channels["ch2"].ManageEnabled)
}

func TestGrantMultipleChannelsWithAuth(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "auth=my-pam-key&channel=ch1%2Cch2&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"user","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channels":{"ch1":{"auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}},"ch2":{"auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		Read(true).Write(true).
		AuthKeys([]string{"my-pam-key"}).
		Channels([]string{"ch1", "ch2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.Channels["ch1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.Channels["ch1"].AuthKeys["my-pam-key"].ManageEnabled)

	assert.True(res.Channels["ch2"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.Channels["ch2"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.Channels["ch2"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantSingleGroup(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "channel-group=cg1&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel-group","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel-groups":{"cg1":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		Read(true).Write(true).
		ChannelGroups([]string{"cg1"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].ManageEnabled)
}

func TestGrantSingleGroupWithAuth(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "auth=my-pam-key&channel-group=cg1&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel-group+auth","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel-groups":"cg1","auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		ChannelGroups([]string{"cg1"}).
		AuthKeys([]string{"my-pam-key"}).
		Write(true).
		Read(true).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ManageEnabled)
}

func TestGrantMultipleGroups(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "channel-group=cg1%2Ccg2&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel-group","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel-groups":{"cg1":{"r":1,"w":1,"m":0,"d":0},"cg2":{"r":1,"w":1,"m":0,"d":0}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		Read(true).Write(true).
		ChannelGroups([]string{"cg1", "cg2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].ManageEnabled)

	assert.True(res.ChannelGroups["cg2"].WriteEnabled)
	assert.True(res.ChannelGroups["cg2"].ReadEnabled)
	assert.False(res.ChannelGroups["cg2"].ManageEnabled)
}

func TestGrantMultipleGroupsWithAuth(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v2/auth/grant/sub-key/%s", pamConfig.SubscribeKey),
		Query:              "auth=my-pam-key&channel-group=cg1%2Ccg2&m=0&r=1&w=1&d=0",
		ResponseBody:       `{"message":"Success","payload":{"level":"channel-group+auth","subscribe_key":"sub-c-b9ab9508-43cf-11e8-9967-869954283fb4","ttl":1440,"channel-groups":{"cg1":{"auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}},"cg2":{"auths":{"my-pam-key":{"r":1,"w":1,"m":0,"d":0}}}}},"service":"Access Manager","status":200}`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk", "signature", "timestamp"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(pamConfigCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Grant().
		Read(true).Write(true).
		AuthKeys([]string{"my-pam-key"}).
		ChannelGroups([]string{"cg1", "cg2"}).
		Execute()

	assert.Nil(err)
	assert.NotNil(res)

	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.ChannelGroups["cg1"].AuthKeys["my-pam-key"].ManageEnabled)

	assert.True(res.ChannelGroups["cg2"].AuthKeys["my-pam-key"].WriteEnabled)
	assert.True(res.ChannelGroups["cg2"].AuthKeys["my-pam-key"].ReadEnabled)
	assert.False(res.ChannelGroups["cg2"].AuthKeys["my-pam-key"].ManageEnabled)
}
