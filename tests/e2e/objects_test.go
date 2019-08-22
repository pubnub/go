package e2e

import (
	//"log"
	//"os"
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/stretchr/testify/assert"
)

func TestObjectsCreateUpdateGetDeleteUser(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	id := fmt.Sprintf("testuser_%d", r.Intn(99999))
	name := "name"
	extid := "extid"
	purl := "profileurl"
	email := "email"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, _, err := pn.CreateUser().Include(incl).Id(id).Name(name).ExternalId(extid).ProfileUrl(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)
	assert.Equal(id, res.Data.Id)
	assert.Equal(name, res.Data.Name)
	assert.Equal(extid, res.Data.ExternalId)
	assert.Equal(purl, res.Data.ProfileUrl)
	assert.Equal(email, res.Data.Email)
	assert.NotNil(res.Data.Created)
	assert.NotNil(res.Data.Updated)
	assert.NotNil(res.Data.ETag)
	assert.Equal("b", res.Data.Custom["a"])
	assert.Equal("d", res.Data.Custom["c"])

	email = "email2"

	res2, _, err2 := pn.UpdateUser().Include(incl).Id(id).Name(name).ExternalId(extid).ProfileUrl(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)
	assert.Equal(id, res2.Data.Id)
	assert.Equal(name, res2.Data.Name)
	assert.Equal(extid, res2.Data.ExternalId)
	assert.Equal(purl, res2.Data.ProfileUrl)
	assert.Equal(email, res2.Data.Email)
	assert.Equal(res.Data.Created, res2.Data.Created)
	assert.NotNil(res2.Data.Updated)
	assert.NotNil(res2.Data.ETag)
	assert.Equal("b", res2.Data.Custom["a"])
	assert.Equal("d", res2.Data.Custom["c"])

	res3, _, err3 := pn.GetUser().Include(incl).Id(id).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)
	assert.Equal(id, res3.Data.Id)
	assert.Equal(name, res3.Data.Name)
	assert.Equal(extid, res3.Data.ExternalId)
	assert.Equal(purl, res3.Data.ProfileUrl)
	assert.Equal(email, res3.Data.Email)
	assert.Equal(res.Data.Created, res3.Data.Created)
	assert.Equal(res2.Data.Updated, res3.Data.Updated)
	assert.Equal(res2.Data.ETag, res3.Data.ETag)
	assert.Equal("b", res3.Data.Custom["a"])
	assert.Equal("d", res3.Data.Custom["c"])

	//getusers
	res6, _, err6 := pn.GetUsers().Include(incl).Limit(100).Count(true).Execute()
	assert.Nil(err6)
	assert.Equal(200, res6.Status)
	assert.True(res6.TotalCount > 0)
	found := false
	for i := range res6.Data {
		if res6.Data[i].Id == id {
			assert.Equal(name, res6.Data[i].Name)
			assert.Equal(extid, res6.Data[i].ExternalId)
			assert.Equal(purl, res6.Data[i].ProfileUrl)
			assert.Equal(email, res6.Data[i].Email)
			assert.Equal(res.Data.Created, res6.Data[i].Created)
			assert.Equal(res2.Data.Updated, res6.Data[i].Updated)
			assert.Equal(res2.Data.ETag, res6.Data[i].ETag)
			assert.Equal("b", res6.Data[i].Custom["a"])
			assert.Equal("d", res6.Data[i].Custom["c"])
			found = true
		}
	}
	assert.True(found)

	//delete
	res5, _, err5 := pn.DeleteUser().Id(id).Execute()
	assert.Nil(err5)
	assert.Equal(200, res5.Status)
	assert.Nil(res5.Data)

	//getuser
	res4, _, err4 := pn.GetUser().Include(incl).Id(id).Execute()
	assert.NotNil(err4)
	assert.Nil(res4)

}

func TestObjectsCreateUpdateGetDeleteSpace(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	id := fmt.Sprintf("testspace_%d", r.Intn(99999))
	name := "name"
	desc := "desc"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, _, err := pn.CreateSpace().Include(incl).Id(id).Name(name).Description(desc).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)
	assert.Equal(id, res.Data.Id)
	assert.Equal(name, res.Data.Name)
	assert.Equal(desc, res.Data.Description)
	assert.NotNil(res.Data.Created)
	assert.NotNil(res.Data.Updated)
	assert.NotNil(res.Data.ETag)
	assert.Equal("b", res.Data.Custom["a"])
	assert.Equal("d", res.Data.Custom["c"])

	desc = "desc2"

	res2, _, err2 := pn.UpdateSpace().Include(incl).Id(id).Name(name).Description(desc).Custom(custom).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)
	assert.Equal(id, res2.Data.Id)
	assert.Equal(name, res2.Data.Name)
	assert.Equal(desc, res2.Data.Description)
	assert.Equal(res.Data.Created, res2.Data.Created)
	assert.NotNil(res2.Data.Updated)
	assert.NotNil(res2.Data.ETag)
	assert.Equal("b", res2.Data.Custom["a"])
	assert.Equal("d", res2.Data.Custom["c"])

	res3, _, err3 := pn.GetSpace().Include(incl).Id(id).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)
	assert.Equal(id, res3.Data.Id)
	assert.Equal(name, res3.Data.Name)
	assert.Equal(desc, res3.Data.Description)
	assert.Equal(res.Data.Created, res3.Data.Created)
	assert.Equal(res2.Data.Updated, res3.Data.Updated)
	assert.Equal(res2.Data.ETag, res3.Data.ETag)
	assert.Equal("b", res3.Data.Custom["a"])
	assert.Equal("d", res3.Data.Custom["c"])

	//getusers
	res6, _, err6 := pn.GetSpaces().Include(incl).Limit(100).Count(true).Execute()
	assert.Nil(err6)
	assert.Equal(200, res6.Status)
	assert.True(res6.TotalCount > 0)
	found := false
	for i := range res6.Data {
		if res6.Data[i].Id == id {
			assert.Equal(name, res6.Data[i].Name)
			assert.Equal(desc, res6.Data[i].Description)
			assert.Equal(res.Data.Created, res6.Data[i].Created)
			assert.Equal(res2.Data.Updated, res6.Data[i].Updated)
			assert.Equal(res2.Data.ETag, res6.Data[i].ETag)
			assert.Equal("b", res6.Data[i].Custom["a"])
			assert.Equal("d", res6.Data[i].Custom["c"])
			found = true
		}
	}
	assert.True(found)

	//delete
	res5, _, err5 := pn.DeleteSpace().Id(id).Execute()
	assert.Nil(err5)
	assert.Equal(200, res5.Status)
	assert.Nil(res5.Data)

	//getuser
	res4, _, err4 := pn.GetSpace().Include(incl).Id(id).Execute()
	assert.NotNil(err4)
	assert.Nil(res4)

}

func TestObjectsMemberships(t *testing.T) {
	assert := assert.New(t)

	limit := 100
	count := true

	pn := pubnub.NewPubNub(configCopy())

	r := GenRandom()

	userid := fmt.Sprintf("testuser_%d", r.Intn(99999))

	name := "name"
	extid := "extid"
	purl := "profileurl"
	email := "email"

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, _, err := pn.CreateUser().Include(incl).Id(userid).Name(name).ExternalId(extid).ProfileUrl(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)
	assert.Equal(userid, res.Data.Id)
	assert.Equal(name, res.Data.Name)
	assert.Equal(extid, res.Data.ExternalId)
	assert.Equal(purl, res.Data.ProfileUrl)
	assert.Equal(email, res.Data.Email)
	assert.NotNil(res.Data.Created)
	assert.NotNil(res.Data.Updated)
	assert.NotNil(res.Data.ETag)
	assert.Equal("b", res.Data.Custom["a"])
	assert.Equal("d", res.Data.Custom["c"])

	spaceid := fmt.Sprintf("testspace_%d", r.Intn(99999))
	desc := "desc"
	custom2 := make(map[string]interface{})
	custom2["a1"] = "b1"
	custom2["c1"] = "d1"

	res2, _, err2 := pn.CreateSpace().Include(incl).Id(spaceid).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)
	assert.Equal(spaceid, res2.Data.Id)
	assert.Equal(name, res2.Data.Name)
	assert.Equal(desc, res2.Data.Description)
	assert.NotNil(res2.Data.Created)
	assert.NotNil(res2.Data.Updated)
	assert.NotNil(res2.Data.ETag)
	assert.Equal("b1", res2.Data.Custom["a1"])
	assert.Equal("d1", res2.Data.Custom["c1"])

	userid2 := fmt.Sprintf("testuser_%d", r.Intn(99999))

	res3, _, err3 := pn.CreateUser().Include(incl).Id(userid2).Name(name).ExternalId(extid).ProfileUrl(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)

	spaceid2 := fmt.Sprintf("testspace_%d", r.Intn(99999))

	res4, _, err4 := pn.CreateSpace().Include(incl).Id(spaceid2).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err4)
	assert.Equal(200, res4.Status)

	inclSm := []pubnub.PNMembersInclude{
		pubnub.PNMembersCustom,
		pubnub.PNMembersUser,
		pubnub.PNMembersUserCustom,
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	in := pubnub.PNMembersInput{
		Id:     userid,
		Custom: custom3,
	}

	inArr := []pubnub.PNMembersInput{
		in,
	}

	//Add Space Memberships

	resAdd, _, errAdd := pn.ManageMembers().SpaceId(spaceid).Add(inArr).Update([]pubnub.PNMembersInput{}).Remove([]pubnub.PNMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errAdd)
	assert.Equal(200, resAdd.Status)
	assert.True(resAdd.TotalCount > 0)
	found := false
	for i := range resAdd.Data {
		if resAdd.Data[i].Id == userid {
			found = true
			assert.Equal("b3", resAdd.Data[i].Custom["a3"])
			assert.Equal("d3", resAdd.Data[i].Custom["c3"])
		}
	}
	assert.True(found)

	//Update Space Memberships

	custom4 := make(map[string]interface{})
	custom4["a2"] = "b2"
	custom4["c2"] = "d2"

	up := pubnub.PNMembersInput{
		Id:     userid,
		Custom: custom4,
	}

	upArr := []pubnub.PNMembersInput{
		up,
	}

	resUp, _, errUp := pn.ManageMembers().SpaceId(spaceid).Add([]pubnub.PNMembersInput{}).Update(upArr).Remove([]pubnub.PNMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errUp)
	assert.Equal(200, resUp.Status)
	assert.True(resUp.TotalCount > 0)
	foundUp := false
	for i := range resUp.Data {
		if resUp.Data[i].Id == userid {
			foundUp = true
			assert.Equal("b2", resUp.Data[i].Custom["a2"])
			assert.Equal("d2", resUp.Data[i].Custom["c2"])
		}
	}
	assert.True(foundUp)

	//Get Space Memberships

	inclMemberships := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsCustom,
		pubnub.PNMembershipsSpace,
		pubnub.PNMembershipsSpaceCustom,
	}

	resGetMem, _, errGetMem := pn.GetMemberships().UserId(userid).Include(inclMemberships).Limit(limit).Count(count).Execute()
	foundGetMem := false
	assert.Nil(errGetMem)
	for i := range resGetMem.Data {
		if resGetMem.Data[i].Id == spaceid {
			foundGetMem = true
			assert.Equal(name, resGetMem.Data[i].Space.Name)
			assert.Equal(desc, resGetMem.Data[i].Space.Description)
			assert.Equal("b1", resGetMem.Data[i].Space.Custom["a1"])
			assert.Equal("d1", resGetMem.Data[i].Space.Custom["c1"])
			assert.Equal("b2", resGetMem.Data[i].Custom["a2"])
			assert.Equal("d2", resGetMem.Data[i].Custom["c2"])
		}
	}
	assert.Equal(200, resGetMem.Status)
	assert.True(foundGetMem)

	//Remove Space Memberships
	re := pubnub.PNMembersRemove{
		Id: userid,
	}

	reArr := []pubnub.PNMembersRemove{
		re,
	}
	resRem, _, errRem := pn.ManageMembers().SpaceId(spaceid).Add([]pubnub.PNMembersInput{}).Update([]pubnub.PNMembersInput{}).Remove(reArr).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errRem)
	assert.Equal(200, resRem.Status)
	foundRem := false
	for i := range resRem.Data {
		if resRem.Data[i].Id == userid {
			foundRem = true
		}
	}
	assert.False(foundRem)

	// //Add user memberships
	// res, status, err := pn.ManageMemberships().UserId(userId).Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()
	// //Update user memberships
	// res, status, err := pn.ManageMemberships().UserId(userId).Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()
	// //Get members
	// res, status, err := pn.GetMembers().SpaceId(id).Include(incl).Limit(limit).Count(count).Execute()

	// //Remove user memberships
	// res, status, err := pn.ManageMemberships().UserId(userId).Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()

	//delete
	res5, _, err5 := pn.DeleteUser().Id(userid).Execute()
	assert.Nil(err5)
	assert.Equal(200, res5.Status)
	assert.Nil(res5.Data)

	//delete
	res6, _, err6 := pn.DeleteSpace().Id(spaceid).Execute()
	assert.Nil(err6)
	assert.Equal(200, res6.Status)
	assert.Nil(res6.Data)

	//delete
	res52, _, err52 := pn.DeleteUser().Id(userid2).Execute()
	assert.Nil(err52)
	assert.Equal(200, res52.Status)
	assert.Nil(res52.Data)

	//delete
	res62, _, err62 := pn.DeleteSpace().Id(spaceid2).Execute()
	assert.Nil(err62)
	assert.Equal(200, res62.Status)
	assert.Nil(res62.Data)

}
