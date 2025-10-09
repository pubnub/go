package e2e

import (
	"log"
	"os"
	"reflect"
	"testing"

	pubnub "github.com/pubnub/go/v7"
	"github.com/stretchr/testify/assert"
)

func TestObjectsV2ChannelMetadataSetUpdateGetRemove(t *testing.T) {
	a := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	id := randomized("testchannel")
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	name := randomized("name")
	desc := "desc"
	custom := map[string]interface{}{"a": "b", "c": "d"}
	status := "active"
	channelType := "public"

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
		pubnub.PNChannelMetadataIncludeStatus,
		pubnub.PNChannelMetadataIncludeType,
	}

	defer removeChannelMetadata(a, pn, id)
	res, st, err := pn.SetChannelMetadata().Include(incl).Channel(id).Name(name).Description(desc).Custom(custom).Status(status).Type(channelType).Execute()

	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(name, res.Data.Name)
		a.Equal(desc, res.Data.Description)
		a.NotNil(res.Data.Updated)
		a.NotNil(res.Data.ETag)
		a.True(reflect.DeepEqual(custom, res.Data.Custom))
		a.Equal(status, res.Data.Status)
		a.Equal(channelType, res.Data.Type)
	}

	desc = "desc2"
	statusUpdated := "inactive"

	res, st, err = pn.SetChannelMetadata().Include(incl).Channel(id).Name(name).Description(desc).Custom(custom).Status(statusUpdated).Type(channelType).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(desc, res.Data.Description)
		a.Equal(statusUpdated, res.Data.Status)
		a.Equal(channelType, res.Data.Type)
	}

	getRes, st, err := pn.GetChannelMetadata().Include(incl).Channel(id).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if getRes != nil {
		a.Equal(statusUpdated, getRes.Data.Status)
		a.Equal(channelType, getRes.Data.Type)
	}

}

func removeChannelMetadata(a *assert.Assertions, pn *pubnub.PubNub, id string) {
	res, st, err := pn.RemoveChannelMetadata().Channel(id).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Nil(res.Data)
	}
}

func TestObjectsV2UUIDMetadataSetUpdateGetRemove(t *testing.T) {
	a := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	id := randomized("testuuid")
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	name := randomized("name")
	email := "go@pubnub.com"
	custom := map[string]interface{}{"a": "b", "c": "d"}
	status := "active"
	uuidType := "public"

	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
		pubnub.PNUUIDMetadataIncludeStatus,
		pubnub.PNUUIDMetadataIncludeType,
	}

	defer removeUUIDMetadata(a, pn, id)
	res, st, err := pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).Email(email).Custom(custom).Status(status).Type(uuidType).Execute()

	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(name, res.Data.Name)
		a.Equal(email, res.Data.Email)
		a.NotNil(res.Data.Updated)
		a.NotNil(res.Data.ETag)
		a.True(reflect.DeepEqual(custom, res.Data.Custom))
		a.Equal(status, res.Data.Status)
		a.Equal(uuidType, res.Data.Type)
	}

	email = "gosdk@pubnub.com"
	statusUpdated := "inactive"

	res, st, err = pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).Email(email).Custom(custom).Status(statusUpdated).Type(uuidType).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(email, res.Data.Email)
		a.Equal(statusUpdated, res.Data.Status)
		a.Equal(uuidType, res.Data.Type)
	}

	getRes, st, err := pn.GetUUIDMetadata().Include(incl).UUID(id).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if getRes != nil {
		a.Equal(statusUpdated, getRes.Data.Status)
		a.Equal(uuidType, getRes.Data.Type)
	}

}

func removeUUIDMetadata(a *assert.Assertions, pn *pubnub.PubNub, id string) {
	res, st, err := pn.RemoveUUIDMetadata().UUID(id).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Nil(res.Data)
	}
}

func TestObjectsV2MembersAddRemove(t *testing.T) {
	a := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	channelid := randomized("channel")
	userid := randomized("uuid")
	inc := []pubnub.PNChannelMembersInclude{pubnub.PNChannelMembersIncludeUUID}

	defer removeChannelMembers(a, pn, channelid, userid)

	res, st, err := pn.
		SetChannelMembers().
		Channel(channelid).
		Set([]pubnub.PNChannelMembersSet{{UUID: pubnub.PNChannelMembersUUID{ID: userid}}}).
		Include(inc).
		Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if err == nil {
		a.True(len(res.Data) > 0)
	}

}

func removeChannelMembers(a *assert.Assertions, pn *pubnub.PubNub, channelid string, userid string) {
	_, st, err := pn.
		RemoveChannelMembers().
		Channel(channelid).
		Remove([]pubnub.PNChannelMembersRemove{{UUID: pubnub.PNChannelMembersUUID{ID: userid}}}).
		Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
}

func TestObjectsV2MembershipAddRemove(t *testing.T) {
	a := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())
	channelid := randomized("channel")
	userid := randomized("uuid")
	inc := []pubnub.PNMembershipsInclude{pubnub.PNMembershipsIncludeChannel}

	defer removeMemberships(a, pn, channelid, userid)

	res, st, err := pn.
		SetMemberships().
		UUID(userid).
		Set([]pubnub.PNMembershipsSet{{Channel: pubnub.PNMembershipsChannel{ID: channelid}}}).
		Include(inc).
		Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if err == nil {
		a.True(len(res.Data) > 0)
	}

}

func removeMemberships(a *assert.Assertions, pn *pubnub.PubNub, channelid string, userid string) {
	_, st, err := pn.
		RemoveMemberships().
		UUID(userid).
		Remove([]pubnub.PNMembershipsRemove{{Channel: pubnub.PNMembershipsChannel{ID: channelid}}}).
		Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
}
