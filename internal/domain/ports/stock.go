package ports

import (
	"context"

	"github.com/vagonaizer/loms/internal/domain/models"
)

type StockRepository interface {
	Get(ctx context.Context, sku uint32) (*models.Stock, error)
	Update(ctx context.Context, stock *models.Stock) error
	Reserve(ctx context.Context, sku uint32, count uint64) error
	Release(ctx context.Context, sku uint32, count uint64) error
}
