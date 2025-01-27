package http

import (
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/repository"
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/service"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/quangdangfit/gocommon/validation"
)

func Routes(app fiber.Router, db db.IDataBase, validator validation.Validation) {
	postRepo := repository.NewPostRepository(db) // Database Collection setting is made here
	postService := service.NewPostService(validator, *postRepo)
	postHandler := NewPostHandler(postService)

	jwtvalidmiddleware := middleware.JWTAuth()

	posts := app.Group("/posts")

	// JWT Token creat and validate
	posts.Post("/createjwt/:name", postHandler.CreateJWT)
	posts.Get("/validatejwt", jwtvalidmiddleware, postHandler.ValidateJWT) // Bunun konumu önemlidir çünkü sistem GetPostById ile karıştırıp hata verir.

	posts.Get("/", postHandler.GetAllPosts)
	posts.Get("/:id", postHandler.GetPostById)
	posts.Post("/", postHandler.CreatePost)
	posts.Put("/:id", postHandler.UpdatePost)
	posts.Delete("/:id", postHandler.DeletePost)

}
