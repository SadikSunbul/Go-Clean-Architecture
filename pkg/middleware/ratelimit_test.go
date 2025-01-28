package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	cfg, err := config.LoadConfig("../config/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	app := fiber.New()
	app.Use(RateLimiter(cfg))

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

func TestRateLimiter_Integration(t *testing.T) {
	// Arrange
	app := fiber.New()

	// Rate limiter config
	cfg := &config.Config{
		Fiber: config.Fiber{
			RateLimit: config.RateLimit{
				Max:        2,
				Expiration: 1,
			},
		},
	}

	// Rate limiter middleware'i ekle
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Fiber.RateLimit.Max,
		Expiration: time.Duration(cfg.Fiber.RateLimit.Expiration) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.ImmutableString("test-key") // Test için sabit key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Çok fazla istek gönderdiniz. Lütfen bekleyin.",
			})
		},
	}))

	// Test endpoint'i
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Act & Assert
	// İlk 2 istek başarılı olmalı
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}

	// 3. istek limit aşımı nedeniyle reddedilmeli
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
}
