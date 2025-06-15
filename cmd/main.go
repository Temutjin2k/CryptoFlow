package main

import (
	"context"
	"marketflow/config"
	"marketflow/internal/app"
	"marketflow/pkg/logger"
)

func main() {
	ctx := context.Background()

	// Init logger
	log := logger.InitLogger(ctx, logger.LevelDebug)

	// Init config
	cfg, err := config.New()
	if err != nil {
		log.Error(ctx, "failed to init config", "error", err)
		return
	}

	// Printing the config
	config.PrintConfig(cfg)

	// Creating application
	app, err := app.NewApplication(ctx, cfg, log)
	if err != nil {
		log.Error(ctx, "failed to init application", "error", err)
		return
	}

	// Running the apllication
	err = app.Run()
	if err != nil {
		log.Error(ctx, "failed to run application", "error", err)
		return
	}
}
