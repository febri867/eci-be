package handler

import (
	"net/http"

	"eci-backend/internal/region/model"
	"eci-backend/internal/region/service"

	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	service *service.RegionService
}

func NewRegionHandler(service *service.RegionService) *RegionHandler {
	return &RegionHandler{
		service: service,
	}
}

func (h *RegionHandler) GetProvinces(c *gin.Context) {
	regions, err := h.service.GetProvinces(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve provinces",
		})
		return
	}

	c.JSON(http.StatusOK, model.RegionResponse{
		Data: toRegionItems(regions),
		Meta: model.RegionMeta{Level: 1},
	})
}

func (h *RegionHandler) GetCities(c *gin.Context) {
	regions, err := h.service.GetCities(
		c.Request.Context(),
		c.Param("provinceId"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.RegionResponse{
		Data: toRegionItems(regions),
		Meta: model.RegionMeta{Level: 2},
	})
}

func (h *RegionHandler) GetDistricts(c *gin.Context) {
	regions, err := h.service.GetDistricts(
		c.Request.Context(),
		c.Param("cityId"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.RegionResponse{
		Data: toRegionItems(regions),
		Meta: model.RegionMeta{Level: 3},
	})
}

func toRegionItems(regions []model.Region) []model.RegionItem {
	items := make([]model.RegionItem, 0, len(regions))

	for _, region := range regions {
		items = append(items, model.RegionItem{
			ID:   region.ID,
			Name: region.Name,
		})
	}

	return items
}
