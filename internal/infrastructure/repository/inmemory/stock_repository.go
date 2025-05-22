package inmemory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/vagonaizer/loms/internal/domain/models"
)

var (
	ErrStockNotFound     = errors.New("stock not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type stockData struct {
	SKU        uint32 `json:"sku"`
	TotalCount uint64 `json:"total_count"`
	Reserved   uint64 `json:"reserved"`
}

type StockRepository struct {
	mu     sync.RWMutex
	stocks map[uint32]*models.Stock
}

func NewStockRepository() (*StockRepository, error) {
	repo := &StockRepository{
		stocks: make(map[uint32]*models.Stock),
	}

	// Load initial data
	data, err := os.ReadFile("internal/infrastructure/repository/inmemory/stock_data.json")
	if err != nil {
		return nil, err
	}

	var stocks []stockData
	if err := json.Unmarshal(data, &stocks); err != nil {
		return nil, err
	}

	for _, s := range stocks {
		repo.stocks[s.SKU] = &models.Stock{
			SKU:        s.SKU,
			TotalCount: s.TotalCount,
			Reserved:   s.Reserved,
		}
	}

	return repo, nil
}

func (r *StockRepository) Get(ctx context.Context, sku uint32) (*models.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, exists := r.stocks[sku]
	if !exists {
		return nil, ErrStockNotFound
	}

	// Return a copy with actual available count
	return &models.Stock{
		SKU:        stock.SKU,
		TotalCount: stock.TotalCount - stock.Reserved, // Return only available count
		Reserved:   stock.Reserved,
	}, nil
}

func (r *StockRepository) Update(ctx context.Context, stock *models.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stocks[stock.SKU] = stock
	return nil
}

func (r *StockRepository) Reserve(ctx context.Context, sku uint32, count uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, exists := r.stocks[sku]
	if !exists {
		return ErrStockNotFound
	}

	if stock.TotalCount-stock.Reserved < count {
		return ErrInsufficientStock
	}

	// Create a copy of the stock to avoid modifying the original
	updatedStock := &models.Stock{
		SKU:        stock.SKU,
		TotalCount: stock.TotalCount,
		Reserved:   stock.Reserved + count,
	}

	// Save the updated stock
	r.stocks[sku] = updatedStock
	return nil
}

func (r *StockRepository) Release(ctx context.Context, sku uint32, count uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stock, exists := r.stocks[sku]
	if !exists {
		return ErrStockNotFound
	}

	if stock.Reserved < count {
		stock.Reserved = 0
		return nil
	}

	// Create a copy of the stock to avoid modifying the original
	updatedStock := &models.Stock{
		SKU:        stock.SKU,
		TotalCount: stock.TotalCount,
		Reserved:   stock.Reserved - count,
	}

	// Save the updated stock
	r.stocks[sku] = updatedStock
	return nil
}
