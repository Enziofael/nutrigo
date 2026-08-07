package handler

import (
	"errors"
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

// =========================================================
// CRUD HANDLERS
// =========================================================

func (h *ExerciseHandler) Create(c *gin.Context) {
	var req models.ExerciseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	exercise, err := h.service.Create(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (h *ExerciseHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.ExerciseDeleteRequest{ID: id}

	err = h.service.Delete(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "exercise not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ExerciseHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.ExerciseGetRequest{ID: id}

	exercise, err := h.service.Get(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "exercise not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *ExerciseHandler) List(c *gin.Context) {
	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort", "rating")
	order := models.Order(c.DefaultQuery("order", "desc"))
	req := models.ExerciseListRequest{
		TgID:   tgID,
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
	}

	exercises, err := h.service.List(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (h *ExerciseHandler) Search(c *gin.Context) {
	if c.Query("search") == "" {
		h.List(c)
	}

	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort", "similarity")
	order := models.Order(c.DefaultQuery("order", "desc"))
	req := models.ExerciseSearchRequest{
		TgID:   tgID,
		Search: search,
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
	}

	searchRes, err := h.service.Search(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, searchRes)
}

func (h *ExerciseHandler) Count(c *gin.Context) {
	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := models.ExerciseCountRequest{TgID: tgID}

	count, err := h.service.Count(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *ExerciseHandler) Patch(c *gin.Context) {
	var req models.ExercisePatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	exercise, err := h.service.Patch(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "exercise not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, exercise)
}
