package pubnub

import "github.com/pubnub/go/utils"

type Config struct {
	PublishKey   string
	SubscribeKey string
	SecretKey    string
	Origin       string
	Uuid         string
	CipherKey    string
	Secure       bool

	// TODO: timeout assignment in http.Client/http.Transport should be
	// completely reviewed due language concerns.

	// net.Dialer.Timeout
	ConnectTimeout int

	// http.Client.Timeout for non-subscribe requests
	NonSubscribeRequestTimeout int

	// http.Client.Timeout for subscribe requests only
	SubscribeRequestTimeout int
}

func NewDemoConfig() *Config {
	demoConfig := NewConfig()

	demoConfig.PublishKey = "demo"
	demoConfig.SubscribeKey = "demo"
	demoConfig.SecretKey = "demo"

	return demoConfig
}

func NewConfig() *Config {
	return &Config{
		Origin:                     "ps.pndsn.com",
		Secure:                     true,
		Uuid:                       utils.Uuid(),
		ConnectTimeout:             10,
		NonSubscribeRequestTimeout: 10,
		SubscribeRequestTimeout:    310,
	}
}
