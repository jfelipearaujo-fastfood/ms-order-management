package payment

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
)

type PaymentRepository struct {
	conn *sql.DB
}

func NewPaymentRepository(conn *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		conn: conn,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *payment_entity.Payment) error {
	sql, params, err := goqu.
		Insert("order_payments").
		Cols("order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at").
		Vals(
			goqu.Vals{
				payment.OrderId,
				payment.PaymentId,
				payment.TotalItems,
				payment.Amount,
				payment.State,
				payment.CreatedAt,
				payment.UpdatedAt,
			},
		).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = r.conn.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *payment_entity.Payment) error {
	sql, params, err := goqu.
		Update("order_payments").
		Set(goqu.Record{
			"state":      payment.State,
			"updated_at": payment.UpdatedAt,
		}).
		Where(goqu.Ex{"payment_id": payment.PaymentId}).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = r.conn.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}
