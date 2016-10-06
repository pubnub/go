package messaging

import (
	"fmt"
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

func (env *subscribeEnvelope) getChannelsAndGroups() (channels, channelGroups []string) {
	if env.Messages != nil {
		count := 0
		for _, msg := range env.Messages {
			count++
			msg.writeMessageLog(count)
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

func (msg *subscribeMessage) writeMessageLog(count int) {
	// start logging
	infoLogger.Println(fmt.Sprintf("-----Message %d-----", count))
	infoLogger.Println(fmt.Sprintf("Channel, %s", msg.Channel))
	infoLogger.Println(fmt.Sprintf("Flags, %d", msg.Flags))
	infoLogger.Println(fmt.Sprintf("IssuingClientId, %s", msg.IssuingClientId))
	infoLogger.Println(fmt.Sprintf("OriginatingTimetoken Region, %d", msg.OriginatingTimetoken.Region))
	infoLogger.Println(fmt.Sprintf("OriginatingTimetoken Timetoken, %s", msg.OriginatingTimetoken.Timetoken))
	infoLogger.Println(fmt.Sprintf("PublishTimetokenMetadata Region, %d", msg.PublishTimetokenMetadata.Region))
	infoLogger.Println(fmt.Sprintf("PublishTimetokenMetadata Timetoken, %s", msg.PublishTimetokenMetadata.Timetoken))

	strPayload, ok := msg.Payload.(string)
	if ok {
		infoLogger.Println(fmt.Sprintf("Payload, %s", strPayload))
	} else {
		infoLogger.Println(fmt.Sprintf("Payload, not converted to string %s", msg.Payload))
	}
	infoLogger.Println(fmt.Sprintf("SequenceNumber, %d", msg.SequenceNumber))
	infoLogger.Println(fmt.Sprintf("Shard, %s", msg.Shard))
	infoLogger.Println(fmt.Sprintf("SubscribeKey, %s", msg.SubscribeKey))
	infoLogger.Println(fmt.Sprintf("SubscriptionMatch, %s", msg.SubscriptionMatch))
	strUserMetadata, ok := msg.UserMetadata.(string)
	if ok {
		infoLogger.Println(fmt.Sprintf("UserMetadata, %s", strUserMetadata))
	} else {
		infoLogger.Println(fmt.Sprintf("UserMetadata, not converted to string"))
	}
	// end logging
}
