package repository

import (
	"context"
	"errors"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// :::::::::::::::::::::::::::::
// 			Repository
// :::::::::::::::::::::::::::::

//go:generate mockery --name=IRepository
type IRepository[T interface{}] interface {
	Create(entity T) (T, error)
	Update(id string, entity bson.M) (T, error)
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

func (r *Repostiroty[T]) Create(entity T) (T, interface{}, error) {
	resp, err := r.collection.InsertOne(context.Background(), entity)
	if err != nil {
		var zero T // T tipinin zero value'su
		return zero, nil, err
	}

	return entity, resp.InsertedID, nil
}

func (r *Repostiroty[T]) Update(id string, entity bson.M) (*mongo.UpdateResult, error) {
	fid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid post id format")
	}
	upresult, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": fid}, entity)
	if err != nil {
		return nil, err
	}
	return upresult, nil
}

func (r *Repostiroty[T]) Delete(id string) (*mongo.DeleteResult, error) {
	fid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid post id format")
	}
	deleteresult, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": fid})
	if err != nil {
		return nil, err
	}
	return deleteresult, nil
}

func (r *Repostiroty[T]) GetById(id string) (T, error) {
	// ID format kontrolü
	fid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var zero T
		return zero, errors.New("invalid post id format")
	}

	var entity T
	err = r.collection.FindOne(context.Background(), bson.M{"_id": fid}).Decode(&entity)
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
