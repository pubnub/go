package e2e

import (
	pubnub "github.com/pubnub/go/v7"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
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

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	defer removeChannelMetadata(a, pn, id)
	res, st, err := pn.SetChannelMetadata().Include(incl).Channel(id).Name(name).Description(desc).Custom(custom).Execute()

	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(name, res.Data.Name)
		a.Equal(desc, res.Data.Description)
		a.NotNil(res.Data.Updated)
		a.NotNil(res.Data.ETag)
		a.True(reflect.DeepEqual(custom, res.Data.Custom))
	}

	desc = "desc2"

	res, st, err = pn.SetChannelMetadata().Include(incl).Channel(id).Name(name).Description(desc).Custom(custom).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(desc, res.Data.Description)
	}

	_, st, err = pn.GetChannelMetadata().Include(incl).Channel(id).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)

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

	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	defer removeUUIDMetadata(a, pn, id)
	res, st, err := pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).Email(email).Custom(custom).Execute()

	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(name, res.Data.Name)
		a.Equal(email, res.Data.Email)
		a.NotNil(res.Data.Updated)
		a.NotNil(res.Data.ETag)
		a.True(reflect.DeepEqual(custom, res.Data.Custom))
	}

	email = "gosdk@pubnub.com"

	res, st, err = pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).Email(email).Custom(custom).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)
	if res != nil {
		a.Equal(id, res.Data.ID)
		a.Equal(email, res.Data.Email)
	}

	_, st, err = pn.GetUUIDMetadata().Include(incl).UUID(id).Execute()
	a.Nil(err)
	a.Equal(200, st.StatusCode)

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
