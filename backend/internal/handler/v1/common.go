package handler

import (
	"net/http"

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
