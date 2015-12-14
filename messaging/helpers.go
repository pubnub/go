package messaging

import (
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
