package jtoken

import (
	"strings"
	"time"

	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/config"
	"github.com/SadikSunbul/Go-Clean-Architecture/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/quangdangfit/gocommon/logger"
)

// :::::::::::::::::
// 		JWT
// :::::::::::::::::

const (
	AccessTokenType  = "x-access"  // 5 minutes
	RefreshTokenType = "x-refresh" // 30 days
)

func GenerateAccessToken(payload map[string]interface{}) string {

	cfg := config.GetConfig()

	payload["type"] = AccessTokenType

	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Duration(cfg.JWT.AccessTokenExpiredTime) * time.Second).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(cfg.JWT.AuthSecret))
	if err != nil {
		logger.Error("Failed to generate access token: ", err)
		return ""
	}
	return token
}

func GenerateRefreshToken(payload map[string]interface{}) string {
	cfg := config.GetConfig()
	payload["type"] = RefreshTokenType

	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Duration(cfg.JWT.RefreshTokenExpiredTime) * time.Second).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(cfg.JWT.AuthSecret))
	if err != nil {
		logger.Error("Failed to generate refresh token: ", err)
		return ""
	}
	return token
}

func ValidateToken(jwtToken string) (map[string]interface{}, error) {

	cfg := config.GetConfig()
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.AuthSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	var data map[string]interface{}
	utils.Copy(&data, tokenData["payload"])

	return data, nil

}
