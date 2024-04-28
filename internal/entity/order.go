package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type Order struct {
	UUID string `json:"id"`

	CustomerID     string    `json:"customer_id"`
	TrackID        TrackID   `json:"track_id"`
	State          State     `json:"state"`
	StateTitle     string    `json:"state_title"`
	StateUpdatedAt time.Time `json:"state_updated_at"`

	TotalItems int     `json:"total_items"`
	TotalPrice float64 `json:"total_price"`

	Items []Item `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrder(customerID string, now time.Time) Order {
	return Order{
		UUID: uuid.NewString(),

		CustomerID:     customerID,
		TrackID:        NewTrackID(),
		State:          Created,
		StateUpdatedAt: now,

		TotalItems: 0,
		TotalPrice: 0,

		Items: make([]Item, 0),

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (o *Order) AddItem(item Item, now time.Time) error {
	for _, i := range o.Items {
		if i.UUID == item.UUID {
			return custom_error.ErrOrderItemAlreadyExists
		}
	}

	o.Items = append(o.Items, item)
	o.UpdatedAt = now

	o.calculateTotals()

	return nil
}

func (o *Order) calculateTotals() {
	o.TotalItems = 0
	o.TotalPrice = 0

	for _, item := range o.Items {
		o.TotalItems += item.Quantity
		o.TotalPrice += item.UnitPrice * float64(item.Quantity)
	}
}

func (o *Order) UpdateState(toState State, now time.Time) error {
	if o.State == toState {
		return nil
	}

	if !o.State.CanTransitionTo(toState) {
		return custom_error.ErrOrderInvalidStateTransition
	}

	o.State = toState
	o.StateTitle = toState.String()
	o.StateUpdatedAt = now
	o.UpdatedAt = now

	return nil
}

func (o *Order) RefreshStateTitle() {
	o.StateTitle = o.State.String()
}

func (o *Order) CanAddItems() bool {
	return o.State == Created
}

func (o *Order) IsCompleted() bool {
	return o.State == Delivered || o.State == Cancelled
}
