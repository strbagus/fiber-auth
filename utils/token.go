package utils

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	d "github.com/strbagus/fiber-auth/database"
)

func IsJWTUsable(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	ctx := context.Background()
	exists, err := d.RD.Exists(ctx, "blacklist:"+tokenString).Result()
	if err != nil {
		return false
	}
	return exists == 0
}

func GenerateToken(payload map[string]interface{}, expires time.Duration) (string, error) {
	hmacSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"exp": time.Now().Add(expires).Unix(),
	}
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecret)
}
