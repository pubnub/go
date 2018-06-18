// +build !go1.8,go1.7

package e2e

import (
	"context"
	"time"
)

const (
	ERR_CONTEXT_CANCELLED = "request canceled"
	ERR_CONTEXT_DEADLINE  = "request canceled"
)

func contextWithTimeout(parent context.Context,
	timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

var backgroundContext = context.Background()
