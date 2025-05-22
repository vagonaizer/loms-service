package inmemory

import (
	"context"
	"sync"

	"github.com/vagonaizer/loms/internal/domain/models"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[int64]*models.Order
	nextID int64
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int64]*models.Order),
		nextID: 1,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = r.nextID
	r.nextID++
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepository) Get(ctx context.Context, orderID int64) (*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[orderID]
	if !exists {
		return nil, nil
	}
	return order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return nil
	}
	r.orders[order.ID] = order
	return nil
}
