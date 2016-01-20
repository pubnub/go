package messaging

import (
	"fmt"
	"strings"
	"time"
)

// Timeout channel for non-subscribe requests
func Timeout() <-chan time.Time {
	return Timeouts(GetNonSubscribeTimeout())
}

// Timeout channel for subscribe requests
func SubscribeTimeout() <-chan time.Time {
	return Timeouts(GetSubscribeTimeout())
}

// Timeout channel with custon timeout value
func Timeouts(seconds uint16) <-chan time.Time {
	return time.After(time.Second * time.Duration(seconds))
}

// Type for presence events
type PresenceEvent struct {
	Action    string  `json:"action"`
	Uuid      string  `json:"uuid"`
	Timestamp float64 `json:"timestamp"`
	Occupancy int     `json:"occupancy"`
}

func stringPresenceOrSubscribe(channel string) string {
	const (
		subscribeMessage string = "Subscription to"
		presenceMessage  string = "Presence notifications for"
	)

	if strings.HasSuffix(channel, presenceSuffix) {
		return presenceMessage
	} else {
		return subscribeMessage
	}
}

func splitItems(items string) []string {
	if items == "" {
		return []string{}
	} else {
		return strings.Split(items, ",")
	}
}

func addPnpresToString(items string) string {

	var presenceItems []string

	itemsSlice := splitItems(items)

	for _, v := range itemsSlice {
		presenceItems = append(presenceItems,
			fmt.Sprintf("%s%s", v, presenceSuffix))
	}

	return strings.Join(presenceItems, ",")
}

func removePnpres(initial string) string {
	return strings.TrimSuffix(initial, presenceSuffix)
}
