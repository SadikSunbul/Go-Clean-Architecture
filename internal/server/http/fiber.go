package http

import (
	"fmt"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/quangdangfit/gocommon/validation"
)

// ::::::::::::::::::::::::::
// 		Fiber Server
// ::::::::::::::::::::::::::

type FiberServer struct {
	app       *fiber.App
	cfg       *config.Config
	validator validation.Validation
	db        db.MongoDB
}

func NewFiberServer(db db.MongoDB, cfg *config.Config, validator validation.Validation) *FiberServer {
	return &FiberServer{
		db:        db,
		cfg:       cfg,
		validator: validator,
	}
}

func (s *FiberServer) Run() error {
	s.app = fiber.New()

	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.Fiber.Port))
}

func (s *FiberServer) GetApp() *fiber.App {
	return s.app
}

func (s *FiberServer) MapRoutes() error {
	v1 := s.app.Group("/api/v1")

	v1.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	return nil
}
