package handler

import (
	"errors"
	"net/http"
	"time"

	"eci-backend/internal/auth/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	TurnstileToken string `json:"turnstileToken" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload",
		})
		return
	}

	user, token, expiresAt, err := h.authService.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.TurnstileToken,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Email atau password salah",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Terjadi kesalahan saat login",
		})
		return
	}

	maxAge := int(time.Until(expiresAt).Seconds())

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("eci_session", token, maxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("eci_session")
	if err == nil && token != "" {
		if err := h.authService.Logout(
			c.Request.Context(),
			token,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal logout",
			})
			return
		}
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("eci_session", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.authService.GetProfile(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
