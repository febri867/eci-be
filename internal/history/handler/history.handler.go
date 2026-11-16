package handler

import (
	"eci-backend/internal/history/repository"
	"net/http"
	"strconv"

	"eci-backend/internal/history/service"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	service *service.HistoryService
}

func NewHistoryHandler(svc *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		service: svc,
	}
}

func (h *HistoryHandler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid pagination parameters",
		})
		return
	}

	filter := repository.HistoryFilter{
		Action:    c.Query("action"),
		ActorID:   c.Query("actorId"),
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
	}

	events, total, summary, err := h.service.FindAll(
		c.Request.Context(),
		page,
		limit,
		filter,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    events,
		"page":    page,
		"limit":   limit,
		"total":   total,
		"hasMore": int64(page*limit) < total,
		"summary": summary,
	})
}
