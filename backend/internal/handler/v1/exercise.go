package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
	"github.com/Enziofael/nutrigo/shared/models/v1"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	service *service.ExerciseService
}

func NewExerciseHandler(service *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{service: service}
}

func (h *ExerciseHandler) GetByID(c *gin.Context) {
	IDStr := c.Param("id")
	if IDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	exercise, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exercise == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exercise not found"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) ListByTgID(c *gin.Context) {
	tgIDStr := c.Param("tg_id")
	if tgIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tg_id is required"})
		return
	}

	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tg_id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort", "id")

	order, err := models.ValidateOrder(c.DefaultQuery("order", "desc"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order"})
	}

	req := models.ExerciseListRequest{
		TgID:   tgID,
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
	}

	exercises, err := h.service.ListByTgID(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exercises == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exercises not found"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (h *ExerciseHandler) CountByTgID(c *gin.Context) {
	tgIDStr := c.Param("tg_id")
	if tgIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tg_id is required"})
		return
	}

	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tg_id"})
		return
	}

	count, err := h.service.CountByTgID(c.Request.Context(), tgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *ExerciseHandler) Create(c *gin.Context) {
	var req models.ExerciseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercise, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (h *ExerciseHandler) Patch(c *gin.Context) {
	IDStr := c.Param("id")
	if IDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ID, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.ExercisePatchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patch request"})
		return
	}

	exercise, err := h.service.Patch(c.Request.Context(), ID, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrExerciseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "exercise not found"})
		case errors.Is(err, service.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update exercise"})
		}
		return
	}
	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) Delete(c *gin.Context) {
	IDStr := c.Param("id")
	if IDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ID, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), ID); err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "exercise not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
