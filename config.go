package pubnub

import "github.com/pubnub/go/utils"

type Config struct {
	PublishKey                 string
	SubscribeKey               string
	SecretKey                  string
	Origin                     string
	Uuid                       string
	Secure                     bool
	Crypto                     bool
	ConnectionTimeout          int
	NonSubscribeRequestTimeout int
	SubscribeRequestTimeout    int
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
		ConnectionTimeout:          10,
		NonSubscribeRequestTimeout: 10,
		SubscribeRequestTimeout:    10,
	}
}
