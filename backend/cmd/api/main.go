package main

import (
	"log"

	"go-shop-app-backend/internal/infra/config"
	"go-shop-app-backend/internal/infra/db"
	httph "go-shop-app-backend/internal/infra/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config error:", err)
	}

	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal("db error:", err)
	}

	r := httph.NewRouter(database, cfg)

	log.Println("server started on port", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
