package app

import (
	"crypto/tls"
	"fmt"

	tlsconfig "github.com/NekNB/CyberNavigate/backend/user-service/internal/app/tls"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/config"
	userAPI "github.com/NekNB/CyberNavigate/backend/user-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/middlewares"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/service/session"
	userService "github.com/NekNB/CyberNavigate/backend/user-service/internal/service/user"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage/postgres"
	"github.com/NekNB/CyberNavigate/swagger"
	"github.com/NekNB/CyberNavigate/swagger/gen/user"
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

	sessionService, err := session.New(cfg, log, postgresStorage)
	if err != nil {
		panic(err)
	}
	userService := userService.New(log, postgresStorage, sessionService)

	app := fiber.New(fiber.Config{
		AppName: "ArStore userServiceV1",
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middlewares.AuthMiddleware(sessionService))

	userAPI := userAPI.New(cfg, log, userService, sessionService)

	user.RegisterHandlersWithOptions(app.Name("API"), userAPI, user.FiberServerOptions{
		BaseURL: "/api/v1",
	})

	//Добавляем Specs директорию
	app.Get("/specs/*", static.New("", static.Config{
		FS:     swagger.SpecsFS,
		Browse: true,
	}))

	// Добавляем swagger директорию
	app.Get("/swagger/*", static.New("/swagger", static.Config{
		FS:     assets.SwaggerUI,
		Browse: true,
	}))
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusPermanentRedirect).Redirect().To("/swagger")
	})

	return &Server{log: log, app: app, cfg: cfg, TLSCfg: TLSCfg}, nil
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
