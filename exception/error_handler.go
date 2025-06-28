package exception

import (
	"errors"
	"fmt"
	"net/http"

	"api.mijkomp.com/helpers/logger"
	"api.mijkomp.com/models/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	if _, ok := err.(LoginError); ok {
		logger.LogWarning(fmt.Sprintf("Login error: %s", err.Error()))
		return ctx.Status(http.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	// if  err == gorm.ErrRecordNotFound {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.LogWarning(fmt.Sprintf("Record not found error: %s", err.Error()))
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	if _, ok := err.(NotFoundError); ok {
		logger.LogWarning(fmt.Sprintf("Not found error: %s", err.Error()))
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		if fiberErr.Code == fiber.StatusNotFound {
			logger.LogWarning(fmt.Sprintf("Fiber not found error: %s", err.Error()))
			return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
		}
	}

	if _, ok := err.(ValidationError); ok {
		logger.LogWarning(fmt.Sprintf("Validation error: %s", err.Error()))
		return ctx.Status(http.StatusBadRequest).JSON(response.NewErrorWebResponse(400, "BAD_REQUEST", err.Error()))
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		logger.LogWarning(fmt.Sprintf("Validator validation errors: %s", err.Error()))
		errMessages := ParseValidationError(err)

		return ctx.Status(http.StatusBadRequest).JSON(response.NewErrorWebResponse(400, "BAD_REQUEST", errMessages...))
	}

	if _, ok := err.(ForbiddenError); ok {
		logger.LogWarning(fmt.Sprintf("Forbidden error: %s", err.Error()))
		return ctx.Status(http.StatusForbidden).JSON(response.NewErrorWebResponse(403, "FORBIDDEN", err.Error()))
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		if fiberErr.Code == http.StatusMethodNotAllowed {
			logger.LogWarning(fmt.Sprintf("Method not allowed error: %s", fiberErr.Error()))
			return ctx.Status(http.StatusMethodNotAllowed).JSON(response.NewErrorWebResponse(405, "METHOD_NOT_ALLOWED", fiberErr.Error()))
		}
	}

	fmt.Print(err.Error())
	logger.LogError(errors.New("Internal server error: " + err.Error()))

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
