package main

import (
	"database/sql"
	"log"

	"github.com/Enziofael/nutrigo/backend/internal/server"
	cfg "github.com/Enziofael/nutrigo/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	//initing

	db, err := sql.Open("postgres", cfg.DbConString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	server.SetupDb(db)

	r := gin.Default()
	server.Setup(r)

	err = r.Run(cfg.GetBackendPort())
}
