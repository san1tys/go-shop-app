package app

import (
	"context"
)

type App struct {
	Container *Container
	Server    *Server
}

func NewApp(c *Container) *App {
	return &App{
		Container: c,
		Server:    NewHTTPServer(c),
	}
}

func (a *App) Run() error {
	return a.Server.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	// Корректно останавливаем worker pool, чтобы дождаться фоновых задач.
	if a.Container != nil && a.Container.WorkerPool != nil {
		a.Container.WorkerPool.Stop()
	}

	return nil
}
