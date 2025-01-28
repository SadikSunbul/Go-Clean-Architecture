package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/dto"
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/service"
	"github.com/SadikSunbul/Go-Clean-Architecture/model/entity"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/jtoken"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

// PostHandler handles HTTP requests for posts
type PostHandler struct {
	service service.IPostService
	cache   redis.IRedis
	cfg     config.Config
}

// NewPostHandler creates a new post handler
func NewPostHandler(service service.IPostService, cache redis.IRedis) *PostHandler {
	return &PostHandler{
		service: service,
		cache:   cache,
		cfg:     *config.GetConfig(),
	}
}

// GetAllPosts godoc
// @Summary Tüm gönderileri listele
// @Description Sistemdeki tüm gönderileri getirir
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response{data=[]entity.Post}
// @Failure 500 {object} dto.Response
// @Router /posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	var posts *[]entity.Post
	cachekey := c.Request().URI().String()
	err := h.cache.Get(cachekey, &posts)
	if err == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   posts,
		})
	}

	posts, err = h.service.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	_ = h.cache.SetWithExpiration(cachekey, posts, time.Duration(h.cfg.Redis.ProductCachingTime))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   posts,
	})
}

// GetPostById godoc
// @Summary Gönderiyi ID ile getir
// @Description Belirtilen ID'ye sahip gönderiyi getirir
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} dto.Response{data=entity.Post}
// @Failure 404 {object} dto.Response
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostById(c *fiber.Ctx) error {
	var post entity.Post
	cachekey := c.Request().URI().String()
	err := h.cache.Get(cachekey, &post)
	if err == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   post,
		})
	}

	id := c.Params("id")
	post, err = h.service.GetById(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	_ = h.cache.SetWithExpiration(cachekey, post, time.Duration(h.cfg.Redis.ProductCachingTime))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   post,
	})
}

// CreatePost godoc
// @Summary Yeni gönderi oluştur
// @Description Yeni bir gönderi oluşturur
// @Tags posts
// @Accept json
// @Produce json
// @Param post body dto.PostDto true "Post bilgileri"
// @Success 201 {object} dto.Response{data=entity.Post}
// @Failure 400,500 {object} dto.Response
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var post dto.PostDto
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	reqpost, err := h.service.Create(&post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts"))
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   reqpost,
	})
}

// UpdatePost godoc
// @Summary Gönderiyi güncelle
// @Description Belirtilen ID'ye sahip gönderiyi günceller
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body dto.PostUpdateDto true "Güncellenecek post bilgileri"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400,404 {object} dto.Response
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post dto.PostUpdateDto
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	post.UpdatedAt = time.Now()
	count, err := h.service.Update(id, &post)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts/%s", id))
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts"))

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"update_count": count,
	})
}

// DeletePost godoc
// @Summary Gönderiyi sil
// @Description Belirtilen ID'ye sahip gönderiyi siler
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 404 {object} dto.Response
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts/%s", id))
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts"))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

// CreateJWT godoc
// @Summary JWT token oluştur
// @Description Verilen isim için JWT token oluşturur
// @Tags auth
// @Accept json
// @Produce json
// @Param name path string true "Kullanıcı adı"
// @Success 200 {object} dto.TokenResponse
// @Router /auth/token/{name} [get]
func (h *PostHandler) CreateJWT(c *fiber.Ctx) error {
	data := make(map[string]interface{})
	name := c.Params("name")
	data["name"] = name
	token := jtoken.GenerateAccessToken(data)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  token,
	})
}

// ValidateJWT godoc
// @Summary JWT token doğrula
// @Description JWT token'ı doğrular
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string
// @Failure 401 {object} dto.Response
// @Router /auth/validate [get]
func (h *PostHandler) ValidateJWT(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Hello, %s", c.Locals("name")))
}
