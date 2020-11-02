package e2e

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func ActivateWithPAMV2() *pubnub.PubNub {
	pn := pubnub.NewPubNub(pamConfigCopy())
	return pn
}

func RunGrantV2(pn *pubnub.PubNub, users, spaces []string, read, write, manage, del, create, createPattern bool) []string {
	r := GenRandom()

	authkey := fmt.Sprintf("authkey_%d", r.Intn(99999))

	res, _, _ := pn.Grant().
		Read(true).Write(true).Manage(true).
		Get(true).Update(true).Join(true).
		UUIDs(append(users, spaces...)).AuthKeys([]string{authkey}).
		Execute()

	if res != nil {
		return []string{authkey}
	}
	return []string{}
}

func SetPNV2(pn, pn2 *pubnub.PubNub, tokens []string) {
	pn.Config.SubscribeKey = pn2.Config.SubscribeKey
	pn.Config.PublishKey = pn2.Config.PublishKey
	pn.Config.SecretKey = ""
	pn.Config.Origin = pn2.Config.Origin
	pn.Config.Secure = pn2.Config.Secure

	pn.Config.AuthKey = tokens[0]
}

func TestObjectsV2CreateUpdateGetDeleteUUID(t *testing.T) {
	ObjectsCreateUpdateGetDeleteUUIDCommon(t, false, false)
}

func TestObjectsV2CreateUpdateGetDeleteUUIDWithPAM(t *testing.T) {
	ObjectsCreateUpdateGetDeleteUUIDCommon(t, true, false)
}

// func TestObjectsV2CreateUpdateGetDeleteUUIDWithPAMWithoutSecKey(t *testing.T) {
// 	ObjectsCreateUpdateGetDeleteUUIDCommon(t, true, true)
// }

func ObjectsCreateUpdateGetDeleteUUIDCommon(t *testing.T, withPAM, runWithoutSecretKey bool) {

	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	id := fmt.Sprintf("testuser_%d", r.Intn(99999))
	if withPAM {
		pn2 := ActivateWithPAMV2()
		if runWithoutSecretKey {
			tokens := RunGrantV2(pn2, []string{id}, []string{}, true, true, true, true, true, true)
			SetPNV2(pn, pn2, tokens)
		} else {
			pn = pn2
			RunGrantV2(pn, []string{id}, []string{}, true, true, false, true, true, false)
		}

	}

	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	name := fmt.Sprintf("name_%d", r.Intn(99999))
	extid := "extid"
	purl := "profileurl"
	email := "email"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	res, st, err := pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, st.StatusCode)
	if res != nil {
		assert.Equal(id, res.Data.ID)
		assert.Equal(name, res.Data.Name)
		assert.Equal(extid, res.Data.ExternalID)
		assert.Equal(purl, res.Data.ProfileURL)
		assert.Equal(email, res.Data.Email)
		// assert.NotNil(res.Data.Created)
		assert.NotNil(res.Data.Updated)
		assert.NotNil(res.Data.ETag)
		assert.Equal("b", res.Data.Custom["a"])
		assert.Equal("d", res.Data.Custom["c"])
	}

	email = "email2"

	res2, st2, err2 := pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err2)
	assert.Equal(200, st2.StatusCode)
	if res2 != nil {
		assert.Equal(id, res2.Data.ID)
		assert.Equal(name, res2.Data.Name)
		assert.Equal(extid, res2.Data.ExternalID)
		assert.Equal(purl, res2.Data.ProfileURL)
		assert.Equal(email, res2.Data.Email)
		// assert.Equal(res.Data.Created, res2.Data.Created)
		assert.NotNil(res2.Data.Updated)
		assert.NotNil(res2.Data.ETag)
		assert.Equal("b", res2.Data.Custom["a"])
		assert.Equal("d", res2.Data.Custom["c"])
	}

	res3, st3, err3 := pn.GetUUIDMetadata().UUID(id).Include(incl).Execute()
	assert.Nil(err3)
	assert.Equal(200, st3.StatusCode)
	if res3 != nil {
		assert.Equal(id, res3.Data.ID)
		assert.Equal(name, res3.Data.Name)
		assert.Equal(extid, res3.Data.ExternalID)
		assert.Equal(purl, res3.Data.ProfileURL)
		assert.Equal(email, res3.Data.Email)
		// assert.Equal(res.Data.Created, res3.Data.Created)
		if res2 != nil {
			assert.Equal(res2.Data.Updated, res3.Data.Updated)
			assert.Equal(res2.Data.ETag, res3.Data.ETag)
		}
		assert.Equal("b", res3.Data.Custom["a"])
		assert.Equal("d", res3.Data.Custom["c"])
	}

	//getusers
	sort := []string{"updated:desc"}
	if withPAM {
		res6, st6, err6 := pn.GetAllUUIDMetadata().Include(incl).Sort(sort).Limit(100).Count(true).Execute()
		assert.Nil(err6)
		assert.Equal(200, st6.StatusCode)
		assert.True(res6.TotalCount > 0)
		found := false
		for i := range res6.Data {
			if res6.Data[i].ID == id {
				assert.Equal(name, res6.Data[i].Name)
				assert.Equal(extid, res6.Data[i].ExternalID)
				assert.Equal(purl, res6.Data[i].ProfileURL)
				assert.Equal(email, res6.Data[i].Email)
				// assert.Equal(res.Data.Created, res6.Data[i].Created)
				if res2 != nil {
					assert.Equal(res2.Data.Updated, res6.Data[i].Updated)
					assert.Equal(res2.Data.ETag, res6.Data[i].ETag)
				}
				assert.Equal("b", res6.Data[i].Custom["a"])
				assert.Equal("d", res6.Data[i].Custom["c"])
				found = true
			}
		}
		assert.True(found)

		res6F, st6F, err6F := pn.GetAllUUIDMetadata().Include(incl).Limit(100).Filter("name == '" + name + "'").Count(true).Execute()
		assert.Nil(err6F)
		assert.Equal(200, st6F.StatusCode)
		assert.True(res6F.TotalCount > 0)
		foundF := false
		for i := range res6F.Data {
			//fmt.Println(res6F.Data[i], id)
			if res6F.Data[i].ID == id {
				assert.Equal(name, res6F.Data[i].Name)
				assert.Equal(extid, res6F.Data[i].ExternalID)
				assert.Equal(purl, res6F.Data[i].ProfileURL)
				assert.Equal(email, res6F.Data[i].Email)
				// assert.Equal(res.Data.Created, res6F.Data[i].Created)
				assert.Equal(res2.Data.Updated, res6F.Data[i].Updated)
				assert.Equal(res2.Data.ETag, res6F.Data[i].ETag)
				assert.Equal("b", res6F.Data[i].Custom["a"])
				assert.Equal("d", res6F.Data[i].Custom["c"])
				foundF = true
			}
		}
		assert.True(foundF)
	}

	//delete
	res5, st5, err5 := pn.RemoveUUIDMetadata().UUID(id).Execute()
	assert.Nil(err5)
	assert.Equal(200, st5.StatusCode)
	if res5 != nil {
		assert.Nil(res5.Data)
	}

	//getuser
	res4, st4, err4 := pn.GetUUIDMetadata().Include(incl).UUID(id).Execute()
	assert.NotNil(err4)
	if res5 != nil {
		assert.Nil(res4)
	}
	assert.Equal(404, st4.StatusCode)

}

func TestObjectsV2CreateUpdateGetDeleteChannel(t *testing.T) {
	ObjectsCreateUpdateGetDeleteChannelCommon(t, false, false)
}

func TestObjectsV2CreateUpdateGetDeleteChannelWithPAM(t *testing.T) {
	ObjectsCreateUpdateGetDeleteChannelCommon(t, true, false)
}

// func TestObjectsV2CreateUpdateGetDeleteChannelWithPAMWithoutSecKey(t *testing.T) {
// 	ObjectsCreateUpdateGetDeleteChannelCommon(t, true, true)
// }

func ObjectsCreateUpdateGetDeleteChannelCommon(t *testing.T, withPAM, runWithoutSecretKey bool) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	id := fmt.Sprintf("testspace_%d", r.Intn(99999))

	if withPAM {
		pn2 := ActivateWithPAMV2()
		if runWithoutSecretKey {
			tokens := RunGrantV2(pn2, []string{}, []string{id}, true, true, false, true, true, true)
			SetPNV2(pn, pn2, tokens)
		} else {
			pn = pn2
			RunGrantV2(pn, []string{}, []string{id}, true, true, false, true, true, false)
		}

	}
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	name := fmt.Sprintf("name_%d", r.Intn(99999))
	desc := "desc"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	res, st, err := pn.SetChannelMetadata().Include(incl).Channel(id).Name(name).Description(desc).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, st.StatusCode)
	if res != nil {
		assert.Equal(id, res.Data.ID)
		assert.Equal(name, res.Data.Name)
		assert.Equal(desc, res.Data.Description)
		// assert.NotNil(res.Data.Created)
		assert.NotNil(res.Data.Updated)
		assert.NotNil(res.Data.ETag)
		assert.Equal("b", res.Data.Custom["a"])
		assert.Equal("d", res.Data.Custom["c"])
	}

	desc = "desc2"

	res2, st2, err2 := pn.SetChannelMetadata().Include(incl).Channel(id).Name(name).Description(desc).Custom(custom).Execute()
	assert.Nil(err2)
	assert.Equal(200, st2.StatusCode)
	if res2 != nil {
		assert.Equal(id, res2.Data.ID)
		assert.Equal(name, res2.Data.Name)
		assert.Equal(desc, res2.Data.Description)
		// assert.Equal(res.Data.Created, res2.Data.Created)
		assert.NotNil(res2.Data.Updated)
		assert.NotNil(res2.Data.ETag)
		assert.Equal("b", res2.Data.Custom["a"])
		assert.Equal("d", res2.Data.Custom["c"])
	}

	res3, st3, err3 := pn.GetChannelMetadata().Include(incl).Channel(id).Execute()
	assert.Nil(err3)
	assert.Equal(200, st3.StatusCode)
	if res3 != nil {
		assert.Equal(id, res3.Data.ID)
		assert.Equal(name, res3.Data.Name)
		assert.Equal(desc, res3.Data.Description)
		// assert.Equal(res.Data.Created, res3.Data.Created)
		assert.Equal(res2.Data.Updated, res3.Data.Updated)
		assert.Equal(res2.Data.ETag, res3.Data.ETag)
		assert.Equal("b", res3.Data.Custom["a"])
		assert.Equal("d", res3.Data.Custom["c"])
	}

	sort := []string{"updated:desc"}
	//getusers
	if withPAM {
		res6, st6, err6 := pn.GetAllChannelMetadata().Include(incl).Sort(sort).Limit(100).Count(true).Execute()
		assert.Nil(err6)
		assert.Equal(200, st6.StatusCode)
		found := false
		if res6 != nil {
			assert.True(res6.TotalCount > 0)

			for i := range res6.Data {
				if res6.Data[i].ID == id {
					assert.Equal(name, res6.Data[i].Name)
					assert.Equal(desc, res6.Data[i].Description)
					// assert.Equal(res.Data.Created, res6.Data[i].Created)
					assert.Equal(res2.Data.Updated, res6.Data[i].Updated)
					assert.Equal(res2.Data.ETag, res6.Data[i].ETag)
					assert.Equal("b", res6.Data[i].Custom["a"])
					assert.Equal("d", res6.Data[i].Custom["c"])
					found = true
				}
			}
		}
		assert.True(found)

		res6F, st6F, err6F := pn.GetAllChannelMetadata().Include(incl).Limit(100).Filter("name like '" + name + "*'").Count(true).Execute()
		assert.Nil(err6F)
		assert.Equal(200, st6F.StatusCode)
		foundF := false
		if res6F != nil {
			assert.True(res6F.TotalCount > 0)

			for i := range res6F.Data {
				if res6F.Data[i].ID == id {
					assert.Equal(name, res6F.Data[i].Name)
					assert.Equal(desc, res6F.Data[i].Description)
					// assert.Equal(res.Data.Created, res6F.Data[i].Created)
					assert.Equal(res2.Data.Updated, res6F.Data[i].Updated)
					assert.Equal(res2.Data.ETag, res6F.Data[i].ETag)
					assert.Equal("b", res6F.Data[i].Custom["a"])
					assert.Equal("d", res6F.Data[i].Custom["c"])
					foundF = true
				}
			}
		}
		assert.True(foundF)

	}

	//delete
	res5, st5, err5 := pn.RemoveChannelMetadata().Channel(id).Execute()
	assert.Nil(err5)
	assert.Equal(200, st5.StatusCode)
	if res5 != nil {
		assert.Nil(res5.Data)
	}

	//getuser
	res4, st4, err4 := pn.GetChannelMetadata().Include(incl).Channel(id).Execute()
	assert.NotNil(err4)
	if res4 != nil {
		assert.Nil(res4)
	}
	assert.Equal(404, st4.StatusCode)

}

func TestObjectsV2SetRemoveMembershipsV2(t *testing.T) {
	ObjectsSetRemoveMembershipsCommonV2(t, false, false)
}

func TestObjectsV2SetRemoveMembershipsV2WithPAM(t *testing.T) {
	ObjectsSetRemoveMembershipsCommonV2(t, true, false)
}

// PASSES after adding PAM checks for Update Members
// func TestObjectsV2SetRemoveMembershipsV2WithPAMWithoutSecKey(t *testing.T) {
// 	ObjectsSetRemoveMembershipsCommonV2(t, true, true)
// }

func ObjectsSetRemoveMembershipsCommonV2(t *testing.T, withPAM, runWithoutSecretKey bool) {
	assert := assert.New(t)

	limit := 100
	count := true

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	userid := fmt.Sprintf("testuser_%d", r.Intn(99999))
	spaceid := fmt.Sprintf("testspace_%d", r.Intn(99999))

	if withPAM {
		pn2 := ActivateWithPAMV2()
		if runWithoutSecretKey {
			tokens := RunGrantV2(pn2, []string{userid}, []string{spaceid}, true, true, true, true, true, true)
			SetPNV2(pn, pn2, tokens)
		} else {
			pn = pn2
			RunGrantV2(pn, []string{userid}, []string{spaceid}, true, true, true, true, true, false)
		}

	}
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	name := fmt.Sprintf("name_%d", r.Intn(99999))
	extid := "extid"
	purl := "profileurl"
	email := "email"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	inclUUID := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	res, st, err := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, st.StatusCode)
	if res != nil {
		assert.Equal(userid, res.Data.ID)
		assert.Equal(name, res.Data.Name)
		assert.Equal(extid, res.Data.ExternalID)
		assert.Equal(purl, res.Data.ProfileURL)
		assert.Equal(email, res.Data.Email)
		// assert.NotNil(res.Data.Created)
		assert.NotNil(res.Data.Updated)
		assert.NotNil(res.Data.ETag)
		assert.Equal("b", res.Data.Custom["a"])
		assert.Equal("d", res.Data.Custom["c"])
	}

	desc := "desc"
	custom2 := make(map[string]interface{})
	custom2["a1"] = "b1"
	custom2["c1"] = "d1"

	inclChannel := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	res2, st2, err2 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err2)
	assert.Equal(200, st2.StatusCode)
	//fmt.Println("res2-->", res2)
	if res2 != nil {
		assert.Equal(spaceid, res2.Data.ID)
		assert.Equal(name, res2.Data.Name)
		assert.Equal(desc, res2.Data.Description)
		// assert.NotNil(res2.Data.Created)
		assert.NotNil(res2.Data.Updated)
		assert.NotNil(res2.Data.ETag)
		assert.Equal("b1", res2.Data.Custom["a1"])
		assert.Equal("d1", res2.Data.Custom["c1"])
	}

	userid2 := fmt.Sprintf("testuser_%d", r.Intn(99999))

	_, st3, err3 := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid2).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err3)
	assert.Equal(200, st3.StatusCode)

	spaceid2 := fmt.Sprintf("testspace_%d", r.Intn(99999))

	_, st4, err4 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid2).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err4)
	assert.Equal(200, st4.StatusCode)

	userid3 := fmt.Sprintf("testuser_%d", r.Intn(99999))

	_, stuser3, erruser3 := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid3).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(erruser3)
	assert.Equal(200, stuser3.StatusCode)

	spaceid3 := fmt.Sprintf("testspace_%d", r.Intn(99999))

	_, stspace3, errspace3 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid3).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(errspace3)
	assert.Equal(200, stspace3.StatusCode)

	inclSm := []pubnub.PNChannelMembersInclude{
		pubnub.PNChannelMembersIncludeUUIDCustom,
		pubnub.PNChannelMembersIncludeUUID,
		pubnub.PNChannelMembersIncludeCustom,
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	uuid := pubnub.PNChannelMembersUUID{
		ID: userid,
	}

	in := pubnub.PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom3,
	}
	uuid3 := pubnub.PNChannelMembersUUID{
		ID: userid3,
	}

	inUser3 := pubnub.PNChannelMembersSet{
		UUID:   uuid3,
		Custom: custom3,
	}

	inArr := []pubnub.PNChannelMembersSet{
		in,
		inUser3,
	}

	//Add Space Memberships
	sort := []string{"updated:desc"}

	resAdd, stAdd, errAdd := pn.SetChannelMembers().Channel(spaceid).Sort(sort).Set(inArr).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errAdd)
	assert.Equal(200, stAdd.StatusCode)
	if errAdd == nil {
		sortMembers1 := false
		sortMembers2 := false

		found := false
		assert.True(resAdd.TotalCount > 0)
		//fmt.Println("resAdd-->", resAdd)

		for i := range resAdd.Data {
			if resAdd.Data[i].UUID.ID == userid {
				found = true
				assert.Equal(custom3["a3"], resAdd.Data[i].Custom["a3"])
				assert.Equal(custom3["c3"], resAdd.Data[i].Custom["c3"])
				assert.Equal(userid, resAdd.Data[i].UUID.ID)
				assert.Equal(name, resAdd.Data[i].UUID.Name)
				assert.Equal(extid, resAdd.Data[i].UUID.ExternalID)
				assert.Equal(purl, resAdd.Data[i].UUID.ProfileURL)
				assert.Equal(email, resAdd.Data[i].UUID.Email)
				assert.Equal(custom["a"], resAdd.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resAdd.Data[i].UUID.Custom["c"])
			}
		}
		if (resAdd.Data != nil) && (len(resAdd.Data) > 1) {
			sortMembers1 = (resAdd.Data[1].UUID.ID == userid)
			sortMembers2 = (resAdd.Data[0].UUID.ID == userid3)
			assert.True(sortMembers1)
			assert.True(sortMembers2)
		} else {
			assert.Fail("Sort ", "resAdd.Data null or ", len(resAdd.Data))
		}

		assert.True(found)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMembers->", errAdd.Error())
		}
	}

	//Update Space Memberships
	if !withPAM {

		custom4 := make(map[string]interface{})
		custom4["a2"] = "b2"
		custom4["c2"] = "d2"

		up := pubnub.PNChannelMembersSet{
			UUID:   uuid,
			Custom: custom4,
		}

		upArr := []pubnub.PNChannelMembersSet{
			up,
		}

		resUp, stUp, errUp := pn.SetChannelMembers().Channel(spaceid).Sort(sort).Set(upArr).Include(inclSm).Limit(limit).Count(count).Execute()
		assert.Nil(errUp)
		assert.Equal(200, stUp.StatusCode)
		if errUp == nil {
			assert.True(resUp.TotalCount > 0)
			foundUp := false
			for i := range resUp.Data {
				if resUp.Data[i].UUID.ID == userid {
					foundUp = true
					assert.Equal("b2", resUp.Data[i].Custom["a2"])
					assert.Equal("d2", resUp.Data[i].Custom["c2"])
					//assert.Equal(userid, resAdd.Data[i].UUID.ID)
					assert.Equal(name, resAdd.Data[i].UUID.Name)
					assert.Equal(extid, resAdd.Data[i].UUID.ExternalID)
					assert.Equal(purl, resAdd.Data[i].UUID.ProfileURL)
					assert.Equal(email, resAdd.Data[i].UUID.Email)
					assert.Equal(custom["a"], resAdd.Data[i].UUID.Custom["a"])
					assert.Equal(custom["c"], resAdd.Data[i].UUID.Custom["c"])

				}
			}
			assert.True(foundUp)
		} else {
			if enableDebuggingInTests {

				fmt.Println("ManageMembers->", errUp.Error())
			}
		}
	}
	//Get Space Memberships

	inclMemberships := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsIncludeCustom,
		pubnub.PNMembershipsIncludeChannel,
		pubnub.PNMembershipsIncludeChannelCustom,
	}

	//fmt.Println("GetMemberships ====>")

	resGetMem, stGetMem, errGetMem := pn.GetMemberships().UUID(userid).Include(inclMemberships).Sort(sort).Limit(limit).Count(count).Execute()
	foundGetMem := false
	assert.Nil(errGetMem)
	if errGetMem == nil {
		for i := range resGetMem.Data {
			if resGetMem.Data[i].Channel.ID == spaceid {
				foundGetMem = true
				assert.Equal(name, resGetMem.Data[i].Channel.Name)
				assert.Equal(desc, resGetMem.Data[i].Channel.Description)
				assert.Equal("b1", resGetMem.Data[i].Channel.Custom["a1"])
				assert.Equal("d1", resGetMem.Data[i].Channel.Custom["c1"])
				if withPAM {
					assert.Equal("b3", resGetMem.Data[i].Custom["a3"])
					assert.Equal("d3", resGetMem.Data[i].Custom["c3"])
				} else {
					assert.Equal("b2", resGetMem.Data[i].Custom["a2"])
					assert.Equal("d2", resGetMem.Data[i].Custom["c2"])
				}
			}
		}
		assert.Equal(200, stGetMem.StatusCode)
		assert.True(foundGetMem)
	} else {
		if enableDebuggingInTests {
			fmt.Println("GetMemberships->", errGetMem.Error())
		}
	}

	//filterExp := fmt.Sprintf("custom.c3 == '%s' || custom.c2 == '%s'", "d3", "d2")
	filterExp := fmt.Sprintf("channel.name == '%s'", name)

	//fmt.Println("GetMemberships ====>", filterExp)

	resGetMemF, stGetMemF, errGetMemF := pn.GetMemberships().UUID(userid).Include(inclMemberships).Filter(filterExp).Limit(limit).Count(count).Execute()
	foundGetMemF := false
	assert.Nil(errGetMemF)
	if errGetMemF == nil {
		for i := range resGetMemF.Data {
			if resGetMemF.Data[i].Channel.ID == spaceid {
				foundGetMemF = true
				assert.Equal(name, resGetMemF.Data[i].Channel.Name)
				assert.Equal(desc, resGetMemF.Data[i].Channel.Description)
				assert.Equal("b1", resGetMemF.Data[i].Channel.Custom["a1"])
				assert.Equal("d1", resGetMemF.Data[i].Channel.Custom["c1"])
				if withPAM {
					assert.Equal("b3", resGetMemF.Data[i].Custom["a3"])
					assert.Equal("d3", resGetMemF.Data[i].Custom["c3"])
				} else {
					assert.Equal("b2", resGetMemF.Data[i].Custom["a2"])
					assert.Equal("d2", resGetMemF.Data[i].Custom["c2"])
				}
			}
		}
		assert.Equal(200, stGetMemF.StatusCode)
		assert.True(foundGetMemF)
	} else {
		if enableDebuggingInTests {

			fmt.Println("GetMemberships->", errGetMemF.Error())
		}
	}

	//Remove Space Memberships
	re := pubnub.PNChannelMembersRemove{
		UUID: uuid,
	}

	reArr := []pubnub.PNChannelMembersRemove{
		re,
	}
	resRem, stRem, errRem := pn.RemoveChannelMembers().Channel(spaceid).Remove(reArr).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errRem)
	assert.Equal(200, stRem.StatusCode)
	//fmt.Println("====>stRem.StatusCode", stRem.StatusCode)
	if errRem == nil {

		foundRem := false
		for i := range resRem.Data {
			if resRem.Data[i].UUID.ID == userid {
				foundRem = true
				assert.Equal("b2", resRem.Data[i].Custom["a2"])
				assert.Equal("d2", resRem.Data[i].Custom["c2"])
				assert.Equal(userid, resRem.Data[i].UUID.ID)
				assert.Equal(name, resRem.Data[i].UUID.Name)
				assert.Equal(extid, resRem.Data[i].UUID.ExternalID)
				assert.Equal(purl, resRem.Data[i].UUID.ProfileURL)
				assert.Equal(email, resRem.Data[i].UUID.Email)
				assert.Equal(custom["a"], resRem.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resRem.Data[i].UUID.Custom["c"])

			}
		}
		assert.False(foundRem)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMembers->", errRem.Error())
		}
	}

	channel2 := pubnub.PNMembershipsChannel{
		ID: spaceid2,
	}

	inMem := pubnub.PNMembershipsSet{
		Channel: channel2,
		Custom:  custom3,
	}

	channel3 := pubnub.PNMembershipsChannel{
		ID: spaceid3,
	}

	inMemSpace3 := pubnub.PNMembershipsSet{
		Channel: channel3,
		Custom:  custom3,
	}

	inArrMem := []pubnub.PNMembershipsSet{
		inMem,
		inMemSpace3,
	}

	//Add user memberships
	resManageMemAdd, stManageMemAdd, errManageMemAdd := pn.SetMemberships().UUID(userid2).Set(inArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	//fmt.Println("resManageMemAdd -->", resManageMemAdd)
	assert.Nil(errManageMemAdd)
	assert.Equal(200, stManageMemAdd.StatusCode)
	if errManageMemAdd == nil {
		foundManageMembers := false
		for i := range resManageMemAdd.Data {
			if resManageMemAdd.Data[i].Channel.ID == spaceid2 {
				assert.Equal(spaceid2, resManageMemAdd.Data[i].Channel.ID)
				assert.Equal(name, resManageMemAdd.Data[i].Channel.Name)
				assert.Equal(desc, resManageMemAdd.Data[i].Channel.Description)
				assert.Equal(custom2["a1"], resManageMemAdd.Data[i].Channel.Custom["a1"])
				assert.Equal(custom2["c1"], resManageMemAdd.Data[i].Channel.Custom["c1"])
				assert.Equal(custom3["a3"], resManageMemAdd.Data[i].Custom["a3"])
				assert.Equal(custom3["c3"], resManageMemAdd.Data[i].Custom["c3"])
				foundManageMembers = true
			}
		}
		assert.True(foundManageMembers)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMemberships->", errManageMemAdd.Error())
		}
	}

	// //Update user memberships

	custom5 := make(map[string]interface{})
	custom5["a5"] = "b5"
	custom5["c5"] = "d5"

	upMem := pubnub.PNMembershipsSet{
		Channel: channel2,
		Custom:  custom5,
	}

	upArrMem := []pubnub.PNMembershipsSet{
		upMem,
	}
	sortMemberships1 := false
	sortMemberships2 := false

	resManageMemUp, stManageMemUp, errManageMemUp := pn.SetMemberships().UUID(userid2).Sort(sort).Set(upArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	//fmt.Println("resManageMemUp -->", resManageMemUp)
	assert.Nil(errManageMemUp)
	assert.Equal(200, stManageMemUp.StatusCode)
	if errManageMemUp == nil {
		foundManageMembersUp := false

		for i := range resManageMemUp.Data {
			//fmt.Println("resManageMemUp.Data[i].ID == spaceid2-->", resRem.Data[i].UUID.ID, spaceid2)
			if resManageMemUp.Data[i].Channel.ID == spaceid2 {
				assert.Equal(spaceid2, resManageMemUp.Data[i].Channel.ID)
				assert.Equal(name, resManageMemUp.Data[i].Channel.Name)
				assert.Equal(desc, resManageMemUp.Data[i].Channel.Description)
				assert.Equal(custom2["a1"], resManageMemAdd.Data[i].Channel.Custom["a1"])
				assert.Equal(custom2["c1"], resManageMemAdd.Data[i].Channel.Custom["c1"])
				assert.Equal(custom5["a5"], resManageMemUp.Data[i].Custom["a5"])
				assert.Equal(custom5["c5"], resManageMemUp.Data[i].Custom["c5"])
				foundManageMembersUp = true
			}
		}
		if (resManageMemUp.Data != nil) && (len(resManageMemUp.Data) > 1) {
			sortMemberships1 = (resManageMemUp.Data[0].Channel.ID == spaceid2)
			sortMemberships2 = (resManageMemUp.Data[1].Channel.ID == spaceid3)
			assert.True(sortMemberships1)
			assert.True(sortMemberships2)
		} else {
			assert.Fail("Sort ", "resManageMemUp.Data null or ", len(resManageMemUp.Data))
		}

		assert.True(foundManageMembersUp)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMemberships->", errManageMemUp.Error())
		}
	}

	// //Get members
	resGetMembers, stGetMembers, errGetMembers := pn.GetChannelMembers().Channel(spaceid2).Include(inclSm).Limit(limit).Count(count).Execute()
	//fmt.Println("resGetMembers -->", resGetMembers)
	assert.Nil(errGetMembers)
	assert.Equal(200, stGetMembers.StatusCode)
	if errGetMembers == nil {
		foundGetMembers := false
		for i := range resGetMembers.Data {
			if resGetMembers.Data[i].UUID.ID == userid2 {
				foundGetMembers = true
				assert.Equal(name, resGetMembers.Data[i].UUID.Name)
				assert.Equal(extid, resGetMembers.Data[i].UUID.ExternalID)
				assert.Equal(purl, resGetMembers.Data[i].UUID.ProfileURL)
				assert.Equal(email, resGetMembers.Data[i].UUID.Email)
				assert.Equal(custom["a"], resGetMembers.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resGetMembers.Data[i].UUID.Custom["c"])
				assert.Equal(custom5["a5"], resGetMembers.Data[i].Custom["a5"])
				assert.Equal(custom5["c5"], resGetMembers.Data[i].Custom["c5"])
			}
		}

		assert.True(foundGetMembers)
	} else {
		if enableDebuggingInTests {

			fmt.Println("GetMembers->", errGetMembers.Error())
		}
	}

	//filterExp2 := fmt.Sprintf("custom.a5 == '%s' || custom.c5 == '%s'", custom5["a5"], custom5["c5"])
	filterExp2 := fmt.Sprintf("uuid.name == '%s'", name)
	//fmt.Println("GetMembers ====>", filterExp2)

	resGetMembersF, stGetMembersF, errGetMembersF := pn.GetChannelMembers().Channel(spaceid2).Include(inclSm).Filter(filterExp2).Limit(limit).Count(count).Execute()
	//fmt.Println("resGetMembers -->", resGetMembersF)
	assert.Nil(errGetMembersF)
	assert.Equal(200, stGetMembersF.StatusCode)
	if errGetMembersF == nil {
		foundGetMembersF := false

		for i := range resGetMembersF.Data {
			if resGetMembersF.Data[i].UUID.ID == userid2 {
				foundGetMembersF = true
				assert.Equal(name, resGetMembersF.Data[i].UUID.Name)
				assert.Equal(extid, resGetMembersF.Data[i].UUID.ExternalID)
				assert.Equal(purl, resGetMembersF.Data[i].UUID.ProfileURL)
				assert.Equal(email, resGetMembersF.Data[i].UUID.Email)
				assert.Equal(custom["a"], resGetMembersF.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resGetMembersF.Data[i].UUID.Custom["c"])
				assert.Equal(custom5["a5"], resGetMembersF.Data[i].Custom["a5"])
				assert.Equal(custom5["c5"], resGetMembersF.Data[i].Custom["c5"])
			}
		}
		assert.True(foundGetMembersF)
	} else {
		if enableDebuggingInTests {

			fmt.Println("GetMembers->", errGetMembersF.Error())
		}
	}

	// //Remove user memberships

	reMem := pubnub.PNMembershipsRemove{
		Channel: channel2,
	}

	reArrMem := []pubnub.PNMembershipsRemove{
		reMem,
	}
	resManageMemRem, stManageMemRem, errManageMemRem := pn.RemoveMemberships().UUID(userid2).Sort(sort).Remove(reArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	assert.Nil(errManageMemRem)
	assert.Equal(200, stManageMemRem.StatusCode)
	if errManageMemRem == nil {

		foundManageMemRem := false
		for i := range resManageMemRem.Data {
			if resManageMemRem.Data[i].Channel.ID == spaceid2 {
				foundManageMemRem = true
			}
		}
		assert.False(foundManageMemRem)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMemberships->", errManageMemRem.Error())
		}
	}

	//delete
	res5, st5, err5 := pn.RemoveUUIDMetadata().UUID(userid).Execute()
	assert.Nil(err5)
	assert.Equal(200, st5.StatusCode)

	assert.Nil(res5.Data)

	//delete
	res6, st6, err6 := pn.RemoveChannelMetadata().Channel(spaceid).Execute()
	assert.Nil(err6)
	assert.Equal(200, st6.StatusCode)
	assert.Nil(res6.Data)

	//delete
	res52, st52, err52 := pn.RemoveUUIDMetadata().UUID(userid2).Execute()
	assert.Nil(err52)
	assert.Equal(200, st52.StatusCode)
	if res52 != nil {
		assert.Nil(res52.Data)
	}

	//delete
	res62, st62, err62 := pn.RemoveChannelMetadata().Channel(spaceid2).Execute()
	assert.Nil(err62)
	assert.Equal(200, st62.StatusCode)
	if res62 != nil {
		assert.Nil(res62.Data)
	}

}

func TestObjectsV2MembershipsV2(t *testing.T) {
	ObjectsMembershipsCommonV2(t, false, false)
}

func TestObjectsV2MembershipsV2WithPAM(t *testing.T) {
	ObjectsMembershipsCommonV2(t, true, false)
}

// PASSES after adding PAM checks for Update Members
// func TestObjectsV2MembershipsV2WithPAMWithoutSecKey(t *testing.T) {
// 	ObjectsMembershipsCommonV2(t, true, true)
// }

func ObjectsMembershipsCommonV2(t *testing.T, withPAM, runWithoutSecretKey bool) {
	assert := assert.New(t)

	limit := 100
	count := true

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	userid := fmt.Sprintf("testuser_%d", r.Intn(99999))
	spaceid := fmt.Sprintf("testspace_%d", r.Intn(99999))

	if withPAM {
		pn2 := ActivateWithPAMV2()
		if runWithoutSecretKey {
			tokens := RunGrantV2(pn2, []string{userid}, []string{spaceid}, true, true, true, true, true, true)
			SetPNV2(pn, pn2, tokens)
		} else {
			pn = pn2
			RunGrantV2(pn, []string{userid}, []string{spaceid}, true, true, true, true, true, false)
		}

	}
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	name := fmt.Sprintf("name_%d", r.Intn(99999))
	extid := "extid"
	purl := "profileurl"
	email := "email"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	inclUUID := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	res, st, err := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, st.StatusCode)
	if res != nil {
		assert.Equal(userid, res.Data.ID)
		assert.Equal(name, res.Data.Name)
		assert.Equal(extid, res.Data.ExternalID)
		assert.Equal(purl, res.Data.ProfileURL)
		assert.Equal(email, res.Data.Email)
		// assert.NotNil(res.Data.Created)
		assert.NotNil(res.Data.Updated)
		assert.NotNil(res.Data.ETag)
		assert.Equal("b", res.Data.Custom["a"])
		assert.Equal("d", res.Data.Custom["c"])
	}

	desc := "desc"
	custom2 := make(map[string]interface{})
	custom2["a1"] = "b1"
	custom2["c1"] = "d1"

	inclChannel := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	res2, st2, err2 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err2)
	assert.Equal(200, st2.StatusCode)
	//fmt.Println("res2-->", res2)
	if res2 != nil {
		assert.Equal(spaceid, res2.Data.ID)
		assert.Equal(name, res2.Data.Name)
		assert.Equal(desc, res2.Data.Description)
		// assert.NotNil(res2.Data.Created)
		assert.NotNil(res2.Data.Updated)
		assert.NotNil(res2.Data.ETag)
		assert.Equal("b1", res2.Data.Custom["a1"])
		assert.Equal("d1", res2.Data.Custom["c1"])
	}

	userid2 := fmt.Sprintf("testuser_%d", r.Intn(99999))

	_, st3, err3 := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid2).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err3)
	assert.Equal(200, st3.StatusCode)

	spaceid2 := fmt.Sprintf("testspace_%d", r.Intn(99999))

	_, st4, err4 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid2).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err4)
	assert.Equal(200, st4.StatusCode)

	userid3 := fmt.Sprintf("testuser_%d", r.Intn(99999))

	_, stuser3, erruser3 := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid3).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(erruser3)
	assert.Equal(200, stuser3.StatusCode)

	spaceid3 := fmt.Sprintf("testspace_%d", r.Intn(99999))

	_, stspace3, errspace3 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid3).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(errspace3)
	assert.Equal(200, stspace3.StatusCode)

	inclSm := []pubnub.PNChannelMembersInclude{
		pubnub.PNChannelMembersIncludeUUIDCustom,
		pubnub.PNChannelMembersIncludeUUID,
		pubnub.PNChannelMembersIncludeCustom,
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	uuid := pubnub.PNChannelMembersUUID{
		ID: userid,
	}

	in := pubnub.PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom3,
	}
	uuid3 := pubnub.PNChannelMembersUUID{
		ID: userid3,
	}

	inUser3 := pubnub.PNChannelMembersSet{
		UUID:   uuid3,
		Custom: custom3,
	}

	inArr := []pubnub.PNChannelMembersSet{
		in,
		inUser3,
	}

	//Add Space Memberships
	sort := []string{"updated:desc"}

	resAdd, stAdd, errAdd := pn.ManageChannelMembers().Channel(spaceid).Sort(sort).Set(inArr).Remove([]pubnub.PNChannelMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errAdd)
	assert.Equal(200, stAdd.StatusCode)
	if errAdd == nil {
		sortMembers1 := false
		sortMembers2 := false

		found := false
		assert.True(resAdd.TotalCount > 0)
		//fmt.Println("resAdd-->", resAdd)

		for i := range resAdd.Data {
			if resAdd.Data[i].UUID.ID == userid {
				found = true
				assert.Equal(custom3["a3"], resAdd.Data[i].Custom["a3"])
				assert.Equal(custom3["c3"], resAdd.Data[i].Custom["c3"])
				assert.Equal(userid, resAdd.Data[i].UUID.ID)
				assert.Equal(name, resAdd.Data[i].UUID.Name)
				assert.Equal(extid, resAdd.Data[i].UUID.ExternalID)
				assert.Equal(purl, resAdd.Data[i].UUID.ProfileURL)
				assert.Equal(email, resAdd.Data[i].UUID.Email)
				assert.Equal(custom["a"], resAdd.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resAdd.Data[i].UUID.Custom["c"])
			}
		}
		if (resAdd.Data != nil) && (len(resAdd.Data) > 1) {
			sortMembers1 = (resAdd.Data[1].UUID.ID == userid)
			sortMembers2 = (resAdd.Data[0].UUID.ID == userid3)
			assert.True(sortMembers1)
			assert.True(sortMembers2)
		} else {
			assert.Fail("Sort ", "resAdd.Data null or ", len(resAdd.Data))
		}

		assert.True(found)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMembers->", errAdd.Error())
		}
	}

	//Update Space Memberships
	if !withPAM {

		custom4 := make(map[string]interface{})
		custom4["a2"] = "b2"
		custom4["c2"] = "d2"

		up := pubnub.PNChannelMembersSet{
			UUID:   uuid,
			Custom: custom4,
		}

		upArr := []pubnub.PNChannelMembersSet{
			up,
		}

		resUp, stUp, errUp := pn.ManageChannelMembers().Channel(spaceid).Sort(sort).Set(upArr).Remove([]pubnub.PNChannelMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
		assert.Nil(errUp)
		assert.Equal(200, stUp.StatusCode)
		if errUp == nil {
			assert.True(resUp.TotalCount > 0)
			foundUp := false
			for i := range resUp.Data {
				if resUp.Data[i].UUID.ID == userid {
					foundUp = true
					assert.Equal("b2", resUp.Data[i].Custom["a2"])
					assert.Equal("d2", resUp.Data[i].Custom["c2"])
					//assert.Equal(userid, resAdd.Data[i].UUID.ID)
					assert.Equal(name, resAdd.Data[i].UUID.Name)
					assert.Equal(extid, resAdd.Data[i].UUID.ExternalID)
					assert.Equal(purl, resAdd.Data[i].UUID.ProfileURL)
					assert.Equal(email, resAdd.Data[i].UUID.Email)
					assert.Equal(custom["a"], resAdd.Data[i].UUID.Custom["a"])
					assert.Equal(custom["c"], resAdd.Data[i].UUID.Custom["c"])

				}
			}
			assert.True(foundUp)
		} else {
			if enableDebuggingInTests {

				fmt.Println("ManageMembers->", errUp.Error())
			}
		}
	}
	//Get Space Memberships

	inclMemberships := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsIncludeCustom,
		pubnub.PNMembershipsIncludeChannel,
		pubnub.PNMembershipsIncludeChannelCustom,
	}

	//fmt.Println("GetMemberships ====>")

	resGetMem, stGetMem, errGetMem := pn.GetMemberships().UUID(userid).Include(inclMemberships).Sort(sort).Limit(limit).Count(count).Execute()
	foundGetMem := false
	assert.Nil(errGetMem)
	if errGetMem == nil {
		for i := range resGetMem.Data {
			if resGetMem.Data[i].Channel.ID == spaceid {
				foundGetMem = true
				assert.Equal(name, resGetMem.Data[i].Channel.Name)
				assert.Equal(desc, resGetMem.Data[i].Channel.Description)
				assert.Equal("b1", resGetMem.Data[i].Channel.Custom["a1"])
				assert.Equal("d1", resGetMem.Data[i].Channel.Custom["c1"])
				if withPAM {
					assert.Equal("b3", resGetMem.Data[i].Custom["a3"])
					assert.Equal("d3", resGetMem.Data[i].Custom["c3"])
				} else {
					assert.Equal("b2", resGetMem.Data[i].Custom["a2"])
					assert.Equal("d2", resGetMem.Data[i].Custom["c2"])
				}
			}
		}
		assert.Equal(200, stGetMem.StatusCode)
		assert.True(foundGetMem)
	} else {
		if enableDebuggingInTests {
			fmt.Println("GetMemberships->", errGetMem.Error())
		}
	}

	//filterExp := fmt.Sprintf("custom.c3 == '%s' || custom.c2 == '%s'", "d3", "d2")
	filterExp := fmt.Sprintf("channel.name == '%s'", name)

	//fmt.Println("GetMemberships ====>", filterExp)

	resGetMemF, stGetMemF, errGetMemF := pn.GetMemberships().UUID(userid).Include(inclMemberships).Filter(filterExp).Limit(limit).Count(count).Execute()
	foundGetMemF := false
	assert.Nil(errGetMemF)
	if errGetMemF == nil {
		for i := range resGetMemF.Data {
			if resGetMemF.Data[i].Channel.ID == spaceid {
				foundGetMemF = true
				assert.Equal(name, resGetMemF.Data[i].Channel.Name)
				assert.Equal(desc, resGetMemF.Data[i].Channel.Description)
				assert.Equal("b1", resGetMemF.Data[i].Channel.Custom["a1"])
				assert.Equal("d1", resGetMemF.Data[i].Channel.Custom["c1"])
				if withPAM {
					assert.Equal("b3", resGetMemF.Data[i].Custom["a3"])
					assert.Equal("d3", resGetMemF.Data[i].Custom["c3"])
				} else {
					assert.Equal("b2", resGetMemF.Data[i].Custom["a2"])
					assert.Equal("d2", resGetMemF.Data[i].Custom["c2"])
				}
			}
		}
		assert.Equal(200, stGetMemF.StatusCode)
		assert.True(foundGetMemF)
	} else {
		if enableDebuggingInTests {

			fmt.Println("GetMemberships->", errGetMemF.Error())
		}
	}

	//Remove Space Memberships
	re := pubnub.PNChannelMembersRemove{
		UUID: uuid,
	}

	reArr := []pubnub.PNChannelMembersRemove{
		re,
	}
	resRem, stRem, errRem := pn.ManageChannelMembers().Channel(spaceid).Set([]pubnub.PNChannelMembersSet{}).Remove(reArr).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errRem)
	assert.Equal(200, stRem.StatusCode)
	if errRem == nil {

		foundRem := false
		for i := range resRem.Data {
			if resRem.Data[i].UUID.ID == userid {
				foundRem = true
				assert.Equal("b2", resRem.Data[i].Custom["a2"])
				assert.Equal("d2", resRem.Data[i].Custom["c2"])
				assert.Equal(userid, resRem.Data[i].UUID.ID)
				assert.Equal(name, resRem.Data[i].UUID.Name)
				assert.Equal(extid, resRem.Data[i].UUID.ExternalID)
				assert.Equal(purl, resRem.Data[i].UUID.ProfileURL)
				assert.Equal(email, resRem.Data[i].UUID.Email)
				assert.Equal(custom["a"], resRem.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resRem.Data[i].UUID.Custom["c"])

			}
		}
		assert.False(foundRem)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMembers->", errRem.Error())
		}
	}

	channel2 := pubnub.PNMembershipsChannel{
		ID: spaceid2,
	}

	inMem := pubnub.PNMembershipsSet{
		Channel: channel2,
		Custom:  custom3,
	}

	channel3 := pubnub.PNMembershipsChannel{
		ID: spaceid3,
	}

	inMemSpace3 := pubnub.PNMembershipsSet{
		Channel: channel3,
		Custom:  custom3,
	}

	inArrMem := []pubnub.PNMembershipsSet{
		inMem,
		inMemSpace3,
	}

	//Add user memberships
	resManageMemAdd, stManageMemAdd, errManageMemAdd := pn.ManageMemberships().UUID(userid2).Set(inArrMem).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	//fmt.Println("resManageMemAdd -->", resManageMemAdd)
	assert.Nil(errManageMemAdd)
	assert.Equal(200, stManageMemAdd.StatusCode)
	if errManageMemAdd == nil {
		foundManageMembers := false
		for i := range resManageMemAdd.Data {
			if resManageMemAdd.Data[i].Channel.ID == spaceid2 {
				assert.Equal(spaceid2, resManageMemAdd.Data[i].Channel.ID)
				assert.Equal(name, resManageMemAdd.Data[i].Channel.Name)
				assert.Equal(desc, resManageMemAdd.Data[i].Channel.Description)
				assert.Equal(custom2["a1"], resManageMemAdd.Data[i].Channel.Custom["a1"])
				assert.Equal(custom2["c1"], resManageMemAdd.Data[i].Channel.Custom["c1"])
				assert.Equal(custom3["a3"], resManageMemAdd.Data[i].Custom["a3"])
				assert.Equal(custom3["c3"], resManageMemAdd.Data[i].Custom["c3"])
				foundManageMembers = true
			}
		}
		assert.True(foundManageMembers)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMemberships->", errManageMemAdd.Error())
		}
	}

	// //Update user memberships

	custom5 := make(map[string]interface{})
	custom5["a5"] = "b5"
	custom5["c5"] = "d5"

	upMem := pubnub.PNMembershipsSet{
		Channel: channel2,
		Custom:  custom5,
	}

	upArrMem := []pubnub.PNMembershipsSet{
		upMem,
	}
	sortMemberships1 := false
	sortMemberships2 := false

	resManageMemUp, stManageMemUp, errManageMemUp := pn.ManageMemberships().UUID(userid2).Sort(sort).Set(upArrMem).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	//fmt.Println("resManageMemUp -->", resManageMemUp)
	assert.Nil(errManageMemUp)
	assert.Equal(200, stManageMemUp.StatusCode)
	if errManageMemUp == nil {
		foundManageMembersUp := false

		for i := range resManageMemUp.Data {
			//fmt.Println("resManageMemUp.Data[i].ID == spaceid2-->", resRem.Data[i].UUID.ID, spaceid2)
			if resManageMemUp.Data[i].Channel.ID == spaceid2 {
				assert.Equal(spaceid2, resManageMemUp.Data[i].Channel.ID)
				assert.Equal(name, resManageMemUp.Data[i].Channel.Name)
				assert.Equal(desc, resManageMemUp.Data[i].Channel.Description)
				assert.Equal(custom2["a1"], resManageMemAdd.Data[i].Channel.Custom["a1"])
				assert.Equal(custom2["c1"], resManageMemAdd.Data[i].Channel.Custom["c1"])
				assert.Equal(custom5["a5"], resManageMemUp.Data[i].Custom["a5"])
				assert.Equal(custom5["c5"], resManageMemUp.Data[i].Custom["c5"])
				foundManageMembersUp = true
			}
		}
		if (resManageMemUp.Data != nil) && (len(resManageMemUp.Data) > 1) {
			sortMemberships1 = (resManageMemUp.Data[0].Channel.ID == spaceid2)
			sortMemberships2 = (resManageMemUp.Data[1].Channel.ID == spaceid3)
			assert.True(sortMemberships1)
			assert.True(sortMemberships2)
		} else {
			assert.Fail("Sort ", "resManageMemUp.Data null or ", len(resManageMemUp.Data))
		}

		assert.True(foundManageMembersUp)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMemberships->", errManageMemUp.Error())
		}
	}

	// //Get members
	resGetMembers, stGetMembers, errGetMembers := pn.GetChannelMembers().Channel(spaceid2).Include(inclSm).Limit(limit).Count(count).Execute()
	//fmt.Println("resGetMembers -->", resGetMembers)
	assert.Nil(errGetMembers)
	assert.Equal(200, stGetMembers.StatusCode)
	if errGetMembers == nil {
		foundGetMembers := false
		for i := range resGetMembers.Data {
			if resGetMembers.Data[i].UUID.ID == userid2 {
				foundGetMembers = true
				assert.Equal(name, resGetMembers.Data[i].UUID.Name)
				assert.Equal(extid, resGetMembers.Data[i].UUID.ExternalID)
				assert.Equal(purl, resGetMembers.Data[i].UUID.ProfileURL)
				assert.Equal(email, resGetMembers.Data[i].UUID.Email)
				assert.Equal(custom["a"], resGetMembers.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resGetMembers.Data[i].UUID.Custom["c"])
				assert.Equal(custom5["a5"], resGetMembers.Data[i].Custom["a5"])
				assert.Equal(custom5["c5"], resGetMembers.Data[i].Custom["c5"])
			}
		}

		assert.True(foundGetMembers)
	} else {
		if enableDebuggingInTests {

			fmt.Println("GetMembers->", errGetMembers.Error())
		}
	}

	//filterExp2 := fmt.Sprintf("custom.a5 == '%s' || custom.c5 == '%s'", custom5["a5"], custom5["c5"])
	filterExp2 := fmt.Sprintf("uuid.name == '%s'", name)
	//fmt.Println("GetMembers ====>", filterExp2)

	resGetMembersF, stGetMembersF, errGetMembersF := pn.GetChannelMembers().Channel(spaceid2).Include(inclSm).Filter(filterExp2).Limit(limit).Count(count).Execute()
	//fmt.Println("resGetMembers -->", resGetMembersF)
	assert.Nil(errGetMembersF)
	assert.Equal(200, stGetMembersF.StatusCode)
	if errGetMembersF == nil {
		foundGetMembersF := false

		for i := range resGetMembersF.Data {
			if resGetMembersF.Data[i].UUID.ID == userid2 {
				foundGetMembersF = true
				assert.Equal(name, resGetMembersF.Data[i].UUID.Name)
				assert.Equal(extid, resGetMembersF.Data[i].UUID.ExternalID)
				assert.Equal(purl, resGetMembersF.Data[i].UUID.ProfileURL)
				assert.Equal(email, resGetMembersF.Data[i].UUID.Email)
				assert.Equal(custom["a"], resGetMembersF.Data[i].UUID.Custom["a"])
				assert.Equal(custom["c"], resGetMembersF.Data[i].UUID.Custom["c"])
				assert.Equal(custom5["a5"], resGetMembersF.Data[i].Custom["a5"])
				assert.Equal(custom5["c5"], resGetMembersF.Data[i].Custom["c5"])
			}
		}
		assert.True(foundGetMembersF)
	} else {
		if enableDebuggingInTests {

			fmt.Println("GetMembers->", errGetMembersF.Error())
		}
	}

	// //Remove user memberships

	reMem := pubnub.PNMembershipsRemove{
		Channel: channel2,
	}

	reArrMem := []pubnub.PNMembershipsRemove{
		reMem,
	}
	resManageMemRem, stManageMemRem, errManageMemRem := pn.ManageMemberships().UUID(userid2).Sort(sort).Set([]pubnub.PNMembershipsSet{}).Remove(reArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	assert.Nil(errManageMemRem)
	assert.Equal(200, stManageMemRem.StatusCode)
	if errManageMemRem == nil {

		foundManageMemRem := false
		for i := range resManageMemRem.Data {
			if resManageMemRem.Data[i].Channel.ID == spaceid2 {
				foundManageMemRem = true
			}
		}
		assert.False(foundManageMemRem)
	} else {
		if enableDebuggingInTests {

			fmt.Println("ManageMemberships->", errManageMemRem.Error())
		}
	}

	//delete
	res5, st5, err5 := pn.RemoveUUIDMetadata().UUID(userid).Execute()
	assert.Nil(err5)
	assert.Equal(200, st5.StatusCode)

	assert.Nil(res5.Data)

	//delete
	res6, st6, err6 := pn.RemoveChannelMetadata().Channel(spaceid).Execute()
	assert.Nil(err6)
	assert.Equal(200, st6.StatusCode)
	assert.Nil(res6.Data)

	//delete
	res52, st52, err52 := pn.RemoveUUIDMetadata().UUID(userid2).Execute()
	assert.Nil(err52)
	assert.Equal(200, st52.StatusCode)
	if res52 != nil {
		assert.Nil(res52.Data)
	}

	//delete
	res62, st62, err62 := pn.RemoveChannelMetadata().Channel(spaceid2).Execute()
	assert.Nil(err62)
	assert.Equal(200, st62.StatusCode)
	if res62 != nil {
		assert.Nil(res62.Data)
	}

}

func TestObjectsV2ListenersV2(t *testing.T) {
	ObjectsListenersCommonV2(t, false, false)
}

func TestObjectsV2ListenersV2WithPAM(t *testing.T) {
	ObjectsListenersCommonV2(t, true, false)
}

// func TestObjectsV2ListenersV2WithPAMWithoutSecKey(t *testing.T) {
// 	ObjectsListenersCommonV2(t, true, true)
// }

func ObjectsListenersCommonV2(t *testing.T, withPAM, runWithoutSecretKey bool) {
	//Create channel names for Space and User
	eventWaitTime := 2
	assert := assert.New(t)

	limit := 100
	count := true

	pn := pubnub.NewPubNub(configCopy())
	pnSub := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	userid := fmt.Sprintf("testlistuser_%d", r.Intn(99999))
	spaceid := fmt.Sprintf("testlistspace_%d", r.Intn(99999))
	if withPAM {
		pn2 := ActivateWithPAMV2()
		if runWithoutSecretKey {
			if enableDebuggingInTests {
				pn2.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
			}
			tokens := RunGrantV2(pn2, []string{userid}, []string{spaceid}, true, true, true, true, true, true)
			SetPNV2(pn, pn2, tokens)
			SetPNV2(pnSub, pn2, tokens)
			//You have to use Grant v2 to subscribe
			pnSub.Config.AuthKey = "authKey"
			pn2.Grant().
				Read(true).Write(true).Manage(true).
				Channels([]string{userid, spaceid}).
				AuthKeys([]string{pnSub.Config.AuthKey}).
				Execute()
		} else {
			pn = pn2
			pnSub = pn2
			if enableDebuggingInTests {
				pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
			}
			RunGrantV2(pn, []string{userid}, []string{spaceid}, true, true, true, true, true, false)
		}
	}
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
		pnSub.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	//Subscribe to the channel names

	listener := pubnub.NewListener()

	var mut sync.RWMutex

	addUserToSpace := false
	addUserToSpace2 := false
	updateUserMem := false
	updateUser := false
	updateSpace := false
	removeUserFromSpace := false
	deleteUser := false
	deleteSpace := false

	doneConnected := make(chan bool)
	exitListener := make(chan bool)

	go func() {
	ExitLabel:
		for {
			//fmt.Println("Running =--->")
			select {

			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnected <- true
				default:
					if enableDebuggingInTests {

						fmt.Println(" --- status: ", status)
					}
				}

			case userEvent := <-listener.UUIDEvent:
				if enableDebuggingInTests {

					fmt.Println(" --- UserEvent: ")
					fmt.Println(fmt.Sprintf("%s", userEvent))
					fmt.Println(fmt.Sprintf("userEvent.Channel: %s", userEvent.Channel))
					fmt.Println(fmt.Sprintf("userEvent.SubscribedChannel: %s", userEvent.SubscribedChannel))
					fmt.Println(fmt.Sprintf("userEvent.Event: %s", userEvent.Event))
					fmt.Println(fmt.Sprintf("userEvent.UUID: %s", userEvent.UUID))
					fmt.Println(fmt.Sprintf("userEvent.Description: %s", userEvent.Description))
					fmt.Println(fmt.Sprintf("userEvent.Timestamp: %s", userEvent.Timestamp))
					fmt.Println(fmt.Sprintf("userEvent.Name: %s", userEvent.Name))
					fmt.Println(fmt.Sprintf("userEvent.ExternalID: %s", userEvent.ExternalID))
					fmt.Println(fmt.Sprintf("userEvent.ProfileURL: %s", userEvent.ProfileURL))
					fmt.Println(fmt.Sprintf("userEvent.Email: %s", userEvent.Email))
					// fmt.Println(fmt.Sprintf("userEvent.Created: %s", userEvent.Created))
					fmt.Println(fmt.Sprintf("userEvent.Updated: %s", userEvent.Updated))
					fmt.Println(fmt.Sprintf("userEvent.ETag: %s", userEvent.ETag))
					fmt.Println(fmt.Sprintf("userEvent.Custom: %v", userEvent.Custom))
				}

				if (userEvent.Event == pubnub.PNObjectsEventRemove) && (userEvent.UUID == userid) {
					mut.Lock()
					deleteUser = true
					mut.Unlock()
				}
				if (userEvent.Event == pubnub.PNObjectsEventSet) && (userEvent.UUID == userid) {
					mut.Lock()
					updateUser = true
					mut.Unlock()
				}
			case spaceEvent := <-listener.ChannelEvent:

				if enableDebuggingInTests {

					fmt.Println(" --- SpaceEvent: ")
					fmt.Println(fmt.Sprintf("%s", spaceEvent))
					fmt.Println(fmt.Sprintf("spaceEvent.SubscribedChannel: %s", spaceEvent.SubscribedChannel))
					fmt.Println(fmt.Sprintf("spaceEvent.Event: %s", spaceEvent.Event))
					fmt.Println(fmt.Sprintf("spaceEvent.ChannelID: %s", spaceEvent.ChannelID))
					fmt.Println(fmt.Sprintf("spaceEvent.Channel: %s", spaceEvent.Channel))
					fmt.Println(fmt.Sprintf("spaceEvent.Description: %s", spaceEvent.Description))
					fmt.Println(fmt.Sprintf("spaceEvent.Timestamp: %s", spaceEvent.Timestamp))
					// fmt.Println(fmt.Sprintf("spaceEvent.Created: %s", spaceEvent.Created))
					fmt.Println(fmt.Sprintf("spaceEvent.Updated: %s", spaceEvent.Updated))
					fmt.Println(fmt.Sprintf("spaceEvent.ETag: %s", spaceEvent.ETag))
					fmt.Println(fmt.Sprintf("spaceEvent.Custom: %v", spaceEvent.Custom))
				}
				if (spaceEvent.Event == pubnub.PNObjectsEventRemove) && (spaceEvent.ChannelID == spaceid) {
					mut.Lock()
					deleteSpace = true
					mut.Unlock()
				}
				if (spaceEvent.Event == pubnub.PNObjectsEventSet) && (spaceEvent.ChannelID == spaceid) {
					mut.Lock()
					updateSpace = true
					mut.Unlock()
				}

			case membershipEvent := <-listener.MembershipEvent:
				if enableDebuggingInTests {

					fmt.Println(" --- MembershipEvent: ")
					fmt.Println(fmt.Sprintf("%s", membershipEvent))
					fmt.Println(fmt.Sprintf("membershipEvent.SubscribedChannel: %s", membershipEvent.SubscribedChannel))
					fmt.Println(fmt.Sprintf("membershipEvent.Event: %s", membershipEvent.Event))
					fmt.Println(fmt.Sprintf("membershipEvent.Channel: %s", membershipEvent.Channel))
					fmt.Println(fmt.Sprintf("membershipEvent.UUID: %s", membershipEvent.UUID))
					fmt.Println(fmt.Sprintf("membershipEvent.ChannelID: %s", membershipEvent.ChannelID))
					fmt.Println(fmt.Sprintf("membershipEvent.Description: %s", membershipEvent.Description))
					fmt.Println(fmt.Sprintf("membershipEvent.Timestamp: %s", membershipEvent.Timestamp))
					fmt.Println(fmt.Sprintf("membershipEvent.Custom: %v", membershipEvent.Custom))
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventSet) && (membershipEvent.ChannelID == spaceid) && (membershipEvent.UUID == userid) && ((membershipEvent.Channel == spaceid) || (membershipEvent.Channel == userid)) {
					mut.Lock()
					addUserToSpace = true
					mut.Unlock()
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventSet) && (membershipEvent.ChannelID == spaceid) && (membershipEvent.UUID == userid) && ((membershipEvent.Channel == spaceid) || (membershipEvent.Channel == userid)) {
					mut.Lock()
					addUserToSpace2 = true
					mut.Unlock()
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventSet) && (membershipEvent.ChannelID == spaceid) && (membershipEvent.UUID == userid) && ((membershipEvent.Channel == spaceid) || (membershipEvent.Channel == userid)) {
					mut.Lock()
					updateUserMem = true
					mut.Unlock()
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventSet) && (membershipEvent.ChannelID == spaceid) && (membershipEvent.UUID == userid) && ((membershipEvent.Channel == spaceid) || (membershipEvent.Channel == userid)) {
					mut.Lock()
					updateUserMem = true
					mut.Unlock()
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventRemove) && (membershipEvent.ChannelID == spaceid) && (membershipEvent.UUID == userid) && ((membershipEvent.Channel == spaceid) || (membershipEvent.Channel == userid)) {
					mut.Lock()
					removeUserFromSpace = true
					mut.Unlock()
				}
			case <-exitListener:
				break ExitLabel

			}

			//fmt.Println("=>>>>>>>>>>>>> restart")

		}

	}()

	pnSub.AddListener(listener)

	pnSub.Subscribe().Channels([]string{userid, spaceid}).Execute()
	tic := time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneConnected:
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	name := "name"
	extid := "extid"
	purl := "profileurl"
	email := "email"
	desc := "desc"

	customUser := make(map[string]interface{})
	customUser["au"] = "bu"
	customUser["cu"] = "du"

	inclUUID := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	//Create User
	_, st, err := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(customUser).Execute()
	assert.Nil(err)
	assert.Equal(200, st.StatusCode)

	//Create Space
	customSpace := make(map[string]interface{})
	customSpace["as"] = "bs"
	customSpace["cs"] = "ds"

	inclChannel := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	_, st4, err4 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid).Name(name).Description(desc).Custom(customSpace).Execute()
	assert.Nil(err4)
	assert.Equal(200, st4.StatusCode)

	time.Sleep(1 * time.Second)

	//Update User
	email = "email2"
	//fmt.Println("SetUUIDMetadata, Update ===> ", userid)

	_, st2, err2 := pn.SetUUIDMetadata().Include(inclUUID).UUID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(customUser).Execute()
	assert.Nil(err2)
	assert.Equal(200, st2.StatusCode)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(updateUser)
	mut.Unlock()

	desc = "desc2"

	//fmt.Println("SetChannelMetadata, Update ===> ", spaceid)

	//Update Space
	_, st3, err3 := pn.SetChannelMetadata().Include(inclChannel).Channel(spaceid).Name(name).Description(desc).Custom(customSpace).Execute()
	assert.Nil(err3)
	assert.Equal(200, st3.StatusCode)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(updateSpace)
	mut.Unlock()

	//Add user to space
	inclSm := []pubnub.PNChannelMembersInclude{
		pubnub.PNChannelMembersIncludeCustom,
		pubnub.PNChannelMembersIncludeUUID,
		pubnub.PNChannelMembersIncludeUUIDCustom,
	}

	if enableDebuggingInTests {

		fmt.Println("inclSm===>", inclSm)
		for k, value := range inclSm {
			fmt.Println("inclSm===>", k, value)
		}
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	uuid := pubnub.PNChannelMembersUUID{
		ID: userid,
	}

	in := pubnub.PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom3,
	}

	inArr := []pubnub.PNChannelMembersSet{
		in,
	}

	_, stAdd, errAdd := pn.ManageChannelMembers().Channel(spaceid).Set(inArr).Remove([]pubnub.PNChannelMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errAdd)
	if enableDebuggingInTests {

		if errAdd != nil {
			fmt.Println("ManageMembers-->", errAdd)
		}
	}
	assert.Equal(200, stAdd.StatusCode)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(addUserToSpace && addUserToSpace2)
	mut.Unlock()

	//Update user membership

	//Read event

	custom5 := make(map[string]interface{})
	custom5["a5"] = "b5"
	custom5["c5"] = "d5"

	channel := pubnub.PNMembershipsChannel{
		ID: spaceid,
	}

	upMem := pubnub.PNMembershipsSet{
		Channel: channel,
		Custom:  custom5,
	}

	upArrMem := []pubnub.PNMembershipsSet{
		upMem,
	}

	inclMemberships := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsIncludeCustom,
		pubnub.PNMembershipsIncludeChannel,
		pubnub.PNMembershipsIncludeChannelCustom,
	}

	resManageMemUp, stManageMemUp, errManageMemUp := pn.ManageMemberships().UUID(userid).Set(upArrMem).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()

	assert.Nil(errManageMemUp)
	if enableDebuggingInTests {

		fmt.Println("resManageMemUp -->", resManageMemUp)
		if errManageMemUp != nil {
			fmt.Println("ManageMemberships-->", errManageMemUp)
		}
	}
	assert.Equal(200, stManageMemUp.StatusCode)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(updateUserMem)
	mut.Unlock()

	//Remove user from space
	reMem := pubnub.PNMembershipsRemove{
		Channel: channel,
	}

	reArrMem := []pubnub.PNMembershipsRemove{
		reMem,
	}
	_, stManageMemRem, errManageMemRem := pn.ManageMemberships().UUID(userid).Set([]pubnub.PNMembershipsSet{}).Remove(reArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	assert.Nil(errManageMemRem)
	if enableDebuggingInTests {

		if errManageMemRem != nil {
			fmt.Println("ManageMemberships-->", errManageMemRem)
		}
	}
	assert.Equal(200, stManageMemRem.StatusCode)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(removeUserFromSpace)
	mut.Unlock()

	//Delete user
	res52, st52, err52 := pn.RemoveUUIDMetadata().UUID(userid).Execute()
	assert.Nil(err52)
	assert.Equal(200, st52.StatusCode)
	assert.Nil(res52.Data)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(deleteUser)
	mut.Unlock()

	//Delete Space
	res62, st62, err62 := pn.RemoveChannelMetadata().Channel(spaceid).Execute()
	assert.Nil(err62)
	assert.Equal(200, st62.StatusCode)
	assert.Nil(res62.Data)

	time.Sleep(1 * time.Second)
	mut.Lock()
	assert.True(deleteSpace)
	mut.Unlock()

	exitListener <- true
}
