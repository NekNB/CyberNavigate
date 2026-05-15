package app

import (
	"fmt"

	"github.com/NekNB/CyberNavigate/backend/user-service/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/config"
	userAPI "github.com/NekNB/CyberNavigate/backend/user-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/middlewares"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/services/session"
	userService "github.com/NekNB/CyberNavigate/backend/user-service/internal/services/user"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage/postgres"
	"github.com/NekNB/CyberNavigate/swagger"
	"github.com/NekNB/CyberNavigate/swagger/gen/user"
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
	// app.Use(encryptcookie.New(encryptcookie.Config{
	// 	Key: encryptcookie.GenerateKey(16),
	// }))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:9000",
			"http://127.0.0.1:9000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "Content-Encoding",
		},
		AllowCredentials: true,
	}))

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
