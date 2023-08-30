package middlewares

import (
	"os"

	"github.com/gofiber/fiber"
)

var API_KEY = os.Getenv("API_KEY")

func ApiKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) {
		apiKey := c.Get("x-api-key")
		if apiKey != API_KEY {
			c.SendStatus(fiber.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
