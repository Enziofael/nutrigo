package routes

import (
	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"
	//handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
	//repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	//service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
)

func RegisterExercisesRoutes(group *gin.RouterGroup, s *server.Server) {

	//repo := repository.NewExercisePostgresRepository(s.DB)
	//service := service.NewExerciseService(repo)
	//handler := handler.NewExerciseHandler(service)

	exercises := group.Group("/exercises")

	// models/v1/exercise.go EXERCISE LEVEL
	// /exercises/:id
	{
		// [models.ExerciseGetRequest]
		exercises.GET("/:id")
		// [models.ExercisePatchRequest]
		exercises.PATCH("/:id")
		// [models.ExerciseDeleteRequest]
		exercises.DELETE("/:id")
	}
	// models/v1/exerciseEntry.go EXERCISE LEVEL
	// /exercises/:exercise_id/entry
	{
		// [models.ExerciseEntryCreateRequest]
		exercises.POST("/:exercise_id/entries") 
		// [models.ExerciseEntryListRequest]
		exercises.GET("/:exercise_id/entries")
	}
}
