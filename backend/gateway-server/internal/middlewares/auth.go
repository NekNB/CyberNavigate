package middlewares

import (
	"crypto/rsa"
	"strings"

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
			// 1. Получаем токен из cookie
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return handleError(c, fiber.ErrUnauthorized, "Invalid token claims")
			}

			// Проверяем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return handleError(c, fiber.ErrUnauthorized, "Invalid authorization header format. Use 'Bearer <token>'")
			}

			tokenString := parts[1]
			if tokenString == "" {
				c.Next()
				return handleError(c, fiber.ErrUnauthorized, "Token not found in cookie")
			}
			claims, err := token.ParseAccessToken(tokenString, publicKey)
			if err != nil {
				return handleError(c, fiber.ErrUnauthorized, "Invalid token claims")
			}

			if policy.Permission == "admin" {
				if !claims.IsAdmin {
					return handleError(c, fiber.ErrForbidden, "access denied")
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
