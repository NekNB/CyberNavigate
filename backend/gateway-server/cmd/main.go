package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log, err := logger.Init(cfg.Env)
	if err != nil {
		panic(err)
	}

	app, err := app.New(cfg, log)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.WithField("signal", sign.String()).Info("stopping application")

	app.Stop()

	log.Info("application stopped")
}
