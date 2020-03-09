package main

import (
	"context"
)

func main() {
	cfg := NewConfig("./config/config.yaml")
	if cfg == nil {
		return
	}

	ctx := context.Background()
	state := NewState(cfg.StaticHosts)
	go state.Run(ctx, cfg.Mikrotic)
	httpState := NewHttpState(cfg.ListenAddress, state.req)
	go httpRun(ctx, httpState)

	<-ctx.Done()
}
