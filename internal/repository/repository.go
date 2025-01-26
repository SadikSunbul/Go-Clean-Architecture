package repository

import (
	"context"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository[T interface{}] interface {
	Create(entity T) (T, error)
	Update(id string, entity T) (T, error)
	Delete(id string) error
	GetById(id string) (T, error)
	GetByField(filter bson.M) (T, error)
	GetAll() ([]T, error)
}

type Repostiroty[T interface{}] struct {
	collection *mongo.Collection
}

func NewRepository[T interface{}](db db.IDataBase, collectionName string) *Repostiroty[T] {
	return &Repostiroty[T]{
		collection: db.GetCollection(collectionName),
	}
}

func (r *Repostiroty[T]) Create(entity T) (T, error) {
	_, err := r.collection.InsertOne(context.Background(), entity)
	if err != nil {
		var zero T // T tipinin zero value'su
		return zero, err
	}
	return entity, nil
}

func (r *Repostiroty[T]) Update(id string, entity T) (*mongo.UpdateResult, error) {
	upresult, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": id}, entity)
	if err != nil {
		return nil, err
	}
	return upresult, nil
}

func (r *Repostiroty[T]) Delete(id string) (*mongo.DeleteResult, error) {
	deleteresult, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return deleteresult, nil
}

func (r *Repostiroty[T]) GetById(id string) (T, error) {
	var entity T
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		var zero T
		return zero, err
	}
	return entity, nil
}

func (r *Repostiroty[T]) GetByField(filter bson.M) (T, error) {
	var entity T
	err := r.collection.FindOne(context.Background(), filter).Decode(&entity)
	if err != nil {
		var zero T
		return zero, err
	}
	return entity, nil
}

func (r *Repostiroty[T]) GetAll() ([]T, error) {
	var entities []T
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
