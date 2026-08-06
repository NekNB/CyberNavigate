package app

import (
	"crypto/tls"
	"fmt"

	tlsconfig "github.com/NekNB/CyberNavigate/backend/article-service/internal/app/tls"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/config"
	articleAPI "github.com/NekNB/CyberNavigate/backend/article-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/middleware"
	articleService "github.com/NekNB/CyberNavigate/backend/article-service/internal/services/article"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/storage/mongo"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/storage/postgres"
	"github.com/NekNB/CyberNavigate/swagger"
	"github.com/NekNB/CyberNavigate/swagger/gen/article"

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
		AppName: "ArStore ArticleServiceV1",
	})
	if cfg.Env != "local" {
		app.Use(recover.New())
	}
	app.Use(middleware.LoggerMiddleware())
	app.Use(logger.New())

	mongoStorage, err := mongo.CreateConnection(log, fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		cfg.Storage.Mongo.User,
		cfg.Storage.Mongo.Password,
		cfg.Storage.Mongo.Host,
		cfg.Storage.Mongo.Port,
	), cfg.Storage.Mongo.Database, cfg.Storage.Mongo.Collection)
	if err != nil {
		log.Error(err)
		return nil, err
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

	articleService := articleService.New(log, mongoStorage, postgresStorage)

	articleAPI := articleAPI.New(log, articleService)

	article.RegisterHandlersWithOptions(app.Name("API"), articleAPI, article.FiberServerOptions{
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

	if err := s.app.Listener(ln); err != nil {
		s.log.Error(err)
		panic(err)
	}
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
