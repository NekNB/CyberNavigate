package app

import (
	"crypto/tls"
	"fmt"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/proxy"
	tlsconfig "github.com/NekNB/CyberNavigate/backend/gateway-server/internal/app/tls"
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
	cfg    *config.Config
	app    *fiber.App
	log    *logrus.Logger
	TLSCfg *tls.Config
}

func New(cfg *config.Config, log *logrus.Logger) (*Server, error) {

	TLSCfg, err := tlsconfig.LoadTLSConfig(
		cfg.Certs.PublicCertPath,
		cfg.Certs.PublicKeyPath,
		cfg.Certs.CaCertPath,
	)
	if err != nil {
		panic(err)
	}

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
	specs.PrintPoliciesByService()

	app := fiber.New(fiber.Config{
		AppName:     "Gateway Server",
		ProxyHeader: fiber.HeaderXForwardedHost,
	})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://localhost:443",
			"https://127.0.0.1:443",
			"https://localhost:9443",
			"https://127.0.0.1:9443",
			"http://localhost:80",
			"http://127.0.0.1:80",
			"http://localhost:9080",
			"http://127.0.0.1:9080",
			"http://127.0.0.1:3777",
			"http://localhost:3777",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "Content-Encoding",
		},
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	// 3. auth middleware
	app.Use(middlewares.AuthorizationMiddleware(cfg, log, specs, publicKey))

	proxy.Register(cfg, log, app.Name("GateWay"))

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
		FS: assets.SwaggerUI,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusPermanentRedirect).Redirect().To("/swagger")
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	return &Server{app: app, cfg: cfg, log: log, TLSCfg: TLSCfg}, nil
}

func (s *Server) HTTPStart() {
	socket := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.HTTPPort)

	if err := s.app.Listen(socket); err != nil {
		s.log.Error(err)
		panic(err)
	}
}

func (s *Server) HTTPSStart() {
	socket := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.HTTPSPort)
	ln, err := tls.Listen(
		"tcp",
		socket,
		s.TLSCfg,
	)
	if err != nil {
		s.log.Error(err)
		panic(err)
	}

	if err := s.app.Listener(ln); err != nil {
		s.log.Error(err)
		panic(err)
	}
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
