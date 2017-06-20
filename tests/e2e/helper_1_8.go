// +build go1.8

package e2e

import (
	"context"
	"time"
)

const (
	ERR_CONTEXT_CANCELLED = "context deadline exceeded"
	ERR_CONTEXT_DEADLINE  = "context deadline exceeded"
)

func contextWithTimeout(parent context.Context,
	timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

var backgroundContext = context.Background()
