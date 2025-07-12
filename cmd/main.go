package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"marketflow/config"
	"marketflow/internal/app"
	"marketflow/pkg/logger"
)

func Run() {
	ctx := context.Background()

	portFlag := flag.Int("port", 0, "Port number")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		fmt.Println(`Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number
  --help       Show help message`)
		os.Exit(0)
	}

	// Init logger
	log := logger.InitLogger(ctx, logger.LevelDebug)

	// Init config
	cfg, err := config.New()
	if err != nil {
		log.Error(ctx, "failed to init config", "error", err)
		return
	}

	if *portFlag != 0 {
		cfg.Server.HTTPServer.Port = *portFlag
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
