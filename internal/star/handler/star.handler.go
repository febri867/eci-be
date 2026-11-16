package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eci-backend/internal/star/service"
)

type StarHandler struct {
	service *service.StarService
}

func NewStarHandler(starService *service.StarService) *StarHandler {
	return &StarHandler{
		service: starService,
	}
}

type generateRequest struct {
	Rows    int    `json:"rows" binding:"required,min=1,max=50"`
	Pattern string `json:"pattern" binding:"required"`
}

func (h *StarHandler) Generate(c *gin.Context) {
	var request generateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	actorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	result, err := h.service.Generate(
		c.Request.Context(),
		actorID.(string),
		request.Rows,
		request.Pattern,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":    request.Rows,
		"pattern": request.Pattern,
		"result":  result,
	})
}
