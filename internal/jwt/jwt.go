package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(config *config.Config) *JWTService {
	return &JWTService{
		secret: []byte(config.JWT_SECRET),
	}
}

func (j *JWTService) GenerateJWT(userID uuid.UUID, typ TokenType, expiry time.Duration) (string, error) {

	if len(j.secret) == 0 {
		return "", fmt.Errorf("invalid jwt secret key")
	}

	claims := jwt.MapClaims{
		"sub":  userID, // subject
		"exp":  time.Now().Add(expiry).Unix(),
		"type": typ,
		"iat":  time.Now().Unix(), // issued at
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) VerifyToken(tokenString string) (*jwt.Token, error) {

	if len(j.secret) == 0 {
		return nil, fmt.Errorf("invalid jwt secret key")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
}
