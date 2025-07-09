package cmd

import (
	"context"
	"marketflow/internal/adapter/exchange"
	"marketflow/pkg/logger"
	"time"
)

func Run_health() {
	ctx := context.Background()
	log := logger.InitLogger(ctx, logger.LevelDebug)
	exClient := exchange.NewExchange("exchange1", "localhost:40101", log)
	for {
		ok, err := exClient.Health(ctx)
		if err != nil {
			log.Error(ctx, "service not available", "err", err)
			return
		}
		log.Info(ctx, "info", "name", exClient.Name, "ok", ok)
		time.Sleep(time.Second * 3)
	}
}
