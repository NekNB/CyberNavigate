package middlewares

import (
	"strings"

	"github.com/NekNB/CyberNavigate/backend/user-service/internal/config"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/services/session"
	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware создает middleware для проверки JWT подписи
func AuthMiddleware(sessionService *session.SessionService) fiber.Handler {

	return func(c fiber.Ctx) error {
		// 1. Получаем токен из cookie
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
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

		// 2. Парсим и проверяем подпись
		claims, err := sessionService.ParseAndVerifyAccessToken(tokenString)
		if err != nil {
			return handleError(c, fiber.ErrUnauthorized, "Invalid token claims")
		}
		c.Locals(config.UserIdPayload, claims.UserID)
		c.Locals(config.IsAdminPayload, claims.IsAdmin)
		c.Locals(config.SessionIdPayload, claims.SessionID)

		// 3. Продолжаем выполнение
		return c.Next()
	}
}

// handleError обрабатывает ошибки аутентификации
func handleError(c fiber.Ctx, err *fiber.Error, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "unauthorized",
		"message": message,
		"status":  fiber.StatusUnauthorized,
	})
}
