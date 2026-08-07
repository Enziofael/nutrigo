package handler

import (
	"errors"
	"net/http"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// =========================================================
// CRUD HANDLERS
// =========================================================

func (h *UserHandler) Create(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	user, err := h.service.Create(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := models.UserDeleteRequest{TgID: tgID}
	err = h.service.Delete(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Get(c *gin.Context) {
	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := models.UserGetRequest{TgID: tgID}
	user, err := h.service.Get(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Patch(c *gin.Context) {
	var req models.UserPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TgID = tgID

	user, err := h.service.Patch(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
