package v1

import (
	handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	"github.com/Enziofael/nutrigo/backend/internal/server"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
	"github.com/gin-gonic/gin"
)

func RegisterCommonRoutes(group *gin.RouterGroup, s *server.Server) {
	repo := repository.NewCommonPostgresRepository(s.DB)
	service := service.NewCommonService(repo)
	handler := handler.NewCommonHandler(service)

	{
		group.GET("/ping", handler.Ping)
	}
}
