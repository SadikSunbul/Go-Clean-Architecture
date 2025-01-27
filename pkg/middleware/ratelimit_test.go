package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	cfg, err := config.LoadConfig("../config/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	app := fiber.New()
	app.Use(RateLimiter())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 51 istek at (limit: 50)
	for i := 0; i < 49; i++ {
		resp, err := app.Test(httptest.NewRequest("GET", "/test", nil))
		assert.NoError(t, err)

		if i < 50 {
			// İlk 50 istek başarılı olmalı
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		} else {
			// 51. istek reddedilmeli
			assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
		}
	}
}
