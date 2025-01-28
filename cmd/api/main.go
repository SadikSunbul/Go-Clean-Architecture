// @title Go Clean Architecture API
// @version 1.0
// @description Bu API, Go Clean Architecture örnek projesidir.
// @termsOfService http://swagger.io/terms/

// @contact.name API Destek
// @contact.email your-email@domain.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @schemes http https
package main

// swag init -g cmd/api/main.go -o docs   | swager belgelrini oluşturu

import (
	"log"

	_ "github.com/SadikSunbul/Go-Clean-Architecture/docs" // swagger docs
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/server/http"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/redis"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
)

func main() {
	// ::::::::::::::::::::::::::::::::::::
	//       Configuration
	// ::::::::::::::::::::::::::::::::::::
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error on load configuration file, error: %v", err)
	}
	logger.Initialize(cfg.Environment)

	// ::::::::::::::::::::::::::::::::::::
	//       Database Connection
	// ::::::::::::::::::::::::::::::::::::
	db, err := db.NewMongoDB()
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	validator := validation.New()

	// ::::::::::::::::::::::::::::::::::::
	//       Redis service Connection
	// ::::::::::::::::::::::::::::::::::::

	cache := redis.New(redis.Config{
		Address:  cfg.Redis.Uri,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Db,
	})

	// ::::::::::::::::::
	//       Server
	// ::::::::::::::::::

	server := http.NewFiberServer(db, cfg, validator, cache)
	if err := server.Run(); err != nil {
		log.Fatalf("Server Error: %v", err)
	}

}
