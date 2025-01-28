package http

import (
	"testing"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	dbMocks "github.com/SadikSunbul/Go-Clean-Architecture/pkg/db/mocks"
	redisMocks "github.com/SadikSunbul/Go-Clean-Architecture/pkg/redis/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_NewFiberServer(t *testing.T) {
	mockDB := new(dbMocks.IDataBase)
	mockRedis := new(redisMocks.IRedis)
	mockConfig := new(config.Config)

	server := NewFiberServer(mockDB, mockConfig, validation.New(), mockRedis)
	assert.NotNil(t, server)
}

func Test_Server_GetApp(t *testing.T) {
	mockDB := new(dbMocks.IDataBase)
	mockedCollection := &mongo.Collection{} // Mock koleksiyon oluşturun
	mockDB.On("GetCollection", "posts").Return(mockedCollection)
	mockRedis := new(redisMocks.IRedis)
	mockConfig := new(config.Config)

	server := NewFiberServer(mockDB, mockConfig, validation.New(), mockRedis)
	assert.NotNil(t, server)
	server.app = fiber.New()

	engine := server.GetApp()
	assert.NotNil(t, engine) // *engine yerine engine kullan
}

func Test_Server_MapRoutes(t *testing.T) {
	mockDB := new(dbMocks.IDataBase)
	mockedCollection := &mongo.Collection{} // Mock koleksiyon oluşturun
	mockDB.On("GetCollection", "posts").Return(mockedCollection)
	mockRedis := new(redisMocks.IRedis)
	mockConfig := new(config.Config)

	server := NewFiberServer(mockDB, mockConfig, validation.New(), mockRedis)
	assert.NotNil(t, server)
	server.app = fiber.New()

	err := server.MapRoutes()
	assert.Nil(t, err)
}
