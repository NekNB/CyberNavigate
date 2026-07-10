package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/app"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/config"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/lib/logger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log, err := logger.Init(cfg.Env)
	if err != nil {
		panic(err)
	}

	app, err := app.New(cfg, log)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	go app.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	app.Stop()

	log.Info("application stopped")
}
