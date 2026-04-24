package util

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wan-system/internal/entity"
	"wan-system/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// TokenUtil handles JWT generation and validation
type TokenUtil struct {
	SecretKey []byte
}

func NewTokenUtil(secretKey string) *TokenUtil {
	return &TokenUtil{
		SecretKey: []byte(secretKey),
	}
}

// GenerateJWT generates a signed JWT for the given user (expires in 5 hours)
func (t *TokenUtil) GenerateJWT(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(5 * time.Hour).Unix(),
	})
	return token.SignedString(t.SecretKey)
}

// getTokenFromRequest extracts JWT string from Authorization header or cookie
func (t *TokenUtil) getTokenFromRequest(ctx *fiber.Ctx) (string, error) {
	authorization := ctx.Get("Authorization")
	var tokenString string

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if ctx.Cookies("token") != "" {
		tokenString = ctx.Cookies("token")
	}

	if tokenString == "" {
		return "", fiber.ErrUnauthorized
	}
	return tokenString, nil
}

func (t *TokenUtil) getToken(ctx *fiber.Ctx) (*jwt.Token, error) {
	tokenString, err := t.getTokenFromRequest(ctx)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.SecretKey, nil
	})
	return token, err
}

// ValidateJWT checks whether the JWT in the request is valid
func (t *TokenUtil) ValidateJWT(ctx *fiber.Ctx) error {
	token, err := t.getToken(ctx)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return fiber.ErrUnauthorized
}

// ParseToken parses a raw JWT string (with optional "Bearer " prefix) into an Auth struct
func (t *TokenUtil) ParseToken(_ context.Context, jwtToken string) (*model.Auth, error) {
	jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return t.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return nil, fiber.ErrUnauthorized
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	return &model.Auth{ID: id}, nil
}
