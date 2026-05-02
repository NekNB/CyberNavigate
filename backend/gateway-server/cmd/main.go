package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/lib/logger"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.MustLoad()

	log, err := logger.Init(cfg.Env)
	if err != nil {
		panic(err)
	}

	app := app.New(ctx, cfg, log)

	go app.Run(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.WithField("signal", sign.String()).Info("stopping application")

	app.Stop(ctx)

	log.Info("application stopped")
}
