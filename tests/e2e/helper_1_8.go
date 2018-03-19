// +build go1.8

package e2e

import (
	"context"
	"time"
)

const (
	ERR_CONTEXT_CANCELLED = "context canceled"
	ERR_CONTEXT_DEADLINE  = "context deadline exceeded"
)

func contextWithTimeout(parent context.Context,
	timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

func contextWithContext(parent context.Context) (
	context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

var backgroundContext = context.Background()
