package middleware

import (
	"api.mijkomp.com/helpers/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

// getRequestID retrieves the request ID from the context
func getRequestID(c *fiber.Ctx) string {
	if requestID := c.Locals("requestId"); requestID != nil {
		return requestID.(string)
	}
	return ""
}

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Catat waktu mulai
		start := time.Now()
		requestID := getRequestID(c)

		// Log incoming request
		logger.LogInfoWithData("Incoming request", map[string]interface{}{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       c.Path(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
			"query":      c.Queries(),
		})

		// Proses request
		err := c.Next()

		// Hitung durasi
		duration := time.Since(start)

		// Log request completion with enhanced data
		logger.LogAPIRequest(
			c.Method(),
			c.Path(),
			c.IP(),
			c.Response().StatusCode(),
			duration,
		)

		// Log additional data for debugging
		logger.LogDebugWithData("Request completed", map[string]interface{}{
			"request_id":   requestID,
			"method":       c.Method(),
			"path":         c.Path(),
			"status_code":  c.Response().StatusCode(),
			"duration_ms":  duration.Milliseconds(),
			"response_size": len(c.Response().Body()),
		})

		return err
	}
}