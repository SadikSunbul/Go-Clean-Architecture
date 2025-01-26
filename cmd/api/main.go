package main

import (
	"github.com/SadikSunbul/Go-Clean-Architecture/internal/server/http"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/db"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"log"
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

	// ::::::::::::::::::
	//       Server
	// ::::::::::::::::::

	server := http.NewFiberServer(*db, cfg, validator)
	if err := server.Run(); err != nil {
		log.Fatalf("Server Error: %v", err)
	}

}
