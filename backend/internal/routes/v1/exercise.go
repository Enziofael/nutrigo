package routes

import (
	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"

	handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
)

func RegisterExerciseRoutes(group *gin.RouterGroup, s *server.Server) {

	repo := repository.NewExercisePostgresRepository(s.DB)
	service := service.NewExerciseService(repo)
	handler := handler.NewExerciseHandler(service)

	exercises := group.Group("/exercises")
	{
		exercises.GET("/:id")
		exercises.PATCH("/:id")
		exercises.DELETE("/:id")
	}
	{
		exercise.POST("/", handler.Create)
	}
	{
		exercise.PATCH("/:id", handler.Patch)
	}
	{
		exercise.DELETE("/:id", handler.Delete)
	}
}
