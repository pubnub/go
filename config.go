package pubnub

import (
	"fmt"
	"github.com/pubnub/go/utils"
	"log"
)

const (
	presenceTimeout = 300
)

type Config struct {
	PublishKey   string
	SubscribeKey string
	SecretKey    string
	AuthKey      string
	Origin       string
	Uuid         string
	CipherKey    string
	Secure       bool

	// net.Dialer.Timeout
	ConnectTimeout int

	// http.Client.Timeout for non-subscribe requests
	NonSubscribeRequestTimeout int

	// http.Client.Timeout for subscribe requests only
	SubscribeRequestTimeout int

	HeartbeatInterval int

	PresenceTimeout int

	MaximumReconnectionRetries int

	MaximumLatencyDataAge int

	FilterExpression string

	PNReconnectionPolicy ReconnectionPolicy

	Log *log.Logger

	SuppressLeaveEvents bool

	DisablePNOtherProcessing bool

	UseHttp2 bool
}

func NewDemoConfig() *Config {
	demoConfig := NewConfig()

	demoConfig.PublishKey = "demo"
	demoConfig.SubscribeKey = "demo"
	demoConfig.SecretKey = "demo"

	return demoConfig

}

func NewConfig() *Config {
	c := Config{
		Origin:                     "ps.pndsn.com",
		Secure:                     true,
		Uuid:                       fmt.Sprintf("pn-%s", utils.Uuid()),
		ConnectTimeout:             10,
		NonSubscribeRequestTimeout: 10,
		SubscribeRequestTimeout:    310,
		MaximumLatencyDataAge:      60,
		MaximumReconnectionRetries: 5,
		SuppressLeaveEvents:        false,
		DisablePNOtherProcessing:   false,
		PNReconnectionPolicy:       PNNonePolicy,
	}

	c.SetPresenceTimeout(presenceTimeout)

	return &c
}

func (c *Config) SetPresenceTimeoutWithCustomInterval(
	timeout, interval int) *Config {
	c.PresenceTimeout = timeout
	c.HeartbeatInterval = interval

	return c
}

func (c *Config) SetPresenceTimeout(timeout int) *Config {
	return c.SetPresenceTimeoutWithCustomInterval(timeout, (timeout/2)-1)
}
