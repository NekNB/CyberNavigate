package app

import (
	"crypto/tls"
	"fmt"

	tlsconfig "github.com/NekNB/CyberNavigate/backend/simulator-service/internal/app/tls"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/config"
	simulatorAPI "github.com/NekNB/CyberNavigate/backend/simulator-service/internal/http"
	simulatorService "github.com/NekNB/CyberNavigate/backend/simulator-service/internal/services/simulator"

	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/storage/postgres"
	"github.com/NekNB/CyberNavigate/swagger"
	"github.com/NekNB/CyberNavigate/swagger/gen/simulator"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/sirupsen/logrus"
)

// Здесь реализуются методы LifeSpan сервера

type Server struct {
	cfg    *config.Config
	app    *fiber.App
	log    *logrus.Logger
	TLSCfg *tls.Config
}

func New(cfg *config.Config, log *logrus.Logger) (*Server, error) {

	TLSCfg, err := tlsconfig.LoadTLSConfig(
		cfg.Certs.ServerCertPath,
		cfg.Certs.ServerKeyPath,
		cfg.Certs.CaCertPath,
	)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "ArStore simulatorServiceV1",
	})
	app.Use(logger.New())
	if cfg.Env != "local" {
		app.Use(recover.New())
	}

	postgresStorage, err := postgres.New(log, fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Storage.Postgres.User,
		cfg.Storage.Postgres.Password,
		cfg.Storage.Postgres.Host,
		cfg.Storage.Postgres.Port,
		cfg.Storage.Postgres.Database,
	))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	simulatorService := simulatorService.New(log, postgresStorage)

	simulatorAPI := simulatorAPI.New(log, simulatorService)

	simulator.RegisterHandlersWithOptions(app.Name("API"), simulatorAPI, simulator.FiberServerOptions{
		BaseURL: "/api/v1",
	})

	//Добавляем Specs директорию
	app.Get("/specs/*", static.New("", static.Config{
		FS:     swagger.SpecsFS,
		Browse: true,
	}))

	app.Get("/redoc/*", static.New("/redoc", static.Config{
		FS: assets.RedocUI,
	}))
	// Добавляем swagger директорию
	app.Get("/swagger/*", static.New("/swagger", static.Config{
		FS:     assets.SwaggerUI,
		Browse: true,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusPermanentRedirect).Redirect().To("/swagger")
	})

	return &Server{app: app, cfg: cfg, TLSCfg: TLSCfg, log: log}, nil
}

func (s *Server) HTTPSStart() {
	ln, err := tls.Listen(
		"tcp",
		fmt.Sprintf("0.0.0.0:%d", s.cfg.HTTP.Port),
		s.TLSCfg,
	)
	if err != nil {
		panic(err)
	}

	panic(s.app.Listener(ln))
}

func (s *Server) HTTPStart() {
	s.log.Debug("Start")
	if err := s.app.Listen(fmt.Sprintf("127.0.0.1:%d", s.cfg.HTTP.Port), fiber.ListenConfig{
		ShutdownTimeout: s.cfg.HTTP.Timeout,
	}); err != nil {
		s.log.Error(err)

		panic(err)
	}
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
