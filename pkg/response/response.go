package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func JSON(c *fiber.Ctx, status int, data interface{}) {
	c.Status(status).JSON(Response{
		Result: data,
	})
}
