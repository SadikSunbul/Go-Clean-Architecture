package db

import (
	"context"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ::::::::::::::::::::
// 		MongoDB
// ::::::::::::::::::::

// mockery --name=IDataBase --output=./mocks     | ile üret

//go:generate mockery --name=IDatabase
type IDataBase interface {
	GetDatabase() *mongo.Database
	CreateCollection(name string) *mongo.Collection
	GetCollection(name string) *mongo.Collection
}

type MongoDB struct {
	Database *mongo.Database
}

func NewMongoDB() (*MongoDB, error) {
	cfg := config.GetConfig()

	clientOptions := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	db := client.Database(cfg.Mongo.DatabaseName)
	return &MongoDB{Database: db}, nil
}

func (db *MongoDB) GetDatabase() *mongo.Database {
	return db.Database
}
func (db *MongoDB) CreateCollection(name string) *mongo.Collection {
	return db.Database.Collection(name)
}
func (db *MongoDB) GetCollection(name string) *mongo.Collection {
	return db.Database.Collection(name)
}
