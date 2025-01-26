package response

import (
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/gofiber/fiber/v2"
)

func Error(c *fiber.Ctx, status int, err error, message string) {
	cfg := config.GetConfig()
	errorRes := map[string]interface{}{
		"message": message,
	}

	if cfg.Environment != config.ProductionEnv {
		errorRes["debug"] = err.Error()
	}

	c.Status(status).JSON(Response{Error: errorRes})
}
