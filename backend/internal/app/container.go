package app

import (
	"database/sql"
	"time"

	"go-shop-app-backend/internal/infra/auth"
	"go-shop-app-backend/internal/infra/config"
	"go-shop-app-backend/internal/infra/db"
	"go-shop-app-backend/internal/orders"
	"go-shop-app-backend/internal/products"
	"go-shop-app-backend/internal/users"
	"go-shop-app-backend/pkg/workerpool"
)

type Container struct {
	Config *config.Config
	DB     *sql.DB
	JWT    *auth.Manager

	// WorkerPool используется для фоновой обработки задач (например, пост-обработки заказов).
	WorkerPool *workerpool.Pool

	UserRepo    users.Repository
	UserService users.Service

	ProductRepo    products.Repository
	ProductService products.Service

	OrderRepo    orders.Repository
	OrderService orders.Service
}

func NewContainer(cfg *config.Config) (*Container, error) {
	database, err := db.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	jwtManager := auth.NewManager(cfg.JWTSecret, 24*time.Hour)
	workerPool := workerpool.New(5)

	c := &Container{
		Config:     cfg,
		DB:         database,
		JWT:        jwtManager,
		WorkerPool: workerPool,
	}

	c.UserRepo = users.NewPostgresRepository(database)
	c.UserService = users.NewService(c.UserRepo, jwtManager)

	c.ProductRepo = products.NewPostgresRepository(database)
	c.ProductService = products.NewService(c.ProductRepo)

	c.OrderRepo = orders.NewPostgresRepository(database)
	c.OrderService = orders.NewService(c.OrderRepo, workerPool)

	return c, nil
}
