// +build go1.7

package pubnub

import (
	"context"
)

func contextWithCancel(parent context.Context) (
	context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

var backgroundContext = context.Background()
