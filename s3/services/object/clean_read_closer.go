package object

import (
	"context"
	"io"
	"time"
)

func WrapCleanReadCloser(rc io.ReadCloser, timeout time.Duration, afterCloseHooks ...func()) io.ReadCloser {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	arc := &cleanReadCloser{
		rc:     rc,
		cancel: cancel,
	}
	go func() {
		<-ctx.Done()
		_ = rc.Close()
		// call after hooks by stack order
		for len(afterCloseHooks) > 0 {
			idx := len(afterCloseHooks) - 1
			f := afterCloseHooks[idx]
			f()
			afterCloseHooks = afterCloseHooks[:idx]
		}
	}()
	return arc
}

type cleanReadCloser struct {
	rc     io.ReadCloser
	cancel context.CancelFunc
}

func (h *cleanReadCloser) Read(p []byte) (n int, err error) {
	return h.rc.Read(p)
}

func (h *cleanReadCloser) Close() error {
	defer h.cancel()
	return h.rc.Close()
}
