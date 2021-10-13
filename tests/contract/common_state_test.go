package contract

import pubnub "github.com/pubnub/go/v6"

type commonStateKey struct{}

type commonState struct {
	pubNub *pubnub.PubNub
	err    error
}

func newCommonState(contractTestConfig contractTestConfig) *commonState {
	config := pubnub.NewConfig()
	config.PublishKey = contractTestConfig.publishKey
	config.SubscribeKey = contractTestConfig.subscribeKey
	config.SecretKey = contractTestConfig.secretKey
	config.Origin = contractTestConfig.hostPort
	config.Secure = contractTestConfig.secure

	return &commonState{
		pubNub: pubnub.NewPubNub(config),
	}
}
