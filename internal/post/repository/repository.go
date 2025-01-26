package repository

import (
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/repository"
	"github.com/SadikSunbul/Go-Clean-Architecture/model/entity"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
)

type IPostRepository interface {
	repository.IRepository[entity.Post]
}

type PostRepository struct {
	repository.Repostiroty[entity.Post]
}

func NewPostRepository(db db.IDataBase) *PostRepository {
	return &PostRepository{
		Repostiroty: *repository.NewRepository[entity.Post](db, "posts"),
	}
}
