package middlewares

import (
	"crypto/rsa"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/lib/parser"
	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/lib/token"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

// AuthorizationMiddleware
// 1. матчим route через trie
// 2. проверяем x-auth policy
// 3. разрешаем/запрещаем
func AuthorizationMiddleware(cfg *config.Config, log *logrus.Logger, specs *parser.ServiceSpecs, publicKey *rsa.PublicKey) fiber.Handler {

	return func(c fiber.Ctx) error {
		policy, serviceName, err := specs.FindPolicy(c.Method(), c.Path())
		if err != nil {
			return c.Next()
		}

		log.Printf("Service: %s, Public: %v, Permission: %s",
			serviceName, policy.Public, policy.Permission)
		if !policy.Public {
			if policy.JWT {
				// 1. Получаем токен из cookie
				accessToken := c.Cookies("accessToken")
				if accessToken == "" {
					return handleError(c, fiber.ErrUnauthorized, "Access token not found")
				}
				claims, err := token.ParseAccessToken(accessToken, publicKey)
				if err != nil {
					return handleError(c, fiber.ErrUnauthorized, "Invalid token claims")
				}
				// Добавляем к запросу заголовок с Id
				c.Request().Header.Add("X-User-Id", claims.ID)

				if policy.Permission == "admin" {
					if !claims.IsAdmin {
						return handleError(c, fiber.ErrForbidden, "access denied")
					}
				}
			}
		}
		return c.Next()
	}
}

// handleError обрабатывает ошибки аутентификации
func handleError(c fiber.Ctx, err *fiber.Error, message string) error {
	return c.Status(err.Code).JSON(fiber.Map{
		"error":   err.Message,
		"message": message,
		"status":  err,
	})
}
