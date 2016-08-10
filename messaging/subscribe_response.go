package messaging

/*import (
	"fmt"
)*/

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
