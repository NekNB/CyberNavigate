package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/NekNB/CyberNavigate/backend/article-service/internal/app"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/config"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/lib/logger"
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

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTtl)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}
