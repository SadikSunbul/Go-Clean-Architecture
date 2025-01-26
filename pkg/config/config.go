package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

const (
	ProductionEnv = "production"
)

// :::::::::::::::::
// 		Config
// :::::::::::::::::

type Config struct {
	Fiber       Fiber  `yaml:"fiber" json:"fiber"`
	Mongo       Mongo  `yaml:"mongo" json:"mongo"`
	JWT         JWT    `yaml:"jwt" json:"jwt"`
	Environment string `yaml:"environment" json:"environment"`
}

type JWT struct {
	AuthSecret              string `yaml:"auth_secret" json:"auth_secret"`
	AccessTokenExpiredTime  int    `yaml:"access_token_expired_time" json:"access_token_expired_time"`
	RefreshTokenExpiredTime int    `yaml:"refresh_token_expired_time" json:"refresh_token_expired_time"`
}

type Fiber struct {
	Port int `yaml:"port" json:"port"`
}

type Mongo struct {
	URI          string `yaml:"uri" json:"uri"`
	DatabaseName string `yaml:"databaseName" json:"databaseName"`
}

var cfg Config

func LoadConfig(fName string) (*Config, error) {

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	configPath := filepath.Join(currentDir, fName)

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config dosyası okunamadı: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config parse edilemedi: %w", err)
	}

	// Çevresel değişkenleri kontrol et
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("çevresel değişkenler parse edilemedi: %w", err)
	}

	return &cfg, nil
}

func GetConfig() *Config {
	return &cfg
}
