package pubnub

type PNConfiguration struct {
	PublishKey                 string
	SubscribeKey               string
	SecretKey                  string
	Origin                     string
	Secure                     bool
	ConnectionTimeout          int
	NonSubscribeRequestTimeout int
	SubscribeRequestTimeout    int
}

func NewPNConfigurationDemo() *PNConfiguration {
	return &PNConfiguration{
		PublishKey:                 "demo",
		SubscribeKey:               "demo",
		SecretKey:                  "demo",
		Origin:                     "ps.pndsn.com",
		Secure:                     true,
		ConnectionTimeout:          2000,
		NonSubscribeRequestTimeout: 2000,
		SubscribeRequestTimeout:    2000,
	}
}

func NewPNConfiguration() *PNConfiguration {
	return &PNConfiguration{
		Origin: "ps.pndsn.com",
		Secure: true,
	}
}
