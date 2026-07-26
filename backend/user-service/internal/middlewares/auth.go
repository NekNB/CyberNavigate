package middlewares

import (
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/config"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/service/session"
	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware создает middleware для проверки JWT подписи
func AuthMiddleware(sessionService *session.SessionService) fiber.Handler {

	return func(c fiber.Ctx) error {
		// 1. Получаем токен из cookie
		accessToken := c.Cookies("accessToken")

		if accessToken != "" {
			// 2. Парсим и проверяем подпись
			claims, err := sessionService.ParseAndVerifyAccessToken(accessToken)
			if err != nil {
				return handleError(c, fiber.ErrUnauthorized, "Invalid token claims")
			}
			c.Locals(config.UserIdPayload, claims.UserID)
			c.Locals(config.IsAdminPayload, claims.IsAdmin)
			c.Locals(config.SessionIdPayload, claims.SessionID)
		}
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
