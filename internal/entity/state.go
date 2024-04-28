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

func IsValidState(s State) bool {
	return s >= Created && s <= Cancelled
}
