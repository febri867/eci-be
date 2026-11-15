package handler

import (
	"net/http"
	"strconv"
	"time"

	"eci-backend/internal/employee/model"
	"eci-backend/internal/employee/service"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

type EmployeeRequest struct {
	Name      string `json:"name" binding:"required"`
	Position  string `json:"position" binding:"required"`
	Grade     string `json:"grade" binding:"required"`
	Salary    int64  `json:"salary" binding:"required,gt=0"`
	Address   string `json:"address" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	BirthDate string `json:"birthDate" binding:"required"`
}

type EmployeeListResponse struct {
	Data    []model.Employee `json:"data"`
	Page    int64            `json:"page"`
	Limit   int64            `json:"limit"`
	Total   int64            `json:"total"`
	HasMore bool             `json:"hasMore"`
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var request EmployeeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid employee data"})
		return
	}

	employee, err := request.toModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.Create(
		c.Request.Context(),
		employee,
		c.GetString("userID"),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) FindAll(c *gin.Context) {
	page := parseQueryInt(c.Query("page"), 1)
	limit := parseQueryInt(c.Query("limit"), 10)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	minAge := parseQueryInt(c.Query("minAge"), 0)
	maxAge := parseQueryInt(c.Query("maxAge"), 0)

	if minAge < 0 || maxAge < 0 || (maxAge > 0 && minAge > maxAge) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid age range",
		})
		return
	}

	filter := service.EmployeeFilter{
		Search:   c.Query("search"),
		Grade:    c.Query("grade"),
		MinAge:   int(minAge),
		MaxAge:   int(maxAge),
		Province: c.Query("province"),
		City:     c.Query("city"),
		District: c.Query("district"),
	}

	employees, total, err := h.service.FindAll(
		c.Request.Context(),
		filter,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve employees",
		})
		return
	}

	c.JSON(http.StatusOK, EmployeeListResponse{
		Data:    employees,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasMore: page*limit < total,
	})
}

func (h *EmployeeHandler) FindByID(c *gin.Context) {
	employee, err := h.service.FindByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	var request EmployeeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid employee data"})
		return
	}

	employee, err := request.toModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		employee,
		c.GetString("userID"),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee updated successfully",
	})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
		c.GetString("userID"),
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee deleted successfully",
	})
}

func (r EmployeeRequest) toModel() (*model.Employee, error) {
	birthDate, err := time.Parse("2006-01-02", r.BirthDate)
	if err != nil {
		return nil, errInvalidBirthDate
	}

	return &model.Employee{
		Name:      r.Name,
		Position:  r.Position,
		Grade:     r.Grade,
		Salary:    r.Salary,
		Address:   r.Address,
		Province:  r.Province,
		City:      r.City,
		District:  r.District,
		BirthDate: birthDate,
	}, nil
}

func parseQueryInt(value string, fallback int64) int64 {
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

var errInvalidBirthDate = &birthDateError{}

type birthDateError struct{}

func (e *birthDateError) Error() string {
	return "Birth date must use YYYY-MM-DD format"
}
