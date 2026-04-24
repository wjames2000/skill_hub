package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hpds/skill-hub/internal/middleware"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/logger"
	"github.com/hpds/skill-hub/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo   *repository.UserRepo
	apiKeyRepo *repository.APIKeyRepo
	jwtSecret  string
	jwtExpire  int
}

func NewAuthHandler(userRepo *repository.UserRepo, apiKeyRepo *repository.APIKeyRepo, jwtSecret string, jwtExpire int) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, apiKeyRepo: apiKeyRepo, jwtSecret: jwtSecret, jwtExpire: jwtExpire}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/register", h.Register)
	rg.GET("/auth/github", h.GitHubOAuth)
	rg.POST("/auth/github/callback", h.GitHubCallback)
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if user == nil {
		user, err = h.userRepo.GetByEmail(req.Username)
		if err != nil {
			response.Error(c, errno.DBError)
			return
		}
	}
	if user == nil {
		response.Error(c, errno.UserNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, errno.InvalidPassword)
		return
	}

	_ = h.userRepo.UpdateLastLogin(user.ID)

	token, err := h.generateJWT(user)
	if err != nil {
		logger.Error("generate jwt", logger.ErrorField(err))
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	existing, _ := h.userRepo.GetByUsername(req.Username)
	if existing != nil {
		response.Error(c, errno.UserExists)
		return
	}

	existing, _ = h.userRepo.GetByEmail(req.Email)
	if existing != nil {
		response.Error(c, errno.UserExists)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("hash password", logger.ErrorField(err))
		response.Error(c, errno.InternalError)
		return
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
		Status:       1,
	}

	if err := h.userRepo.Create(user); err != nil {
		response.Error(c, errno.DBError)
		return
	}

	token, err := h.generateJWT(user)
	if err != nil {
		logger.Error("generate jwt", logger.ErrorField(err))
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
		},
	})
}

func (h *AuthHandler) GitHubOAuth(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		response.Error(c, errno.ParamError)
		return
	}

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=user:email,read:user",
		clientID,
	)
	c.JSON(http.StatusOK, gin.H{"redirect_url": redirectURL})
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errno.ParamError)
		return
	}

	user, err := h.userRepo.GetByGitHubID(req.Code)
	if err != nil {
		response.Error(c, errno.DBError)
		return
	}
	if user == nil {
		user = &model.User{
			Username: fmt.Sprintf("gh_%s", sha256Hex(req.Code)[:12]),
			Email:    fmt.Sprintf("%s@github.user", sha256Hex(req.Code)[:12]),
			GitHubID: req.Code,
			Role:     "user",
			Status:   1,
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(sha256Hex(req.Code)), bcrypt.DefaultCost)
		user.PasswordHash = string(hash)
		if err := h.userRepo.Create(user); err != nil {
			response.Error(c, errno.DBError)
			return
		}
	}

	_ = h.userRepo.UpdateLastLogin(user.ID)

	token, err := h.generateJWT(user)
	if err != nil {
		response.Error(c, errno.InternalError)
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
		},
	})
}

func (h *AuthHandler) generateJWT(user *model.User) (string, error) {
	expire := time.Duration(h.jwtExpire) * time.Hour
	claims := &middleware.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "skill-hub",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
