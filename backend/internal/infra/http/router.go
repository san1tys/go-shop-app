package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-shop-app-backend/internal/infra/auth"
	"go-shop-app-backend/internal/infra/config"
	infraDB "go-shop-app-backend/internal/infra/db"
	"go-shop-app-backend/internal/orders"
	"go-shop-app-backend/internal/products"
	"go-shop-app-backend/internal/users"
)

func NewRouter(db *sql.DB, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		if err := infraDB.HealthCheck(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	jwtManager := auth.NewManager(cfg.JWTSecret, 24*time.Hour)

	authRequired := v1.Group("/")
	authRequired.Use(AuthMiddleware(jwtManager))

	adminGroup := v1.Group("/admin")
	adminGroup.Use(AuthMiddleware(jwtManager), AdminOnly())

	userRepo := users.NewPostgresRepository(db)
	userService := users.NewService(userRepo, jwtManager)
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(v1)

	productRepo := products.NewPostgresRepository(db)
	productService := products.NewService(productRepo)
	productHandler := products.NewHandler(productService)
	productHandler.RegisterRoutes(v1)

	productHandler.RegisterRoutes(adminGroup)

	orderRepo := orders.NewPostgresRepository(db)
	orderService := orders.NewService(orderRepo)
	orderHandler := orders.NewHandler(orderService)
	orderHandler.RegisterRoutes(authRequired)

	return r
}
