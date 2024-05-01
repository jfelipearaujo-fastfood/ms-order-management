package order_entity

type OrderState int

const (
	None       OrderState = iota
	Created               // When the order is created
	Received              // When the order is received and is ready to be processed
	Processing            // When the order is being processed by the kitchen
	Completed             // When the order is completed and ready to be delivered
	Delivered             // When the order is delivered to the customer
	Cancelled             // When the order is cancelled
)

var (
	order_state_machine = map[OrderState][]OrderState{
		None:       {Created},
		Created:    {Received, Cancelled},
		Received:   {Processing, Cancelled},
		Processing: {Completed, Cancelled},
		Completed:  {Delivered},
	}
)

func (s OrderState) CanTransitionTo(to OrderState) bool {
	for _, allowed := range order_state_machine[s] {
		if to == allowed {
			return true
		}
	}
	return false
}

func (s OrderState) String() string {
	switch s {
	case None:
		return "None"
	case Created:
		return "Created"
	case Received:
		return "Received"
	case Processing:
		return "Processing"
	case Completed:
		return "Completed"
	case Delivered:
		return "Delivered"
	case Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func IsValidState(s OrderState) bool {
	return s >= Created && s <= Cancelled
}
