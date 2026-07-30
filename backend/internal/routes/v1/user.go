package routes

import (
	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"

	handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
)

func RegisterUserRoutes(group *gin.RouterGroup, s *server.Server) {

	repo := repository.NewUserPostgresRepository(s.DB)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	user := group.Group("/user")
	{
		user.GET("/:tg_id", handler.GetByTgID)
	}
	{
		user.POST("/", handler.Create)
	}
	{
		user.PATCH("/:tg_id", handler.Patch)
	}
	{
		user.DELETE("/:tg_id", handler.Delete)
	}
}
