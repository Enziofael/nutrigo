package v1

import (
	"net/http"
	"strconv"

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

func (h *UserHandler) GetByTgID(c *gin.Context) {
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

	user, err := h.service.GetByTgID(c.Request.Context(), int64(tgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
