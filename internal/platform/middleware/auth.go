package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"eci-backend/internal/auth/repository"

	"github.com/gin-gonic/gin"
)

func Auth(sessionRepository *repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("eci_session")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])

		session, err := sessionRepository.FindByTokenHash(
			c.Request.Context(),
			tokenHash,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired session",
			})
			return
		}

		c.Set("userID", session.UserID.Hex())
		c.Next()
	}
}
