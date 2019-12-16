package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushPayload(t *testing.T) {
	PushPayloadCommon(t, false, true, true, true, true, true, false)
}

func TestPushPayloadWithCtx(t *testing.T) {
	PushPayloadCommon(t, true, true, true, true, true, true, false)
}

func PushPayloadCommon(t *testing.T, withContext, withAPNS, withAPNS2, withFCM, withMPNS, withCommonPayload, setAPNSAlert bool) {
	assert := assert.New(t)

	o := newPublishPushHelperBuilder(pubnub)
	if withContext {
		o = newPublishPushHelperBuilderWithContext(pubnub, backgroundContext)
	}

	aps := PNAPSData{
		Alert: "apns alert",
		Badge: 1,
		Sound: "ding",
		Custom: map[string]interface{}{
			"aps_key1": "aps_value1",
			"aps_key2": "aps_value2",
		},
	}
	if !setAPNSAlert {
		aps.Alert = nil
		aps.Title = "title"
		aps.Subtitle = "subtitle"
		aps.Body = "body"
	}

	apns := PNAPNSData{
		APS: aps,
		Custom: map[string]interface{}{
			"apns_key1": "apns_value1",
			"apns_key2": "apns_value2",
		},
	}

	apns2One := PNAPNS2Data{
		CollapseID: "invitations",
		Expiration: "2019-12-13T22:06:09Z",
		Version:    "v1",
		Targets: []PNPushTarget{
			PNPushTarget{
				Environment: PNPushEnvironmentDevelopment,
				Topic:       "com.meetings.chat.app",
				ExcludeDevices: []string{
					"device1",
					"device2",
				},
			},
		},
	}

	apns2Two := PNAPNS2Data{
		CollapseID: "invitations",
		Expiration: "2019-12-15T22:06:09Z",
		Version:    "v2",
		Targets: []PNPushTarget{
			PNPushTarget{
				Environment: PNPushEnvironmentProduction,
				Topic:       "com.meetings.chat.app",
				ExcludeDevices: []string{
					"device3",
					"device4",
				},
			},
		},
	}

	apns2 := []PNAPNS2Data{apns2One, apns2Two}

	if withAPNS2 || withAPNS {
		o.SetAPNSPayload(apns, nil)
		if withAPNS2 {
			o.SetAPNSPayload(apns, apns2)
		}
	}

	mpns := PNMPNSData{
		Title:       "title",
		Type:        "type",
		Count:       1,
		BackTitle:   "BackTitle",
		BackContent: "BackContent",
		Custom: map[string]interface{}{
			"mpns_key1": "mpns_value1",
			"mpns_key2": "mpns_value2",
		},
	}

	if withMPNS {
		o.SetMPNSPayload(mpns)
	}

	fcm := PNFCMData{
		Data: PNFCMDataFields{
			Summary: "summary",
			Custom: map[string]interface{}{
				"fcm_data_key1": "fcm_data_value1",
				"fcm_data_key2": "fcm_data_value2",
			},
		},
		Custom: map[string]interface{}{
			"fcm_key1": "fcm_value1",
			"fcm_key2": "fcm_value2",
		},
	}

	if withFCM {
		o.SetFCMPayload(fcm)
	}

	CommonPayload := map[string]interface{}{
		"a": map[string]interface{}{
			"common_key1": "common_value1",
			"common_key2": "common_value2",
		},
		"b": "val",
	}

	if withCommonPayload {
		o.SetCommonPayload(CommonPayload)
	}

	result := o.BuildPayload()
	assert.NotNil(result)
	if result != nil {
		if withAPNS2 || withAPNS {
			resAPNS := result["pn_apns"].(map[string]interface{})
			resAPS := resAPNS["aps"].(map[string]interface{})
			assert.Equal(apns.APS.Badge, resAPS["badge"])
			assert.Equal(apns.APS.Sound, resAPS["sound"])
			if setAPNSAlert {
				assert.Equal(apns.APS.Alert, resAPS["alert"])
			} else {
				resAlert := resAPS["alert"].(map[string]interface{})
				assert.Equal(apns.APS.Title, resAlert["title"])
				assert.Equal(apns.APS.Subtitle, resAlert["subtitle"])
				assert.Equal(apns.APS.Body, resAlert["body"])
			}
			assert.Equal(apns.APS.Custom["aps_key1"], resAPS["aps_key1"])
			assert.Equal(apns.APS.Custom["aps_key2"], resAPS["aps_key2"])
			assert.Equal(apns.Custom["apns_key1"], resAPNS["apns_key1"])
			assert.Equal(apns.Custom["apns_key2"], resAPNS["apns_key2"])

			if withAPNS2 {
				resAPNS2 := result["pn_push"].([]PNAPNS2Data)
				assert.Equal(apns2[0].CollapseID, resAPNS2[0].CollapseID)
				assert.Equal(apns2[0].Expiration, resAPNS2[0].Expiration)
				assert.Equal(apns2[0].Version, resAPNS2[0].Version)
				assert.True(apns2[0].Targets[0].Environment == resAPNS2[0].Targets[0].Environment)
				assert.Equal(apns2[0].Targets[0].Topic, resAPNS2[0].Targets[0].Topic)
				assert.Equal(apns2[0].Targets[0].ExcludeDevices[0], resAPNS2[0].Targets[0].ExcludeDevices[0])
				assert.Equal(apns2[0].Targets[0].ExcludeDevices[1], resAPNS2[0].Targets[0].ExcludeDevices[1])

				assert.Equal(apns2[1].CollapseID, resAPNS2[1].CollapseID)
				assert.Equal(apns2[1].Expiration, resAPNS2[1].Expiration)
				assert.Equal(apns2[1].Version, resAPNS2[1].Version)
				assert.True(apns2[1].Targets[0].Environment == resAPNS2[1].Targets[0].Environment)
				assert.Equal(apns2[1].Targets[0].Topic, resAPNS2[1].Targets[0].Topic)
				assert.Equal(apns2[1].Targets[0].ExcludeDevices[0], resAPNS2[1].Targets[0].ExcludeDevices[0])
				assert.Equal(apns2[1].Targets[0].ExcludeDevices[1], resAPNS2[1].Targets[0].ExcludeDevices[1])

			}
		}

		if withMPNS {
			resMPNS := result["pn_mpns"].(map[string]interface{})
			assert.Equal(mpns.Title, resMPNS["title"])
			assert.Equal(mpns.Type, resMPNS["type"])
			assert.Equal(mpns.Count, resMPNS["count"])
			assert.Equal(mpns.BackTitle, resMPNS["back_title"])
			assert.Equal(mpns.BackContent, resMPNS["back_content"])
			assert.Equal(mpns.Custom["mpns_key1"], resMPNS["mpns_key1"])
			assert.Equal(mpns.Custom["mpns_key2"], resMPNS["mpns_key2"])
		}

		if withFCM {
			resFCM := result["pn_gcm"].(map[string]interface{})
			resFCMData := resFCM["data"].(map[string]interface{})

			assert.Equal(resFCMData["summary"], resFCMData["summary"])

			assert.Equal(fcm.Data.Custom["fcm_data_key1"], resFCMData["fcm_data_key1"])
			assert.Equal(fcm.Data.Custom["fcm_data_key2"], resFCMData["fcm_data_key2"])

			assert.Equal(fcm.Custom["fcm_key1"], resFCM["fcm_key1"])
			assert.Equal(fcm.Custom["fcm_key2"], resFCM["fcm_key2"])
		}

		if withCommonPayload {
			resCommonPayloadA := result["a"].(map[string]interface{})
			CommonPayloadA := CommonPayload["a"].(map[string]interface{})
			assert.Equal(CommonPayloadA["common_key1"], resCommonPayloadA["common_key1"])
			assert.Equal(CommonPayloadA["common_key2"], resCommonPayloadA["common_key2"])
		}
	}
}

func TestPushAPNSPayload(t *testing.T) {
	PushPayloadCommon(t, false, true, false, false, false, false, false)
}

func TestPushAPNSPayloadWithSetAlert(t *testing.T) {
	PushPayloadCommon(t, false, true, false, false, false, false, true)
}

func TestPushAPNS2Payload(t *testing.T) {
	PushPayloadCommon(t, false, false, true, false, false, false, false)
}

func TestPushMPNSPayload(t *testing.T) {
	PushPayloadCommon(t, false, false, false, false, true, false, false)
}

func TestPushFCMPayload(t *testing.T) {
	PushPayloadCommon(t, false, false, false, true, false, false, false)
}

func TestPushCommonPayload(t *testing.T) {
	PushPayloadCommon(t, false, false, false, false, false, true, false)
}

func TestPushAPNSPayloadWithCtx(t *testing.T) {
	PushPayloadCommon(t, false, true, false, false, false, false, false)
}

func TestPushAPNS2PayloadWithCtx(t *testing.T) {
	PushPayloadCommon(t, false, false, true, false, false, false, false)
}

func TestPushMPNSPayloadWithCtx(t *testing.T) {
	PushPayloadCommon(t, false, false, false, false, true, false, false)
}

func TestPushFCMPayloadWithCtx(t *testing.T) {
	PushPayloadCommon(t, false, false, false, true, false, false, false)
}

func TestPushCommonPayloadWithCtx(t *testing.T) {
	PushPayloadCommon(t, false, false, false, false, false, true, false)
}
