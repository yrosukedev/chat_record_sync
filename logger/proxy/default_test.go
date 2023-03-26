package proxy

import (
	"context"
	"io"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	ctx := context.Background()
	proxy := NewLoggerProxy(LogLevelDebug, NewDefaultLogger())
	proxy.Debug(ctx, "[xxx] xx yyy xxxx, err: %s", "io.ErrShortBuffer")
	proxy.Info(ctx, "[xxx] xx yyy xxxx, err: %v", io.ErrShortBuffer)
	proxy.Warn(ctx, "[xxx] xx yyy xxxx, err: %#v", io.ErrShortBuffer)
	proxy.Error(ctx, "[xxx] xx yyy xxxx, err: %s", io.ErrShortBuffer)
}
