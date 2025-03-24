package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/personal/task-management/pkg/server"
)

type App struct {
	servers []server.Server
	name    string
}

type Option func(*App)

func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithServer(server server.Server) Option {
	return func(a *App) {
		a.servers = append(a.servers, server)
	}
}

func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

func (a *App) Run() error {
	log.Printf("Starting %s", a.name)

	for _, s := range a.servers {
		if err := s.Start(context.Background()); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	for _, s := range a.servers {
		if err := s.Stop(context.Background()); err != nil {
			log.Printf("Failed to stop server: %v", err)
		}
	}

	log.Printf("Shutting down %s", a.name)
	return nil
}

func (a *App) Stop() error {
	for _, s := range a.servers {
		if err := s.Stop(context.Background()); err != nil {
			log.Printf("Failed to stop server: %v", err)
		}
	}
	return nil
}
