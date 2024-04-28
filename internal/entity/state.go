package entity

type State int

const (
	None State = iota
	Created
	Received
	Processing
	Completed
	Delivered
	Cancelled
)

var (
	state_machine = map[State][]State{
		None:       {Created},
		Created:    {Received, Cancelled},
		Received:   {Processing, Cancelled},
		Processing: {Completed, Cancelled},
		Completed:  {Delivered},
	}
)

func (s State) CanTransitionTo(to State) bool {
	for _, allowed := range state_machine[s] {
		if to == allowed {
			return true
		}
	}
	return false
}

func (s State) String() string {
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

func IsValidState(s State) bool {
	return s >= Created && s <= Cancelled
}
