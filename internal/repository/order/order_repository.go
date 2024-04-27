package order_repository

import (
	"context"
	"database/sql"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
)

type OrderRepository struct {
	conn *sql.DB
}

func NewOrderRepository(conn *sql.DB) *OrderRepository {
	return &OrderRepository{
		conn: conn,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *entity.Order) error {
	queryInsertOrder := `
		INSERT INTO orders (id, customer_id, track_id, state, state_updated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	queryInsertOrderItems := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4);
	`

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		queryInsertOrder,
		order.UUID,
		order.CustomerID,
		order.TrackID,
		order.State,
		order.StateUpdatedAt,
		order.CreatedAt,
		order.UpdatedAt)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx,
			queryInsertOrderItems,
			order.UUID,
			item.UUID,
			item.Quantity,
			item.UnitPrice)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return errTx
			}
			return err
		}
	}

	return tx.Commit()
}
