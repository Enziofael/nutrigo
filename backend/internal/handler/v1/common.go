package handler

import (
	"errors"
	"net/http"
	"strconv"

	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	service *service.CommonService
}

func NewCommonHandler(service *service.CommonService) *CommonHandler {
	return &CommonHandler{service: service}
}

func (h *CommonHandler) Ping(c *gin.Context) {
	_, err := h.service.Ping(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}

func parseTgID(c *gin.Context) (int64, error) {
	idStr := c.Param("tg_id")
	if idStr == "" {
		return 0, errors.New("tg_id is required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid tg_id")
	}
	return id, nil
}

func parseID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, errors.New("id is required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
