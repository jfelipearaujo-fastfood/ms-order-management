package payment_entity

import "time"

type Payment struct {
	OrderId   string `json:"order_id"`
	PaymentId string `json:"payment_id"`

	TotalItems int          `json:"total_items"`
	Amount     float64      `json:"amount"`
	State      PaymentState `json:"state"`
	StateTitle string       `json:"state_title"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPayment(orderId, paymentId string, totalItems int, amount float64, now time.Time) Payment {
	return Payment{
		OrderId:   orderId,
		PaymentId: paymentId,

		TotalItems: totalItems,
		Amount:     amount,
		State:      WaitingForApproval,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Payment) IsInState(states ...PaymentState) bool {
	for _, state := range states {
		if p.State == state {
			return true
		}
	}

	return false
}

func (p *Payment) RefreshStateTitle() {
	p.StateTitle = p.State.String()
}
