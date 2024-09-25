package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request is processed
		log.Printf(
			"%s - [%s] %s %s %d %dms",
			c.IP(),
			time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start).Milliseconds(),
		)

		return err
	}
}
