package app

import (
	"fmt"

	"github.com/NekNB/CyberNavigate/backend/article-service/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/config"
	articleAPI "github.com/NekNB/CyberNavigate/backend/article-service/internal/http"
	articleService "github.com/NekNB/CyberNavigate/backend/article-service/internal/services/article"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/storage/mongo"
	"github.com/NekNB/CyberNavigate/backend/article-service/internal/storage/postgres"
	"github.com/NekNB/CyberNavigate/swagger"
	"github.com/NekNB/CyberNavigate/swagger/gen/article"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/sirupsen/logrus"
)

// Здесь реализуются методы LifeSpan сервера

type Server struct {
	cfg *config.Config
	app *fiber.App
}

func New(cfg *config.Config, log *logrus.Logger) (*Server, error) {

	app := fiber.New(fiber.Config{
		AppName: "ArStore ArticleServiceV1",
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:9000",
			"http://localhost:8000",
			"http://127.0.0.1:9000",
			"http://127.0.0.1:8000",
			"http://cyber-navigate_gateway-server:9000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "Content-Encoding",
		},
		AllowCredentials: true,
	}))

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

	return &Server{app: app, cfg: cfg}, nil
}

func (s *Server) Start() {
	s.app.Listen(
		fmt.Sprintf("0.0.0.0:%d", s.cfg.HTTP.Port),
	)
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
