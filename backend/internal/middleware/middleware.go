package middleware

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Enziofael/nutrigo/backend/internal/server"
	cfg "github.com/Enziofael/nutrigo/shared/config"
	"github.com/gin-gonic/gin"
)

func Register(s *server.Server) {
	s.GIN.Use(databaseMiddleware(s.DB))
	s.GIN.Use(apiTokenMiddleware())
}

func apiTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "unauthorized",
				"message": "X-API-Key header is empty",
			})
			c.Abort()
			return
		}
		if cfg.GetBackendAPIToken() != c.GetHeader("X-API-Key") {

			c.JSON(http.StatusForbidden, gin.H{
				"status":  "forbidden",
				"message": "invalid api token",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func databaseMiddleware(db *sql.DB) gin.HandlerFunc {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
