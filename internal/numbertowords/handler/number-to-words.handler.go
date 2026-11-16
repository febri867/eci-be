package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eci-backend/internal/numbertowords/service"
)

type NumberToWordsHandler struct {
	service *service.NumberToWordsService
}

func NewNumberToWordsHandler(
	numberToWordsService *service.NumberToWordsService,
) *NumberToWordsHandler {
	return &NumberToWordsHandler{
		service: numberToWordsService,
	}
}

type convertRequest struct {
	Number int64 `json:"number" binding:"required,min=0"`
}

func (h *NumberToWordsHandler) Convert(c *gin.Context) {
	var request convertRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid number",
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

	result, err := h.service.Convert(
		c.Request.Context(),
		actorID.(string),
		request.Number,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number": request.Number,
		"text":   result,
	})
}
