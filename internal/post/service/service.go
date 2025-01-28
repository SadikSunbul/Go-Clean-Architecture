package service

import (
	"errors"

	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/dto"
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/post/repository"
	"github.com/SadikSunbul/Go-Clean-Architecture/model/entity"
	"github.com/quangdangfit/gocommon/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockery --name=IPostService
type IPostService interface {
	Create(post *dto.PostDto) (*entity.Post, error)
	Update(id string, post *dto.PostUpdateDto) (int64, error)
	Delete(id string) error
	GetById(id string) (entity.Post, error)
	GetAll() (*[]entity.Post, error)
}

type PostService struct {
	validator  validation.Validation
	repository repository.PostRepository
}

func NewPostService(validator validation.Validation, repository repository.PostRepository) *PostService {
	return &PostService{
		validator:  validator,
		repository: repository,
	}
}

func (s *PostService) Create(post *dto.PostDto) (*entity.Post, error) {
	if err := s.validator.ValidateStruct(post); err != nil {
		return nil, err
	}

	result, id, err := s.repository.Create(post.ToPost())
	if err != nil {
		return nil, err
	}
	result.ID = id.(primitive.ObjectID)
	return &result, nil
}

func (s *PostService) Update(id string, post *dto.PostUpdateDto) (int64, error) {
	if err := s.validator.ValidateStruct(post); err != nil {
		return 0, err
	}

	result, err := s.repository.Update(id, bson.M{"$set": post})
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (s *PostService) Delete(id string) error {

	deleterespons, err := s.repository.Delete(id)
	if err != nil {
		return err
	}

	if deleterespons.DeletedCount == 0 {
		return errors.New("Post not found")
	}

	return nil
}

func (s *PostService) GetById(id string) (entity.Post, error) {

	getpost, err := s.repository.GetById(id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return entity.Post{}, errors.New("post not found")
		}
		return entity.Post{}, err
	}

	return getpost, nil
}

func (s *PostService) GetAll() (*[]entity.Post, error) {
	getposts, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return &getposts, nil
}
