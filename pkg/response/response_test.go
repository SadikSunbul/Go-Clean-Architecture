package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_JSON(t *testing.T) {
	app := fiber.New()

	// Test için bir Fiber context oluştur
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	// Test verisi
	data := map[string]string{"key": "value"}

	// Test: JSON fonksiyonunu çağır
	JSON(c, 200, data)

	// Assert: Beklenen sonuçları kontrol et
	assert.Equal(t, 200, c.Response().StatusCode())
	assert.Contains(t, string(c.Response().Body()), `"result":{"key":"value"}`)
}

func Test_JSON_Success(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		JSON(c, http.StatusOK, fiber.Map{
			"message": "success",
			"data":    []string{"item1", "item2"},
		})
		return nil
	})

	// Act
	resp, err := app.Test(httptest.NewRequest("GET", "/test", nil))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Result struct {
			Message string   `json:"message"`
			Data    []string `json:"data"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "success", result.Result.Message)
	assert.Equal(t, []string{"item1", "item2"}, result.Result.Data)
}

func Test_JSON_EmptyResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/empty", func(c *fiber.Ctx) error {
		JSON(c, http.StatusOK, nil)
		return nil
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/empty", nil))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_JSON_ErrorResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/error", func(c *fiber.Ctx) error {
		JSON(c, http.StatusBadRequest, fiber.Map{
			"error": "invalid request",
		})
		return nil
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/error", nil))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_Error_ValidationError(t *testing.T) {
	app := fiber.New()
	app.Post("/validate", func(c *fiber.Ctx) error {
		Error(c, http.StatusBadRequest, errors.New("invalid input"), "invalid input")
		return nil
	})

	resp, err := app.Test(httptest.NewRequest("POST", "/validate", nil))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_Error_NotFoundError(t *testing.T) {
	app := fiber.New()
	app.Get("/notfound", func(c *fiber.Ctx) error {
		Error(c, http.StatusNotFound, errors.New("resource not found"), "resource not found")
		return nil
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/notfound", nil))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_Error_InternalError(t *testing.T) {
	app := fiber.New()
	app.Get("/internal", func(c *fiber.Ctx) error {
		Error(c, http.StatusInternalServerError, errors.New("internal server error"), "internal server error")
		return nil
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/internal", nil))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Error_CustomError(t *testing.T) {
	app := fiber.New()
	app.Get("/custom", func(c *fiber.Ctx) error {
		Error(c, http.StatusTeapot, errors.New("I'm a teapot"), "I'm a teapot")
		return nil
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/custom", nil))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}
