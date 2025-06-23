package middleware

import (
	"os"
	"strconv"

	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuthMiddleware(userService service.UserService) func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		user, err := helpers.TokenValid(c)

		if err != nil {

			return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
				Code:    401,
				Status:  "Unauthorized, Token invalid",
				Content: "",
			})
		}

		// var userToken domain.UserToken
		userId, _ := strconv.ParseUint(user[0], 10, 64)
		guidToken, _ := uuid.Parse(user[1])

		hasToken := userService.HasToken(uint(userId), guidToken)
		// userData := userService.GetById(uint(userId))

		if !hasToken {
			return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
				Code:    401,
				Status:  "Unauthorized",
				Content: "",
			})
		}

		c.Locals("userId", strconv.FormatUint(uint64(userId), 10))

		return c.Next()
	}
}

func PaymentMiddlewareXendit() func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		callbackToken := c.Get("x-callback-token")
		webhookToken := os.Getenv("XENDIT_WEBHOOK_TOKEN")

		if callbackToken != webhookToken {
			return c.Status(fiber.StatusUnauthorized).JSON(response.WebResponse{
				Code:    401,
				Status:  "Unauthorized",
				Content: "",
			})
		}

		return c.Next()
	}
}
