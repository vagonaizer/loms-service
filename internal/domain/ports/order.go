package ports

import (
	"context"

	"github.com/vagonaizer/loms/internal/domain/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	Get(ctx context.Context, orderID int64) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
}
