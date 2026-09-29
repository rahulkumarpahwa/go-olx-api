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

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func getSecretKey(config *config.Config) []byte {
	return []byte(config.JWT_SECRET)
}

func GenerateJWT(config *config.Config, userID uuid.UUID, typ TokenType, expiry time.Duration) (string, error) {

	secret := getSecretKey(config)

	if len(secret) == 0 {
		return "", fmt.Errorf("invalid jwt secret key")
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
		UserID:    userID,
		TokenType: typ,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func VerifyToken(
	config *config.Config,
	tokenString string,
) (*Claims, error) {

	secret := getSecretKey(config)

	if len(secret) == 0 {
		return nil, fmt.Errorf("invalid jwt secret key")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
