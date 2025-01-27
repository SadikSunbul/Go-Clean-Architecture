package http

import (
	"net/http/httptest"
	"testing"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/middleware"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock yapıları
type MockDB struct {
	mock.Mock
	db.IDataBase
}

type MockRedis struct {
	mock.Mock
	redis.IRedis
}

type MockValidator struct {
	mock.Mock
	validation.Validation
}

func TestNewFiberServer(t *testing.T) {
	// Arrange
	mockDB := new(MockDB)
	mockRedis := new(MockRedis)
	mockValidator := new(MockValidator)
	cfg := &config.Config{
		Fiber: config.Fiber{
			Port: 3000,
			RateLimit: config.RateLimit{
				Max:        50,
				Expiration: 1,
			},
		},
	}

	// Act
	server := NewFiberServer(mockDB, cfg, mockValidator, mockRedis)

	// Assert
	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.cfg)
	assert.Equal(t, mockDB, server.db)
	assert.Equal(t, mockRedis, server.cache)
	assert.Equal(t, mockValidator, server.validator)
}

func TestFiberServer_GetApp(t *testing.T) {
	// Arrange
	server := &FiberServer{
		app: fiber.New(),
	}

	// Act
	app := server.GetApp()

	// Assert
	assert.NotNil(t, app)
	assert.Equal(t, server.app, app)
}

func TestFiberServer_HealthCheck(t *testing.T) {
	// Arrange
	mockDB := new(MockDB)
	mockRedis := new(MockRedis)
	mockValidator := new(MockValidator)
	cfg := &config.Config{
		Fiber: config.Fiber{
			Port: 3000,
		},
	}
	server := NewFiberServer(mockDB, cfg, mockValidator, mockRedis)
	server.app = fiber.New()

	// Health check endpoint'ini ekle
	server.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Act
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := server.app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestFiberServer_RateLimiter(t *testing.T) {
	// Arrange
	mockDB := new(MockDB)
	mockRedis := new(MockRedis)
	mockValidator := new(MockValidator)
	cfg := &config.Config{
		Fiber: config.Fiber{
			Port: 3000,
			RateLimit: config.RateLimit{
				Max:        2, // 2 istek limiti
				Expiration: 1,
			},
		},
	}
	server := NewFiberServer(mockDB, cfg, mockValidator, mockRedis)
	server.app = fiber.New()

	// Rate limiter ekle
	server.app.Use(middleware.RateLimiter(cfg))
	// Test endpoint'i ekle
	server.app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Rate limiter ekle
	server.app.Use(middleware.RateLimiter(cfg))

	// Act & Assert
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1") // IP'yi header üzerinden belirt

	// İlk 2 istek başarılı olmalı (limit: 2)
	for i := 0; i < 2; i++ {
		resp, err := server.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}

	// 3. istek reddedilmeli (limit aşıldı)
	respp, err := server.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusTooManyRequests, respp.StatusCode)
}
