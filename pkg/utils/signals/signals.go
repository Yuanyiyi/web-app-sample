package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	sigkillCtx *context.Context
)

// NewSigKillContext returns a Context that cancels when os.Interrupt or os.Kill is received
func newSigKillContext() *context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	return &ctx
}

func GetSigKillContext() *context.Context {
	if sigkillCtx == nil {
		sigkillCtx = newSigKillContext()
	}
	return sigkillCtx
}
