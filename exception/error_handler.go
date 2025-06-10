package exception

import (
	"errors"
	"fmt"
	"net/http"

	"api.mijkomp.com/models/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	if _, ok := err.(LoginError); ok {
		return ctx.Status(http.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	// if  err == gorm.ErrRecordNotFound {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	if _, ok := err.(NotFoundError); ok {
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		// log.Println(err.Error())
		if fiberErr.Code == fiber.StatusNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(response.NewErrorWebResponse(404, "NOT_FOUND", err.Error()))
		}
	}

	if _, ok := err.(ValidationError); ok {
		return ctx.Status(http.StatusBadRequest).JSON(response.NewErrorWebResponse(400, "BAD_REQUEST", err.Error()))
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		errMessages := ParseValidationError(err)

		return ctx.Status(http.StatusBadRequest).JSON(response.NewErrorWebResponse(400, "BAD_REQUEST", errMessages...))
	}

	if _, ok := err.(ForbiddenError); ok {
		return ctx.Status(http.StatusForbidden).JSON(response.NewErrorWebResponse(403, http.StatusText(http.StatusForbidden), err.Error()))
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		if fiberErr.Code == http.StatusMethodNotAllowed {
			return ctx.Status(http.StatusMethodNotAllowed).JSON(response.NewErrorWebResponse(405, "METHOD_NOT_ALLOWED", fiberErr.Error()))
		}
	}

	// TODO: insert error ke log
	fmt.Println(err.Error())

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
