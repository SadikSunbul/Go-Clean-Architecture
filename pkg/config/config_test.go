package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) SetupTest() {
	// Her test öncesi config'i sıfırla
	cfg = Config{} // Global cfg'yi sıfırla
}

func (s *ConfigTestSuite) TestLoadConfig_Success() {
	// Config dosyasını yükle
	config, err := LoadConfig("config.yaml")

	// Test Kontrolleri:
	assert.NoError(s.T(), err)                  // Hata olmamalı
	assert.NotNil(s.T(), config)                // Config nil olmamalı
	assert.NotEmpty(s.T(), config.Mongo.URI)    // MongoDB URI'si boş olmamalı
	assert.Greater(s.T(), config.Fiber.Port, 0) // Port 0'dan büyük olmalı
}

func (s *ConfigTestSuite) TestLoadConfig_FileNotFound() {
	// Var olmayan dosyayı yüklemeye çalış
	config, err := LoadConfig("nonexistent.yaml")

	// Test Kontrolleri:
	assert.Error(s.T(), err)                                        // Hata dönmeli
	assert.Nil(s.T(), config)                                       // Config nil olmalı
	assert.Contains(s.T(), err.Error(), "config dosyası okunamadı") // Hata mesajı kontrol
}

func (s *ConfigTestSuite) TestGetConfig_WithoutLoad() {
	// Config yüklemeden GetConfig çağır
	config := GetConfig()

	// Test Kontrolleri:
	assert.NotNil(s.T(), config)          // Config nil olmamalı
	assert.Empty(s.T(), config.Mongo.URI) // URI boş olmalı
	assert.Zero(s.T(), config.Fiber.Port) // Port 0 olmalı
}

func (s *ConfigTestSuite) TestGetConfig_AfterLoad() {
	// Önce config'i yükle
	loadedConfig, err := LoadConfig("config.yaml")
	assert.NoError(s.T(), err) // Yükleme hatası olmamalı

	// GetConfig ile al
	config := GetConfig()

	// Test Kontrolleri:
	assert.Equal(s.T(), loadedConfig.Mongo.URI, config.Mongo.URI)   // URI'ler eşit olmalı
	assert.Equal(s.T(), loadedConfig.Fiber.Port, config.Fiber.Port) // Portlar eşit olmalı
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
