package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vagonaizer/loms/internal/domain/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var orderID int64
	err = tx.QueryRowContext(ctx, `INSERT INTO orders (user_id, status, created_at) VALUES ($1, $2, NOW()) RETURNING id`, order.UserID, order.Status).Scan(&orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, `INSERT INTO order_items (order_id, sku, count) VALUES ($1, $2, $3)`, orderID, item.SKU, item.Count)
		if err != nil {
			return err
		}
	}

	order.ID = orderID
	return tx.Commit()
}

func (r *OrderRepository) Get(ctx context.Context, orderID int64) (*models.Order, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, status, created_at FROM orders WHERE id = $1`, orderID)
	var order models.Order
	if err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `SELECT sku, count FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.SKU, &item.Count); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET status = $1 WHERE id = $2`, order.Status, order.ID)
	return err
}
