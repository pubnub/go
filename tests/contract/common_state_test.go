package contract

import pubnub "github.com/pubnub/go/v6"

type commonStateKey struct{}

type commonState struct {
	contractTestConfig         contractTestConfig
	pubNub         *pubnub.PubNub
	err            error
	statusResponse pubnub.StatusResponse
}

func newCommonState(contractTestConfig contractTestConfig) *commonState {


	return &commonState{
		contractTestConfig: contractTestConfig,
	}
}
