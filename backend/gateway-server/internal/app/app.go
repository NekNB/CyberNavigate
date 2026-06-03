package app

import (
	"fmt"
	"net/http"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/proxy"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/assets"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/lib/parser"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/lib/token"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/middlewares"

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
	log *logrus.Logger
}

func New(cfg *config.Config, log *logrus.Logger) (*Server, error) {
	publicKey, err := token.GetPublicKeyFromFile(cfg.PublicKeyPath)
	if err != nil {
		panic(err)
	}
	// Загружаем спецификации из внешнего пакета
	// rootDir = "" если файлы в корне FS
	specs, err := parser.LoadSpecsFromFS(
		swagger.SpecsFS, // embed.FS из внешнего пакета
		"/api/v1",       // префикс для путей
		"docs",          // корневая директория в FS
	)
	if err != nil {
		log.Fatalf("Failed to load specs: %v", err)
	}

	// 2. Вывести политики по сервисам
	fmt.Println("\n=== POLICIES BY SERVICE ===")
	specs.PrintPoliciesByService()

	app := fiber.New(fiber.Config{
		AppName:     "Gateway Server",
		ProxyHeader: fiber.HeaderXForwardedHost,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	// 3. auth middleware
	app.Use(middlewares.AuthorizationMiddleware(cfg, log, specs, publicKey))

	// app.Use(func(c fiber.Ctx) error {
	// 	log.Info(c.Request())
	// 	return c.Next()
	// })
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

	proxy.Register(cfg, log, app.Name("GateWay"))

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

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	return &Server{app: app, cfg: cfg, log: log}, nil
}

func (s *Server) Run() {
	socket := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)

	s.log.Infof("gateway on %s", socket)
	if err := s.app.Listen(socket); err != nil && err != http.ErrServerClosed {
		s.log.Error(err)

		panic(err)
	}
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
