package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func SetUserId() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		userId, ok := c.GetHeaders()["X-User-Id"]
		if ok && len(userId) == 1 {
			c.Locals("userId", userId[0])
		}
		return
	}
}
