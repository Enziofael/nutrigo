package routes

import (
	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"
	//handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
	//repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	//service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
)

func RegisterExerciseEntriesRoutes(group *gin.RouterGroup, s *server.Server) {

	//repo := repository.NewExercisePostgresRepository(s.DB)
	//service := service.NewExerciseService(repo)
	//handler := handler.NewExerciseHandler(service)

	exerciseEntries := group.Group("/exercise_entries")

	// models/v1/exerciseEntry.go EXERCISE ENTRY LEVEL
	// /exercise_entries/:id
	{
		// [models.ExerciseGetRequest]
		exerciseEntries.GET("/:id")
		// [models.ExercisePatchRequest]
		exerciseEntries.PATCH("/:id")
		// [models.ExerciseDeleteRequest]
		exerciseEntries.DELETE("/:id")
	}
}
