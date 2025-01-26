package http

import (
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
	reqpost, err := h.service.Create(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(*reqpost)
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post dto.PostDto
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
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
