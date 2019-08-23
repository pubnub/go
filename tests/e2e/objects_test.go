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
	fmt.Println("resAdd-->", resAdd)
	found := false
	for i := range resAdd.Data {
		if resAdd.Data[i].Id == userid {
			found = true
			assert.Equal(custom3["a3"], resAdd.Data[i].Custom["a3"])
			assert.Equal(custom3["c3"], resAdd.Data[i].Custom["c3"])
			assert.Equal(userid, resAdd.Data[0].User.Id)
			assert.Equal(name, resAdd.Data[0].User.Name)
			assert.Equal(extid, resAdd.Data[0].User.ExternalId)
			assert.Equal(purl, resAdd.Data[0].User.ProfileUrl)
			assert.Equal(email, resAdd.Data[0].User.Email)
			assert.Equal(custom["a"], resAdd.Data[0].User.Custom["a"])
			assert.Equal(custom["c"], resAdd.Data[0].User.Custom["c"])
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
			assert.Equal(userid, resAdd.Data[0].User.Id)
			assert.Equal(name, resAdd.Data[0].User.Name)
			assert.Equal(extid, resAdd.Data[0].User.ExternalId)
			assert.Equal(purl, resAdd.Data[0].User.ProfileUrl)
			assert.Equal(email, resAdd.Data[0].User.Email)
			assert.Equal(custom["a"], resAdd.Data[0].User.Custom["a"])
			assert.Equal(custom["c"], resAdd.Data[0].User.Custom["c"])

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
			assert.Equal("b2", resUp.Data[i].Custom["a2"])
			assert.Equal("d2", resUp.Data[i].Custom["c2"])
			assert.Equal(userid, resUp.Data[0].User.Id)
			assert.Equal(name, resUp.Data[0].User.Name)
			assert.Equal(extid, resUp.Data[0].User.ExternalId)
			assert.Equal(purl, resUp.Data[0].User.ProfileUrl)
			assert.Equal(email, resUp.Data[0].User.Email)
			assert.Equal(custom["a"], resUp.Data[0].User.Custom["a"])
			assert.Equal(custom["c"], resUp.Data[0].User.Custom["c"])

		}
	}
	assert.False(foundRem)

	inMem := pubnub.PNMembershipsInput{
		Id:     spaceid2,
		Custom: custom3,
	}

	inArrMem := []pubnub.PNMembershipsInput{
		inMem,
	}

	// //Add user memberships
	resManageMemAdd, _, errManageMemAdd := pn.ManageMemberships().UserId(userid2).Add(inArrMem).Update([]pubnub.PNMembershipsInput{}).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	fmt.Println("resManageMemAdd -->", resManageMemAdd)
	assert.Nil(errManageMemAdd)
	assert.Equal(200, resManageMemAdd.Status)
	foundManageMembers := false
	for i := range resManageMemAdd.Data {
		if resManageMemAdd.Data[i].Id == spaceid2 {
			assert.Equal(spaceid2, resManageMemAdd.Data[i].Space.Id)
			assert.Equal(name, resManageMemAdd.Data[i].Space.Name)
			assert.Equal(desc, resManageMemAdd.Data[i].Space.Description)
			assert.Equal(custom2["a1"], resManageMemAdd.Data[i].Space.Custom["a1"])
			assert.Equal(custom2["c1"], resManageMemAdd.Data[i].Space.Custom["c1"])
			assert.Equal(custom3["a3"], resManageMemAdd.Data[i].Custom["a3"])
			assert.Equal(custom3["c3"], resManageMemAdd.Data[i].Custom["c3"])
			foundManageMembers = true
		}
	}
	assert.True(foundManageMembers)

	// //Update user memberships

	custom5 := make(map[string]interface{})
	custom5["a5"] = "b5"
	custom5["c5"] = "d5"

	upMem := pubnub.PNMembershipsInput{
		Id:     spaceid2,
		Custom: custom5,
	}

	upArrMem := []pubnub.PNMembershipsInput{
		upMem,
	}

	resManageMemUp, _, errManageMemUp := pn.ManageMemberships().UserId(userid2).Add([]pubnub.PNMembershipsInput{}).Update(upArrMem).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	fmt.Println("resManageMemUp -->", resManageMemUp)
	assert.Nil(errManageMemUp)
	assert.Equal(200, resManageMemUp.Status)
	foundManageMembersUp := false
	for i := range resManageMemUp.Data {
		if resManageMemUp.Data[i].Id == spaceid2 {
			assert.Equal(spaceid2, resManageMemUp.Data[i].Space.Id)
			assert.Equal(name, resManageMemUp.Data[i].Space.Name)
			assert.Equal(desc, resManageMemUp.Data[i].Space.Description)
			assert.Equal(custom2["a1"], resManageMemAdd.Data[i].Space.Custom["a1"])
			assert.Equal(custom2["c1"], resManageMemAdd.Data[i].Space.Custom["c1"])
			assert.Equal(custom5["a5"], resManageMemUp.Data[i].Custom["a5"])
			assert.Equal(custom5["c5"], resManageMemUp.Data[i].Custom["c5"])
			foundManageMembersUp = true
		}
	}
	assert.True(foundManageMembersUp)

	// //Get members
	resGetMembers, _, errGetMembers := pn.GetMembers().SpaceId(spaceid2).Include(inclSm).Limit(limit).Count(count).Execute()
	fmt.Println("resGetMembers -->", resGetMembers)
	assert.Nil(errGetMembers)
	assert.Equal(200, resGetMembers.Status)
	foundGetMembers := false
	for i := range resGetMembers.Data {
		if resGetMembers.Data[i].Id == userid2 {
			foundGetMembers = true
			assert.Equal(name, resGetMembers.Data[i].User.Name)
			assert.Equal(extid, resGetMembers.Data[i].User.ExternalId)
			assert.Equal(purl, resGetMembers.Data[i].User.ProfileUrl)
			assert.Equal(email, resGetMembers.Data[i].User.Email)
			assert.Equal(custom["a"], resGetMembers.Data[i].User.Custom["a"])
			assert.Equal(custom["c"], resGetMembers.Data[i].User.Custom["c"])
			assert.Equal(custom5["a5"], resGetMembers.Data[i].Custom["a5"])
			assert.Equal(custom5["c5"], resGetMembers.Data[i].Custom["c5"])
		}
	}
	assert.True(foundGetMembers)

	// //Remove user memberships

	reMem := pubnub.PNMembershipsRemove{
		Id: spaceid2,
	}

	reArrMem := []pubnub.PNMembershipsRemove{
		reMem,
	}
	resManageMemRem, _, errManageMemRem := pn.ManageMemberships().UserId(userid2).Add([]pubnub.PNMembershipsInput{}).Update([]pubnub.PNMembershipsInput{}).Remove(reArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	assert.Nil(errManageMemRem)
	assert.Equal(200, resManageMemRem.Status)

	foundManageMemRem := false
	for i := range resManageMemRem.Data {
		if resManageMemRem.Data[i].Id == spaceid2 {
			foundManageMemRem = true
		}
	}
	assert.False(foundManageMemRem)

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
