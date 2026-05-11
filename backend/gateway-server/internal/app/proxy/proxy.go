package proxy

import (
	"fmt"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func Register(cfg *config.Config, router fiber.Router) {
	services := config.Normalize(cfg)
	for _, service := range services {
		router.All(
			service.Cfg.Path+"/*",
			func(c fiber.Ctx) error {
				target := fmt.Sprintf("%s://%s:%d",
					service.Cfg.Protocol,
					service.Cfg.Host,
					service.Cfg.Port,
				)

				fmt.Println(c.Request())

				path := c.Params("*")
				return proxy.Do(c, target+service.Cfg.Path+path)
			})

	}
}
