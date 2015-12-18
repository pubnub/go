package messaging

import (
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
		return presenceSuffix
	} else {
		return subscribeMessage
	}
}
