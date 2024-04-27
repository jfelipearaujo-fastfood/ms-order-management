package order_repository

import (
	"context"
	"database/sql"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
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

func (r *OrderRepository) GetByID(ctx context.Context, id string) (entity.Order, error) {
	order := entity.Order{}

	queryGetOrder := `
		SELECT id, customer_id, track_id, state, state_updated_at, created_at, updated_at
		FROM orders
		WHERE id = $1;	
	`

	queryGetOrderItems := `
		SELECT product_id, quantity, price
		FROM order_items
		WHERE order_id = $1;
	`

	statement, err := r.conn.QueryContext(ctx, queryGetOrder, id)
	if err != nil {
		return entity.Order{}, err
	}
	defer statement.Close()

	for statement.Next() {
		err = statement.Scan(
			&order.UUID,
			&order.CustomerID,
			&order.TrackID,
			&order.State,
			&order.StateUpdatedAt,
			&order.CreatedAt,
			&order.UpdatedAt)
		if err != nil {
			return entity.Order{}, err
		}
	}

	if order.UUID == "" {
		return entity.Order{}, repository.ErrOrderNotFound
	}

	order.Items = []entity.Item{}

	statement, err = r.conn.QueryContext(ctx, queryGetOrderItems, id)
	if err != nil {
		return entity.Order{}, err
	}
	defer statement.Close()

	for statement.Next() {
		item := entity.Item{}
		err = statement.Scan(
			&item.UUID,
			&item.Quantity,
			&item.UnitPrice)
		if err != nil {
			return entity.Order{}, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}
