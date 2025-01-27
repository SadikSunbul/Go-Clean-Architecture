package http

import (
	"fmt"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/jtoken"
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/dto"
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/service"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service service.IPostService
}

func NewPostHandler(service service.IPostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

func (h *PostHandler) GetPostById(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := h.service.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
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
	return c.JSON(fiber.Map{"success": true, "update count": count})
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
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
