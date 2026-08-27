package routes

import (
	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"

	//"github.com/Enziofael/nutrigo/shared/models/v1"

	handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
)

func RegisterUsersRoutes(group *gin.RouterGroup, s *server.Server) {
	userR := repository.NewUserPostgresRepository(s.DB)
	exerciseR := repository.NewExercisePostgresRepository(s.DB)

	userS := service.NewUserService(userR, exerciseR)
	exerciseS := service.NewExerciseService(exerciseR, userR)

	uh := handler.NewUserHandler(userS)
	eh := handler.NewExerciseHandler(exerciseS)

	users := group.Group("/users")

	// models/v1/user.go GLOBAL LEVEL
	// /users
	{
		// [models.UserCreateRequest]
		users.POST("", uh.Create)
		// [models.UserListRequest]
		users.GET("", uh.List)
	}
	// models/v1/user.go USER LEVEL
	// /users/:id
	{
		// [models.UserGetRequest]
		users.GET("/:id", uh.Get)
		// [models.UserPatchRequest]
		users.PATCH("/:id", uh.Patch)
		// [models.UserDeleteRequest]
		users.DELETE("/:id", uh.Delete)
	}
	// models/v1/exercise.go USER LEVEL
	// /users/:user_id/exercises
	{
		// [models.ExerciseCreateRequest]
		users.POST("/:user_id/exercises", eh.Create)
		// [models.ExerciseListRequest]
		users.GET("/:user_id/exercises", eh.List)
	}
}
