package middleware

import (
	"fmt"
	"time"

	"oauth/Go/common/config"

	"github.com/golang-jwt/jwt/v4"
)

var cfg, _ = config.LoadConfig()

func GenerateToken(UserID string) (string, error) {
	payload := jwt.MapClaims{
		"user_id": UserID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	if cfg.JWT.Secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).
		SignedString([]byte(cfg.JWT.Secret))
}

func ValidateToken(tokenStr string) (map[string]interface{}, error) {
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

