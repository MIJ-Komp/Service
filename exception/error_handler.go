package exception

import (
	"errors"
	"net/http"

	"api.mijkomp.com/helpers/logger"
	"api.mijkomp.com/models/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

// getRequestID retrieves the request ID from the context
func getRequestID(c *fiber.Ctx) string {
	if requestID := c.Locals("requestId"); requestID != nil {
		return requestID.(string)
	}
	return ""
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	requestID := getRequestID(ctx)

	if _, ok := err.(LoginError); ok {
		logger.LogWarningWithData("Login error occurred", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
			"ip":         ctx.IP(),
			"path":       ctx.Path(),
		})
		return ctx.Status(http.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	// if  err == gorm.ErrRecordNotFound {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.LogWarningWithData("Record not found in database", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
			"ip":         ctx.IP(),
			"path":       ctx.Path(),
		})
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	if _, ok := err.(NotFoundError); ok {
		logger.LogWarningWithData("Resource not found", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
			"ip":         ctx.IP(),
			"path":       ctx.Path(),
		})
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		if fiberErr.Code == fiber.StatusNotFound {
			logger.LogWarningWithData("Fiber route not found", map[string]interface{}{
				"request_id":  requestID,
				"error":       err.Error(),
				"status_code": fiberErr.Code,
				"ip":          ctx.IP(),
				"path":        ctx.Path(),
				"method":      ctx.Method(),
			})
			return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
		}
	}

	if _, ok := err.(ValidationError); ok {
		logger.LogWarningWithData("Validation error occurred", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
			"ip":         ctx.IP(),
			"path":       ctx.Path(),
			"method":     ctx.Method(),
		})
		return ctx.Status(http.StatusBadRequest).JSON(response.NewErrorWebResponse(400, "BAD_REQUEST", err.Error()))
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		errMessages := ParseValidationError(err)
		logger.LogWarningWithData("Validator validation errors", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
			"messages":   errMessages,
			"ip":         ctx.IP(),
			"path":       ctx.Path(),
			"method":     ctx.Method(),
		})

		return ctx.Status(http.StatusBadRequest).JSON(response.NewErrorWebResponse(400, "BAD_REQUEST", errMessages...))
	}

	if _, ok := err.(ForbiddenError); ok {
		logger.LogWarningWithData("Access forbidden", map[string]interface{}{
			"request_id": requestID,
			"error":      err.Error(),
			"ip":         ctx.IP(),
			"path":       ctx.Path(),
			"method":     ctx.Method(),
		})
		return ctx.Status(http.StatusForbidden).JSON(response.NewErrorWebResponse(403, "FORBIDDEN", err.Error()))
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		if fiberErr.Code == http.StatusMethodNotAllowed {
			logger.LogWarningWithData("Method not allowed", map[string]interface{}{
				"request_id":  requestID,
				"error":       fiberErr.Error(),
				"status_code": fiberErr.Code,
				"ip":          ctx.IP(),
				"path":        ctx.Path(),
				"method":      ctx.Method(),
			})
			return ctx.Status(http.StatusMethodNotAllowed).JSON(response.NewErrorWebResponse(405, "METHOD_NOT_ALLOWED", fiberErr.Error()))
		}
	}

	// Log internal server error with detailed information
	logger.LogErrorWithData("Internal server error occurred", map[string]interface{}{
		"request_id": requestID,
		"error":      err.Error(),
		"ip":         ctx.IP(),
		"path":       ctx.Path(),
		"method":     ctx.Method(),
		"headers":    ctx.GetReqHeaders(),
	})

	return ctx.Status(http.StatusInternalServerError).JSON(response.NewErrorWebResponse(500, "INTERNAL_SERVER_ERROR", "Maaf terjadi kesalahan pada server."))
}

func ParseValidationError(err error) []string {
	messages := []string{}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			// fmt.Println(e.Namespace())
			// fmt.Println(e.Field())
			// fmt.Println(e.StructNamespace())
			// fmt.Println(e.StructField())
			// fmt.Println(e.Tag())
			// fmt.Println(e.ActualTag())
			// fmt.Println(e.Kind())
			// fmt.Println(e.Type())
			// fmt.Println(e.Value())
			// fmt.Println(e.Param())
			// fmt.Println()

			messages = append(messages, e.Field()+" "+e.Tag()+" "+e.Param()+" ")
		}
	}

	return messages
}
