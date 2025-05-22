package app

import (
	"fmt"

	"github.com/vagonaizer/loms/internal/config"
	"github.com/vagonaizer/loms/internal/infrastructure/api/grpc"
	lomsclient "github.com/vagonaizer/loms/internal/infrastructure/client/loms"
	"github.com/vagonaizer/loms/internal/infrastructure/repository/inmemory"
	"github.com/vagonaizer/loms/internal/usecase/loms"
)

type App struct {
	config *config.Config
	server *grpc.Server
	loms   *lomsclient.Client
}

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize LOMS client
	lomsClient, err := lomsclient.NewClient(cfg.LOMS.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to create LOMS client: %w", err)
	}

	// Initialize repositories
	orderRepo := inmemory.NewOrderRepository()
	stockRepo, err := inmemory.NewStockRepository()
	if err != nil {
		return nil, err
	}

	// Initialize service
	service := loms.NewService(orderRepo, stockRepo)

	// Initialize gRPC server
	server := grpc.NewServer(cfg.GRPC.Port, service)

	return &App{
		config: cfg,
		server: server,
		loms:   lomsClient,
	}, nil
}

func (a *App) Run() error {
	return a.server.Start()
}

func (a *App) Stop() {
	if a.loms != nil {
		a.loms.Close()
	}
	a.server.Stop()
}
