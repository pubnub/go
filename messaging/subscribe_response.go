package messaging

import (
	"strings"
)

type subscribeEnvelope struct {
	Messages      []subscribeMessage `json:"m"`
	TimetokenMeta timetokenMetadata  `json:"t"`
}

type timetokenMetadata struct {
	Timetoken string `json:"t"`
	Region    int    `json:"r"`
}

type subscribeMessage struct {
	Shard                    string            `json:"a"`
	SubscriptionMatch        string            `json:"b"`
	Channel                  string            `json:"c"`
	Payload                  interface{}       `json:"d"`
	Flags                    int               `json:"f"`
	IssuingClientId          string            `json:"i"`
	SubscribeKey             string            `json:"k"`
	SequenceNumber           uint64            `json:"s"`
	OriginatingTimetoken     timetokenMetadata `json:"o"`
	PublishTimetokenMetadata timetokenMetadata `json:"p"`
	UserMetadata             interface{}       `json:"u"`
	//WaypointList string `json:"w"`
	//EatAfterReading bool `json:"ear"`
	//ReplicationMap interface{} `json:"r"`
}

type subscribeMessageResponseV2 struct {
	ChannelGroup             string            `json:"ChannelGroup"`
	Channel                  string            `json:"Channel"`
	Payload                  interface{}       `json:"Payload"`
	IssuingClientId          string            `json:"IssuingClientId"`
	OriginatingTimetoken     timetokenMetadata `json:"OriginatingTimetoken"`
	PublishTimetokenMetadata timetokenMetadata `json:"PublishTimetokenMetadata"`
	UserMetadata             interface{}       `json:"UserMetadata"`
	//SubscribeKey             string            `json:"SubscribeKey"`
	//SequenceNumber           uint64            `json:"SequenceNumber"`
}

type presenceMessageResponseV2 struct {
	ChannelGroup         string            `json:"ChannelGroup"`
	Channel              string            `json:"Channel"`
	IssuingClientId      string            `json:"IssuingClientId"`
	OriginatingTimetoken timetokenMetadata `json:"OriginatingTimetoken"`
	UserMetadata         interface{}       `json:"UserMetadata"`
	State                interface{}       `json:"State"`
	Event                string            `json:"Event"`
	UUID                 string            `json:"UUID"`
	Timestamp            string            `json:"Timestamp"`
	Occupancy            int               `json:"Occupancy"`
}

type statusResponse struct {
}

func (msg *subscribeMessage) getMessageResponse() *subscribeMessageResponse {
	res := &subscribeMessageResponse{}
	res.Channel = msg.Channel
	res.IssuingClientId = msg.IssuingClientId
	res.OriginatingTimetoken = msg.OriginatingTimetoken
	res.Payload = msg.Payload
	res.PublishTimetokenMetadata = msg.PublishTimetokenMetadata
	res.SequenceNumber = msg.SequenceNumber
	res.SubscribeKey = msg.SubscribeKey
	res.ChannelGroup = msg.SubscriptionMatch
	res.UserMetadata = msg.UserMetadata
	return res
}

func (env *subscribeEnvelope) getChannelsAndGroups(pub *Pubnub) (channels, channelGroups []string) {
	if env.Messages != nil {
		count := 0
		for _, msg := range env.Messages {
			count++
			msg.writeMessageLog(count, pub)
			channels = append(channels, msg.Channel)
			if (msg.Channel != msg.SubscriptionMatch) &&
				(!strings.Contains(msg.SubscriptionMatch, ".*")) &&
				(msg.SubscriptionMatch != "") {
				channelGroups = append(channelGroups, msg.SubscriptionMatch)
			}
		}
	}
	return channels, channelGroups
}

func (msg *subscribeMessage) writeMessageLog(count int, pub *Pubnub) {
	// start logging
	infoLogger.Printf("INFO: -----Message %d-----", count)
	infoLogger.Printf("INFO: Channel, %s", msg.Channel)
	infoLogger.Printf("INFO: Flags, %d", msg.Flags)
	infoLogger.Printf("INFO: IssuingClientId, %s", msg.IssuingClientId)
	infoLogger.Printf("INFO: OriginatingTimetoken Region, %d", msg.OriginatingTimetoken.Region)
	infoLogger.Printf("INFO: OriginatingTimetoken Timetoken, %s", msg.OriginatingTimetoken.Timetoken)
	infoLogger.Printf("INFO: PublishTimetokenMetadata Region, %d", msg.PublishTimetokenMetadata.Region)
	infoLogger.Printf("INFO: PublishTimetokenMetadata Timetoken, %s", msg.PublishTimetokenMetadata.Timetoken)

	strPayload, ok := msg.Payload.(string)
	if ok {
		infoLogger.Printf("INFO: Payload, %s", strPayload)
	} else {
		infoLogger.Printf("INFO: Payload, not converted to string %s", msg.Payload)
	}
	infoLogger.Printf("INFO: SequenceNumber, %d", msg.SequenceNumber)
	infoLogger.Printf("INFO: Shard, %s", msg.Shard)
	infoLogger.Printf("INFO: SubscribeKey, %s", msg.SubscribeKey)
	infoLogger.Printf("INFO: SubscriptionMatch, %s", msg.SubscriptionMatch)
	strUserMetadata, ok := msg.UserMetadata.(string)
	if ok {
		infoLogger.Printf("INFO: UserMetadata, %s", strUserMetadata)
	} else {
		infoLogger.Printf("INFO: UserMetadata, not converted to string")
	}
	// end logging
}
