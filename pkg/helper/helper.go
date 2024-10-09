package helper

import (
	"errors"
	"fmt"
	"taskmanagementapi/pkg/config"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserClaims struct {
	Id    string
	Email string
	jwt.StandardClaims
}

func GetTokenFromHeader(header string) string {
	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}

	return header
}

func ExtractUserIDFromToken(tokenString string) (string, string, error) {
	cfg, _ := config.LoadConfig()
	token, err := jwt.ParseWithClaims(tokenString, &AuthUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(cfg.JwtSecretKey), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*AuthUserClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	return claims.Id, claims.Email, nil

}

func PasswordHash(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("error from password hashing")
	}
	hash := string(hashPassword)
	return hash, nil
}