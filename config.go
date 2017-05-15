package pubnub

type Config struct {
	PublishKey                 string
	SubscribeKey               string
	SecretKey                  string
	Origin                     string
	Secure                     bool
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
		Origin: "ps.pndsn.com",
		Secure: true,
	}
}
