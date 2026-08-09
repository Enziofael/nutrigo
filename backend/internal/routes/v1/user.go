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

	users := group.Group("/users")
	{
		users.POST("/", handler.Create)
		users.GET("/", handler.Get)
	}
	{
		users.GET("/:id", handler.Get)
		users.PATCH("/:id", handler.Patch)
		users.DELETE("/:id", handler.Delete)
	}
	{
		users.POST("/:id/exercises")
		users.GET("/:id/exercises")
	}
}
