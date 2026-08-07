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

	exercise := group.Group("/exercise")
	{
		exercise.GET("/:id", handler.Get)
		exercise.GET("/u/:tg_id", handler.List)
		exercise.GET("/u/:tg_id/count", handler.Count)
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
