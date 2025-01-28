package response

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_Error(t *testing.T) {
	app := fiber.New()

	// Test için bir Fiber context oluştur
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	// Hata mesajı
	err := errors.New("test error")
	message := "an error occurred"

	// Test: Hata fonksiyonunu çağır
	Error(c, 400, err, message)

	// Assert: Beklenen sonuçları kontrol et
	assert.Equal(t, 400, c.Response().StatusCode())
	assert.Contains(t, string(c.Response().Body()), message)
	assert.Contains(t, string(c.Response().Body()), "test error")
}
