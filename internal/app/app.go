package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vagonaizer/loms/internal/config"
	"github.com/vagonaizer/loms/internal/infrastructure/api/grpc"
	lomsclient "github.com/vagonaizer/loms/internal/infrastructure/client/loms"
	"github.com/vagonaizer/loms/internal/infrastructure/repository/postgres"
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

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Initialize repositories
	orderRepo := postgres.NewOrderRepository(db)
	stockRepo := postgres.NewStockRepository(db)

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
