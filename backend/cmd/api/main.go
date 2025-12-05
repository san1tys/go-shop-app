package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-shop-app-backend/internal/app"
	"go-shop-app-backend/internal/infra/config"
	"go-shop-app-backend/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatalf("container init error: %v", err)
	}

	application := app.NewApp(container)

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
