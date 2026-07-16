package main

import (
	"database/sql"
	"log"

	"github.com/Enziofael/nutrigo/backend/internal/middleware"
	"github.com/Enziofael/nutrigo/backend/internal/routes"
	"github.com/Enziofael/nutrigo/backend/internal/server"
	cfg "github.com/Enziofael/nutrigo/shared/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	//initing

	db, err := sql.Open("postgres", cfg.DBConString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	
	s := server.New(r, db)

	middleware.Register(s)
	routes.Register(s)

	err = r.Run(cfg.GetBackendPort())
	if err != nil {
		log.Fatalf("Fatal error at main.go: %v", err)
	}
}
