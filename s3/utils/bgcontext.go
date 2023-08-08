package utils

import (
	"context"
	"time"
)

// BgContext returns a context that can be used for async operations.
// Cancellation/timeouts are removed, so parent cancellations/timeout will
// not propagate from parent.
// Context values are preserved.
// This can be used for goroutines that live beyond the parent context.
func BgContext(parent context.Context) context.Context {
	return bgCtx{parent: parent}
}

type bgCtx struct {
	parent context.Context
}

func (a bgCtx) Done() <-chan struct{} {
	return nil
}

func (a bgCtx) Err() error {
	return nil
}

func (a bgCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (a bgCtx) Value(key interface{}) interface{} {
	return a.parent.Value(key)
}
