package v1

import (
	"database/sql"

	"github.com/Enziofael/nutrigo/backend/internal/server"
	"github.com/gin-gonic/gin"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
	handler "github.com/Enziofael/nutrigo/backend/internal/handler/v1"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

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
}
