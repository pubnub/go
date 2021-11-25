package contract

import (
	"context"
	pubnub "github.com/pubnub/go/v6"
)

func iHaveAKeysetWithAccessManagerEnabled(ctx context.Context) error {
	state := getCommonState(ctx)
	config := pubnub.NewConfig()
	config.PublishKey = state.contractTestConfig.publishKey
	config.SubscribeKey = state.contractTestConfig.subscribeKey
	config.SecretKey = state.contractTestConfig.secretKey
	config.Origin = state.contractTestConfig.hostPort
	config.Secure = state.contractTestConfig.secure
	state.pubNub = pubnub.NewPubNub(config)
	return nil
}

func iHaveAKeysetWithAccessManagerEnabledWithoutSecretKey(ctx context.Context) error {
	state := getCommonState(ctx)
	config := pubnub.NewConfig()
	config.PublishKey = state.contractTestConfig.publishKey
	config.SubscribeKey = state.contractTestConfig.subscribeKey
	config.Origin = state.contractTestConfig.hostPort
	config.Secure = state.contractTestConfig.secure
	state.pubNub = pubnub.NewPubNub(config)
	return nil
}
