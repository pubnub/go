package pubnub

import (
//"fmt"
//"strconv"

)

type PNMPNSData struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Count       int    `json:"count"`
	BackTitle   string `json:"back_title"`
	BackContent string `json:"back_content"`
	Custom      map[string]interface{}
}

type PNFCMData struct {
	Data   PNFCMDataFields `json:"data"`
	Custom map[string]interface{}
}

type PNFCMDataFields struct {
	Summary interface{} `json:"summary"`
	Custom  map[string]interface{}
}

type PNAPSData struct {
	Alert    interface{} `json:"alert"`
	Badge    int         `json:"badge"`
	Sound    string      `json:"sound"`
	Title    string      `json:"title"`
	Subtitle string      `json:"subtitle"`
	Body     string      `json:"body"`
	Custom   map[string]interface{}
}

type PNAPNSData struct {
	APS    PNAPSData `json:"aps"`
	Custom map[string]interface{}
}

type PNAPNS2Data struct {
	CollapseId string         `json:"collapseId"`
	Expiration string         `json:"expiration"`
	Targets    []PNPushTarget `json:"targets"`
	Version    string         `json:"version"`
}

type PNPushTarget struct {
	Topic          string            `json:"topic"`
	ExcludeDevices []string          `json:"exclude_devices"`
	Environment    PNPushEnvironment `json:"environment"`
}

type publishPushHelperBuilder struct {
	opts *publishPushHelperOpts
}

func newPublishPushHelperBuilder(pubnub *PubNub) *publishPushHelperBuilder {
	builder := publishPushHelperBuilder{
		opts: &publishPushHelperOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newPublishPushHelperBuilderWithContext(pubnub *PubNub,
	context Context) *publishPushHelperBuilder {
	builder := publishPushHelperBuilder{
		opts: &publishPushHelperOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// func (b *publishPushHelperBuilder) APNS2SupportData(pushAPNS2SupportData []PushAPNS2Data) *publishPushHelperBuilder {
// 	b.opts.PushAPNS2Data = pushAPNS2SupportData

// 	return b
// }

// // APNS2Data APNS2 Support Data
// func (b *publishPushHelperBuilder) PushTypes(pushTypes []PNPushType) *publishPushHelperBuilder {
// 	b.opts.PushTypes = pushTypes

// 	return b
// }

// func (b *publishPushHelperBuilder) Title(title string) *publishPushHelperBuilder {
// 	b.opts.PushTitle = title

// 	return b
// }

// func (b *publishPushHelperBuilder) Body(body string) *publishPushHelperBuilder {
// 	b.opts.PushBody = body

// 	return b
// }

// func (b *publishPushHelperBuilder) Badge(badge int) *publishPushHelperBuilder {
// 	b.opts.PushBadge = badge

// 	return b
// }

// func (b *publishPushHelperBuilder) Sound(sound string) *publishPushHelperBuilder {
// 	b.opts.PushSound = sound

// 	return b
// }

// func (b *publishPushHelperBuilder) Email(email string) *publishPushHelperBuilder {
// 	b.opts.Email = email

// 	return b
// }

// func (b *publishPushHelperBuilder) Custom(custom map[string]interface{}) *publishPushHelperBuilder {
// 	b.opts.PushCustomData = custom

// 	return b
// }

// GetPayload runs the publishPushHelper request.
func (b *publishPushHelperBuilder) BuildPayload() (map[string]interface{}, error) {

	// err := b.opts.validate()
	// if err != nil {
	// 	return emptyPNPublishPushHelperResponse, err
	// }

	response := make(map[string]interface{})
	apns := b.opts.buildAPNSPayload()
	if apns != nil {
		response["pn_apns"] = apns
		apns2 := b.opts.PushAPNS2Data
		if apns2 != nil {
			response["pn_push"] = apns2
		}
	}

	mpns := b.opts.buildMPNSPayload()
	if mpns != nil {
		response["pn_mpns"] = mpns
	}

	fcm := b.opts.buildFCMPayload()
	if fcm != nil {
		response["pn_gcm"] = fcm
	}

	if b.opts.CommonPayload != nil {
		for key, value := range b.opts.CommonPayload {
			response[key] = value
		}
	}

	// for value := range b.opts.PushTypes {
	// 	switch value {
	// 	case PNPushTypeAPNS:
	// 		response.PushPayload["pn_apns"] = b.buildAPNSPayload()
	// 	case PNPushTypeAPNS2:
	// 		response.PushPayload["pn_push"] = b.opts.PushAPNS2Data
	// 		response.PushPayload["pn_apns"] = b.buildAPNSPayload()
	// 	case PNPushTypeMPNS:
	// 		response.PushPayload["pn_mpns"] = b.buildMPNSPayload()
	// 	case PNPushTypeGCM:
	// 		response.PushPayload["pn_gcm"] = b.buildFCMPayload()
	// 	default:
	// 	}
	// }

	return response, nil
}

func (o *publishPushHelperOpts) buildAPNSPayload() map[string]interface{} {
	apns := make(map[string]interface{})
	if o.PushAPNSData != nil {
		aps := &o.PushAPNSData.APS
		if aps != nil {
			apsData := make(map[string]interface{})
			if aps.Alert != nil {
				apsData["alert"] = aps.Alert
			} else if aps.Subtitle != "" || aps.Body != "" || aps.Title != "" {
				alert := make(map[string]interface{})
				if aps.Subtitle != "" {
					alert["subtitle"] = aps.Subtitle
				}
				if aps.Title != "" {
					alert["title"] = aps.Title
				}
				if aps.Body != "" {
					alert["body"] = aps.Body
				}
				apsData["alert"] = alert
			}

			apsData["badge"] = aps.Badge

			if aps.Sound != "" {
				apsData["sound"] = aps.Sound
			}

			if aps.Custom != nil {
				for key, value := range aps.Custom {
					apsData[key] = value
				}
			}
			apns["aps"] = apsData
		}
		custom := o.PushAPNSData.Custom
		if custom != nil {
			for key, value := range custom {
				apns[key] = value
			}
		}
	}

	return apns
}

func (o *publishPushHelperOpts) buildMPNSPayload() map[string]interface{} {
	mpns := make(map[string]interface{})
	if o.PushMPNSData != nil {
		if o.PushMPNSData.Title != "" {
			mpns["title"] = o.PushMPNSData.Title
		}
		if o.PushMPNSData.Type != "" {
			mpns["type"] = o.PushMPNSData.Type
		}
		if o.PushMPNSData.BackTitle != "" {
			mpns["back_title"] = o.PushMPNSData.BackTitle
		}
		if o.PushMPNSData.BackContent != "" {
			mpns["back_content"] = o.PushMPNSData.BackContent
		}
		mpns["count"] = o.PushMPNSData.Count

		custom := o.PushMPNSData.Custom
		if custom != nil {
			for key, value := range custom {
				mpns[key] = value
			}
		}
	}

	return mpns
}

func (o *publishPushHelperOpts) buildFCMPayload() map[string]interface{} {
	fcm := make(map[string]interface{})
	if o.PushFCMData != nil {
		data := &o.PushFCMData.Data
		if data != nil {
			fcmData := make(map[string]interface{})
			if data.Summary != nil {
				fcmData["summary"] = data.Summary
			}

			custom := data.Custom
			if custom != nil {
				for key, value := range custom {
					fcmData[key] = value
				}
			}
			fcm["data"] = fcmData
		}

		custom := o.PushFCMData.Custom
		if custom != nil {
			for key, value := range custom {
				fcm[key] = value
			}
		}
	}

	return fcm
}

func (b *publishPushHelperBuilder) SetAPNSPayload(pnAPNSData PNAPNSData, pnAPNS2Data []PNAPNS2Data) *publishPushHelperBuilder {
	b.opts.PushAPNSData = &pnAPNSData
	b.opts.PushAPNS2Data = pnAPNS2Data

	return b
}

func (b *publishPushHelperBuilder) SetMPNSPayload(pnMPNSData PNMPNSData) *publishPushHelperBuilder {
	b.opts.PushMPNSData = &pnMPNSData

	return b
}

func (b *publishPushHelperBuilder) SetCommonPayload(commonPayload map[string]interface{}) *publishPushHelperBuilder {
	b.opts.CommonPayload = commonPayload

	return b
}

func (b *publishPushHelperBuilder) SetFCMPayload(pnFCMData PNFCMData) *publishPushHelperBuilder {
	b.opts.PushFCMData = &pnFCMData

	return b
}

type publishPushHelperOpts struct {
	pubnub *PubNub

	//PushTypes      []PNPushType
	// PushTitle      string
	// PushBody       string
	// PushBadge      int
	// PushSound      string
	PushAPNS2Data  []PNAPNS2Data
	PushAPNSData   *PNAPNSData
	PushMPNSData   *PNMPNSData
	PushFCMData    *PNFCMData
	CommonPayload  map[string]interface{}
	PushCustomData map[string]interface{}

	ctx Context
}

func (o *publishPushHelperOpts) context() Context {
	return o.ctx
}

func (o *publishPushHelperOpts) validate() error {
	// if len(o.PushTypes) <= 0 {
	// 	return newValidationError(o, StrMissingPushType)
	// }

	// if len(o.PushTitle) <= 0 {
	// 	return newValidationError(o, StrMissingPushTitle)
	// }

	return nil
}
