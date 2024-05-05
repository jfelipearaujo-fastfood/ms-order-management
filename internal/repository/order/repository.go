package order_repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
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

func (r *OrderRepository) Create(ctx context.Context, order *order_entity.Order) error {
	queryInsertOrder := `
		INSERT INTO orders (id, customer_id, track_id, state, state_updated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	queryInsertOrderItems := `
		INSERT INTO order_items (order_id, product_id, name, quantity, price)
		VALUES ($1, $2, $3, $4, $5);
	`

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		queryInsertOrder,
		order.Id,
		order.CustomerId,
		order.TrackId,
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
			order.Id,
			item.Id,
			item.Name,
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

func (r *OrderRepository) GetByID(ctx context.Context, id string) (order_entity.Order, error) {
	return r.getBy(ctx, "id", id)
}

func (r *OrderRepository) GetByTrackID(ctx context.Context, trackId string) (order_entity.Order, error) {
	return r.getBy(ctx, "track_id", trackId)
}

func (r *OrderRepository) getBy(ctx context.Context, column string, value string) (order_entity.Order, error) {
	order := order_entity.Order{}

	sql, params, err := goqu.
		From("orders").
		Select("id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at").
		Where(goqu.Ex{column: value}).
		ToSQL()
	if err != nil {
		return order_entity.Order{}, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return order_entity.Order{}, err
	}
	defer statement.Close()

	for statement.Next() {
		err = statement.Scan(
			&order.Id,
			&order.CustomerId,
			&order.TrackId,
			&order.State,
			&order.StateUpdatedAt,
			&order.CreatedAt,
			&order.UpdatedAt)
		if err != nil {
			return order_entity.Order{}, err
		}
	}

	if order.Id == "" {
		return order_entity.Order{}, custom_error.ErrOrderNotFound
	}

	sql, params, err = goqu.
		From("order_payments").
		Select("order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at").
		Where(goqu.Ex{"order_id": order.Id}).
		Order(goqu.I("created_at").Asc()).
		ToSQL()
	if err != nil {
		return order_entity.Order{}, err
	}

	statement, err = r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return order_entity.Order{}, err
	}
	defer statement.Close()

	for statement.Next() {
		payment := payment_entity.Payment{}
		err = statement.Scan(
			&payment.OrderId,
			&payment.PaymentId,
			&payment.TotalItems,
			&payment.Amount,
			&payment.State,
			&payment.CreatedAt,
			&payment.UpdatedAt)
		if err != nil {
			return order_entity.Order{}, err
		}

		payment.RefreshStateTitle()

		order.Payments = append(order.Payments, payment)
	}

	order.Items = []order_entity.Item{}

	sql, params, err = goqu.
		From("order_items").
		Select("order_items.product_id", "order_items.name", "order_items.quantity", "order_items.price").
		LeftJoin(goqu.T("orders"), goqu.On(goqu.I("order_items.order_id").Eq(goqu.I("orders.id")))).
		Where(goqu.ExOr{
			"order_items.order_id": value,
			"orders.track_id":      value,
		}).
		ToSQL()
	if err != nil {
		return order_entity.Order{}, err
	}

	statement, err = r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return order_entity.Order{}, err
	}
	defer statement.Close()

	for statement.Next() {
		item := order_entity.Item{}
		err = statement.Scan(
			&item.Id,
			&item.Name,
			&item.Quantity,
			&item.UnitPrice)
		if err != nil {
			return order_entity.Order{}, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderRepository) GetAll(
	ctx context.Context,
	pagination common.Pagination,
	filter repository.GetAllOrdersFilter,
) (int, []order_entity.Order, error) {
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

	orders := []order_entity.Order{}

	for statement.Next() {
		order := order_entity.Order{}
		err = statement.Scan(
			&order.Id,
			&order.CustomerId,
			&order.TrackId,
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

func (r *OrderRepository) Update(ctx context.Context, order *order_entity.Order, updateItems bool) error {
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
		INSERT INTO order_items (order_id, product_id, name, quantity, price)
		VALUES ($1, $2, $3, $4, $5);
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
		order.Id)
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
			order.Id)
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
				order.Id,
				item.Id,
				item.Name,
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
