package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Server struct {
	GIN *gin.Engine
	DB  *sql.DB
}

func New(gin *gin.Engine, db *sql.DB) *Server {
	return &Server{
		GIN: gin,
		DB:  db,
	}
}
