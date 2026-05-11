package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

func DetailedLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Получаем полный URL
		fullURL := c.BaseURL() + c.OriginalURL()

		fmt.Printf("\n┌─────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│ REQUEST\n")
		fmt.Printf("├─────────────────────────────────────────────────────────────┤\n")
		fmt.Printf("│ Method: %s\n", c.Method())
		fmt.Printf("│ URL:    %s\n", fullURL)
		fmt.Printf("│ Path:   %s\n", c.OriginalURL())
		fmt.Printf("│ Body:   %s\n", string(c.Body()))

		err := c.Next()

		fmt.Printf("├─────────────────────────────────────────────────────────────┤\n")
		fmt.Printf("│ RESPONSE\n")
		fmt.Printf("├─────────────────────────────────────────────────────────────┤\n")
		fmt.Printf("│ Status: %d\n", c.Response().StatusCode())
		fmt.Printf("│ Time:   %v\n", time.Since(start))
		fmt.Printf("│ Body:   %s\n", string(c.Response().Body()))
		fmt.Printf("└─────────────────────────────────────────────────────────────┘\n\n")

		return err
	}
}
