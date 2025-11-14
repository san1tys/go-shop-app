package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"go-shop-app-backend/internal/infra/config"
)

func NewRouter(db *sql.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	h := NewHandler(db, cfg)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	return r
}
