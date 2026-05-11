package app

import (
	"fmt"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/proxy"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/NekNB/CyberNavigate/swagger"
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
		AppName: "Gateway Server",
	})
	app.Use(logger.New())
	app.Use(recover.New())
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

	proxy.Register(cfg, app.Name("GateWay"))

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

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})

	return &Server{app: app, cfg: cfg}, nil
}

func (s *Server) Start() {
	s.app.Listen(
		fmt.Sprintf(":%d", s.cfg.Server.Port),
	)
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
