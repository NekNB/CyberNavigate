package app

import (
	"time"

	"github.com/NekNB/AuthService/sso/internal/services/auth"
	"github.com/NekNB/AuthService/sso/internal/storage/sqlite"
	grpcapp "github.com/NekNB/CyberNavigate/backend/article-service/internal/app/grpc"
	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *logrus.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
