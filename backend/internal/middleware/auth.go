package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/logger"
	"github.com/hpds/skill-hub/pkg/response"
)

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errno.Unauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			response.Error(c, errno.Unauthorized)
			c.Abort()
			return
		}

		claims, err := parseJWT(tokenStr, jwtSecret)
		if err != nil {
			logger.Warn("jwt parse failed", logger.String("error", err.Error()))
			response.Error(c, errno.Unauthorized)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errno.Unauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			response.Error(c, errno.Unauthorized)
			c.Abort()
			return
		}

		claims, err := parseJWT(tokenStr, jwtSecret)
		if err != nil {
			response.Error(c, errno.Unauthorized)
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			response.Error(c, errno.Forbidden)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.Next()
			return
		}

		claims, err := parseJWT(tokenStr, jwtSecret)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

type apiKeyRepo interface {
	GetByKey(string) (*model.APIKey, error)
	UpdateLastUsed(int64) error
}

func APIKeyAuth(apiKeyRepo apiKeyRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.Next()
			return
		}

		key, err := apiKeyRepo.GetByKey(apiKey)
		if err != nil || key == nil || key.IsRevoked {
			c.Next()
			return
		}

		c.Set("user_id", key.UserID)
		c.Set("auth_type", "api_key")
		_ = apiKeyRepo.UpdateLastUsed(key.ID)
		c.Next()
	}
}

func parseJWT(tokenStr, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
