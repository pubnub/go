package e2e

import (
	"fmt"
	//"log"
	//"os"
	"testing"
	"time"

	pubnub "github.com/sprucehealth/pubnub-go"
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

	res, _, err := pn.CreateUser().Include(incl).ID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)
	assert.Equal(id, res.Data.ID)
	assert.Equal(name, res.Data.Name)
	assert.Equal(extid, res.Data.ExternalID)
	assert.Equal(purl, res.Data.ProfileURL)
	assert.Equal(email, res.Data.Email)
	assert.NotNil(res.Data.Created)
	assert.NotNil(res.Data.Updated)
	assert.NotNil(res.Data.ETag)
	assert.Equal("b", res.Data.Custom["a"])
	assert.Equal("d", res.Data.Custom["c"])

	email = "email2"

	res2, _, err2 := pn.UpdateUser().Include(incl).ID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)
	assert.Equal(id, res2.Data.ID)
	assert.Equal(name, res2.Data.Name)
	assert.Equal(extid, res2.Data.ExternalID)
	assert.Equal(purl, res2.Data.ProfileURL)
	assert.Equal(email, res2.Data.Email)
	assert.Equal(res.Data.Created, res2.Data.Created)
	assert.NotNil(res2.Data.Updated)
	assert.NotNil(res2.Data.ETag)
	assert.Equal("b", res2.Data.Custom["a"])
	assert.Equal("d", res2.Data.Custom["c"])

	res3, _, err3 := pn.GetUser().Include(incl).ID(id).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)
	assert.Equal(id, res3.Data.ID)
	assert.Equal(name, res3.Data.Name)
	assert.Equal(extid, res3.Data.ExternalID)
	assert.Equal(purl, res3.Data.ProfileURL)
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
		if res6.Data[i].ID == id {
			assert.Equal(name, res6.Data[i].Name)
			assert.Equal(extid, res6.Data[i].ExternalID)
			assert.Equal(purl, res6.Data[i].ProfileURL)
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
	res5, _, err5 := pn.DeleteUser().ID(id).Execute()
	assert.Nil(err5)
	assert.Equal(200, res5.Status)
	assert.Nil(res5.Data)

	//getuser
	res4, _, err4 := pn.GetUser().Include(incl).ID(id).Execute()
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

	res, _, err := pn.CreateSpace().Include(incl).ID(id).Name(name).Description(desc).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)
	assert.Equal(id, res.Data.ID)
	assert.Equal(name, res.Data.Name)
	assert.Equal(desc, res.Data.Description)
	assert.NotNil(res.Data.Created)
	assert.NotNil(res.Data.Updated)
	assert.NotNil(res.Data.ETag)
	assert.Equal("b", res.Data.Custom["a"])
	assert.Equal("d", res.Data.Custom["c"])

	desc = "desc2"

	res2, _, err2 := pn.UpdateSpace().Include(incl).ID(id).Name(name).Description(desc).Custom(custom).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)
	assert.Equal(id, res2.Data.ID)
	assert.Equal(name, res2.Data.Name)
	assert.Equal(desc, res2.Data.Description)
	assert.Equal(res.Data.Created, res2.Data.Created)
	assert.NotNil(res2.Data.Updated)
	assert.NotNil(res2.Data.ETag)
	assert.Equal("b", res2.Data.Custom["a"])
	assert.Equal("d", res2.Data.Custom["c"])

	res3, _, err3 := pn.GetSpace().Include(incl).ID(id).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)
	assert.Equal(id, res3.Data.ID)
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
		if res6.Data[i].ID == id {
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
	res5, _, err5 := pn.DeleteSpace().ID(id).Execute()
	assert.Nil(err5)
	assert.Equal(200, res5.Status)
	assert.Nil(res5.Data)

	//getuser
	res4, _, err4 := pn.GetSpace().Include(incl).ID(id).Execute()
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

	res, _, err := pn.CreateUser().Include(incl).ID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)
	assert.Equal(userid, res.Data.ID)
	assert.Equal(name, res.Data.Name)
	assert.Equal(extid, res.Data.ExternalID)
	assert.Equal(purl, res.Data.ProfileURL)
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

	res2, _, err2 := pn.CreateSpace().Include(incl).ID(spaceid).Name(name).Description(desc).Custom(custom2).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)
	assert.Equal(spaceid, res2.Data.ID)
	assert.Equal(name, res2.Data.Name)
	assert.Equal(desc, res2.Data.Description)
	assert.NotNil(res2.Data.Created)
	assert.NotNil(res2.Data.Updated)
	assert.NotNil(res2.Data.ETag)
	assert.Equal("b1", res2.Data.Custom["a1"])
	assert.Equal("d1", res2.Data.Custom["c1"])

	userid2 := fmt.Sprintf("testuser_%d", r.Intn(99999))

	res3, _, err3 := pn.CreateUser().Include(incl).ID(userid2).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)

	spaceid2 := fmt.Sprintf("testspace_%d", r.Intn(99999))

	res4, _, err4 := pn.CreateSpace().Include(incl).ID(spaceid2).Name(name).Description(desc).Custom(custom2).Execute()
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
		ID:     userid,
		Custom: custom3,
	}

	inArr := []pubnub.PNMembersInput{
		in,
	}

	//Add Space Memberships

	resAdd, _, errAdd := pn.ManageMembers().SpaceID(spaceid).Add(inArr).Update([]pubnub.PNMembersInput{}).Remove([]pubnub.PNMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errAdd)
	assert.Equal(200, resAdd.Status)
	assert.True(resAdd.TotalCount > 0)
	fmt.Println("resAdd-->", resAdd)
	found := false
	for i := range resAdd.Data {
		if resAdd.Data[i].ID == userid {
			found = true
			assert.Equal(custom3["a3"], resAdd.Data[i].Custom["a3"])
			assert.Equal(custom3["c3"], resAdd.Data[i].Custom["c3"])
			assert.Equal(userid, resAdd.Data[0].User.ID)
			assert.Equal(name, resAdd.Data[0].User.Name)
			assert.Equal(extid, resAdd.Data[0].User.ExternalID)
			assert.Equal(purl, resAdd.Data[0].User.ProfileURL)
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
		ID:     userid,
		Custom: custom4,
	}

	upArr := []pubnub.PNMembersInput{
		up,
	}

	resUp, _, errUp := pn.ManageMembers().SpaceID(spaceid).Add([]pubnub.PNMembersInput{}).Update(upArr).Remove([]pubnub.PNMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errUp)
	assert.Equal(200, resUp.Status)
	assert.True(resUp.TotalCount > 0)
	foundUp := false
	for i := range resUp.Data {
		if resUp.Data[i].ID == userid {
			foundUp = true
			assert.Equal("b2", resUp.Data[i].Custom["a2"])
			assert.Equal("d2", resUp.Data[i].Custom["c2"])
			assert.Equal(userid, resAdd.Data[0].User.ID)
			assert.Equal(name, resAdd.Data[0].User.Name)
			assert.Equal(extid, resAdd.Data[0].User.ExternalID)
			assert.Equal(purl, resAdd.Data[0].User.ProfileURL)
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

	resGetMem, _, errGetMem := pn.GetMemberships().UserID(userid).Include(inclMemberships).Limit(limit).Count(count).Execute()
	foundGetMem := false
	assert.Nil(errGetMem)
	for i := range resGetMem.Data {
		if resGetMem.Data[i].ID == spaceid {
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
		ID: userid,
	}

	reArr := []pubnub.PNMembersRemove{
		re,
	}
	resRem, _, errRem := pn.ManageMembers().SpaceID(spaceid).Add([]pubnub.PNMembersInput{}).Update([]pubnub.PNMembersInput{}).Remove(reArr).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errRem)
	assert.Equal(200, resRem.Status)
	foundRem := false
	for i := range resRem.Data {
		if resRem.Data[i].ID == userid {
			foundRem = true
			assert.Equal("b2", resUp.Data[i].Custom["a2"])
			assert.Equal("d2", resUp.Data[i].Custom["c2"])
			assert.Equal(userid, resUp.Data[0].User.ID)
			assert.Equal(name, resUp.Data[0].User.Name)
			assert.Equal(extid, resUp.Data[0].User.ExternalID)
			assert.Equal(purl, resUp.Data[0].User.ProfileURL)
			assert.Equal(email, resUp.Data[0].User.Email)
			assert.Equal(custom["a"], resUp.Data[0].User.Custom["a"])
			assert.Equal(custom["c"], resUp.Data[0].User.Custom["c"])

		}
	}
	assert.False(foundRem)

	inMem := pubnub.PNMembershipsInput{
		ID:     spaceid2,
		Custom: custom3,
	}

	inArrMem := []pubnub.PNMembershipsInput{
		inMem,
	}

	//Add user memberships
	resManageMemAdd, _, errManageMemAdd := pn.ManageMemberships().UserID(userid2).Add(inArrMem).Update([]pubnub.PNMembershipsInput{}).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	fmt.Println("resManageMemAdd -->", resManageMemAdd)
	assert.Nil(errManageMemAdd)
	assert.Equal(200, resManageMemAdd.Status)
	foundManageMembers := false
	for i := range resManageMemAdd.Data {
		if resManageMemAdd.Data[i].ID == spaceid2 {
			assert.Equal(spaceid2, resManageMemAdd.Data[i].Space.ID)
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
		ID:     spaceid2,
		Custom: custom5,
	}

	upArrMem := []pubnub.PNMembershipsInput{
		upMem,
	}

	resManageMemUp, _, errManageMemUp := pn.ManageMemberships().UserID(userid2).Add([]pubnub.PNMembershipsInput{}).Update(upArrMem).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	fmt.Println("resManageMemUp -->", resManageMemUp)
	assert.Nil(errManageMemUp)
	assert.Equal(200, resManageMemUp.Status)
	foundManageMembersUp := false
	for i := range resManageMemUp.Data {
		if resManageMemUp.Data[i].ID == spaceid2 {
			assert.Equal(spaceid2, resManageMemUp.Data[i].Space.ID)
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
	resGetMembers, _, errGetMembers := pn.GetMembers().SpaceID(spaceid2).Include(inclSm).Limit(limit).Count(count).Execute()
	fmt.Println("resGetMembers -->", resGetMembers)
	assert.Nil(errGetMembers)
	assert.Equal(200, resGetMembers.Status)
	foundGetMembers := false
	for i := range resGetMembers.Data {
		if resGetMembers.Data[i].ID == userid2 {
			foundGetMembers = true
			assert.Equal(name, resGetMembers.Data[i].User.Name)
			assert.Equal(extid, resGetMembers.Data[i].User.ExternalID)
			assert.Equal(purl, resGetMembers.Data[i].User.ProfileURL)
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
		ID: spaceid2,
	}

	reArrMem := []pubnub.PNMembershipsRemove{
		reMem,
	}
	resManageMemRem, _, errManageMemRem := pn.ManageMemberships().UserID(userid2).Add([]pubnub.PNMembershipsInput{}).Update([]pubnub.PNMembershipsInput{}).Remove(reArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	assert.Nil(errManageMemRem)
	assert.Equal(200, resManageMemRem.Status)

	foundManageMemRem := false
	for i := range resManageMemRem.Data {
		if resManageMemRem.Data[i].ID == spaceid2 {
			foundManageMemRem = true
		}
	}
	assert.False(foundManageMemRem)

	//delete
	res5, _, err5 := pn.DeleteUser().ID(userid).Execute()
	assert.Nil(err5)
	assert.Equal(200, res5.Status)
	assert.Nil(res5.Data)

	//delete
	res6, _, err6 := pn.DeleteSpace().ID(spaceid).Execute()
	assert.Nil(err6)
	assert.Equal(200, res6.Status)
	assert.Nil(res6.Data)

	//delete
	res52, _, err52 := pn.DeleteUser().ID(userid2).Execute()
	assert.Nil(err52)
	assert.Equal(200, res52.Status)
	assert.Nil(res52.Data)

	//delete
	res62, _, err62 := pn.DeleteSpace().ID(spaceid2).Execute()
	assert.Nil(err62)
	assert.Equal(200, res62.Status)
	assert.Nil(res62.Data)

}

func TestObjectsListeners(t *testing.T) {
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
	//pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	//pnSub.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	//Subscribe to the channel names

	listener := pubnub.NewListener()

	doneConnected := make(chan bool)
	doneUpdateUser := make(chan bool)
	doneUpdateSpace := make(chan bool)
	doneAddUserToSpace := make(chan bool)
	doneAddUserToSpace2 := make(chan bool)
	doneUpdateUserMem := make(chan bool)
	doneRemoveUserFromSpace := make(chan bool)
	doneDeleteUser := make(chan bool)
	doneDeleteSpace := make(chan bool)

	go func() {
		for {
			fmt.Println("Running =--->")
			select {

			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnected <- true
				default:
					fmt.Println(" --- status: ", status)
				}

			case userEvent := <-listener.UserEvent:

				fmt.Println(" --- UserEvent: ")
				fmt.Println(fmt.Sprintf("%s", userEvent))
				fmt.Println(fmt.Sprintf("userEvent.Channel: %s", userEvent.Channel))
				fmt.Println(fmt.Sprintf("userEvent.SubscribedChannel: %s", userEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("userEvent.Event: %s", userEvent.Event))
				fmt.Println(fmt.Sprintf("userEvent.UserID: %s", userEvent.UserID))
				fmt.Println(fmt.Sprintf("userEvent.Description: %s", userEvent.Description))
				fmt.Println(fmt.Sprintf("userEvent.Timestamp: %s", userEvent.Timestamp))
				fmt.Println(fmt.Sprintf("userEvent.Name: %s", userEvent.Name))
				fmt.Println(fmt.Sprintf("userEvent.ExternalID: %s", userEvent.ExternalID))
				fmt.Println(fmt.Sprintf("userEvent.ProfileURL: %s", userEvent.ProfileURL))
				fmt.Println(fmt.Sprintf("userEvent.Email: %s", userEvent.Email))
				fmt.Println(fmt.Sprintf("userEvent.Created: %s", userEvent.Created))
				fmt.Println(fmt.Sprintf("userEvent.Updated: %s", userEvent.Updated))
				fmt.Println(fmt.Sprintf("userEvent.ETag: %s", userEvent.ETag))
				fmt.Println(fmt.Sprintf("userEvent.Custom: %v", userEvent.Custom))

				if (userEvent.Event == pubnub.PNObjectsEventDelete) && (userEvent.UserID == userid) {
					doneDeleteUser <- true
				}
				if (userEvent.Event == pubnub.PNObjectsEventUpdate) && (userEvent.UserID == userid) {
					doneUpdateUser <- true
				}
			case spaceEvent := <-listener.SpaceEvent:

				fmt.Println(" --- SpaceEvent: ")
				fmt.Println(fmt.Sprintf("%s", spaceEvent))
				fmt.Println(fmt.Sprintf("spaceEvent.Channel: %s", spaceEvent.Channel))
				fmt.Println(fmt.Sprintf("spaceEvent.SubscribedChannel: %s", spaceEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("spaceEvent.Event: %s", spaceEvent.Event))
				fmt.Println(fmt.Sprintf("spaceEvent.SpaceID: %s", spaceEvent.SpaceID))
				fmt.Println(fmt.Sprintf("spaceEvent.Description: %s", spaceEvent.Description))
				fmt.Println(fmt.Sprintf("spaceEvent.Timestamp: %s", spaceEvent.Timestamp))
				fmt.Println(fmt.Sprintf("spaceEvent.Created: %s", spaceEvent.Created))
				fmt.Println(fmt.Sprintf("spaceEvent.Updated: %s", spaceEvent.Updated))
				fmt.Println(fmt.Sprintf("spaceEvent.ETag: %s", spaceEvent.ETag))
				fmt.Println(fmt.Sprintf("spaceEvent.Custom: %v", spaceEvent.Custom))
				if (spaceEvent.Event == pubnub.PNObjectsEventDelete) && (spaceEvent.SpaceID == spaceid) {
					doneDeleteSpace <- true
				}
				if (spaceEvent.Event == pubnub.PNObjectsEventUpdate) && (spaceEvent.SpaceID == spaceid) {
					doneUpdateSpace <- true
				}

			case membershipEvent := <-listener.MembershipEvent:

				fmt.Println(" --- MembershipEvent: ")
				fmt.Println(fmt.Sprintf("%s", membershipEvent))
				fmt.Println(fmt.Sprintf("membershipEvent.Channel: %s", membershipEvent.Channel))
				fmt.Println(fmt.Sprintf("membershipEvent.SubscribedChannel: %s", membershipEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("membershipEvent.Event: %s", membershipEvent.Event))
				fmt.Println(fmt.Sprintf("membershipEvent.SpaceID: %s", membershipEvent.SpaceID))
				fmt.Println(fmt.Sprintf("membershipEvent.UserID: %s", membershipEvent.UserID))
				fmt.Println(fmt.Sprintf("membershipEvent.Description: %s", membershipEvent.Description))
				fmt.Println(fmt.Sprintf("membershipEvent.Timestamp: %s", membershipEvent.Timestamp))
				fmt.Println(fmt.Sprintf("membershipEvent.Custom: %v", membershipEvent.Custom))
				if (membershipEvent.Event == pubnub.PNObjectsEventCreate) && (membershipEvent.SpaceID == spaceid) && (membershipEvent.UserID == userid) && (membershipEvent.Channel == spaceid) {
					doneAddUserToSpace <- true
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventCreate) && (membershipEvent.SpaceID == spaceid) && (membershipEvent.UserID == userid) && ((membershipEvent.Channel == userid) || (membershipEvent.Channel == fmt.Sprintf("pnuser-%s", userid))) {
					doneAddUserToSpace2 <- true
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventUpdate) && (membershipEvent.SpaceID == spaceid) && (membershipEvent.UserID == userid) && (membershipEvent.Channel == spaceid) {
					doneUpdateUserMem <- true
				}
				if (membershipEvent.Event == pubnub.PNObjectsEventDelete) && (membershipEvent.SpaceID == spaceid) && (membershipEvent.UserID == userid) && (membershipEvent.Channel == spaceid) {
					doneRemoveUserFromSpace <- true
				}

			}
		}

	}()

	pnSub.AddListener(listener)

	pnSub.Subscribe().Channels([]string{fmt.Sprintf("pnuser-%s", userid), userid, spaceid}).Execute()
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

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	//Create User
	res, _, err := pn.CreateUser().Include(incl).ID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(customUser).Execute()
	assert.Nil(err)
	assert.Equal(200, res.Status)

	//Create Space
	customSpace := make(map[string]interface{})
	customSpace["as"] = "bs"
	customSpace["cs"] = "ds"

	res4, _, err4 := pn.CreateSpace().Include(incl).ID(spaceid).Name(name).Description(desc).Custom(customSpace).Execute()
	assert.Nil(err4)
	assert.Equal(200, res4.Status)

	time.Sleep(1 * time.Second)

	//Update User
	email = "email2"

	res2, _, err2 := pn.UpdateUser().Include(incl).ID(userid).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(customUser).Execute()
	assert.Nil(err2)
	assert.Equal(200, res2.Status)

	//Read event
	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneUpdateUser:
		assert.True(true)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	time.Sleep(1 * time.Second)

	desc = "desc2"

	//Update Space
	res3, _, err3 := pn.UpdateSpace().Include(incl).ID(spaceid).Name(name).Description(desc).Custom(customSpace).Execute()
	assert.Nil(err3)
	assert.Equal(200, res3.Status)

	//Read event
	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneUpdateSpace:
		assert.True(true)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	//Add user to space
	inclSm := []pubnub.PNMembersInclude{
		pubnub.PNMembersCustom,
		pubnub.PNMembersUser,
		pubnub.PNMembersUserCustom,
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	in := pubnub.PNMembersInput{
		ID:     userid,
		Custom: custom3,
	}

	inArr := []pubnub.PNMembersInput{
		in,
	}
	time.Sleep(1 * time.Second)

	resAdd, _, errAdd := pn.ManageMembers().SpaceID(spaceid).Add(inArr).Update([]pubnub.PNMembersInput{}).Remove([]pubnub.PNMembersRemove{}).Include(inclSm).Limit(limit).Count(count).Execute()
	assert.Nil(errAdd)
	assert.Equal(200, resAdd.Status)

	//Read event
	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	addUserToSpace := false
	addUserToSpace2 := false

	runfor := true
	waitforfunc := make(chan bool)

	go func() {
	LabelBreak:
		for runfor {

			select {
			case <-doneAddUserToSpace:
				addUserToSpace = true
				if addUserToSpace2 {
					runfor = false
					fmt.Println("break 1")
					waitforfunc <- true
					break LabelBreak
				}
			case <-doneAddUserToSpace2:
				addUserToSpace2 = true
				if addUserToSpace {
					runfor = false
					fmt.Println("break 2")
					waitforfunc <- true
					break LabelBreak
				}
			case <-tic.C:
				tic.Stop()
				assert.Fail("timeout")
				waitforfunc <- true
				break LabelBreak
			}

		}

	}()

	<-waitforfunc

	assert.True(addUserToSpace && addUserToSpace2)

	//Update user membership

	custom5 := make(map[string]interface{})
	custom5["a5"] = "b5"
	custom5["c5"] = "d5"

	upMem := pubnub.PNMembershipsInput{
		ID:     spaceid,
		Custom: custom5,
	}

	upArrMem := []pubnub.PNMembershipsInput{
		upMem,
	}

	inclMemberships := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsCustom,
		pubnub.PNMembershipsSpace,
		pubnub.PNMembershipsSpaceCustom,
	}

	resManageMemUp, _, errManageMemUp := pn.ManageMemberships().UserID(userid).Add([]pubnub.PNMembershipsInput{}).Update(upArrMem).Remove([]pubnub.PNMembershipsRemove{}).Include(inclMemberships).Limit(limit).Count(count).Execute()
	fmt.Println("resManageMemUp -->", resManageMemUp)
	assert.Nil(errManageMemUp)
	assert.Equal(200, resManageMemUp.Status)

	//Read event
	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneUpdateUserMem:
		assert.True(true)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	//Remove user from space
	reMem := pubnub.PNMembershipsRemove{
		ID: spaceid,
	}

	reArrMem := []pubnub.PNMembershipsRemove{
		reMem,
	}
	resManageMemRem, _, errManageMemRem := pn.ManageMemberships().UserID(userid).Add([]pubnub.PNMembershipsInput{}).Update([]pubnub.PNMembershipsInput{}).Remove(reArrMem).Include(inclMemberships).Limit(limit).Count(count).Execute()
	assert.Nil(errManageMemRem)
	assert.Equal(200, resManageMemRem.Status)

	//Read event
	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneRemoveUserFromSpace:
		assert.True(true)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	//Delete user
	res52, _, err52 := pn.DeleteUser().ID(userid).Execute()
	assert.Nil(err52)
	assert.Equal(200, res52.Status)
	assert.Nil(res52.Data)

	//Read event

	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneDeleteUser:
		assert.True(true)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}

	//Delete Space
	res62, _, err62 := pn.DeleteSpace().ID(spaceid).Execute()
	assert.Nil(err62)
	assert.Equal(200, res62.Status)
	assert.Nil(res62.Data)

	//Read event

	tic = time.NewTicker(time.Duration(eventWaitTime) * time.Second)
	select {
	case <-doneDeleteSpace:
		assert.True(true)
	case <-tic.C:
		tic.Stop()
		assert.Fail("timeout")
	}
}
