package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vagonaizer/loms/internal/domain/models"
)

var (
	ErrStockNotFound     = errors.New("stock not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) Get(ctx context.Context, sku uint32) (*models.Stock, error) {
	row := r.db.QueryRowContext(ctx, `SELECT sku, total_count, reserved FROM stock WHERE sku = $1`, sku)
	var stock models.Stock
	if err := row.Scan(&stock.SKU, &stock.TotalCount, &stock.Reserved); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStockNotFound
		}
		return nil, err
	}
	return &stock, nil
}

func (r *StockRepository) Update(ctx context.Context, stock *models.Stock) error {
	_, err := r.db.ExecContext(ctx, `UPDATE stock SET total_count = $1, reserved = $2 WHERE sku = $3`, stock.TotalCount, stock.Reserved, stock.SKU)
	return err
}

func (r *StockRepository) Reserve(ctx context.Context, sku uint32, count uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var available, reserved, total uint64
	err = tx.QueryRowContext(ctx, `SELECT total_count, reserved FROM stock WHERE sku = $1 FOR UPDATE`, sku).Scan(&total, &reserved)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStockNotFound
		}
		return err
	}
	available = total - reserved
	if available < count {
		return ErrInsufficientStock
	}
	_, err = tx.ExecContext(ctx, `UPDATE stock SET reserved = reserved + $1 WHERE sku = $2`, count, sku)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *StockRepository) Release(ctx context.Context, sku uint32, count uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var reserved uint64
	err = tx.QueryRowContext(ctx, `SELECT reserved FROM stock WHERE sku = $1 FOR UPDATE`, sku).Scan(&reserved)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStockNotFound
		}
		return err
	}
	var toRelease uint64
	if reserved < count {
		toRelease = reserved
	} else {
		toRelease = count
	}
	_, err = tx.ExecContext(ctx, `UPDATE stock SET reserved = reserved - $1 WHERE sku = $2`, toRelease, sku)
	if err != nil {
		return err
	}
	return tx.Commit()
}
