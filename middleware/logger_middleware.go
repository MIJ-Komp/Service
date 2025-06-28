package middleware

import (
	"api.mijkomp.com/helpers/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Catat waktu mulai
		start := time.Now()

		// Proses request
		err := c.Next()

		// Hitung durasi
		duration := time.Since(start)

		// Log request
		logger.LogAPIRequest(
			c.Method(),
			c.Path(),
			c.IP(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}