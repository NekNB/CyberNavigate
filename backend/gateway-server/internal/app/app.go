package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/proxy"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/specs"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/swagger"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/middlewares"

	"github.com/sirupsen/logrus"
)

type App struct {
	ctx    context.Context
	cfg    *config.Config
	log    *logrus.Logger
	server *http.Server
	port   int
}

func New(ctx context.Context, cfg *config.Config, log *logrus.Logger) *App {

	return &App{
		ctx:  ctx,
		cfg:  cfg,
		log:  log,
		port: cfg.Server.Port,
	}
}

func (a *App) Run(ctx context.Context) {

	mux := http.NewServeMux()

	specsHandler := specs.NewSpecs()
	mux.Handle("/specs/", http.StripPrefix("/specs/", specsHandler))

	swaggerHandler := swagger.NewSwagger()
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", swaggerHandler))

	HTTPHandler, err := proxy.New(a.cfg)
	if err != nil {
		panic(err)
	}
	mux.Handle("/api/", HTTPHandler)

	mainHandler := middlewares.Recovery(a.log)(middlewares.RequestLogging(a.log)(mux))

	a.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", a.cfg.Server.Port),
		Handler: mainHandler,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	a.log.Infof("gateway on : %s", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Error(err)

		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	if a.server == nil {
		return nil
	}

	a.log.Infof("shutting down gateway...")

	return a.server.Shutdown(ctx)
}
