package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a new UUID for this request
		requestID := uuid.New().String()
		
		// Set the request ID in the context
		c.Locals("requestId", requestID)
		
		// Add the request ID to the response header
		c.Set("X-Request-ID", requestID)
		
		return c.Next()
	}
}