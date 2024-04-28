package order_repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
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
	return r.getBy(ctx, "id", id)
}

func (r *OrderRepository) GetByTrackID(ctx context.Context, trackId string) (entity.Order, error) {
	return r.getBy(ctx, "track_id", trackId)
}

func (r *OrderRepository) getBy(ctx context.Context, column string, value string) (entity.Order, error) {
	order := entity.Order{}

	sql, params, err := goqu.
		From("orders").
		Select("id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at").
		Where(goqu.Ex{column: value}).
		ToSQL()
	if err != nil {
		return entity.Order{}, err
	}

	queryGetOrderItems := `
		SELECT oi.product_id, oi.quantity, oi.price
		FROM order_items oi
		LEFT JOIN orders o ON o.id = oi.order_id
		WHERE oi.order_id = $1 OR o.track_id = $1;
	`

	statement, err := r.conn.QueryContext(ctx, sql, params...)
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
		return entity.Order{}, custom_error.ErrOrderNotFound
	}

	order.Items = []entity.Item{}

	statement, err = r.conn.QueryContext(ctx, queryGetOrderItems, value)
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

func (r *OrderRepository) GetAll(
	ctx context.Context,
	pagination common.Pagination,
	filter repository.GetAllOrdersFilter,
) (int, []entity.Order, error) {
	skip := pagination.Page*pagination.Size - pagination.Size

	stateFilter := goqu.And(
		goqu.Ex{"customer_id": filter.CustomerID},
		goqu.Ex{"state": goqu.Op{
			"gte": filter.StateFrom,
		}},
		goqu.Ex{"state": goqu.Op{
			"lt": filter.StateTo,
		}},
	)

	sql, params, err := goqu.
		From("orders").
		Select(goqu.COUNT("id")).
		Where(goqu.And(stateFilter)).
		Limit(uint(pagination.Size)).
		Offset(uint(skip)).
		ToSQL()

	if err != nil {
		return 0, nil, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return 0, nil, err
	}

	var count int

	for statement.Next() {
		err = statement.Scan(&count)
		if err != nil {
			return 0, nil, err
		}
	}

	sql, params, err = goqu.
		From("orders").
		Select("id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at").
		Where(stateFilter).
		Order(goqu.I("created_at").Asc()).
		Limit(uint(pagination.Size)).
		Offset(uint(skip)).
		ToSQL()
	if err != nil {
		return 0, nil, err
	}

	statement, err = r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return 0, nil, err
	}

	orders := []entity.Order{}

	for statement.Next() {
		order := entity.Order{}
		err = statement.Scan(
			&order.UUID,
			&order.CustomerID,
			&order.TrackID,
			&order.State,
			&order.StateUpdatedAt,
			&order.CreatedAt,
			&order.UpdatedAt)
		if err != nil {
			return 0, nil, err
		}
		orders = append(orders, order)
	}

	return count, orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *entity.Order, updateItems bool) error {
	queryUpdateOrder := `
		UPDATE orders
		SET state = $1, state_updated_at = $2, updated_at = $3
		WHERE id = $4;
	`

	queryDeleteOrderItems := `
		DELETE FROM order_items
		WHERE order_id = $1;
	`

	queryInsertOrderItems := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4);
	`

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx,
		queryUpdateOrder,
		order.State,
		order.StateUpdatedAt,
		order.UpdatedAt,
		order.UUID)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return err
	}

	if rowsAffected == 0 {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return custom_error.ErrOrderNotFound
	}

	if updateItems {
		_, err = tx.ExecContext(ctx,
			queryDeleteOrderItems,
			order.UUID)
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
	}

	return tx.Commit()
}
