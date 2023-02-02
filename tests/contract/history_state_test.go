package contract

import (
	"context"
	pubnub "github.com/pubnub/go/v7"
)

type historyStateKey struct{}

type historyState struct {
	fetchResponse *pubnub.FetchResponse
}

func newHistoryState() *historyState {
	return &historyState{}
}

func getHistoryState(ctx context.Context) *historyState {
	return ctx.Value(historyStateKey{}).(*historyState)
}
