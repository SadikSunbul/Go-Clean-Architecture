package middleware

import (
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter(cfg *config.Config) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.Fiber.RateLimit.Max,
		Expiration: time.Duration(cfg.Fiber.RateLimit.Expiration) * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // IP bazlı limitleme
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Çok fazla istek gönderdiniz. Lütfen 1 dakika bekleyin.",
				"limit": 50,
				"reset": time.Now().Add(1 * time.Minute).Unix(),
			})
		},
		Storage: nil, // Memory storage kullanılacak (default)
	})
}
