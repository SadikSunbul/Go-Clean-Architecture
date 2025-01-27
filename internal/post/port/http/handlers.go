package http

import (
	"fmt"
	"github.com/SadikSunbul/Go-Clean-Architecture/model/entity"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/jtoken"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/redis"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/response"
	"net/http"
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/dto"
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/service"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service service.IPostService
	cache   redis.IRedis
	cfg     config.Config
}

func NewPostHandler(service service.IPostService, cache redis.IRedis) *PostHandler {
	return &PostHandler{
		service: service,
		cache:   cache,
		cfg:     *config.GetConfig(),
	}
}

func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {

	var posts *[]entity.Post
	cachekey := c.Request().URI().String()
	err := h.cache.Get(cachekey, &posts) // rediste varmı yokmu kontrol
	if err == nil {
		response.JSON(c, http.StatusOK, posts)
		return nil
	}

	posts, err = h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	_ = h.cache.SetWithExpiration(cachekey, posts, time.Duration(h.cfg.Redis.ProductCachingTime)) // redise veriyi ekleme
	return c.JSON(posts)
}

func (h *PostHandler) GetPostById(c *fiber.Ctx) error {

	var post entity.Post
	cachekey := c.Request().URI().String()
	err := h.cache.Get(cachekey, &post) // rediste varmı yokmu kontrol
	if err == nil {
		response.JSON(c, http.StatusOK, post)
		return nil
	}

	id := c.Params("id")
	post, err = h.service.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	_ = h.cache.SetWithExpiration(cachekey, post, time.Duration(h.cfg.Redis.ProductCachingTime)) // redise veriyi ekleme
	return c.JSON(post)
}

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var post dto.PostDto
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	reqpost, err := h.service.Create(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts"))
	return c.JSON(*reqpost)
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post dto.PostUpdateDto
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	post.UpdatedAt = time.Now()
	count, err := h.service.Update(id, &post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts/%s", id))
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts"))

	return c.JSON(fiber.Map{"success": true, "update count": count})
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts/%s", id))
	_ = h.cache.RemovePattern(fmt.Sprintf("*/posts"))
	return c.JSON(fiber.Map{"success": true})
}

// ::::::::::::::::
// 		JWT
// ::::::::::::::::

func (h *PostHandler) CreateJWT(c *fiber.Ctx) error {
	data := make(map[string]interface{})
	name := c.Params("name")
	data["name"] = name
	token := jtoken.GenerateAccessToken(data)
	return c.JSON(fiber.Map{"token": token})
}

func (h *PostHandler) ValidateJWT(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Hello, %s", c.Locals("name")))
}
