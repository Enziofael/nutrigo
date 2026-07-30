package routes

import (
	apiv1 "github.com/Enziofael/nutrigo/backend/internal/routes/v1"
	"github.com/Enziofael/nutrigo/backend/internal/server"
)

func Register(server *server.Server) {
	router := server.GIN

	v1 := router.Group("/v1")
	{
		apiv1.RegisterCommonRoutes(v1, server)
		apiv1.RegisterUserRoutes(v1, server)
		apiv1.RegisterExerciseRoutes(v1, server)
		// apiv1.Register...Routes(v1, server)
		// ...
	}
	/*
		v2 := router.Group("/v2")
		{
			apiv2.Register...Routes(v2, server)
			...
		}
	*/
}
