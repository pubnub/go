package pubnub

type PNConfiguration struct {
	PublishKey   string
	SubscribeKey string
	SecretKey    string
}

func NewPNConfigurationDemo() *PNConfiguration {
	return &PNConfiguration{"demo", "demo", "demo"}
}

func NewPNConfiguration() *PNConfiguration {
	return &PNConfiguration{}
}
