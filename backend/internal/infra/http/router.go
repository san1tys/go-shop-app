package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"go-shop-app-backend/internal/infra/config"
	"go-shop-app-backend/internal/products"
)

func NewRouter(db *sql.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	h := NewHandler(db, cfg)

	// --- AUTH ROUTES ---
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	// --- PRODUCTS MODULE ---
	productRepo := products.NewPostgresRepository(db)
	productService := products.NewService(productRepo)
	productHandler := products.NewHandler(productService)
	productHandler.RegisterRoutes(r)
	// ----------------------

	return r
}
