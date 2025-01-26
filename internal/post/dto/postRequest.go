package dto

import (
	"github.com/SadikSunbul/Go-Clean-Architecture/model/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostDto struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (p *PostDto) ToPost() entity.Post {
	return entity.Post{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *PostDto) FromPost(post entity.Post) PostDto {
	return PostDto{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
