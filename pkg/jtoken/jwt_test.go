package jtoken

import (
	"testing"
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JWTTestSuite struct {
	suite.Suite
}

func (s *JWTTestSuite) SetupTest() {
	// Her test öncesi config yükle
	_, err := config.LoadConfig("../config/config.yaml")
	assert.NoError(s.T(), err)
}

func (s *JWTTestSuite) TestGenerateAccessToken_Success() {
	// Test verisi
	payload := map[string]interface{}{
		"user_id": "123",
		"role":    "admin",
	}

	// Token oluştur
	token := GenerateAccessToken(payload)

	// Kontroller
	assert.NotEmpty(s.T(), token)

	// Token'ı doğrula
	decodedPayload, err := ValidateToken(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), payload["user_id"], decodedPayload["user_id"])
	assert.Equal(s.T(), payload["role"], decodedPayload["role"])
	assert.Equal(s.T(), AccessTokenType, decodedPayload["type"])
}

func (s *JWTTestSuite) TestGenerateRefreshToken_Success() {
	// Test verisi
	payload := map[string]interface{}{
		"user_id": "123",
	}

	// Token oluştur
	token := GenerateRefreshToken(payload)

	// Kontroller
	assert.NotEmpty(s.T(), token)

	// Token'ı doğrula
	decodedPayload, err := ValidateToken(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), payload["user_id"], decodedPayload["user_id"])
	assert.Equal(s.T(), RefreshTokenType, decodedPayload["type"])
}

func (s *JWTTestSuite) TestInvalidToken() {
	// Geçersiz token doğrulama
	payload, err := ValidateToken("invalid.token.here")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), payload)
}

func (s *JWTTestSuite) TestEmptyToken() {
	// Boş token doğrulama
	payload, err := ValidateToken("")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), payload)
}

func (s *JWTTestSuite) TestTokenWithoutBearer() {
	// Bearer prefix'i olmayan token oluştur
	token := GenerateAccessToken(map[string]interface{}{"user_id": "123"})

	// Token'ı doğrula
	payload, err := ValidateToken(token)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), payload)
}

func (s *JWTTestSuite) TestTokenWithBearer() {
	// Bearer prefix'li token oluştur
	token := "Bearer " + GenerateAccessToken(map[string]interface{}{"user_id": "123"})

	// Token'ı doğrula
	payload, err := ValidateToken(token)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), payload)
}

func (s *JWTTestSuite) TestTokenPayloadContent() {
	// Karmaşık payload ile test
	payload := map[string]interface{}{
		"user_id":    "123",
		"role":       "admin",
		"email":      "test@test.com",
		"is_active":  true,
		"login_time": time.Now().Unix(),
	}

	// Token oluştur ve doğrula
	token := GenerateAccessToken(payload)
	decodedPayload, err := ValidateToken(token)

	// Kontroller
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), decodedPayload)
	assert.Equal(s.T(), payload["user_id"], decodedPayload["user_id"])
	assert.Equal(s.T(), payload["role"], decodedPayload["role"])
	assert.Equal(s.T(), payload["email"], decodedPayload["email"])
	assert.Equal(s.T(), payload["is_active"], decodedPayload["is_active"])
}

func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

func Test_ExpiredToken(t *testing.T) {
	token := GenerateAccessToken(map[string]interface{}{"testuser": "sadik"})
	// Daha uzun bir süre bekleyin
	time.Sleep(time.Duration(config.GetConfig().JWT.AccessTokenExpiredTime+1) * time.Second)
	payload, err := ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, payload)
}
