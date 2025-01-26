package middleware

import (
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/jtoken"
	"github.com/gofiber/fiber/v2"
)

// ::::::::::::::::::::::::::::::::::::
// 		AUTH JWT MIDDLEWARE
// ::::::::::::::::::::::::::::::::::::

func JWTAuth() func(c *fiber.Ctx) error {
	return JWT(jtoken.AccessTokenType)
}
func JWTRefreshAuth() func(c *fiber.Ctx) error {
	return JWT(jtoken.RefreshTokenType)
}

func JWT(tokenType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.GetReqHeaders()["Authorization"][0]
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token is empty",
			})
		}

		payload, err := jtoken.ValidateToken(token)
		if err != nil || payload == nil || payload["type"] != tokenType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token is invalid",
			})
		}

		c.Locals("userId", payload["id"])
		c.Locals("role", payload["role"])

		return c.Next()
	}
}
